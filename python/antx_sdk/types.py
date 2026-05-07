from dataclasses import dataclass, field
from typing import Any, Dict, List, Optional, Union

from .enums import (
    AssetSnapshotTimeTag,
    KLineType,
    MarginMode,
    MarketPriceType,
    OrderCancelReason,
    OrderStatus,
    PriceType,
    TimeInForce,
    TriggerType,
)

# Enum-valued fields below are typed as `Union[<Enum>, int]` (or
# `Union[<StrEnum>, str]`) so callers may pass either an enum member or a
# raw int/str. See `antx_sdk.enums` for the full value tables and the
# authoritative .proto / .api sources.
#
# Filter list fields (filter*List) are HTTP query strings. Multiple values
# are comma-joined; build them with `enums.filter_int_list([...])` or
# `enums.filter_bool_list([...])`. Empty string = "no filter".

# =============================== Open TP/SL parameter ===============================
# Used both as a query response field on Order and as a request field
# on CreateOrderParam / CreateOrderBatchDetail.
@dataclass
class OpenTpSlParam:
    price: str = ""                                     # Order price, "0" or "" for market order
    size: str = ""                                      # Order size
    clientOrderId: str = ""                             # Client custom ID (for idempotency, max length 64)
    triggerPriceType: Union[PriceType, int] = PriceType.UNSPECIFIED
    triggerPrice: str = ""                              # Trigger price
    expireTime: int = 0                                 # Expiration time, unit: milliseconds


# Base response
@dataclass
class BaseResp:
    code: str
    msg: str


@dataclass
class IndexerPageOffsetData:
    createTime: str
    itemId: str


# Market data models
@dataclass
class KLine:
    klineId: str
    exchangeId: str
    klineType: Union[KLineType, str]
    priceType: Union[MarketPriceType, str]
    klineTime: int
    trades: str
    size: str
    value: str
    high: str
    low: str
    open: str
    close: str
    makerBuySize: str
    makerBuyValue: str


@dataclass
class TickerData:
    exchangeId: str
    lastPrice: str
    markPrice: str
    indexPrice: str
    oraclePrice: str
    priceChange: str
    priceChangePercent: str
    high: str
    low: str
    open: str
    close: str
    size: str
    value: str
    openInterest: str
    fundingRate: str
    fundingTime: str
    nextFundingTime: str
    startTime: str
    endTime: str
    highTime: str
    lowTime: str
    trades: str


@dataclass
class BookOrder:
    price: str
    size: str


@dataclass
class DepthData:
    exchangeId: str
    bids: List[BookOrder]
    asks: List[BookOrder]
    updatedTime: int


@dataclass
class Ticket:
    exchangeId: str
    price: str
    size: str
    value: str
    isBuy: bool
    time: str


@dataclass
class FundingRate:
    exchangeId: str
    fundingRate: str
    oraclePrice: str
    indexPrice: str
    fundingTime: int
    isSettlement: bool
    updatedTime: int


@dataclass
class Price:
    exchangeId: str
    price: str
    priceTime: int
    createdTime: int


# Coin / Exchange
@dataclass
class Coin:
    id: str
    symbol: str
    stepSizeScale: int
    assetChainId: str
    assetContractAddress: str


@dataclass
class RiskTier:
    maxLeverage: int
    maintenanceMarginRatioPpm: int
    positionValueUpperBound: str


@dataclass
class Perpetual:
    supportMarginModeList: List[int]
    riskTierList: List[RiskTier]
    liquidateFeeRatePpm: int
    defaultLeverage: int
    enableOrderCreate: bool
    enableOrderFill: bool
    enablePositionOpen: bool
    fundingInterestRatePpm: int
    fundingImpactMarginNotional: str
    fundingRateAbsMaxPpm: int
    fundingRateIntervalMinutes: int


@dataclass
class Exchange:
    id: str
    symbol: str
    baseCoinId: str
    quoteCoinId: str
    stepSizeScale: int
    tickSizeScale: int
    orderPriceMaxRatioPpm: int
    orderPriceMinRatioPpm: int
    orderSizeMax: str
    perpetual: Optional[Perpetual] = None


# Gateway responses
@dataclass
class GetCoinListRespData:
    coinList: List[Coin]


@dataclass
class GetCoinListResponse:
    code: str
    msg: str
    data: GetCoinListRespData


@dataclass
class GetExchangeListRespData:
    exchangeList: List[Exchange]


@dataclass
class GetExchangeListResponse:
    code: str
    msg: str
    data: GetExchangeListRespData


# Kline
@dataclass
class GetKLineReq:
    exchangeId: str
    klineType: Union[KLineType, str]                    # see enums.KLineType
    priceType: Union[MarketPriceType, str]              # see enums.MarketPriceType (string form)
    size: int = 100
    offsetData: str = ""
    filterBeginKlineTimeInclusive: int = 0
    filterEndKlineTimeExclusive: int = 0


@dataclass
class GetKLineRespData:
    klineList: List[KLine]
    nextPageOffsetData: str


@dataclass
class GetKLineResp:
    code: str
    msg: str
    data: GetKLineRespData


# Funding history
@dataclass
class GetFundingHistoryReq:
    exchangeId: str
    size: int = 100
    offsetData: str = ""
    filterSettlementFundingRate: bool = False
    filterBeginTimeInclusive: int = 0
    filterEndTimeExclusive: int = 0


@dataclass
class GetFundingHistoryRespData:
    fundingRateList: List[FundingRate]
    nextPageOffsetData: str


@dataclass
class GetFundingHistoryResp:
    code: str
    msg: str
    data: GetFundingHistoryRespData


# Trading queries
# Field set is the authoritative one from api-gateway/api/def_order.api.

@dataclass
class GetActiveOrderReq:
    subaccountId: str
    size: int
    pageOffsetDataCreatedTime: str = ""
    pageOffsetDataItemId: str = ""
    pageOffsetData: str = ""
    filterExchangeIdList: str = ""
    filterOrderStatusList: str = ""
    filterIsLiquidateList: str = ""
    filterIsDeleverageList: str = ""
    filterIsPositionTpslList: str = ""
    filterOrderIdList: str = ""
    filterStartCreatedTimeInclusive: int = 0
    filterEndCreatedTimeExclusive: int = 0


@dataclass
class Order:
    id: str
    subaccountId: str
    coinId: str
    exchangeId: str
    isBuy: bool
    price: str
    size: str
    clientOrderId: str
    timeInForce: Union[TimeInForce, int]
    reduceOnly: bool
    expireTime: int
    isPositionTp: bool
    isPositionSl: bool
    isLiquidate: bool
    isDeleverage: bool
    triggerType: Union[TriggerType, int]
    triggerPriceType: Union[PriceType, int]
    triggerPrice: str
    openTpSlParentOrderId: str
    isSetOpenTp: bool
    openTpParam: OpenTpSlParam
    isSetOpenSl: bool
    openSlParam: OpenTpSlParam
    marginMode: Union[MarginMode, int]
    leverage: int
    takerFeeRatePpm: int
    makerFeeRatePpm: int
    liquidateFeeRatePpm: int
    addOrderBookBlockHeight: int
    addOrderBookBlockTime: int
    addOrderBookTransactionIndex: str
    addOrderBookOperationIndex: str
    status: Union[OrderStatus, int]
    cancelReason: Union[OrderCancelReason, int]
    cumFillSize: str
    cumFillValue: str
    cumFillFee: str
    cumLiquidateFee: str
    maxFillPrice: str
    minFillPrice: str
    cumRealizePnl: str
    createdTime: int
    updatedTime: int


@dataclass
class GetActiveOrderRespData:
    orderList: List[Order]
    pageOffsetData: IndexerPageOffsetData


@dataclass
class GetActiveOrderResp:
    code: str
    msg: str
    data: GetActiveOrderRespData


@dataclass
class GetHistoryOrderReq:
    subaccountId: str
    size: int
    pageOffsetDataCreatedTime: str = ""
    pageOffsetDataItemId: str = ""
    pageOffsetData: str = ""
    filterExchangeIdList: str = ""
    filterIsBuyList: str = ""
    filterOrderStatusList: str = ""
    filterIsLiquidateList: str = ""
    filterIsDeleverageList: str = ""
    filterIsPositionTpList: str = ""
    filterIsPositionSlList: str = ""
    filterOrderIdList: str = ""
    filterStartCreatedTimeInclusive: int = 0
    filterEndCreatedTimeExclusive: int = 0


@dataclass
class GetHistoryOrderRespData:
    orderList: List[Order]
    pageOffsetData: IndexerPageOffsetData


@dataclass
class GetHistoryOrderResp:
    code: str
    msg: str
    data: GetHistoryOrderRespData


# Account asset / transactions (selected shells)
@dataclass
class GetPerpetualAccountAssetReq:
    subaccountId: str


@dataclass
class PerpetualCollateral:
    subaccountId: str
    coinId: str
    amount: str
    legacyAmount: str
    cumDepositAmount: str
    cumWithdrawAmount: str
    cumTransferInAmount: str
    cumTransferOutAmount: str
    cumCrossPositionOpenLongAmount: str
    cumCrossPositionOpenShortAmount: str
    cumCrossPositionCloseLongAmount: str
    cumCrossPositionCloseShortAmount: str
    cumIsolatedPositionOpenAmount: str
    cumIsolatedPositionCloseAmount: str
    cumIsolatedPositionMarginUpdateAmount: str


@dataclass
class PositionStat:
    cumOpenSize: str
    cumOpenValue: str
    cumOpenFee: str
    cumCloseSize: str
    cumCloseValue: str
    cumCloseFee: str
    cumFundingFee: str
    cumLiquidateFee: str
    createdTime: int
    updatedTime: int


@dataclass
class PerpetualPosition:
    subaccountId: str
    coinId: str
    exchangeId: str
    marginMode: Union[MarginMode, int]
    openSize: str
    openValue: str
    openFee: str
    fundingFee: str
    isolatedMarginAmount: str
    isolatedCollateralAmount: str
    cacheFundingIndex: str
    latestFundingIndex: str
    termCount: int
    longTermStat: PositionStat
    shortTermStat: PositionStat
    longTotalStat: PositionStat
    shortTotalStat: PositionStat
    createdTime: int
    updatedTime: int


@dataclass
class GetPerpetualAccountAssetRespData:
    subaccountId: str
    collateralList: List[PerpetualCollateral]
    positionList: List[PerpetualPosition]
    lastHandledBlockHeight: int
    lastHandledBlockTime: int
    lastHandledTransactionIndex: str
    lastHandledEventIndex: str


@dataclass
class GetPerpetualAccountAssetResp:
    code: str
    msg: str
    data: GetPerpetualAccountAssetRespData


# Generic envelope for send raw tx
@dataclass
class SendRawTxRequest:
    typeUrl: str
    rawTx: str
    accountNumber: int


@dataclass
class SendRawTxResponseData:
    txHash: str = ""
    rawTx: str = ""
    resultData: str = ""
    hash: str = ""
    txId: str = ""


@dataclass
class SendRawTxResponse:
    code: str
    msg: str
    data: SendRawTxResponseData


# =============================== Transaction detail (explorer) ===============================
# Returned by GET /explorer/tx/{hash}. `block == 0` means the tx has not been
# included on chain yet (still in mempool / async pending). `status == True`
# means included AND executed successfully; `status == False` means included
# but rejected at execution (errorCode tells you why).

@dataclass
class ChainTransactionDetail:
    rawTx: str = ""
    block: int = 0                       # 0 == pending / not yet included
    hash: str = ""
    fromAddress: str = ""                # JSON key is "from", aliased here
    status: bool = False
    error: Any = None
    errorCode: int = 0
    actionList: List[Any] = field(default_factory=list)
    resultData: str = ""


@dataclass
class GetTransactionDetailResp:
    code: str
    msg: str
    data: ChainTransactionDetail


# =============================== Additional query request types ===============================
# Mirrors of golang/types/trading.go GetXxxReq structs. All fields use camelCase
# to match the HTTP query parameter names the gateway expects.

@dataclass
class GetPositionTransactionReq:
    subaccountId: str
    size: int
    pageOffsetDataCreatedTime: str = ""
    pageOffsetDataItemId: str = ""
    pageOffsetData: str = ""
    filterExchangeIdList: str = ""           # comma-separated exchange IDs
    filterTypeList: str = ""                 # comma-separated transaction types
    filterMarginModeList: str = ""           # comma-separated margin modes
    filterStartCreatedTimeInclusive: int = 0
    filterEndCreatedTimeExclusive: int = 0


@dataclass
class GetCollateralTransactionReq:
    subaccountId: str
    size: int
    pageOffsetDataCreatedTime: str = ""
    pageOffsetDataItemId: str = ""
    pageOffsetData: str = ""
    filterExchangeIdList: str = ""           # comma-separated exchange IDs
    filterTypeList: str = ""                 # comma-separated transaction types
    filterStartCreatedTimeInclusive: int = 0
    filterEndCreatedTimeExclusive: int = 0


@dataclass
class GetAssetSnapshotReq:
    subaccountId: str
    size: int
    pageOffsetDataCreatedTime: str = ""
    pageOffsetDataItemId: str = ""
    pageOffsetData: str = ""
    filterCoinId: str = ""
    filterTimeTag: Union[AssetSnapshotTimeTag, str] = ""   # see enums.AssetSnapshotTimeTag
    filterStartCreatedTimeInclusive: int = 0
    filterEndCreatedTimeExclusive: int = 0


@dataclass
class GetHistoryOrderFillTransactionReq:
    subaccountId: str
    size: int
    pageOffsetDataCreatedTime: str = ""
    pageOffsetDataItemId: str = ""
    pageOffsetData: str = ""
    filterExchangeIdList: str = ""
    filterIsBuyList: str = ""
    filterIsLiquidateList: str = ""
    filterIsDeleverageList: str = ""
    filterIsPositionTpList: str = ""
    filterIsPositionSlList: str = ""
    filterOrderIdList: str = ""
    filterStartCreatedTimeInclusive: int = 0
    filterEndCreatedTimeExclusive: int = 0


@dataclass
class GetHistoryPositionTermReq:
    subaccountId: str
    size: int
    pageOffsetDataCreatedTime: str = ""
    pageOffsetDataItemId: str = ""
    pageOffsetData: str = ""
    filterExchangeIdList: str = ""
    filterStartCreatedTimeInclusive: int = 0
    filterEndCreatedTimeExclusive: int = 0


# =============================== Order / cancellation parameter types ===============================
# Mirrors of golang/types/gateway.go (CreateOrderParam, CreateOrderBatchDetail,
# CreateOrderBatchParam, CancelOrderParam, CancelOrderByClientIdParam,
# CancelAllOrderParam, CloseAllPositionParam).
#
# Notes:
# - subaccountId / exchangeId are Go uint64; in JSON / dict form they are
#   typically passed as strings. The Python client does not coerce types, so
#   pass whichever form the gateway accepts (str is the safer default).
# - Price/size are split into a (scale, value) pair. value is the raw integer
#   mantissa, scale is the decimal exponent (e.g. price = value * 10^scale).
# - agentAddress is filled in by AntxClient automatically — leave it empty
#   when constructing the dataclass.

@dataclass
class CreateOrderParam:
    subaccountId: str = ""
    exchangeId: str = ""
    marginMode: Union[MarginMode, int] = MarginMode.UNSPECIFIED
    leverage: int = 1
    isBuy: bool = False
    priceScale: int = 0
    priceValue: int = 0
    sizeScale: int = 0
    sizeValue: int = 0
    clientOrderId: str = ""
    timeInForce: Union[TimeInForce, int] = TimeInForce.GOOD_TIL_CANCEL
    reduceOnly: bool = False
    expireTime: int = 0                      # ms since epoch
    isMarket: bool = False
    isPositionTp: bool = False
    isPositionSl: bool = False
    triggerType: Union[TriggerType, int] = TriggerType.UNSPECIFIED
    triggerPriceType: Union[PriceType, int] = PriceType.UNSPECIFIED
    triggerPriceValue: int = 0
    openTpslParentOrderId: int = 0
    isSetOpenTp: bool = False
    openTpParam: Optional[OpenTpSlParam] = None
    isSetOpenSl: bool = False
    openSlParam: Optional[OpenTpSlParam] = None
    agentAddress: str = ""                   # auto-filled by client


@dataclass
class CreateOrderBatchDetail:
    isBuy: bool = False
    priceScale: int = 0
    priceValue: int = 0
    sizeScale: int = 0
    sizeValue: int = 0
    clientOrderId: str = ""
    timeInForce: Union[TimeInForce, int] = TimeInForce.GOOD_TIL_CANCEL
    reduceOnly: bool = False
    expireTime: int = 0
    isMarket: bool = False
    isPositionTp: bool = False
    isPositionSl: bool = False
    triggerType: Union[TriggerType, int] = TriggerType.UNSPECIFIED
    triggerPriceType: Union[PriceType, int] = PriceType.UNSPECIFIED
    triggerPriceValue: int = 0
    isSetOpenTp: bool = False
    openTpParam: Optional[OpenTpSlParam] = None
    isSetOpenSl: bool = False
    openSlParam: Optional[OpenTpSlParam] = None


@dataclass
class CreateOrderBatchParam:
    subaccountId: str = ""
    exchangeId: str = ""
    marginMode: Union[MarginMode, int] = MarginMode.UNSPECIFIED
    leverage: int = 1
    createOrderParam: List[CreateOrderBatchDetail] = field(default_factory=list)
    agentAddress: str = ""                   # auto-filled by client


@dataclass
class CancelOrderParam:
    subaccountId: str = ""
    orderIdList: List[str] = field(default_factory=list)
    agentAddress: str = ""                   # auto-filled by client


@dataclass
class CancelOrderByClientIdParam:
    subaccountId: str = ""
    clientOrderIdList: List[str] = field(default_factory=list)
    agentAddress: str = ""                   # auto-filled by client


@dataclass
class CancelAllOrderParam:
    subaccountId: str = ""
    filterExchangeIdList: List[str] = field(default_factory=list)
    agentAddress: str = ""                   # auto-filled by client


@dataclass
class CloseAllPositionParam:
    subaccountId: str = ""
    filterExchangeIdList: List[str] = field(default_factory=list)
    agentAddress: str = ""                   # auto-filled by client

