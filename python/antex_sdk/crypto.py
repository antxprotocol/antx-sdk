import binascii
from typing import Tuple

from bech32 import bech32_decode, bech32_encode, convertbits
from eth_account import Account
from eth_account.messages import encode_defunct
from eth_utils import to_checksum_address
from ecdsa import SECP256k1, SigningKey
import hashlib

from .constants import ACCOUNT_HRP


def _bech32_to_bytes(addr: str) -> Tuple[str, bytes]:
    hrp, data = bech32_decode(addr)
    if hrp is None or data is None:
        raise ValueError("invalid bech32 address")
    decoded = convertbits(data, 5, 8, False)
    if decoded is None:
        raise ValueError("failed to convert bech32 data")
    return hrp, bytes(decoded)


def _bytes_to_bech32(hrp: str, payload: bytes) -> str:
    data5 = convertbits(list(payload), 8, 5, True)
    if data5 is None:
        raise ValueError("failed to convert to bech32 bits")
    return bech32_encode(hrp, data5)


def convert_to_eth_addr(addr: str, account_hrp: str = ACCOUNT_HRP) -> str:
    if not addr:
        raise ValueError("addr can't be empty")
    if addr.startswith("0x") and len(addr) == 42:
        return to_checksum_address(addr)
    # try bech32
    hrp, payload = _bech32_to_bytes(addr)
    return to_checksum_address("0x" + payload.hex())


def convert_to_antex_addr(addr: str, account_hrp: str = ACCOUNT_HRP) -> str:
    if not addr:
        raise ValueError("addr can't be empty")
    if addr.startswith("0x") and len(addr) == 42:
        raw = bytes.fromhex(addr[2:])
        return _bytes_to_bech32(account_hrp, raw)
    # already bech32
    hrp, payload = _bech32_to_bytes(addr)
    if hrp != account_hrp:
        # re-encode to target hrp
        return _bytes_to_bech32(account_hrp, payload)
    return addr


def eth_personal_sign(message: str, eth_private_key_hex: str) -> str:
    key = eth_private_key_hex[2:] if eth_private_key_hex.startswith("0x") else eth_private_key_hex
    if len(key) != 64:
        raise ValueError("invalid eth private key length")
    acct = Account.from_key(bytes.fromhex(key))
    signable = encode_defunct(text=message)
    signed = Account.sign_message(signable, private_key=acct.key)
    return "0x" + signed.signature.hex()


def verify_eth_personal_signature(address: str, data: bytes, sig_hex: str) -> bool:
    try:
        msg = encode_defunct(text=data.decode() if isinstance(data, (bytes, bytearray)) else str(data))
        recovered = Account.recover_message(msg, signature=sig_hex)
        return to_checksum_address(recovered) == to_checksum_address(address)
    except Exception:
        return False


def secp256k1_pubkey_compressed(private_key_bytes: bytes) -> bytes:
    sk = SigningKey.from_string(private_key_bytes, curve=SECP256k1)
    vk = sk.verifying_key
    # Compressed format: 33 bytes
    return vk.to_string("compressed")


def _ripemd160(data: bytes) -> bytes:
    try:
        h = hashlib.new('ripemd160')
        h.update(data)
        return h.digest()
    except Exception:
        try:
            from Crypto.Hash import RIPEMD  # type: ignore
            h = RIPEMD.new()
            h.update(data)
            return h.digest()
        except Exception as e:
            raise RuntimeError(f"ripemd160 not available: {e}")


def derive_antex_bech32_address(private_key_bytes: bytes, account_hrp: str = ACCOUNT_HRP) -> str:
    pub = secp256k1_pubkey_compressed(private_key_bytes)
    sha = hashlib.sha256(pub).digest()
    ripe = _ripemd160(sha)
    return _bytes_to_bech32(account_hrp, ripe)


