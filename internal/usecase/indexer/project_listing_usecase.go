package usecase

import (
	"context"
	"fmt"
	"generative-xyz-search-engine/pkg/driver/algolia"
	"generative-xyz-search-engine/pkg/entity"
	"generative-xyz-search-engine/pkg/logger"
	"generative-xyz-search-engine/pkg/model"
	"generative-xyz-search-engine/utils"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func (uc *indexerUsecase) indexProjectListingData(ctx context.Context, isDelta bool) error {
	logger.AtLog.Infof("START indexProjectListingData algolia data %v", time.Now())
	defer logger.AtLog.Infof("DONE indexProjectListingData algolia data %v", time.Now())
	if time.Now().Minute() < 45 {
		return nil
	}

	wG := &sync.WaitGroup{}
	projectMapData := map[string]*entity.ProjectAlgolia{}
	var err error

	wG.Add(1)
	go func() {
		defer wG.Done()
		projectMapData, err = uc.fetchAllProjectData(ctx)
		if err != nil {
			logger.AtLog.Error(err)
		}
	}()

	userMapData := map[string]*entity.UserAlgolia{}
	wG.Add(1)
	go func() {
		defer wG.Done()
		userMapData, err = uc.fetchAllUserData(ctx)
		if err != nil {
			logger.AtLog.Error(err)
		}
	}()

	wG.Add(1)
	btcVolumesMap := map[string]model.AggregateProjectItemResp{}
	go func() {
		defer wG.Done()
		btcVolumes, _ := uc.dexBtcListingRepo.AggregateBTCVolumn()
		for _, i := range btcVolumes {
			btcVolumesMap[i.ProjectID] = i
		}
	}()

	wG.Wait()

	projectListingMapData := make(map[string]*entity.ProjectListing)
	for _, project := range projectMapData {
		skip := false
		for _, c := range project.Categories {
			if c == "63f8325a1460b1502544101b" {
				skip = true
				break
			}
		}

		if skip {
			continue
		}

		if _, ok := projectListingMapData[project.TokenID]; !ok && !skip {
			listing := &entity.ProjectListing{
				ObjectID: project.TokenID,
				Project: &entity.ProjectInfo{
					Name:            project.Name,
					TokenId:         project.TokenID,
					Thumbnail:       project.Image,
					ContractAddress: project.ContractAddress,
					CreatorAddress:  project.CreatorAddrr,
					MaxSupply:       project.MaxSupply,
					MintingInfo: &entity.ProjectMintingInfo{
						Index:        project.Index,
						IndexReverse: project.IndexReverse,
					},
				},
				IsHidden:  project.IsHidden,
				MintPrice: project.MintPrice,
			}

			if owner, ok := userMapData[project.CreatorAddrr]; ok {
				listing.Owner = &entity.OwnerInfo{
					WalletAddress:           owner.WalletAddress,
					WalletAddressPayment:    owner.WalletAddressPayment,
					WalletAddressBTC:        owner.WalletAddressBTC,
					WalletAddressBTCTaproot: owner.WalletAddressBTCTaproot,
					DisplayName:             owner.DisplayName,
					Avatar:                  owner.Avatar,
				}
			}
			projectListingMapData[project.TokenID] = listing
		}
	}

	data := []*entity.ProjectListing{}
	client := resty.New()

	for _, btc := range projectListingMapData {
		if btc.Project.IsHidden {
			continue
		}

		projectID := btc.Project.TokenId
		logger.AtLog.Infof("processing: %s", projectID)

		filter := &algolia.AlgoliaFilter{
			Page: 0, Limit: 500, FilterStrs: []string{fmt.Sprintf("projectID:%s", projectID)},
		}

		addresses := []string{}
		wG1 := &sync.WaitGroup{}
		for {
			var oTokens []*entity.TokenUriAlgolia
			resp, err := uc.algoliaClient.Search("token-uris", filter)
			if err != nil {
				logger.AtLog.Error(err)
				return err
			}

			resp.UnmarshalHits(&oTokens)
			if len(oTokens) == 0 {
				break
			}

			for _, t := range oTokens {
				wG1.Add(1)
				go func(t *entity.TokenUriAlgolia) {
					defer wG1.Done()
					r := &InscriptionDetail{}
					_, err = client.R().SetResult(&r).
						Get(fmt.Sprintf("%s/inscription/%s", viper.GetString("GENERATIVE_EXPLORER_API"), t.TokenID))
					if err != nil {
						logger.AtLog.Logger.Error("Get list inscriptions error", zap.Error(err))
						return
					}
					addresses = append(addresses, r.Address)
				}(t)
			}
			wG1.Wait()
			filter.Page += 1
		}

		btc.NumberOwners = int64(len(utils.RemoveDuplicateValues(addresses)))
		floorPrice := uint64(0)
		if btc.Project.MintingInfo.Index < btc.Project.MaxSupply {
			num, err := strconv.ParseUint(btc.MintPrice, 10, 64)
			if err == nil {
				floorPrice = num
			}
		} else {
			floorPrice, _ = uc.dexBtcListingRepo.RetrieveFloorPriceOfCollection(projectID)
		}

		project := projectMapData[projectID]
		hidden := false
		if project != nil && project.IsHidden {
			hidden = true
		}

		if floorPrice <= 0 && btc.Project.MintingInfo.Index < btc.Project.MaxSupply {
			hidden = true
		}

		if hidden {
			btc.IsHidden = true
		} else {
			currentListing, _ := uc.tokenUriRepo.ProjectGetCurrentListingNumber(projectID)
			volume, _ := uc.dexBtcListingRepo.ProjectGetListingVolume(projectID)

			volumeCEX, _ := uc.dexBtcListingRepo.ProjectGetCEXVolume(projectID)
			mintVolume, _ := uc.tokenUriRepo.ProjectGetMintVolume(projectID)

			firstSaleVolume := float64(0)
			if firstVolume, ok := btcVolumesMap[projectID]; ok {
				firstSaleVolume = firstVolume.Amount
			}

			btc.ProjectMarketplaceData = &entity.ProjectMarketplaceData{
				FloorPrice:      floorPrice,
				Listed:          currentListing,
				TotalVolume:     volume + mintVolume + volumeCEX + uint64(firstSaleVolume),
				MintVolume:      mintVolume,
				FirstSaleVolume: firstSaleVolume,
			}

			btc.TotalVolume = volume + mintVolume + volumeCEX + uint64(firstSaleVolume)
			btc.Priority = 3
			if btc.ProjectMarketplaceData.FloorPrice > 0 && btc.ProjectMarketplaceData.TotalVolume > 0 {
				btc.IsBuyable = true
				btc.Priority = 1
			} else if btc.ProjectMarketplaceData.FloorPrice > 0 {
				btc.Priority = 2
			}
		}

		data = append(data, btc)
		if len(data) == 500 {
			uc.algoliaClient.BulkIndexer("project-listing", data)
			data = []*entity.ProjectListing{}
		}
	}

	if len(data) > 0 {
		uc.algoliaClient.BulkIndexer("project-listing", data)
	}

	return nil
}

func (uc indexerUsecase) fetchAllProjectData(ctx context.Context) (map[string]*entity.ProjectAlgolia, error) {
	filter := algolia.AlgoliaFilter{
		Page: 0, Limit: 500,
	}
	projectMapData := make(map[string]*entity.ProjectAlgolia)
	for {
		var projects []*entity.ProjectAlgolia
		resp, err := uc.algoliaClient.Search("projects", &filter)
		if err != nil {
			logger.AtLog.Error(err)
			return nil, err
		}
		resp.UnmarshalHits(&projects)

		if len(projects) == 0 {
			break
		}

		for _, project := range projects {
			if _, ok := projectMapData[project.TokenID]; !ok {
				projectMapData[project.TokenID] = project
			}
		}
		filter.Page += 1
	}
	return projectMapData, nil
}

func (uc indexerUsecase) fetchAllDexBtcListingData(ctx context.Context) (map[string]*entity.DexBtcListingAlgolia, error) {
	filter := algolia.AlgoliaFilter{
		Page: 0, Limit: 500,
	}

	mapData := make(map[string]*entity.DexBtcListingAlgolia)
	for {
		var dexBtcListings []*entity.DexBtcListingAlgolia
		resp, err := uc.algoliaClient.Search("dex_btc_listing", &filter)
		if err != nil {
			logger.AtLog.Error(err)
			return nil, err
		}
		resp.UnmarshalHits(&dexBtcListings)

		if len(dexBtcListings) == 0 {
			break
		}

		for _, dex := range dexBtcListings {
			if _, ok := mapData[dex.InscriptionID]; !ok {
				mapData[dex.InscriptionID] = dex
			}
		}
		filter.Page += 1
	}
	return mapData, nil
}

func (uc indexerUsecase) fetchAllUserData(ctx context.Context) (map[string]*entity.UserAlgolia, error) {
	filter := algolia.AlgoliaFilter{
		Page: 0, Limit: 500,
	}

	mapData := make(map[string]*entity.UserAlgolia)
	for {
		var users []*entity.UserAlgolia
		resp, err := uc.algoliaClient.Search("users", &filter)
		if err != nil {
			logger.AtLog.Error(err)
			return nil, err
		}
		resp.UnmarshalHits(&users)

		if len(users) == 0 {
			break
		}

		for _, u := range users {
			if _, ok := mapData[u.WalletAddress]; !ok {
				mapData[u.WalletAddress] = u
			}
		}
		filter.Page += 1
	}
	return mapData, nil
}
