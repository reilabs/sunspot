from ecdsa import SigningKey
# from ecdsa import NIST256p
from ecdsa import SECP256k1
from ecdsa.util import sigencode_string
import hashlib

# Fixed private key (32 bytes, test only)
priv_hex = "1" * 64
priv_bytes = bytes.fromhex(priv_hex)

# sk = SigningKey.from_string(priv_bytes, curve=NIST256p)
sk = SigningKey.from_string(priv_bytes, curve=SECP256k1)
vk = sk.get_verifying_key()

msg = b"hello world"

# SHA-256 digest of the message
msg_hash = hashlib.sha256(msg).digest()

# Deterministic ECDSA signature, raw r||s (64 bytes)
sig_raw = sk.sign_deterministic(
    msg,
    hashfunc=hashlib.sha256,
    sigencode=sigencode_string
)

# Split into r and s
r_bytes = sig_raw[:32]
s_bytes = sig_raw[32:]

# Public key X and Y (uncompressed coordinates)
x_bytes = vk.to_string()[:32]
y_bytes = vk.to_string()[32:]

# Compressed public key (33 bytes)
prefix = b'\x02' if int.from_bytes(y_bytes, "big") % 2 == 0 else b'\x03'
compressed_pub = prefix + x_bytes

# Helper for pretty byte array output
def to_byte_array(data: bytes) -> str:
    return "[" + ", ".join(f"0x{b:02x}" for b in data) + "]"

print("Public key X (32B):", to_byte_array(x_bytes))
print("Public key Y (32B):", to_byte_array(y_bytes))
print("Message SHA-256 (32B):", to_byte_array(msg_hash))
print("Signature (64B r||s):", to_byte_array(sig_raw))

# Verify the raw r||s signature
assert vk.verify(
    sig_raw,
    msg,
    hashfunc=hashlib.sha256,
    sigdecode=lambda s, order: (
        int.from_bytes(s[:32], 'big'),
        int.from_bytes(s[32:], 'big')
    )
)
print("verified ✅")
