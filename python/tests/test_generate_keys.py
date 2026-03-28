#!/usr/bin/env python3
"""
Test: generate ethPrivateKey and agentPrivateKey
"""

import os
import secrets
import sys

script_dir = os.path.dirname(os.path.abspath(__file__))
parent_dir = os.path.dirname(script_dir)
sys.path.insert(0, parent_dir)

from eth_account import Account
from antx_sdk.crypto import derive_antx_bech32_address


def test_generate_keys():
    # Generate ethPrivateKey
    eth_private_key = "0x" + secrets.token_hex(32)
    acct = Account.from_key(eth_private_key)
    eth_address = acct.address

    print("\n--- ETH Private Key ---")
    print(f"ETH_PRIVATE_KEY={eth_private_key}")
    print(f"ETH_ADDRESS={eth_address}")

    # Generate agentPrivateKey
    agent_private_key_bytes = secrets.token_bytes(32)
    agent_private_key = agent_private_key_bytes.hex()
    agent_address = derive_antx_bech32_address(agent_private_key_bytes)

    print("\n--- Agent Private Key ---")
    print(f"AGENT_PRIVATE_KEY={agent_private_key}")
    print(f"AGENT_ADDRESS={agent_address}")

    print("\n=== Done ===")
    print("Export these as environment variables before running examples:")
    print(f"  export ETH_PRIVATE_KEY={eth_private_key}")
    print(f"  export ETH_ADDRESS={eth_address}")
    print(f"  export AGENT_PRIVATE_KEY={agent_private_key}")

    assert len(eth_private_key) == 66  # 0x + 64 hex chars
    assert eth_address.startswith("0x") and len(eth_address) == 42
    assert len(agent_private_key) == 64
    assert agent_address.startswith("omni")
