package graph

import pb "github.com/ocfl-archive/dlza-manager-handler/handlerproto"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AllowedTenants     []string
	ClientClerkHandler pb.ClerkHandlerServiceClient
}
