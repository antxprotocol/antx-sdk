# Antex SDK for Python

This package provides HTTP, WebSocket, and on-chain trading capabilities for Antex Protocol, mirroring the features of the Go SDK.

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
from antex_sdk.client import AntexClient
from antex_sdk.constants import KLINE_TYPE_MINUTE_1, PRICE_TYPE_LAST

client = AntexClient(base_url="https://testnet.antex.ai", ws_url="wss://testnet.antex.ai/api/v1/ws")
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
- Python protobuf files are expected under `python/antex_proto/`. They can be generated from `golang/vendor/github.com/antexprotocol/antex-proto/proto/**`.

License: Apache-2.0


