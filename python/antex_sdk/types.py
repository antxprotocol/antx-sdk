from dataclasses import dataclass, field
from typing import Any, Dict, List, Optional


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
    klineType: str
    priceType: str
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
    klineType: str
    priceType: str
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


# Trading queries (selected)
@dataclass
class GetActiveOrderReq:
    subaccountId: str
    size: int
    offsetData: str = ""
    pageOffsetDataCreatedTime: str = ""
    pageOffsetDataItemId: str = ""
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
    timeInForce: int
    reduceOnly: bool
    expireTime: int
    isPositionTp: bool
    isPositionSl: bool
    isLiquidate: bool
    isDeleverage: bool
    triggerType: int
    triggerPriceType: int
    triggerPrice: str
    openTpSlParentOrderId: str
    isSetOpenTp: bool
    openTpParam: Dict[str, Any]
    isSetOpenSl: bool
    openSlParam: Dict[str, Any]
    marginMode: int
    leverage: int
    takerFeeRatePpm: int
    makerFeeRatePpm: int
    liquidateFeeRatePpm: int
    addOrderBookBlockHeight: int
    addOrderBookBlockTime: int
    addOrderBookTransactionIndex: str
    addOrderBookOperationIndex: str
    status: int
    cancelReason: int
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
class GetHistoryOrderReq(GetActiveOrderReq):
    pass


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
    marginMode: int
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


