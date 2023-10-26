package service

import (
	"context"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/graph/model"
	pb "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/proto"

	"emperror.dev/errors"
)

func GetTenants(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, options *model.TenantListOptions) (*model.TenantList, error) {
	sortKey := "ID"
	sortDirection := "ASC"
	take := 10
	skip := 0
	if options != nil {
		if options.SortKey != nil {
			sortKey = options.SortKey.String()
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				sortDirection = "DESC"
			}
		}
		if *options.Take > 1000 {
			return nil, errors.New("You could not retrieve more than 1000 tenants")
		}
		if options.Take != nil {
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
	}
	tenantsPb, err := clientClerkHandler.FindAllTenantsPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not FindAllTenants: %v", err)
	}
	tenants := make([]*model.Tenant, 0)
	for _, tenantPb := range tenantsPb.Tenants {
		tenant := model.Tenant{}
		tenant.Name = tenantPb.Name
		tenant.Alias = tenantPb.Alias
		tenant.ID = tenantPb.Id
		tenant.Email = tenantPb.Email
		tenant.Person = tenantPb.Person
		tenants = append(tenants, &tenant)
	}
	return &model.TenantList{Items: tenants, TotalItems: len(tenants)}, nil
}

func GetStorageLocationsForTenant(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, obj *model.Tenant, options *model.StorageLocationListOptions) (*model.StorageLocationList, error) {
	sortKey := "ID"
	sortDirection := "ASC"
	take := 10
	skip := 0
	if options != nil {
		if options.SortKey != nil {
			sortKey = options.SortKey.String()
		}
		if options.SortKey != nil {
			sortKey = options.SortKey.String()
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				sortDirection = "DESC"
			}
		}
		if *options.Take > 1000 {
			return nil, errors.New("You could not retrieve more than 1000 storageLocations")
		}
		if options.Take != nil {
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
	}
	storageLocationsPb, err := clientClerkHandler.GetStorageLocationsByTenantIdPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey, Id: obj.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStorageLocationsByTenantID: %v", err)
	}
	storageLocations := make([]*model.StorageLocation, 0)
	for _, storageLocationPb := range storageLocationsPb.StorageLocations {
		storageLocation := model.StorageLocation{}
		storageLocation.TenantID = storageLocationPb.TenantId
		storageLocation.Alias = storageLocationPb.Alias
		storageLocation.ID = storageLocationPb.Id
		storageLocation.Connection = storageLocationPb.Connection
		storageLocation.Quality = int(storageLocationPb.Quality)
		storageLocation.FillFirst = storageLocationPb.FillFirst
		storageLocation.NumberOfThreads = int(storageLocationPb.NumberOfThreads)
		storageLocation.OcflType = storageLocationPb.OcflType
		storageLocation.Price = int(storageLocationPb.Price)
		storageLocation.Type = storageLocationPb.Type
		storageLocation.SecurityCompliency = storageLocationPb.SecurityCompliency
		storageLocation.Vault = storageLocationPb.Vault
		storageLocations = append(storageLocations, &storageLocation)
	}
	return &model.StorageLocationList{Items: storageLocations, TotalItems: len(storageLocations)}, nil
}

func GetCollectionsForTenant(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, obj *model.Tenant, options *model.CollectionListOptions) (*model.CollectionList, error) {
	sortKey := "ID"
	sortDirection := "ASC"
	take := 10
	skip := 0
	if options != nil {
		if options.SortKey != nil {
			sortKey = options.SortKey.String()
		}
		if options.SortKey != nil {
			sortKey = options.SortKey.String()
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				sortDirection = "DESC"
			}
		}
		if *options.Take > 1000 {
			return nil, errors.New("You could not retrieve more than 1000 collections")
		}
		if options.Take != nil {
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
	}
	collectionsPb, err := clientClerkHandler.GetCollectionsByTenantIdPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey, Id: obj.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionsByTenantID: %v", err)
	}
	collections := make([]*model.Collection, 0)
	for _, collectionPb := range collectionsPb.Collections {
		collection := model.Collection{}
		collection.Name = collectionPb.Name
		collection.Alias = collectionPb.Alias
		collection.ID = collectionPb.Id
		collection.TenantID = collectionPb.TenantId
		collection.Description = collectionPb.Description
		collection.Owner = collectionPb.Owner
		collection.OwnerMail = collectionPb.OwnerMail
		collection.Quality = int(collectionPb.Quality)
		collections = append(collections, &collection)
	}
	return &model.CollectionList{Items: collections, TotalItems: len(collections)}, nil
}
