package models

type CopyQuery struct {
	SortQuery
	Count uint `json:"count"`
}

type SortQuery struct {
	SystemName string `json:"systemName"`
	Where      string `json:"where"`
}

type AddSystemQuery struct {
	SystemName string `json:"systemName"`
	Type       bool   `json:"type"`
}
