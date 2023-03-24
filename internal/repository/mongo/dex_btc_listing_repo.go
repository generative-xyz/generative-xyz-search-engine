package mongo

import (
	"context"
	"generative-xyz-search-engine/internal/core/port"
	"generative-xyz-search-engine/pkg/driver/mongodb"
	"generative-xyz-search-engine/pkg/model"
	"generative-xyz-search-engine/utils"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ port.IDexBtcListingRepository = (*dexBtcListingRepository)(nil)

type dexBtcListingRepository struct {
	mongodb.BaseRepository
}

func NewDexBtcListingRepository(db *mongo.Database) port.IDexBtcListingRepository {
	return &dexBtcListingRepository{
		BaseRepository: mongodb.BaseRepository{
			CollectionName: "dex_btc_listing",
			DB:             db,
		},
	}
}

func (r dexBtcListingRepository) RetrieveFloorPriceOfCollection() ([]*model.MarketplaceBTCListingFloorPrice, error) {
	resp := []*model.MarketplaceBTCListingFloorPrice{}
	cursor, err := r.DB.Collection("dex_btc_listing").Aggregate(context.TODO(), bson.A{
		bson.D{
			{"$project",
				bson.D{
					{"_id", 1},
					{"amount", 1},
					{"inscription_id", 1},
					{"matched", 1},
					{"cancelled", 1},
				},
			},
		},
		bson.D{
			{"$match",
				bson.D{
					{"matched", false},
					{"cancelled", false},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "token_uri"},
					{"localField", "inscription_id"},
					{"foreignField", "token_id"},
					{"let", bson.D{{"id", "$_id"}}},
					{"pipeline",
						bson.A{
							bson.D{{"$project", bson.D{{"project_id", 1}}}},
						},
					},
					{"as", "collections"},
				},
			},
		},
		bson.D{{"$unwind", "$collections"}},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$collections.project_id"},
					{"min_amount", bson.D{{"$min", "$amount"}}},
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.TODO(), &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (r dexBtcListingRepository) ProjectGetCEXVolume() ([]*model.TokenUriListingVolume, error) {
	result := []*model.TokenUriListingVolume{}
	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"isSold", true}}}},
		bson.D{{"$addFields", bson.D{{"price", bson.D{{"$toDouble", "$amount"}}}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "token_uri"},
					{"localField", "inscriptionID"},
					{"foreignField", "token_id"},
					{"as", "collections"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$collections"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$collections.project_id"},
					{"total_amount", bson.D{{"$sum", "$price"}}},
				},
			},
		},
	}

	cursor, err := r.DB.Collection("marketplace_btc_listing").Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = cursor.All((context.TODO()), &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (r dexBtcListingRepository) ProjectGetListingVolume() ([]*model.TokenUriListingVolume, error) {
	result := []*model.TokenUriListingVolume{}
	pipeline := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"matched", true},
					{"cancelled", false},
					{"buyer", bson.D{{"$exists", true}}},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"diffbuyer",
						bson.D{
							{"$ne",
								bson.A{
									"$buyer",
									"$seller_address",
								},
							},
						},
					},
				},
			},
		},
		bson.D{{"$match", bson.D{{"diffbuyer", true}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "token_uri"},
					{"localField", "inscription_id"},
					{"foreignField", "token_id"},
					{"as", "collections"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$collections"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$collections.project_id"},
					{"total_amount", bson.D{{"$sum", "$amount"}}},
				},
			},
		},
	}

	cursor, err := r.DB.Collection("dex_btc_listing").Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = cursor.All((context.TODO()), &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (r dexBtcListingRepository) AggregateBTCVolumn() ([]model.AggregateProjectItemResp, error) {
	//resp := &entity.AggregateWalletAddres{}
	confs := []model.AggregateProjectItemResp{}

	calculate := bson.M{"$sum": "$project_mint_price"}
	// PayType *string
	// ReferreeIDs []string
	matchStage := bson.M{"$match": bson.M{"$and": bson.A{
		bson.M{"isMinted": true},
		// bson.M{"projectID": projectID},
	}}}

	pipeLine := bson.A{
		matchStage,
		bson.M{"$group": bson.M{"_id": bson.M{"projectID": "$projectID"},
			"amount": calculate,
			"minted": bson.M{"$sum": 1},
		}},
		bson.M{"$sort": bson.M{"_id": -1}},
	}

	cursor, err := r.DB.Collection("mint_nft_btc").Aggregate(context.TODO(), pipeLine, nil)
	if err != nil {
		return nil, err
	}

	// display the results
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	for _, item := range results {
		res := &model.AggregateProjectItem{}
		err = utils.Transform(item, res)
		if err != nil {
			return nil, err
		}

		tmp := model.AggregateProjectItemResp{
			ProjectID: res.ID.ProjectID,
			Paytype:   res.ID.Paytype,
			BtcRate:   res.ID.BtcRate,
			EthRate:   res.ID.EthRate,
			MintPrice: res.ID.MintPrice,
			Amount:    res.Amount,
			Minted:    res.Minted,
		}
		confs = append(confs, tmp)
	}

	return confs, nil
}
