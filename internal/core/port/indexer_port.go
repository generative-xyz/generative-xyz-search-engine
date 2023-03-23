package port

import (
	"context"
	"generative-xyz-search-engine/pkg/driver/mongodb"
	"generative-xyz-search-engine/pkg/model"

	"go.mongodb.org/mongo-driver/bson"
)

type IIndexerUsecase interface {
	Schedule()
	ProcessIndexDataAlgolia(context.Context, bool) error
}

type IProjectRepository interface {
	mongodb.Repository
	SelectedProjectFields() bson.M
}

type ITokenUriRepository interface {
	mongodb.Repository
	ProjectGetCurrentListingNumber(projectID string) (uint64, error)
	ProjectGetMintVolume(projectID string) (uint64, error)
}

type IUserRepository interface {
	mongodb.Repository
}

type IDexBtcListingRepository interface {
	mongodb.Repository
	ProjectGetListingVolume(projectID string) (uint64, error)
	ProjectGetCEXVolume(projectID string) (uint64, error)
	RetrieveFloorPriceOfCollection(collectionID string) (uint64, error)
	AggregateBTCVolumn() ([]model.AggregateProjectItemResp, error)
}
