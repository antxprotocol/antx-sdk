package constants

// =============================== API Path Constants ===============================

const (
	// Base API path
	BaseAPIPath = "/api/v1"

	// Address related
	GetAddressInfoPath = BaseAPIPath + "/address/getAddressInfo"

	// Subaccount related
	GetSubaccountPath = BaseAPIPath + "/subaccount/getSubaccount"

	// Trading related
	GetCoinListPath         = BaseAPIPath + "/trade/getCoinList"
	GetExchangeListPath     = BaseAPIPath + "/trade/getExchangeList"
	SendTransactionPath     = BaseAPIPath + "/trade/sendTransaction"
	SendSyncTransactionPath = BaseAPIPath + "/trade/sendSyncTransaction"

	// Market data related
	GetKlinePath          = BaseAPIPath + "/trade/getKline"
	GetTickerPath         = BaseAPIPath + "/trade/getTicker"
	GetDepthPath          = BaseAPIPath + "/trade/getDepth"
	GetTradePath          = BaseAPIPath + "/trade/getTrade"
	GetFundingHistoryPath = BaseAPIPath + "/trade/getFundingHistory"
	GetPricePath          = BaseAPIPath + "/trade/getPrice"

	// Trading query related
	GetActiveOrderPath                 = BaseAPIPath + "/trade/getActiveOrder"
	GetHistoryOrderPath                = BaseAPIPath + "/trade/getHistoryOrder"
	GetPerpetualAccountAssetPath       = BaseAPIPath + "/trade/getPerpetualAccountAsset"
	GetPositionTransactionPath         = BaseAPIPath + "/trade/getPositionTransaction"
	GetCollateralTransactionPath       = BaseAPIPath + "/trade/getCollateralTransaction"
	GetAssetSnapshotPath               = BaseAPIPath + "/trade/getAssetSnapshot"
	GetHistoryOrderFillTransactionPath = BaseAPIPath + "/trade/getHistoryOrderFillTransaction"
	GetHistoryPositionTermPath         = BaseAPIPath + "/trade/getHistoryPositionTerm"

	// Blockchain explorer related
	GetTransactionPath = BaseAPIPath + "/explorer/tx"

	// WebSocket related
	WebSocketPath = "/api/v1/ws"
)

// =============================== K-line Type Constants ===============================

const (
	KlineTypeMinute1  = "MINUTE_1"  // 1 minute
	KlineTypeMinute5  = "MINUTE_5"  // 5 minutes
	KlineTypeMinute15 = "MINUTE_15" // 15 minutes
	KlineTypeMinute30 = "MINUTE_30" // 30 minutes
	KlineTypeHour1    = "HOUR_1"    // 1 hour
	KlineTypeHour2    = "HOUR_2"    // 2 hours
	KlineTypeHour4    = "HOUR_4"    // 4 hours
	KlineTypeHour6    = "HOUR_6"    // 6 hours
	KlineTypeHour8    = "HOUR_8"    // 8 hours
	KlineTypeHour12   = "HOUR_12"   // 12 hours
	KlineTypeDay1     = "DAY_1"     // 1 day
	KlineTypeWeek1    = "WEEK_1"    // 1 week
	KlineTypeMonth1   = "MONTH_1"   // 1 month
)

// =============================== Price Type Constants ===============================

const (
	PriceTypeLast    = "PRICE_TYPE_LAST"     // Latest price
	PriceTypeAskBest = "PRICE_TYPE_ASK_BEST" // Best ask price
	PriceTypeBidBest = "PRICE_TYPE_BID_BEST" // Best bid price
	PriceTypeMark    = "PRICE_TYPE_MARK"     // Mark price
	PriceTypeOracle  = "PRICE_TYPE_ORACLE"   // Oracle price
)

// =============================== Order Status Constants ===============================

const (
	OrderStatusUnknown         = 0 // Unknown
	OrderStatusPending         = 1 // Pending
	OrderStatusFilled          = 2 // Filled
	OrderStatusCancelled       = 3 // Cancelled
	OrderStatusExpired         = 4 // Expired
	OrderStatusRejected        = 5 // Rejected
	OrderStatusPartiallyFilled = 6 // Partially filled
	OrderStatusLiquidated      = 7 // Liquidated
	OrderStatusDeleveraged     = 8 // Deleveraged
)

// =============================== Transaction Message Type Constants ===============================

const (
	// Order related message types
	MsgCreateOrderTypeURL           = "/antx.chain.order.MsgCreateOrder"
	MsgCreateOrderBatchTypeURL      = "/antx.chain.order.MsgCreateOrderBatch"
	MsgCancelOrderTypeURL           = "/antx.chain.order.MsgCancelOrder"
	MsgCancelOrderByClientIdTypeURL = "/antx.chain.order.MsgCancelOrderByClientId"
	MsgCancelAllOrderTypeURL        = "/antx.chain.order.MsgCancelAllOrder"
	MsgCloseAllPositionTypeURL      = "/antx.chain.order.MsgCloseAllPosition"

	// Agent related message types
	MsgBindAgentTypeURL = "/antx.chain.agent.MsgBindAgent"
)
