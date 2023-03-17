package mongo

import (
	"context"
	"generative-xyz-search-engine/internal/core/port"
	"generative-xyz-search-engine/pkg/driver/mongodb"
	"generative-xyz-search-engine/pkg/model"

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

func (r dexBtcListingRepository) RetrieveFloorPriceOfCollection(collectionID string) (uint64, error) {
	resp := []model.MarketplaceBTCListingFloorPrice{}
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
					{"as", "collection_id"},
				},
			},
		},
		bson.D{{"$unwind", "$collection_id"}},
		bson.D{
			{"$match",
				bson.D{
					{"$expr",
						bson.D{
							{"$eq",
								bson.A{
									bson.D{
										{"$getField",
											bson.D{
												{"field", bson.D{{"$literal", "project_id"}}},
												{"input", "$collection_id"},
											},
										},
									},
									collectionID,
								},
							},
						},
					},
				},
			},
		},
		bson.D{{"$sort", bson.D{{"amount", 1}}}},
		bson.D{{"$limit", 1}},
	})

	if err != nil {
		return 0, err
	}

	if err = cursor.All(context.TODO(), &resp); err != nil {
		return 0, err
	}

	if len(resp) == 0 {
		return 0, nil
	}

	return resp[0].Price, nil
}

func (r dexBtcListingRepository) ProjectGetCEXVolume(projectID string) (uint64, error) {
	result := []model.TokenUriListingVolume{}
	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"isSold", true}}}},
		bson.D{{"$addFields", bson.D{{"price", bson.D{{"$toDouble", "$amount"}}}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "token_uri"},
					{"localField", "inscriptionID"},
					{"foreignField", "token_id"},
					{"let", bson.D{{"id", "$_id"}}},
					{"pipeline",
						bson.A{
							bson.D{{"$match", bson.D{{"project_id", projectID}}}},
						},
					},
					{"as", "collection_id"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$collection_id"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", ""},
					{"Amount", bson.D{{"$sum", "$price"}}},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 0},
					{"totalAmount", "$Amount"},
				},
			},
		},
	}

	cursor, err := r.DB.Collection("marketplace_btc_listing").Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if err = cursor.All((context.TODO()), &result); err != nil {
		return 0, errors.WithStack(err)
	}
	if len(result) > 0 {
		return uint64(result[0].TotalAmount), nil
	}

	return 0, nil
}

func (r dexBtcListingRepository) ProjectGetListingVolume(projectID string) (uint64, error) {
	result := []model.TokenUriListingVolume{}
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
					{"let", bson.D{{"id", "$_id"}}},
					{"pipeline",
						bson.A{
							bson.D{{"$match", bson.D{{"project_id", projectID}}}},
						},
					},
					{"as", "collection_id"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$collection_id"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", ""},
					{"Amount", bson.D{{"$sum", "$amount"}}},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 0},
					{"totalAmount", "$Amount"},
				},
			},
		},
	}

	cursor, err := r.DB.Collection("dex_btc_listing").Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if err = cursor.All((context.TODO()), &result); err != nil {
		return 0, errors.WithStack(err)
	}
	if len(result) > 0 {
		return uint64(result[0].TotalAmount), nil
	}

	return 0, nil
}
