package usecase

import (
	"context"
	"fmt"
	"generative-xyz-search-engine/pkg/logger"
	"sync"

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

func (uc *indexerUsecase) inscriptionIndexingData(ctx context.Context, isDelta bool) error {
	// Create a Resty Client
	client := resty.New()
	result := &ListInscriptionResponse{}
	index := 0
	for {
		_, err := client.R().
			EnableTrace().
			SetResult(result).
			Get(fmt.Sprintf("%s/inscriptions/%d", viper.GetString("GENERATIVE_EXPLORER_API"), index))

		if err != nil {
			logger.AtLog.Logger.Error("Get detail inscription error", zap.Error(err))
			return err
		}

		if result.Next == 0 || index > 10_000 {
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

					resp.ObjectID = id
					data = append(data, resp)
				}(r)
			}
			wg.Wait()
			uc.algoliaClient.BulkIndexer("inscriptions", data)
			logger.AtLog.Infof("INDEXING %d", index)
		}(result)

		index += 100
	}
	return nil
}