package port

import (
	"context"
	"generative-xyz-search-engine/pkg/driver/mongodb"

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
}

type IUserRepository interface {
	mongodb.Repository
}
