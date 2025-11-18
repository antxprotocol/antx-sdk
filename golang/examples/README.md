# Antx SDK Examples

This directory contains complete usage examples for the Antx SDK.

## Example Files

### complete_example.go

This is a complete example file that demonstrates all major features of the Antx SDK:

1. **Basic Functions Demo**
   - Get coin list
   - Get exchange list

2. **Market Data Functions Demo**
   - Get K-line data
   - Get funding rate history

3. **WebSocket Real-time Market Data Demo**
   - Subscribe to Ticker data
   - Real-time market data parsing

4. **Trading Functions Demo**
   - Bind agent
   - Get subaccount list
   - Create order
   - Cancel order

5. **Trading Query Functions Demo**
   - Active order query
   - History order query
   - Account asset query
   - Position transaction query
   - Collateral transaction query
   - Asset snapshot query
   - Order fill transaction query
   - History position term query

## Running Examples

```bash
# Navigate to examples directory
cd examples

# Run complete example
go run complete_example.go
```

## Configuration

Before running examples, please ensure:

1. **Update configuration**: Modify configuration constants in `complete_example.go`
   ```go
   const (
       chainID         = "antex-testnet"  // Chain ID
       ethPrivateKey   = "your_eth_private_key"   // Ethereum private key
       agentPrivateKey = "your_agent_private_key" // Agent private key
   )
   ```

2. **Ensure network connectivity**: Examples need to connect to Antx testnet

3. **Prepare test funds**: Some functions (like creating orders) require sufficient account balance

## Function Reference

### Basic Functions
- `GetCoinList()` - Get supported coin list
- `GetExchangeList()` - Get exchange list
- `GetSubaccountList()` - Get subaccount list

### Market Data Functions
- `GetKline()` - Get K-line data
- `GetFundingHistory()` - Get funding rate history
- WebSocket real-time subscription functions

### Trading Functions
- `BindAgent()` - Bind agent
- `CreateOrder()` - Create order
- `CancelOrderByClientId()` - Cancel order
- `GetTransactionResult()` - Query transaction result

### Trading Query Functions
- `GetActiveOrder()` - Get active orders
- `GetHistoryOrder()` - Get history orders
- `GetPerpetualAccountAsset()` - Get perpetual contract account assets
- `GetPositionTransaction()` - Get position transactions
- `GetCollateralTransaction()` - Get collateral transactions
- `GetAssetSnapshot()` - Get asset snapshots
- `GetHistoryOrderFillTransaction()` - Get history order fill transactions
- `GetHistoryPositionTerm()` - Get history position terms

## Numeric Processing

### Using Decimal for Precise Calculations
For all fields involving amounts, prices, quantities, etc. that require precise calculations, it is recommended to use the `github.com/shopspring/decimal` library:

```go
import "github.com/shopspring/decimal"

// Parse string to decimal
total, err := decimal.NewFromString(asset.Total)
if err != nil {
    log.Printf("Failed to parse total: %v", err)
    return
}

// Perform precise calculations
available, err := decimal.NewFromString(asset.Available)
if err != nil {
    log.Printf("Failed to parse available: %v", err)
    return
}

// Calculate frozen assets
frozen := total.Sub(available)

// Convert to string
fmt.Printf("Frozen assets: %s", frozen.String())
```

## Usage Examples

### Basic Query Example

```go
// Query active orders
activeOrderReq := sdk.GetActiveOrderReq{
    SubaccountId: "123456", // Subaccount ID
    Size:         10,       // Number of records
}

activeOrderResp, err := client.GetActiveOrder(activeOrderReq)
if err != nil {
    log.Printf("Failed to get active orders: %v", err)
    return
}

fmt.Printf("Retrieved %d active orders\n", len(activeOrderResp.Data.OrderList))
for _, order := range activeOrderResp.Data.OrderList {
    fmt.Printf("Order ID: %s, Exchange: %s, Direction: %s, Price: %d, Size: %d\n",
        order.Id, order.ExchangeId,
        map[bool]string{true: "Buy", false: "Sell"}[order.IsBuy],
        order.PriceValue, order.SizeValue)
}
```

### Advanced Filter Query

```go
// Query active orders for specific exchanges
activeOrderReq := sdk.GetActiveOrderReq{
    SubaccountId:         "123456",
    Size:                 20,
    FilterExchangeIdList: "200001,200002", // Only query BTC-USDT and ETH-USDT
    FilterOrderStatusList: "1,2",          // Only query pending and filled orders
    FilterStartCreatedTimeInclusive: uint64(time.Now().Add(-24*time.Hour).UnixMilli()),
    FilterEndCreatedTimeExclusive:   uint64(time.Now().UnixMilli()),
}

// Query position transactions
positionReq := sdk.GetPositionTransactionReq{
    SubaccountId:         "123456",
    Size:                 50,
    FilterExchangeIdList: "200001", // Only query BTC-USDT position transactions
    FilterTypeList:       "1,2",    // Only query open and close transactions
    FilterStartCreatedTimeInclusive: uint64(time.Now().Add(-7*24*time.Hour).UnixMilli()),
}
```

## Notes

1. **Private Key Security**: Private keys in examples are for demonstration only, do not use in production
2. **Network Environment**: Ensure access to Antx testnet
3. **Error Handling**: Examples include basic error handling, extend as needed for actual use
4. **Concurrency Safety**: WebSocket client is concurrency-safe and can be used in multiple goroutines
5. **Numeric Precision**: All numeric fields are in string format, use decimal library for precise calculations
6. **Pagination Limit**: Maximum 100 records per query
7. **Time Format**: Time parameters use millisecond timestamps

## More Information

- For detailed API documentation, refer to `MARKET_DATA_README.md`
- SDK source code is located in the project root directory
- For issues, please check project documentation or submit an Issue
