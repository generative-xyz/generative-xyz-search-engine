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
	ProjectGetCurrentListingNumber() ([]*model.TokenUriListingPage, error)
	ProjectGetMintVolume() ([]*model.TokenUriListingVolume, error)
}

type IUserRepository interface {
	mongodb.Repository
}

type IDexBtcListingRepository interface {
	mongodb.Repository
	ProjectGetListingVolume() ([]*model.TokenUriListingVolume, error)
	ProjectGetCEXVolume() ([]*model.TokenUriListingVolume, error)
	RetrieveFloorPriceOfCollection() ([]*model.MarketplaceBTCListingFloorPrice, error)
	AggregateBTCVolumn() ([]model.AggregateProjectItemResp, error)
}
