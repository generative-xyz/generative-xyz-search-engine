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

func (r tokenRepository) ProjectGetMintVolume(projectID string) (uint64, error) {
	result := []model.TokenUriListingVolume{}
	pipeline := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"isMinted", true},
					{"projectID", projectID},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id", ""},
					{"Amount", bson.D{{"$sum", "$project_mint_price"}}},
				},
			},
		},
	}

	cursor, err := r.DB.Collection("mint_nft_btc").Aggregate(context.TODO(), pipeline)
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

func (r tokenRepository) ProjectGetCurrentListingNumber(projectID string) (uint64, error) {
	result := []model.TokenUriListingPage{}
	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"project_id", projectID}}}},
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
			{"$facet",
				bson.D{
					{"totalCount",
						bson.A{
							bson.D{{"$count", "count"}},
						},
					},
				},
			},
		},
	}

	cursor, err := r.DB.Collection("token_uri").Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if err = cursor.All((context.TODO()), &result); err != nil {
		return 0, errors.WithStack(err)
	}
	if len(result) > 0 {
		if len(result[0].TotalCount) > 0 {
			return uint64(result[0].TotalCount[0].Count), nil
		}
		return 0, nil
	}

	return 0, nil
}
