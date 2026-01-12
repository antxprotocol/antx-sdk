package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	sdk "github.com/antxprotocol/antx-sdk-golang"
	"github.com/antxprotocol/antx-sdk-golang/constants"
	"github.com/antxprotocol/antx-sdk-golang/types"
	"github.com/shopspring/decimal"
)

var (
	// Basic configuration
	gatewayURL = "https://testnet.antxfi.com"
	wsURL      = "wss://testnet.antxfi.com/api/v1/ws"
	chainID    = "antx-testnet"

	// Credential configuration (example uses real test credentials)
	ethPrivateKey   = ""
	agentPrivateKey = ""
	ethAddress      = ""

	// Example default parameters
	defaultExchangeId = "200001"
)

func init() {
	if key, exist := os.LookupEnv("ETH_PRIVATE_KEY"); exist {
		ethPrivateKey = key
	}
	if key, exist := os.LookupEnv("AGENT_PRIVATE_KEY"); exist {
		agentPrivateKey = key
	}
	if key, exist := os.LookupEnv("ETH_ADDRESS"); exist {
		ethAddress = key
	}
}

func main() {
	// Create SDK client (using top-level unified variables)
	client, err := sdk.NewAntxClient(sdk.Config{
		GatewayHost:     gatewayURL,
		ChainID:         chainID,
		EthPrivateKey:   ethPrivateKey,
		AgentPrivateKey: agentPrivateKey,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	client.SetGateway(gatewayURL, wsURL)

	fmt.Println("=== Antx SDK Complete Example ===")
	fmt.Printf("SDK client created, base URL: %s\n", gatewayURL)

	// Call various demo functions
	demoBasicFunctions(client)
	demoMarketData(client)
	demoWebSocketRealtime(client)
	demoTradingFunctions(client)
	demoTradingQueries(client)

	fmt.Println("\n=== Example Complete ===")
}

// demoBasicFunctions demonstrates basic functions
func demoBasicFunctions(client *sdk.AntxClient) {
	fmt.Println("\n=== 1. Basic Functions Demo ===")

	// Get coin list
	fmt.Println("\n1.1 Getting coin list:")
	coinList, err := client.GetCoinList()
	if err != nil {
		log.Printf("Failed to get coin list: %v", err)
	} else {
		fmt.Printf("Retrieved %d coins:\n", len(coinList))
		for i, coin := range coinList {
			if i >= 3 { // Only show first 3
				break
			}
			fmt.Printf("  Coin %d: ID=%s, Symbol=%s, Precision=%d\n",
				i+1, coin.Id, coin.Symbol, coin.StepSizeScale)
		}
	}

	// Get exchange list
	fmt.Println("\n1.2 Getting exchange list:")
	exchangeList, err := client.GetExchangeList()
	if err != nil {
		log.Printf("Failed to get exchange list: %v", err)
	} else {
		fmt.Printf("Retrieved %d exchanges:\n", len(exchangeList))
		for i, exchange := range exchangeList {
			if i >= 3 { // Only show first 3
				break
			}
			fmt.Printf("  Exchange %d: ID=%s, Symbol=%s, BaseCoin=%s, QuoteCoin=%s\n",
				i+1, exchange.Id, exchange.Symbol, exchange.BaseCoinId, exchange.QuoteCoinId)
		}
	}
}

// demoMarketData demonstrates market data functions
func demoMarketData(client *sdk.AntxClient) {
	fmt.Println("\n=== 2. Market Data Functions Demo ===")

	// Get K-line data
	fmt.Println("\n2.1 Getting K-line data:")
	klineReq := types.GetKLineReq{
		ExchangeId: defaultExchangeId,          // Exchange ID
		KlineType:  constants.KlineTypeMinute1, // 1-minute K-line
		PriceType:  constants.PriceTypeLast,    // Latest price
		Size:       10,                         // Get 10 records
	}

	klineResp, err := client.GetKline(klineReq)
	if err != nil {
		log.Printf("Failed to get kline data: %v", err)
	} else {
		fmt.Printf("Retrieved %d K-line records:\n", len(klineResp.Data.KlineList))
		for i, kline := range klineResp.Data.KlineList {
			if i >= 3 { // Only show first 3
				break
			}
			fmt.Printf("  K-line %d: Time=%s, Open=%s, High=%s, Low=%s, Close=%s, Volume=%s\n",
				i+1,
				time.Unix(int64(kline.KlineTime/1000), int64(kline.KlineTime%1000)*1e6).Format("2006-01-02 15:04:05"),
				kline.Open,
				kline.High,
				kline.Low,
				kline.Close,
				kline.Size,
			)
		}
	}

	// Get funding rate history
	fmt.Println("\n2.2 Getting funding rate history:")
	fundingReq := types.GetFundingHistoryReq{
		ExchangeId: defaultExchangeId,
		Size:       5,
	}

	fundingResp, err := client.GetFundingHistory(fundingReq)
	if err != nil {
		log.Printf("Failed to get funding history: %v", err)
	} else {
		fmt.Printf("Retrieved %d funding rate records:\n", len(fundingResp.Data.FundingRateList))
		for i, rate := range fundingResp.Data.FundingRateList {
			if i >= 3 { // Only show first 3
				break
			}
			fmt.Printf("  Funding Rate %d: Time=%s, Rate=%s, OraclePrice=%s\n",
				i+1,
				time.Unix(int64(rate.FundingTime/1000), 0).Format("2006-01-02 15:04:05"),
				rate.FundingRate,
				rate.OraclePrice,
			)
		}
	}
}

// demoWebSocketRealtime demonstrates WebSocket real-time market data
func demoWebSocketRealtime(client *sdk.AntxClient) {
	fmt.Println("\n=== 3. WebSocket Real-time Market Data Demo ===")

	// Connect WebSocket
	messageHandler := func(data []byte) {
		fmt.Printf("Received WebSocket message: %s\n", string(data))
	}
	errorHandler := func(err error) {
		fmt.Printf("WebSocket error: %v\n", err)
	}

	// If empty, use wsURL configured during initialization
	err := client.ConnectWebSocket(messageHandler, errorHandler)
	if err != nil {
		log.Printf("Failed to connect WebSocket: %v", err)
	} else {
		fmt.Println("WebSocket connected successfully")

		// Subscribe to Ticker data
		tickerChan, err := client.SubscribeToTicker("200001")
		if err != nil {
			log.Printf("Failed to subscribe to Ticker: %v", err)
		} else {
			fmt.Println("Subscribed to Ticker data, waiting for data...")

			// Set timeout to avoid program waiting indefinitely
			timeout := time.After(5 * time.Second)

			select {
			case data := <-tickerChan:
				fmt.Println("Received Ticker data:")
				tickerData, err := client.ParseTickerData(data)
				if err != nil {
					log.Printf("Failed to parse Ticker data: %v", err)
				} else {
					fmt.Printf("  Exchange: %s, LastPrice: %s, 24h Change: %s%%\n",
						tickerData.ExchangeId,
						tickerData.LastPrice,
						tickerData.PriceChangePercent,
					)
				}
			case <-timeout:
				fmt.Println("Timeout waiting for Ticker data")
			}
		}

		// Disconnect WebSocket
		client.DisconnectWebSocket()
		fmt.Println("WebSocket disconnected")
	}
}

// demoTradingFunctions demonstrates trading functions
func demoTradingFunctions(client *sdk.AntxClient) {
	fmt.Println("\n=== 4. Trading Functions Demo ===")
	fmt.Println("Note: Trading functions require agent binding first")

	// Perform bindAgent operation first
	fmt.Println("\n4.0 Binding agent:")
	fmt.Printf("Using ETH address: %s\n", ethAddress)
	agentAddress := client.GetAgentAddress()
	fmt.Printf("Agent address: %s\n", agentAddress)

	// Bind agent, set expiration time to 1 hour
	txHash, err := client.BindAgent(ethPrivateKey, chainID, 3600)
	if err != nil {
		log.Printf("Failed to bind agent: %v", err)
		fmt.Println("Skipping trading functions demo")
		return
	}
	fmt.Printf("Agent bound successfully, transaction hash: %s\n", txHash)

	// Dynamically get a valid subaccount ID
	testSubaccountId := ""
	subList, err := client.GetSubaccountList(1, ethAddress, agentAddress)
	if err != nil {
		log.Printf("Failed to get subaccount list: %v", err)
		fmt.Println("Skipping trading functions demo")
		return
	} else if len(subList) == 0 {
		log.Printf("No available subaccounts, skipping trading functions demo")
		return
	}
	testSubaccountId = subList[0].Id
	fmt.Printf("Found subaccount: %s\n", testSubaccountId)

	// Convert subaccount ID to uint64
	subaccountIdUint, err := strconv.ParseUint(testSubaccountId, 10, 64)
	if err != nil {
		log.Printf("Failed to convert subaccount ID: %v", err)
		return
	}

	// Convert exchange ID to uint64
	exchangeIdUint, err := strconv.ParseUint(defaultExchangeId, 10, 64)
	if err != nil {
		log.Printf("Failed to convert exchange ID: %v", err)
		return
	}

	// 4.1 Create limit buy order
	fmt.Println("\n4.1 Creating limit buy order:")
	createOrderReq := types.CreateOrderParam{
		SubaccountId:      subaccountIdUint,
		ExchangeId:        exchangeIdUint,
		MarginMode:        1, // Full margin mode
		Leverage:          1, // 1x leverage
		IsBuy:             true,
		PriceScale:        2,      // Price precision: 2 decimal places
		PriceValue:        100000, // Price 1000.00 (100000/100)
		SizeScale:         3,      // Size precision: 3 decimal places
		SizeValue:         100,    // Size 0.100 (100/1000)
		ClientOrderId:     "test-order-001",
		TimeInForce:       1, // GTC
		ReduceOnly:        false,
		ExpireTime:        uint64(time.Now().Add(24 * time.Hour).Unix()), // Expires in 24 hours
		IsMarket:          false,
		IsPositionTp:      false,
		IsPositionSl:      false,
		TriggerType:       0,
		TriggerPriceType:  0,
		TriggerPriceValue: 0,
		IsSetOpenTp:       false,
		IsSetOpenSl:       false,
	}

	orderTxHash, err := client.CreateOrder(&createOrderReq)
	if err != nil {
		log.Printf("Failed to create order: %v", err)
	} else {
		fmt.Printf("Order created successfully, transaction hash: %s\n", orderTxHash)
		// Wait for transaction confirmation and brief pause to avoid account sequence conflicts in subsequent transactions
		if orderTxHash != "" {
			fmt.Println("Waiting for transaction confirmation and brief delay...")
			time.Sleep(3 * time.Second)
		}
	}

	// 4.2 Create market sell order
	fmt.Println("\n4.2 Creating market sell order:")
	marketOrderReq := types.CreateOrderParam{
		SubaccountId:      subaccountIdUint,
		ExchangeId:        exchangeIdUint,
		MarginMode:        1, // Full margin mode
		Leverage:          1, // 1x leverage
		IsBuy:             false,
		PriceScale:        2,
		PriceValue:        0, // Market order price is 0
		SizeScale:         3,
		SizeValue:         50, // Size 0.050
		ClientOrderId:     "test-market-order-001",
		TimeInForce:       3, // IOC more suitable for market orders
		ReduceOnly:        false,
		ExpireTime:        uint64(time.Now().Add(24 * time.Hour).Unix()), // Expires in 24 hours
		IsMarket:          true,                                          // Market order
		IsPositionTp:      false,
		IsPositionSl:      false,
		TriggerType:       0,
		TriggerPriceType:  0,
		TriggerPriceValue: 0,
		IsSetOpenTp:       false,
		IsSetOpenSl:       false,
	}

	marketOrderTxHash, err := client.CreateOrder(&marketOrderReq)
	if err != nil {
		log.Printf("Failed to create market order: %v", err)
	} else {
		fmt.Printf("Market order created successfully, transaction hash: %s\n", marketOrderTxHash)
	}

	// 4.3 Cancel order
	fmt.Println("\n4.3 Canceling order:")
	orderIdUint, err := strconv.ParseUint("188531408901", 10, 64)
	if err != nil {
		log.Printf("Failed to convert order ID: %v", err)
	} else {
		cancelOrderReq := types.CancelOrderParam{
			SubaccountId: subaccountIdUint,
			OrderIdList:  []uint64{orderIdUint},
		}

		cancelTxHash, err := client.CancelOrder(&cancelOrderReq)
		if err != nil {
			log.Printf("Failed to cancel order: %v", err)
		} else {
			fmt.Printf("Order canceled successfully, transaction hash: %s\n", cancelTxHash)
		}
	}

	// 4.4 Create batch orders
	fmt.Println("\n4.4 Creating batch orders:")
	batchOrderReq := types.CreateOrderBatchParam{
		AgentAddress: client.GetAgentAddress(),
		SubaccountId: subaccountIdUint,
		ExchangeId:   exchangeIdUint,
		MarginMode:   1,
		Leverage:     1,
		CreateOrderParam: []*types.CreateOrderBatchDetail{
			{
				IsBuy:             true,
				PriceScale:        2,
				PriceValue:        95000, // Price 950.00
				SizeScale:         3,
				SizeValue:         200, // Size 0.200
				ClientOrderId:     "batch-order-001",
				TimeInForce:       1,
				ReduceOnly:        false,
				ExpireTime:        uint64(time.Now().Add(24 * time.Hour).Unix()), // Expires in 24 hours
				IsMarket:          false,
				IsPositionTp:      false,
				IsPositionSl:      false,
				TriggerType:       0,
				TriggerPriceType:  0,
				TriggerPriceValue: 0,
				IsSetOpenTp:       false,
				IsSetOpenSl:       false,
			},
			{
				IsBuy:             false,
				PriceScale:        2,
				PriceValue:        105000, // Price 1050.00
				SizeScale:         3,
				SizeValue:         150, // Size 0.150
				ClientOrderId:     "batch-order-002",
				TimeInForce:       1,
				ReduceOnly:        false,
				ExpireTime:        uint64(time.Now().Add(24 * time.Hour).Unix()), // Expires in 24 hours
				IsMarket:          false,
				IsPositionTp:      false,
				IsPositionSl:      false,
				TriggerType:       0,
				TriggerPriceType:  0,
				TriggerPriceValue: 0,
				IsSetOpenTp:       false,
				IsSetOpenSl:       false,
			},
		},
	}

	batchTxHash, err := client.CreateOrderBatch(&batchOrderReq)
	if err != nil {
		log.Printf("Failed to create batch orders: %v", err)
	} else {
		fmt.Printf("Batch orders created successfully, transaction hash: %s\n", batchTxHash)
	}
}

// demoTradingQueries demonstrates trading query functions
func demoTradingQueries(client *sdk.AntxClient) {
	fmt.Println("\n=== 5. Trading Query Functions Demo ===")

	// Re-fetch subaccount (if not obtained in trading functions demo)
	var testSubaccountId string
	agentAddress := client.GetAgentAddress()
	subList, err := client.GetSubaccountList(1, ethAddress, agentAddress)
	if err != nil {
		log.Printf("Failed to get subaccount list, skipping subsequent trading queries: %v", err)
	} else if len(subList) == 0 {
		log.Printf("No available subaccounts, skipping subsequent trading queries")
	} else {
		testSubaccountId = subList[0].Id
		fmt.Printf("Found subaccount: %s\n", testSubaccountId)
	}

	// Get active orders (skip if no valid subaccount)
	fmt.Println("\n5.1 Getting active orders:")
	if testSubaccountId == "" {
		fmt.Println("No valid subaccount, skipping 5.1 ~ 5.8 trading query demo")
		return
	}
	activeOrderReq := types.GetActiveOrderReq{
		SubaccountId: testSubaccountId,
		Size:         10,
	}

	activeOrderResp, err := client.GetActiveOrder(activeOrderReq)
	if err != nil {
		log.Printf("Failed to get active orders: %v", err)
	} else if len(activeOrderResp.Data.OrderList) == 0 {
		fmt.Println("No active orders currently")
	} else {
		fmt.Printf("Retrieved %d active orders:\n", len(activeOrderResp.Data.OrderList))
		for i, order := range activeOrderResp.Data.OrderList {
			if i >= 3 { // Only show first 3
				break
			}
			// Price/size are strings, display directly with decimal
			price, _ := decimal.NewFromString(order.Price)
			size, _ := decimal.NewFromString(order.Size)
			fmt.Printf("  Order %d: ID=%s, Exchange=%s, Direction=%s, Price=%s, Size=%s, Status=%d\n",
				i+1, order.Id, order.ExchangeId,
				map[bool]string{true: "Buy", false: "Sell"}[order.IsBuy],
				price.String(), size.String(), order.Status)
		}
	}

	// Get history orders
	fmt.Println("\n5.2 Getting history orders:")
	historyOrderReq := types.GetHistoryOrderReq{
		SubaccountId: testSubaccountId,
		Size:         10,
	}

	historyOrderResp, err := client.GetHistoryOrder(historyOrderReq)
	if err != nil {
		log.Printf("Failed to get history orders: %v", err)
	} else {
		fmt.Printf("Retrieved %d history orders:\n", len(historyOrderResp.Data.OrderList))
		for i, order := range historyOrderResp.Data.OrderList {
			if i >= 3 { // Only show first 3
				break
			}
			price, _ := decimal.NewFromString(order.Price)
			size, _ := decimal.NewFromString(order.Size)
			fmt.Printf("  History Order %d: ID=%s, Exchange=%s, Direction=%s, Price=%s, Size=%s, Status=%d\n",
				i+1, order.Id, order.ExchangeId,
				map[bool]string{true: "Buy", false: "Sell"}[order.IsBuy],
				price.String(), size.String(), order.Status)
		}
	}

	// Get perpetual contract account assets
	fmt.Println("\n5.3 Getting perpetual contract account assets:")
	assetReq := types.GetPerpetualAccountAssetReq{
		SubaccountId: testSubaccountId,
	}

	assetResp, err := client.GetPerpetualAccountAsset(assetReq)
	if err != nil {
		log.Printf("Failed to get perpetual contract account assets: %v", err)
	} else {
		fmt.Printf("Retrieved %d collaterals, %d positions:\n", len(assetResp.Data.CollateralList), len(assetResp.Data.PositionList))

		// Display collateral information
		for i, collateral := range assetResp.Data.CollateralList {
			if i >= 3 { // Only show first 3
				break
			}
			fmt.Printf("  Collateral %d: CoinId=%s, Amount=%s\n",
				i+1, collateral.CoinId, collateral.Amount)
		}

		// Display position information
		for i, position := range assetResp.Data.PositionList {
			if i >= 3 { // Only show first 3
				break
			}
			fmt.Printf("  Position %d: Exchange=%s, OpenSize=%s, OpenValue=%s, FundingFee=%s\n",
				i+1, position.ExchangeId, position.OpenSize, position.OpenValue, position.FundingFee)
		}
	}

	// Get position transactions
	fmt.Println("\n5.4 Getting position transactions:")
	positionReq := types.GetPositionTransactionReq{
		SubaccountId: testSubaccountId,
		Size:         10,
	}

	positionResp, err := client.GetPositionTransaction(positionReq)
	if err != nil {
		log.Printf("Failed to get position transactions: %v", err)
	} else {
		fmt.Printf("Retrieved %d position transactions:\n", len(positionResp.Data.PositionTransactionList))
		for i, position := range positionResp.Data.PositionTransactionList {
			if i >= 3 { // Only show first 3
				break
			}
			// Handle empty value display
			deltaOpenSize := position.DeltaOpenSize
			if deltaOpenSize == "" {
				deltaOpenSize = "0"
			}
			deltaOpenValue := position.DeltaOpenValue
			if deltaOpenValue == "" {
				deltaOpenValue = "0"
			}
			fillSize := position.FillSize
			if fillSize == "" {
				fillSize = "0"
			}
			fillValue := position.FillValue
			if fillValue == "" {
				fillValue = "0"
			}
			fmt.Printf("  Position Transaction %d: Exchange=%s, DeltaOpenSize=%s, DeltaOpenValue=%s, FillSize=%s, FillValue=%s\n",
				i+1, position.ExchangeId, deltaOpenSize, deltaOpenValue, fillSize, fillValue)
		}
	}

	// Get collateral transactions
	fmt.Println("\n5.5 Getting collateral transactions:")
	collateralReq := types.GetCollateralTransactionReq{
		SubaccountId: testSubaccountId,
		Size:         10,
	}

	collateralResp, err := client.GetCollateralTransaction(collateralReq)
	if err != nil {
		log.Printf("Failed to get collateral transactions: %v", err)
	} else {
		fmt.Printf("Retrieved %d collateral transactions:\n", len(collateralResp.Data.CollateralTransactionList))
		for i, collateral := range collateralResp.Data.CollateralTransactionList {
			if i >= 3 { // Only show first 3
				break
			}
			// Handle empty value display
			deltaAmount := collateral.DeltaAmount
			if deltaAmount == "" {
				deltaAmount = "0"
			}
			fmt.Printf("  Collateral Transaction %d: CoinId=%s, DeltaAmount=%s, Type=%d\n",
				i+1, collateral.CoinId, deltaAmount, collateral.Type)
		}
	}

	// Get asset snapshots
	fmt.Println("\n5.6 Getting asset snapshots:")
	snapshotReq := types.GetAssetSnapshotReq{
		SubaccountId: testSubaccountId,
		Size:         10,
	}

	snapshotResp, err := client.GetAssetSnapshot(snapshotReq)
	if err != nil {
		log.Printf("Failed to get asset snapshots: %v", err)
	} else {
		fmt.Printf("Retrieved %d asset snapshots:\n", len(snapshotResp.Data.AssetSnapshotList))
		for i, snapshot := range snapshotResp.Data.AssetSnapshotList {
			if i >= 3 { // Only show first 3
				break
			}
			fmt.Printf("  Asset Snapshot %d: CoinId=%s, TotalEquity=%s, TotalRealizePnl=%s, Time=%d\n",
				i+1, snapshot.CoinId, snapshot.TotalEquity, snapshot.TotalRealizePnl, snapshot.SnapshotTime)
		}
	}

	// Get history order fill transactions
	fmt.Println("\n5.7 Getting history order fill transactions:")
	fillReq := types.GetHistoryOrderFillTransactionReq{
		SubaccountId: testSubaccountId,
		Size:         10,
	}

	fillResp, err := client.GetHistoryOrderFillTransaction(fillReq)
	if err != nil {
		log.Printf("Failed to get history order fill transactions: %v", err)
	} else {
		fmt.Printf("Retrieved %d order fill transactions:\n", len(fillResp.Data.OrderFillTransactionList))
		for i, fill := range fillResp.Data.OrderFillTransactionList {
			if i >= 3 { // Only show first 3
				break
			}
			fmt.Printf("  Order Fill Transaction %d: %#v\n", i+1, fill)
		}
	}

	// Get history position terms
	fmt.Println("\n5.8 Getting history position terms:")
	termReq := types.GetHistoryPositionTermReq{
		SubaccountId: testSubaccountId,
		Size:         10,
	}

	termResp, err := client.GetHistoryPositionTerm(termReq)
	if err != nil {
		log.Printf("Failed to get history position terms: %v", err)
	} else {
		fmt.Printf("Retrieved %d history position terms:\n", len(termResp.Data.PositionTermList))
		for i, term := range termResp.Data.PositionTermList {
			if i >= 3 { // Only show first 3
				break
			}
			fmt.Printf("  Position Term %d: Exchange=%s, TermCount=%d, CumOpenSize=%s, CumCloseSize=%s\n",
				i+1, term.ExchangeId, term.TermCount, term.CumOpenSize, term.CumCloseSize)
		}
	}
}
