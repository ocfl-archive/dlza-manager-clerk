package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ocfl-archive/dlza-manager-clerk/graph/model"
	"github.com/ocfl-archive/dlza-manager-clerk/middleware"
	"github.com/ocfl-archive/dlza-manager-clerk/service"
)

// Objects is the resolver for the objects field.
func (r *collectionResolver) Objects(ctx context.Context, obj *model.Collection, options *model.ObjectListOptions) (*model.ObjectList, error) {
	collections, err := service.GetObjectsForCollection(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetObjectsForCollection: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return collections, nil
}

// Files is the resolver for the files field.
func (r *collectionResolver) Files(ctx context.Context, obj *model.Collection, options *model.FileListOptions) (*model.FileList, error) {
	files, err := service.GetFilesForCollection(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetFilesForCollection: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return files, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, code string) (*model.User, error) {
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		fmt.Println("gc err", err, "\n", ctx)
		return nil, middleware.GraphqlErrorWrapper(err, ctx, http.StatusInternalServerError)
	}

	err = middleware.Callback(ctx, gc, code)
	if err != nil {
		fmt.Println("Callback err", err, "\n", ctx)
		return nil, middleware.GraphqlErrorWrapper(err, ctx, http.StatusInternalServerError)
	}
	userClaim, err := middleware.GetUser(gc)
	if err != nil {
		fmt.Println("GetUser err", err, "\n", ctx)
		middleware.GraphqlErrorWrapper(err, ctx, http.StatusUnauthorized)
	}
	user := model.User{
		Username: userClaim.PreferredUsername,
		Email:    userClaim.Email,
		ID:       userClaim.Sub,
	}

	return &user, nil
}

// Logout is the resolver for the logout field.
func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}
	middleware.ResetSession(gc)
	return true, nil
}

// ObjectInstances is the resolver for the objectInstances field.
func (r *objectResolver) ObjectInstances(ctx context.Context, obj *model.Object, options *model.ObjectInstanceListOptions) (*model.ObjectInstanceList, error) {
	objectInstances, err := service.GetObjectInstancesForObject(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetObjectInstancesForObject: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return objectInstances, nil
}

// Files is the resolver for the files field.
func (r *objectResolver) Files(ctx context.Context, obj *model.Object, options *model.FileListOptions) (*model.FileList, error) {
	files, err := service.GetFilesForObject(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetFilesForObject: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return files, nil
}

// ObjectInstanceChecks is the resolver for the objectInstanceChecks field.
func (r *objectInstanceResolver) ObjectInstanceChecks(ctx context.Context, obj *model.ObjectInstance, options *model.ObjectInstanceCheckListOptions) (*model.ObjectInstanceCheckList, error) {
	objectInstanceChecks, err := service.GetObjectInstanceChecksForObjectInstance(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetObjectInstanceChecksForObjectInstance: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return objectInstanceChecks, nil
}

// Auth is the resolver for the auth field.
func (r *queryResolver) Auth(ctx context.Context) (*model.Auth, error) {
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(err, ctx, http.StatusInternalServerError)
	}

	authCodeUrl, err := middleware.GetAuthCodeURL(gc)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(err, ctx, http.StatusInternalServerError)
	}
	return &model.Auth{
		AuthCodeURL: authCodeUrl,
	}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	//TODO do something to retrieve username and user email from gincontexts
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(err, ctx, http.StatusInternalServerError)
	}

	userClaim, err := middleware.GetUser(gc)
	if err != nil {
		return nil, nil
	}
	user := model.User{
		Username: userClaim.PreferredUsername,
		Email:    userClaim.Email,
		ID:       userClaim.Sub,
	}
	return &user, nil
}

// Tenants is the resolver for the tenants field.
func (r *queryResolver) Tenants(ctx context.Context, options *model.TenantListOptions) (*model.TenantList, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	tenants, err := service.GetTenants(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not FindAllTenants: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return tenants, nil
}

// Tenant is the resolver for the tenant field.
func (r *queryResolver) Tenant(ctx context.Context, id string) (*model.Tenant, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	tenant, err := service.GetTenantById(ctx, r.ClientClerkHandler, id, r.AllowedTenants)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetTenantById: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return tenant, nil
}

// Collections is the resolver for the collections field.
func (r *queryResolver) Collections(ctx context.Context, options *model.CollectionListOptions) (*model.CollectionList, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	collections, err := service.GetCollectionsForTenantId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetCollectionsForTenantId: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return collections, nil
}

// Collection is the resolver for the collection field.
func (r *queryResolver) Collection(ctx context.Context, id string) (*model.Collection, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	collection, err := service.GetCollectionById(ctx, r.ClientClerkHandler, id)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetCollectionById: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return collection, nil
}

// Objects is the resolver for the objects field.
func (r *queryResolver) Objects(ctx context.Context, options *model.ObjectListOptions) (*model.ObjectList, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	objects, err := service.GetObjectsForCollectionId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetObjectsForCollectionId: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return objects, nil
}

// Object is the resolver for the object field.
func (r *queryResolver) Object(ctx context.Context, id string) (*model.Object, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	object, err := service.GetObjectById(ctx, r.ClientClerkHandler, id)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetObjectById: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return object, nil
}

// ObjectInstances is the resolver for the objectInstances field.
func (r *queryResolver) ObjectInstances(ctx context.Context, options *model.ObjectInstanceListOptions) (*model.ObjectInstanceList, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	objectInstances, err := service.GetObjectInstancesForObjectId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetObjectInstancesForObjectId: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return objectInstances, nil
}

// ObjectInstance is the resolver for the objectInstance field.
func (r *queryResolver) ObjectInstance(ctx context.Context, id string) (*model.ObjectInstance, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	objectInstance, err := service.GetObjectInstanceById(ctx, r.ClientClerkHandler, id)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetObjectInstanceById: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return objectInstance, nil
}

// ObjectInstanceChecks is the resolver for the objectInstanceChecks field.
func (r *queryResolver) ObjectInstanceChecks(ctx context.Context, options *model.ObjectInstanceCheckListOptions) (*model.ObjectInstanceCheckList, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	objectInstanceChecks, err := service.GetObjectInstanceChecksForObjectInstanceId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetObjectInstanceChecksForObjectInstanceId: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return objectInstanceChecks, nil
}

// ObjectInstanceCheck is the resolver for the objectInstanceCheck field.
func (r *queryResolver) ObjectInstanceCheck(ctx context.Context, id string) (*model.ObjectInstanceCheck, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	objectInstanceCheck, err := service.GetObjectInstanceCheckById(ctx, r.ClientClerkHandler, id)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetObjectInstanceCheckById: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return objectInstanceCheck, nil
}

// Files is the resolver for the files field.
func (r *queryResolver) Files(ctx context.Context, options *model.FileListOptions) (*model.FileList, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	files, err := service.GetFilesForObjectId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetFilesForObjectId: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return files, nil
}

// File is the resolver for the file field.
func (r *queryResolver) File(ctx context.Context, id string) (*model.File, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	file, err := service.GetFileById(ctx, r.ClientClerkHandler, id)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetFileById: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return file, nil
}

// StorageLocations is the resolver for the storageLocations field.
func (r *queryResolver) StorageLocations(ctx context.Context, options *model.StorageLocationListOptions) (*model.StorageLocationList, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	storageLocations, err := service.GetStorageLocationsForTenantId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetStorageLocationsForTenantId: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return storageLocations, nil
}

// StorageLocation is the resolver for the storageLocation field.
func (r *queryResolver) StorageLocation(ctx context.Context, id string) (*model.StorageLocation, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	storageLocation, err := service.GetStorageLocationById(ctx, r.ClientClerkHandler, id)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetStorageLocationById: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return storageLocation, nil
}

// StoragePartitions is the resolver for the storagePartitions field.
func (r *queryResolver) StoragePartitions(ctx context.Context, options *model.StoragePartitionListOptions) (*model.StoragePartitionList, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	storagePartitions, err := service.GetStoragePartitionsForLocationId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetStoragePartitionsForLocationId: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return storagePartitions, nil
}

// StoragePartition is the resolver for the storagePartition field.
func (r *queryResolver) StoragePartition(ctx context.Context, id string) (*model.StoragePartition, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	storagePartition, err := service.GetStoragePartitionById(ctx, r.ClientClerkHandler, id)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetStoragePartitionById: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return storagePartition, nil
}

// MimeTypes is the resolver for the mimeTypes field.
func (r *queryResolver) MimeTypes(ctx context.Context, options *model.MimeTypeListOptions) (*model.MimeTypeList, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	mimeTypes, err := service.GetMimeTypesForCollectionId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetMimeTypesForCollectionId: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return mimeTypes, nil
}

// PronomIds is the resolver for the pronomIds field.
func (r *queryResolver) PronomIds(ctx context.Context, options *model.PronomIDListOptions) (*model.PronomIDList, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	pronoms, err := service.GetPronomsForCollectionId(ctx, r.ClientClerkHandler, options, r.AllowedTenants)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetPronomsForCollectionId: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return pronoms, nil
}

// StoragePartitions is the resolver for the storagePartitions field.
func (r *storageLocationResolver) StoragePartitions(ctx context.Context, obj *model.StorageLocation, options *model.StoragePartitionListOptions) (*model.StoragePartitionList, error) {
	storagePartitions, err := service.GetStoragePartitionsForLocation(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetStoragePartitionsForLocation: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return storagePartitions, nil
}

// ObjectInstances is the resolver for the objectInstances field.
func (r *storagePartitionResolver) ObjectInstances(ctx context.Context, obj *model.StoragePartition, options *model.ObjectInstanceListOptions) (*model.ObjectInstanceList, error) {
	objectInstances, err := service.GetObjectInstancesForStoragePartition(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetObjectInstancesForStoragePartition: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return objectInstances, nil
}

// Collections is the resolver for the collections field.
func (r *tenantResolver) Collections(ctx context.Context, obj *model.Tenant, options *model.CollectionListOptions) (*model.CollectionList, error) {
	collections, err := service.GetCollectionsForTenant(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetCollectionsForTenant: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return collections, nil
}

// StorageLocations is the resolver for the storageLocations field.
func (r *tenantResolver) StorageLocations(ctx context.Context, obj *model.Tenant, options *model.StorageLocationListOptions) (*model.StorageLocationList, error) {
	storageLocations, err := service.GetStorageLocationsForTenant(ctx, r.ClientClerkHandler, obj, options)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not GetStorageLocationsForTenant: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return storageLocations, nil
}

// Tenants is the resolver for the tenants field.
func (r *userResolver) Tenants(ctx context.Context, obj *model.User) ([]*model.Tenant, error) {
	if errM := middleware.GraphqlVerifyToken(ctx); errM != nil {
		return nil, middleware.GraphqlErrorWrapper(errM, ctx, http.StatusUnauthorized)
	}
	tenants, err := service.GetTenants(ctx, r.ClientClerkHandler, nil, r.AllowedTenants)
	if err != nil {
		return nil, middleware.GraphqlErrorWrapper(errors.New("Could not FindAllTenants: "+err.Error()), ctx, http.StatusInternalServerError)
	}
	return tenants.Items, err
}

// Collection returns CollectionResolver implementation.
func (r *Resolver) Collection() CollectionResolver { return &collectionResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

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

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type collectionResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type objectResolver struct{ *Resolver }
type objectInstanceResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type storageLocationResolver struct{ *Resolver }
type storagePartitionResolver struct{ *Resolver }
type tenantResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
