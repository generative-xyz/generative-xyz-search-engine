package mongo

import (
	"generative-xyz-search-engine/internal/core/port"
	"generative-xyz-search-engine/pkg/driver/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ port.IProjectRepository = (*projectRepository)(nil)

type projectRepository struct {
	mongodb.BaseRepository
}

func NewProjectRepository(db *mongo.Database) port.IProjectRepository {
	return &projectRepository{
		BaseRepository: mongodb.BaseRepository{
			CollectionName: "projects",
			DB:             db,
		},
	}
}

func (r *projectRepository) SelectedProjectFields() bson.M {
	f := bson.M{
		"_id":                    1,
		"contractAddress":        1,
		"tokenid":                1,
		"tokenIDInt":             1,
		"maxSupply":              1,
		"mintPrice":              1,
		"name":                   1,
		"creatorName":            1,
		"creatorAddress":         1,
		"creatorAddrrBTC":        1,
		"categories":             1,
		"thumbnail":              1,
		"mintFee":                1,
		"openMintUnixTimestamp":  1,
		"closeMintUnixTimestamp": 1,
		"genNFTAddr":             1,
		"mintTokenAddress":       1,
		"minted_time":            1,
		"license":                1,
		"description":            1,
		"stats":                  1,
		"status":                 1,
		"isSynced":               1,
		"completeTime":           1,
		"block_number_minted":    1,
		"traitsStat":             1,
		"priority":               1,
		"isHidden":               1,
		"tokenDescription":       1,
		"socialWeb":              1,
		"socialTwitter":          1,
		"socialDiscord":          1,
		"socialMedium":           1,
		"socialInstagram":        1,
		"index":                  1,
		"indexReverse":           1,
		"creatorProfile":         1,
		"images":                 1,
		"whiteListEthContracts":  1,
		"isFullChain":            1,
		"inscription_icon":       1,
		"source":                 1,
		"reportUsers":            1,
		"mintpriceeth":           1,
		"fromAuthentic":          1,
		"tokenAddress":           1,
		"ownerOf":                1,
		"ordinalsTx":             1,
		"inscribedBy":            1,
	}
	return f
}
