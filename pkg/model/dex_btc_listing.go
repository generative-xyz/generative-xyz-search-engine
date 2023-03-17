package model

import "time"

type DexBtcListing struct {
	Model         `bson:"inline"`
	RawPSBT       string     `bson:"raw_psbt" json:"rawPsbt"`
	SplitTx       string     `bson:"split_tx" json:"splitTx"`
	InscriptionID string     `bson:"inscription_id" json:"inscriptionId"`
	Amount        uint64     `bson:"amount" json:"amount"`
	SellerAddress string     `bson:"seller_address" json:"sellerAddress"`
	Verified      bool       `bson:"verified" json:"verified"`
	CancelAt      *time.Time `bson:"cancel_at" json:"cancelAt"`
	Cancelled     bool       `bson:"cancelled" json:"cancelled"`
	CancelTx      string     `bson:"cancel_tx" json:"cancelTx"`
	Inputs        []string   `bson:"inputs" json:"inputs"`
	Matched       bool       `bson:"matched" json:"matched"`
	MatchedTx     string     `bson:"matched_tx" json:"matchedTx"`
	MatchAt       *time.Time `bson:"matched_at" json:"matchAt"`
	Buyer         string     `bson:"buyer" json:"buyer"`
}

type TokenUriListingVolume struct {
	TotalAmount uint64 `bson:"totalAmount" json:"totalAmount"`
}

type MarketplaceBTCListingFloorPrice struct {
	ID    string `bson:"_id"`
	Price uint64 `bson:"amount"`
}
