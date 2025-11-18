#!/usr/bin/env python3
"""
Complete example for Antex SDK Python

This example demonstrates:
1. Basic functions (coin list, exchange list)
2. Market data (kline, funding history)
3. WebSocket real-time data
4. Trading functions (bind agent, create limit/market orders, cancel orders, batch orders)
5. Trading queries (active orders, history orders, account assets, position/collateral transactions, snapshots, etc.)

This example matches the Go SDK's complete_example.go functionality.
"""

import os
import time
from typing import Optional

from antex_sdk.client import AntexClient
from antex_sdk.constants import (
    ACCOUNT_HRP,
    KLINE_TYPE_MINUTE_1,
    PRICE_TYPE_LAST,
)


# Configuration
GATEWAY_URL = os.environ.get("ANTEX_GATEWAY", "https://testnet.antex.ai")
WS_URL = os.environ.get("ANTEX_WS", "wss://testnet.antex.ai/api/v1/ws")
CHAIN_ID = os.environ.get("ANTEX_CHAIN_ID", "antex-testnet")

# Credentials (set via environment variables or replace with your keys)
ETH_PRIVATE_KEY = os.environ.get("ETH_PRIVATE_KEY", "")
AGENT_PRIVATE_KEY = os.environ.get("AGENT_PRIVATE_KEY", "")

# Example parameters
DEFAULT_EXCHANGE_ID = "200001"


def demo_basic_functions(client: AntexClient):
    """Demo basic functions: coin list, exchange list"""
    print("\n=== 1. Basic Functions Demo ===")

    print("\n1.1 Getting coin list:")
    try:
        coins = client.get_coin_list()
        print(f"✓ Retrieved {len(coins)} coins")
        for i, coin in enumerate(coins[:3]):  # Show first 3
            print(f"  Coin {i+1}: ID={coin.get('id', 'N/A')}, Symbol={coin.get('symbol', 'N/A')}, "
                  f"StepSizeScale={coin.get('stepSizeScale', 'N/A')}")
    except Exception as e:
        print(f"⚠ Failed to get coin list: {e}")

    print("\n1.2 Getting exchange list:")
    try:
        exchanges = client.get_exchange_list()
        print(f"✓ Retrieved {len(exchanges)} exchanges")
        for i, exchange in enumerate(exchanges[:3]):  # Show first 3
            print(f"  Exchange {i+1}: ID={exchange.get('id', 'N/A')}, Symbol={exchange.get('symbol', 'N/A')}, "
                  f"BaseCoinId={exchange.get('baseCoinId', 'N/A')}, QuoteCoinId={exchange.get('quoteCoinId', 'N/A')}")
    except Exception as e:
        print(f"⚠ Failed to get exchange list: {e}")


def demo_market_data(client: AntexClient):
    """Demo market data: kline, funding history"""
    print("\n=== 2. Market Data Demo ===")

    print("\n2.1 Getting kline data:")
    try:
        kline_req = {
            "exchangeId": DEFAULT_EXCHANGE_ID,
            "klineType": KLINE_TYPE_MINUTE_1,
            "priceType": PRICE_TYPE_LAST,
            "size": 10,
        }
        kline_resp = client.get_kline(kline_req)
        kline_list = kline_resp.get("data", {}).get("klineList", [])
        print(f"✓ Retrieved {len(kline_list)} kline records")
        for i, kline in enumerate(kline_list[:3]):  # Show first 3
            kline_time = kline.get("klineTime", 0)
            time_str = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(kline_time / 1000)) if kline_time else "N/A"
            print(f"  Kline {i+1}: Time={time_str}, Open={kline.get('open', 'N/A')}, "
                  f"High={kline.get('high', 'N/A')}, Low={kline.get('low', 'N/A')}, "
                  f"Close={kline.get('close', 'N/A')}, Size={kline.get('size', 'N/A')}")
    except Exception as e:
        print(f"⚠ Failed to get kline data: {e}")

    print("\n2.2 Getting funding history:")
    try:
        funding_req = {
            "exchangeId": DEFAULT_EXCHANGE_ID,
            "size": 5,
        }
        funding_resp = client.get_funding_history(funding_req)
        funding_list = funding_resp.get("data", {}).get("fundingRateList", [])
        print(f"✓ Retrieved {len(funding_list)} funding rate records")
        for i, rate in enumerate(funding_list[:3]):  # Show first 3
            funding_time = rate.get("fundingTime", 0)
            time_str = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(funding_time / 1000)) if funding_time else "N/A"
            print(f"  Funding Rate {i+1}: Time={time_str}, Rate={rate.get('fundingRate', 'N/A')}, "
                  f"OraclePrice={rate.get('oraclePrice', 'N/A')}")
    except Exception as e:
        print(f"⚠ Failed to get funding history: {e}")


def demo_websocket_realtime(client: AntexClient):
    """Demo WebSocket real-time market data"""
    print("\n=== 3. WebSocket Real-time Data Demo ===")

    def on_msg(data: bytes):
        print(f"  Received WebSocket message: {data[:80]}...")

    def on_err(e: Exception):
        print(f"  WebSocket error: {e}")

    try:
        print("\n3.1 Connecting WebSocket...")
        client.connect_websocket(on_msg, on_err)
        time.sleep(1)
        if client._ws and client._ws.is_connected():
            print("✓ WebSocket connected")

            print("\n3.2 Subscribing to ticker...")
            ticker_chan = client.subscribe_to_ticker(DEFAULT_EXCHANGE_ID)
            print("✓ Subscribed to ticker data, waiting for data...")

            # Wait for data with timeout
            timeout = time.time() + 5
            received = False
            while time.time() < timeout:
                try:
                    data = ticker_chan.get(timeout=2)
                    print("✓ Received ticker data:")
                    ticker_data = client.parse_ticker_data(data)
                    print(f"  ExchangeId: {ticker_data.get('exchangeId', 'N/A')}, "
                          f"LastPrice: {ticker_data.get('lastPrice', 'N/A')}, "
                          f"PriceChangePercent: {ticker_data.get('priceChangePercent', 'N/A')}%")
                    received = True
                    break
                except Exception:
                    continue

            if not received:
                print("⚠ No ticker data received within 5 seconds")

            client.disconnect_websocket()
            print("✓ WebSocket disconnected")
        else:
            print("⚠ WebSocket connection failed or timeout")
    except Exception as e:
        print(f"⚠ WebSocket demo failed: {e}")


def demo_trading_functions(client: AntexClient):
    """Demo trading functions: bind agent, create/cancel orders"""
    print("\n=== 4. Trading Functions Demo ===")
    print("Note: Trading functions require valid credentials and agent binding")

    if not ETH_PRIVATE_KEY or not AGENT_PRIVATE_KEY:
        print("⚠ ETH_PRIVATE_KEY and AGENT_PRIVATE_KEY not set, skipping trading demo")
        print("   Set them via environment variables or modify the script")
        return

    try:
        print("\n4.0 Setting credentials...")
        client.set_credentials(agent_private_key_hex=AGENT_PRIVATE_KEY, chain_id=CHAIN_ID, account_hrp=ACCOUNT_HRP)
        agent_addr = client._agent_address_bech32
        print(f"✓ Agent address: {agent_addr}")

        print("\n4.1 Binding agent...")
        try:
            tx_hash = client.bind_agent(eth_private_key_hex=ETH_PRIVATE_KEY, chain_id=CHAIN_ID, expire_seconds=3600)
            print(f"✓ Bind agent successful, tx_hash: {tx_hash}")
            time.sleep(3)
        except Exception as e:
            print(f"⚠ Bind agent failed: {e}")
            print("  (This may be expected if agent is already bound)")

        # Get subaccount list (example)
        print("\n4.2 Getting subaccount list...")
        try:
            from eth_account import Account
            from antex_sdk.crypto import convert_to_eth_addr
            eth_addr = Account.from_key(bytes.fromhex(ETH_PRIVATE_KEY.replace("0x", ""))).address
            eth_addr_checksum = convert_to_eth_addr(eth_addr)
            sub_list = client.get_subaccount_list(chain_type=1, chain_address=eth_addr_checksum, agent_address=agent_addr)
            if sub_list:
                test_subaccount_id = sub_list[0].get("id", "")
                print(f"✓ Found subaccount: {test_subaccount_id}")
                
                # Convert subaccount ID to int
                test_subaccount_id_int = int(test_subaccount_id)
                exchange_id_int = int(DEFAULT_EXCHANGE_ID)
                
                print("\n4.3 Creating limit buy order...")
                try:
                    create_order_req = {
                        "subaccountId": test_subaccount_id_int,
                        "exchangeId": exchange_id_int,
                        "marginMode": 1,  # Full margin
                        "leverage": 1,
                        "isBuy": True,
                        "priceScale": 2,
                        "priceValue": 100000,  # Price 1000.00
                        "sizeScale": 3,
                        "sizeValue": 100,  # Size 0.100
                        "clientOrderId": "py-test-001",
                        "timeInForce": 1,  # GTC
                        "reduceOnly": False,
                        "expireTime": int(time.time()) + 86400,  # 24 hours
                        "isMarket": False,
                        "isPositionTp": False,
                        "isPositionSl": False,
                        "triggerType": 0,
                        "triggerPriceType": 0,
                        "triggerPriceValue": 0,
                        "isSetOpenTp": False,
                        "isSetOpenSl": False,
                    }
                    order_tx_hash = client.create_order(create_order_req)
                    print(f"✓ Create order successful, tx_hash: {order_tx_hash}")
                    # Wait for transaction confirmation
                    if order_tx_hash:
                        print("Waiting for transaction confirmation...")
                        time.sleep(3)
                except Exception as e:
                    print(f"⚠ Create order failed: {e}")

                print("\n4.4 Creating market sell order...")
                try:
                    market_order_req = {
                        "subaccountId": test_subaccount_id_int,
                        "exchangeId": exchange_id_int,
                        "marginMode": 1,  # Full margin
                        "leverage": 1,
                        "isBuy": False,
                        "priceScale": 2,
                        "priceValue": 0,  # Market order price is 0
                        "sizeScale": 3,
                        "sizeValue": 50,  # Size 0.050
                        "clientOrderId": "py-market-001",
                        "timeInForce": 3,  # IOC (more suitable for market orders)
                        "reduceOnly": False,
                        "expireTime": int(time.time()) + 86400,  # 24 hours
                        "isMarket": True,  # Market order
                        "isPositionTp": False,
                        "isPositionSl": False,
                        "triggerType": 0,
                        "triggerPriceType": 0,
                        "triggerPriceValue": 0,
                        "isSetOpenTp": False,
                        "isSetOpenSl": False,
                    }
                    market_order_tx_hash = client.create_order(market_order_req)
                    print(f"✓ Create market order successful, tx_hash: {market_order_tx_hash}")
                except Exception as e:
                    print(f"⚠ Create market order failed: {e}")

                print("\n4.5 Canceling order...")
                try:
                    # Use a sample order ID (replace with actual order ID from previous steps)
                    # For demo purposes, using a placeholder - in real usage, get from active orders
                    sample_order_id = "188531408901"  # Example order ID
                    cancel_order_req = {
                        "subaccountId": test_subaccount_id_int,
                        "orderIdList": [int(sample_order_id)],
                    }
                    cancel_tx_hash = client.cancel_order(cancel_order_req)
                    print(f"✓ Cancel order successful, tx_hash: {cancel_tx_hash}")
                except Exception as e:
                    print(f"⚠ Cancel order failed: {e}")

                print("\n4.6 Creating batch orders...")
                try:
                    batch_order_req = {
                        "agentAddress": agent_addr,
                        "subaccountId": test_subaccount_id_int,
                        "exchangeId": exchange_id_int,
                        "marginMode": 1,
                        "leverage": 1,
                        "createOrderParam": [
                            {
                                "isBuy": True,
                                "priceScale": 2,
                                "priceValue": 95000,  # Price 950.00
                                "sizeScale": 3,
                                "sizeValue": 200,  # Size 0.200
                                "clientOrderId": "batch-order-001",
                                "timeInForce": 1,
                                "reduceOnly": False,
                                "expireTime": int(time.time()) + 86400,  # 24 hours
                                "isMarket": False,
                                "isPositionTp": False,
                                "isPositionSl": False,
                                "triggerType": 0,
                                "triggerPriceType": 0,
                                "triggerPriceValue": 0,
                                "isSetOpenTp": False,
                                "isSetOpenSl": False,
                            },
                            {
                                "isBuy": False,
                                "priceScale": 2,
                                "priceValue": 105000,  # Price 1050.00
                                "sizeScale": 3,
                                "sizeValue": 150,  # Size 0.150
                                "clientOrderId": "batch-order-002",
                                "timeInForce": 1,
                                "reduceOnly": False,
                                "expireTime": int(time.time()) + 86400,  # 24 hours
                                "isMarket": False,
                                "isPositionTp": False,
                                "isPositionSl": False,
                                "triggerType": 0,
                                "triggerPriceType": 0,
                                "triggerPriceValue": 0,
                                "isSetOpenTp": False,
                                "isSetOpenSl": False,
                            },
                        ],
                    }
                    batch_tx_hash = client.create_order_batch(batch_order_req)
                    print(f"✓ Create batch orders successful, tx_hash: {batch_tx_hash}")
                except Exception as e:
                    print(f"⚠ Create batch orders failed: {e}")
            else:
                print("⚠ No subaccounts found")
        except Exception as e:
            print(f"⚠ Get subaccount list failed: {e}")

    except Exception as e:
        print(f"⚠ Trading functions demo failed: {e}")
        import traceback
        traceback.print_exc()


def demo_trading_queries(client: AntexClient):
    """Demo trading queries: active orders, account assets, etc."""
    print("\n=== 5. Trading Queries Demo ===")

    if not AGENT_PRIVATE_KEY:
        print("⚠ AGENT_PRIVATE_KEY not set, skipping trading queries demo")
        return

    try:
        client.set_credentials(agent_private_key_hex=AGENT_PRIVATE_KEY, chain_id=CHAIN_ID, account_hrp=ACCOUNT_HRP)
        agent_addr = client._agent_address_bech32

        # Get subaccount
        try:
            from eth_account import Account
            from antex_sdk.crypto import convert_to_eth_addr
            # Use ETH_PRIVATE_KEY to get ETH address (same as demo_trading_functions)
            eth_addr = Account.from_key(bytes.fromhex(ETH_PRIVATE_KEY.replace("0x", ""))).address
            eth_addr_checksum = convert_to_eth_addr(eth_addr)
            sub_list = client.get_subaccount_list(chain_type=1, chain_address=eth_addr_checksum, agent_address=agent_addr)
            if not sub_list:
                print("⚠ No subaccounts found, skipping queries")
                return
            test_subaccount_id = sub_list[0].get("id", "")
        except Exception as e:
            print(f"⚠ Failed to get subaccount: {e}")
            return

        # Get active orders (skip if no valid subaccount)
        print("\n5.1 Getting active orders...")
        if test_subaccount_id == "":
            print("⚠ No valid subaccount, skipping 5.1 ~ 5.8 trading queries demo")
            return
        try:
            active_order_req = {
                "subaccountId": test_subaccount_id,
                "size": 10,
            }
            active_order_resp = client.get_active_order(active_order_req)
            order_list = active_order_resp.get("data", {}).get("orderList", [])
            if len(order_list) == 0:
                print("⚠ No active orders found")
            else:
                print(f"✓ Retrieved {len(order_list)} active orders")
                for i, order in enumerate(order_list[:3]):  # Show first 3
                    print(f"  Order {i+1}: ID={order.get('id', 'N/A')}, ExchangeId={order.get('exchangeId', 'N/A')}, "
                          f"Direction={'Buy' if order.get('isBuy') else 'Sell'}, "
                          f"Price={order.get('price', 'N/A')}, Size={order.get('size', 'N/A')}, "
                          f"Status={order.get('status', 'N/A')}")
        except Exception as e:
            print(f"⚠ Failed to get active orders: {e}")

        # Get history orders
        print("\n5.2 Getting history orders...")
        try:
            history_order_req = {
                "subaccountId": test_subaccount_id,
                "size": 10,
            }
            history_order_resp = client.get_history_order(history_order_req)
            history_order_list = history_order_resp.get("data", {}).get("orderList", [])
            print(f"✓ Retrieved {len(history_order_list)} history orders")
            for i, order in enumerate(history_order_list[:3]):  # Show first 3
                print(f"  History Order {i+1}: ID={order.get('id', 'N/A')}, ExchangeId={order.get('exchangeId', 'N/A')}, "
                      f"Direction={'Buy' if order.get('isBuy') else 'Sell'}, "
                      f"Price={order.get('price', 'N/A')}, Size={order.get('size', 'N/A')}, "
                      f"Status={order.get('status', 'N/A')}")
        except Exception as e:
            print(f"⚠ Failed to get history orders: {e}")

        # Get perpetual account asset
        print("\n5.3 Getting perpetual account asset...")
        try:
            asset_req = {
                "subaccountId": test_subaccount_id,
            }
            asset_resp = client.get_perpetual_account_asset(asset_req)
            data = asset_resp.get("data", {})
            collateral_list = data.get("collateralList", [])
            position_list = data.get("positionList", [])
            print(f"✓ Retrieved {len(collateral_list)} collaterals, {len(position_list)} positions")
            for i, collateral in enumerate(collateral_list[:3]):  # Show first 3
                print(f"  Collateral {i+1}: CoinId={collateral.get('coinId', 'N/A')}, "
                      f"Amount={collateral.get('amount', 'N/A')}")
            for i, position in enumerate(position_list[:3]):  # Show first 3
                print(f"  Position {i+1}: ExchangeId={position.get('exchangeId', 'N/A')}, "
                      f"OpenSize={position.get('openSize', 'N/A')}, "
                      f"OpenValue={position.get('openValue', 'N/A')}, "
                      f"FundingFee={position.get('fundingFee', 'N/A')}")
        except Exception as e:
            print(f"⚠ Failed to get account asset: {e}")

        # Get position transaction
        print("\n5.4 Getting position transaction...")
        try:
            position_req = {
                "subaccountId": test_subaccount_id,
                "size": 10,
            }
            position_resp = client.get_position_transaction(position_req)
            position_tx_list = position_resp.get("data", {}).get("positionTransactionList", [])
            print(f"✓ Retrieved {len(position_tx_list)} position transactions")
            for i, tx in enumerate(position_tx_list[:3]):  # Show first 3
                delta_open_size = tx.get("deltaOpenSize", "0")
                if delta_open_size == "":
                    delta_open_size = "0"
                delta_open_value = tx.get("deltaOpenValue", "0")
                if delta_open_value == "":
                    delta_open_value = "0"
                fill_size = tx.get("fillSize", "0")
                if fill_size == "":
                    fill_size = "0"
                fill_value = tx.get("fillValue", "0")
                if fill_value == "":
                    fill_value = "0"
                print(f"  Position Transaction {i+1}: ExchangeId={tx.get('exchangeId', 'N/A')}, "
                      f"DeltaOpenSize={delta_open_size}, DeltaOpenValue={delta_open_value}, "
                      f"FillSize={fill_size}, FillValue={fill_value}")
        except Exception as e:
            print(f"⚠ Failed to get position transaction: {e}")

        # Get collateral transaction
        print("\n5.5 Getting collateral transaction...")
        try:
            collateral_req = {
                "subaccountId": test_subaccount_id,
                "size": 10,
            }
            collateral_resp = client.get_collateral_transaction(collateral_req)
            collateral_tx_list = collateral_resp.get("data", {}).get("collateralTransactionList", [])
            print(f"✓ Retrieved {len(collateral_tx_list)} collateral transactions")
            for i, tx in enumerate(collateral_tx_list[:3]):  # Show first 3
                delta_amount = tx.get("deltaAmount", "0")
                if delta_amount == "":
                    delta_amount = "0"
                print(f"  Collateral Transaction {i+1}: CoinId={tx.get('coinId', 'N/A')}, "
                      f"DeltaAmount={delta_amount}, Type={tx.get('type', 'N/A')}")
        except Exception as e:
            print(f"⚠ Failed to get collateral transaction: {e}")

        # Get asset snapshot
        print("\n5.6 Getting asset snapshot...")
        try:
            snapshot_req = {
                "subaccountId": test_subaccount_id,
                "size": 10,
            }
            snapshot_resp = client.get_asset_snapshot(snapshot_req)
            snapshot_list = snapshot_resp.get("data", {}).get("assetSnapshotList", [])
            print(f"✓ Retrieved {len(snapshot_list)} asset snapshots")
            for i, snapshot in enumerate(snapshot_list[:3]):  # Show first 3
                print(f"  Asset Snapshot {i+1}: CoinId={snapshot.get('coinId', 'N/A')}, "
                      f"TotalEquity={snapshot.get('totalEquity', 'N/A')}, "
                      f"TotalRealizePnl={snapshot.get('totalRealizePnl', 'N/A')}, "
                      f"SnapshotTime={snapshot.get('snapshotTime', 'N/A')}")
        except Exception as e:
            print(f"⚠ Failed to get asset snapshot: {e}")

        # Get history order fill transaction
        print("\n5.7 Getting history order fill transaction...")
        try:
            fill_req = {
                "subaccountId": test_subaccount_id,
                "size": 10,
            }
            fill_resp = client.get_history_order_fill_transaction(fill_req)
            fill_tx_list = fill_resp.get("data", {}).get("orderFillTransactionList", [])
            print(f"✓ Retrieved {len(fill_tx_list)} order fill transactions")
            for i, fill in enumerate(fill_tx_list[:3]):  # Show first 3
                print(f"  Order Fill Transaction {i+1}: {fill}")
        except Exception as e:
            print(f"⚠ Failed to get history order fill transaction: {e}")

        # Get history position term
        print("\n5.8 Getting history position term...")
        try:
            term_req = {
                "subaccountId": test_subaccount_id,
                "size": 10,
            }
            term_resp = client.get_history_position_term(term_req)
            term_list = term_resp.get("data", {}).get("positionTermList", [])
            print(f"✓ Retrieved {len(term_list)} position terms")
            for i, term in enumerate(term_list[:3]):  # Show first 3
                print(f"  Position Term {i+1}: ExchangeId={term.get('exchangeId', 'N/A')}, "
                      f"TermCount={term.get('termCount', 'N/A')}, "
                      f"CumOpenSize={term.get('cumOpenSize', 'N/A')}, "
                      f"CumCloseSize={term.get('cumCloseSize', 'N/A')}")
        except Exception as e:
            print(f"⚠ Failed to get history position term: {e}")

    except Exception as e:
        print(f"⚠ Trading queries demo failed: {e}")


def main():
    """Run all demos"""
    print("=" * 60)
    print("Antex SDK Python - Complete Example")
    print("=" * 60)
    print(f"Gateway: {GATEWAY_URL}")
    print(f"Chain ID: {CHAIN_ID}")

    client = AntexClient(base_url=GATEWAY_URL, ws_url=WS_URL)

    # Run demos
    demo_basic_functions(client)
    demo_market_data(client)
    demo_websocket_realtime(client)
    demo_trading_functions(client)
    demo_trading_queries(client)

    print("\n" + "=" * 60)
    print("Example Complete!")
    print("=" * 60)


if __name__ == "__main__":
    main()
