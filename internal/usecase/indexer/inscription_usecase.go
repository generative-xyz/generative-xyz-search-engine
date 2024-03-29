package usecase

import (
	"context"
	"fmt"
	"generative-xyz-search-engine/pkg/driver/algolia"
	"generative-xyz-search-engine/pkg/logger"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ListInscriptionResponse struct {
	Inscriptions []string
	Prev         int64
	Next         int64
}

type InscriptionDetail struct {
	ObjectID      string                 `json:"objectID"`
	Chain         string                 `json:"chain"`
	GenesisFee    int64                  `json:"genesis_fee"`
	GenesisHeight int64                  `json:"genesis_height"`
	Address       string                 `json:"address"`
	ContentType   string                 `json:"content_type"`
	InscriptionId string                 `json:"inscription_id"`
	Next          string                 `json:"next"`
	Number        int64                  `json:"number"`
	Output        map[string]interface{} `json:"output"`
	Previous      string                 `json:"previous"`
	Sat           int64                  `json:"sat"`
	Satpoint      string                 `json:"satpoint"`
	Timestamp     string                 `json:"timestamp"`
}

func (uc *indexerUsecase) fixInscriptionData(ctx context.Context) error {
	var inscriptions []*InscriptionDetail
	filter := algolia.AlgoliaFilter{
		Page: 0, Limit: 500,
		FilterStrs: []string{"sat=0"},
	}
	resp, err := uc.algoliaClient.Search("inscriptions", &filter)
	if err != nil {
		logger.AtLog.Error(err)
		return err
	}
	resp.UnmarshalHits(&inscriptions)

	// data := []*InscriptionDetail{}

	// client := resty.New()
	for _, i := range inscriptions {
		uc.algoliaClient.DeleteObject("inscriptions", i.ObjectID)
		// resp := &InscriptionDetail{}
		// _, err := client.R().
		// 	EnableTrace().
		// 	SetResult(&resp).
		// 	Get(fmt.Sprintf("%s/inscription/%s", viper.GetString("GENERATIVE_EXPLORER_API"), i.InscriptionId))
		// if err != nil {
		// 	logger.AtLog.Logger.Error("Get list inscriptions error", zap.Error(err))
		// }
		// if i.InscriptionId == "" {
		// 	continue
		// }

		// resp.ObjectID = i.InscriptionId
		// data = append(data, resp)
	}
	// uc.algoliaClient.BulkIndexer("inscriptions", data)
	return nil
}

func (uc *indexerUsecase) inscriptionIndexingData(ctx context.Context, isDelta bool) error {
	logger.AtLog.Infof("START inscriptionIndexingData algolia data %v", time.Now())
	// Create a Resty Client
	client := resty.New()
	index := int(0)

	if err := uc.redis.Get(ctx, "Inscription_Index_Count", &index); err != nil {
		uc.redis.Set(ctx, "Inscription_Index_Count", 260000, time.Duration(time.Hour*1))
		index = 260000
	}

	for {
		result := &ListInscriptionResponse{}
		_, err := client.R().
			EnableTrace().
			SetResult(result).
			Get(fmt.Sprintf("%s/inscriptions/%d", viper.GetString("GENERATIVE_EXPLORER_API"), index))

		if err != nil {
			logger.AtLog.Logger.Error("Get detail inscription error", zap.Error(err))
			return err
		}

		if result.Next == 0 {
			break
		}

		func(res *ListInscriptionResponse) {
			data := []*InscriptionDetail{}
			wg := &sync.WaitGroup{}
			for _, r := range res.Inscriptions {
				wg.Add(1)
				go func(id string) {
					defer wg.Done()
					resp := &InscriptionDetail{}
					_, err := client.R().
						EnableTrace().
						SetResult(&resp).
						Get(fmt.Sprintf("%s/inscription/%s", viper.GetString("GENERATIVE_EXPLORER_API"), id))
					if err != nil {
						logger.AtLog.Logger.Error("Get list inscriptions error", zap.Error(err))
					}

					if resp.Sat == 0 || resp.Address == "" {
						return
					}

					resp.ObjectID = id
					data = append(data, resp)
				}(r)
			}
			wg.Wait()

			uc.algoliaClient.BulkIndexer("inscriptions", data)
			uc.redis.Set(ctx, "Inscription_Index_Count", index-100, time.Duration(time.Hour*24*30))
			logger.AtLog.Infof("INDEXING %d", index)
		}(result)
		index += 100
	}

	// if err := uc.redis.Set(ctx, "Inscription_Index_Count", index, time.Duration(time.Hour*24*30)); err != nil {
	// 	return err
	// }
	logger.AtLog.Infof("DONE inscriptionIndexingData algolia data %v", time.Now())
	return nil
}
