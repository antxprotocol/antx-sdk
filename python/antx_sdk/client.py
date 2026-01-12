import base64
import time
from typing import Any, Dict, Optional, Tuple

from . import constants as C
from .constants import ACCOUNT_HRP
from .http import HTTPClient
from .ws import WebSocketClient, parse_wrapped_first
from .crypto import (
    eth_personal_sign,
    convert_to_eth_addr,
    convert_to_antx_addr,
    secp256k1_pubkey_compressed,
    derive_antx_bech32_address,
)
from .tx import (
    Signer,
    build_auth_info,
    build_tx_body,
    encode_tx_base64,
    pack_any,
    sign_tx,
)

# Proto messages (require generated modules under python/antx_proto/**)
try:
    from antx_proto.antx.chain.agent import tx_pb2 as agent_tx
    from antx_proto.antx.chain.order import tx_pb2 as order_tx
except Exception:  # noqa: BLE001
    agent_tx = None
    order_tx = None


class AntxClient:
    def __init__(self, base_url: str = "", ws_url: str = "") -> None:
        self._http = HTTPClient(base_url) if base_url else HTTPClient("")
        self._base_url = base_url
        self._ws_url = ws_url
        self._ws: Optional[WebSocketClient] = None
        # signing context
        self._chain_id: str = ""
        self._agent_priv_hex: str = ""
        self._agent_priv_bytes: Optional[bytes] = None
        self._agent_pubkey_bytes: Optional[bytes] = None
        self._agent_address_bech32: Optional[str] = None
        self._account_number: Optional[int] = None

    def set_gateway(self, base_url: str, ws_url: str) -> None:
        self._base_url = base_url
        self._ws_url = ws_url
        self._http.set_base_url(base_url)

    # --------------- Credentials & account ---------------
    def set_credentials(self, agent_private_key_hex: str, chain_id: str, account_hrp: str = ACCOUNT_HRP) -> None:
        key_hex = agent_private_key_hex[2:] if agent_private_key_hex.startswith("0x") else agent_private_key_hex
        if len(key_hex) != 64:
            raise ValueError("invalid agent private key length")
        self._agent_priv_hex = key_hex
        self._agent_priv_bytes = bytes.fromhex(key_hex)
        self._agent_pubkey_bytes = secp256k1_pubkey_compressed(self._agent_priv_bytes)
        # auto-derive bech32 address from pubkey hash
        try:
            self._agent_address_bech32 = derive_antx_bech32_address(self._agent_priv_bytes, account_hrp)
        except Exception:
            self._agent_address_bech32 = None
        self._chain_id = chain_id
        # defer address resolution until first account query

    def _ensure_account_number(self) -> None:
        if self._account_number is not None and self._agent_address_bech32:
            return
        if not self._agent_priv_bytes:
            raise RuntimeError("credentials not set")
        # We need bech32 address; allow querying by hex address through converter
        # First compute raw eth-like address from pubkey (20 bytes) via convert_to_antx_addr helper if available
        # Here we must query gateway using bech32 agent address derived from private key
        # Since we don't have SDK Config to derive bech32 directly, we require user to later call set_agent_address if needed
        # Fallback: ask GetAddressInfo with any address to obtain normalized bech32
        if not self._agent_address_bech32:
            # Not ideal: we cannot derive bech32 without HRP and hashing rules; rely on caller to set it explicitly if needed
            pass
        # If address is still None, user should provide externally; otherwise continue

    # --------------- HTTP helpers ---------------
    def _http_get(self, path: str, params: Dict[str, Any]) -> Dict[str, Any]:
        try:
            return self._http.get(path, params)
        except Exception as e:
            # Preserve request_details if available
            if hasattr(e, 'request_details'):
                # Re-raise with preserved details
                raise
            raise

    def _http_post(self, path: str, data: Dict[str, Any]) -> Dict[str, Any]:
        return self._http.post(path, data)

    def _get_account_number_and_sequence(self, address: str) -> Tuple[int, int]:
        params = {"address": address}
        resp = self._http_get(C.GET_ADDRESS_INFO_PATH, params)
        if resp.get("code") != "0":
            raise RuntimeError(f"get account info failed: {resp.get('msg')}")
        data = resp.get("data", {})
        acc = int(data.get("accountNumber", "0"))
        seq = int(data.get("sequence", "0"))
        return acc, seq

    def get_account_info(self, address: Optional[str] = None) -> Dict[str, int]:
        """
        Get account number and sequence for the given address.
        
        Args:
            address: Bech32 address. If None, uses the agent address.
        
        Returns:
            Dict with 'accountNumber' and 'sequence' keys.
        
        Example:
            >>> info = client.get_account_info()
            >>> print(f"Account Number: {info['accountNumber']}, Sequence: {info['sequence']}")
        """
        if address is None:
            if not self._agent_address_bech32:
                raise RuntimeError("agent address not set; call set_credentials first or provide address")
            address = self._agent_address_bech32
        acc, seq = self._get_account_number_and_sequence(address)
        return {"accountNumber": acc, "sequence": seq}

    # --------------- Queries ---------------
    def get_coin_list(self) -> Any:
        resp = self._http_get(C.GET_COIN_LIST_PATH, {})
        if resp.get("code") != "0":
            raise RuntimeError(f"get coin list failed: {resp.get('msg')}")
        return resp["data"]["coinList"]

    def get_exchange_list(self) -> Any:
        resp = self._http_get(C.GET_EXCHANGE_LIST_PATH, {})
        if resp.get("code") != "0":
            raise RuntimeError(f"get exchange list failed: {resp.get('msg')}")
        return resp["data"]["exchangeList"]

    def get_subaccount_list(self, chain_type: int, chain_address: str, agent_address: str) -> Any:
        params = {
            "chainType": str(chain_type),
            "chainAddress": chain_address,
            "agentAddress": agent_address,
        }
        resp = self._http_get(C.GET_SUBACCOUNT_PATH, params)
        if resp.get("code") != "0":
            raise RuntimeError(f"get subaccount list failed: {resp.get('msg')}")
        return resp["data"]["subaccountList"]

    def get_eth_address(self, antx_bech32_address: str) -> str:
        return convert_to_eth_addr(antx_bech32_address)

    def get_kline(self, req: Dict[str, Any]) -> Any:
        params = {k: v for k, v in req.items() if v not in (None, "", 0)}
        resp = self._http_get(C.GET_KLINE_PATH, params)
        if resp.get("code") != "0":
            raise RuntimeError(f"get kline failed: {resp.get('msg')}")
        return resp

    def get_funding_history(self, req: Dict[str, Any]) -> Any:
        params = {k: v for k, v in req.items() if v not in (None, "", 0, False)}
        resp = self._http_get(C.GET_FUNDING_HISTORY_PATH, params)
        if resp.get("code") != "0":
            raise RuntimeError(f"get funding history failed: {resp.get('msg')}")
        return resp

    def get_active_order(self, req: Dict[str, Any]) -> Any:
        params = {k: v for k, v in req.items() if v not in (None, "")}
        resp = self._http_get(C.GET_ACTIVE_ORDER_PATH, params)
        if resp.get("code") != "0":
            raise RuntimeError(f"get active order failed: {resp.get('msg')}")
        return resp

    def get_history_order(self, req: Dict[str, Any]) -> Any:
        params = {k: v for k, v in req.items() if v not in (None, "")}
        resp = self._http_get(C.GET_HISTORY_ORDER_PATH, params)
        if resp.get("code") != "0":
            raise RuntimeError(f"get history order failed: {resp.get('msg')}")
        return resp

    def get_perpetual_account_asset(self, req: Dict[str, Any]) -> Any:
        params = {k: v for k, v in req.items() if v not in (None, "")}
        resp = self._http_get(C.GET_PERPETUAL_ACCOUNT_ASSET_PATH, params)
        if resp.get("code") != "0":
            raise RuntimeError(f"get perpetual account asset failed: {resp.get('msg')}")
        return resp

    def get_position_transaction(self, req: Dict[str, Any]) -> Any:
        params = {k: v for k, v in req.items() if v not in (None, "", 0)}
        resp = self._http_get(C.GET_POSITION_TRANSACTION_PATH, params)
        if resp.get("code") != "0":
            raise RuntimeError(f"get position transaction failed: {resp.get('msg')}")
        return resp

    def get_collateral_transaction(self, req: Dict[str, Any]) -> Any:
        params = {k: v for k, v in req.items() if v not in (None, "", 0)}
        resp = self._http_get(C.GET_COLLATERAL_TRANSACTION_PATH, params)
        if resp.get("code") != "0":
            raise RuntimeError(f"get collateral transaction failed: {resp.get('msg')}")
        return resp

    def get_asset_snapshot(self, req: Dict[str, Any]) -> Any:
        params = {k: v for k, v in req.items() if v not in (None, "", 0)}
        resp = self._http_get(C.GET_ASSET_SNAPSHOT_PATH, params)
        if resp.get("code") != "0":
            raise RuntimeError(f"get asset snapshot failed: {resp.get('msg')}")
        return resp

    def get_history_order_fill_transaction(self, req: Dict[str, Any]) -> Any:
        params = {k: v for k, v in req.items() if v not in (None, "", 0)}
        resp = self._http_get(C.GET_HISTORY_ORDER_FILL_TRANSACTION_PATH, params)
        if resp.get("code") != "0":
            raise RuntimeError(f"get history order fill transaction failed: {resp.get('msg')}")
        return resp

    def get_history_position_term(self, req: Dict[str, Any]) -> Any:
        params = {k: v for k, v in req.items() if v not in (None, "", 0)}
        resp = self._http_get(C.GET_HISTORY_POSITION_TERM_PATH, params)
        if resp.get("code") != "0":
            raise RuntimeError(f"get history position term failed: {resp.get('msg')}")
        return resp

    # --------------- WebSocket ---------------
    def connect_websocket(self, message_handler=None, error_handler=None) -> None:
        if self._ws:  # close previous
            try:
                self._ws.disconnect()
            except Exception:
                pass
        if not self._ws_url:
            raise ValueError("ws_url is not set")
        self._ws = WebSocketClient(self._ws_url, message_handler, error_handler)
        self._ws.connect()

    def subscribe_to_ticker(self, exchange_id: str):
        if not self._ws:
            raise RuntimeError("websocket not connected")
        return self._ws.subscribe_to_ticker(exchange_id)

    def subscribe_to_kline(self, price_type: str, exchange_id: str, kline_type: str):
        if not self._ws:
            raise RuntimeError("websocket not connected")
        return self._ws.subscribe_to_kline(price_type, exchange_id, kline_type)

    def disconnect_websocket(self) -> None:
        if self._ws:
            self._ws.disconnect()

    # --------------- Parsers ---------------
    def parse_ticker_data(self, data: bytes):
        return parse_wrapped_first(data, "ticker")

    def parse_kline_data(self, data: bytes):
        return parse_wrapped_first(data, "kline")

    def parse_depth_data(self, data: bytes):
        return parse_wrapped_first(data, "depth")

    def parse_trade_data(self, data: bytes):
        return parse_wrapped_first(data, "trade")

    def parse_funding_rate_data(self, data: bytes):
        return parse_wrapped_first(data, "funding rate")

    def parse_price_data(self, data: bytes):
        return parse_wrapped_first(data, "price")

    # --------------- Transactions ---------------
    def _sign_and_send_tx(self, type_url: str, msg, unordered: bool) -> str:
        if agent_tx is None or order_tx is None:
            raise RuntimeError("protobuf modules not available; generate python/antx_proto first")
        if not self._agent_priv_bytes or not self._agent_pubkey_bytes:
            raise RuntimeError("credentials not set; call set_credentials first")
        if not self._agent_address_bech32:
            raise RuntimeError("agent address not set; please set after credentials")

        if self._account_number is None:
            acc, _ = self._get_account_number_and_sequence(self._agent_address_bech32)
            self._account_number = acc
        
        if unordered:
            sequence = 0
            timeout_timestamp_ns = int((time.time() + 10) * 1_000_000_000)
        else:
            _, seq = self._get_account_number_and_sequence(self._agent_address_bech32)
            sequence = seq
            timeout_timestamp_ns = 0
        
        account_number = self._account_number

        any_msg = pack_any(msg, type_url)
        body = build_tx_body([any_msg], unordered=unordered, timeout_timestamp_ns=timeout_timestamp_ns)
        auth = build_auth_info(self._agent_pubkey_bytes, sequence, gas_limit=200000, fee_amounts=[])
        
        tx_bytes = sign_tx(body, auth, self._chain_id, account_number, Signer(
            private_key_bytes=self._agent_priv_bytes,
            public_key_bytes=self._agent_pubkey_bytes,
        ))
        raw_b64 = encode_tx_base64(tx_bytes)

        req = {
            "typeUrl": type_url,
            "rawTx": raw_b64,
            "accountNumber": int(account_number),
        }
        
        resp = self._http_post(C.SEND_TRANSACTION_PATH, req)
        if resp.get("code") != "0":
            raise RuntimeError(f"failed to send transaction: {resp.get('msg')}")
        data = resp.get("data", {})
        tx_hash = data.get("txHash") or data.get("hash") or data.get("txId") or ""
        return tx_hash

    def set_agent_address(self, bech32_address: str) -> None:
        self._agent_address_bech32 = bech32_address

    # Bind agent
    def bind_agent(self, eth_private_key_hex: str, chain_id: str, expire_seconds: int) -> str:
        if agent_tx is None:
            raise RuntimeError("protobuf modules not available; generate python/antx_proto first")
        if not self._agent_address_bech32:
            raise RuntimeError("agent address not set; call set_agent_address")
        create_ms = int(time.time() * 1000)
        expire_ms = int((time.time() + expire_seconds) * 1000)
        message = (
            f"Action:BindAgent\n"
            f"AgentAddress:{self._agent_address_bech32}\n"
            f"CreateTime:{create_ms}\n"
            f"ExpireTime:{expire_ms}\n"
            f"ChainId:{chain_id}"
        )
        signature = eth_personal_sign(message, eth_private_key_hex)
        from eth_account import Account as EthAccount
        eth_addr_raw = EthAccount.from_key(bytes.fromhex(eth_private_key_hex.replace("0x", ""))).address
        eth_addr = convert_to_eth_addr(eth_addr_raw)
        msg = agent_tx.MsgBindAgent(
            agent_address=self._agent_address_bech32,
            chain_type=1,  # EVM
            chain_address=eth_addr,
            create_time=create_ms,
            expire_time=expire_ms,
            chain_signature=signature,
        )
        return self._sign_and_send_tx(C.MSG_BIND_AGENT_TYPE_URL, msg, unordered=False)

    # Orders
    def create_order(self, params: Dict[str, Any]) -> str:
        if order_tx is None:
            raise RuntimeError("protobuf modules not available; generate python/antx_proto first")
        msg = order_tx.MsgCreateOrder(
            agent_address=self._agent_address_bech32,
            subaccount_id=params["subaccountId"],
            exchange_id=params["exchangeId"],
            margin_mode=params.get("marginMode", 0),
            leverage=params.get("leverage", 1),
            is_buy=params["isBuy"],
            price_scale=params.get("priceScale", 0),
            price_value=params.get("priceValue", 0),
            size_scale=params.get("sizeScale", 0),
            size_value=params.get("sizeValue", 0),
            client_order_id=params.get("clientOrderId", ""),
            time_in_force=params.get("timeInForce", 1),
            reduce_only=params.get("reduceOnly", False),
            expire_time=params.get("expireTime", 0),
            is_market=params.get("isMarket", False),
            is_position_tp=params.get("isPositionTp", False),
            is_position_sl=params.get("isPositionSl", False),
            trigger_type=params.get("triggerType", 0),
            trigger_price_type=params.get("triggerPriceType", 0),
            trigger_price_value=params.get("triggerPriceValue", 0),
            is_set_open_tp=params.get("isSetOpenTp", False),
            open_tp_param=params.get("openTpParam"),
            is_set_open_sl=params.get("isSetOpenSl", False),
            open_sl_param=params.get("openSlParam"),
        )
        return self._sign_and_send_tx(C.MSG_CREATE_ORDER_TYPE_URL, msg, unordered=True)

    def create_order_batch(self, params: Dict[str, Any]) -> str:
        if order_tx is None:
            raise RuntimeError("protobuf modules not available; generate python/antx_proto first")
        batch_list = []
        for o in params.get("createOrderParam", []):
            # Convert camelCase to snake_case for protobuf
            batch_order = order_tx.CreateOrderParam(
                is_buy=o.get("isBuy", False),
                price_scale=o.get("priceScale", 0),
                price_value=o.get("priceValue", 0),
                size_scale=o.get("sizeScale", 0),
                size_value=o.get("sizeValue", 0),
                client_order_id=o.get("clientOrderId", ""),
                time_in_force=o.get("timeInForce", 1),
                reduce_only=o.get("reduceOnly", False),
                expire_time=o.get("expireTime", 0),
                is_market=o.get("isMarket", False),
                is_position_tp=o.get("isPositionTp", False),
                is_position_sl=o.get("isPositionSl", False),
                trigger_type=o.get("triggerType", 0),
                trigger_price_type=o.get("triggerPriceType", 0),
                trigger_price_value=o.get("triggerPriceValue", 0),
                is_set_open_tp=o.get("isSetOpenTp", False),
                open_tp_param=o.get("openTpParam"),
                is_set_open_sl=o.get("isSetOpenSl", False),
                open_sl_param=o.get("openSlParam"),
            )
            batch_list.append(batch_order)
        msg = order_tx.MsgCreateOrderBatch(
            agent_address=self._agent_address_bech32,
            subaccount_id=params["subaccountId"],
            exchange_id=params["exchangeId"],
            margin_mode=params.get("marginMode", 0),
            leverage=params.get("leverage", 1),
            create_order_param=batch_list,
        )
        return self._sign_and_send_tx(C.MSG_CREATE_ORDER_BATCH_TYPE_URL, msg, unordered=True)

    def cancel_order(self, params: Dict[str, Any]) -> str:
        if order_tx is None:
            raise RuntimeError("protobuf modules not available; generate python/antx_proto first")
        msg = order_tx.MsgCancelOrder(
            agent_address=self._agent_address_bech32,
            subaccount_id=params["subaccountId"],
            order_id=params.get("orderIdList", []),
        )
        return self._sign_and_send_tx(C.MSG_CANCEL_ORDER_TYPE_URL, msg, unordered=True)

    def cancel_order_by_client_id(self, params: Dict[str, Any]) -> str:
        if order_tx is None:
            raise RuntimeError("protobuf modules not available; generate python/antx_proto first")
        msg = order_tx.MsgCancelOrderByClientId(
            agent_address=self._agent_address_bech32,
            subaccount_id=params["subaccountId"],
            client_order_id=params.get("clientOrderIdList", []),
        )
        return self._sign_and_send_tx(C.MSG_CANCEL_ORDER_BY_CLIENT_ID_TYPE_URL, msg, unordered=True)

    def cancel_all_order(self, params: Dict[str, Any]) -> str:
        if order_tx is None:
            raise RuntimeError("protobuf modules not available; generate python/antx_proto first")
        msg = order_tx.MsgCancelAllOrder(
            agent_address=self._agent_address_bech32,
            subaccount_id=params["subaccountId"],
            filter_exchange_id=params.get("filterExchangeIdList", []),
        )
        return self._sign_and_send_tx(C.MSG_CANCEL_ALL_ORDER_TYPE_URL, msg, unordered=True)

    def close_all_position(self, params: Dict[str, Any]) -> str:
        if order_tx is None:
            raise RuntimeError("protobuf modules not available; generate python/antx_proto first")
        msg = order_tx.MsgCloseAllPosition(
            agent_address=self._agent_address_bech32,
            subaccount_id=params["subaccountId"],
            filter_exchange_id=params.get("filterExchangeIdList", []),
        )
        return self._sign_and_send_tx(C.MSG_CLOSE_ALL_POSITION_TYPE_URL, msg, unordered=True)


