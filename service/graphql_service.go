package service

import (
	"context"

	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/graph/model"
	pb "gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/proto"

	"slices"

	"emperror.dev/errors"
)

func GetTenants(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, options *model.TenantListOptions, allowedTenants []string) (*model.TenantList, error) {
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
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 tenants")
			}
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
	}
	tenantsPb, err := clientClerkHandler.FindAllTenantsPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey, AllowedTenants: allowedTenants})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not FindAllTenants: %v", err)
	}
	tenants := make([]*model.Tenant, 0)
	for _, tenantPb := range tenantsPb.Tenants {
		tenant := tenantToGraphQlTenant(tenantPb)
		tenants = append(tenants, tenant)
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
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				sortDirection = "DESC"
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 storageLocations")
			}
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
		storageLocation := storageLocationToGraphQlStorageLocation(storageLocationPb)
		storageLocation.Tenant = obj
		storageLocations = append(storageLocations, storageLocation)
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
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				sortDirection = "DESC"
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 collections")
			}
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
		collection := collectionToGraphQlCollection(collectionPb)
		collection.Tenant = obj
		collections = append(collections, collection)
	}
	return &model.CollectionList{Items: collections, TotalItems: len(collections)}, nil
}

func GetCollectionsForTenantId(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, options *model.CollectionListOptions, allowedTenants []string) (*model.CollectionList, error) {
	sortKey := "ID"
	sortDirection := "ASC"
	take := 10
	skip := 0
	tenantId := ""
	if options != nil {
		if options.SortKey != nil {
			sortKey = options.SortKey.String()
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				sortDirection = "DESC"
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 collections")
			}
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
		if options.TenantID != nil {
			tenantId = *options.TenantID
		}
	}
	collectionsPb, err := clientClerkHandler.GetCollectionsByTenantIdPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey, Id: tenantId, AllowedTenants: allowedTenants})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionsByTenantID: %v", err)
	}
	tenantsMap := make(map[string]*model.Tenant)
	collections := make([]*model.Collection, 0)
	for _, collectionPb := range collectionsPb.Collections {
		collection := collectionToGraphQlCollection(collectionPb)
		if tenantsMap[collection.TenantID] == nil {
			tenantPb, err := clientClerkHandler.FindTenantById(ctx, &pb.Id{Id: collection.TenantID})
			if err != nil {
				return nil, errors.Wrapf(err, "Could not FindTenantById: %v", err)
			}
			tenantsMap[collection.TenantID] = tenantToGraphQlTenant(tenantPb)
		}
		collection.Tenant = tenantsMap[collection.TenantID]
		collections = append(collections, collection)
	}
	return &model.CollectionList{Items: collections, TotalItems: len(collections)}, nil
}

func GetObjectsForCollection(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, obj *model.Collection, options *model.ObjectListOptions) (*model.ObjectList, error) {
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
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 objects")
			}
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
	}
	objectsPb, err := clientClerkHandler.GetObjectsByCollectionIdPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey, Id: obj.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionsByTenantID: %v", err)
	}
	objects := make([]*model.Object, 0)
	for _, objectPb := range objectsPb.Objects {
		object := objectToGraphQlObject(objectPb)
		object.Collection = obj
		objects = append(objects, object)
	}
	return &model.ObjectList{Items: objects, TotalItems: len(objects)}, nil
}

func GetObjectsForCollectionId(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, options *model.ObjectListOptions, allowedTenants []string) (*model.ObjectList, error) {
	sortKey := "ID"
	sortDirection := "ASC"
	take := 10
	skip := 0
	collectionId := ""

	if options != nil {
		if options.SortKey != nil {
			sortKey = options.SortKey.String()
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				sortDirection = "DESC"
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 objects")
			}
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
		if options.CollectionID != nil {
			collectionId = *options.CollectionID
		}
	}
	objectsPb, err := clientClerkHandler.GetObjectsByCollectionIdPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey, Id: collectionId, AllowedTenants: allowedTenants})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionsByTenantID: %v", err)
	}

	collectionsMap := make(map[string]*model.Collection)
	objects := make([]*model.Object, 0)
	for _, objectPb := range objectsPb.Objects {
		object := objectToGraphQlObject(objectPb)
		if collectionsMap[object.CollectionID] == nil {
			collectionPb, err := clientClerkHandler.GetCollectionById(ctx, &pb.Id{Id: object.CollectionID})
			if err != nil {
				return nil, errors.Wrapf(err, "Could not GetCollectionById: %v", err)
			}
			collectionsMap[object.CollectionID] = collectionToGraphQlCollection(collectionPb)
		}
		object.Collection = collectionsMap[object.CollectionID]
		objects = append(objects, object)
	}
	return &model.ObjectList{Items: objects, TotalItems: len(objects)}, nil
}

func GetObjectInstancesForObject(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, obj *model.Object, options *model.ObjectInstanceListOptions) (*model.ObjectInstanceList, error) {
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
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 objectInstances")
			}
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
	}
	objectInstancesPb, err := clientClerkHandler.GetObjectInstancesByObjectIdPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey, Id: obj.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstancesByObjectIdPaginated: %v", err)
	}
	objectInstances := make([]*model.ObjectInstance, 0)
	for _, objectInstancePb := range objectInstancesPb.ObjectInstances {
		objectInstance := objectInstanceToGraphQlObjectInstance(objectInstancePb)
		objectInstance.Object = obj
		objectInstances = append(objectInstances, objectInstance)
	}
	return &model.ObjectInstanceList{Items: objectInstances, TotalItems: len(objectInstances)}, nil
}

func GetFilesForObject(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, obj *model.Object, options *model.FileListOptions) (*model.FileList, error) {
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
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 Files")
			}
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
	}
	filesPb, err := clientClerkHandler.GetFilesByObjectIdPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey, Id: obj.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetFilesByObjectIdPaginated: %v", err)
	}
	files := make([]*model.File, 0)
	for _, filePb := range filesPb.Files {
		file := fileToGraphQlFile(filePb)
		file.Object = obj
		files = append(files, file)
	}
	return &model.FileList{Items: files, TotalItems: len(files)}, nil
}

func GetObjectInstancesForObjectId(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, options *model.ObjectInstanceListOptions, allowedTenants []string) (*model.ObjectInstanceList, error) {
	sortKey := "ID"
	sortDirection := "ASC"
	take := 10
	skip := 0
	objectId := ""

	if options != nil {
		if options.SortKey != nil {
			sortKey = options.SortKey.String()
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				sortDirection = "DESC"
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 objectInstances")
			}
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
		if options.ObjectID != nil {
			objectId = *options.ObjectID
		}
	}
	objectInstancesPb, err := clientClerkHandler.GetObjectInstancesByObjectIdPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey, Id: objectId, AllowedTenants: allowedTenants})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstancesByObjectIdPaginated: %v", err)
	}
	partitionsMap := make(map[string]*model.StoragePartition)
	objectsMap := make(map[string]*model.Object)
	objectInstances := make([]*model.ObjectInstance, 0)
	for _, objectInstancePb := range objectInstancesPb.ObjectInstances {
		objectInstance := objectInstanceToGraphQlObjectInstance(objectInstancePb)
		if objectsMap[objectInstance.ObjectID] == nil {
			objectPb, err := clientClerkHandler.GetObjectById(ctx, &pb.Id{Id: objectInstance.ObjectID})
			if err != nil {
				return nil, errors.Wrapf(err, "Could not GetObjectById: %v", err)
			}
			objectsMap[objectInstance.ObjectID] = objectToGraphQlObject(objectPb)
		}
		if partitionsMap[objectInstance.StoragePartitionID] == nil {
			storagePartitionPb, err := clientClerkHandler.GetStoragePartitionById(ctx, &pb.Id{Id: objectInstance.StoragePartitionID})
			if err != nil {
				return nil, errors.Wrapf(err, "Could not GetStoragePartitionById: %v", err)
			}
			partitionsMap[objectInstance.StoragePartitionID] = storagePartitionToGraphQlStoragePartition(storagePartitionPb)
		}
		objectInstance.Object = objectsMap[objectInstance.ObjectID]
		objectInstance.StoragePartition = partitionsMap[objectInstance.StoragePartitionID]
		objectInstances = append(objectInstances, objectInstance)
	}
	return &model.ObjectInstanceList{Items: objectInstances, TotalItems: len(objectInstances)}, nil
}

func GetFilesForObjectId(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, options *model.FileListOptions, allowedTenants []string) (*model.FileList, error) {
	sortKey := "ID"
	sortDirection := "ASC"
	take := 10
	skip := 0
	objectId := ""

	if options != nil {
		if options.SortKey != nil {
			sortKey = options.SortKey.String()
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				sortDirection = "DESC"
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 files")
			}
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
		if options.ObjectID != nil {
			objectId = *options.ObjectID
		}
	}
	filesPb, err := clientClerkHandler.GetFilesByObjectIdPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey, Id: objectId, AllowedTenants: allowedTenants})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetFilesByObjectIdPaginated: %v", err)
	}
	objectsMap := make(map[string]*model.Object)
	files := make([]*model.File, 0)
	for _, filePb := range filesPb.Files {
		file := fileToGraphQlFile(filePb)
		if objectsMap[file.ObjectID] == nil {
			objectPb, err := clientClerkHandler.GetObjectById(ctx, &pb.Id{Id: file.ObjectID})
			if err != nil {
				return nil, errors.Wrapf(err, "Could not GetObjectById: %v", err)
			}
			objectsMap[file.ObjectID] = objectToGraphQlObject(objectPb)
		}
		file.Object = objectsMap[file.ObjectID]
		files = append(files, file)
	}
	return &model.FileList{Items: files, TotalItems: len(files)}, nil
}

func GetObjectInstanceChecksForObjectInstance(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, obj *model.ObjectInstance, options *model.ObjectInstanceCheckListOptions) (*model.ObjectInstanceCheckList, error) {
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
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 objectInstanceChecks")
			}
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
	}
	objectInstanceChecksPb, err := clientClerkHandler.GetObjectInstanceChecksByObjectInstanceIdPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey, Id: obj.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstanceChecksByObjectInstanceIdPaginated: %v", err)
	}
	objectInstanceChecks := make([]*model.ObjectInstanceCheck, 0)
	for _, objectInstanceCheckPb := range objectInstanceChecksPb.ObjectInstanceChecks {
		objectInstanceCheck := objectInstanceCheckToGraphQlObjectInstanceCheck(objectInstanceCheckPb)
		objectInstanceCheck.ObjectInstance = obj
		objectInstanceChecks = append(objectInstanceChecks, objectInstanceCheck)
	}
	return &model.ObjectInstanceCheckList{Items: objectInstanceChecks, TotalItems: len(objectInstanceChecks)}, nil
}

func GetObjectInstanceChecksForObjectInstanceId(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, options *model.ObjectInstanceCheckListOptions, allowedTenants []string) (*model.ObjectInstanceCheckList, error) {
	sortKey := "ID"
	sortDirection := "ASC"
	take := 10
	skip := 0
	objectInstanceId := ""

	if options != nil {
		if options.SortKey != nil {
			sortKey = options.SortKey.String()
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				sortDirection = "DESC"
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 objectInstanceChecks")
			}
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
		if options.ObjectInstanceID != nil {
			objectInstanceId = *options.ObjectInstanceID
		}
	}
	objectInstanceChecksPb, err := clientClerkHandler.GetObjectInstanceChecksByObjectInstanceIdPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey, Id: objectInstanceId, AllowedTenants: allowedTenants})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstanceChecksByObjectInstanceIdPaginated: %v", err)
	}
	objectInstancesMap := make(map[string]*model.ObjectInstance)
	objectInstanceChecks := make([]*model.ObjectInstanceCheck, 0)
	for _, objectInstanceCheckPb := range objectInstanceChecksPb.ObjectInstanceChecks {
		objectInstanceCheck := objectInstanceCheckToGraphQlObjectInstanceCheck(objectInstanceCheckPb)
		if objectInstancesMap[objectInstanceCheck.ObjectInstanceID] == nil {
			objectInstancePb, err := clientClerkHandler.GetObjectInstanceById(ctx, &pb.Id{Id: objectInstanceCheck.ObjectInstanceID})
			if err != nil {
				return nil, errors.Wrapf(err, "Could not GetObjectInstanceById: %v", err)
			}
			objectInstancesMap[objectInstanceCheck.ObjectInstanceID] = objectInstanceToGraphQlObjectInstance(objectInstancePb)
		}
		objectInstanceCheck.ObjectInstance = objectInstancesMap[objectInstanceCheck.ObjectInstanceID]
		objectInstanceChecks = append(objectInstanceChecks, objectInstanceCheck)
	}
	return &model.ObjectInstanceCheckList{Items: objectInstanceChecks, TotalItems: len(objectInstanceChecks)}, nil
}

func GetStorageLocationsForTenantId(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, options *model.StorageLocationListOptions, allowedTenants []string) (*model.StorageLocationList, error) {
	sortKey := "ID"
	sortDirection := "ASC"
	take := 10
	skip := 0
	tenantId := ""

	if options != nil {
		if options.SortKey != nil {
			sortKey = options.SortKey.String()
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				sortDirection = "DESC"
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 storageLocations")
			}
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
		if options.TenantID != nil {
			tenantId = *options.TenantID
		}
	}
	storageLocationsPb, err := clientClerkHandler.GetStorageLocationsByTenantIdPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey, Id: tenantId, AllowedTenants: allowedTenants})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStorageLocationsByTenantIdPaginated: %v", err)
	}
	tenantsMap := make(map[string]*model.Tenant)
	storageLocations := make([]*model.StorageLocation, 0)
	for _, storageLocationPb := range storageLocationsPb.StorageLocations {
		storageLocation := storageLocationToGraphQlStorageLocation(storageLocationPb)
		if tenantsMap[storageLocation.TenantID] == nil {
			tenantPb, err := clientClerkHandler.FindTenantById(ctx, &pb.Id{Id: storageLocation.TenantID})
			if err != nil {
				return nil, errors.Wrapf(err, "Could not FindTenantById: %v", err)
			}
			tenantsMap[storageLocation.TenantID] = tenantToGraphQlTenant(tenantPb)
		}
		storageLocation.Tenant = tenantsMap[storageLocation.TenantID]
		storageLocations = append(storageLocations, storageLocation)
	}
	return &model.StorageLocationList{Items: storageLocations, TotalItems: len(storageLocations)}, nil
}

func GetStoragePartitionsForLocationId(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, options *model.StoragePartitionListOptions, allowedTenants []string) (*model.StoragePartitionList, error) {
	sortKey := "ID"
	sortDirection := "ASC"
	take := 10
	skip := 0
	storageLocationId := ""

	if options != nil {
		if options.SortKey != nil {
			sortKey = options.SortKey.String()
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				sortDirection = "DESC"
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 storagePartitions")
			}
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
		if options.StorageLocationID != nil {
			storageLocationId = *options.StorageLocationID
		}
	}
	storagePartitionsPb, err := clientClerkHandler.GetStoragePartitionsByLocationIdPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey, Id: storageLocationId, AllowedTenants: allowedTenants})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStoragePartitionsByLocationIdPaginated: %v", err)
	}
	storageLocationsMap := make(map[string]*model.StorageLocation)
	storagePartitions := make([]*model.StoragePartition, 0)
	for _, storagePartitionPb := range storagePartitionsPb.StoragePartitions {
		storagePartition := storagePartitionToGraphQlStoragePartition(storagePartitionPb)
		if storageLocationsMap[storagePartition.StorageLocationID] == nil {
			storageLocationPb, err := clientClerkHandler.GetStorageLocationById(ctx, &pb.Id{Id: storagePartition.StorageLocationID})
			if err != nil {
				return nil, errors.Wrapf(err, "Could not GetStorageLocationById: %v", err)
			}
			storageLocationsMap[storagePartition.StorageLocationID] = storageLocationToGraphQlStorageLocation(storageLocationPb)
		}
		storagePartition.StorageLocation = storageLocationsMap[storagePartition.StorageLocationID]
		storagePartitions = append(storagePartitions, storagePartition)
	}
	return &model.StoragePartitionList{Items: storagePartitions, TotalItems: len(storagePartitions)}, nil
}

func GetStoragePartitionsForLocation(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, obj *model.StorageLocation, options *model.StoragePartitionListOptions) (*model.StoragePartitionList, error) {
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
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 storagePartitions")
			}
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
	}
	storagePartitionsPb, err := clientClerkHandler.GetStoragePartitionsByLocationIdPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey, Id: obj.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStoragePartitionsByLocationIdPaginated: %v", err)
	}
	storagePartitions := make([]*model.StoragePartition, 0)
	for _, storagePartitionPb := range storagePartitionsPb.StoragePartitions {
		storagePartition := storagePartitionToGraphQlStoragePartition(storagePartitionPb)
		storagePartition.StorageLocation = obj
		storagePartitions = append(storagePartitions, storagePartition)
	}
	return &model.StoragePartitionList{Items: storagePartitions, TotalItems: len(storagePartitions)}, nil
}

func GetObjectInstancesForStoragePartition(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, obj *model.StoragePartition, options *model.ObjectInstanceListOptions) (*model.ObjectInstanceList, error) {
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
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 objectInstances")
			}
			take = *options.Take
		}
		if options.Skip != nil {
			skip = *options.Skip
		}
	}
	objectInstancesPb, err := clientClerkHandler.GetObjectInstancesByStoragePartitionIdPaginated(ctx, &pb.Pagination{Skip: int32(skip), Take: int32(take), SortDirection: sortDirection, SortKey: sortKey, Id: obj.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstancesByStoragePartitionIdPaginated: %v", err)
	}
	objectInstances := make([]*model.ObjectInstance, 0)
	for _, objectInstancePb := range objectInstancesPb.ObjectInstances {
		objectInstance := objectInstanceToGraphQlObjectInstance(objectInstancePb)
		objectInstance.StoragePartition = obj
		objectInstances = append(objectInstances, objectInstance)
	}
	return &model.ObjectInstanceList{Items: objectInstances, TotalItems: len(objectInstances)}, nil
}

func GetTenantById(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, id string, allowedTenants []string) (*model.Tenant, error) {
	if len(allowedTenants) != 0 {
		if !slices.Contains(allowedTenants, id) {
			return nil, errors.New("This user hasn't rights to access this information")
		}
	}
	tenantPb, err := clientClerkHandler.FindTenantById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not FindTenantById: %v", err)
	}
	tenant := tenantToGraphQlTenant(tenantPb)
	return tenant, err
}

func GetCollectionById(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, id string) (*model.Collection, error) {
	collectionPb, err := clientClerkHandler.GetCollectionById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionById: %v", err)
	}
	collection := collectionToGraphQlCollection(collectionPb)
	tenant, err := clientClerkHandler.FindTenantById(ctx, &pb.Id{Id: collection.TenantID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not FindTenantById: %v", err)
	}
	collection.Tenant = tenantToGraphQlTenant(tenant)
	return collection, err
}

func GetObjectById(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, id string) (*model.Object, error) {
	objectPb, err := clientClerkHandler.GetObjectById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionById: %v", err)
	}
	object := objectToGraphQlObject(objectPb)
	collectionPb, err := clientClerkHandler.GetCollectionById(ctx, &pb.Id{Id: object.CollectionID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionById: %v", err)
	}
	collection := collectionToGraphQlCollection(collectionPb)
	object.Collection = collection
	return object, err
}

func GetObjectInstanceById(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, id string) (*model.ObjectInstance, error) {
	objectInstancePb, err := clientClerkHandler.GetObjectInstanceById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstanceById: %v", err)
	}
	objectInstance := objectInstanceToGraphQlObjectInstance(objectInstancePb)
	objectPb, err := clientClerkHandler.GetObjectById(ctx, &pb.Id{Id: objectInstance.ObjectID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionById: %v", err)
	}
	object := objectToGraphQlObject(objectPb)
	objectInstance.Object = object
	return objectInstance, err
}

func GetObjectInstanceCheckById(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, id string) (*model.ObjectInstanceCheck, error) {
	objectInstanceCheckPb, err := clientClerkHandler.GetObjectInstanceCheckById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstanceCheckById: %v", err)
	}
	objectInstanceCheck := objectInstanceCheckToGraphQlObjectInstanceCheck(objectInstanceCheckPb)
	objectInstancePb, err := clientClerkHandler.GetObjectInstanceById(ctx, &pb.Id{Id: objectInstanceCheck.ObjectInstanceID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstanceById: %v", err)
	}
	objectInstance := objectInstanceToGraphQlObjectInstance(objectInstancePb)
	objectInstanceCheck.ObjectInstance = objectInstance
	return objectInstanceCheck, err
}

func GetFileById(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, id string) (*model.File, error) {
	filePb, err := clientClerkHandler.GetFileById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetFileById: %v", err)
	}
	file := fileToGraphQlFile(filePb)
	objectPb, err := clientClerkHandler.GetObjectById(ctx, &pb.Id{Id: file.ObjectID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionById: %v", err)
	}
	object := objectToGraphQlObject(objectPb)
	file.Object = object
	return file, err
}

func GetStorageLocationById(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, id string) (*model.StorageLocation, error) {
	storageLocationPb, err := clientClerkHandler.GetStorageLocationById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStorageLocationById: %v", err)
	}
	storageLocation := storageLocationToGraphQlStorageLocation(storageLocationPb)
	tenantPb, err := clientClerkHandler.FindTenantById(ctx, &pb.Id{Id: storageLocation.TenantID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not FindTenantById: %v", err)
	}
	tenant := tenantToGraphQlTenant(tenantPb)
	storageLocation.Tenant = tenant
	return storageLocation, err
}

func GetStoragePartitionById(ctx context.Context, clientClerkHandler pb.ClerkHandlerServiceClient, id string) (*model.StoragePartition, error) {
	storagePartitionPb, err := clientClerkHandler.GetStoragePartitionById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStoragePartitionById: %v", err)
	}
	storagePartition := storagePartitionToGraphQlStoragePartition(storagePartitionPb)
	storageLocationPb, err := clientClerkHandler.GetStorageLocationById(ctx, &pb.Id{Id: storagePartition.StorageLocationID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStorageLocationById: %v", err)
	}
	storageLocation := storageLocationToGraphQlStorageLocation(storageLocationPb)
	storagePartition.StorageLocation = storageLocation
	return storagePartition, err
}

func tenantToGraphQlTenant(tenantPb *pb.Tenant) *model.Tenant {
	var tenant model.Tenant
	tenant.ID = tenantPb.Id
	tenant.Email = tenantPb.Email
	tenant.Person = tenantPb.Person
	tenant.Name = tenantPb.Name
	tenant.Alias = tenantPb.Alias
	return &tenant
}

func collectionToGraphQlCollection(collectionPb *pb.Collection) *model.Collection {
	var collection model.Collection
	collection.ID = collectionPb.Id
	collection.Alias = collectionPb.Alias
	collection.Description = collectionPb.Description
	collection.Owner = collectionPb.Owner
	collection.Name = collectionPb.Name
	collection.OwnerMail = collectionPb.OwnerMail
	collection.Quality = int(collectionPb.Quality)
	collection.TenantID = collectionPb.TenantId
	return &collection
}

func objectToGraphQlObject(objectPb *pb.Object) *model.Object {
	var object model.Object
	object.Signature = objectPb.Signature
	object.Sets = objectPb.Sets
	object.Identifiers = objectPb.Identifiers
	object.Title = objectPb.Title
	object.AlternativeTitles = objectPb.AlternativeTitles
	object.Description = objectPb.Description
	object.Keywords = objectPb.Keywords
	object.References = objectPb.References
	object.IngestWorkflow = objectPb.IngestWorkflow
	object.User = objectPb.User
	object.Address = objectPb.Address
	object.Created = objectPb.Created
	object.LastChanged = objectPb.LastChanged
	object.Size = int(objectPb.Size)
	object.ID = objectPb.Id
	object.CollectionID = objectPb.CollectionId
	object.Checksum = objectPb.Checksum
	return &object
}

func objectInstanceToGraphQlObjectInstance(objectInstancePb *pb.ObjectInstance) *model.ObjectInstance {
	var objectInstance model.ObjectInstance
	objectInstance.Path = objectInstancePb.Path
	objectInstance.Size = int(objectInstancePb.Size)
	objectInstance.Created = objectInstancePb.Created
	objectInstance.Status = objectInstancePb.Status
	objectInstance.ID = objectInstancePb.Id
	objectInstance.StoragePartitionID = objectInstancePb.StoragePartitionId
	objectInstance.ObjectID = objectInstancePb.ObjectId
	return &objectInstance
}

func fileToGraphQlFile(filePb *pb.File) *model.File {
	var file model.File
	file.Checksum = filePb.Checksum
	file.Name = filePb.Name
	file.Size = int(filePb.Size)
	file.Mimetype = filePb.Mimetype
	file.Pronom = filePb.Pronom
	file.Width = int(filePb.Width)
	file.Height = int(filePb.Height)
	file.Duration = int(filePb.Duration)
	file.ID = filePb.Id
	file.ObjectID = filePb.ObjectId
	return &file
}

func objectInstanceCheckToGraphQlObjectInstanceCheck(objectInstanceCheckPb *pb.ObjectInstanceCheck) *model.ObjectInstanceCheck {
	var objectInstanceCheck model.ObjectInstanceCheck
	objectInstanceCheck.Checktime = objectInstanceCheckPb.CheckTime
	objectInstanceCheck.Error = objectInstanceCheckPb.Error
	objectInstanceCheck.Message = objectInstanceCheckPb.Message
	objectInstanceCheck.ID = objectInstanceCheckPb.Id
	objectInstanceCheck.ObjectInstanceID = objectInstanceCheckPb.ObjectInstanceId
	return &objectInstanceCheck
}

func storageLocationToGraphQlStorageLocation(storageLocationPb *pb.StorageLocation) *model.StorageLocation {
	var storageLocation model.StorageLocation
	storageLocation.ID = storageLocationPb.Id
	storageLocation.Alias = storageLocationPb.Alias
	storageLocation.Type = storageLocationPb.Type
	storageLocation.Vault = storageLocationPb.Vault
	storageLocation.Connection = storageLocationPb.Connection
	storageLocation.Quality = int(storageLocationPb.Quality)
	storageLocation.Price = int(storageLocationPb.Price)
	storageLocation.SecurityCompliency = storageLocationPb.SecurityCompliency
	storageLocation.FillFirst = storageLocationPb.FillFirst
	storageLocation.OcflType = storageLocationPb.OcflType
	storageLocation.NumberOfThreads = int(storageLocationPb.NumberOfThreads)
	storageLocation.TenantID = storageLocationPb.TenantId
	return &storageLocation
}

func storagePartitionToGraphQlStoragePartition(storagePartitionPb *pb.StoragePartition) *model.StoragePartition {
	var storagePartition model.StoragePartition
	storagePartition.Alias = storagePartitionPb.Alias
	storagePartition.Name = storagePartitionPb.Name
	storagePartition.MaxSize = int(storagePartitionPb.MaxSize)
	storagePartition.MaxObjects = int(storagePartitionPb.MaxObjects)
	storagePartition.CurrentSize = int(storagePartitionPb.CurrentSize)
	storagePartition.CurrentObjects = int(storagePartitionPb.CurrentObjects)
	storagePartition.ID = storagePartitionPb.Id
	storagePartition.StorageLocationID = storagePartitionPb.StorageLocationId

	return &storagePartition
}
