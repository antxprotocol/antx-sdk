package types

import (
	exchangetypes "github.com/antxprotocol/antx-proto/gen/go/antx/chain/exchange"
	ordertypes "github.com/antxprotocol/antx-proto/gen/go/antx/chain/order"
	pricetypes "github.com/antxprotocol/antx-proto/gen/go/antx/chain/price"
)

// =============================== Base Response Types ===============================
// BaseResp is already defined in base.go, not repeated here

// =============================== Account Related Types ===============================

// GetAccountNumberAndSequenceResponseData get account number and sequence response data
type GetAccountNumberAndSequenceResponseData struct {
	Exist         bool   `json:"exist"`
	AccountNumber string `json:"accountNumber"`
	Sequence      string `json:"sequence"`
}

// GetAccountNumberAndSequenceResponse get account number and sequence response
type GetAccountNumberAndSequenceResponse struct {
	BaseResp
	Data GetAccountNumberAndSequenceResponseData `json:"data"`
}

// =============================== Subaccount Related Types ===============================

// GetSubaccountListResponse get subaccount list response
type GetSubaccountListResponse struct {
	BaseResp
	Data GetSubaccountListResponseData `json:"data"`
}

// GetSubaccountListResponseData get subaccount list response data
type GetSubaccountListResponseData struct {
	SubaccountList []Subaccount `json:"subaccountList"`
}

// Subaccount subaccount information
type Subaccount struct {
	Id              string         `json:"id"`              // Subaccount ID, must be greater than 0
	ChainType       int32          `json:"chainType"`       // Chain type
	ChainAddress    string         `json:"chainAddress"`    // Chain address
	ClientAccountId string         `json:"clientAccountId"` // Client custom ID, for idempotency check, max length 64
	IsSystemAccount bool           `json:"isSystemAccount"` // Whether it is a system account
	TakerFeeRatePpm uint32         `json:"takerFeeRatePpm"` // Taker fee rate, unit: parts per million
	MakerFeeRatePpm uint32         `json:"makerFeeRatePpm"` // Maker fee rate, unit: parts per million
	TradeSetting    []TradeSetting `json:"tradeSetting"`    // Perpetual contract trading settings
}

// TradeSetting trading settings
type TradeSetting struct {
	ExchangeId string `json:"exchangeId"` // Exchange ID
	MarginMode uint32 `json:"marginMode"` // Margin mode 0: Unknown 1: Cross 2: Isolated
	Leverage   uint32 `json:"leverage"`   // Leverage multiplier
}

// =============================== Coin Related Types ===============================

// GetCoinListResponse get coin list response
type GetCoinListResponse struct {
	BaseResp
	Data GetCoinListRespData `json:"data,omitempty"`
}

// GetCoinListRespData get coin list response data
type GetCoinListRespData struct {
	CoinList []Coin `json:"coinList"`
}

// Coin coin information
type Coin struct {
	Id                   string `json:"id"`                   // Coin ID, range [1001, 9999], incrementally set
	Symbol               string `json:"symbol"`               // Coin symbol, max length 32, e.g., 'BTC'
	StepSizeScale        int32  `json:"stepSizeScale"`        // Decimal places for minimum unit quantity, i.e., minimum unit = 10^-step_size_scale, e.g., step_size_scale=6 means 0.000001, range [-100, 100]
	AssetChainId         string `json:"assetChainId"`         // Asset chain ID
	AssetContractAddress string `json:"assetContractAddress"` // Asset contract address
}

// =============================== Exchange Related Types ===============================

// GetExchangeListResponse get exchange list response
type GetExchangeListResponse struct {
	BaseResp
	Data GetExchangeListRespData `json:"data,omitempty"`
}

// GetExchangeListRespData get exchange list response data
type GetExchangeListRespData struct {
	ExchangeList []Exchange `json:"exchangeList"`
}

// Exchange exchange information
type Exchange struct {
	Id                    string    `json:"id"`                    // Exchange ID, must be incrementally set, spot range [100001, 109999], perpetual range [200001, 209999]
	Symbol                string    `json:"symbol"`                // Exchange symbol, e.g., BTC/USDT, BTC-USDT
	BaseCoinId            string    `json:"baseCoinId"`            // Base coin ID, e.g., BTC
	QuoteCoinId           string    `json:"quoteCoinId"`           // Quote coin ID, e.g., USDT
	StepSizeScale         int32     `json:"stepSizeScale"`         // Decimal places for minimum position unit, i.e., minimum unit = 10^-step_size_scale, e.g., step_size_scale=6 means 0.000001
	TickSizeScale         int32     `json:"tickSizeScale"`         // Decimal places for minimum price unit, i.e., minimum unit = 10^-tick_size_scale, e.g., tick_size_scale=-1 means 10
	OrderPriceMaxRatioPpm uint32    `json:"orderPriceMaxRatioPpm"` // Maximum order price limit ratio (compared to oracle price), unit: parts per million
	OrderPriceMinRatioPpm uint32    `json:"orderPriceMinRatioPpm"` // Minimum order price limit ratio (compared to oracle price), unit: parts per million
	OrderSizeMax          string    `json:"orderSizeMax"`          // Maximum order size
	Perpetual             Perpetual `json:"perpetual,omitempty"`   // Perpetual contract trading information
}

// Perpetual perpetual contract information
type Perpetual struct {
	SupportMarginModeList       []uint32   `json:"supportMarginModeList"`       // Supported margin modes
	RiskTierList                []RiskTier `json:"riskTierList"`                // Risk tier list (must be sorted by max leverage descending, maintenance margin ratio ascending, position value ascending)
	LiquidateFeeRatePpm         uint32     `json:"liquidateFeeRatePpm"`         // Default liquidation fee rate, unit: parts per million
	DefaultLeverage             uint32     `json:"defaultLeverage"`             // Default leverage multiplier
	EnableOrderCreate           bool       `json:"enableOrderCreate"`           // Whether order creation is allowed
	EnableOrderFill             bool       `json:"enableOrderFill"`             // Whether order fill is allowed
	EnablePositionOpen          bool       `json:"enablePositionOpen"`          // Whether position opening is allowed
	FundingInterestRatePpm      uint32     `json:"fundingInterestRatePpm"`      // Funding interest rate
	FundingImpactMarginNotional string     `json:"fundingImpactMarginNotional"` // Funding impact margin notional
	FundingRateAbsMaxPpm        uint32     `json:"fundingRateAbsMaxPpm"`        // Maximum absolute funding rate
	FundingRateIntervalMinutes  uint32     `json:"fundingRateIntervalMinutes"`  // Funding rate calculation interval
}

// RiskTier risk tier
type RiskTier struct {
	MaxLeverage               uint32 `json:"maxLeverage"`               // Maximum leverage multiplier
	MaintenanceMarginRatioPpm uint32 `json:"maintenanceMarginRatioPpm"` // Maintenance margin ratio
	PositionValueUpperBound   string `json:"positionValueUpperBound"`   // Position value upper bound
}

// =============================== Trading Related Types ===============================

// SendRawTxRequest send raw transaction request
type SendRawTxRequest struct {
	TypeURL       string `json:"typeUrl"`
	RawTx         string `json:"rawTx"`
	AccountNumber uint64 `json:"accountNumber"`
}

// SendRawTxResponse send raw transaction response
type SendRawTxResponse struct {
	BaseResp
	Data SendRawTxResponseData `json:"data"`
}

// SendRawTxResponseData send raw transaction response data
type SendRawTxResponseData struct {
	TxHash     string `json:"txHash"`     // Transaction hash
	RawTx      string `json:"rawTx"`      // Transaction raw tx
	ResultData string `json:"resultData"` // Data
	// Possible alternative field names
	Hash string `json:"hash"` // Transaction hash (alternative field name)
	TxID string `json:"txId"` // Transaction ID (alternative field name)
}

// SendSyncTransactionResponse send sync transaction response
type SendSyncTransactionResponse struct {
	BaseResp
	Data string `json:"data"`
}

// =============================== Blockchain Explorer Related Types ===============================

// GetTransactionResultRequest get transaction result request
type GetTransactionResultRequest struct {
	Hash string `json:"hash"`
}

// GetTransactionResultResponse get transaction result response
type GetTransactionResultResponse struct {
	Code         string                       `json:"code"`
	Msg          string                       `json:"msg"`
	Params       map[string]interface{}       `json:"params,omitempty"`
	RequestTime  string                       `json:"requestTime,omitempty"`
	ResponseTime string                       `json:"responseTime,omitempty"`
	TraceId      string                       `json:"traceId,omitempty"`
	Data         GetTransactionResultRespData `json:"data,omitempty"`
}

// GetTransactionResultRespData get transaction result response data
type GetTransactionResultRespData struct {
	RawTx      string             `json:"rawTx"`      // Raw data
	Block      uint64             `json:"block"`      // Block height
	Hash       string             `json:"hash"`       // Transaction hash
	From       string             `json:"from"`       // Sender
	Status     bool               `json:"status"`     // Status
	Error      interface{}        `json:"error"`      // Error
	ActionList []ExplorerTxAction `json:"action"`     // Actions
	ResultData string             `json:"resultData"` // Data
}

// ExplorerTxAction blockchain explorer transaction action
type ExplorerTxAction struct {
	TypeUrl string      `json:"typeUrl"` // Type
	Detail  interface{} `json:"detail"`  // Details
}

// =============================== Order Related Types ===============================

// CreateOrderParam create order parameter
type CreateOrderParam struct {
	AgentAddress          string
	SubaccountId          uint64
	ExchangeId            uint64
	MarginMode            exchangetypes.MarginMode
	Leverage              uint32
	IsBuy                 bool
	PriceScale            int32
	PriceValue            uint64
	SizeScale             int32
	SizeValue             uint64
	ClientOrderId         string
	TimeInForce           ordertypes.TimeInForce
	ReduceOnly            bool
	ExpireTime            uint64
	IsMarket              bool
	IsPositionTp          bool
	IsPositionSl          bool
	TriggerType           ordertypes.TriggerType
	TriggerPriceType      pricetypes.PriceType
	TriggerPriceValue     uint64
	OpenTpslParentOrderId uint64
	IsSetOpenTp           bool
	OpenTpParam           ordertypes.OpenTpSlParam
	IsSetOpenSl           bool
	OpenSlParam           ordertypes.OpenTpSlParam
}

// CreateOrderBatchParam create order batch parameter
type CreateOrderBatchParam struct {
	AgentAddress     string
	SubaccountId     uint64
	ExchangeId       uint64
	MarginMode       exchangetypes.MarginMode
	Leverage         uint32
	CreateOrderParam []*CreateOrderBatchDetail
}

// CreateOrderBatchDetail create order batch detail
type CreateOrderBatchDetail struct {
	IsBuy             bool
	PriceScale        int32
	PriceValue        uint64
	SizeScale         int32
	SizeValue         uint64
	ClientOrderId     string
	TimeInForce       ordertypes.TimeInForce
	ReduceOnly        bool
	ExpireTime        uint64
	IsMarket          bool
	IsPositionTp      bool
	IsPositionSl      bool
	TriggerType       ordertypes.TriggerType
	TriggerPriceType  pricetypes.PriceType
	TriggerPriceValue uint64
	IsSetOpenTp       bool
	OpenTpParam       ordertypes.OpenTpSlParam
	IsSetOpenSl       bool
	OpenSlParam       ordertypes.OpenTpSlParam
}

// CancelOrderParam cancel order parameter
type CancelOrderParam struct {
	AgentAddress string
	SubaccountId uint64
	OrderIdList  []uint64
}

// CancelOrderByClientIdParam cancel order by client ID parameter
type CancelOrderByClientIdParam struct {
	AgentAddress      string
	SubaccountId      uint64
	ClientOrderIdList []string
}

// CancelAllOrderParam cancel all orders parameter
type CancelAllOrderParam struct {
	AgentAddress         string
	SubaccountId         uint64
	FilterExchangeIdList []uint64
}

// CloseAllPositionParam close all positions parameter
type CloseAllPositionParam struct {
	AgentAddress         string
	SubaccountId         uint64
	FilterExchangeIdList []uint64
}
