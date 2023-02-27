package model

import "time"

type Project struct {
	Model                  `bson:"inline"`
	ContractAddress        string `bson:"contractAddress" json:"contractAddress"`
	TokenID                string `bson:"tokenid" json:"tokenID"`
	MaxSupply              int64  `bson:"maxSupply" json:"maxSupply"`
	MintPrice              string `bson:"mintPrice" json:"mintPrice"`
	MintPriceEth           string
	NetworkFeeEth          string
	NetworkFee             string             `bson:"networkFee" json:"networkFee"`
	Name                   string             `bson:"name" json:"name"`
	CreatorName            string             `bson:"creatorName" json:"creatorName"`
	CreatorAddrr           string             `bson:"creatorAddress" json:"creatorAddrr"`
	CreatorAddrrBTC        string             `bson:"creatorAddrrBTC" json:"creatorAddrrBTC"`
	Description            string             `bson:"description" json:"description"`
	OpenMintUnixTimestamp  int                `bson:"openMintUnixTimestamp" json:"openMintUnixTimestamp"`
	CloseMintUnixTimestamp int                `bson:"closeMintUnixTimestamp" json:"closeMintUnixTimestamp"`
	Thumbnail              string             `bson:"thumbnail" json:"thumbnail"`
	ReservationList        []string           `bson:"reservationList" json:"reservationList"`
	MintFee                int                `bson:"mintFee" json:"mintFee"`
	TokenDescription       string             `bson:"tokenDescription" json:"tokenDescription"`
	Styles                 string             `bson:"styles" json:"styles"`
	Royalty                int                `bson:"royalty" json:"royalty"`
	SocialWeb              string             `bson:"socialWeb" json:"socialWeb"`
	SocialTwitter          string             `bson:"socialTwitter" json:"socialTwitter"`
	SocialDiscord          string             `bson:"socialDiscord" json:"socialDiscord"`
	SocialMedium           string             `bson:"socialMedium" json:"socialMedium"`
	SocialInstagram        string             `bson:"socialInstagram" json:"socialInstagram"`
	License                string             `bson:"license" json:"license"`
	GenNFTAddr             string             `bson:"genNFTAddr" json:"genNFTAddr"`
	MintTokenAddress       string             `bson:"mintTokenAddress" json:"mintTokenAddress"`
	Tags                   []string           `bson:"tags" json:"tags"`
	Categories             []string           `bson:"categories" json:"categories"`
	Status                 bool               `bson:"status" json:"status"`
	IsSynced               bool               `bson:"isSynced" json:"isSynced"`
	MintingInfo            ProjectMintingInfo `bson:"inline" json:"mintingInfo"`
	CompleteTime           int64              `bson:"completeTime" json:"completeTime"`
	// CreatorProfile          User               `bson:"creatorProfile" json:"creatorProfile"`
	BlockNumberMinted *string     `bson:"block_number_minted" json:"block_number_minted"`
	MintedTime        *time.Time  `bson:"minted_time" json:"minted_time"`
	Stats             ProjectStat `bson:"stats" json:"stats"`
	TraitsStat        []TraitStat `bson:"traitsStat" json:"traitsStat"`
	Priority          *int        `bson:"priority" json:"priority"`
	IsHidden          bool        `bson:"isHidden" json:"isHidden"`
	//if user uses links instead of animation URL
	WhiteListEthContracts []string         `bson:"whiteListEthContracts" json:"whiteListEthContracts"` //if user uses links instead of animation URL
	IsFullChain           bool             `bson:"isFullChain" json:"isFullChain"`
	ReportUsers           []*ReportProject `bson:"reportUsers" json:"reportUsers"`
	InscriptionIcon       string           `bson:"inscription_icon" json:"inscriptionIcon"`
	Source                string           `bson:"source" json:"source"`
}

type ReportProject struct {
	OriginalLink      string `bson:"originalLink" json:"originalLink"`
	ReportUserAddress string `bson:"reportUserAddress" json:"reportUserAddress"`
}

type ProjectMintingInfo struct {
	Index        int64 `bson:"index"`
	IndexReverse int64 `bson:"indexReverse"`
}

type TraitStat struct {
	TraitName       string           `bson:"traitName" json:"traitName"`
	TraitValuesStat []TraitValueStat `bson:"traitValuesStat" json:"traitValuesStat"`
}

type ProjectStat struct {
	LastTimeSynced   *time.Time `bson:"lastTimeSynced" json:"lastTimeSynced"`
	UniqueOwnerCount uint32     `bson:"uniqueOwnerCount" json:"uniqueOwnerCount"`
	// TODO add other stats here
	TotalTradingVolumn string `bson:"totalTradingVolumn" json:"totalTradingVolumn"`
	FloorPrice         string `bson:"floorPrice" json:"floorPrice"`
	BestMakeOfferPrice string `bson:"bestMakeOfferPrice" json:"bestMakeOfferPrice"`
	ListedPercent      int32  `bson:"listedPercent" json:"listedPercent"`
	MintedCount        int32  `bson:"minted_count" json:"mintedCount"`
	TrendingScore      int64  `bson:"trending_score" json:"trendingScore"`
}

type TraitValueStat struct {
	Value  string `bson:"value" json:"value"`
	Rarity int32  `bson:"rarity" json:"rarity"`
}
