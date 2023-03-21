package entity

import (
	"generative-xyz-search-engine/pkg/model"
	"time"
)

type ProjectListing struct {
	ObjectID               string                  `json:"objectID"`
	ContractAddress        string                  `json:"contractAddress"`
	Project                *ProjectInfo            `json:"project"`
	TotalSupply            int64                   `json:"totalSupply"`
	NumberOwners           int64                   `json:"numberOwners"`
	NumberOwnersPercentage float64                 `json:"numberOwnersPercentage"`
	FloorPrice             string                  `json:"floorPrice"`
	FloorPriceOneDay       *VolumneObject          `json:"floorPriceOneDay"`
	FloorPriceOneWeek      *VolumneObject          `json:"floorPriceOneWeek"`
	VolumeFifteenMinutes   *VolumneObject          `json:"volumeFifteenMinutes"`
	VolumeOneDay           *VolumneObject          `json:"volumeOneDay"`
	VolumeOneWeek          *VolumneObject          `json:"volumeOneWeek"`
	ProjectMarketplaceData *ProjectMarketplaceData `json:"projectMarketplaceData"`
	Owner                  *OwnerInfo              `json:"owner"`
	IsHidden               bool                    `json:"isHidden"`
	TotalVolume            uint64                  `json:"totalVolume"`
	IsBuyable              bool                    `json:"isBuyable"`
	Priority               int                     `json:"priority"`
}

type ProjectMarketplaceData struct {
	Listed          uint64  `json:"listed"`
	FloorPrice      uint64  `json:"floor_price"`
	TotalVolume     uint64  `json:"volume"`
	MintVolume      uint64  `json:"mint_volume"`
	CEXVolume       uint64  `json:"cex_volume"`
	FirstSaleVolume float64 `json:"first_sale_volume"`
}

type OwnerInfo struct {
	WalletAddress           string `json:"walletAddress,omitempty"`
	WalletAddressPayment    string `json:"walletAddress_payment,omitempty"`
	WalletAddressBTC        string `json:"walletAddress_btc,omitempty"`
	WalletAddressBTCTaproot string `json:"walletAddress_btc_taproot,omitempty"`
	DisplayName             string `json:"displayName,omitempty"`
	Avatar                  string `json:"avatar"`
}

type ProjectMintingInfo struct {
	Index        int64 `json:"index"`
	IndexReverse int64 `json:"index_reverse"`
}

type ProjectInfo struct {
	Name            string              `json:"name"`
	TokenId         string              `json:"tokenId"`
	Thumbnail       string              `json:"thumbnail"`
	ContractAddress string              `json:"contractAddress"`
	CreatorAddress  string              `json:"creatorAddress"`
	MaxSupply       int64               `json:"maxSupply"`
	MintingInfo     *ProjectMintingInfo `json:"mintingInfo"`
	IsMintedOut     bool                `json:"isMintedOut"`
}

type VolumneObject struct {
	Amount            string  `json:"amount"`
	PercentageChanged float64 `json:"percentageChanged"`
}

type DexBtcListingAlgolia struct {
	ObjectID string `json:"objectID"`
	model.DexBtcListing
}

type UserAlgolia struct {
	ObjectID                string        `json:"objectID"`
	WalletAddress           string        `json:"walletAddress"`
	WalletAddressPayment    string        `json:"walletAddressPayment"`
	WalletAddressBTC        string        `json:"walletAddressBtc,omitempty"`
	WalletAddressBTCTaproot string        `json:"walletAddressBtcTaproot,omitempty"`
	DisplayName             string        `json:"displayName"`
	Bio                     string        `json:"bio"`
	Avatar                  string        `json:"avatar"`
	CreatedAt               *time.Time    `json:"createdAt"`
	ProfileSocial           ProfileSocial `json:"profileSocial"`
	Stats                   UserStats     `json:"stats"`
}

type UserStats struct {
	CollectionCreated int32   `json:"collectionCreated"`
	NftMinted         int32   `json:"nftMinted"`
	OutputMinted      int32   `json:"outputMinted"`
	VolumeMinted      float64 `json:"volumeMinted"`
}

type ProfileSocial struct {
	Web             string `json:"web"`
	Twitter         string `json:"twitter"`
	Discord         string `json:"discord"`
	Medium          string `json:"medium"`
	Instagram       string `json:"instagram"`
	EtherScan       string `json:"etherScan"`
	TwitterVerified bool   `json:"twitterVerified"`
}

type ProjectAlgolia struct {
	ObjectID        string     `json:"objectID"`
	TokenID         string     `json:"tokenID"`
	Name            string     `json:"name"`
	CreatorAddrr    string     `json:"creatorAddrr"`
	CreatorName     string     `json:"creatorName"`
	DeletedAt       *time.Time `json:"deletedAt"`
	Image           string     `json:"image"`
	ContractAddress string     `json:"contractAddress"`
	ItemDesc        string     `json:"itemDesc"`
	Index           int64      `json:"index"`
	IndexReverse    int64      `json:"indexReverse"`
	MintPrice       string     `json:"mintPrice"`
	MaxSupply       int64      `json:"maxSupply"`
	Categories      []string   `json:"categories"`
	Status          bool       `json:"status"`
	IsHidden        bool       `json:"isHidden"`
	IsSynced        bool       `json:"isSynced"`
}

type TokenUriAlgolia struct {
	ObjectID         string `json:"objectID"`
	TokenID          string `json:"tokenID"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Image            string `json:"image"`
	InscriptionIndex string `json:"inscriptionIndex"`

	ProjectName string `json:"projectName"`
	ProjectID   string `json:"projectID"`
	Thumbnail   string `json:"thumbnail"`
}
