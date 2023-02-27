package mongo

import (
	"generative-xyz-search-engine/internal/core/port"
	"generative-xyz-search-engine/pkg/driver/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

var _ port.ITokenUriRepository = (*tokenRepository)(nil)

type tokenRepository struct {
	mongodb.BaseRepository
}

func NewTokenRepository(db *mongo.Database) port.ITokenUriRepository {
	return &tokenRepository{
		BaseRepository: mongodb.BaseRepository{
			CollectionName: "token_uri",
			DB:             db,
		},
	}
}
