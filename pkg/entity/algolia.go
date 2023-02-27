package entity

import "time"

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
	AnimationURL     string `json:"animationUrl"`
	InscriptionIndex string `json:"inscriptionIndex"`
	PriceBTC         string `json:"priceBTC"`
	OrderID          string `json:"orderID"`
	ProjectName      string `json:"projectName"`
	ProjectID        string `json:"projectID"`
	Thumbnail        string `json:"thumbnail"`
	Priority         int    `json:"priority"`
}
