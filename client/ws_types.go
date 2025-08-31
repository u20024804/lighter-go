package client

import "time"

// WebSocket message types
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

type WSSubscribeMessage struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Symbol  string `json:"symbol,omitempty"`
}

type WSUnsubscribeMessage struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Symbol  string `json:"symbol,omitempty"`
}

// Market data types
type WSOrderBookUpdate struct {
	Symbol    string      `json:"symbol"`
	Bids      [][]string  `json:"bids"`
	Asks      [][]string  `json:"asks"`
	Timestamp int64       `json:"timestamp"`
}

type WSTickerUpdate struct {
	Symbol    string `json:"symbol"`
	Bid       string `json:"bid"`
	Ask       string `json:"ask"`
	BidSize   string `json:"bid_size"`
	AskSize   string `json:"ask_size"`
	Timestamp int64  `json:"timestamp"`
}

type WSTradeUpdate struct {
	Symbol    string `json:"symbol"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	Side      string `json:"side"`
	Timestamp int64  `json:"timestamp"`
}

// Account data types
type WSAccountUpdate struct {
	AccountIndex     int64                 `json:"account_index"`
	AvailableBalance string                `json:"available_balance"`
	MarketStats      []AccountMarketStats  `json:"market_stats,omitempty"`
	Timestamp        int64                 `json:"timestamp"`
}

type WSOrderUpdate struct {
	AccountIndex     int64  `json:"account_index"`
	OrderId          string `json:"order_id"`
	ClientOrderIndex int64  `json:"client_order_index"`
	MarketId         uint8  `json:"market_id"`
	Status           string `json:"status"`
	BaseQuantity     string `json:"base_quantity"`
	FilledQuantity   string `json:"filled_quantity"`
	Price            string `json:"price"`
	IsAsk            uint8  `json:"is_ask"`
	Timestamp        int64  `json:"timestamp"`
}

// Channel constants
const (
	ChannelOrderBook    = "orderbook"
	ChannelTicker       = "ticker"
	ChannelTrades       = "trades"
	ChannelAccount      = "account"
	ChannelOrders       = "orders"
	ChannelMarkPrice    = "markprice"
)

// Message types
const (
	MessageTypeSubscribe   = "subscribe"
	MessageTypeUnsubscribe = "unsubscribe"
	MessageTypePing        = "ping"
	MessageTypePong        = "pong"
)

// WebSocket configuration
type WSConfig struct {
	URL             string
	ReconnectDelay  time.Duration
	PingInterval    time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxReconnects   int
}

func DefaultWSConfig() *WSConfig {
	return &WSConfig{
		URL:             "wss://api.lighter.xyz/ws",
		ReconnectDelay:  5 * time.Second,
		PingInterval:    30 * time.Second,
		ReadTimeout:     60 * time.Second,
		WriteTimeout:    10 * time.Second,
		MaxReconnects:   10,
	}
}