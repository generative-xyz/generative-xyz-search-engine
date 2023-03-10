package model

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
