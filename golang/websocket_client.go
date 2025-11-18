package sdk

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/antxprotocol/antx-sdk-golang/types"
	"github.com/gorilla/websocket"
)

// WsReqBase WebSocket request base structure
type WsReqBase struct {
	Method string `json:"method"` // Request method
}

// WsRegisterReq WebSocket subscription registration request structure
type WsRegisterReq struct {
	Channel      string `json:"channel"`                // Channel
	ChainType    int32  `json:"chainType,omitempty"`    // Chain type
	ChainAddress string `json:"chainAddress,omitempty"` // ETH address
}

// WsSubscribeReq WebSocket subscription request structure
type WsSubscribeReq struct {
	WsReqBase
	Subscription WsRegisterReq `json:"subscription"` // Subscription
}

// WsRespBase WebSocket response base structure
type WsRespBase struct {
	Channel string `json:"channel"`         // Channel
	Event   string `json:"event,omitempty"` // Event
	User    string `json:"user,omitempty"`  // ETH address
}

// WebSocketClient encapsulates WebSocket connection
type WebSocketClient struct {
	conn           *websocket.Conn
	url            string
	messageHandler func([]byte)
	errorHandler   func(error)
	isConnected    bool
}

// NewWebSocketClient creates a new WebSocket client
func NewWebSocketClient(wsURL string, messageHandler func([]byte), errorHandler func(error)) *WebSocketClient {
	// If a complete URL is passed, use it directly; otherwise use old logic
	var u url.URL
	if strings.HasPrefix(wsURL, "ws://") || strings.HasPrefix(wsURL, "wss://") {
		parsedURL, err := url.Parse(wsURL)
		if err != nil {
			// If parsing fails, fallback
			u = url.URL{Scheme: "ws", Host: wsURL, Path: "/api/v1/ws"}
		} else {
			u = *parsedURL
		}
	} else {
		u = url.URL{Scheme: "ws", Host: wsURL, Path: "/api/v1/ws"}
	}
	return &WebSocketClient{
		url:            u.String(),
		messageHandler: messageHandler,
		errorHandler:   errorHandler,
	}
}

// Connect establishes WebSocket connection
func (c *WebSocketClient) Connect() error {
	log.Printf("connecting to %s", c.url)

	// Set request headers to avoid WAF blocking
	header := make(http.Header)
	header.Set("X-App-Token", "ANTECH-APP-SECRET-KEY-001")
	header.Set("User-Agent", "Mozilla/5.0 (Mobile; FlutterApp/1.0)")
	header.Set("Origin", c.getOriginFromURL())

	conn, _, err := websocket.DefaultDialer.Dial(c.url, header)
	if err != nil {
		c.isConnected = false
		return fmt.Errorf("websocket dial error: %w", err)
	}
	c.conn = conn
	c.isConnected = true
	log.Println("websocket connected")

	go c.listenForMessages()
	return nil
}

// getOriginFromURL extracts Origin from WebSocket URL
func (c *WebSocketClient) getOriginFromURL() string {
	u, err := url.Parse(c.url)
	if err != nil {
		return ""
	}
	// Convert ws:// or wss:// to http:// or https://
	scheme := "http"
	if u.Scheme == "wss" {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, u.Host)
}

// listenForMessages listens for WebSocket messages
func (c *WebSocketClient) listenForMessages() {
	defer func() {
		c.isConnected = false
		if c.conn != nil {
			c.conn.Close()
		}
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if c.errorHandler != nil {
				c.errorHandler(fmt.Errorf("websocket read error: %w", err))
			}
			return
		}
		if c.messageHandler != nil {
			c.messageHandler(message)
		}
	}
}

// Subscribe subscribes to WebSocket channel
func (c *WebSocketClient) Subscribe(channel string) error {
	if !c.isConnected {
		return fmt.Errorf("websocket not connected")
	}

	req := WsSubscribeReq{
		WsReqBase: WsReqBase{
			Method: "subscribe",
		},
		Subscription: WsRegisterReq{
			Channel: channel,
		},
	}

	return c.conn.WriteJSON(req)
}

// Unsubscribe unsubscribes from WebSocket channel
func (c *WebSocketClient) Unsubscribe(channel string) error {
	if !c.isConnected {
		return fmt.Errorf("websocket not connected")
	}

	req := WsSubscribeReq{
		WsReqBase: WsReqBase{
			Method: "unsubscribe",
		},
		Subscription: WsRegisterReq{
			Channel: channel,
		},
	}

	return c.conn.WriteJSON(req)
}

// SubscribeToTicker subscribes to Ticker data
func (c *WebSocketClient) SubscribeToTicker(exchangeId string) (<-chan []byte, error) {
	channel := fmt.Sprintf("ticker.%s", exchangeId)
	err := c.Subscribe(channel)
	if err != nil {
		return nil, err
	}

	// Create a channel to receive data
	tickerChan := make(chan []byte, 100)

	// Set message handler
	originalHandler := c.messageHandler
	c.messageHandler = func(msg []byte) {
		// Parse message, check if it's ticker data
		var resp WsRespBase
		if err := json.Unmarshal(msg, &resp); err == nil {
			if resp.Channel == channel {
				select {
				case tickerChan <- msg:
				default:
					// If channel is full, drop message
				}
			}
		}

		// Call original handler
		if originalHandler != nil {
			originalHandler(msg)
		}
	}

	return tickerChan, nil
}

// SubscribeToKline subscribes to K-line data
func (c *WebSocketClient) SubscribeToKline(priceType, exchangeId, klineType string) (<-chan []byte, error) {
	channel := fmt.Sprintf("kline.%s.%s.%s", priceType, exchangeId, klineType)
	err := c.Subscribe(channel)
	if err != nil {
		return nil, err
	}

	// Create a channel to receive data
	klineChan := make(chan []byte, 100)

	// Set message handler
	originalHandler := c.messageHandler
	c.messageHandler = func(msg []byte) {
		// Parse message, check if it's kline data
		var resp WsRespBase
		if err := json.Unmarshal(msg, &resp); err == nil {
			if resp.Channel == channel {
				select {
				case klineChan <- msg:
				default:
					// If channel is full, drop message
				}
			}
		}

		// Call original handler
		if originalHandler != nil {
			originalHandler(msg)
		}
	}

	return klineChan, nil
}

// Disconnect disconnects WebSocket connection
func (c *WebSocketClient) Disconnect() error {
	if c.conn != nil {
		c.isConnected = false
		return c.conn.Close()
	}
	return nil
}

// IsConnected checks connection status
func (c *WebSocketClient) IsConnected() bool {
	return c.isConnected
}

// ParseTickerData parses Ticker data
func ParseTickerData(data []byte) (*types.TickerData, error) {
	// WebSocket returns wrapped data structure, need to parse outer structure first
	var wsResponse struct {
		Channel string             `json:"channel"`
		Event   string             `json:"event"`
		Data    []types.TickerData `json:"data"`
	}

	if err := json.Unmarshal(data, &wsResponse); err != nil {
		return nil, fmt.Errorf("failed to parse websocket response: %w", err)
	}

	// Check if there is data
	if len(wsResponse.Data) == 0 {
		return nil, fmt.Errorf("no ticker data in response")
	}

	// Return first ticker data
	return &wsResponse.Data[0], nil
}

// ParseKlineData parses K-line data
func ParseKlineData(data []byte) (*types.KLine, error) {
	// WebSocket returns wrapped data structure, need to parse outer structure first
	var wsResponse struct {
		Channel string        `json:"channel"`
		Event   string        `json:"event"`
		Data    []types.KLine `json:"data"`
	}

	if err := json.Unmarshal(data, &wsResponse); err != nil {
		return nil, fmt.Errorf("failed to parse websocket response: %w", err)
	}

	// Check if there is data
	if len(wsResponse.Data) == 0 {
		return nil, fmt.Errorf("no kline data in response")
	}

	// Return first kline data
	return &wsResponse.Data[0], nil
}
