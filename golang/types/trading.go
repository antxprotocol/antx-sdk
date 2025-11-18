package types

// =============================== Trading Query Related Structures ===============================

// Order order
type Order struct {
	Id                           string        `json:"id"`                           // Order ID
	SubaccountId                 string        `json:"subaccountId"`                 // Subaccount ID
	CoinId                       string        `json:"coinId"`                       // Trading coin ID
	ExchangeId                   string        `json:"exchangeId"`                   // Exchange ID
	IsBuy                        bool          `json:"isBuy"`                        // Whether it is a buy order
	Price                        string        `json:"price"`                        // Order price, if price=0 then it's a market order
	Size                         string        `json:"size"`                         // Order size
	ClientOrderId                string        `json:"clientOrderId"`                // Client custom ID, for idempotency check, max length 64
	TimeInForce                  uint32        `json:"timeInForce"`                  // Order execution strategy
	ReduceOnly                   bool          `json:"reduceOnly"`                   // Whether it is a reduce-only order
	ExpireTime                   uint64        `json:"expireTime"`                   // Expiration time, unit: milliseconds
	IsPositionTp                 bool          `json:"isPositionTp"`                 // Whether it is a position take-profit/stop-loss order
	IsPositionSl                 bool          `json:"isPositionSl"`                 // Whether it is a position take-profit/stop-loss order
	IsLiquidate                  bool          `json:"isLiquidate"`                  // Whether it is a liquidation order
	IsDeleverage                 bool          `json:"isDeleverage"`                 // Whether it is an auto-deleverage order
	TriggerType                  uint32        `json:"triggerType"`                  // Conditional order trigger type
	TriggerPriceType             uint32        `json:"triggerPriceType"`             // Conditional order trigger price type
	TriggerPrice                 string        `json:"triggerPrice"`                 // Trigger price
	OpenTpSlParentOrderId        string        `json:"openTpSlParentOrderId"`        // Open order ID for open take-profit/stop-loss orders
	IsSetOpenTp                  bool          `json:"isSetOpenTp"`                  // Whether to set open take-profit
	OpenTpParam                  OpenTpSlParam `json:"openTpParam"`                  // Open take-profit parameters, only meaningful when is_set_open_tp=true
	IsSetOpenSl                  bool          `json:"isSetOpenSl"`                  // Whether to set open stop-loss
	OpenSlParam                  OpenTpSlParam `json:"openSlParam"`                  // Open stop-loss parameters, only meaningful when is_set_open_sl=true
	MarginMode                   uint32        `json:"marginMode"`                   // Margin mode when placing order
	Leverage                     uint32        `json:"leverage"`                     // Leverage multiplier when placing order
	TakerFeeRatePpm              uint32        `json:"takerFeeRatePpm"`              // Taker fee rate when placing order, unit: parts per million
	MakerFeeRatePpm              uint32        `json:"makerFeeRatePpm"`              // Maker fee rate when placing order, unit: parts per million
	LiquidateFeeRatePpm          uint32        `json:"liquidateFeeRatePpm"`          // Liquidation fee rate when placing order, unit: parts per million
	AddOrderBookBlockHeight      uint64        `json:"addOrderBookBlockHeight"`      // Block height when order was added to order book, if 0, not triggered yet
	AddOrderBookBlockTime        uint64        `json:"addOrderBookBlockTime"`        // Block time when order was added to order book, if 0, not triggered yet
	AddOrderBookTransactionIndex string        `json:"addOrderBookTransactionIndex"` // Transaction index in block when order was added to order book
	AddOrderBookOperationIndex   string        `json:"addOrderBookOperationIndex"`   // Operation index in transaction when order was added to order book
	Status                       uint32        `json:"status"`                       // Order status
	CancelReason                 uint32        `json:"cancelReason"`                 // Order cancellation reason
	CumFillSize                  string        `json:"cumFillSize"`                  // Cumulative filled size, actual type is decimal
	CumFillValue                 string        `json:"cumFillValue"`                 // Cumulative filled value, actual type is decimal
	CumFillFee                   string        `json:"cumFillFee"`                   // Cumulative filled fee, actual type is decimal
	CumLiquidateFee              string        `json:"cumLiquidateFee"`              // Cumulative liquidation fee, actual type is decimal
	MaxFillPrice                 string        `json:"maxFillPrice"`                 // Maximum fill price for current order, actual type is decimal
	MinFillPrice                 string        `json:"minFillPrice"`                 // Minimum fill price for current order, actual type is decimal
	CumRealizePnl                string        `json:"cumRealizePnl"`                // Cumulative realized PnL, actual type is decimal
	CreatedTime                  uint64        `json:"createdTime"`                  // Created time
	UpdatedTime                  uint64        `json:"updatedTime"`                  // Updated time
}

// OpenTpSlParam open take-profit/stop-loss parameters
type OpenTpSlParam struct {
	Price            string `json:"price"`            // Order price, market order fill 0
	Size             string `json:"size"`             // Order size
	ClientOrderId    string `json:"clientOrderId"`    // Client custom ID, for idempotency check, max length 64
	TriggerPriceType uint32 `json:"triggerPriceType"` // Trigger price type
	TriggerPrice     string `json:"triggerPrice"`     // Trigger price
	ExpireTime       uint64 `json:"expireTime"`       // Expiration time
}

// PerpetualCollateral perpetual contract collateral
type PerpetualCollateral struct {
	SubaccountId                          string `json:"subaccountId"`                          // Subaccount ID
	CoinId                                string `json:"coinId"`                                // Collateral coin ID
	Amount                                string `json:"amount"`                                // Collateral amount
	LegacyAmount                          string `json:"legacyAmount"`                          // Legacy accounting balance field, display only, not used in calculations
	CumDepositAmount                      string `json:"cumDepositAmount"`                      // Cumulative deposit amount
	CumWithdrawAmount                     string `json:"cumWithdrawAmount"`                     // Cumulative withdrawal amount
	CumTransferInAmount                   string `json:"cumTransferInAmount"`                   // Cumulative transfer-in amount
	CumTransferOutAmount                  string `json:"cumTransferOutAmount"`                  // Cumulative transfer-out amount
	CumCrossPositionOpenLongAmount        string `json:"cumCrossPositionOpenLongAmount"`        // Cumulative cross position open long deducted collateral amount
	CumCrossPositionOpenShortAmount       string `json:"cumCrossPositionOpenShortAmount"`       // Cumulative cross position open short added collateral amount
	CumCrossPositionCloseLongAmount       string `json:"cumCrossPositionCloseLongAmount"`       // Cumulative cross position close long added collateral amount
	CumCrossPositionCloseShortAmount      string `json:"cumCrossPositionCloseShortAmount"`      // Cumulative cross position close short deducted collateral amount
	CumIsolatedPositionOpenAmount         string `json:"cumIsolatedPositionOpenAmount"`         // Cumulative isolated position open deducted collateral amount
	CumIsolatedPositionCloseAmount        string `json:"cumIsolatedPositionCloseAmount"`        // Cumulative isolated position close added collateral amount
	CumIsolatedPositionMarginUpdateAmount string `json:"cumIsolatedPositionMarginUpdateAmount"` // Cumulative isolated position margin update collateral amount
}

// PositionStat position statistics
type PositionStat struct {
	CumOpenSize     string `json:"cumOpenSize"`     // Current open size (positive for long, negative for short)
	CumOpenValue    string `json:"cumOpenValue"`    // Current open value (accumulates on open, proportionally decreases on close)
	CumOpenFee      string `json:"cumOpenFee"`      // Current open fee after allocation (accumulates on open, proportionally decreases on close)
	CumCloseSize    string `json:"cumCloseSize"`    // Current close size (positive for long, negative for short)
	CumCloseValue   string `json:"cumCloseValue"`   // Current close value (accumulates on close, proportionally decreases on open)
	CumCloseFee     string `json:"cumCloseFee"`     // Current close fee after allocation (accumulates on close, proportionally decreases on open)
	CumFundingFee   string `json:"cumFundingFee"`   // Current position funding fee after allocation (accumulates on settlement, proportionally decreases on close)
	CumLiquidateFee string `json:"cumLiquidateFee"` // Current position liquidation fee after allocation (accumulates on settlement, proportionally decreases on close)
	CreatedTime     uint64 `json:"createdTime"`     // Created time
	UpdatedTime     uint64 `json:"updatedTime"`     // Updated time
}

// PerpetualPosition perpetual contract position
type PerpetualPosition struct {
	SubaccountId             string       `json:"subaccountId"`             // Subaccount ID
	CoinId                   string       `json:"coinId"`                   // Collateral coin ID
	ExchangeId               string       `json:"exchangeId"`               // Exchange ID, must be perpetual contract
	MarginMode               uint32       `json:"marginMode"`               // Margin mode
	OpenSize                 string       `json:"openSize"`                 // Current open size (positive for long, negative for short)
	OpenValue                string       `json:"openValue"`                // Current open value (accumulates on open, proportionally decreases on close)
	OpenFee                  string       `json:"openFee"`                  // Current open fee after allocation (accumulates on open, proportionally decreases on close)
	FundingFee               string       `json:"fundingFee"`               // Current position funding fee after allocation (accumulates on settlement, proportionally decreases on close)
	IsolatedMarginAmount     string       `json:"isolatedMarginAmount"`     // Isolated margin amount, meaningful when perpetual contract is in isolated mode
	IsolatedCollateralAmount string       `json:"isolatedCollateralAmount"` // Isolated collateral amount, meaningful when perpetual contract is in isolated mode
	CacheFundingIndex        string       `json:"cacheFundingIndex"`        // Cached funding rate index, updated when asset is updated
	LatestFundingIndex       string       `json:"latestFundingIndex"`       // Latest updated funding rate index
	TermCount                int32        `json:"termCount"`                // Long position term count, starts from 1, increments after complete close
	LongTermStat             PositionStat `json:"longTermStat"`             // Long position term cumulative statistics, cleared after complete close
	ShortTermStat            PositionStat `json:"shortTermStat"`            // Short position term cumulative statistics, cleared after complete close
	LongTotalStat            PositionStat `json:"longTotalStat"`            // Long position total cumulative statistics
	ShortTotalStat           PositionStat `json:"shortTotalStat"`           // Short position total cumulative statistics
	CreatedTime              uint64       `json:"createdTime"`              // Created time
	UpdatedTime              uint64       `json:"updatedTime"`              // Updated time
}

// PerpetualPositionTransaction perpetual contract position transaction
type PerpetualPositionTransaction struct {
	Id                             string `json:"id"`                             // Unique identifier
	SubaccountId                   string `json:"subaccountId"`                   // Subaccount ID
	CoinId                         string `json:"coinId"`                         // Coin ID
	ExchangeId                     string `json:"exchangeId"`                     // Contract ID
	TermCount                      uint32 `json:"termCount"`                      // Position term count
	MarginMode                     uint32 `json:"marginMode"`                     // Margin mode
	Type                           uint32 `json:"type"`                           // Transaction type
	DeltaOpenSize                  string `json:"deltaOpenSize"`                  // Position size change
	DeltaOpenValue                 string `json:"deltaOpenValue"`                 // Open value change
	DeltaOpenFee                   string `json:"deltaOpenFee"`                   // Open fee change
	DeltaFundingFee                string `json:"deltaFundingFee"`                // Funding fee change
	DeltaIsolatedMarginAmount      string `json:"deltaIsolatedMarginAmount"`      // Isolated margin amount change
	DeltaIsolatedCollateralAmount  string `json:"deltaIsolatedCollateralAmount"`  // Isolated collateral amount change
	BeforeOpenSize                 string `json:"beforeOpenSize"`                 // Position size before change
	BeforeOpenValue                string `json:"beforeOpenValue"`                // Open value before change
	BeforeOpenFee                  string `json:"beforeOpenFee"`                  // Open fee before change
	BeforeFundingFee               string `json:"beforeFundingFee"`               // Funding fee before change
	BeforeIsolatedMarginAmount     string `json:"beforeIsolatedMarginAmount"`     // Isolated margin amount before change
	BeforeIsolatedCollateralAmount string `json:"beforeIsolatedCollateralAmount"` // Isolated collateral amount before change
	FillSize                       string `json:"fillSize"`                       // Fill size (positive for buy, negative for sell)
	FillValue                      string `json:"fillValue"`                      // Fill value (positive for buy, negative for sell)
	FillFee                        string `json:"fillFee"`                        // Fill fee (usually zero or negative)
	FillPrice                      string `json:"fillPrice"`                      // Fill price (not precise, for display only)
	LiquidateFee                   string `json:"liquidateFee"`                   // Liquidation fee (exists when there is close fill, usually zero or negative)
	RealizePnl                     string `json:"realizePnl"`                     // Realized PnL (exists when there is close fill, not precise, for display only)
	IsPositionTp                   bool   `json:"isPositionTp"`                   // Whether it is a position take-profit/stop-loss order
	IsPositionSl                   bool   `json:"isPositionSl"`                   // Whether it is a position take-profit/stop-loss order
	IsLiquidate                    bool   `json:"isLiquidate"`                    // Whether it is a liquidation order
	IsDeleverage                   bool   `json:"isDeleverage"`                   // Whether it is an auto-deleverage order
	FundingTime                    uint64 `json:"fundingTime"`                    // Funding rate settlement time
	FundingRate                    string `json:"fundingRate"`                    // Funding rate
	FundingMarkPrice               string `json:"fundingMarkPrice"`               // Funding rate related index price
	FundingOraclePrice             string `json:"fundingOraclePrice"`             // Funding rate related oracle price
	FundingPositionSize            string `json:"fundingPositionSize"`            // Position size at funding fee settlement (positive for long, negative for short)
	OrderId                        string `json:"orderId"`                        // Associated order ID
	OrderFillTransactionId         string `json:"orderFillTransactionId"`         // Associated order fill transaction ID
	CollateralTransactionId        string `json:"collateralTransactionId"`        // Associated collateral transaction ID
	BlockHeight                    uint64 `json:"blockHeight"`                    // Block height
	BlockTime                      uint64 `json:"blockTime"`                      // Block time
	TransactionIndex               string `json:"transactionIndex"`               // Transaction index
	EventIndex                     string `json:"eventIndex"`                     // Event index
	CreatedTime                    uint64 `json:"createdTime"`                    // Created time
	UpdatedTime                    uint64 `json:"updatedTime"`                    // Updated time
}

// CollateralTransaction collateral transaction
type CollateralTransaction struct {
	Id                       string `json:"id"`                       // Unique identifier
	SubaccountId             string `json:"subaccountId"`             // Subaccount ID
	CoinId                   string `json:"coinId"`                   // Coin ID
	Type                     uint32 `json:"type"`                     // Transaction type
	DeltaAmount              string `json:"deltaAmount"`              // Collateral change amount
	DeltaLegacyAmount        string `json:"deltaLegacyAmount"`        // Legacy accounting balance field change amount
	BeforeAmount             string `json:"beforeAmount"`             // Collateral amount before change
	BeforeLegacyAmount       string `json:"beforeLegacyAmount"`       // Legacy accounting balance field before change
	TransferPeerSubaccountId string `json:"transferPeerSubaccountId"` // Transfer peer subaccount ID
	TransferPeerExchangeType uint32 `json:"transferPeerExchangeType"` // Transfer peer account exchange type
	TransferReason           uint32 `json:"transferReason"`           // Transfer reason
	TransferRemark           string `json:"transferRemark"`           // Transfer remark
	FillSize                 string `json:"fillSize"`                 // Fill size (positive for buy, negative for sell)
	FillValue                string `json:"fillValue"`                // Fill value (positive for buy, negative for sell)
	FillFee                  string `json:"fillFee"`                  // Fill fee (usually zero or negative)
	FillPrice                string `json:"fillPrice"`                // Fill price (not precise, for display only)
	LiqFee                   string `json:"liqFee"`                   // Liquidation fee (exists when there is close fill, usually zero or negative)
	RealizePnl               string `json:"realizePnl"`               // Realized PnL (exists when there is close fill, not precise, for display only)
	IsPositionTp             bool   `json:"isPositionTp"`             // Whether it is a position take-profit/stop-loss order
	IsPositionSl             bool   `json:"isPositionSl"`             // Whether it is a position take-profit/stop-loss order
	IsLiquidate              bool   `json:"isLiquidate"`              // Whether it is a liquidation order
	IsDeleverage             bool   `json:"isDeleverage"`             // Whether it is an auto-deleverage order
	FundingTime              uint64 `json:"fundingTime"`              // Funding rate settlement time
	FundingRate              string `json:"fundingRate"`              // Funding rate
	FundingIndexPrice        string `json:"fundingIndexPrice"`        // Funding rate related index price
	FundingOraclePrice       string `json:"fundingOraclePrice"`       // Funding rate related oracle price
	FundingPositionSize      string `json:"fundingPositionSize"`      // Position size at funding fee settlement (positive for long, negative for short)
	ExchangeId               string `json:"exchangeId"`               // Associated position contract ID
	OrderId                  string `json:"orderId"`                  // Associated order ID
	OrderFillTransactionId   string `json:"orderFillTransactionId"`   // Associated order fill transaction ID
	OrderSubaccountId        string `json:"orderSubaccountId"`        // Associated order subaccount ID
	PositionTransactionId    string `json:"positionTransactionId"`    // Associated position transaction ID
	BlockHeight              uint64 `json:"blockHeight"`              // Block height
	BlockTime                uint64 `json:"blockTime"`                // Block time
	TransactionIndex         string `json:"transactionIndex"`         // Transaction index
	EventIndex               string `json:"eventIndex"`               // Event index
	CreatedTime              uint64 `json:"createdTime"`              // Created time
	UpdatedTime              uint64 `json:"updatedTime"`              // Updated time
}

// AssetSnapshot asset snapshot
type AssetSnapshot struct {
	SubaccountId       string `json:"subaccountId"`       // Subaccount ID
	CoinId             string `json:"coinId"`             // Coin ID
	SnapshotTime       uint64 `json:"snapshotTime"`       // Snapshot time
	TotalEquity        string `json:"totalEquity"`        // Total collateral value
	TotalRealizePnl    string `json:"totalRealizePnl"`    // Total realized PnL
	TermRealizePnl     string `json:"termRealizePnl"`     // Term realized PnL
	TermFillValue      string `json:"termFillValue"`      // Term fill value (currently only returned when time_tag is 1)
	TermDepositAmount  string `json:"termDepositAmount"`  // Term deposit amount
	TermWithdrawAmount string `json:"termWithdrawAmount"` // Term withdrawal amount
}

// PerpetualPositionTerm perpetual contract position term
type PerpetualPositionTerm struct {
	SubaccountId    string `json:"subaccountId"`    // Subaccount ID
	CoinId          string `json:"coinId"`          // Collateral coin ID
	ExchangeId      string `json:"exchangeId"`      // Perpetual contract ID
	TermCount       uint32 `json:"termCount"`       // Term count, starts from 1, increments after complete close and open
	IsIsolated      bool   `json:"isIsolated"`      // Whether it is isolated
	CumOpenSize     string `json:"cumOpenSize"`     // Cumulative open size
	CumOpenValue    string `json:"cumOpenValue"`    // Cumulative open value
	CumOpenFee      string `json:"cumOpenFee"`      // Cumulative open fee
	CumCloseSize    string `json:"cumCloseSize"`    // Cumulative close size
	CumCloseValue   string `json:"cumCloseValue"`   // Cumulative close value
	CumCloseFee     string `json:"cumCloseFee"`     // Cumulative close fee
	CumFundingFee   string `json:"cumFundingFee"`   // Cumulative settled funding fee
	CumLiquidateFee string `json:"cumLiquidateFee"` // Cumulative liquidation fee
	CloseLeverage   string `json:"closeLeverage"`   // Leverage multiplier at complete close, actual type is decimal
	CreatedTime     uint64 `json:"createdTime"`     // Created time
	UpdatedTime     uint64 `json:"updatedTime"`     // Updated time
}

// OrderFillTransaction order fill transaction
type OrderFillTransaction struct {
	Id                                    string `json:"id"`                                    // Unique identifier
	SubaccountId                          string `json:"subaccountId"`                          // Subaccount ID
	CoinId                                string `json:"coinId"`                                // Trading coin ID
	ExchangeId                            string `json:"exchangeId"`                            // Exchange ID
	OrderId                               string `json:"orderId"`                               // Order ID
	IsBuy                                 bool   `json:"isBuy"`                                 // Buy/sell direction
	FillSize                              string `json:"fillSize"`                              // Actual fill size
	FillValue                             string `json:"fillValue"`                             // Actual fill value
	FillFee                               string `json:"fillFee"`                               // Actual fill fee
	FillPrice                             string `json:"fillPrice"`                             // Fill price (not precise, for display only)
	LiquidateFee                          string `json:"liquidateFee"`                          // If it's a liquidation (forced close) fill, this field is the liquidation fee
	RealizePnl                            string `json:"realizePnl"`                            // Actual realized PnL (only has value when fill includes close)
	IsMaker                               bool   `json:"isMaker"`                               // Actual fill direction, whether it is a maker fill
	IsPositionTp                          bool   `json:"isPositionTp"`                          // Whether it is a position take-profit/stop-loss order
	IsPositionSl                          bool   `json:"isPositionSl"`                          // Whether it is a position take-profit/stop-loss order
	IsLiquidate                           bool   `json:"isLiquidate"`                           // Whether it is a liquidation (forced close) order
	IsDeleverage                          bool   `json:"isDeleverage"`                          // Whether it is an auto-deleverage order
	SpotAssetTransactionId                string `json:"spotAssetTransactionId"`                // Associated spot asset transaction ID
	ClosePerpetualPositionTransactionId   string `json:"closePerpetualPositionTransactionId"`   // Associated close position transaction ID
	ClosePerpetualCollateralTransactionId string `json:"closePerpetualCollateralTransactionId"` // Associated close collateral transaction ID
	OpenPerpetualPositionTransactionId    string `json:"openPerpetualPositionTransactionId"`    // Associated open position transaction ID
	OpenPerpetualCollateralTransactionId  string `json:"openPerpetualCollateralTransactionId"`  // Associated open collateral transaction ID
	BlockHeight                           uint64 `json:"blockHeight"`                           // Block height
	BlockTime                             uint64 `json:"blockTime"`                             // Block time
	TransactionIndex                      string `json:"transactionIndex"`                      // Transaction index
	EventIndex                            string `json:"eventIndex"`                            // Event index
	CreatedTime                           uint64 `json:"createdTime"`                           // Created time
	UpdatedTime                           uint64 `json:"updatedTime"`                           // Updated time
}

// =============================== Request and Response Structures ===============================

// GetActiveOrderReq get active orders request
type GetActiveOrderReq struct {
	SubaccountId                    string `form:"subaccountId"`                             // Subaccount ID
	Size                            uint32 `form:"size"`                                     // Number of records, must be greater than 0 and less than or equal to 100
	OffsetData                      string `form:"offsetData,optional"`                      // Offset data
	PageOffsetDataCreatedTime       string `form:"pageOffsetDataCreatedTime,optional"`       // Pagination offset data, creation time
	PageOffsetDataItemId            string `form:"pageOffsetDataItemId,optional"`            // Pagination offset data, itemId
	FilterExchangeIdList            string `form:"filterExchangeIdList,optional"`            // Filter active orders for corresponding contracts, if empty get all contracts' active orders
	FilterOrderStatusList           string `form:"filterOrderStatusList,optional"`           // Filter orders with specified status, if empty get all status orders
	FilterIsLiquidateList           string `form:"filterIsLiquidateList,optional"`           // Filter orders with specified liquidation status, if empty get all orders
	FilterIsDeleverageList          string `form:"filterIsDeleverageList,optional"`          // Filter orders with specified deleverage status, if empty get all orders
	FilterIsPositionTpslList        string `form:"filterIsPositionTpslList,optional"`        // Filter orders with specified position take-profit/stop-loss status, if empty get all orders
	FilterOrderIdList               string `form:"filterOrderIdList,optional"`               // Filter orders with specified order IDs, if empty get all orders
	FilterStartCreatedTimeInclusive uint64 `form:"filterStartCreatedTimeInclusive,optional"` // Filter orders created at or after specified start time, if empty or 0 start from earliest
	FilterEndCreatedTimeExclusive   uint64 `form:"filterEndCreatedTimeExclusive,optional"`   // Filter orders created before specified end time, if empty or 0 get until latest
}

// GetActiveOrderResp get active orders response
type GetActiveOrderResp struct {
	BaseResp
	Data GetActiveOrderRespData `json:"data,omitempty"`
}

// GetActiveOrderRespData get active orders response data
type GetActiveOrderRespData struct {
	OrderList      []Order               `json:"orderList"`      // Order list
	PageOffsetData IndexerPageOffsetData `json:"pageOffsetData"` // Next page offset data
}

// GetHistoryOrderReq get history orders request
type GetHistoryOrderReq struct {
	SubaccountId                    string `form:"subaccountId"`                             // Subaccount ID
	Size                            uint32 `form:"size"`                                     // Number of records, must be greater than 0 and less than or equal to 100
	OffsetData                      string `form:"offsetData,optional"`                      // Offset data
	PageOffsetDataCreatedTime       string `form:"pageOffsetDataCreatedTime,optional"`       // Pagination offset data, creation time
	PageOffsetDataItemId            string `form:"pageOffsetDataItemId,optional"`            // Pagination offset data, itemId
	FilterExchangeIdList            string `form:"filterExchangeIdList,optional"`            // Filter history orders for corresponding contracts, if empty get all contracts' history orders
	FilterOrderStatusList           string `form:"filterOrderStatusList,optional"`           // Filter orders with specified status, if empty get all status orders
	FilterIsLiquidateList           string `form:"filterIsLiquidateList,optional"`           // Filter orders with specified liquidation status, if empty get all orders
	FilterIsDeleverageList          string `form:"filterIsDeleverageList,optional"`          // Filter orders with specified deleverage status, if empty get all orders
	FilterIsPositionTpslList        string `form:"filterIsPositionTpslList,optional"`        // Filter orders with specified position take-profit/stop-loss status, if empty get all orders
	FilterOrderIdList               string `form:"filterOrderIdList,optional"`               // Filter orders with specified order IDs, if empty get all orders
	FilterStartCreatedTimeInclusive uint64 `form:"filterStartCreatedTimeInclusive,optional"` // Filter orders created at or after specified start time, if empty or 0 start from earliest
	FilterEndCreatedTimeExclusive   uint64 `form:"filterEndCreatedTimeExclusive,optional"`   // Filter orders created before specified end time, if empty or 0 get until latest
}

// GetHistoryOrderResp get history orders response
type GetHistoryOrderResp struct {
	BaseResp
	Data GetHistoryOrderRespData `json:"data,omitempty"`
}

// GetHistoryOrderRespData get history orders response data
type GetHistoryOrderRespData struct {
	OrderList      []Order               `json:"orderList"`      // Order list
	PageOffsetData IndexerPageOffsetData `json:"pageOffsetData"` // Next page offset data
}

// GetPerpetualAccountAssetReq get perpetual contract account assets request
type GetPerpetualAccountAssetReq struct {
	SubaccountId string `form:"subaccountId"` // Subaccount ID
}

// GetPerpetualAccountAssetResp get perpetual contract account assets response
type GetPerpetualAccountAssetResp struct {
	BaseResp
	Data GetPerpetualAccountAssetRespData `json:"data,omitempty"`
}

// GetPerpetualAccountAssetRespData get perpetual contract account assets response data
type GetPerpetualAccountAssetRespData struct {
	SubaccountId                string                `json:"subaccountId"`                // Subaccount ID
	CollateralList              []PerpetualCollateral `json:"collateralList"`              // Collateral list
	PositionList                []PerpetualPosition   `json:"positionList"`                // Position list
	LastHandledBlockHeight      uint64                `json:"lastHandledBlockHeight"`      // Last handled block height
	LastHandledBlockTime        uint64                `json:"lastHandledBlockTime"`        // Last handled block time
	LastHandledTransactionIndex string                `json:"lastHandledTransactionIndex"` // Last handled transaction index
	LastHandledEventIndex       string                `json:"lastHandledEventIndex"`       // Last handled event index
}

// GetPositionTransactionReq get position transactions request
type GetPositionTransactionReq struct {
	SubaccountId                    string `form:"subaccountId"`                             // Subaccount ID
	Size                            uint32 `form:"size"`                                     // Number of records
	PageOffsetDataCreatedTime       string `form:"pageOffsetDataCreatedTime,optional"`       // Pagination offset data, creation time
	PageOffsetDataItemId            string `form:"pageOffsetDataItemId,optional"`            // Pagination offset data, itemId
	FilterExchangeIdList            string `form:"filterExchangeIdList,optional"`            // Exchange IDs, multiple exchange IDs separated by commas
	FilterTypeList                  string `form:"filterTypeList,optional"`                  // Transaction types, multiple transaction types separated by commas
	FilterMarginModeList            string `form:"filterMarginModeList,optional"`            // Margin modes, multiple margin modes separated by commas
	FilterStartCreatedTimeInclusive uint64 `form:"filterStartCreatedTimeInclusive,optional"` // Filter position transactions created at or after specified start time, if empty or 0 start from earliest
	FilterEndCreatedTimeExclusive   uint64 `form:"filterEndCreatedTimeExclusive,optional"`   // Filter position transactions created before specified end time, if empty or 0 get until latest
}

// GetPositionTransactionResp get position transactions response
type GetPositionTransactionResp struct {
	BaseResp
	Data GetPositionTransactionRespData `json:"data,omitempty"`
}

// GetPositionTransactionRespData get position transactions response data
type GetPositionTransactionRespData struct {
	PositionTransactionList []PerpetualPositionTransaction `json:"positionTransactionList"` // Position transaction list
	PageOffsetData          IndexerPageOffsetData          `json:"pageOffsetData"`          // Next page offset data
}

// GetCollateralTransactionReq get collateral transactions request
type GetCollateralTransactionReq struct {
	SubaccountId                    string `form:"subaccountId"`                             // Subaccount ID
	Size                            uint32 `form:"size"`                                     // Number of records
	PageOffsetDataCreatedTime       string `form:"pageOffsetDataCreatedTime,optional"`       // Pagination offset data, creation time
	PageOffsetDataItemId            string `form:"pageOffsetDataItemId,optional"`            // Pagination offset data, itemId
	FilterCoinId                    string `form:"filterCoinId,optional"`                    // Coin IDs, multiple coin IDs separated by commas
	FilterTypeList                  string `form:"filterTypeList,optional"`                  // Transaction types, multiple transaction types separated by commas
	FilterStartCreatedTimeInclusive uint64 `form:"filterStartCreatedTimeInclusive,optional"` // Filter collateral transactions created at or after specified start time, if empty or 0 start from earliest
	FilterEndCreatedTimeExclusive   uint64 `form:"filterEndCreatedTimeExclusive,optional"`   // Filter collateral transactions created before specified end time, if empty or 0 get until latest
}

// GetCollateralTransactionResp get collateral transactions response
type GetCollateralTransactionResp struct {
	BaseResp
	Data GetCollateralTransactionRespData `json:"data,omitempty"`
}

// GetCollateralTransactionRespData get collateral transactions response data
type GetCollateralTransactionRespData struct {
	CollateralTransactionList []CollateralTransaction `json:"collateralTransactionList"` // Collateral transaction list
	PageOffsetData            IndexerPageOffsetData   `json:"pageOffsetData"`            // Next page offset data
}

// GetAssetSnapshotReq get asset snapshots request
type GetAssetSnapshotReq struct {
	SubaccountId                    string `form:"subaccountId"`                             // Subaccount ID
	Size                            uint32 `form:"size"`                                     // Number of records
	PageOffsetDataCreatedTime       string `form:"pageOffsetDataCreatedTime,optional"`       // Pagination offset data, creation time
	PageOffsetDataItemId            string `form:"pageOffsetDataItemId,optional"`            // Pagination offset data, itemId
	FilterCoinId                    string `form:"filterCoinId,optional"`                    // Filter asset snapshots for corresponding coins, if empty get all coins' asset snapshots
	FilterTimeTag                   string `form:"filterTimeTag,optional"`                   // Filter asset snapshots by time type, 0 means query by hour, 1 means query by day
	FilterStartCreatedTimeInclusive uint64 `form:"filterStartCreatedTimeInclusive,optional"` // Filter asset snapshots created at or after specified start time, if empty or 0 start from earliest
	FilterEndCreatedTimeExclusive   uint64 `form:"filterEndCreatedTimeExclusive,optional"`   // Filter asset snapshots created before specified end time, if empty or 0 get until latest
}

// GetAssetSnapshotResp get asset snapshots response
type GetAssetSnapshotResp struct {
	BaseResp
	Data GetAssetSnapshotRespData `json:"data,omitempty"`
}

// GetAssetSnapshotRespData get asset snapshots response data
type GetAssetSnapshotRespData struct {
	AssetSnapshotList []AssetSnapshot       `json:"assetSnapshotList"` // Asset snapshot list
	PageOffsetData    IndexerPageOffsetData `json:"pageOffsetData"`    // Next page offset data
}

// GetHistoryOrderFillTransactionReq get history order fill transactions request
type GetHistoryOrderFillTransactionReq struct {
	SubaccountId                    string `form:"subaccountId"`                             // Subaccount ID
	Size                            uint32 `form:"size"`                                     // Number of records, must be greater than 0 and less than or equal to 100
	PageOffsetDataCreatedTime       string `form:"pageOffsetDataCreatedTime,optional"`       // Pagination offset data, creation time
	PageOffsetDataItemId            string `form:"pageOffsetDataItemId,optional"`            // Pagination offset data, itemId
	FilterExchangeIdList            string `form:"filterExchangeIdList,optional"`            // Exchange IDs, multiple exchange IDs separated by commas
	FilterCoinIdList                string `form:"filterCoinIdList,optional"`                // Coin IDs, multiple coin IDs separated by commas
	FilterOrderIdList               string `form:"filterOrderIdList,optional"`               // Order IDs, multiple order IDs separated by commas
	FilterStartCreatedTimeInclusive uint64 `form:"filterStartCreatedTimeInclusive,optional"` // Filter order fill transactions created at or after specified start time, if empty or 0 start from earliest
	FilterEndCreatedTimeExclusive   uint64 `form:"filterEndCreatedTimeExclusive,optional"`   // Filter order fill transactions created before specified end time, if empty or 0 get until latest
}

// GetHistoryOrderFillTransactionResp get history order fill transactions response
type GetHistoryOrderFillTransactionResp struct {
	BaseResp
	Data GetHistoryOrderFillTransactionRespData `json:"data,omitempty"`
}

// GetHistoryOrderFillTransactionRespData get history order fill transactions response data
type GetHistoryOrderFillTransactionRespData struct {
	OrderFillTransactionList []OrderFillTransaction `json:"orderFillTransactionList"` // Order fill transaction list
	PageOffsetData           IndexerPageOffsetData  `json:"pageOffsetData"`           // Next page offset data
}

// GetHistoryPositionTermReq get history position terms request
type GetHistoryPositionTermReq struct {
	SubaccountId                    string `form:"subaccountId"`                             // Subaccount ID
	Size                            uint32 `form:"size"`                                     // Number of records
	PageOffsetDataCreatedTime       string `form:"pageOffsetDataCreatedTime,optional"`       // Pagination offset data, creation time
	PageOffsetDataItemId            string `form:"pageOffsetDataItemId,optional"`            // Pagination offset data, itemId
	FilterExchangeIdList            string `form:"filterExchangeIdList,optional"`            // Exchange IDs, multiple exchange IDs separated by commas
	FilterStartCreatedTimeInclusive uint64 `form:"filterStartCreatedTimeInclusive,optional"` // Filter position terms created at or after specified start time, if empty or 0 start from earliest
	FilterEndCreatedTimeExclusive   uint64 `form:"filterEndCreatedTimeExclusive,optional"`   // Filter position terms created before specified end time, if empty or 0 get until latest
}

// GetHistoryPositionTermResp get history position terms response
type GetHistoryPositionTermResp struct {
	BaseResp
	Data GetHistoryPositionTermRespData `json:"data,omitempty"`
}

// GetHistoryPositionTermRespData get history position terms response data
type GetHistoryPositionTermRespData struct {
	PositionTermList []PerpetualPositionTerm `json:"positionTermList"` // Position term list
	PageOffsetData   IndexerPageOffsetData   `json:"pageOffsetData"`   // Next page offset data
}

// =============================== Helper Methods ===============================
