package usecase

import (
	"context"
	"fmt"
	"generative-xyz-search-engine/pkg/driver/algolia"
	"generative-xyz-search-engine/pkg/entity"
	"generative-xyz-search-engine/pkg/logger"
	"generative-xyz-search-engine/utils"
	"strings"
	"time"
)

func (uc *indexerUsecase) indexProjectListingData(ctx context.Context, isDelta bool) error {
	logger.AtLog.Infof("START indexProjectListingData algolia data %v", time.Now())
	filterToken := algolia.AlgoliaFilter{
		Page: 0, Limit: 1000,
	}

	projectMapData, err := uc.fetchAllProjectData(ctx)
	if err != nil {
		logger.AtLog.Error(err)
		return err
	}

	userMapData, err := uc.fetchAllUserData(ctx)
	if err != nil {
		logger.AtLog.Error(err)
		return err
	}
	projectListingMapData := make(map[string]*entity.ProjectListing)
	for {
		logger.AtLog.Info(filterToken)
		var tokens []*entity.TokenUriAlgolia
		resp, err := uc.algoliaClient.Search("token-uris", &filterToken)
		if err != nil {
			logger.AtLog.Error(err)
			return err
		}
		resp.UnmarshalHits(&tokens)

		if len(tokens) == 0 {
			break
		}

		for _, token := range tokens {
			if project, ok := projectMapData[token.ProjectID]; ok {
				skip := false
				for _, c := range project.Categories {
					if c == "63f8325a1460b1502544101b" {
						skip = true
						break
					}
				}

				if _, ok := projectListingMapData[token.ProjectID]; !ok && !skip {
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
					projectListingMapData[token.ProjectID] = listing
				}
			}
		}
		filterToken.Page += 1
	}

	data := []*entity.ProjectListing{}

	for _, btc := range projectListingMapData {
		projectID := btc.Project.TokenId
		logger.AtLog.Infof("processing: %s", projectID)
		var oTokens []*entity.TokenUriAlgolia
		resp, err := uc.algoliaClient.Search("token-uris", &algolia.AlgoliaFilter{
			Limit: 10_000, FilterStrs: []string{fmt.Sprintf("projectID:%s", projectID)},
		})

		if err != nil {
			logger.AtLog.Error(err)
			return err
		}

		resp.UnmarshalHits(&oTokens)
		addresses := []string{}
		filters := []string{}
		for _, t := range oTokens {
			filters = append(filters, fmt.Sprintf("inscription_id:%s", t.TokenID))
		}

		resp, err = uc.algoliaClient.Search("inscriptions", &algolia.AlgoliaFilter{
			Limit: 10_000, FilterStrs: []string{strings.Join(filters, " OR ")},
		})

		if err == nil {
			for _, r := range resp.Hits {
				if r["address"] != nil {
					addresses = append(addresses, r["address"].(string))
				}

			}
		}
		btc.NumberOwners = int64(len(utils.RemoveDuplicateValues(addresses)))

		floorPrice, _ := uc.dexBtcListingRepo.RetrieveFloorPriceOfCollection(projectID)

		if floorPrice <= 0 && btc.Project.MintingInfo.Index < btc.Project.MaxSupply {
			btc.IsHidden = true
		} else {
			currentListing, _ := uc.tokenUriRepo.ProjectGetCurrentListingNumber(projectID)
			volume, _ := uc.dexBtcListingRepo.ProjectGetListingVolume(projectID)

			volumeCEX, _ := uc.dexBtcListingRepo.ProjectGetCEXVolume(projectID)
			mintVolume, _ := uc.tokenUriRepo.ProjectGetMintVolume(projectID)

			btc.ProjectMarketplaceData = &entity.ProjectMarketplaceData{
				FloorPrice:  floorPrice,
				Listed:      currentListing,
				TotalVolume: volume + mintVolume + volumeCEX,
				MintVolume:  mintVolume,
			}
			btc.TotalVolume = volume + mintVolume + volumeCEX
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

	logger.AtLog.Infof("DONE indexProjectListingData algolia data %v", time.Now())
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
