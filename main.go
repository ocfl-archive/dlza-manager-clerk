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
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/config"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/controller"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/data/certs"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/data/web"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/models"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/router"
	graphqlServer "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/server"
	handlerClient "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-handler/client"
	storageHandlerClient "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-storage-handler/client"
	ubLogger "gitlab.switch.ch/ub-unibas/go-ublogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var configParam = flag.String("config", "", "config file in toml format, no need for filetype for this param")

//go:embed all:dlza-frontend/build
var uiFS embed.FS

//go:embed graph/schema.graphqls
var schemaFS embed.FS

func main() {

	flag.Parse()
	conf, err := config.GetConfig(*configParam)
	if err != nil {
		log.Fatal(err)
	}

	//////ClerkStorageHandler gRPC connection
	clerkStorageHandlerServiceClient, connectionClerkStorageHandler, err := storageHandlerClient.NewStorageHandlerClient(conf.StorageHandler.Host+":"+strconv.Itoa(conf.StorageHandler.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer connectionClerkStorageHandler.Close()

	//////ClerkHandler gRPC connection
	clerkHandlerServiceClient, connectionClerkHandler, err := handlerClient.NewHandlerClient(conf.Handler.Host+":"+strconv.Itoa(conf.Handler.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer connectionClerkHandler.Close()

	tenantController := controller.NewTenantController(clerkHandlerServiceClient)
	storageLocationController := controller.NewStorageLocationController(clerkHandlerServiceClient)
	storagePartitionController := controller.NewStoragePartitionController(clerkStorageHandlerServiceClient)
	collectionController := controller.NewCollectionController(clerkHandlerServiceClient)
	statusController := controller.NewStatusController(clerkHandlerServiceClient)
	routes := router.NewRouter(tenantController, storageLocationController, collectionController, storagePartitionController, statusController)

	logger, logStash, logFile := ubLogger.CreateUbMultiLoggerTLS(
		conf.GraphQLConfig.Logging.TraceLevel, conf.GraphQLConfig.Logging.Filename,
		ubLogger.SetLogStash(conf.GraphQLConfig.Logging.StashHost, conf.GraphQLConfig.Logging.StashPortNb, conf.GraphQLConfig.Logging.Namespace, conf.GraphQLConfig.Logging.StashTraceLevel))
	if logStash != nil {
		defer logStash.Close()
	}
	if logFile != nil {
		defer logFile.Close()
	}

	// logger, logStash, logFile := ubLogger.CreateUbMultiLogger(
	// 	cfg.Logging.StashHost,
	// 	cfg.Logging.StashPortNb,
	// 	cfg.Logging.TraceLevel,
	// 	cfg.Logging.StashTraceLevel,
	// 	cfg.Logging.Filename,
	// 	cfg.Logging.Filename)
	// if logStash != nil {
	// 	defer logStash.Close()
	// }
	// if logFile != nil {
	// 	defer logFile.Close()
	// }

	// output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	// logger := zerolog.New(output).With().Str("timestamp", time.Now().String()).Logger()

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
	}, clerkHandlerServiceClient, routes, conf.GraphQLConfig.Domain)
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
