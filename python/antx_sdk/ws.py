import json
import logging
import queue
import ssl
import threading
from typing import Callable, Optional, Tuple
from urllib.parse import urlparse

from websocket import WebSocketApp

from .constants import WEBSOCKET_PATH


class WebSocketClient:
    def __init__(self, ws_url: str, message_handler: Optional[Callable[[bytes], None]] = None,
                 error_handler: Optional[Callable[[Exception], None]] = None) -> None:
        if ws_url.startswith("ws://") or ws_url.startswith("wss://"):
            self.url = ws_url
        else:
            self.url = f"ws://{ws_url}{WEBSOCKET_PATH}"
        self._message_handler = message_handler
        self._error_handler = error_handler
        self._app: Optional[WebSocketApp] = None
        self._thread: Optional[threading.Thread] = None
        self._connected = False

    def _get_origin_from_url(self) -> str:
        """Extract Origin from WebSocket URL"""
        try:
            parsed = urlparse(self.url)
            scheme = "https" if parsed.scheme == "wss" else "http"
            return f"{scheme}://{parsed.netloc}"
        except Exception:  # noqa: BLE001
            return ""

    def connect(self) -> None:
        def on_open(ws):  # noqa: ANN001
            self._connected = True
            logging.info("websocket connected: %s", self.url)

        def on_message(ws, message):  # noqa: ANN001
            if self._message_handler:
                try:
                    self._message_handler(message if isinstance(message, bytes) else message.encode())
                except Exception as e:  # noqa: BLE001
                    logging.exception("message handler error: %s", e)

        def on_error(ws, error):  # noqa: ANN001
            self._connected = False
            if self._error_handler:
                try:
                    self._error_handler(Exception(error))
                except Exception as e:  # noqa: BLE001
                    logging.exception("error handler error: %s", e)

        def on_close(ws, code, msg):  # noqa: ANN001
            self._connected = False
            logging.info("websocket closed: %s %s", code, msg)

        # Set request headers to avoid WAF blocking (same as Go version)
        header = {
            "X-App-Token": "ANTECH-APP-SECRET-KEY-001",
            "User-Agent": "Mozilla/5.0 (Mobile; FlutterApp/1.0)",
            "Origin": self._get_origin_from_url(),
        }

        self._app = WebSocketApp(
            self.url,
            on_open=on_open,
            on_message=on_message,
            on_error=on_error,
            on_close=on_close,
            header=header,
        )
        # Configure SSL options to skip certificate verification for testnet
        sslopt = {"cert_reqs": ssl.CERT_NONE} if self.url.startswith("wss://") else {}
        self._thread = threading.Thread(
            target=self._app.run_forever,
            kwargs={"sslopt": sslopt},
            daemon=True
        )
        self._thread.start()

    def is_connected(self) -> bool:
        return self._connected

    def disconnect(self) -> None:
        if self._app is not None:
            self._app.close()
        self._connected = False

    def _send_json(self, obj) -> None:
        if not self._connected or self._app is None:
            raise RuntimeError("websocket not connected")
        self._app.send(json.dumps(obj))

    def subscribe(self, channel: str) -> None:
        req = {
            "method": "subscribe",
            "subscription": {
                "channel": channel,
            },
        }
        self._send_json(req)

    def unsubscribe(self, channel: str) -> None:
        req = {
            "method": "unsubscribe",
            "subscription": {
                "channel": channel,
            },
        }
        self._send_json(req)

    def subscribe_to_ticker(self, exchange_id: str) -> queue.Queue:
        channel = f"ticker.{exchange_id}"
        q: queue.Queue = queue.Queue(maxsize=100)

        original = self._message_handler

        def handler(msg: bytes) -> None:
            try:
                payload = json.loads(msg)
                if payload.get("channel") == channel:
                    try:
                        q.put_nowait(msg)
                    except queue.Full:
                        pass
            finally:
                if original:
                    try:
                        original(msg)
                    except Exception:  # noqa: BLE001
                        pass

        self._message_handler = handler
        self.subscribe(channel)
        return q

    def subscribe_to_kline(self, price_type: str, exchange_id: str, kline_type: str) -> queue.Queue:
        channel = f"kline.{price_type}.{exchange_id}.{kline_type}"
        q: queue.Queue = queue.Queue(maxsize=100)

        original = self._message_handler

        def handler(msg: bytes) -> None:
            try:
                payload = json.loads(msg)
                if payload.get("channel") == channel:
                    try:
                        q.put_nowait(msg)
                    except queue.Full:
                        pass
            finally:
                if original:
                    try:
                        original(msg)
                    except Exception:  # noqa: BLE001
                        pass

        self._message_handler = handler
        self.subscribe(channel)
        return q


def parse_wrapped_first(data: bytes, key_type: str):
    try:
        obj = json.loads(data)
    except Exception as e:  # noqa: BLE001
        raise ValueError(f"failed to parse websocket response: {e}")
    arr = obj.get("data") or []
    if not arr:
        raise ValueError(f"no {key_type} data in response")
    return arr[0]


