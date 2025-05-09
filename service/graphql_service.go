package service

import (
	"context"
	"encoding/json"
	"github.com/je4/utils/v2/pkg/zLogger"
	"github.com/ocfl-archive/dlza-manager-clerk/models"
	pb "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	"golang.org/x/exp/maps"
	"regexp"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/ocfl-archive/dlza-manager-clerk/graph/model"
	"github.com/ocfl-archive/dlza-manager-clerk/middleware"
	pbHandler "github.com/ocfl-archive/dlza-manager-handler/handlerproto"
	"slices"
)

const (
	sortDirectionDescending string = "DESC NULLS LAST"
	sortDirectionAscending  string = "ASC"
)

func GetTenants(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, options *model.TenantListOptions, allowedTenants []string) (*model.TenantList, error) {

	keyCloakGroup, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}

	if (len(tenantList) == 0) && (!slices.Contains(keyCloakGroup, "dlza-admin")) {
		return nil, errors.New("You are not allowed to retrieve datas")
	} else if len(tenantList) > 0 {
		for _, tenant := range tenantList {
			allowedTenants = append(allowedTenants, tenant.Id)
		}
		// allowedTenants = tenantList
	}
	if slices.Contains(keyCloakGroup, "dlza-admin") {
		allowedTenants = []string{}
	}
	optionsPb := pb.Pagination{
		Take:           10,
		SortDirection:  sortDirectionAscending,
		AllowedTenants: allowedTenants,
		SortKey:        "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 tenants")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}
	tenantsPb, err := clientClerkHandler.FindAllTenantsPaginated(ctx, &optionsPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not FindAllTenants: %v", err)
	}
	tenants := make([]*model.Tenant, 0)
	for _, tenantPb := range tenantsPb.Tenants {
		tenant := tenantToGraphQlTenant(tenantPb)
		amountAndSize, err := clientClerkHandler.GetAmountOfObjectsAndTotalSizeByTenantId(ctx, &pb.Id{Id: tenant.ID})
		if err != nil {
			return nil, errors.Wrapf(err, "Could not GetAmountOfObjectsAndTotalSizeByTenantId: %v", err)
		}
		tenant.TotalAmountOfObjects = int(amountAndSize.Amount)
		tenant.TotalSize = float64(amountAndSize.Size)
		tenant.Permissions = make([]string, 0)
		if len(tenantList) > 0 {
			for _, tenantKL := range tenantList {
				if tenantKL.Id == tenant.ID {
					if tenantKL.Update && tenantKL.Delete && tenantKL.Create && tenantKL.Read {
						tenant.Permissions = append(tenant.Permissions, "collection", "storageLocation", "storagePartition")
					}
				}
			}
		}
		tenants = append(tenants, tenant)
	}
	return &model.TenantList{Items: tenants, TotalItems: int(tenantsPb.TotalItems)}, nil
}

func GetStorageLocationsForTenant(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, obj *model.Tenant, options *model.StorageLocationListOptions) (*model.StorageLocationList, error) {
	optionsPb := pb.Pagination{
		Take:          10,
		SortDirection: sortDirectionAscending,
		Id:            obj.ID,
		SortKey:       "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 storage locations")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.CollectionID != nil {
			optionsPb.SecondId = *options.CollectionID
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}
	storageLocationsPb, err := clientClerkHandler.GetStorageLocationsByTenantOrCollectionIdPaginated(ctx, &optionsPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStorageLocationsByTenantID: %v", err)
	}
	storageLocations := make([]*model.StorageLocation, 0)
	for _, storageLocationPb := range storageLocationsPb.StorageLocations {
		storageLocation := storageLocationToGraphQlStorageLocation(storageLocationPb)
		amountOfObjects, err := clientClerkHandler.GetAmountOfObjectsForStorageLocationId(ctx, &pb.Id{Id: storageLocation.ID})
		if err != nil {
			return nil, errors.Wrapf(err, "Could not GetAmountOfObjectsForStorageLocationId: %v", err)
		}
		amountOfErrors, err := clientClerkHandler.GetAmountOfErrorsForStorageLocationId(ctx, &pb.Id{Id: storageLocation.ID})
		if err != nil {
			return nil, errors.Wrapf(err, "Could not GetAmountOfErrorsForStorageLocationId: %v", err)
		}
		storageLocation.AmountOfErrors = int(amountOfErrors.Size)
		storageLocation.AmountOfObjects = int(amountOfObjects.Size)
		storageLocation.Tenant = obj
		storageLocations = append(storageLocations, storageLocation)
	}
	return &model.StorageLocationList{Items: storageLocations, TotalItems: int(storageLocationsPb.TotalItems)}, nil
}

func GetCollectionsForTenant(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, obj *model.Tenant, options *model.CollectionListOptions) (*model.CollectionList, error) {
	optionsPb := pb.Pagination{
		Take:          10,
		SortDirection: sortDirectionAscending,
		Id:            obj.ID,
		SortKey:       "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 collections")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}
	collectionsPb, err := clientClerkHandler.GetCollectionsByTenantIdPaginated(ctx, &optionsPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionsByTenantID: %v", err)
	}
	collections := make([]*model.Collection, 0)
	for _, collectionPb := range collectionsPb.Collections {
		collection := collectionToGraphQlCollection(collectionPb)
		amountOfErrors, err := clientClerkHandler.GetAmountOfErrorsByCollectionId(ctx, &pb.Id{Id: collection.ID})
		if err != nil {
			return nil, errors.Wrapf(err, "Could not GetAmountOfErrorsByCollectionId: %v", err)
		}
		collection.AmountOfErrors = int(amountOfErrors.Size)
		collection.Tenant = obj
		collections = append(collections, collection)
	}
	return &model.CollectionList{Items: collections, TotalItems: int(collectionsPb.TotalItems)}, nil
}

func GetCollectionsForTenantId(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, options *model.CollectionListOptions, allowedTenants []string) (*model.CollectionList, error) {
	keyCloakGroup, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if (len(tenantList) == 0) && (!slices.Contains(keyCloakGroup, "dlza-admin")) {
		return nil, errors.New("You are not allowed to retrieve datas")
	} else if len(tenantList) > 0 {
		for _, tenant := range tenantList {
			allowedTenants = append(allowedTenants, tenant.Id)
		}
		// allowedTenants = tenantList
	}
	if slices.Contains(keyCloakGroup, "dlza-admin") {
		allowedTenants = []string{}
	}
	optionsPb := pb.Pagination{
		Take:           10,
		SortDirection:  sortDirectionAscending,
		AllowedTenants: allowedTenants,
		SortKey:        "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 collections")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.TenantID != nil {
			optionsPb.Id = *options.TenantID
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}

	collectionsPb, err := clientClerkHandler.GetCollectionsByTenantIdPaginated(ctx, &optionsPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionsByTenantID: %v", err)
	}
	tenantsMap := make(map[string]*model.Tenant)
	collections := make([]*model.Collection, 0)
	for _, collectionPb := range collectionsPb.Collections {
		collection := collectionToGraphQlCollection(collectionPb)
		amountOfErrors, err := clientClerkHandler.GetAmountOfErrorsByCollectionId(ctx, &pb.Id{Id: collection.ID})
		if err != nil {
			return nil, errors.Wrapf(err, "Could not GetAmountOfErrorsByCollectionId: %v", err)
		}
		if tenantsMap[collection.TenantID] == nil {
			tenantPb, err := clientClerkHandler.FindTenantById(ctx, &pb.Id{Id: collection.TenantID})
			if err != nil {
				return nil, errors.Wrapf(err, "Could not FindTenantById: %v", err)
			}
			tenantsMap[collection.TenantID] = tenantToGraphQlTenant(tenantPb)
		}
		collection.AmountOfErrors = int(amountOfErrors.Size)
		collection.Tenant = tenantsMap[collection.TenantID]
		collections = append(collections, collection)
	}
	return &model.CollectionList{Items: collections, TotalItems: int(collectionsPb.TotalItems)}, nil
}

func GetObjectsForCollection(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, obj *model.Collection, options *model.ObjectListOptions, logger zLogger.ZLogger) (*model.ObjectList, error) {
	optionsPb := pb.Pagination{
		Take:          10,
		SortDirection: sortDirectionAscending,
		Id:            obj.ID,
		SortKey:       "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 objects")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.TenantID != nil {
			optionsPb.SecondId = *options.TenantID
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}
	logger.Debug().Msg("grpc function calling objects were executed")
	objectsPb, err := clientClerkHandler.GetObjectsByCollectionIdPaginated(ctx, &optionsPb)
	logger.Debug().Msg("grpc function calling objects returned objects")
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectsByCollectionIdPaginated: %v", err)
	}
	objects := make([]*model.Object, 0)
	for _, objectPb := range objectsPb.Objects {
		object := objectToGraphQlObject(objectPb)
		status, err := clientClerkHandler.GetStatusForObjectId(ctx, &pb.Id{Id: object.ID})
		if err != nil {
			return nil, errors.Wrapf(err, "Could not GetStatusForObjectId: %v", err)
		}
		object.Status = int(status.Size)
		object.Collection = obj
		objects = append(objects, object)
	}
	logger.Debug().Msg("returning list of objects in service method")
	return &model.ObjectList{Items: objects, TotalItems: int(objectsPb.TotalItems)}, nil
}

func GetFilesForCollection(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, obj *model.Collection, options *model.FileListOptions) (*model.FileList, error) {
	optionsPb := pb.Pagination{
		Take:          10,
		SortDirection: sortDirectionAscending,
		Id:            obj.ID,
		SortKey:       "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 files")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.TenantID != nil {
			optionsPb.SecondId = *options.TenantID
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}
	filesPb, err := clientClerkHandler.GetFilesByCollectionIdPaginated(ctx, &optionsPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetFilesByCollectionIdPaginated: %v", err)
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
	return &model.FileList{Items: files, TotalItems: int(filesPb.TotalItems)}, nil
}

func GetObjectsForCollectionId(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, options *model.ObjectListOptions, allowedTenants []string, logger zLogger.ZLogger) (*model.ObjectList, error) {
	keyCloakGroup, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if (len(tenantList) == 0) && (!slices.Contains(keyCloakGroup, "dlza-admin")) {
		return nil, errors.New("You are not allowed to retrieve datas")
	} else if len(tenantList) > 0 {
		for _, tenant := range tenantList {
			allowedTenants = append(allowedTenants, tenant.Id)
		}
		// allowedTenants = tenantList
	}
	if slices.Contains(keyCloakGroup, "dlza-admin") {
		allowedTenants = []string{}
	}
	optionsPb := pb.Pagination{
		Take:           10,
		SortDirection:  sortDirectionAscending,
		AllowedTenants: allowedTenants,
		SortKey:        "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 collections")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.CollectionID != nil {
			optionsPb.Id = *options.CollectionID
		}
		if options.TenantID != nil {
			optionsPb.SecondId = *options.TenantID
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}
	logger.Debug().Msgf("grpc function calling objects were executed %s", time.Now())
	objectsPb, err := clientClerkHandler.GetObjectsByCollectionIdPaginated(ctx, &optionsPb)
	logger.Debug().Msgf("grpc function calling objects returned objects%s", time.Now())
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionsByTenantID: %v", err)
	}

	collectionsMap := make(map[string]*model.Collection)
	objects := make([]*model.Object, 0)
	for _, objectPb := range objectsPb.Objects {
		object := objectToGraphQlObject(objectPb)
		status, err := clientClerkHandler.GetStatusForObjectId(ctx, &pb.Id{Id: object.ID})
		if err != nil {
			return nil, errors.Wrapf(err, "Could not GetStatusForObjectId: %v", err)
		}
		object.Status = int(status.Size)
		if collectionsMap[object.CollectionID] == nil {
			collectionPb, err := clientClerkHandler.GetCollectionByIdFromMv(ctx, &pb.Id{Id: object.CollectionID})
			if err != nil {
				return nil, errors.Wrapf(err, "Could not GetCollectionByIdFromMv: %v", err)
			}
			collectionsMap[object.CollectionID] = collectionToGraphQlCollection(collectionPb)
		}
		object.Collection = collectionsMap[object.CollectionID]
		objects = append(objects, object)
	}
	logger.Debug().Msgf("returning list of objects in service method%s", time.Now())
	return &model.ObjectList{Items: objects, TotalItems: int(objectsPb.TotalItems)}, nil
}

func GetObjectInstancesForObject(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, obj *model.Object, options *model.ObjectInstanceListOptions) (*model.ObjectInstanceList, error) {
	optionsPb := pb.Pagination{
		Take:          10,
		SortDirection: sortDirectionAscending,
		Id:            obj.ID,
		SortKey:       "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 object instances")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.TenantID != nil {
			optionsPb.SecondId = *options.TenantID
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}
	objectInstancesPb, err := clientClerkHandler.GetObjectInstancesByObjectIdPaginated(ctx, &optionsPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstancesByObjectIdPaginated: %v", err)
	}
	objectInstances := make([]*model.ObjectInstance, 0)
	for _, objectInstancePb := range objectInstancesPb.ObjectInstances {
		objectInstance := objectInstanceToGraphQlObjectInstance(objectInstancePb)
		objectInstance.Object = obj
		objectInstances = append(objectInstances, objectInstance)
	}
	return &model.ObjectInstanceList{Items: objectInstances, TotalItems: int(objectInstancesPb.TotalItems)}, nil
}

func GetFilesForObject(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, obj *model.Object, options *model.FileListOptions) (*model.FileList, error) {
	optionsPb := pb.Pagination{
		Take:          10,
		SortDirection: sortDirectionAscending,
		Id:            obj.ID,
		SortKey:       "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 files")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.TenantID != nil {
			optionsPb.SecondId = *options.TenantID
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}
	filesPb, err := clientClerkHandler.GetFilesByObjectIdPaginated(ctx, &optionsPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetFilesByObjectIdPaginated: %v", err)
	}
	files := make([]*model.File, 0)
	for _, filePb := range filesPb.Files {
		file := fileToGraphQlFile(filePb)
		file.Object = obj
		files = append(files, file)
	}
	return &model.FileList{Items: files, TotalItems: int(filesPb.TotalItems)}, nil
}

func GetObjectInstancesForObjectId(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, options *model.ObjectInstanceListOptions, allowedTenants []string) (*model.ObjectInstanceList, error) {
	keyCloakGroup, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if (len(tenantList) == 0) && (!slices.Contains(keyCloakGroup, "dlza-admin")) {
		return nil, errors.New("You are not allowed to retrieve datas")
	} else if len(tenantList) > 0 {
		for _, tenant := range tenantList {
			allowedTenants = append(allowedTenants, tenant.Id)
		}
		// allowedTenants = tenantList
	}
	if slices.Contains(keyCloakGroup, "dlza-admin") {
		allowedTenants = []string{}
	}
	optionsPb := pb.Pagination{
		Take:           10,
		SortDirection:  sortDirectionAscending,
		AllowedTenants: allowedTenants,
		SortKey:        "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 object instances")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.TenantID != nil {
			optionsPb.SecondId = *options.TenantID
		}
		if options.ObjectID != nil {
			optionsPb.Id = *options.ObjectID
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}
	objectInstancesPb, err := clientClerkHandler.GetObjectInstancesByObjectIdPaginated(ctx, &optionsPb)
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
	return &model.ObjectInstanceList{Items: objectInstances, TotalItems: int(objectInstancesPb.TotalItems)}, nil
}

func GetFilesForObjectId(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, options *model.FileListOptions, allowedTenants []string) (*model.FileList, error) {
	keyCloakGroup, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if (len(tenantList) == 0) && (!slices.Contains(keyCloakGroup, "dlza-admin")) {
		return nil, errors.New("You are not allowed to retrieve datas")
	} else if len(tenantList) > 0 {
		for _, tenant := range tenantList {
			allowedTenants = append(allowedTenants, tenant.Id)
		}
		// allowedTenants = tenantList
	}
	if slices.Contains(keyCloakGroup, "dlza-admin") {
		allowedTenants = []string{}
	}
	optionsPb := pb.Pagination{
		Take:           10,
		SortDirection:  sortDirectionAscending,
		AllowedTenants: allowedTenants,
		SortKey:        "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 files")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.TenantID != nil {
			optionsPb.SecondId = *options.TenantID
		}
		if options.ObjectID != nil {
			optionsPb.Id = *options.ObjectID
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}
	filesPb, err := clientClerkHandler.GetFilesByObjectIdPaginated(ctx, &optionsPb)
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
	return &model.FileList{Items: files, TotalItems: int(filesPb.TotalItems)}, nil
}

func GetObjectInstanceChecksForObjectInstance(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, obj *model.ObjectInstance, options *model.ObjectInstanceCheckListOptions) (*model.ObjectInstanceCheckList, error) {
	optionsPb := pb.Pagination{
		Take:          10,
		SortDirection: sortDirectionAscending,
		Id:            obj.ID,
		SortKey:       "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 object instances")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.TenantID != nil {
			optionsPb.SecondId = *options.TenantID
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}
	objectInstanceChecksPb, err := clientClerkHandler.GetObjectInstanceChecksByObjectInstanceIdPaginated(ctx, &optionsPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstanceChecksByObjectInstanceIdPaginated: %v", err)
	}
	objectInstanceChecks := make([]*model.ObjectInstanceCheck, 0)
	for _, objectInstanceCheckPb := range objectInstanceChecksPb.ObjectInstanceChecks {
		objectInstanceCheck := objectInstanceCheckToGraphQlObjectInstanceCheck(objectInstanceCheckPb)
		objectInstanceCheck.ObjectInstance = obj
		objectInstanceChecks = append(objectInstanceChecks, objectInstanceCheck)
	}
	return &model.ObjectInstanceCheckList{Items: objectInstanceChecks, TotalItems: int(objectInstanceChecksPb.TotalItems)}, nil
}

func GetObjectInstanceChecksForObjectInstanceId(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, options *model.ObjectInstanceCheckListOptions, allowedTenants []string) (*model.ObjectInstanceCheckList, error) {
	keyCloakGroup, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if (len(tenantList) == 0) && (!slices.Contains(keyCloakGroup, "dlza-admin")) {
		return nil, errors.New("You are not allowed to retrieve datas")
	} else if len(tenantList) > 0 {
		for _, tenant := range tenantList {
			allowedTenants = append(allowedTenants, tenant.Id)
		}
		// allowedTenants = tenantList
	}
	if slices.Contains(keyCloakGroup, "dlza-admin") {
		allowedTenants = []string{}
	}
	optionsPb := pb.Pagination{
		Take:           10,
		SortDirection:  sortDirectionAscending,
		AllowedTenants: allowedTenants,
		SortKey:        "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 object instance checks")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.TenantID != nil {
			optionsPb.SecondId = *options.TenantID
		}
		if options.ObjectInstanceID != nil {
			optionsPb.Id = *options.ObjectInstanceID
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}
	objectInstanceChecksPb, err := clientClerkHandler.GetObjectInstanceChecksByObjectInstanceIdPaginated(ctx, &optionsPb)
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
	return &model.ObjectInstanceCheckList{Items: objectInstanceChecks, TotalItems: int(objectInstanceChecksPb.TotalItems)}, nil
}

func GetStorageLocationsForTenantOrCollectionId(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, options *model.StorageLocationListOptions, allowedTenants []string) (*model.StorageLocationList, error) {
	keyCloakGroup, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if (len(tenantList) == 0) && (!slices.Contains(keyCloakGroup, "dlza-admin")) {
		return nil, errors.New("You are not allowed to retrieve datas")
	} else if len(tenantList) > 0 {
		for _, tenant := range tenantList {
			allowedTenants = append(allowedTenants, tenant.Id)
		}
		// allowedTenants = tenantList
	}
	if slices.Contains(keyCloakGroup, "dlza-admin") {
		allowedTenants = []string{}
	}
	optionsPb := pb.Pagination{
		Take:           10,
		SortDirection:  sortDirectionAscending,
		AllowedTenants: allowedTenants,
		SortKey:        "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 storage locations")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.TenantID != nil {
			optionsPb.Id = *options.TenantID
		}
		if options.CollectionID != nil {
			optionsPb.SecondId = *options.CollectionID
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}
	var storageLocationsPb *pb.StorageLocations
	storageLocationsPb, err = clientClerkHandler.GetStorageLocationsByTenantOrCollectionIdPaginated(ctx, &optionsPb)

	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStorageLocationsForTenantOrCollectionId: %v", err)
	}
	tenantsMap := make(map[string]*model.Tenant)
	storageLocations := make([]*model.StorageLocation, 0)
	for _, storageLocationPb := range storageLocationsPb.StorageLocations {
		storageLocation := storageLocationToGraphQlStorageLocation(storageLocationPb)
		amountOfObjects, err := clientClerkHandler.GetAmountOfObjectsForStorageLocationId(ctx, &pb.Id{Id: storageLocation.ID})
		if err != nil {
			return nil, errors.Wrapf(err, "Could not GetAmountOfObjectsForStorageLocationId: %v", err)
		}
		amountOfErrors, err := clientClerkHandler.GetAmountOfErrorsForStorageLocationId(ctx, &pb.Id{Id: storageLocation.ID})
		if err != nil {
			return nil, errors.Wrapf(err, "Could not GetAmountOfErrorsForStorageLocationId: %v", err)
		}
		storageLocation.AmountOfErrors = int(amountOfErrors.Size)
		storageLocation.AmountOfObjects = int(amountOfObjects.Size)
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
	return &model.StorageLocationList{Items: storageLocations, TotalItems: int(storageLocationsPb.TotalItems)}, nil
}

func GetStoragePartitionsForLocationId(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, options *model.StoragePartitionListOptions, allowedTenants []string) (*model.StoragePartitionList, error) {
	keyCloakGroup, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if (len(tenantList) == 0) && (!slices.Contains(keyCloakGroup, "dlza-admin")) {
		return nil, errors.New("You are not allowed to retrieve datas")
	} else if len(tenantList) > 0 {
		for _, tenant := range tenantList {
			allowedTenants = append(allowedTenants, tenant.Id)
		}
		// allowedTenants = tenantList
	}
	if slices.Contains(keyCloakGroup, "dlza-admin") {
		allowedTenants = []string{}
	}
	optionsPb := pb.Pagination{
		Take:           10,
		SortDirection:  sortDirectionAscending,
		AllowedTenants: allowedTenants,
		SortKey:        "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 storage partitions")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.StorageLocationID != nil {
			optionsPb.Id = *options.StorageLocationID
		}
		if options.TenantID != nil {
			optionsPb.SecondId = *options.TenantID
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}
	storagePartitionsPb, err := clientClerkHandler.GetStoragePartitionsByLocationIdPaginated(ctx, &optionsPb)
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
	return &model.StoragePartitionList{Items: storagePartitions, TotalItems: int(storagePartitionsPb.TotalItems)}, nil
}

func GetStoragePartitionsForLocation(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, obj *model.StorageLocation, options *model.StoragePartitionListOptions) (*model.StoragePartitionList, error) {
	optionsPb := pb.Pagination{
		Take:          10,
		SortDirection: sortDirectionAscending,
		Id:            obj.ID,
		SortKey:       "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 storage partitions")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.TenantID != nil {
			optionsPb.SecondId = *options.TenantID
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}
	storagePartitionsPb, err := clientClerkHandler.GetStoragePartitionsByLocationIdPaginated(ctx, &optionsPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStoragePartitionsByLocationIdPaginated: %v", err)
	}
	storagePartitions := make([]*model.StoragePartition, 0)
	for _, storagePartitionPb := range storagePartitionsPb.StoragePartitions {
		storagePartition := storagePartitionToGraphQlStoragePartition(storagePartitionPb)
		storagePartition.StorageLocation = obj
		storagePartitions = append(storagePartitions, storagePartition)
	}
	return &model.StoragePartitionList{Items: storagePartitions, TotalItems: int(storagePartitionsPb.TotalItems)}, nil
}

func GetObjectInstancesForStoragePartition(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, obj *model.StoragePartition, options *model.ObjectInstanceListOptions) (*model.ObjectInstanceList, error) {
	optionsPb := pb.Pagination{
		Take:          10,
		SortDirection: sortDirectionAscending,
		Id:            obj.ID,
		SortKey:       "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 object instances")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.TenantID != nil {
			optionsPb.SecondId = *options.TenantID
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.Search != nil {
			optionsPb.SearchField = strings.ToLower(*options.Search)
		}
	}
	objectInstancesPb, err := clientClerkHandler.GetObjectInstancesByStoragePartitionIdPaginated(ctx, &optionsPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstancesByStoragePartitionIdPaginated: %v", err)
	}
	objectInstances := make([]*model.ObjectInstance, 0)
	for _, objectInstancePb := range objectInstancesPb.ObjectInstances {
		objectInstance := objectInstanceToGraphQlObjectInstance(objectInstancePb)
		objectInstance.StoragePartition = obj
		objectInstances = append(objectInstances, objectInstance)
	}
	return &model.ObjectInstanceList{Items: objectInstances, TotalItems: int(objectInstancesPb.TotalItems)}, nil
}

func GetTenantById(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, id string, allowedTenants []string) (*model.Tenant, error) {
	keyCloakGroup, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if (len(tenantList) == 0) && (!slices.Contains(keyCloakGroup, "dlza-admin")) {
		return nil, errors.New("You are not allowed to retrieve datas")
	} else if len(tenantList) > 0 {
		for _, tenant := range tenantList {
			allowedTenants = append(allowedTenants, tenant.Id)
		}
		// allowedTenants = tenantList
	}
	if slices.Contains(keyCloakGroup, "dlza-admin") {
		allowedTenants = []string{}
	}
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
	tenant.Permissions = make([]string, 0)
	if len(tenantList) > 0 {
		for _, tenantKL := range tenantList {
			if tenantKL.Id == id {
				if tenantKL.Update && tenantKL.Delete && tenantKL.Create && tenantKL.Read {
					tenant.Permissions = append(tenant.Permissions, "collection", "storageLocation", "storagePartition")
				}
			}
		}
	}
	amountAndSize, err := clientClerkHandler.GetAmountOfObjectsAndTotalSizeByTenantId(ctx, &pb.Id{Id: tenant.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetAmountOfObjectsAndTotalSizeByTenantId: %v", err)
	}
	tenant.TotalAmountOfObjects = int(amountAndSize.Amount)
	tenant.TotalSize = float64(amountAndSize.Size)
	return tenant, err
}

func GetCollectionById(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, id string) (*model.Collection, error) {
	collectionPb, err := clientClerkHandler.GetCollectionByIdFromMv(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionByIdFromMv: %v", err)
	}
	collection := collectionToGraphQlCollection(collectionPb)
	amountOfErrors, err := clientClerkHandler.GetAmountOfErrorsByCollectionId(ctx, &pb.Id{Id: collection.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetAmountOfErrorsByCollectionId: %v", err)
	}
	collection.AmountOfErrors = int(amountOfErrors.Size)
	tenant, err := clientClerkHandler.FindTenantById(ctx, &pb.Id{Id: collection.TenantID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not FindTenantById: %v", err)
	}
	collection.Tenant = tenantToGraphQlTenant(tenant)
	sizeForAllObjectInstances, err := clientClerkHandler.GetSizeForAllObjectInstancesByCollectionId(ctx, &pb.Id{Id: collection.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetSizeForAllObjectInstances: %v", err)
	}
	collection.TotalObjectSizeForAllObjectInstances = float64(sizeForAllObjectInstances.Size)
	return collection, err
}

func GetObjectById(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, id string) (*model.Object, error) {
	objectPb, err := clientClerkHandler.GetObjectById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectById: %v", err)
	}
	object := objectToGraphQlObject(objectPb)
	status, err := clientClerkHandler.GetStatusForObjectId(ctx, &pb.Id{Id: object.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStatusForObjectId: %v", err)
	}
	object.Status = int(status.Size)
	collectionPb, err := clientClerkHandler.GetCollectionByIdFromMv(ctx, &pb.Id{Id: object.CollectionID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionByIdFromMv: %v", err)
	}
	collection := collectionToGraphQlCollection(collectionPb)
	object.Collection = collection
	return object, err
}

func GetObjectInstanceById(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, id string) (*model.ObjectInstance, error) {
	objectInstancePb, err := clientClerkHandler.GetObjectInstanceById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstanceById: %v", err)
	}
	objectInstance := objectInstanceToGraphQlObjectInstance(objectInstancePb)
	objectPb, err := clientClerkHandler.GetObjectById(ctx, &pb.Id{Id: objectInstance.ObjectID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectById: %v", err)
	}
	object := objectToGraphQlObject(objectPb)
	objectInstance.Object = object
	return objectInstance, err
}

func GetObjectInstanceCheckById(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, id string) (*model.ObjectInstanceCheck, error) {
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

func GetFileById(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, id string) (*model.File, error) {
	filePb, err := clientClerkHandler.GetFileById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetFileById: %v", err)
	}
	file := fileToGraphQlFile(filePb)
	objectPb, err := clientClerkHandler.GetObjectById(ctx, &pb.Id{Id: file.ObjectID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectById: %v", err)
	}
	object := objectToGraphQlObject(objectPb)
	file.Object = object
	return file, err
}

func GetStorageLocationById(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, id string) (*model.StorageLocation, error) {
	storageLocationPb, err := clientClerkHandler.GetStorageLocationById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStorageLocationById: %v", err)
	}
	storageLocation := storageLocationToGraphQlStorageLocation(storageLocationPb)
	amountOfObjects, err := clientClerkHandler.GetAmountOfObjectsForStorageLocationId(ctx, &pb.Id{Id: storageLocation.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetAmountOfObjectsForStorageLocationId: %v", err)
	}
	amountOfErrors, err := clientClerkHandler.GetAmountOfErrorsForStorageLocationId(ctx, &pb.Id{Id: storageLocation.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetAmountOfErrorsForStorageLocationId: %v", err)
	}
	storageLocation.AmountOfErrors = int(amountOfErrors.Size)
	storageLocation.AmountOfObjects = int(amountOfObjects.Size)
	tenantPb, err := clientClerkHandler.FindTenantById(ctx, &pb.Id{Id: storageLocation.TenantID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not FindTenantById: %v", err)
	}
	tenant := tenantToGraphQlTenant(tenantPb)
	storageLocation.Tenant = tenant
	return storageLocation, err
}

func GetStoragePartitionById(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, id string) (*model.StoragePartition, error) {
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

//Statistic

func GetMimeTypesForCollectionId(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, options *model.MimeTypeListOptions, allowedTenants []string) (*model.MimeTypeList, error) {
	keyCloakGroup, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if (len(tenantList) == 0) && (!slices.Contains(keyCloakGroup, "dlza-admin")) {
		return nil, errors.New("You are not allowed to retrieve datas")
	} else if len(tenantList) > 0 {
		for _, tenant := range tenantList {
			allowedTenants = append(allowedTenants, tenant.Id)
		}
		// allowedTenants = tenantList
	}
	if slices.Contains(keyCloakGroup, "dlza-admin") {
		allowedTenants = []string{}
	}
	optionsPb := pb.Pagination{
		Take:           10,
		SortDirection:  sortDirectionAscending,
		AllowedTenants: allowedTenants,
		SortKey:        "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 mime types")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.TenantID != nil {
			optionsPb.SecondId = *options.TenantID
		}
		if options.CollectionID != nil {
			optionsPb.Id = *options.CollectionID
		}
	}
	mimeTypesPb, err := clientClerkHandler.GetMimeTypesForCollectionId(ctx, &optionsPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetMimeTypesForCollectionId: %v", err)
	}

	mimeTypes := make([]*model.MimeType, 0)
	for _, mimeTypePb := range mimeTypesPb.MimeTypes {
		mimeType := model.MimeType{ID: mimeTypePb.Id, FileCount: int(mimeTypePb.FileCount), FilesSize: float64(mimeTypePb.FilesSize)}
		mimeTypes = append(mimeTypes, &mimeType)
	}
	return &model.MimeTypeList{Items: mimeTypes, TotalItems: int(mimeTypesPb.TotalItems)}, nil
}

func GetPronomsForCollectionId(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, options *model.PronomIDListOptions, allowedTenants []string) (*model.PronomIDList, error) {
	keyCloakGroup, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if (len(tenantList) == 0) && (!slices.Contains(keyCloakGroup, "dlza-admin")) {
		return nil, errors.New("You are not allowed to retrieve datas")
	} else if len(tenantList) > 0 {
		for _, tenant := range tenantList {
			allowedTenants = append(allowedTenants, tenant.Id)
		}
		// allowedTenants = tenantList
	}
	if slices.Contains(keyCloakGroup, "dlza-admin") {
		allowedTenants = []string{}
	}
	optionsPb := pb.Pagination{
		Take:           10,
		SortDirection:  sortDirectionAscending,
		AllowedTenants: allowedTenants,
		SortKey:        "ID",
	}
	if options != nil {
		if options.SortKey != nil {
			optionsPb.SortKey = toSnakeCase(options.SortKey.String())
		}
		if options.SortDirection != nil {
			if *options.SortDirection == model.SortDirectionDescending {
				optionsPb.SortDirection = sortDirectionDescending
			}
		}
		if options.Take != nil {
			if *options.Take > 1000 {
				return nil, errors.New("You could not retrieve more than 1000 pronoms")
			}
			optionsPb.Take = int32(*options.Take)
		}
		if options.Skip != nil {
			optionsPb.Skip = int32(*options.Skip)
		}
		if options.TenantID != nil {
			optionsPb.SecondId = *options.TenantID
		}
		if options.CollectionID != nil {
			optionsPb.Id = *options.CollectionID
		}
	}
	pronomsPb, err := clientClerkHandler.GetPronomsForCollectionId(ctx, &optionsPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetMimeTypesForCollectionId: %v", err)
	}

	pronoms := make([]*model.PronomID, 0)
	for _, pronomPb := range pronomsPb.Pronoms {
		pronom := model.PronomID{ID: pronomPb.Id, FileCount: int(pronomPb.FileCount), FilesSize: float64(pronomPb.FilesSize)}
		pronoms = append(pronoms, &pronom)
	}
	return &model.PronomIDList{Items: pronoms, TotalItems: int(pronomsPb.TotalItems)}, nil
}

func CreateCollection(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, input *model.CollectionInput) (*model.Collection, error) {
	_, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if len(tenantList) > 0 {
		for count, tenant := range tenantList {
			if tenant.Id == input.TenantID {
				if tenant.Create && tenant.Read && tenant.Delete {
					break
				} else {
					return nil, errors.New("You are not allowed to proceed with creating collection")
				}
			}
			if count == len(tenantList)-1 {
				return nil, errors.New("You are not allowed to proceed with creating collection")
			}
		}
	} else {
		return nil, errors.New("You are not allowed to proceed with creating collection")
	}
	collectionPb := collectionInputToGrpcCollection(*input)
	idPb, err := clientClerkHandler.CreateCollection(ctx, collectionPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not CreateCollection: %v", err)
	}
	collectionPb.Id = idPb.Id
	collectionG := collectionToGraphQlCollection(collectionPb)
	return collectionG, nil
}

func UpdateCollection(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, input *model.CollectionInput) (*model.Collection, error) {
	_, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if len(tenantList) > 0 {
		for count, tenant := range tenantList {
			if tenant.Id == input.TenantID {
				if tenant.Create && tenant.Read && tenant.Delete {
					break
				} else {
					return nil, errors.New("You are not allowed to proceed with updating collection")
				}
			}
			if count == len(tenantList)-1 {
				return nil, errors.New("You are not allowed to proceed with updating collection")
			}
		}
	} else {
		return nil, errors.New("You are not allowed to proceed with updating collection")
	}
	collectionPb := collectionInputToGrpcCollection(*input)
	_, err = clientClerkHandler.UpdateCollection(ctx, collectionPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not UpdateCollection: %v", err)
	}
	collectionG := collectionToGraphQlCollection(collectionPb)
	return collectionG, nil
}

func DeleteCollection(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, id string) (*model.Collection, error) {
	collectionPb, err := clientClerkHandler.GetCollectionById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionById: %v", err)
	}
	_, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if len(tenantList) > 0 {
		for count, tenant := range tenantList {
			if tenant.Id == collectionPb.TenantId {
				if tenant.Create && tenant.Read && tenant.Delete {
					break
				} else {
					return nil, errors.New("You are not allowed to proceed with deleting the collection")
				}
			}
			if count == len(tenantList)-1 {
				return nil, errors.New("You are not allowed to proceed with deleting the collection")
			}
		}
	} else {
		return nil, errors.New("You are not allowed to proceed with deleting the collection")
	}
	collection := collectionToGraphQlCollection(collectionPb)
	_, err = clientClerkHandler.DeleteCollectionById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not DeleteCollection: %v", err)
	}
	return collection, nil
}

func CreateStorageLocation(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, input *model.StorageLocationInput) (*model.StorageLocation, error) {
	_, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if len(tenantList) > 0 {
		for count, tenant := range tenantList {
			if tenant.Id == input.TenantID {
				if tenant.Create && tenant.Read && tenant.Delete {
					break
				} else {
					return nil, errors.New("You are not allowed to proceed with creating storage location")
				}
			}
			if count == len(tenantList)-1 {
				return nil, errors.New("You are not allowed to proceed with creating storage location")
			}
		}
	} else {
		return nil, errors.New("You are not allowed to proceed with creating storage location")
	}
	storageLocationPb := storageLocationInputToGrpcStorageLocation(input)
	idPb, err := clientClerkHandler.SaveStorageLocation(ctx, storageLocationPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not CreateStorageLocation: %v", err)
	}
	storageLocationPb.Id = idPb.Id
	storageLocationG := storageLocationToGraphQlStorageLocation(storageLocationPb)
	return storageLocationG, nil
}

func UpdateStorageLocation(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, input *model.StorageLocationInput) (*model.StorageLocation, error) {
	_, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if len(tenantList) > 0 {
		for count, tenant := range tenantList {
			if tenant.Id == input.TenantID {
				if tenant.Create && tenant.Read && tenant.Delete {
					break
				} else {
					return nil, errors.New("You are not allowed to proceed with updating storage location")
				}
			}
			if count == len(tenantList)-1 {
				return nil, errors.New("You are not allowed to proceed with updating storage location")
			}
		}
	} else {
		return nil, errors.New("You are not allowed to proceed with updating storage location")
	}
	storageLocationPb := storageLocationInputToGrpcStorageLocation(input)
	_, err = clientClerkHandler.UpdateStorageLocation(ctx, storageLocationPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not UpdateStorageLocation: %v", err)
	}
	storageLocationG := storageLocationToGraphQlStorageLocation(storageLocationPb)
	return storageLocationG, nil
}

func DeleteStorageLocation(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, id string) (*model.StorageLocation, error) {
	StorageLocationPb, err := clientClerkHandler.GetStorageLocationById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStorageLocationById: %v", err)
	}
	_, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if len(tenantList) > 0 {
		for count, tenant := range tenantList {
			if tenant.Id == StorageLocationPb.TenantId {
				if tenant.Create && tenant.Read && tenant.Delete {
					break
				} else {
					return nil, errors.New("You are not allowed to proceed with deleting the storage location")
				}
			}
			if count == len(tenantList)-1 {
				return nil, errors.New("You are not allowed to proceed with deleting the storage location")
			}
		}
	} else {
		return nil, errors.New("You are not allowed to proceed with deleting the storage location")
	}
	storageLocationPb, err := clientClerkHandler.GetStorageLocationById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStorageLocationById: %v", err)
	}
	_, err = clientClerkHandler.DeleteStorageLocationById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not DeleteStorageLocationById: %v", err)
	}
	storageLocationG := storageLocationToGraphQlStorageLocation(storageLocationPb)
	return storageLocationG, nil
}

func CreateStoragePartition(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, input *model.StoragePartitionInput) (*model.StoragePartition, error) {
	storageLocationPb, err := clientClerkHandler.GetStorageLocationById(ctx, &pb.Id{Id: input.StorageLocationID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStorageLocationById: %v", err)
	}
	_, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if len(tenantList) > 0 {
		for count, tenant := range tenantList {
			if tenant.Id == storageLocationPb.TenantId {
				if tenant.Create && tenant.Read && tenant.Delete {
					break
				} else {
					return nil, errors.New("You are not allowed to proceed with creating storage partition")
				}
			}
			if count == len(tenantList)-1 {
				return nil, errors.New("You are not allowed to proceed with updating storage partition")
			}
		}
	} else {
		return nil, errors.New("You are not allowed to proceed with creating storage partition")
	}
	storagePartitionPb := storagePartitionInputToGrpcStoragePartition(input)
	idPb, err := clientClerkHandler.CreateStoragePartition(ctx, storagePartitionPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not CreateStoragePartition: %v", err)
	}
	storagePartitionPb.Id = idPb.Id
	storagePartitionG := storagePartitionToGraphQlStoragePartition(storagePartitionPb)
	return storagePartitionG, nil
}

func UpdateStoragePartition(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, input *model.StoragePartitionInput) (*model.StoragePartition, error) {
	storageLocationPb, err := clientClerkHandler.GetStorageLocationById(ctx, &pb.Id{Id: input.StorageLocationID})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStorageLocationById: %v", err)
	}
	_, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if len(tenantList) > 0 {
		for count, tenant := range tenantList {
			if tenant.Id == storageLocationPb.TenantId {
				if tenant.Create && tenant.Read && tenant.Delete {
					break
				} else {
					return nil, errors.New("You are not allowed to proceed with updating storage partition")
				}
			}
			if count == len(tenantList)-1 {
				return nil, errors.New("You are not allowed to proceed with updating storage partition")
			}
		}
	} else {
		return nil, errors.New("You are not allowed to proceed with updating storage partition")
	}
	storagePartitionPb := storagePartitionInputToGrpcStoragePartition(input)
	_, err = clientClerkHandler.UpdateStoragePartition(ctx, storagePartitionPb)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not UpdateStoragePartition: %v", err)
	}
	storagePartitionG := storagePartitionToGraphQlStoragePartition(storagePartitionPb)
	return storagePartitionG, nil
}

func DeleteStoragePartition(ctx context.Context, clientClerkHandler pbHandler.ClerkHandlerServiceClient, id string) (*model.StoragePartition, error) {
	StoragePartitionPb, err := clientClerkHandler.GetStoragePartitionById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStoragePartitionById: %v", err)
	}
	StorageLocationPb, err := clientClerkHandler.GetStorageLocationById(ctx, &pb.Id{Id: StoragePartitionPb.StorageLocationId})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStorageLocationById: %v", err)
	}
	_, tenantList, err := middleware.TenantGroups(ctx)
	if err != nil {
		return nil, err
	}
	if len(tenantList) > 0 {
		for count, tenant := range tenantList {
			if tenant.Id == StorageLocationPb.TenantId {
				if tenant.Create && tenant.Read && tenant.Delete {
					break
				} else {
					return nil, errors.New("You are not allowed to proceed with deleting the storage partition")
				}
			}
			if count == len(tenantList)-1 {
				return nil, errors.New("You are not allowed to proceed with deleting the storage partition")
			}
		}
	} else {
		return nil, errors.New("You are not allowed to proceed with deleting the storage partition")
	}
	storagePartition, err := GetStoragePartitionById(ctx, clientClerkHandler, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStoragePartitionById: %v", err)
	}
	_, err = clientClerkHandler.DeleteStoragePartitionById(ctx, &pb.Id{Id: id})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not DeleteStoragePartitionById: %v", err)
	}
	return storagePartition, nil
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
	collection.TotalFileSize = float64(collectionPb.TotalFileSize)
	collection.TotalFileCount = int(collectionPb.TotalFileCount)
	collection.TotalObjectCount = int(collectionPb.TotalObjectCount)
	return &collection
}

func collectionInputToGrpcCollection(collection model.CollectionInput) *pb.Collection {
	var collectionPb pb.Collection
	collectionPb.Id = collection.ID
	collectionPb.Alias = collection.Alias
	collectionPb.Description = collection.Description
	collectionPb.Owner = collection.Owner
	collectionPb.Name = collection.Name
	collectionPb.OwnerMail = collection.OwnerMail
	collectionPb.Quality = int32(collection.Quality)
	collectionPb.TenantId = collection.TenantID
	return &collectionPb
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
	object.Size = float64(objectPb.Size)
	object.ID = objectPb.Id
	object.CollectionID = objectPb.CollectionId
	object.Checksum = objectPb.Checksum
	object.TotalFileSize = float64(objectPb.TotalFileSize)
	object.TotalFileCount = int(objectPb.TotalFileCount)
	object.Authors = objectPb.Authors
	object.Holding = objectPb.Holding
	object.Expiration = objectPb.Expiration
	object.Head = objectPb.Head
	var versionsMap models.Versions
	err := json.Unmarshal([]byte(objectPb.Versions), &versionsMap)
	if err != nil {
		object.Versions = objectPb.Versions
	} else {
		var versionsString string
		for _, version := range maps.Keys(versionsMap) {
			var coma string
			var created string
			if versionsMap[version].Created != "" {
				created = " : " + versionsMap[version].Created
			}
			if versionsString != "" {
				coma = ", "
			}
			versionsString = versionsString + version + created + coma
		}
		object.Versions = versionsString
	}
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
	file.MimeType = filePb.MimeType
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
	storageLocation.Connection = "xxxxxxxxxxxxxx"
	storageLocation.Quality = int(storageLocationPb.Quality)
	storageLocation.Price = int(storageLocationPb.Price)
	storageLocation.SecurityCompliency = storageLocationPb.SecurityCompliency
	storageLocation.FillFirst = storageLocationPb.FillFirst
	storageLocation.OcflType = storageLocationPb.OcflType
	storageLocation.NumberOfThreads = int(storageLocationPb.NumberOfThreads)
	storageLocation.TenantID = storageLocationPb.TenantId
	storageLocation.TotalExistingVolume = float64(storageLocationPb.TotalExistingVolume)
	storageLocation.TotalFilesSize = float64(storageLocationPb.TotalFilesSize)
	return &storageLocation
}

func storageLocationInputToGrpcStorageLocation(storageLocationInput *model.StorageLocationInput) *pb.StorageLocation {
	var storageLocationPb pb.StorageLocation
	storageLocationPb.Id = storageLocationInput.ID
	storageLocationPb.Alias = storageLocationInput.Alias
	storageLocationPb.Type = storageLocationInput.Type
	storageLocationPb.Vault = storageLocationInput.Vault
	storageLocationPb.Connection = storageLocationInput.Connection
	storageLocationPb.Quality = int32(storageLocationInput.Quality)
	storageLocationPb.Price = int32(storageLocationInput.Price)
	storageLocationPb.SecurityCompliency = storageLocationInput.SecurityCompliency
	storageLocationPb.FillFirst = storageLocationInput.FillFirst
	storageLocationPb.OcflType = storageLocationInput.OcflType
	storageLocationPb.NumberOfThreads = int32(storageLocationInput.NumberOfThreads)
	storageLocationPb.TenantId = storageLocationInput.TenantID
	return &storageLocationPb
}

func storagePartitionInputToGrpcStoragePartition(storagePartitionInput *model.StoragePartitionInput) *pb.StoragePartition {
	var storagePartitionPb pb.StoragePartition
	storagePartitionPb.Alias = storagePartitionInput.Alias
	storagePartitionPb.Name = storagePartitionInput.Name
	storagePartitionPb.MaxSize = int64(storagePartitionInput.MaxSize)
	storagePartitionPb.MaxObjects = int64(storagePartitionInput.MaxObjects)
	storagePartitionPb.CurrentSize = int64(storagePartitionInput.CurrentSize)
	storagePartitionPb.CurrentObjects = int64(storagePartitionInput.CurrentObjects)
	storagePartitionPb.Id = storagePartitionInput.ID
	storagePartitionPb.StorageLocationId = storagePartitionInput.StorageLocationID
	return &storagePartitionPb
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

func toSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
