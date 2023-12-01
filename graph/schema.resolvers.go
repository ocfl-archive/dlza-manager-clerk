package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"emperror.dev/errors"

	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/graph/model"
	"gitlab.switch.ch/ub-unibas/dlza/microservices/dlza-manager-clerk/service"
)

// Objects is the resolver for the objects field.
func (r *collectionResolver) Objects(ctx context.Context, obj *model.Collection, options *model.ObjectListOptions) (*model.ObjectList, error) {
	collections, err := service.GetObjectsForCollection(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectsForCollection: %v", err)
	}
	return collections, nil
}

// Files is the resolver for the files field.
func (r *collectionResolver) Files(ctx context.Context, obj *model.Collection, options *model.FileListOptions) (*model.FileList, error) {
	files, err := service.GetFilesForCollection(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetFilesForCollection: %v", err)
	}
	return files, nil
}

// ObjectInstances is the resolver for the objectInstances field.
func (r *objectResolver) ObjectInstances(ctx context.Context, obj *model.Object, options *model.ObjectInstanceListOptions) (*model.ObjectInstanceList, error) {
	objectInstances, err := service.GetObjectInstancesForObject(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstancesForObject: %v", err)
	}
	return objectInstances, nil
}

// Files is the resolver for the files field.
func (r *objectResolver) Files(ctx context.Context, obj *model.Object, options *model.FileListOptions) (*model.FileList, error) {
	files, err := service.GetFilesForObject(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetFilesForObject: %v", err)
	}
	return files, nil
}

// ObjectInstanceChecks is the resolver for the objectInstanceChecks field.
func (r *objectInstanceResolver) ObjectInstanceChecks(ctx context.Context, obj *model.ObjectInstance, options *model.ObjectInstanceCheckListOptions) (*model.ObjectInstanceCheckList, error) {
	objectInstanceChecks, err := service.GetObjectInstanceChecksForObjectInstance(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstanceChecksForObjectInstance: %v", err)
	}
	return objectInstanceChecks, nil
}

// Tenants is the resolver for the tenants field.
func (r *queryResolver) Tenants(ctx context.Context, options *model.TenantListOptions) (*model.TenantList, error) {
	tenants, err := service.GetTenants(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not FindAllTenants: %v", err)
	}
	return tenants, nil
}

// Tenant is the resolver for the tenant field.
func (r *queryResolver) Tenant(ctx context.Context, id string) (*model.Tenant, error) {
	tenant, err := service.GetTenantById(ctx, r.ClientClerkHandler, id, r.AllowedTenants)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetTenantById: %v", err)
	}
	return tenant, nil
}

// Collections is the resolver for the collections field.
func (r *queryResolver) Collections(ctx context.Context, options *model.CollectionListOptions) (*model.CollectionList, error) {
	collections, err := service.GetCollectionsForTenantId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionsForTenantId: %v", err)
	}
	return collections, nil
}

// Collection is the resolver for the collection field.
func (r *queryResolver) Collection(ctx context.Context, id string) (*model.Collection, error) {
	collection, err := service.GetCollectionById(ctx, r.ClientClerkHandler, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionById: %v", err)
	}
	return collection, nil
}

// Objects is the resolver for the objects field.
func (r *queryResolver) Objects(ctx context.Context, options *model.ObjectListOptions) (*model.ObjectList, error) {
	objects, err := service.GetObjectsForCollectionId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectsForCollectionId: %v", err)
	}
	return objects, nil
}

// Object is the resolver for the object field.
func (r *queryResolver) Object(ctx context.Context, id string) (*model.Object, error) {
	object, err := service.GetObjectById(ctx, r.ClientClerkHandler, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectById: %v", err)
	}
	return object, nil
}

// ObjectInstances is the resolver for the objectInstances field.
func (r *queryResolver) ObjectInstances(ctx context.Context, options *model.ObjectInstanceListOptions) (*model.ObjectInstanceList, error) {
	objectInstances, err := service.GetObjectInstancesForObjectId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstancesForObjectId: %v", err)
	}
	return objectInstances, nil
}

// ObjectInstance is the resolver for the objectInstance field.
func (r *queryResolver) ObjectInstance(ctx context.Context, id string) (*model.ObjectInstance, error) {
	objectInstance, err := service.GetObjectInstanceById(ctx, r.ClientClerkHandler, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstanceById: %v", err)
	}
	return objectInstance, nil
}

// ObjectInstanceChecks is the resolver for the objectInstanceChecks field.
func (r *queryResolver) ObjectInstanceChecks(ctx context.Context, options *model.ObjectInstanceCheckListOptions) (*model.ObjectInstanceCheckList, error) {
	objectInstanceChecks, err := service.GetObjectInstanceChecksForObjectInstanceId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstanceChecksForObjectInstance: %v", err)
	}
	return objectInstanceChecks, nil
}

// ObjectInstanceCheck is the resolver for the objectInstanceCheck field.
func (r *queryResolver) ObjectInstanceCheck(ctx context.Context, id string) (*model.ObjectInstanceCheck, error) {
	objectInstanceCheck, err := service.GetObjectInstanceCheckById(ctx, r.ClientClerkHandler, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstanceCheckById: %v", err)
	}
	return objectInstanceCheck, nil
}

// Files is the resolver for the files field.
func (r *queryResolver) Files(ctx context.Context, options *model.FileListOptions) (*model.FileList, error) {
	files, err := service.GetFilesForObjectId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetFilesForObjectId: %v", err)
	}
	return files, nil
}

// File is the resolver for the file field.
func (r *queryResolver) File(ctx context.Context, id string) (*model.File, error) {
	file, err := service.GetFileById(ctx, r.ClientClerkHandler, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetFileById: %v", err)
	}
	return file, nil
}

// StorageLocations is the resolver for the storageLocations field.
func (r *queryResolver) StorageLocations(ctx context.Context, options *model.StorageLocationListOptions) (*model.StorageLocationList, error) {
	storageLocations, err := service.GetStorageLocationsForTenantId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStorageLocationsForTenantId: %v", err)
	}
	return storageLocations, nil
}

// StorageLocation is the resolver for the storageLocation field.
func (r *queryResolver) StorageLocation(ctx context.Context, id string) (*model.StorageLocation, error) {
	storageLocation, err := service.GetStorageLocationById(ctx, r.ClientClerkHandler, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStorageLocationById: %v", err)
	}
	return storageLocation, nil
}

// StoragePartitions is the resolver for the storagePartitions field.
func (r *queryResolver) StoragePartitions(ctx context.Context, options *model.StoragePartitionListOptions) (*model.StoragePartitionList, error) {
	storagePartitions, err := service.GetStoragePartitionsForLocationId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStoragePartitionsForTenantId: %v", err)
	}
	return storagePartitions, nil
}

// StoragePartition is the resolver for the storagePartition field.
func (r *queryResolver) StoragePartition(ctx context.Context, id string) (*model.StoragePartition, error) {
	storagePartition, err := service.GetStoragePartitionById(ctx, r.ClientClerkHandler, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStoragePartitionById: %v", err)
	}
	return storagePartition, nil
}

// MimeTypes is the resolver for the mimeTypes field.
func (r *queryResolver) MimeTypes(ctx context.Context, options *model.MimeTypeListOptions) (*model.MimeTypeList, error) {
	mimeTypes, err := service.GetMimeTypesForCollectionId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetMimeTypesForCollectionId: %v", err)
	}
	return mimeTypes, nil
}

// PronomIds is the resolver for the pronomIds field.
func (r *queryResolver) PronomIds(ctx context.Context, options *model.PronomIDListOptions) (*model.PronomIDList, error) {
	pronoms, err := service.GetPronomsForCollectionId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetPronomsForCollectionId: %v", err)
	}
	return pronoms, nil
}

// StoragePartitions is the resolver for the storagePartitions field.
func (r *storageLocationResolver) StoragePartitions(ctx context.Context, obj *model.StorageLocation, options *model.StoragePartitionListOptions) (*model.StoragePartitionList, error) {
	storagePartitions, err := service.GetStoragePartitionsForLocation(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStorageLocationsForTenant: %v", err)
	}
	return storagePartitions, nil
}

// ObjectInstances is the resolver for the objectInstances field.
func (r *storagePartitionResolver) ObjectInstances(ctx context.Context, obj *model.StoragePartition, options *model.ObjectInstanceListOptions) (*model.ObjectInstanceList, error) {
	objectInstances, err := service.GetObjectInstancesForStoragePartition(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetObjectInstancesForStoragePartition: %v", err)
	}
	return objectInstances, nil
}

// Collections is the resolver for the collections field.
func (r *tenantResolver) Collections(ctx context.Context, obj *model.Tenant, options *model.CollectionListOptions) (*model.CollectionList, error) {
	collections, err := service.GetCollectionsForTenant(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetCollectionsForTenant: %v", err)
	}
	return collections, nil
}

// StorageLocations is the resolver for the storageLocations field.
func (r *tenantResolver) StorageLocations(ctx context.Context, obj *model.Tenant, options *model.StorageLocationListOptions) (*model.StorageLocationList, error) {
	storageLocations, err := service.GetStorageLocationsForTenant(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not GetStorageLocationsForTenant: %v", err)
	}
	return storageLocations, nil
}

// Collection returns CollectionResolver implementation.
func (r *Resolver) Collection() CollectionResolver { return &collectionResolver{r} }

// Object returns ObjectResolver implementation.
func (r *Resolver) Object() ObjectResolver { return &objectResolver{r} }

// ObjectInstance returns ObjectInstanceResolver implementation.
func (r *Resolver) ObjectInstance() ObjectInstanceResolver { return &objectInstanceResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// StorageLocation returns StorageLocationResolver implementation.
func (r *Resolver) StorageLocation() StorageLocationResolver { return &storageLocationResolver{r} }

// StoragePartition returns StoragePartitionResolver implementation.
func (r *Resolver) StoragePartition() StoragePartitionResolver { return &storagePartitionResolver{r} }

// Tenant returns TenantResolver implementation.
func (r *Resolver) Tenant() TenantResolver { return &tenantResolver{r} }

type collectionResolver struct{ *Resolver }
type objectResolver struct{ *Resolver }
type objectInstanceResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type storageLocationResolver struct{ *Resolver }
type storagePartitionResolver struct{ *Resolver }
type tenantResolver struct{ *Resolver }
