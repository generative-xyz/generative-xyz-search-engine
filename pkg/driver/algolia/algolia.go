package algolia

import (
	"generative-xyz-search-engine/pkg/logger"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type AlgoliaFilter struct {
	Page       int
	Limit      int
	SearchStr  string
	FilterStrs []string
}

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

func (al *GenerativeAlgolia) Search(indexName string, builder *AlgoliaFilter) (search.QueryRes, error) {
	index := al.Client.InitIndex(indexName)
	if builder.Limit == 0 {
		builder.Limit = 10
	}

	if builder.Page == 0 {
		builder.Page = 1
	}

	opts := []interface{}{
		opt.Page(builder.Page - 1),
		opt.HitsPerPage(builder.Limit),
		opt.AttributesToRetrieve("*"),
	}

	if len(builder.FilterStrs) > 0 {
		// opts = append(opts, opt.TypoTolerance(false))
		// opts = append(opts, opt.RestrictSearchableAttributes(builder.SearchField))
		for _, s := range builder.FilterStrs {
			opts = append(opts, opt.Filters(s))
		}
	}

	res, err := index.Search(builder.SearchStr, opts...)
	return res, err
}
