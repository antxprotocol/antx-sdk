package types

import "encoding/json"

// =============================== Market Data Related Structures ===============================

// KLine K-line data
type KLine struct {
	KlineId       string `json:"klineId"`       // K-line ID
	ExchangeId    string `json:"exchangeId"`    // Exchange ID
	KlineType     string `json:"klineType"`     // K-line type
	PriceType     string `json:"priceType"`     // K-line price type
	KlineTime     uint64 `json:"klineTime"`     // K-line time
	Trades        string `json:"trades"`        // Number of trades
	Size          string `json:"size"`          // Volume
	Value         string `json:"value"`         // Turnover
	High          string `json:"high"`          // Highest price
	Low           string `json:"low"`           // Lowest price
	Open          string `json:"open"`          // Open price
	Close         string `json:"close"`         // Close price
	MakerBuySize  string `json:"makerBuySize"`  // Maker buy volume
	MakerBuyValue string `json:"makerBuyValue"` // Maker buy turnover
}

// TickerData Ticker data
type TickerData struct {
	ExchangeId         string `json:"exchangeId"`         // Exchange ID
	LastPrice          string `json:"lastPrice"`          // Last price
	MarkPrice          string `json:"markPrice"`          // Mark price
	IndexPrice         string `json:"indexPrice"`         // Index price
	OraclePrice        string `json:"oraclePrice"`        // Oracle price
	PriceChange        string `json:"priceChange"`        // Price change
	PriceChangePercent string `json:"priceChangePercent"` // Price change percentage
	High               string `json:"high"`               // 24h highest price
	Low                string `json:"low"`                // 24h lowest price
	Open               string `json:"open"`               // Open price
	Close              string `json:"close"`              // Close price
	Size               string `json:"size"`               // Volume
	Value              string `json:"value"`              // 24h turnover
	OpenInterest       string `json:"openInterest"`       // Open interest
	FundingRate        string `json:"fundingRate"`        // Funding rate
	FundingTime        string `json:"fundingTime"`        // Funding rate time
	NextFundingTime    string `json:"nextFundingTime"`    // Next funding rate time
	StartTime          string `json:"startTime"`          // Start time
	EndTime            string `json:"endTime"`            // End time
	HighTime           string `json:"highTime"`           // Highest price time
	LowTime            string `json:"lowTime"`            // Lowest price time
	Trades             string `json:"trades"`             // Number of trades
}

// DepthData depth data
type DepthData struct {
	ExchangeId  string      `json:"exchangeId"`  // Exchange ID
	Bids        []BookOrder `json:"bids"`        // Buy order list
	Asks        []BookOrder `json:"asks"`        // Sell order list
	UpdatedTime uint64      `json:"updatedTime"` // Updated time
}

// BookOrder order book order
type BookOrder struct {
	Price string `json:"price"` // Price
	Size  string `json:"size"`  // Size
}

// Ticket trade data
type Ticket struct {
	ExchangeId string `json:"exchangeId"` // Exchange ID
	Price      string `json:"price"`      // Trade price
	Size       string `json:"size"`       // Trade size
	Value      string `json:"value"`      // Trade value
	IsBuy      bool   `json:"isBuy"`      // Whether it is a buy order
	Time       string `json:"time"`       // Trade time
}

// FundingRate funding rate
type FundingRate struct {
	ExchangeId   string `json:"exchangeId"`   // Exchange ID
	FundingRate  string `json:"fundingRate"`  // Funding rate
	OraclePrice  string `json:"oraclePrice"`  // Oracle price
	IndexPrice   string `json:"indexPrice"`   // Index price
	FundingTime  uint64 `json:"fundingTime"`  // Funding rate time
	IsSettlement bool   `json:"isSettlement"` // Whether it is a settlement
	UpdatedTime  uint64 `json:"updatedTime"`  // Updated time
}

// Price price data
type Price struct {
	ExchangeId  string `json:"exchangeId"`  // Exchange ID
	Price       string `json:"price"`       // Price
	PriceTime   uint64 `json:"priceTime"`   // Price time
	CreatedTime uint64 `json:"createdTime"` // Created time
}

// =============================== Request and Response Structures ===============================

// GetKLineReq get K-line information request
type GetKLineReq struct {
	ExchangeId                    string `form:"exchangeId"`                             // Exchange ID
	KlineType                     string `form:"klineType"`                              // K-line type
	PriceType                     string `form:"priceType"`                              // Price type
	Size                          uint32 `form:"size,optional,default=100"`              // Number of records, default 100
	OffsetData                    string `form:"offsetData,optional"`                    // Pagination offset, if empty, get first page
	FilterBeginKlineTimeInclusive int64  `form:"filterBeginKlineTimeInclusive,optional"` // Start time, if empty, get oldest data
	FilterEndKlineTimeExclusive   int64  `form:"filterEndKlineTimeExclusive,optional"`   // End time, if empty, get latest data
}

// GetKLineRespData get K-line information response data
type GetKLineRespData struct {
	KlineList          []KLine `json:"klineList"`          // K-line list
	NextPageOffsetData string  `json:"nextPageOffsetData"` // Next page offset, if no next page, empty string
}

// GetKLineResp get K-line information response
type GetKLineResp struct {
	BaseResp
	Data GetKLineRespData `json:"data,omitempty"`
}

// GetFundingHistoryReq get funding rate history request
type GetFundingHistoryReq struct {
	ExchangeId                  string `form:"exchangeId"`                           // Exchange ID
	Size                        uint32 `form:"size,optional,default=100"`            // Number of records, default 100
	OffsetData                  string `form:"offsetData,optional"`                  // Pagination offset, if empty, get first page
	FilterSettlementFundingRate bool   `form:"filterSettlementFundingRate,optional"` // Whether to only get settlement funding rates
	FilterBeginTimeInclusive    uint64 `form:"filterBeginTimeInclusive,optional"`    // Start time, if empty, get oldest data
	FilterEndTimeExclusive      uint64 `form:"filterEndTimeExclusive,optional"`      // End time, if empty, get latest data
}

// GetFundingHistoryRespData get funding rate history response data
type GetFundingHistoryRespData struct {
	FundingRateList    []FundingRate `json:"fundingRateList"`    // Funding rate list
	NextPageOffsetData string        `json:"nextPageOffsetData"` // Next page offset, if no next page, empty string
}

// GetFundingHistoryResp get funding rate history response
type GetFundingHistoryResp struct {
	BaseResp
	Data GetFundingHistoryRespData `json:"data,omitempty"`
}

// =============================== Helper Methods ===============================

// =============================== Helper Methods ===============================

// ToJSON converts to JSON string
func (k *KLine) ToJSON() (string, error) {
	data, err := json.Marshal(k)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToJSON converts to JSON string
func (t *TickerData) ToJSON() (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToJSON converts to JSON string
func (d *DepthData) ToJSON() (string, error) {
	data, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToJSON converts to JSON string
func (t *Ticket) ToJSON() (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToJSON converts to JSON string
func (f *FundingRate) ToJSON() (string, error) {
	data, err := json.Marshal(f)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToJSON converts to JSON string
func (p *Price) ToJSON() (string, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
