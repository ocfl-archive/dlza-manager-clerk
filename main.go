package main

import (
	"crypto/tls"
	"crypto/x509"
	"embed"
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"encoding/pem"
	"flag"
	"fmt"
	configutil "github.com/je4/utils/v2/pkg/config"
	"github.com/je4/utils/v2/pkg/zLogger"
	"github.com/ocfl-archive/dlza-manager-clerk/certs"
	"github.com/ocfl-archive/dlza-manager-clerk/config"
	"github.com/ocfl-archive/dlza-manager-clerk/controller"
	"github.com/ocfl-archive/dlza-manager-clerk/data/web"
	"github.com/ocfl-archive/dlza-manager-clerk/models"
	"github.com/ocfl-archive/dlza-manager-clerk/router"
	graphqlServer "github.com/ocfl-archive/dlza-manager-clerk/server"
	handlerClientProto "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	ublogger "gitlab.switch.ch/ub-unibas/go-ublogger/v2"
	"go.ub.unibas.ch/cloud/certloader/v2/pkg/loader"
	"go.ub.unibas.ch/cloud/miniresolver/v2/pkg/resolver"
	"io"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var configFile = flag.String("config", "", "config file in toml format")

//go:embed all:dlza-frontend/build
var uiFS embed.FS

//go:embed graph/schema.graphqls
var schemaFS embed.FS

func main() {

	flag.Parse()

	var cfgFS fs.FS
	var cfgFile string
	if *configFile != "" {
		cfgFS = os.DirFS(filepath.Dir(*configFile))
		cfgFile = filepath.Base(*configFile)
	} else {
		cfgFS = config.ConfigFS
		cfgFile = "clerk.toml"
	}

	conf := &config.Config{
		LocalAddr: "localhost:8443",
		//ResolverTimeout: config.Duration(10 * time.Minute),
		ExternalAddr:            "https://localhost:8443",
		ResolverTimeout:         configutil.Duration(10 * time.Minute),
		ResolverNotFoundTimeout: configutil.Duration(10 * time.Second),
		ServerTLS: &loader.Config{
			Type: "DEV",
		},
		ClientTLS: &loader.Config{
			Type: "DEV",
		},
	}
	if err := config.LoadConfig(cfgFS, cfgFile, conf); err != nil {
		log.Fatalf("cannot load toml from [%v] %s: %v", cfgFS, cfgFile, err)
	}
	// create logger instance
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("cannot get hostname: %v", err)
	}

	var loggerTLSConfig *tls.Config
	var loggerLoader io.Closer
	if conf.Log.Stash.TLS != nil {
		loggerTLSConfig, loggerLoader, err = loader.CreateClientLoader(conf.Log.Stash.TLS, nil)
		if err != nil {
			log.Fatalf("cannot create client loader: %v", err)
		}
		defer loggerLoader.Close()
	}

	_logger, _logstash, _logfile, err := ublogger.CreateUbMultiLoggerTLS(conf.Log.Level, conf.Log.File,
		ublogger.SetDataset(conf.Log.Stash.Dataset),
		ublogger.SetLogStash(conf.Log.Stash.LogstashHost, conf.Log.Stash.LogstashPort, conf.Log.Stash.Namespace, conf.Log.Stash.LogstashTraceLevel),
		ublogger.SetTLS(conf.Log.Stash.TLS != nil),
		ublogger.SetTLSConfig(loggerTLSConfig),
	)
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}
	if _logstash != nil {
		defer _logstash.Close()
	}

	if _logfile != nil {
		defer _logfile.Close()
	}

	l2 := _logger.With().Timestamp().Str("host", hostname).Logger() //.Output(output)
	var logger zLogger.ZLogger = &l2

	clientCert, clientLoader, err := loader.CreateClientLoader(conf.ClientTLS, logger)
	if err != nil {
		logger.Panic().Msgf("cannot create client loader: %v", err)
	}
	defer clientLoader.Close()

	logger.Info().Msgf("resolver address is %s", conf.ResolverAddr)
	resolverClient, err := resolver.NewMiniresolverClient(conf.ResolverAddr, conf.GRPCClient, clientCert, nil, time.Duration(conf.ResolverTimeout), time.Duration(conf.ResolverNotFoundTimeout), logger)
	if err != nil {
		logger.Fatal().Msgf("cannot create resolver client: %v", err)
	}
	defer resolverClient.Close()

	//////ClerkHandler gRPC connection

	clientClerkHandler, err := resolver.NewClient[handlerClientProto.ClerkHandlerServiceClient](
		resolverClient,
		handlerClientProto.NewClerkHandlerServiceClient,
		handlerClientProto.ClerkHandlerService_ServiceDesc.ServiceName, conf.Domain)
	if err != nil {
		logger.Panic().Msgf("cannot create clientClerkHandler grpc client: %v", err)
	}

	tenantController := controller.NewTenantController(clientClerkHandler)
	storageLocationController := controller.NewStorageLocationController(clientClerkHandler)
	collectionController := controller.NewCollectionController(clientClerkHandler)
	statusController := controller.NewStatusController(clientClerkHandler)
	objectInstanceController := controller.NewObjectInstanceController(clientClerkHandler)
	objectController := controller.NewObjectController(clientClerkHandler)
	routes := router.NewRouter(conf.Jwt, tenantController, storageLocationController, collectionController, statusController, objectInstanceController, objectController)

	// find static fs
	var staticFS fs.FS
	if conf.GraphQLConfig.WebStatic == "" {
		staticFS, err = fs.Sub(web.WebFS, "ui/build")
		if err != nil {
			emperror.Panic(errors.Wrapf(err, "cannot create subFS of %v/%s", web.WebFS, "ui/build"))
		}
	} else {
		staticFS = os.DirFS(conf.GraphQLConfig.WebStatic)
	}

	var cert tls.Certificate
	var addCA = []*x509.Certificate{}
	if conf.GraphQLConfig.TLSCert == "" {
		certBytes, err := fs.ReadFile(certs.CertFS, "localhost.cert.pem")
		if err != nil {
			emperror.Panic(errors.Wrapf(err, "cannot read internal cert %v/%s", certs.CertFS, "localhost.cert.pem"))
		}
		keyBytes, err := fs.ReadFile(certs.CertFS, "localhost.key.pem")
		if err != nil {
			emperror.Panic(errors.Wrapf(err, "cannot read internal key %v/%s", certs.CertFS, "localhost.key.pem"))
		}
		if cert, err = tls.X509KeyPair(certBytes, keyBytes); err != nil {
			emperror.Panic(errors.Wrap(err, "cannot create internal cert"))
		}

		rootCABytes, err := fs.ReadFile(certs.CertFS, "ca.cert.pem")
		if err != nil {
			emperror.Panic(errors.Wrapf(err, "cannot read root ca %v/%s", certs.CertFS, "ca.cert.pem"))
		}
		block, _ := pem.Decode(rootCABytes)
		if block == nil {
			emperror.Panic(errors.Wrapf(err, "cannot decode root ca"))
		}
		rootCA, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			emperror.Panic(errors.Wrap(err, "cannot parse root ca"))
		}
		addCA = append(addCA, rootCA)
	} else {
		if cert, err = tls.LoadX509KeyPair(conf.GraphQLConfig.TLSCert, conf.GraphQLConfig.TLSKey); err != nil {
			emperror.Panic(errors.Wrapf(err, "cannot load key pair %s - %s", conf.GraphQLConfig.TLSCert, conf.GraphQLConfig.TLSKey))
		}
		if conf.GraphQLConfig.RootCA != nil {
			for _, caName := range conf.GraphQLConfig.RootCA {
				rootCABytes, err := os.ReadFile(caName)
				if err != nil {
					emperror.Panic(errors.Wrapf(err, "cannot read root ca %s", caName))
				}
				block, _ := pem.Decode(rootCABytes)
				if block == nil {
					emperror.Panic(errors.Wrapf(err, "cannot decode root ca"))
				}
				rootCA, err := x509.ParseCertificate(block.Bytes)
				if err != nil {
					emperror.Panic(errors.Wrap(err, "cannot parse root ca"))
				}
				addCA = append(addCA, rootCA)
			}
		}
	}
	graphqlServer.UiFS = uiFS
	graphqlServer.SchemaFS = schemaFS
	srv, err := graphqlServer.NewServer(conf.GraphQLConfig.Addr, conf.GraphQLConfig.ExtAddr, cert, addCA, staticFS, logger, models.Keycloak{
		Addr:         conf.GraphQLConfig.Keycloak.Addr,
		Realm:        conf.GraphQLConfig.Keycloak.Realm,
		Callback:     conf.GraphQLConfig.Keycloak.Callback,
		ClientId:     conf.GraphQLConfig.Keycloak.ClientId,
		ClientSecret: conf.GraphQLConfig.Keycloak.ClientSecret,
	}, clientClerkHandler, routes, conf.GraphQLConfig.Domain)
	if err != nil {
		emperror.Panic(errors.Wrap(err, "cannot create server"))
	}
	cancel, err := srv.Startup()
	if err != nil {
		emperror.Panic(errors.Wrap(err, "cannot start server"))
	}
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	fmt.Println("press ctrl+c to stop server")
	s := <-done
	fmt.Println("got signal:", s)

}
