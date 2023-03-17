package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenUri struct {
	Model       `bson:"inline"`
	TokenID     string `bson:"token_id" json:"token_id"`
	TokenIDInt  int    `bson:"token_id_int" json:"token_id_int"`
	TokenIDMini *int   `bson:"token_id_mini" json:"token_id_mini"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Image       string `bson:"image" json:"image"`

	ProjectID         string  `bson:"project_id" json:"project_id"`
	ProjectIDInt      int64   `bson:"project_id_int" json:"project_id_int"`
	BlockNumberMinted *string `bson:"block_number_minted" json:"block_number_minted"`
	Thumbnail         string  `bson:"thumbnail"`

	OwnerAddr     string  `bson:"owner_addrress"`
	CreatorAddr   string  `bson:"creator_address"`
	Priority      *int    `bson:"priority"`
	MinterAddress *string `bson:"minter_address"`

	IsOnchain                      bool     `bson:"isOnchain"`
	InscriptionIndex               string   `bson:"inscription_index"`
	SyncedInscriptionInfo          bool     `bson:"synced_inscription_info"`
	CreatedByCollectionInscription bool     `bson:"created_by_collection_inscription"`
	Source                         string   `bson:"source" json:"source"`
	Project                        *Project `bson:"project"`
	OrderInscriptionIndex          int64    `bson:"order_inscription_index" json:"orderInscriptionIndex"`
}

type TokenUriListingPage struct {
	TotalData  []TokenUriListingFilter `bson:"totalData" json:"totalData"`
	TotalCount []struct {
		Count int64 `bson:"count" json:"count"`
	} `bson:"totalCount" json:"totalCount"`
}

type TokenUriListingFilter struct {
	ID                    primitive.ObjectID `bson:"_id" json:"_id"`
	TokenID               string             `bson:"token_id" json:"tokenID"`
	Name                  string             `bson:"name" json:"name"`
	Image                 string             `bson:"image" json:"image"`
	ContractAddress       string             `bson:"contract_address" json:"contract_address"`
	AnimationURL          string             `bson:"animation_url" json:"animation_url"`
	AnimationHtml         *string            `bson:"animation_html"`
	ProjectID             string             `bson:"project_id" json:"projectID"`
	MintedTime            *time.Time         `bson:"minted_time" json:"minted_time"`
	GenNFTAddr            string             `bson:"gen_nft_addrress" json:"genNFTAddr"`
	Thumbnail             string             `bson:"thumbnail" json:"thumbnail"`
	InscriptionIndex      string             `bson:"inscription_index" json:"inscriptionIndex"`
	OrderInscriptionIndex int                `bson:"order_inscription_index" json:"orderInscriptionIndex"`
	OrderID               primitive.ObjectID `bson:"orderID" json:"orderID"`
	Price                 int64              `bson:"priceBTC" json:"priceBTC"`
	PriceETH              string             `bson:"priceETH" json:"priceETH"`
	Buyable               bool               `bson:"buyable" json:"buyable"`
	SellVerified          bool               `bson:"sell_verified" json:"sell_verified"`
	Project               struct {
		TokenID string `bson:"tokenid" json:"tokenID"`
		Royalty int64  `bson:"royalty" json:"royalty"`
	} `bson:"project" json:"project"`
}
