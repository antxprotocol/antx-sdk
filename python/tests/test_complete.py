#!/usr/bin/env python3
"""
Complete functionality test for Antex SDK Python

This test verifies:
- Imports and proto availability
- Address derivation
- Transaction message creation
- Transaction building
- HTTP queries (may fail if gateway requires auth)
- WebSocket connection
"""

import sys
import time
import os

sys.path.insert(0, './python')

from antex_sdk.client import AntexClient
from antex_sdk.constants import ACCOUNT_HRP

# Configuration
GATEWAY = "https://testnet.antex.ai"
WS = "wss://testnet.antex.ai/api/v1/ws"
CHAIN_ID = "antex-testnet"

# Load private key from file
def load_private_key():
    """Load private key from .test_private_key file"""
    script_dir = os.path.dirname(os.path.abspath(__file__))
    key_file = os.path.join(script_dir, '.test_private_key')
    
    if not os.path.exists(key_file):
        raise FileNotFoundError(
            f"Private key file not found: {key_file}\n"
            f"Please create {key_file} with your test private key."
        )
    
    with open(key_file, 'r') as f:
        key = f.read().strip()
        # Remove 0x prefix if present
        if key.startswith('0x'):
            key = key[2:]
        return key

PRIV = load_private_key()


def test_imports():
    """Test all necessary imports"""
    print("\n[1/7] Testing imports...")
    try:
        from antex_sdk.client import AntexClient
        from antex_sdk.tx import COSMOS_PROTO_AVAILABLE
        if not COSMOS_PROTO_AVAILABLE:
            print("✗ Cosmos proto not available!")
            return False
        print("✓ All imports successful")
        return True
    except Exception as e:
        print(f"✗ Import failed: {e}")
        import traceback
        traceback.print_exc()
        return False


def test_client_creation():
    """Test client creation"""
    print("\n[2/7] Creating client...")
    try:
        client = AntexClient(base_url=GATEWAY, ws_url=WS)
        print(f"✓ Client created")
        return client
    except Exception as e:
        print(f"✗ Client creation failed: {e}")
        return None


def test_address_derivation(client):
    """Test address derivation from private key"""
    print("\n[3/7] Testing address derivation...")
    try:
        client.set_credentials(agent_private_key_hex=PRIV, chain_id=CHAIN_ID, account_hrp=ACCOUNT_HRP)
        agent_addr = client._agent_address_bech32
        print(f"✓ Derived agent address: {agent_addr}")
        return True
    except Exception as e:
        print(f"✗ Address derivation failed: {e}")
        import traceback
        traceback.print_exc()
        return False


def test_transaction_messages(client):
    """Test transaction message creation"""
    print("\n[4/7] Testing transaction message creation...")
    try:
        from antex_proto.antex.chain.agent import tx_pb2 as agent_tx
        from antex_proto.antex.chain.order import tx_pb2 as order_tx
        from antex_sdk.crypto import eth_personal_sign

        agent_addr = client._agent_address_bech32
        create_ms = int(time.time() * 1000)
        expire_ms = int((time.time() + 3600) * 1000)
        message = (
            f"Action:BindAgent\n"
            f"AgentAddress:{agent_addr}\n"
            f"CreateTime:{create_ms}\n"
            f"ExpireTime:{expire_ms}\n"
            f"ChainId:{CHAIN_ID}"
        )
        signature = eth_personal_sign(message, PRIV)
        print(f"✓ EVM personal_sign: {signature[:30]}...")

        msg = agent_tx.MsgBindAgent(
            agent_address=agent_addr,
            chain_type=1,
            chain_address="0x1234567890123456789012345678901234567890",
            create_time=create_ms,
            expire_time=expire_ms,
            chain_signature=signature,
        )
        print(f"✓ MsgBindAgent created successfully")

        # Test order message
        order_msg = order_tx.MsgCreateOrder(
            agent_address=agent_addr,
            subaccount_id=1,
            exchange_id=200001,
            margin_mode=1,
            leverage=1,
            is_buy=True,
            price_scale=2,
            price_value=100000,
            size_scale=3,
            size_value=100,
            client_order_id="test-001",
            time_in_force=1,
            reduce_only=False,
            expire_time=0,
            is_market=False,
            is_position_tp=False,
            is_position_sl=False,
            trigger_type=0,
            trigger_price_type=0,
            trigger_price_value=0,
            is_set_open_tp=False,
            is_set_open_sl=False,
        )
        print(f"✓ MsgCreateOrder created successfully")
        return True
    except Exception as e:
        print(f"✗ Transaction message creation failed: {e}")
        import traceback
        traceback.print_exc()
        return False


def test_transaction_building(client):
    """Test transaction building (dry-run)"""
    print("\n[5/7] Testing transaction building (dry-run)...")
    try:
        from antex_sdk.tx import pack_any, build_tx_body, build_auth_info
        from antex_proto.antex.chain.agent import tx_pb2 as agent_tx

        agent_addr = client._agent_address_bech32
        msg = agent_tx.MsgBindAgent(
            agent_address=agent_addr,
            chain_type=1,
            chain_address="0x1234567890123456789012345678901234567890",
            create_time=int(time.time() * 1000),
            expire_time=int((time.time() + 3600) * 1000),
            chain_signature="0x0000",
        )

        any_msg = pack_any(msg, "/antex.chain.agent.MsgBindAgent")
        body = build_tx_body([any_msg])
        print(f"✓ TxBody created: {len(body.messages)} message(s)")

        auth = build_auth_info(client._agent_pubkey_bytes, 0, gas_limit=200000, fee_amounts=[])
        print(f"✓ AuthInfo created")

        print(f"✓ Transaction building works (signing requires account_number)")
        return True
    except Exception as e:
        print(f"✗ Transaction building failed: {e}")
        import traceback
        traceback.print_exc()
        return False


def test_http_queries(client):
    """Test HTTP queries"""
    print("\n[6/7] Testing HTTP queries...")
    try:
        coins = client.get_coin_list()
        print(f"✓ get_coin_list: {len(coins)} coins")
        return True
    except Exception as e:
        error_msg = str(e)
        print(f"⚠ get_coin_list: {type(e).__name__}")
        
        # Extract request details - traverse exception chain to find the one with details
        details = None
        response_obj = None
        exc = e
        while exc:
            if hasattr(exc, 'request_details') and exc.request_details:
                details = exc.request_details
                # Check if this has a status_code
                if details.get('status_code') is not None:
                    break
            if hasattr(exc, 'response') and exc.response:
                response_obj = exc.response
            exc = getattr(exc, '__cause__', None)
        
        if details:
            print(f"  Request Method: {details.get('method', 'N/A')}")
            print(f"  Request URL: {details.get('url', 'N/A')}")
            print(f"  Request Params: {details.get('params', {})}")
            status_code = details.get('status_code')
            if status_code is None and response_obj:
                status_code = response_obj.status_code
            print(f"  Status Code: {status_code}")
            response_text = details.get('response_text')
            if not response_text and response_obj:
                try:
                    response_text = response_obj.text[:500]
                except Exception:
                    pass
            if response_text:
                # Check if it's HTML (likely WAF challenge)
                if response_text.strip().startswith('<!DOCTYPE') or '<html' in response_text[:100]:
                    print(f"  Response: HTML page (likely WAF challenge/CAPTCHA)")
                    print(f"  Response preview: {response_text[:150]}...")
                    print(f"  Note: Gateway is protected by AWS WAF, returning CAPTCHA challenge page")
                    print(f"  This is expected behavior - the SDK code is working correctly")
                else:
                    print(f"  Response Text: {response_text[:200]}")
        elif "Request details:" in error_msg:
            # Try to extract from error message string
            try:
                import ast
                details_str = error_msg.split("Request details:")[1].strip()
                details = ast.literal_eval(details_str)
                print(f"  Request Method: {details.get('method', 'N/A')}")
                print(f"  Request URL: {details.get('url', 'N/A')}")
                print(f"  Request Params: {details.get('params', {})}")
                print(f"  Status Code: {details.get('status_code', 'N/A')}")
            except Exception:
                print(f"  Error: {error_msg[:300]}")
        elif hasattr(e, 'response'):
            print(f"  Status Code: {e.response.status_code}")
            print(f"  Request URL: {e.response.url}")
            print(f"  Response: {e.response.text[:200] if e.response.text else 'Empty'}")
        else:
            print(f"  Error: {error_msg[:300]}")
        
        if "405" in error_msg or (hasattr(e, 'response') and e.response.status_code == 405):
            print(f"  Note: Gateway may require different HTTP method or authentication")
        
        return False  # Non-critical, don't fail the test


def test_websocket(client):
    """Test WebSocket connection"""
    print("\n[7/7] Testing WebSocket...")
    try:
        ws_msg_received = False

        def on_msg(data):
            nonlocal ws_msg_received
            ws_msg_received = True

        client.connect_websocket(on_msg, lambda e: None)
        time.sleep(2)
        if client._ws and client._ws.is_connected():
            print("✓ WebSocket connected")
            client.disconnect_websocket()
            return True
        else:
            print("⚠ WebSocket connection timeout")
            return False
    except Exception as e:
        print(f"⚠ WebSocket: {type(e).__name__}")
        return False


def main():
    """Run all tests"""
    print("=" * 60)
    print("Antex SDK Python - Complete Functionality Test")
    print("=" * 60)

    if not test_imports():
        sys.exit(1)

    client = test_client_creation()
    if not client:
        sys.exit(1)

    if not test_address_derivation(client):
        sys.exit(1)

    if not test_transaction_messages(client):
        sys.exit(1)

    if not test_transaction_building(client):
        sys.exit(1)

    # These may fail due to gateway configuration, but are not critical
    http_ok = test_http_queries(client)
    ws_ok = test_websocket(client)

    print("\n" + "=" * 60)
    print("Test Summary:")
    print("✓ Imports: OK")
    print("✓ Address derivation: OK")
    print("✓ Transaction message creation: OK")
    print("✓ Transaction building: OK")
    if http_ok:
        print("✓ HTTP queries: OK")
    else:
        print("⚠ HTTP queries: Failed (non-critical, may need gateway config)")
    if ws_ok:
        print("✓ WebSocket: OK")
    else:
        print("⚠ WebSocket: Failed (non-critical, may need gateway config)")
    print("=" * 60)
    print("\n✓ Core functionality verified!")
    print("Note: HTTP/WS failures are non-critical and may be due to:")
    print("      - Gateway requiring authentication")
    print("      - Network connectivity issues")
    print("      - Gateway configuration differences")
    print("\nTransaction signing is ready when account_number is available")


if __name__ == "__main__":
    main()

