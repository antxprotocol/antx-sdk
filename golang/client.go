package sdk

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/antxprotocol/antx-sdk-golang/constants"
	"github.com/antxprotocol/antx-sdk-golang/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	globalClient *AntxClient
)

// Config client configuration
type Config struct {
	GatewayHost     string // Gateway URI, e.g., "http://127.0.0.1:8080"
	ChainID         string // Chain ID, e.g., "antx-devnet"
	EthPrivateKey   string // Private key in hexadecimal string
	AgentPrivateKey string // Private key in hexadecimal string
}

// AntxClient encapsulates the client for interacting with Antx chain
type AntxClient struct {
	clientCtx       client.Context
	ethPrivateKey   *ecdsa.PrivateKey
	ethAddress      ethCommon.Address
	agentPrivateKey cryptotypes.PrivKey
	agentAddress    sdk.AccAddress
	chainID         string
	gatewayHost     string
	accountNumber   uint64
	// merged HTTP/WebSocket capabilities
	baseURL    string
	wsURL      string
	httpClient *http.Client
	wsClient   *WebSocketClient
}

// NewAntxClient creates a new Antx client
func NewAntxClient(config Config) (*AntxClient, error) {
	// Validate configuration parameters
	if config.ChainID == "" {
		return nil, fmt.Errorf("chain ID cannot be empty")
	}
	if config.EthPrivateKey == "" {
		return nil, fmt.Errorf("eth private key cannot be empty")
	}
	if config.AgentPrivateKey == "" {
		return nil, fmt.Errorf("agent private key cannot be empty")
	}

	// Parse private keys
	ethPrivateKeyHex := strings.TrimPrefix(config.EthPrivateKey, "0x")
	agentPrivateKeyHex := strings.TrimPrefix(config.AgentPrivateKey, "0x")
	if len(ethPrivateKeyHex) != 64 {
		return nil, fmt.Errorf("invalid eth private key length: expected 64 characters, got %d", len(ethPrivateKeyHex))
	}
	if len(agentPrivateKeyHex) != 64 {
		return nil, fmt.Errorf("invalid agent private key length: expected 64 characters, got %d", len(agentPrivateKeyHex))
	}

	agentPrivateKeyBytes, err := hex.DecodeString(agentPrivateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode agent private key: %w", err)
	}

	// Create secp256k1 private key
	ethPrivatekeyHex := strings.TrimPrefix(config.EthPrivateKey, "0x")
	ethPrivateKey, err := ethCrypto.HexToECDSA(ethPrivatekeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode eth private key: %w", err)
	}
	agentPrivateKey := &secp256k1.PrivKey{Key: agentPrivateKeyBytes}

	// Get addresses
	ethAddress := ethCrypto.PubkeyToAddress(ethPrivateKey.PublicKey)
	agentAddress := sdk.AccAddress(agentPrivateKey.PubKey().Address())

	// Create interface registry
	interfaceRegistry := codectypes.NewInterfaceRegistry()

	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)

	// Create codec
	cdc := codec.NewProtoCodec(interfaceRegistry)

	// Create client context
	clientCtx := client.Context{}.
		WithCodec(cdc).
		WithInterfaceRegistry(interfaceRegistry).
		WithBroadcastMode(flags.BroadcastSync).
		WithChainID(config.ChainID).
		WithFromAddress(agentAddress).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithTxConfig(authtx.NewTxConfig(cdc, authtx.DefaultSignModes))

	client := &AntxClient{
		clientCtx:       clientCtx,
		ethPrivateKey:   ethPrivateKey,
		ethAddress:      ethAddress,
		agentPrivateKey: agentPrivateKey,
		agentAddress:    agentAddress,
		chainID:         config.ChainID,
		gatewayHost:     config.GatewayHost,
	}

	// initialize http client and baseURL
	client.httpClient = &http.Client{Timeout: 30 * time.Second}
	client.baseURL = config.GatewayHost

	if config.GatewayHost != "" {
		accountNumber, _, err := client.GetAccountNumberAndSequence(agentAddress.String())
		if err != nil {
			return nil, fmt.Errorf("failed to get account number and sequence: %w", err)
		}
		client.accountNumber, err = strconv.ParseUint(accountNumber, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse account number: %w", err)
		}
	}
	return client, nil
}

// NewAntxQueryClient creates a lightweight client for HTTP queries and WebSocket only (no on-chain signing configuration required)
func NewAntxQueryClient(baseURL, wsURL string) *AntxClient {
	return &AntxClient{
		baseURL:    baseURL,
		wsURL:      wsURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// GetAgentAddress gets the agent address
func (c *AntxClient) GetAgentAddress() string {
	return c.agentAddress.String()
}

// SetGateway sets the HTTP and WebSocket gateway addresses
func (c *AntxClient) SetGateway(baseURL, wsURL string) {
	c.baseURL = baseURL
	c.wsURL = wsURL
	if c.httpClient == nil {
		c.httpClient = &http.Client{Timeout: 30 * time.Second}
	}
}

// =============================== HTTP Request Methods (merged) ===============================

func (c *AntxClient) httpGet(path string, params map[string]string, result interface{}) error {
	if c.baseURL == "" {
		return fmt.Errorf("gateway baseURL is not set")
	}
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	if c.httpClient == nil {
		c.httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create GET request: %w", err)
	}
	// Set request headers to avoid WAF blocking
	req.Header.Set("X-App-Token", "ANTECH-APP-SECRET-KEY-001")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Mobile; FlutterApp/1.0)")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w, body: %s", err, string(body))
	}
	return nil
}

func (c *AntxClient) httpPost(path string, data interface{}, result interface{}) error {
	if c.baseURL == "" {
		return fmt.Errorf("gateway baseURL is not set")
	}
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request data: %w", err)
	}
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed to create POST request: %w", err)
	}
	// Set request headers to avoid WAF blocking
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Token", "ANTECH-APP-SECRET-KEY-001")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Mobile; FlutterApp/1.0)")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w, body: %s", err, string(body))
	}
	return nil
}

// GetAccountNumberAndSequence gets the account number and sequence
func (c *AntxClient) GetAccountNumberAndSequence(address string) (string, string, error) {
	if c.baseURL == "" {
		return "0", "0", nil
	}

	var result types.GetAccountNumberAndSequenceResponse
	params := map[string]string{
		"address": address,
	}
	if err := c.httpGet(constants.GetAddressInfoPath, params, &result); err != nil {
		return "", "", err
	}

	if result.BaseResp.Code != "0" {
		return "", "", fmt.Errorf("get account info failed: %s", result.BaseResp.Msg)
	}

	return result.Data.AccountNumber, result.Data.Sequence, nil
}

// SendRawTx sends a raw transaction
func (c *AntxClient) SendRawTx(req types.SendRawTxRequest) (*types.SendRawTxResponse, error) {
	if c.baseURL == "" {
		return &types.SendRawTxResponse{
			BaseResp: types.BaseResp{Code: "0", Msg: "success"},
			Data: types.SendRawTxResponseData{
				TxHash:     "mock_tx_hash",
				RawTx:      req.RawTx,
				ResultData: "mock_result",
			},
		}, nil
	}

	var result types.SendRawTxResponse
	if err := c.httpPost(constants.SendTransactionPath, req, &result); err != nil {
		return nil, err
	}

	// Add debug information
	if result.Data.TxHash != "" {
		logx.Infof("SendRawTx response: txHash=%s", result.Data.TxHash)
	}

	return &result, nil
}

// GetEthAddress gets the Ethereum address
func (c *AntxClient) GetEthAddress() string {
	addr, _ := ConvertToEthAddr(c.ethAddress.String())
	return addr
}

func (c *AntxClient) SignAndSendTx(typeURL string, msg sdk.Msg, unordered bool) (string, error) {
	return c.signAndSendTx(typeURL, msg, unordered)
}

func (c *AntxClient) signAndSendTx(typeURL string, msg sdk.Msg, unordered bool) (string, error) {
	// Create transaction builder
	txBuilder := c.clientCtx.TxConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msg); err != nil {
		logx.Errorf("failed to set messages: %w", err)
		return "", fmt.Errorf("failed to set messages: %w", err)
	}
	timeoutInt := time.Now().Add(10 * time.Second).UnixNano()
	timeout := time.Unix(timeoutInt/1e9, timeoutInt%1e9)
	if unordered {
		txBuilder.SetUnordered(unordered)
		txBuilder.SetTimeoutTimestamp(timeout)
	}

	// Set gas and fee
	txBuilder.SetGasLimit(200000)
	txBuilder.SetFeeAmount(sdk.NewCoins()) // No fee

	// Create in-memory keyring for signing
	kr := keyring.NewInMemory(c.clientCtx.Codec)

	// Import private key directly to keyring
	keyName := "temp-key"
	privKeyHex := hex.EncodeToString(c.agentPrivateKey.Bytes())
	err := kr.ImportPrivKeyHex(keyName, privKeyHex, "secp256k1")
	if err != nil {
		logx.Errorf("failed to import private key to keyring: %w", err)
		return "", fmt.Errorf("failed to import private key to keyring: %w", err)
	}

	// Create transaction factory
	txFactory := tx.Factory{}.
		WithChainID(c.chainID).
		WithTxConfig(c.clientCtx.TxConfig).
		WithAccountNumber(c.accountNumber).
		WithSignMode(authtx.DefaultSignModes[0]).
		WithKeybase(kr)

	if !unordered {
		_, sequence, err := c.GetAccountNumberAndSequence(c.agentAddress.String())
		if err != nil {
			logx.Errorf("failed to get account number and sequence: %w", err)
			return "", fmt.Errorf("failed to get account number and sequence: %w", err)
		}
		sequenceUint, err := strconv.ParseUint(sequence, 10, 64)
		if err != nil {
			logx.Errorf("failed to parse sequence: %w", err)
			return "", fmt.Errorf("failed to parse sequence: %w", err)
		}
		txFactory = txFactory.WithSequence(sequenceUint)
	}

	// Sign transaction using tx.Sign
	if err := tx.Sign(context.Background(), txFactory, keyName, txBuilder, true); err != nil {
		logx.Errorf("failed to sign transaction: %w, ttl: %v", err, timeout.Format(time.RFC3339))
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	txBytes, err := c.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		logx.Errorf("failed to encode transaction: %w, ttl: %v", err, timeout.Format(time.RFC3339))
		return "", fmt.Errorf("failed to encode transaction: %w, ttl: %v", err, timeout.Format(time.RFC3339))
	}
	logx.Infof("rawTx: %s", base64.StdEncoding.EncodeToString(txBytes))

	// Send transaction
	req := types.SendRawTxRequest{
		TypeURL:       typeURL,
		RawTx:         base64.StdEncoding.EncodeToString(txBytes),
		AccountNumber: c.accountNumber,
	}
	resp, err := c.SendRawTx(req)
	if err != nil {
		logx.Errorf("failed to send transaction: %w, ttl: %v", err, timeout.Format(time.RFC3339))
		return "", fmt.Errorf("failed to send transaction: %w, ttl: %v", err, timeout.Format(time.RFC3339))
	}
	// Try to get transaction hash, support multiple field names
	txHash := resp.Data.TxHash
	if txHash == "" {
		txHash = resp.Data.Hash
	}
	if txHash == "" {
		txHash = resp.Data.TxID
	}

	return txHash, nil
}

// =============================== Market Data and Trading Queries (merged from SDKClient) ===============================

// GetCoinList gets the coin list
func (c *AntxClient) GetCoinList() ([]types.Coin, error) {
	var result types.GetCoinListResponse
	if err := c.httpGet(constants.GetCoinListPath, map[string]string{}, &result); err != nil {
		return nil, err
	}
	if result.BaseResp.Code != "0" {
		return nil, fmt.Errorf("get coin list failed: %s", result.BaseResp.Msg)
	}
	return result.Data.CoinList, nil
}

// GetSubaccountList gets the subaccount list
func (c *AntxClient) GetSubaccountList(chainType int32, chainAddress, agentAddress string) ([]types.Subaccount, error) {
	var result types.GetSubaccountListResponse
	params := map[string]string{
		"chainType":    strconv.FormatInt(int64(chainType), 10),
		"chainAddress": chainAddress,
		"agentAddress": agentAddress,
	}
	if err := c.httpGet(constants.GetSubaccountPath, params, &result); err != nil {
		return nil, err
	}
	if result.BaseResp.Code != "0" {
		return nil, fmt.Errorf("get subaccount list failed: %s", result.BaseResp.Msg)
	}
	return result.Data.SubaccountList, nil
}

// GetExchangeList gets the exchange list
func (c *AntxClient) GetExchangeList() ([]types.Exchange, error) {
	var result types.GetExchangeListResponse
	if err := c.httpGet(constants.GetExchangeListPath, map[string]string{}, &result); err != nil {
		return nil, err
	}
	if result.BaseResp.Code != "0" {
		return nil, fmt.Errorf("get exchange list failed: %s", result.BaseResp.Msg)
	}
	return result.Data.ExchangeList, nil
}

// GetKline gets K-line data
func (c *AntxClient) GetKline(req types.GetKLineReq) (*types.GetKLineResp, error) {
	var result types.GetKLineResp
	params := map[string]string{
		"exchangeId": req.ExchangeId,
		"klineType":  req.KlineType,
		"priceType":  req.PriceType,
	}
	if req.Size > 0 {
		params["size"] = strconv.FormatUint(uint64(req.Size), 10)
	}
	if req.OffsetData != "" {
		params["offsetData"] = req.OffsetData
	}
	if req.FilterBeginKlineTimeInclusive > 0 {
		params["filterBeginKlineTimeInclusive"] = strconv.FormatInt(req.FilterBeginKlineTimeInclusive, 10)
	}
	if req.FilterEndKlineTimeExclusive > 0 {
		params["filterEndKlineTimeExclusive"] = strconv.FormatInt(req.FilterEndKlineTimeExclusive, 10)
	}
	if err := c.httpGet(constants.GetKlinePath, params, &result); err != nil {
		return nil, err
	}
	if result.BaseResp.Code != "0" {
		return nil, fmt.Errorf("get kline failed: %s", result.BaseResp.Msg)
	}
	return &result, nil
}

// GetFundingHistory gets funding rate history
func (c *AntxClient) GetFundingHistory(req types.GetFundingHistoryReq) (*types.GetFundingHistoryResp, error) {
	var result types.GetFundingHistoryResp
	params := map[string]string{
		"exchangeId": req.ExchangeId,
		"size":       strconv.FormatUint(uint64(req.Size), 10),
	}
	if req.OffsetData != "" {
		params["offsetData"] = req.OffsetData
	}
	if req.FilterSettlementFundingRate {
		params["filterSettlementFundingRate"] = "true"
	}
	if req.FilterBeginTimeInclusive > 0 {
		params["filterBeginTimeInclusive"] = strconv.FormatUint(req.FilterBeginTimeInclusive, 10)
	}
	if req.FilterEndTimeExclusive > 0 {
		params["filterEndTimeExclusive"] = strconv.FormatUint(req.FilterEndTimeExclusive, 10)
	}
	if err := c.httpGet(constants.GetFundingHistoryPath, params, &result); err != nil {
		return nil, err
	}
	if result.BaseResp.Code != "0" {
		return nil, fmt.Errorf("get funding history failed: %s", result.BaseResp.Msg)
	}
	return &result, nil
}

// GetActiveOrder gets active orders
func (c *AntxClient) GetActiveOrder(req types.GetActiveOrderReq) (*types.GetActiveOrderResp, error) {
	var result types.GetActiveOrderResp
	params := map[string]string{
		"subaccountId": req.SubaccountId,
		"size":         strconv.FormatUint(uint64(req.Size), 10),
	}
	if req.OffsetData != "" {
		params["offsetData"] = req.OffsetData
	}
	if req.PageOffsetDataCreatedTime != "" {
		params["pageOffsetDataCreatedTime"] = req.PageOffsetDataCreatedTime
	}
	if req.PageOffsetDataItemId != "" {
		params["pageOffsetDataItemId"] = req.PageOffsetDataItemId
	}
	if req.FilterExchangeIdList != "" {
		params["filterExchangeIdList"] = req.FilterExchangeIdList
	}
	if req.FilterOrderStatusList != "" {
		params["filterOrderStatusList"] = req.FilterOrderStatusList
	}
	if req.FilterIsLiquidateList != "" {
		params["filterIsLiquidateList"] = req.FilterIsLiquidateList
	}
	if req.FilterIsDeleverageList != "" {
		params["filterIsDeleverageList"] = req.FilterIsDeleverageList
	}
	if req.FilterIsPositionTpslList != "" {
		params["filterIsPositionTpslList"] = req.FilterIsPositionTpslList
	}
	if req.FilterOrderIdList != "" {
		params["filterOrderIdList"] = req.FilterOrderIdList
	}
	if req.FilterStartCreatedTimeInclusive > 0 {
		params["filterStartCreatedTimeInclusive"] = strconv.FormatUint(req.FilterStartCreatedTimeInclusive, 10)
	}
	if req.FilterEndCreatedTimeExclusive > 0 {
		params["filterEndCreatedTimeExclusive"] = strconv.FormatUint(req.FilterEndCreatedTimeExclusive, 10)
	}
	// Add debug information
	logx.Infof("GetActiveOrder request params: %+v", params)

	if err := c.httpGet(constants.GetActiveOrderPath, params, &result); err != nil {
		return nil, err
	}
	if result.BaseResp.Code != "0" {
		return nil, fmt.Errorf("get active order failed: %s", result.BaseResp.Msg)
	}
	return &result, nil
}

// GetHistoryOrder gets history orders
func (c *AntxClient) GetHistoryOrder(req types.GetHistoryOrderReq) (*types.GetHistoryOrderResp, error) {
	var result types.GetHistoryOrderResp
	params := map[string]string{
		"subaccountId": req.SubaccountId,
		"size":         strconv.FormatUint(uint64(req.Size), 10),
	}
	if req.OffsetData != "" {
		params["offsetData"] = req.OffsetData
	}
	if req.PageOffsetDataCreatedTime != "" {
		params["pageOffsetDataCreatedTime"] = req.PageOffsetDataCreatedTime
	}
	if req.PageOffsetDataItemId != "" {
		params["pageOffsetDataItemId"] = req.PageOffsetDataItemId
	}
	if req.FilterExchangeIdList != "" {
		params["filterExchangeIdList"] = req.FilterExchangeIdList
	}
	if req.FilterOrderStatusList != "" {
		params["filterOrderStatusList"] = req.FilterOrderStatusList
	}
	if req.FilterIsLiquidateList != "" {
		params["filterIsLiquidateList"] = req.FilterIsLiquidateList
	}
	if req.FilterIsDeleverageList != "" {
		params["filterIsDeleverageList"] = req.FilterIsDeleverageList
	}
	if req.FilterIsPositionTpslList != "" {
		params["filterIsPositionTpslList"] = req.FilterIsPositionTpslList
	}
	if req.FilterOrderIdList != "" {
		params["filterOrderIdList"] = req.FilterOrderIdList
	}
	if req.FilterStartCreatedTimeInclusive > 0 {
		params["filterStartCreatedTimeInclusive"] = strconv.FormatUint(req.FilterStartCreatedTimeInclusive, 10)
	}
	if req.FilterEndCreatedTimeExclusive > 0 {
		params["filterEndCreatedTimeExclusive"] = strconv.FormatUint(req.FilterEndCreatedTimeExclusive, 10)
	}
	if err := c.httpGet(constants.GetHistoryOrderPath, params, &result); err != nil {
		return nil, err
	}
	if result.BaseResp.Code != "0" {
		return nil, fmt.Errorf("get history order failed: %s", result.BaseResp.Msg)
	}
	return &result, nil
}

// GetPerpetualAccountAsset gets perpetual contract account assets
func (c *AntxClient) GetPerpetualAccountAsset(req types.GetPerpetualAccountAssetReq) (*types.GetPerpetualAccountAssetResp, error) {
	var result types.GetPerpetualAccountAssetResp
	params := map[string]string{"subaccountId": req.SubaccountId}
	if err := c.httpGet(constants.GetPerpetualAccountAssetPath, params, &result); err != nil {
		return nil, err
	}
	if result.BaseResp.Code != "0" {
		return nil, fmt.Errorf("get perpetual account asset failed: %s", result.BaseResp.Msg)
	}
	return &result, nil
}

// GetPositionTransaction gets position transactions
func (c *AntxClient) GetPositionTransaction(req types.GetPositionTransactionReq) (*types.GetPositionTransactionResp, error) {
	var result types.GetPositionTransactionResp
	params := map[string]string{
		"subaccountId": req.SubaccountId,
		"size":         strconv.FormatUint(uint64(req.Size), 10),
	}
	if req.PageOffsetDataCreatedTime != "" {
		params["pageOffsetDataCreatedTime"] = req.PageOffsetDataCreatedTime
	}
	if req.PageOffsetDataItemId != "" {
		params["pageOffsetDataItemId"] = req.PageOffsetDataItemId
	}
	if req.FilterExchangeIdList != "" {
		params["filterExchangeIdList"] = req.FilterExchangeIdList
	}
	if req.FilterTypeList != "" {
		params["filterTypeList"] = req.FilterTypeList
	}
	if req.FilterMarginModeList != "" {
		params["filterMarginModeList"] = req.FilterMarginModeList
	}
	if req.FilterStartCreatedTimeInclusive > 0 {
		params["filterStartCreatedTimeInclusive"] = strconv.FormatUint(req.FilterStartCreatedTimeInclusive, 10)
	}
	if req.FilterEndCreatedTimeExclusive > 0 {
		params["filterEndCreatedTimeExclusive"] = strconv.FormatUint(req.FilterEndCreatedTimeExclusive, 10)
	}
	if err := c.httpGet(constants.GetPositionTransactionPath, params, &result); err != nil {
		return nil, err
	}
	if result.BaseResp.Code != "0" {
		return nil, fmt.Errorf("get position transaction failed: %s", result.BaseResp.Msg)
	}
	return &result, nil
}

// GetCollateralTransaction gets collateral transactions
func (c *AntxClient) GetCollateralTransaction(req types.GetCollateralTransactionReq) (*types.GetCollateralTransactionResp, error) {
	var result types.GetCollateralTransactionResp
	params := map[string]string{
		"subaccountId": req.SubaccountId,
		"size":         strconv.FormatUint(uint64(req.Size), 10),
	}
	if req.PageOffsetDataCreatedTime != "" {
		params["pageOffsetDataCreatedTime"] = req.PageOffsetDataCreatedTime
	}
	if req.PageOffsetDataItemId != "" {
		params["pageOffsetDataItemId"] = req.PageOffsetDataItemId
	}
	if req.FilterCoinId != "" {
		params["filterCoinId"] = req.FilterCoinId
	}
	if req.FilterTypeList != "" {
		params["filterTypeList"] = req.FilterTypeList
	}
	if req.FilterStartCreatedTimeInclusive > 0 {
		params["filterStartCreatedTimeInclusive"] = strconv.FormatUint(req.FilterStartCreatedTimeInclusive, 10)
	}
	if req.FilterEndCreatedTimeExclusive > 0 {
		params["filterEndCreatedTimeExclusive"] = strconv.FormatUint(req.FilterEndCreatedTimeExclusive, 10)
	}
	if err := c.httpGet(constants.GetCollateralTransactionPath, params, &result); err != nil {
		return nil, err
	}
	if result.BaseResp.Code != "0" {
		return nil, fmt.Errorf("get collateral transaction failed: %s", result.BaseResp.Msg)
	}
	return &result, nil
}

// GetAssetSnapshot gets asset snapshots
func (c *AntxClient) GetAssetSnapshot(req types.GetAssetSnapshotReq) (*types.GetAssetSnapshotResp, error) {
	var result types.GetAssetSnapshotResp
	params := map[string]string{
		"subaccountId": req.SubaccountId,
		"size":         strconv.FormatUint(uint64(req.Size), 10),
	}
	if req.PageOffsetDataCreatedTime != "" {
		params["pageOffsetDataCreatedTime"] = req.PageOffsetDataCreatedTime
	}
	if req.PageOffsetDataItemId != "" {
		params["pageOffsetDataItemId"] = req.PageOffsetDataItemId
	}
	if req.FilterCoinId != "" {
		params["filterCoinId"] = req.FilterCoinId
	}
	if req.FilterTimeTag != "" {
		params["filterTimeTag"] = req.FilterTimeTag
	}
	if req.FilterStartCreatedTimeInclusive > 0 {
		params["filterStartCreatedTimeInclusive"] = strconv.FormatUint(req.FilterStartCreatedTimeInclusive, 10)
	}
	if req.FilterEndCreatedTimeExclusive > 0 {
		params["filterEndCreatedTimeExclusive"] = strconv.FormatUint(req.FilterEndCreatedTimeExclusive, 10)
	}
	if err := c.httpGet(constants.GetAssetSnapshotPath, params, &result); err != nil {
		return nil, err
	}
	if result.BaseResp.Code != "0" {
		return nil, fmt.Errorf("get asset snapshot failed: %s", result.BaseResp.Msg)
	}
	return &result, nil
}

// GetHistoryOrderFillTransaction gets history order fill transactions
func (c *AntxClient) GetHistoryOrderFillTransaction(req types.GetHistoryOrderFillTransactionReq) (*types.GetHistoryOrderFillTransactionResp, error) {
	var result types.GetHistoryOrderFillTransactionResp
	params := map[string]string{
		"subaccountId": req.SubaccountId,
		"size":         strconv.FormatUint(uint64(req.Size), 10),
	}
	if req.PageOffsetDataCreatedTime != "" {
		params["pageOffsetDataCreatedTime"] = req.PageOffsetDataCreatedTime
	}
	if req.PageOffsetDataItemId != "" {
		params["pageOffsetDataItemId"] = req.PageOffsetDataItemId
	}
	if req.FilterExchangeIdList != "" {
		params["filterExchangeIdList"] = req.FilterExchangeIdList
	}
	if req.FilterCoinIdList != "" {
		params["filterCoinIdList"] = req.FilterCoinIdList
	}
	if req.FilterOrderIdList != "" {
		params["filterOrderIdList"] = req.FilterOrderIdList
	}
	if req.FilterStartCreatedTimeInclusive > 0 {
		params["filterStartCreatedTimeInclusive"] = strconv.FormatUint(req.FilterStartCreatedTimeInclusive, 10)
	}
	if req.FilterEndCreatedTimeExclusive > 0 {
		params["filterEndCreatedTimeExclusive"] = strconv.FormatUint(req.FilterEndCreatedTimeExclusive, 10)
	}
	if err := c.httpGet(constants.GetHistoryOrderFillTransactionPath, params, &result); err != nil {
		return nil, err
	}
	if result.BaseResp.Code != "0" {
		return nil, fmt.Errorf("get history order fill transaction failed: %s", result.BaseResp.Msg)
	}
	return &result, nil
}

// GetHistoryPositionTerm gets history position terms
func (c *AntxClient) GetHistoryPositionTerm(req types.GetHistoryPositionTermReq) (*types.GetHistoryPositionTermResp, error) {
	var result types.GetHistoryPositionTermResp
	params := map[string]string{
		"subaccountId": req.SubaccountId,
		"size":         strconv.FormatUint(uint64(req.Size), 10),
	}
	if req.PageOffsetDataCreatedTime != "" {
		params["pageOffsetDataCreatedTime"] = req.PageOffsetDataCreatedTime
	}
	if req.PageOffsetDataItemId != "" {
		params["pageOffsetDataItemId"] = req.PageOffsetDataItemId
	}
	if req.FilterExchangeIdList != "" {
		params["filterExchangeIdList"] = req.FilterExchangeIdList
	}
	if req.FilterStartCreatedTimeInclusive > 0 {
		params["filterStartCreatedTimeInclusive"] = strconv.FormatUint(req.FilterStartCreatedTimeInclusive, 10)
	}
	if req.FilterEndCreatedTimeExclusive > 0 {
		params["filterEndCreatedTimeExclusive"] = strconv.FormatUint(req.FilterEndCreatedTimeExclusive, 10)
	}
	if err := c.httpGet(constants.GetHistoryPositionTermPath, params, &result); err != nil {
		return nil, err
	}
	if result.BaseResp.Code != "0" {
		return nil, fmt.Errorf("get history position term failed: %s", result.BaseResp.Msg)
	}
	return &result, nil
}

// =============================== WebSocket Integration and Parsing ===============================

// ConnectWebSocket establishes connection
func (c *AntxClient) ConnectWebSocket(messageHandler func([]byte), errorHandler func(error)) error {
	if c.wsClient != nil {
		_ = c.wsClient.Disconnect()
	}
	if c.wsURL == "" {
		return fmt.Errorf("wsURL is not set")
	}
	c.wsClient = NewWebSocketClient(c.wsURL, messageHandler, errorHandler)
	return c.wsClient.Connect()
}

// SubscribeToTicker subscribes to Ticker
func (c *AntxClient) SubscribeToTicker(exchangeId string) (<-chan []byte, error) {
	if c.wsClient == nil {
		return nil, fmt.Errorf("websocket not connected")
	}
	return c.wsClient.SubscribeToTicker(exchangeId)
}

// SubscribeToKline subscribes to K-line
func (c *AntxClient) SubscribeToKline(priceType, exchangeId, klineType string) (<-chan []byte, error) {
	if c.wsClient == nil {
		return nil, fmt.Errorf("websocket not connected")
	}
	return c.wsClient.SubscribeToKline(priceType, exchangeId, klineType)
}

// DisconnectWebSocket disconnects
func (c *AntxClient) DisconnectWebSocket() error {
	if c.wsClient != nil {
		return c.wsClient.Disconnect()
	}
	return nil
}

// ParseTickerData parses Ticker
func (c *AntxClient) ParseTickerData(data []byte) (*types.TickerData, error) {
	// WebSocket push format: {"channel":"...","event":"payload","data":[{...}]}
	var wsResp struct {
		Channel string             `json:"channel"`
		Event   string             `json:"event"`
		Data    []types.TickerData `json:"data"`
	}
	if err := json.Unmarshal(data, &wsResp); err != nil {
		return nil, fmt.Errorf("failed to parse websocket response: %w", err)
	}
	if len(wsResp.Data) == 0 {
		return nil, fmt.Errorf("no ticker data in response")
	}
	return &wsResp.Data[0], nil
}

// ParseKlineData parses K-line
func (c *AntxClient) ParseKlineData(data []byte) (*types.KLine, error) {
	// WebSocket push format: {"channel":"...","event":"payload","data":[{...}]}
	var wsResp struct {
		Channel string        `json:"channel"`
		Event   string        `json:"event"`
		Data    []types.KLine `json:"data"`
	}
	if err := json.Unmarshal(data, &wsResp); err != nil {
		return nil, fmt.Errorf("failed to parse websocket response: %w", err)
	}
	if len(wsResp.Data) == 0 {
		return nil, fmt.Errorf("no kline data in response")
	}
	return &wsResp.Data[0], nil
}

// ParseDepthData parses depth
func (c *AntxClient) ParseDepthData(data []byte) (*types.DepthData, error) {
	// WebSocket push format: {"channel":"...","event":"payload","data":[{...}]}
	var wsResp struct {
		Channel string            `json:"channel"`
		Event   string            `json:"event"`
		Data    []types.DepthData `json:"data"`
	}
	if err := json.Unmarshal(data, &wsResp); err != nil {
		return nil, fmt.Errorf("failed to parse websocket response: %w", err)
	}
	if len(wsResp.Data) == 0 {
		return nil, fmt.Errorf("no depth data in response")
	}
	return &wsResp.Data[0], nil
}

// ParseTradeData parses trade
func (c *AntxClient) ParseTradeData(data []byte) (*types.Ticket, error) {
	// WebSocket push format: {"channel":"...","event":"payload","data":[{...}]}
	var wsResp struct {
		Channel string         `json:"channel"`
		Event   string         `json:"event"`
		Data    []types.Ticket `json:"data"`
	}
	if err := json.Unmarshal(data, &wsResp); err != nil {
		return nil, fmt.Errorf("failed to parse websocket response: %w", err)
	}
	if len(wsResp.Data) == 0 {
		return nil, fmt.Errorf("no trade data in response")
	}
	return &wsResp.Data[0], nil
}

// ParseFundingRateData parses funding rate
func (c *AntxClient) ParseFundingRateData(data []byte) (*types.FundingRate, error) {
	// WebSocket push format: {"channel":"...","event":"payload","data":[{...}]}
	var wsResp struct {
		Channel string              `json:"channel"`
		Event   string              `json:"event"`
		Data    []types.FundingRate `json:"data"`
	}
	if err := json.Unmarshal(data, &wsResp); err != nil {
		return nil, fmt.Errorf("failed to parse websocket response: %w", err)
	}
	if len(wsResp.Data) == 0 {
		return nil, fmt.Errorf("no funding rate data in response")
	}
	return &wsResp.Data[0], nil
}

// ParsePriceData parses price
func (c *AntxClient) ParsePriceData(data []byte) (*types.Price, error) {
	// WebSocket push format: {"channel":"...","event":"payload","data":[{...}]}
	var wsResp struct {
		Channel string        `json:"channel"`
		Event   string        `json:"event"`
		Data    []types.Price `json:"data"`
	}
	if err := json.Unmarshal(data, &wsResp); err != nil {
		return nil, fmt.Errorf("failed to parse websocket response: %w", err)
	}
	if len(wsResp.Data) == 0 {
		return nil, fmt.Errorf("no price data in response")
	}
	return &wsResp.Data[0], nil
}
