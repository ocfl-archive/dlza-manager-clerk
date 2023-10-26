package service

import (
	"context"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/models"
	pb "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/proto"
	"time"

	"github.com/pkg/errors"
)

type OrderServiceImpl struct {
	ClientClerkIngestService pb.ClerkIngestServiceClient
}

func NewOrderService(clientClerkIngestService pb.ClerkIngestServiceClient) OrderService {
	return &OrderServiceImpl{ClientClerkIngestService: clientClerkIngestService}
}

func (o *OrderServiceImpl) CopyFiles(incomingOrder models.IncomingOrder) (*pb.Status, error) {

	objectPathsPb := make([]*pb.ObjectPath, 0)
	for _, objectPath := range incomingOrder.ObjectPaths {
		objectPathPb := &pb.ObjectPath{}
		objectPathPb.FilePath = objectPath.FilePath
		objectPathPb.InfoFilePath = objectPath.InfoFilePath

		objectPathsPb = append(objectPathsPb, objectPathPb)
	}

	incomingOrderPb := &pb.IncomingOrder{CollectionAlias: incomingOrder.CollectionAlias, ObjectPaths: objectPathsPb}
	c := context.Background()
	cont, cancel := context.WithTimeout(c, 10000*time.Second)
	defer cancel()
	status, err := o.ClientClerkIngestService.CopyFiles(cont, incomingOrderPb)
	if err != nil {
		return &pb.Status{Ok: false}, errors.Wrapf(err, "cannot copy file for collection '%s'", incomingOrder.CollectionAlias)
	}

	return status, nil
}
