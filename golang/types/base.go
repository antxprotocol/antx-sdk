package types

// BaseResp base response structure
type BaseResp struct {
	Code string `json:"code"` // Response code
	Msg  string `json:"msg"`  // Response message
}

// IndexerPageOffsetData pagination offset data
type IndexerPageOffsetData struct {
	CreateTime string `json:"createTime"` // Next page offset data, creation time
	ItemId     string `json:"itemId"`     // Next page offset data, itemId
}
