package usecase

import (
	"context"
	"fmt"
	"generative-xyz-search-engine/pkg/driver/algolia"
	"generative-xyz-search-engine/pkg/entity"
	"generative-xyz-search-engine/pkg/logger"
	"generative-xyz-search-engine/pkg/model"
	"generative-xyz-search-engine/utils"
	"generative-xyz-search-engine/utils/constants"
	"strconv"
	"strings"
	"sync"
	"time"
)

var specialProjectPrice = map[string]uint64{
	"1002573": 263000000,
}

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
	btcVolumesMap := map[string]*model.AggregateProjectItemResp{}
	go func() {
		defer wG.Done()
		btcVolumes, _ := uc.dexBtcListingRepo.AggregateBTCVolumn()
		for _, i := range btcVolumes {
			btcVolumesMap[i.ProjectID] = i
		}
	}()

	wG.Add(1)
	mapFloorPrice := map[string]*model.MarketplaceBTCListingFloorPrice{}
	go func() {
		defer wG.Done()
		floorPrices, _ := uc.dexBtcListingRepo.RetrieveFloorPriceOfCollection()
		for _, i := range floorPrices {
			mapFloorPrice[i.ProjectId] = i
		}
	}()

	wG.Add(1)
	mapCurrentListing := map[string]*model.TokenUriListingPage{}
	go func() {
		defer wG.Done()
		currentListings, _ := uc.tokenUriRepo.ProjectGetCurrentListingNumber()
		for _, i := range currentListings {
			mapCurrentListing[i.ProjectId] = i
		}
	}()

	// wG.Add(1)
	// mapMintVolume := map[string]*model.TokenUriListingVolume{}
	// go func() {
	// 	defer wG.Done()
	// 	mintVolumes, _ := uc.tokenUriRepo.ProjectGetMintVolume()
	// 	for _, i := range mintVolumes {
	// 		mapMintVolume[i.ProjectId] = i
	// 	}
	// }()

	wG.Add(1)
	mapVolumeCEX := map[string]*model.TokenUriListingVolume{}
	go func() {
		defer wG.Done()
		volumeCEXs, _ := uc.dexBtcListingRepo.ProjectGetCEXVolume()
		for _, i := range volumeCEXs {
			mapVolumeCEX[i.ProjectId] = i
		}
	}()

	wG.Add(1)
	mapVolume := map[string]*model.TokenUriListingVolume{}
	go func() {
		defer wG.Done()
		volumes, _ := uc.dexBtcListingRepo.ProjectGetListingVolume()
		for _, i := range volumes {
			mapVolume[i.ProjectId] = i
		}
	}()

	wG.Add(1)
	mapOldETHVolume := map[string]*model.AggregateProjectItemResp{}
	go func() {
		defer wG.Done()
		ethVolumes, _ := uc.dexBtcListingRepo.AggregationETHWalletAddress()
		for _, i := range ethVolumes {
			mapOldETHVolume[i.ProjectID] = i
		}
	}()

	wG.Add(1)
	mapOldBTCVolume := map[string]*model.AggregateProjectItemResp{}
	go func() {
		defer wG.Done()
		btcVolumes, _ := uc.dexBtcListingRepo.AggregationBTCWalletAddress()
		for _, i := range btcVolumes {
			mapOldBTCVolume[i.ProjectID] = i
		}
	}()

	wG.Wait()

	projectListingMapData := make(map[string]*entity.ProjectListing)
	for _, project := range projectMapData {
		skip := false
		for _, c := range project.Categories {
			if c == constants.CategoryUnverifiedIdStr {
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
	// client := resty.New()
	// apiUrl := viper.GetString("GENERATIVE_EXPLORER_API")

	for _, btc := range projectListingMapData {
		projectID := btc.Project.TokenId
		// if projectID != "1000001" {
		// 	continue
		// }

		floorPrice := uint64(0)
		if btc.Project.MintingInfo.Index < btc.Project.MaxSupply {
			if num, err := strconv.ParseUint(btc.MintPrice, 10, 64); err == nil {
				floorPrice = num
			}
		} else {
			if price, ok := mapFloorPrice[projectID]; ok {
				floorPrice = price.FloorPrice
			}
		}

		project := projectMapData[projectID]
		hidden := false
		if project != nil && project.IsHidden {
			hidden = true
		}

		if floorPrice <= 0 && btc.Project.MintingInfo.Index < btc.Project.MaxSupply {
			hidden = true
		}

		// logger.AtLog.Infof("processing: %s - isHidden %v", projectID, hidden)
		if hidden {
			btc.IsHidden = true
		} else {
			filter := &algolia.AlgoliaFilter{
				Page: 0, Limit: 500, FilterStrs: []string{fmt.Sprintf("projectID:%s", projectID)},
			}
			addresses := []string{}
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
					addresses = utils.RemoveDuplicateValues(addresses)
				}

				// wG1.Add(1)
				// 	go func(t *entity.TokenUriAlgolia) {
				// 		defer wG1.Done()
				// 		r := &InscriptionDetail{}
				// 		_, err = client.R().SetResult(&r).
				// 			Get(fmt.Sprintf("%s/inscription/%s", apiUrl, t.TokenID))
				// 		if err != nil {
				// 			logger.AtLog.Logger.Error("Get list inscriptions error", zap.Error(err))
				// 			return
				// 		}
				// 		addresses = append(addresses, r.Address)
				// 	}(t)

				// wG1.Wait()
				filter.Page += 1
			}
			btc.NumberOwners = int64(len(utils.RemoveDuplicateValues(addresses)))

			currentListing := uint64(0)
			if v, ok := mapCurrentListing[projectID]; ok {
				currentListing = v.Count
			}

			// mintVolume := uint64(0)
			// if v, ok := mapMintVolume[projectID]; ok {
			// 	mintVolume = v.TotalAmount
			// }

			volumeCEX := uint64(0)
			if v, ok := mapVolumeCEX[projectID]; ok {
				volumeCEX = v.TotalAmount
			}

			volume := uint64(0)
			if v, ok := mapVolume[projectID]; ok {
				volume = v.TotalAmount
			}

			firstSaleVolume := float64(0)
			if firstVolume, ok := btcVolumesMap[projectID]; ok {
				firstSaleVolume += firstVolume.Amount
			}

			if v, ok := mapOldETHVolume[projectID]; ok {
				firstSaleVolume += v.Amount / float64(v.BtcRate/v.EthRate)
			}

			if v, ok := mapOldBTCVolume[projectID]; ok {
				firstSaleVolume += v.Amount
			}

			totalVolume := volume + volumeCEX + uint64(firstSaleVolume)
			// override total volume for special project
			if price, has := specialProjectPrice[projectID]; has {
				totalVolume += price
			}
			btc.ProjectMarketplaceData = &entity.ProjectMarketplaceData{
				FloorPrice:  floorPrice,
				Listed:      currentListing,
				TotalVolume: totalVolume,
				// MintVolume:      mintVolume,
				FirstSaleVolume: firstSaleVolume,
			}

			btc.TotalVolume = totalVolume
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
		logger.AtLog.Infof("update total %d project listing", len(data))
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
