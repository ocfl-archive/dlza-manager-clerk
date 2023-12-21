package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"emperror.dev/errors"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/tbaehler/gin-keycloak/pkg/ginkeycloak"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/constants"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/graph"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/models"
	pb "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/proto"
	ubLogger "gitlab.switch.ch/ub-unibas/go-ublogger"
	"golang.org/x/net/http2"
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

	var tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{srv.cert},
		RootCAs:      rootCAs,
	}
	var keycloakConfig = ginkeycloak.KeycloakConfig{
		Url:   srv.keycloak.Addr,
		Realm: srv.keycloak.Realm,
	}

	// router := gin.Default()
	router := srv.router
	router.GET("/schema", func(ctx *gin.Context) {
		ctx.FileFromFS("graph/schema.graphqls", http.FS(SchemaFS))
	})
	router.GET("/graphql", func(ctx *gin.Context) {
		ctx.FileFromFS("dlza-frontend/build/playground.html", http.FS(UiFS))
	})
	graphql := router.Group("/graphql").Use(ginkeycloak.Auth(ginkeycloak.AuthCheck(), keycloakConfig))
	{
		graphql.POST("", srv.graphqlHandler(srv.ClientClerkHandler))
	}
	router.Use(func(ctx *gin.Context) {
		fsys, _ := fs.Sub(UiFS, "dlza-frontend/build")
		path := ctx.Request.URL.Path

		ctx.FileFromFS(path, http.FS(fsys))
	})

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

		rawAccessToken := c.Request.Header.Get("Authorization")
		parts := strings.Split(rawAccessToken, " ")
		if len(parts) != 2 {
			c.Writer.WriteHeader(400)
			fmt.Println("error 400", len(parts))
			return
		}

		var userClaim models.KeyCloakToken
		provider, err := oidc.NewProvider(c, srv.keycloak.Addr+"realms/"+srv.keycloak.Realm)
		if err != nil {
			panic(err)
		}

		oidcConfig := &oidc.Config{
			ClientID: srv.keycloak.ClientId,
		}

		verifier := provider.Verifier(oidcConfig)

		_, err = verifier.Verify(c, parts[1])
		if err != nil {
			c.Error(errors.Errorf("Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError))
			return
		}

		_, err = jwt.ParseWithClaims(parts[1], &userClaim, nil)
		if err != nil && err.Error() != "no Keyfunc was provided." {
			fmt.Println("error 400", err)
			c.Writer.WriteHeader(400)
		}
		ctx := context.WithValue(c, constants.Needed, "Needed to attach context")
		c.Set("keycloak_group", userClaim.Groups)
		c.Set("tenant_list", userClaim.TenantList)

		h.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
	}
}
