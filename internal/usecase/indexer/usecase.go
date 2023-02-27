package usecase

import (
	"context"
	"generative-xyz-search-engine/internal/core/port"
	"generative-xyz-search-engine/pkg/driver/algolia"
	"generative-xyz-search-engine/pkg/entity"
	"generative-xyz-search-engine/pkg/logger"
	"generative-xyz-search-engine/pkg/model"
	"generative-xyz-search-engine/utils"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"

	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type indexerUsecase struct {
	algoliaClient *algolia.GenerativeAlgolia
	projectRepo   port.IProjectRepository
	tokenUriRepo  port.ITokenUriRepository

	ch chan struct{}
}

func (uc *indexerUsecase) Schedule() {
	s := gocron.NewScheduler(time.Local)
	s = s.Cron(viper.GetString("ALGOLIA_INDEX_SCAN_CRON"))
	s.StartImmediately()
	_, err := s.Do(func() {
		uc.ProcessIndexDataAlgolia(context.Background())
	})

	if err != nil {
		logger.AtLog.Logger.Fatal("indexerUsecase.Schedule", zap.Error(err))
		return
	}

	s.StartBlocking()
	uc.ch <- struct{}{}
}

func (uc *indexerUsecase) ProcessIndexDataAlgolia(rootCtx context.Context) error {
	var err error
	span, ctx := tracer.StartSpanFromContext(rootCtx, "indexProduct.ProcessMessage")

	defer func() {
		var spanOpts []tracer.FinishOption
		if err != nil {
			spanOpts = append(spanOpts, tracer.WithError(err))
		}
		span.Finish(spanOpts...)
	}()

	logger.AtLog.Infof("START indexing algolia data %v", time.Now())
	mainW := &sync.WaitGroup{}

	if err = uc.indexingProjectData(ctx, mainW); err != nil {
		return err
	}

	if err = uc.indexingTokenUriData(ctx, mainW); err != nil {
		return err
	}

	logger.AtLog.Infof("DONE indexing algolia data %v", time.Now())
	return nil
}

func (uc *indexerUsecase) indexingProjectData(ctx context.Context, mainW *sync.WaitGroup) error {
	limit := int64(500)
	page := int64(1)
	for {
		var projects []*model.Project
		filters := make(map[string]interface{})

		_, err := uc.projectRepo.Filter(ctx, filters, []string{}, []int{}, page, limit, &projects)
		if err != nil {
			logger.AtLog.Logger.Error(err.Error(), zap.Error(err))
			return err
		}

		if len(projects) == 0 {
			break
		}

		data := make([]*entity.ProjectAlgolia, 0)
		for _, p := range projects {
			d := &entity.ProjectAlgolia{}
			if p.TokenID == "" {
				continue
			}

			if err := utils.Copy(d, p); err != nil {
				logger.AtLog.Logger.Error(err.Error(), zap.Error(err))
				return err
			}

			d.ObjectID = p.TokenID
			d.DeletedAt = p.DeletedAt
			d.Image = p.Thumbnail
			data = append(data, d)

		}

		mainW.Add(1)
		go uc.algoliaClient.BulkIndexer("projects", data, mainW)
		page += 1
	}
	return nil
}

func (uc *indexerUsecase) indexingTokenUriData(ctx context.Context, mainW *sync.WaitGroup) error {
	limit := int64(500)
	page := int64(1)
	for {
		var projects []*model.TokenUri
		filters := make(map[string]interface{})

		_, err := uc.tokenUriRepo.Filter(ctx, filters, []string{}, []int{}, page, limit, &projects)
		if err != nil {
			logger.AtLog.Logger.Error(err.Error(), zap.Error(err))
			return err
		}

		if len(projects) == 0 {
			break
		}

		data := make([]*entity.TokenUriAlgolia, 0)
		for _, p := range projects {
			d := &entity.TokenUriAlgolia{}
			if p.TokenID == "" {
				continue
			}

			if err := utils.Copy(d, p); err != nil {
				logger.AtLog.Logger.Error(err.Error(), zap.Error(err))
				return err
			}

			d.ObjectID = p.TokenID
			d.Image = p.Thumbnail
			data = append(data, d)

		}

		mainW.Add(1)
		go uc.algoliaClient.BulkIndexer("token-uris", data, mainW)
		page += 1
	}
	return nil

}

func NewProjectIndexerUsecase(client *algolia.GenerativeAlgolia, repo port.IProjectRepository, tokenUriRepo port.ITokenUriRepository, ch chan struct{}) port.IIndexerUsecase {
	return &indexerUsecase{
		algoliaClient: client,
		projectRepo:   repo,
		tokenUriRepo:  tokenUriRepo,
		ch:            ch,
	}
}
