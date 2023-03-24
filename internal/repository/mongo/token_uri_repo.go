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

var _ port.ITokenUriRepository = (*tokenRepository)(nil)

type tokenRepository struct {
	mongodb.BaseRepository
}

func NewTokenRepository(db *mongo.Database) port.ITokenUriRepository {
	return &tokenRepository{
		BaseRepository: mongodb.BaseRepository{
			CollectionName: "token_uri",
			DB:             db,
		},
	}
}

func (r tokenRepository) ProjectGetMintVolume() ([]*model.TokenUriListingVolume, error) {
	result := []*model.TokenUriListingVolume{}
	pipeline := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"isMinted", true},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$projectID"},
					{"total_amount", bson.D{{"$sum", "$project_mint_price"}}},
				},
			},
		},
	}

	cursor, err := r.DB.Collection("mint_nft_btc").Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = cursor.All((context.TODO()), &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

func (r tokenRepository) ProjectGetCurrentListingNumber() ([]*model.TokenUriListingPage, error) {
	result := []*model.TokenUriListingPage{}
	pipeline := bson.A{
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "dex_btc_listing"},
					{"localField", "token_id"},
					{"foreignField", "inscription_id"},
					{"let",
						bson.D{
							{"cancelled", "$cancelled"},
							{"matched", "$matched"},
						},
					},
					{"pipeline",
						bson.A{
							bson.D{
								{"$match",
									bson.D{
										{"matched", false},
										{"cancelled", false},
									},
								},
							},
						},
					},
					{"as", "listing"},
				},
			},
		},
		bson.D{
			{"$unwind",
				bson.D{
					{"path", "$listing"},
					{"preserveNullAndEmptyArrays", false},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", "$project_id"},
					{"count", bson.D{{"$sum", 1}}},
				},
			},
		},
	}

	cursor, err := r.DB.Collection("token_uri").Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = cursor.All((context.TODO()), &result); err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}
