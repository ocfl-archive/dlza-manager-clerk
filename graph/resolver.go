package graph

import (
	"github.com/je4/utils/v2/pkg/zLogger"
	pb "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	storagepb "github.com/ocfl-archive/dlza-manager-storage-handler/storagehandlerproto"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AllowedTenants            []string
	ClientClerkHandler        pb.ClerkHandlerServiceClient
	ClientClerkStorageHandler storagepb.ClerkStorageHandlerServiceClient
	Logger                    zLogger.ZLogger
}
