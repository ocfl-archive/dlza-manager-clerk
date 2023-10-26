package service

import (
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/models"
	pb "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/proto"
)

type OrderService interface {
	CopyFiles(incomingOrder models.IncomingOrder) (*pb.Status, error)
}
