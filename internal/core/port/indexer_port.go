package port

import (
	"context"
	"generative-xyz-search-engine/pkg/driver/mongodb"
)

type IIndexerUsecase interface {
	Schedule()
	ProcessIndexDataAlgolia(context.Context, bool) error
}

type IProjectRepository interface {
	mongodb.Repository
}

type ITokenUriRepository interface {
	mongodb.Repository
}

type IUserRepository interface {
	mongodb.Repository
}
