import base64
from dataclasses import dataclass
from typing import List, Optional

from google.protobuf.any_pb2 import Any as ProtoAny
from google.protobuf import timestamp_pb2

# Expect these modules to be generated into python/antex_proto/**
# Try importing cosmos proto, fallback to None if not available
try:
    from antex_proto.cosmos.tx.v1beta1.tx_pb2 import TxBody, AuthInfo, TxRaw, SignDoc
    from antex_proto.cosmos.tx.signing.v1beta1.signing_pb2 import SignMode
    from antex_proto.cosmos.base.v1beta1.coin_pb2 import Coin
    from antex_proto.cosmos.crypto.secp256k1.keys_pb2 import PubKey as Secp256k1PubKey
    from antex_proto.cosmos.tx.v1beta1.tx_pb2 import SignerInfo, ModeInfo, Fee
    COSMOS_PROTO_AVAILABLE = True
except ImportError:
    # Cosmos proto not generated yet - will fail at runtime when tx functions are called
    COSMOS_PROTO_AVAILABLE = False
    TxBody = None
    AuthInfo = None
    TxRaw = None
    SignDoc = None
    SignMode = None
    Coin = None
    Secp256k1PubKey = None
    SignerInfo = None
    ModeInfo = None
    Fee = None


@dataclass
class Signer:
    private_key_bytes: bytes  # 32 bytes secp256k1
    public_key_bytes: bytes   # 33/65 bytes compressed/uncompressed supported by PubKey key field


def pack_any(msg, type_url: str) -> ProtoAny:
    any_msg = ProtoAny()
    any_msg.Pack(msg)
    # Force chain-specific type_url
    any_msg.type_url = type_url
    return any_msg


def build_tx_body(msgs: List[ProtoAny], memo: str = "", timeout_height: int = 0, 
                  unordered: bool = False, timeout_timestamp_ns: int = 0) -> TxBody:
    if not COSMOS_PROTO_AVAILABLE:
        raise RuntimeError("cosmos proto modules not available; ensure proto generation includes cosmos dependencies")
    body = TxBody(messages=msgs, memo=memo)
    if timeout_height:
        body.timeout_height = timeout_height
    if unordered:
        body.unordered = True
    if timeout_timestamp_ns:
        timestamp = timestamp_pb2.Timestamp()
        timestamp.FromNanoseconds(timeout_timestamp_ns)
        body.timeout_timestamp.CopyFrom(timestamp)
    return body


def build_auth_info(pubkey_bytes: bytes, sequence: int, gas_limit: int = 200000, fee_amounts: Optional[List[Coin]] = None) -> AuthInfo:
    pub = Secp256k1PubKey(key=pubkey_bytes)
    pub_any = ProtoAny()
    pub_any.Pack(pub)
    pub_any.type_url = "/cosmos.crypto.secp256k1.PubKey"
    signer_info = SignerInfo(
        public_key=pub_any,
        mode_info=ModeInfo(single=ModeInfo.Single(mode=SignMode.SIGN_MODE_DIRECT)),
        sequence=sequence,
    )
    fee = Fee(gas_limit=gas_limit, amount=fee_amounts or [])
    return AuthInfo(signer_infos=[signer_info], fee=fee)


def sign_tx(body: TxBody, auth_info: AuthInfo, chain_id: str, account_number: int, signer: Signer) -> bytes:
    sign_doc = SignDoc(
        body_bytes=body.SerializeToString(),
        auth_info_bytes=auth_info.SerializeToString(),
        chain_id=chain_id,
        account_number=account_number,
    )
    sign_bytes = sign_doc.SerializeToString()

    import hashlib
    msg_hash = hashlib.sha256(sign_bytes).digest()
    
    try:
        import coincurve
        privkey = coincurve.PrivateKey(signer.private_key_bytes)
        signature = privkey.sign_recoverable(msg_hash, hasher=None)[:64]
    except ImportError:
        try:
            from eth_keys import keys
            private_key = keys.PrivateKey(signer.private_key_bytes)
            sig_obj = private_key.sign_msg_hash(msg_hash)
            r_bytes = sig_obj.r.to_bytes(32, byteorder='big')
            s_bytes = sig_obj.s.to_bytes(32, byteorder='big')
            signature = r_bytes + s_bytes
        except (ImportError, AttributeError, Exception):
            try:
                from ecdsa import SigningKey, SECP256k1
                from ecdsa.util import sigdecode_der
                
                sk = SigningKey.from_string(signer.private_key_bytes, curve=SECP256k1)
                der_sig = sk.sign_deterministic(sign_bytes, hashfunc=hashlib.sha256)
                r, s = sigdecode_der(der_sig, SECP256k1.order)
                r_bytes = r.to_bytes(32, byteorder='big')
                s_bytes = s.to_bytes(32, byteorder='big')
                signature = r_bytes + s_bytes
            except Exception as e:  # noqa: BLE001
                raise RuntimeError(
                    f"failed to sign transaction: {e}. "
                    "Please install 'coincurve' for best compatibility: pip install coincurve"
                )

    if len(signature) != 64:
        raise RuntimeError(f"invalid signature length: expected 64 bytes, got {len(signature)}")

    tx_raw = TxRaw(
        body_bytes=sign_doc.body_bytes,
        auth_info_bytes=sign_doc.auth_info_bytes,
        signatures=[signature],
    )
    return tx_raw.SerializeToString()


def encode_tx_base64(tx_bytes: bytes) -> str:
    return base64.b64encode(tx_bytes).decode()


