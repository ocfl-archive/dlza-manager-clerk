package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"slices"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/coreos/go-oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/constants"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/graph"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/middleware"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/models"
	pb "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/proto"
	ubLogger "gitlab.switch.ch/ub-unibas/go-ublogger"
	"golang.org/x/net/http2"
	"golang.org/x/oauth2"
)

func NewServer(addr, extAddr string, cert tls.Certificate, addCAs []*x509.Certificate, staticFS fs.FS, logger *ubLogger.Logger, keycloak models.Keycloak, clientClerkHandler pb.ClerkHandlerServiceClient, router *gin.Engine) (*Server, error) {
	server := &Server{
		addr:               addr,
		extAddr:            extAddr,
		cert:               cert,
		addCAs:             addCAs,
		staticFS:           staticFS,
		logger:             logger,
		keycloak:           keycloak,
		ClientClerkHandler: clientClerkHandler,
		router:             router,
	}
	return server, nil
}

type Server struct {
	extAddr            string
	server             http.Server
	staticFS           fs.FS
	addr               string
	cert               tls.Certificate
	addCAs             []*x509.Certificate
	logger             *ubLogger.Logger
	keycloak           models.Keycloak
	ClientClerkHandler pb.ClerkHandlerServiceClient
	router             *gin.Engine
}

var UiFS embed.FS
var SchemaFS embed.FS

func (srv *Server) Startup() (context.CancelFunc, error) {
	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	for _, ca := range srv.addCAs {
		rootCAs.AddCert(ca)
	}

	// Keycloak configuration
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, srv.keycloak.Addr+srv.keycloak.Realm)
	if err != nil {
		panic(err)
	}
	var claims struct {
		EndSessionURL string `json:"end_session_endpoint"`
	}
	err = provider.Claims(&claims)
	if err != nil {
		return nil, err
	}

	fmt.Println("srv.keycloak.ClientSecret", srv.keycloak.ClientSecret)
	oauth2Config := oauth2.Config{
		ClientID:     srv.keycloak.ClientId,
		ClientSecret: srv.keycloak.ClientSecret,
		RedirectURL:  srv.keycloak.Callback + "auth/callback",
		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	oidcConfig := &oidc.Config{
		ClientID: srv.keycloak.ClientId,
	}
	verifier := provider.Verifier(oidcConfig)

	var tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{srv.cert},
		RootCAs:      rootCAs,
	}
	// var keycloakConfig = ginkeycloak.KeycloakConfig{
	// 	Url:   srv.keycloak.Addr,
	// 	Realm: srv.keycloak.Realm,
	// }

	// router := gin.Default()
	router := srv.router
	store := cookie.NewStore([]byte(middleware.GenerateStateOauth()))
	router.Use(sessions.Sessions("mysession", store))

	router.GET("/logout", func(c *gin.Context) {
		c.Redirect(http.StatusFound, claims.EndSessionURL+"?user.id_token_hint="+srv.keycloak.ClientId+"&post_logout_redirect_uri="+srv.keycloak.Callback+"login")

	})
	router.GET("/auth/login", func(c *gin.Context) {
		state := middleware.GenerateStateOauth()
		if err != nil {
			c.Error(errors.Errorf("Internal error:"+err.Error(), http.StatusInternalServerError))
			return
		}
		nonce := middleware.GenerateStateOauth()
		if err != nil {
			c.Error(errors.Errorf("Invalid or malformed token:"+err.Error(), http.StatusInternalServerError))
			return
		}

		c.SetCookie("state", state, 60, "/", "localhost", false, true)
		c.SetCookie("nonce", nonce, 60, "/", "localhost", false, true)
		c.Redirect(http.StatusFound, oauth2Config.AuthCodeURL(state, oidc.Nonce(nonce)))
	})
	router.GET("/auth/callback", func(c *gin.Context) {
		session := sessions.Default(c)
		state, err := c.Cookie("state")
		if err != nil {
			c.Error(errors.Errorf("state not found:"+err.Error(), http.StatusBadRequest))
			c.Redirect(http.StatusFound, "/auth/login")
			return
		}
		if c.Request.URL.Query().Get("state") != state {
			c.Error(errors.Errorf("state did not match : "+err.Error(), http.StatusBadRequest))
			return
		}

		oauth2Token, err := oauth2Config.Exchange(ctx, c.Request.URL.Query().Get("code"))
		if err != nil {
			c.Error(errors.Errorf("Failed to exchange token:"+err.Error(), http.StatusInternalServerError))
			c.Redirect(http.StatusFound, "/auth/login")
			return
		}
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			c.Error(errors.Errorf("No id_token field in oauth2 token:", http.StatusInternalServerError))
			c.Redirect(http.StatusFound, "/auth/login")
			return
		}
		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			c.Error(errors.Errorf("callback Failed to verify ID Token:"+err.Error(), http.StatusInternalServerError))
			c.Redirect(http.StatusFound, "/auth/login")
			return
		}

		nonce, err := c.Cookie("nonce")
		if err != nil {
			c.Error(errors.Errorf("nonce not found:"+err.Error(), http.StatusBadRequest))
			c.Redirect(http.StatusFound, "/auth/login")
			return
		}
		if idToken.Nonce != nonce {
			c.Redirect(http.StatusFound, "/auth/login")
			c.Error(errors.Errorf("nonce did not match.", http.StatusBadRequest))
			return
		}
		resp := struct {
			OAuth2Token   *oauth2.Token
			IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
		}{oauth2Token, new(json.RawMessage)}

		if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
			c.Error(errors.Errorf(err.Error(), http.StatusInternalServerError))
			return
		}

		c.SetCookie("access_token", resp.OAuth2Token.AccessToken, int(time.Until(resp.OAuth2Token.Expiry).Seconds()), "/", "localhost", false, true)
		session.Set("refresh_token", resp.OAuth2Token.RefreshToken)
		session.Save()

		c.Redirect(http.StatusFound, "/")
	})
	graphql := router.Group("/graphql").Use(middleware.VerifyToken(ctx, srv.keycloak, verifier, oauth2Config, srv.keycloak))
	{
		graphql.POST("", srv.graphqlHandler(srv.ClientClerkHandler))
	}

	router.Use(middleware.VerifyToken(ctx, srv.keycloak, verifier, oauth2Config, srv.keycloak)).Use(func(ctx *gin.Context) {
		fsys, _ := fs.Sub(UiFS, "dlza-frontend/build")
		path := ctx.Request.URL.Path
		if slices.Contains([]string{"/collections", "/tenants", "/objects", "/files"}, path) {
			ctx.Redirect(http.StatusMovedPermanently, "/")
		}
		ctx.FileFromFS(path, http.FS(fsys))
	})

	router.GET("/schema", func(ctx *gin.Context) {
		ctx.FileFromFS("graph/schema.graphqls", http.FS(SchemaFS))
	}).Use(middleware.VerifyToken(ctx, srv.keycloak, verifier, oauth2Config, srv.keycloak))

	srv.server = http.Server{
		Addr:      srv.addr,
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	if err := http2.ConfigureServer(&srv.server, nil); err != nil {
		return nil, errors.Wrap(err, "cannot configure http2 server")
	}

	go func() {
		srv.logger.Info().Msgf("Starting server (%s): %s", srv.addr, srv.extAddr)
		if err := srv.server.ListenAndServeTLS("", ""); err != nil {
			srv.logger.Error().Msgf("server stopped: %v", err)
		} else {
			srv.logger.Info().Msg("server shut down")
		}
	}()
	return func() {
		if err := srv.server.Close(); err != nil {
			srv.logger.Error().Msgf("error closing server: %v", err)
		}
	}, nil
}

// Defining the Graphql handler
func (srv *Server) graphqlHandler(clientClerkHandler pb.ClerkHandlerServiceClient) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{ClientClerkHandler: clientClerkHandler}}))
	return func(c *gin.Context) {
		rawAccessToken, errT := c.Cookie("access_token")
		if errT != nil {
			c.Writer.WriteHeader(400)
			return
		}

		parts := strings.Split(rawAccessToken, " ")

		var userClaim models.KeyCloakToken
		provider, err := oidc.NewProvider(c, srv.keycloak.Addr+srv.keycloak.Realm)
		if err != nil {
			panic(err)
		}

		oidcConfig := &oidc.Config{
			ClientID: srv.keycloak.ClientId,
		}

		verifier := provider.Verifier(oidcConfig)

		_, err = verifier.Verify(c, parts[0])
		if err != nil {
			c.Error(errors.Errorf("Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError))
			return
		}

		_, err = jwt.ParseWithClaims(parts[0], &userClaim, nil)
		if err != nil && err.Error() != "no Keyfunc was provided." {
			c.Writer.WriteHeader(400)
		}
		ctx := context.WithValue(c, constants.Needed, "Needed to attach context")
		c.Set("keycloak_group", userClaim.Groups)
		c.Set("tenant_list", userClaim.TenantList)
		h.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
	}
}
