package entity

import "time"

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
	IsSynced        bool       `json:"isSynced"`
	DeletedAt       *time.Time `json:"deletedAt"`
	Image           string     `json:"image"`
	ContractAddress string     `json:"contractAddress"`
	ItemDesc        string     `json:"itemDesc"`
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
