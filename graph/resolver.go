package graph

import pb "gitlab.switch.ch/ub-unibas/dlza/dlza-manager/proto"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AllowedTenants     []string
	ClientClerkHandler pb.ClerkHandlerServiceClient
}
