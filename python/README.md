# Antx SDK for Python

This package provides HTTP, WebSocket, and on-chain trading capabilities for Antx Protocol, mirroring the features of the Go SDK.

Features:
- HTTP APIs: coins, exchanges, klines, funding, orders, assets, transactions
- WebSocket (sync): subscribe ticker and klines with compatible payloads
- On-chain trading (planned in this package): bind agent, create/cancel orders

Quick start:
```bash
pip install -e .
```

Basic usage:
```python
from antx_sdk.client import AntxClient
from antx_sdk.constants import KLINE_TYPE_MINUTE_1, PRICE_TYPE_LAST

client = AntxClient(base_url="https://testnet.antxfi.com", ws_url="wss://testnet.antxfi.com/api/v1/ws")
coins = client.get_coin_list()
exchanges = client.get_exchange_list()

kline = client.get_kline({
    "exchangeId": "200001",
    "klineType": KLINE_TYPE_MINUTE_1,
    "priceType": PRICE_TYPE_LAST,
    "size": 10,
})
print(kline)
```

Protobufs:
- Python protobuf files are expected under `python/antx_proto/`. They can be generated from `proto/antx/**`.

License: Apache-2.0


