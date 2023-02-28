package mongo

import (
	"generative-xyz-search-engine/internal/core/port"
	"generative-xyz-search-engine/pkg/driver/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

var _ port.IUserRepository = (*tokenRepository)(nil)

type userRepository struct {
	mongodb.BaseRepository
}

func NewUserRepository(db *mongo.Database) port.IUserRepository {
	return &userRepository{
		BaseRepository: mongodb.BaseRepository{
			CollectionName: "users",
			DB:             db,
		},
	}
}
