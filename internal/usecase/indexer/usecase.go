package usecase

import (
	"context"
	"fmt"
	"generative-xyz-search-engine/internal/core/port"
	"generative-xyz-search-engine/pkg/driver/algolia"
	"generative-xyz-search-engine/pkg/entity"
	"generative-xyz-search-engine/pkg/logger"
	"generative-xyz-search-engine/pkg/model"
	"generative-xyz-search-engine/pkg/redis"
	"generative-xyz-search-engine/utils"
	"generative-xyz-search-engine/utils/constants"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type indexerUsecase struct {
	algoliaClient *algolia.GenerativeAlgolia
	redis         redis.Client
	projectRepo   port.IProjectRepository
	tokenUriRepo  port.ITokenUriRepository
	userRepo      port.IUserRepository

	ch chan struct{}
}

func (uc *indexerUsecase) Schedule() {
	s := gocron.NewScheduler(time.Local)
	s = s.Cron(viper.GetString("ALGOLIA_INDEX_FULL_SCAN_CRON"))
	// s.StartImmediately()
	_, err := s.Do(func() {
		uc.ProcessIndexDataAlgolia(context.Background(), false)
	})

	if err != nil {
		logger.AtLog.Logger.Fatal("indexerUsecase.Schedule.ALGOLIA_INDEX_FULL_SCAN_CRON", zap.Error(err))
		return
	}

	s = s.Cron(viper.GetString("ALGOLIA_INDEX_SCAN_CRON"))
	// s.StartImmediately()
	_, err = s.Do(func() {
		uc.ProcessIndexDataAlgolia(context.Background(), true)
	})

	if err != nil {
		logger.AtLog.Logger.Fatal("indexerUsecase.Schedule.ALGOLIA_INDEX_SCAN_CRON", zap.Error(err))
		return
	}

	s.StartBlocking()
	uc.ch <- struct{}{}
}

func (uc *indexerUsecase) ProcessIndexDataAlgolia(rootCtx context.Context, isDelta bool) error {
	if len(uc.ch) > 0 {
		logger.AtLog.Warn("ProcessIndexDataAlgolia.Execute was skipped.")
		return nil
	}

	var err error
	span, ctx := tracer.StartSpanFromContext(rootCtx, "indexProduct.ProcessMessage")

	defer func() {
		var spanOpts []tracer.FinishOption
		if err != nil {
			spanOpts = append(spanOpts, tracer.WithError(err))
		}
		span.Finish(spanOpts...)
	}()

	if err = uc.indexingUserData(ctx, isDelta); err != nil {
		return err
	}

	if err = uc.indexingProjectData(ctx, isDelta); err != nil {
		return err
	}

	if err = uc.indexingTokenUriData(ctx, isDelta); err != nil {
		return err
	}

	if err = uc.inscriptionIndexingData(ctx, isDelta); err != nil {
		return err
	}

	return nil
}

func (uc *indexerUsecase) indexingUserData(ctx context.Context, isDelta bool) error {
	limit := int64(500)
	lastId := ""
	now := time.Now()

	logger.AtLog.Infof("START indexingUserData algolia data %v", time.Now())
	for {
		var users []*model.User
		filters := make(map[string]interface{})
		if isDelta {
			filters["updated_at"] = bson.M{"$gte": now.Add(constants.DeltaIndexingDataHours)}
		}

		if lastId != "" {
			if id, err := primitive.ObjectIDFromHex(lastId); err == nil {
				filters["_id"] = bson.M{"$lt": id}
			}
		}

		_, err := uc.userRepo.Filter(ctx, filters, []string{"_id"}, []int{-1}, 0, limit, nil, &users)
		if err != nil {
			logger.AtLog.Logger.Error(err.Error(), zap.Error(err))
			return err
		}

		if len(users) == 0 {
			break
		}

		data := make([]*entity.UserAlgolia, 0)
		for _, u := range users {
			d := &entity.UserAlgolia{}
			if err := utils.Copy(d, u); err != nil {
				logger.AtLog.Logger.Error(err.Error(), zap.Error(err))
				return err
			}

			d.ObjectID = u.Id.Hex()
			d.Stats = entity.UserStats(u.Stats)
			d.ProfileSocial = entity.ProfileSocial(u.ProfileSocial)
			data = append(data, d)
		}
		lastId = users[len(users)-1].Id.Hex()

		uc.algoliaClient.BulkIndexer("users", data)
	}

	logger.AtLog.Infof("DONE indexingUserData algolia data %v", time.Now())
	return nil
}

func (uc *indexerUsecase) indexingProjectData(ctx context.Context, isDelta bool) error {
	limit := int64(500)
	lastId := ""
	now := time.Now()

	logger.AtLog.Infof("START indexingProjectData algolia data %v", time.Now())
	for {
		var projects []*model.Project
		filters := make(map[string]interface{})
		if isDelta {
			filters["updated_at"] = bson.M{"$gte": now.Add(constants.DeltaIndexingDataHours)}
		}

		if lastId != "" {
			if id, err := primitive.ObjectIDFromHex(lastId); err == nil {
				filters["_id"] = bson.M{"$lt": id}
			}
		}

		_, err := uc.projectRepo.Filter(
			ctx, filters, []string{"_id"}, []int{-1}, 0, limit, uc.projectRepo.SelectedProjectFields(), &projects,
		)
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

			d.ObjectID = p.Id.Hex()
			d.DeletedAt = p.DeletedAt
			d.Image = p.Thumbnail
			data = append(data, d)
		}
		lastId = projects[len(projects)-1].Id.Hex()

		uc.algoliaClient.BulkIndexer("projects", data)
	}

	logger.AtLog.Infof("DONE indexingProjectData algolia data %v", time.Now())
	return nil
}

func (uc *indexerUsecase) indexingTokenUriData(ctx context.Context, isDelta bool) error {
	limit := int64(500)
	logger.AtLog.Infof("START indexingTokenUriData algolia data %v", time.Now())
	lastId := ""
	count := 0
	now := time.Now()

	for {
		var tokens []*model.TokenUri
		filters := make(map[string]interface{})
		if isDelta {
			filters["updated_at"] = bson.M{"$gte": now.Add(constants.DeltaIndexingDataHours)}
		}

		if lastId != "" {
			if id, err := primitive.ObjectIDFromHex(lastId); err == nil {
				filters["_id"] = bson.M{"$lt": id}
			}
		}

		_, err := uc.tokenUriRepo.Filter(ctx, filters, []string{"_id"}, []int{-1}, 0, limit, nil, &tokens)
		if err != nil {
			logger.AtLog.Logger.Error(err.Error(), zap.Error(err))
			return err
		}

		if len(tokens) == 0 {
			break
		}

		data := make([]*entity.TokenUriAlgolia, 0)
		for _, p := range tokens {
			d := &entity.TokenUriAlgolia{}
			if p.TokenID == "" {
				continue
			}

			d.ObjectID = p.Id.Hex()
			d.TokenID = p.TokenID
			d.Name = p.Name
			d.Description = p.Description
			d.Image = p.Thumbnail
			d.InscriptionIndex = p.InscriptionIndex

			filters := make(map[string]interface{})
			filters["tokenid"] = p.ProjectID

			project := &model.Project{}
			err := uc.projectRepo.FindOne(ctx, filters, &project)
			if err != nil {
				logger.AtLog.Logger.Error(err.Error(), zap.Error(err))
			}

			if project != nil {
				d.ProjectName = fmt.Sprintf("%s #%d", project.Name, p.OrderInscriptionIndex)
				d.ProjectID = p.ProjectID
			}

			d.Thumbnail = p.Thumbnail
			data = append(data, d)
		}
		uc.algoliaClient.BulkIndexer("token-uris", data)
		lastId = tokens[len(tokens)-1].Id.Hex()

		count += len(tokens)
		logger.AtLog.Infof("Count: %d", count)
	}

	logger.AtLog.Infof("DONE indexingTokenUriData algolia data %v", time.Now())
	return nil
}

func NewProjectIndexerUsecase(
	client *algolia.GenerativeAlgolia,
	redis redis.Client,
	repo port.IProjectRepository,
	tokenUriRepo port.ITokenUriRepository,
	userRepo port.IUserRepository,
	ch chan struct{},
) port.IIndexerUsecase {
	return &indexerUsecase{
		algoliaClient: client,
		redis:         redis,
		projectRepo:   repo,
		tokenUriRepo:  tokenUriRepo,
		userRepo:      userRepo,
		ch:            ch,
	}
}
