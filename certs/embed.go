package certs

import "embed"

//go:embed localhost.key.pem
//go:embed localhost.cert.pem
//go:embed ca.cert.pem
var CertFS embed.FS
