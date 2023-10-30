package main

import (
	"crypto/tls"
	"crypto/x509"
	"embed"
	"encoding/pem"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/config"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/controller"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/data/certs"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/data/web"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/models"
	pb "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/proto"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/router"
	graphqlServer "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/server"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/server/configstruct"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/service"

	"emperror.dev/emperror"
	"emperror.dev/errors"
	ubLogger "gitlab.switch.ch/ub-unibas/go-ublogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	portIngest  = ":8082"
	portHandler = ":8083"
)

var configParam = flag.String("config", "", "config file in toml format")

//go:embed all:ui/build
var uiFS embed.FS

func main() {

	// //////ClerkIngest gRPC connection
	// connectionClerkIngest, err := grpc.Dial("localhost"+portIngest, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }

	// //////ClerkHandler gRPC connection
	// connectionClerkHandler, err := grpc.Dial("localhost"+portHandler, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }

	//////ClerkIngest gRPC connection
	// connectionClerkIngest, err := grpc.Dial("dlza-manager-ingester-svc:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }

	// //////ClerkHandler gRPC connection
	// connectionClerkHandler, err := grpc.Dial("dlza-manager-handler-svc:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }

	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	connectionClerkIngest, err := grpc.Dial(conf.Handler, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	//////ClerkHandler gRPC connection
	connectionClerkHandler, err := grpc.Dial(conf.Ingester, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	clerkIngestServiceClient := pb.NewClerkIngestServiceClient(connectionClerkIngest)
	orderService := service.NewOrderService(clerkIngestServiceClient)

	clerkHandlerServiceClient := pb.NewClerkHandlerServiceClient(connectionClerkHandler)

	orderController := controller.NewOrderController(orderService)
	tenantController := controller.NewTenantController(clerkHandlerServiceClient)
	storageLocationController := controller.NewStorageLocationController(clerkHandlerServiceClient)
	storagePartitionController := controller.NewStoragePartitionController(clerkIngestServiceClient)
	collectionController := controller.NewCollectionController(clerkHandlerServiceClient)
	routes := router.NewRouter(orderController, tenantController, storageLocationController, collectionController, storagePartitionController)

	server := &http.Server{
		Addr:    "localhost:8086",
		Handler: routes,
	}
	go func() {
		err = server.ListenAndServe()

		if err != nil {
			log.Fatalf("error: %s", err.Error())
		}
	}()

	//////GraphQl
	flag.Parse()

	var cfgData []byte
	// read embedded toml if no external one is given
	if *configParam == "" {
		cfgData, err = fs.ReadFile(config.ConfigFS, "server.toml")
		if err != nil {
			emperror.Panic(errors.Wrapf(err, "cannot read %v/%s", config.ConfigFS, "server.toml"))
		}
	} else {
		cfgData, err = os.ReadFile(*configParam)
		if err != nil {
			emperror.Panic(errors.Wrapf(err, "cannot read %s", *configParam))
		}
	}
	// load config file
	cfg, err := configstruct.LoadGraphQLConfig(cfgData)
	if err != nil {
		emperror.Panic(errors.Wrap(err, "cannot decode config file"))
	}

	logger, logStash, logFile := ubLogger.CreateUbMultiLogger(
		cfg.Logging.StashHost,
		cfg.Logging.StashPortNb,
		cfg.Logging.TraceLevel,
		cfg.Logging.StashTraceLevel,
		cfg.Logging.Filename,
		cfg.Logging.Filename)
	if logStash != nil {
		defer logStash.Close()
	}
	if logFile != nil {
		defer logFile.Close()
	}

	// output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	// logger := zerolog.New(output).With().Str("timestamp", time.Now().String()).Logger()

	// find static fs
	var staticFS fs.FS
	if cfg.WebStatic == "" {
		staticFS, err = fs.Sub(web.WebFS, "ui/build")
		if err != nil {
			emperror.Panic(errors.Wrapf(err, "cannot create subFS of %v/%s", web.WebFS, "ui/build"))
		}
	} else {
		staticFS = os.DirFS(cfg.WebStatic)
	}

	var cert tls.Certificate
	var addCA = []*x509.Certificate{}
	if cfg.TLSCert == "" {
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
		if cert, err = tls.LoadX509KeyPair(cfg.TLSCert, cfg.TLSKey); err != nil {
			emperror.Panic(errors.Wrapf(err, "cannot load key pair %s - %s", cfg.TLSCert, cfg.TLSKey))
		}
		if cfg.RootCA != nil {
			for _, caName := range cfg.RootCA {
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
	srv, err := graphqlServer.NewServer(cfg.Addr, cfg.ExtAddr, cert, addCA, staticFS, logger, models.Keycloak{
		Addr:         cfg.Keycloak.Addr,
		Realm:        cfg.Keycloak.Realm,
		Callback:     cfg.Keycloak.Callback,
		ClientId:     cfg.Keycloak.ClientId,
		ClientSecret: cfg.Keycloak.ClientSecret,
	}, clerkHandlerServiceClient)
	if err != nil {
		emperror.Panic(errors.Wrap(err, "cannot create server"))
	}
	cancel, err := srv.Startup()
	if err != nil {
		emperror.Panic(errors.Wrap(err, "cannot start server"))
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	fmt.Println("press ctrl+c to stop server")
	s := <-done
	fmt.Println("got signal:", s)

	cancel()

}
