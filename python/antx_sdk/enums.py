"""Enum definitions for Antx SDK request / response fields.

All integer enums inherit from :class:`enum.IntEnum`, so they are fully
interchangeable with plain ``int``:

    >>> int(MarginMode.CROSS) == 1
    True
    >>> MarginMode.CROSS == 1
    True
    >>> dataclasses.asdict(CreateOrderParam(marginMode=MarginMode.CROSS))['marginMode']
    1

Authoritative sources:
  - antx-proto/proto/antx/chain/order/order.proto
        TimeInForce, TriggerType, OrderCancelReason
  - antx-proto/proto/antx/chain/exchange/exchange.proto
        MarginMode, ExchangeType
  - antx-proto/proto/antx/chain/price/price.proto
        PriceType
  - antx-proto/proto/antx/server/indexer/indexer_type.proto
        OrderStatus, PerpetualCollateralTransactionType,
        PerpetualPositionTransactionType
  - antx-api-gateway/api/def_exchange.api
        KLineType, MarketPriceType (string-valued)
"""

from enum import Enum, IntEnum


# ============================== chain / order ==============================

class TimeInForce(IntEnum):
    UNSPECIFIED = 0
    GOOD_TIL_CANCEL = 1          # GTC, default for limit orders
    FILL_OR_KILL = 2             # FOK
    IMMEDIATE_OR_CANCEL = 3      # IOC, recommended for market orders
    POST_ONLY = 4                # Maker-only


class TriggerType(IntEnum):
    UNSPECIFIED = 0              # Not a conditional order
    STOP_LOSS = 1
    TAKE_PROFIT = 2


class OrderCancelReason(IntEnum):
    UNSPECIFIED = 0
    FILLED = 1
    USER = 2
    EXPIRE = 3
    TIME_IN_FORCE = 4
    REDUCE_ONLY = 5
    AVAILABLE_AMOUNT_NOT_ENOUGH = 6
    MAX_POSITION_VALUE_EXCEED = 7
    SELF_TRADE = 8
    LIQUIDATING = 9
    FILL_FAILED = 10
    PRICE_LIMIT = 11
    ADMIN = 12
    SYSTEM = 13


# ============================ chain / exchange ============================

class MarginMode(IntEnum):
    UNSPECIFIED = 0
    CROSS = 1                    # Cross margin
    ISOLATED = 2                 # Isolated margin


class ExchangeType(IntEnum):
    UNSPECIFIED = 0
    SPOT = 1
    PERPETUAL = 2


# ============================== chain / price ==============================

class PriceType(IntEnum):
    """On-chain price-type enum.

    Used by conditional-order ``triggerPriceType`` and OpenTpSlParam.
    For market-data REST ``GetKLineReq.priceType`` use the *string* form
    instead — see :class:`MarketPriceType`.
    """
    UNSPECIFIED = 0
    LAST = 1
    ASK_BEST = 2
    BID_BEST = 3
    ORACLE = 4
    INDEX = 5


# ============================ indexer / order =============================

class OrderStatus(IntEnum):
    """Indexer-reported order status (the ``status`` field on Order).

    Note: ``antx_sdk.constants.ORDER_STATUS_*`` uses a different legacy
    numbering preserved for backward compatibility — prefer this enum.
    """
    UNSPECIFIED = 0
    OPEN = 1                     # Resting on book, possibly partially filled
    FILLED = 2                   # Fully filled (terminal)
    CANCELED = 3                 # Canceled, possibly after partial fill (terminal)
    UNTRIGGERED = 4              # Conditional order, not yet triggered


class PerpetualCollateralTransactionType(IntEnum):
    UNSPECIFIED = 0
    TRANSFER_IN = 1
    TRANSFER_OUT = 2
    CROSS_POSITION_OPEN_LONG = 3        # cross buy-open-long: deducts collateral
    CROSS_POSITION_OPEN_SHORT = 4       # cross sell-open-short: adds collateral
    CROSS_POSITION_CLOSE_LONG = 5       # cross sell-close-long: adds collateral
    CROSS_POSITION_CLOSE_SHORT = 6      # cross buy-close-short: deducts collateral
    ISOLATED_POSITION_OPEN = 7
    ISOLATED_POSITION_CLOSE = 8
    ISOLATED_POSITION_MARGIN_UPDATE = 9
    POSITION_FUNDING = 10               # funding-rate settlement
    FILL_FEE_INCOME = 11                # fee-account only
    LIQUIDATE_FEE_INCOME = 12           # fee-account only
    TRANSFER_FEE_INCOME = 13


class PerpetualPositionTransactionType(IntEnum):
    UNSPECIFIED = 0
    OPEN_LONG = 1
    OPEN_SHORT = 2
    CLOSE_LONG = 3
    CLOSE_SHORT = 4
    FUNDING = 5                  # funding settlement
    ISOLATED_MARGIN_UPDATE = 6


# =========================== api-gateway (string) ===========================

class KLineType(str, Enum):
    """String values accepted by ``GetKLineReq.klineType``."""
    MINUTE_1 = "MINUTE_1"
    MINUTE_5 = "MINUTE_5"
    MINUTE_15 = "MINUTE_15"
    MINUTE_30 = "MINUTE_30"
    HOUR_1 = "HOUR_1"
    HOUR_2 = "HOUR_2"
    HOUR_4 = "HOUR_4"
    HOUR_6 = "HOUR_6"
    HOUR_8 = "HOUR_8"
    HOUR_12 = "HOUR_12"
    DAY_1 = "DAY_1"
    WEEK_1 = "WEEK_1"
    MONTH_1 = "MONTH_1"


class MarketPriceType(str, Enum):
    """String values accepted by ``GetKLineReq.priceType``.

    Distinct from the on-chain :class:`PriceType` IntEnum — the kline /
    market-data REST endpoints use these uppercase strings instead.
    """
    LAST = "PRICE_TYPE_LAST"
    ASK_BEST = "PRICE_TYPE_ASK_BEST"
    BID_BEST = "PRICE_TYPE_BID_BEST"
    MARK = "PRICE_TYPE_MARK"
    ORACLE = "PRICE_TYPE_ORACLE"


class AssetSnapshotTimeTag(str, Enum):
    """Values accepted by ``GetAssetSnapshotReq.filterTimeTag``."""
    HOURLY = "0"
    DAILY = "1"


# =============================== Helpers ===============================

def filter_int_list(values) -> str:
    """Encode a sequence of IntEnum / int values as a comma-separated
    filter list. Empty input -> empty string (== "no filter").

        filter_int_list([MarginMode.CROSS, MarginMode.ISOLATED]) -> "1,2"
    """
    return ",".join(str(int(v)) for v in values)


def filter_bool_list(values) -> str:
    """Encode a sequence of booleans as a comma-separated filter list.

        filter_bool_list([True])         -> "true"
        filter_bool_list([True, False])  -> "true,false"
    """
    return ",".join("true" if bool(v) else "false" for v in values)
