package mongo

import (
	"generative-xyz-search-engine/internal/core/port"
	"generative-xyz-search-engine/pkg/driver/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

var _ port.IProjectRepository = (*projectRepository)(nil)

type projectRepository struct {
	mongodb.BaseRepository
}

func NewProjectRepository(db *mongo.Database) port.IProjectRepository {
	return &projectRepository{
		BaseRepository: mongodb.BaseRepository{
			CollectionName: "projects",
			DB:             db,
		},
	}
}
