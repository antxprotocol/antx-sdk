# Antx SDK Python Tests

## Running Tests

### Prerequisites

Install the SDK and dependencies (use the project venv or your own):

```bash
cd python
python3 -m venv .venv   # if not already created
.venv/bin/pip install -e .
```

### Complete Functionality Test

Run the complete test suite **with the venv** (so that `requests` and other deps are available):

```bash
cd python
.venv/bin/python tests/test_complete.py
```

Or activate the venv first, then run:

```bash
cd python
source .venv/bin/activate
python tests/test_complete.py
```

To exit the venv when done, run: `deactivate`.

This test verifies:
- ✓ Imports and proto availability
- ✓ Address derivation from private key
- ✓ Transaction message creation (MsgBindAgent, MsgCreateOrder)
- ✓ Transaction building (TxBody, AuthInfo)
- ⚠ HTTP queries (may fail if gateway requires authentication)
- ⚠ WebSocket connection (may timeout)

## Test Configuration

The test uses the following default configuration:
- Gateway: `https://devnet.antxfi.com`
- WebSocket: `wss://devnet.antxfi.com/api/v1/ws`
- Chain ID: `omni-devnet`
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

