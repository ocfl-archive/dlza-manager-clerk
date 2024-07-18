package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"embed"
	"encoding/gob"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"emperror.dev/errors"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/ocfl-archive/dlza-manager-clerk/constants"
	"github.com/ocfl-archive/dlza-manager-clerk/graph"
	"github.com/ocfl-archive/dlza-manager-clerk/middleware"
	"github.com/ocfl-archive/dlza-manager-clerk/models"
	pb "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	ubLogger "gitlab.switch.ch/ub-unibas/go-ublogger"
	"golang.org/x/net/http2"
)

func NewServer(addr, extAddr string, cert tls.Certificate, addCAs []*x509.Certificate, staticFS fs.FS, logger *ubLogger.Logger, keycloak models.Keycloak, clientClerkHandler pb.ClerkHandlerServiceClient, router *gin.Engine, domain string) (*Server, error) {
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
		domain:             domain,
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
	domain             string
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

	// // Keycloak configuration
	// ctx := context.Background()

	provider := middleware.GetProvider(srv.keycloak)
	var claims struct {
		EndSessionURL string `json:"end_session_endpoint"`
	}

	// oauth2Config := middleware.GetOauth2Config(srv.keycloak)

	err := provider.Claims(&claims)
	if err != nil {
		return nil, err
	}

	// oidcConfig := &oidc.Config{
	// 	ClientID: srv.keycloak.ClientId,
	// }
	// verifier := provider.Verifier(oidcConfig)

	var tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{srv.cert},
		RootCAs:      rootCAs,
	}

	router := srv.router

	store := memstore.NewStore([]byte("secretToBeChangedWithKubernetesSecret"))
	// store.Options(sessions.Options{Secure: true, SameSite: http.SameSiteLaxMode, HttpOnly: true})
	store.Options(sessions.Options{Secure: true, SameSite: http.SameSiteNoneMode, HttpOnly: true})
	gob.Register(models.KeyCloakToken{})
	router.Use(sessions.Sessions("mysession", store))
	// router.NoRoute(func(c *gin.Context) {
	// 	fmt.Println("test")
	// 	fmt.Printf("%s doesn't exists, redirect on / ", c.Request.URL.Path)
	// 	c.Redirect(http.StatusMovedPermanently, "/")
	// })
	router.Use(middleware.GinContextToContextMiddleware())
	router.GET("/playground", playgroundHandler())
	// router.GET("/auth/login", func(c *gin.Context) {
	// 	session := sessions.Default(c)
	// 	state := middleware.GenerateStateOauth()
	// 	if err != nil {
	// 		c.Error(errors.Errorf("Internal error:"+err.Error(), http.StatusInternalServerError))
	// 		return
	// 	}
	// 	nonce := middleware.GenerateStateOauth()
	// 	if err != nil {
	// 		c.Error(errors.Errorf("Invalid or malformed token:"+err.Error(), http.StatusInternalServerError))
	// 		return
	// 	}

	// 	session.Set("state", state)
	// 	session.Set("nonce", nonce)
	// 	session.Save()

	// 	output := map[string]any{
	// 		"auth_code_url": oauth2Config.AuthCodeURL(state, oidc.Nonce(nonce)),
	// 	}
	// 	c.JSON(http.StatusFound, output)
	// })
	// router.GET("/auth/callback", func(c *gin.Context) {
	// 	output := map[string]any{
	// 		"code": c.Request.URL.Query().Get("code"),
	// 	}
	// 	c.JSON(http.StatusFound, output)
	// })

	graphql := router.Group("/graphql")
	{
		graphql.POST("", srv.graphqlHandler(srv.ClientClerkHandler))
		graphql.OPTIONS("", srv.graphqlHandler(srv.ClientClerkHandler))
	}

	router.Use(static.Serve("/", static.EmbedFolder(UiFS, "dlza-frontend/build")))

	router.Use(func(ctx *gin.Context) {
		// if ctx.Request.URL.Path == "/playground" {
		// 	ctx.Next()
		// } else {
		fsys, _ := fs.Sub(UiFS, "dlza-frontend/build")
		if ctx.Request.URL.Path != "/" {
			fullPath := "dlza-frontend/build" + strings.TrimPrefix(path.Clean(ctx.Request.URL.Path), "/")
			_, err := os.Stat(fullPath)
			if err != nil {
				if !os.IsNotExist(err) {
					panic(err)
				}
				// Requested file does not exist so we return the default (resolves to index.html)
				ctx.Request.URL.Path = "/"
			}
		}
		path := ctx.Request.URL.Path

		ctx.FileFromFS(path, http.FS(fsys))
		// }

	})

	// router.GET("/playground", playgroundHandler())
	// router.GET("/schema", func(ctx *gin.Context) {
	// 	ctx.FileFromFS("graph/schema.graphqls", http.FS(SchemaFS))
	// })
	// .Use(middleware.VerifyToken(ctx, srv.keycloak, verifier, oauth2Config, srv.domain))
	// schema := router.Group("/schema").Use(middleware.VerifyToken(ctx, srv.keycloak, verifier, oauth2Config, srv.domain))
	// {
	// 	schema.GET("", func(ctx *gin.Context) {
	// 		ctx.FileFromFS("graph/schema.graphqls", http.FS(SchemaFS))
	// 	})
	// }
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

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Graphql handler
func (srv *Server) graphqlHandler(clientClerkHandler pb.ClerkHandlerServiceClient) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file

	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{ClientClerkHandler: clientClerkHandler}}))
	return func(c *gin.Context) {
		// fmt.Println("test before")
		// c.Header("Access-Control-Allow-Origin", "https://localhost:9087")
		// // c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		// // c.Header("Access-Control-Allow-Origin", "http://localhost:4173")
		// c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
		// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// fmt.Println("test after")
		ctx := context.WithValue(c, constants.Needed, "Needed to attach context")
		c.Set("keycloak", srv.keycloak)
		h.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
	}
}
