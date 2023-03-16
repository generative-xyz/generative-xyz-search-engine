package model

import "time"

type User struct {
	Model                   `bson:"inline"`
	WalletAddress           string        `bson:"wallet_address" json:"walletAddress,omitempty"`                       // eth wallet define user in platform by connect wallet and sign
	WalletAddressPayment    string        `bson:"wallet_address_payment" json:"walletAddressPayment,omitempty"`        // eth wallet artist receive royalty
	WalletAddressBTC        string        `bson:"wallet_address_btc" json:"walletAddressBtc,omitempty"`                // btc wallet artist receive royalty
	WalletAddressBTCTaproot string        `bson:"wallet_address_btc_taproot" json:"walletAddressBtcTaproot,omitempty"` // btc wallet receive minted nft
	DisplayName             string        `bson:"display_name" json:"displayName,omitempty"`
	Bio                     string        `bson:"bio" json:"bio,omitempty"`
	Avatar                  string        `bson:"avatar" json:"avatar"`
	CreatedAt               *time.Time    `bson:"created_at" json:"createdAt,omitempty"`
	ProfileSocial           ProfileSocial `json:"profileSocial,omitempty" bson:"profile_social"`
	Stats                   UserStats     `bson:"stats" json:"stats"`
	IsAdmin                 bool          `bson:"isAdmin" json:"isAdmin"`
}

type UserStats struct {
	CollectionCreated int32   `bson:"collection_created" json:"collectionCreated"`
	NftMinted         int32   `bson:"nft_minted" json:"nftMinted"`
	OutputMinted      int32   `bson:"output_minted" json:"outputMinted"`
	VolumeMinted      float64 `bson:"volume_minted" json:"volumeMinted"`
}

type ProfileSocial struct {
	Web             string `bson:"web" json:"web,omitempty"`
	Twitter         string `bson:"twitter" json:"twitter,omitempty"`
	Discord         string `bson:"discord" json:"discord,omitempty"`
	Medium          string `bson:"medium" json:"medium,omitempty"`
	Instagram       string `bson:"instagram" json:"instagram,omitempty"`
	EtherScan       string `bson:"etherScan" json:"ether_scan,omitempty"`
	TwitterVerified bool   `bson:"twitter_verified" json:"twitterVerified,omitempty"`
}
