package algolia

import (
	"generative-xyz-search-engine/pkg/logger"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type GenerativeAlgolia struct {
	Client *search.Client
}

func NewAlgoliaClient() *GenerativeAlgolia {
	client := search.NewClient(
		viper.GetString("ALGOLIA_APPLICATION_ID"), viper.GetString("ALGOLIA_API_KEY"),
	)
	return &GenerativeAlgolia{Client: client}
}

func (al *GenerativeAlgolia) BulkIndexer(indexName string, objects interface{}) {
	index := al.Client.InitIndex(indexName)
	// Add objects to the index
	_, err := index.SaveObjects(objects)
	if err != nil {
		logger.AtLog.Logger.Error(err.Error(), zap.Error(err))
	}
}

func (al *GenerativeAlgolia) DeleteObject(indexName string, id string) {
	index := al.Client.InitIndex(indexName)
	// Add objects to the index
	_, err := index.DeleteObject(id)
	if err != nil {
		logger.AtLog.Logger.Error(err.Error(), zap.Error(err))
	}
}
