# Antx SDK Python Tests

## Running Tests

### Complete Functionality Test

Run the complete test suite:

```bash
cd python
python tests/test_complete.py
```

Or directly:

```bash
python python/tests/test_complete.py
```

This test verifies:
- ✓ Imports and proto availability
- ✓ Address derivation from private key
- ✓ Transaction message creation (MsgBindAgent, MsgCreateOrder)
- ✓ Transaction building (TxBody, AuthInfo)
- ⚠ HTTP queries (may fail if gateway requires authentication)
- ⚠ WebSocket connection (may timeout)

## Test Configuration

The test uses the following default configuration:
- Gateway: `https://testnet.antxfi.com`
- WebSocket: `wss://testnet.antxfi.com/api/v1/ws`
- Chain ID: `antx-testnet`
- Test private key: Loaded from `tests/.test_private_key` file

### Setting Up Private Key

Before running tests, create a private key file:

```bash
cd python/tests
echo "xxxxxxxxx" > .test_private_key
```

**Note:** The `.test_private_key` file is ignored by git (see `.gitignore`). 
Never commit private keys to the repository.

You can modify other configuration values in `test_complete.py` if needed.

