package models

import (
	"crypto/tls"
	"crypto/x509"
	"io/fs"
	"net/http"

	ub_logger "gitlab.switch.ch/ub-unibas/go-ublogger"
)

type Server struct {
	ExtAddr  string
	Server   http.Server
	StaticFS fs.FS
	Addr     string
	Cert     tls.Certificate
	AddCAs   []*x509.Certificate
	Logger   *ub_logger.Logger
	Keycloak Keycloak
}
