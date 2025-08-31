package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// Market Data Streams

// StreamOrderBook subscribes to order book updates for a market ID
func (ws *WSClient) StreamOrderBook(ctx context.Context, marketId uint8, callback func(*WSOrderBookUpdate) error) error {
	handler := func(data []byte) error {
		var response struct {
			Type string             `json:"type"`
			Data *WSOrderBookUpdate `json:"data"`
		}
		
		if err := json.Unmarshal(data, &response); err != nil {
			return fmt.Errorf("failed to unmarshal order book update: %v", err)
		}
		
		if response.Type != ChannelOrderBook || response.Data == nil {
			return nil // Not an order book message
		}
		
		return callback(response.Data)
	}
	
	ws.AddHandler(ChannelOrderBook, handler)
	return ws.Subscribe(ChannelOrderBook, fmt.Sprintf("%d", marketId))
}

// StreamTicker subscribes to ticker updates for a market ID
func (ws *WSClient) StreamTicker(ctx context.Context, marketId uint8, callback func(*WSTickerUpdate) error) error {
	handler := func(data []byte) error {
		var response struct {
			Type string           `json:"type"`
			Data *WSTickerUpdate  `json:"data"`
		}
		
		if err := json.Unmarshal(data, &response); err != nil {
			return fmt.Errorf("failed to unmarshal ticker update: %v", err)
		}
		
		if response.Type != ChannelTicker || response.Data == nil {
			return nil
		}
		
		return callback(response.Data)
	}
	
	ws.AddHandler(ChannelTicker, handler)
	return ws.Subscribe(ChannelTicker, fmt.Sprintf("%d", marketId))
}

// StreamTrades subscribes to trade updates for a market ID
func (ws *WSClient) StreamTrades(ctx context.Context, marketId uint8, callback func(*WSTradeUpdate) error) error {
	handler := func(data []byte) error {
		var response struct {
			Type string         `json:"type"`
			Data *WSTradeUpdate `json:"data"`
		}
		
		if err := json.Unmarshal(data, &response); err != nil {
			return fmt.Errorf("failed to unmarshal trade update: %v", err)
		}
		
		if response.Type != ChannelTrades || response.Data == nil {
			return nil
		}
		
		return callback(response.Data)
	}
	
	ws.AddHandler(ChannelTrades, handler)
	return ws.Subscribe(ChannelTrades, fmt.Sprintf("%d", marketId))
}

// StreamMarkPrice subscribes to mark price updates for a market ID
func (ws *WSClient) StreamMarkPrice(ctx context.Context, marketId uint8, callback func(*WSTickerUpdate) error) error {
	handler := func(data []byte) error {
		var response struct {
			Type string           `json:"type"`
			Data *WSTickerUpdate  `json:"data"`
		}
		
		if err := json.Unmarshal(data, &response); err != nil {
			return fmt.Errorf("failed to unmarshal mark price update: %v", err)
		}
		
		if response.Type != ChannelMarkPrice || response.Data == nil {
			return nil
		}
		
		return callback(response.Data)
	}
	
	ws.AddHandler(ChannelMarkPrice, handler)
	return ws.Subscribe(ChannelMarkPrice, fmt.Sprintf("%d", marketId))
}

// Account Data Streams (require authentication)

// StreamAccount subscribes to account balance updates
func (ws *WSClient) StreamAccount(ctx context.Context, callback func(*WSAccountUpdate) error) error {
	if ws.authToken == "" {
		return fmt.Errorf("authentication token required for account streams")
	}
	
	handler := func(data []byte) error {
		var response struct {
			Type string            `json:"type"`
			Data *WSAccountUpdate  `json:"data"`
		}
		
		if err := json.Unmarshal(data, &response); err != nil {
			return fmt.Errorf("failed to unmarshal account update: %v", err)
		}
		
		if response.Type != ChannelAccount || response.Data == nil {
			return nil
		}
		
		return callback(response.Data)
	}
	
	ws.AddHandler(ChannelAccount, handler)
	return ws.Subscribe(ChannelAccount, "")
}

// StreamOrders subscribes to order status updates
func (ws *WSClient) StreamOrders(ctx context.Context, callback func(*WSOrderUpdate) error) error {
	if ws.authToken == "" {
		return fmt.Errorf("authentication token required for order streams")
	}
	
	handler := func(data []byte) error {
		var response struct {
			Type string         `json:"type"`
			Data *WSOrderUpdate `json:"data"`
		}
		
		if err := json.Unmarshal(data, &response); err != nil {
			return fmt.Errorf("failed to unmarshal order update: %v", err)
		}
		
		if response.Type != ChannelOrders || response.Data == nil {
			return nil
		}
		
		return callback(response.Data)
	}
	
	ws.AddHandler(ChannelOrders, handler)
	return ws.Subscribe(ChannelOrders, "")
}

// Unsubscribe methods

// UnsubscribeOrderBook unsubscribes from order book updates
func (ws *WSClient) UnsubscribeOrderBook(marketId uint8) error {
	ws.RemoveHandler(ChannelOrderBook)
	return ws.Unsubscribe(ChannelOrderBook, fmt.Sprintf("%d", marketId))
}

// UnsubscribeTicker unsubscribes from ticker updates
func (ws *WSClient) UnsubscribeTicker(marketId uint8) error {
	ws.RemoveHandler(ChannelTicker)
	return ws.Unsubscribe(ChannelTicker, fmt.Sprintf("%d", marketId))
}

// UnsubscribeTrades unsubscribes from trade updates
func (ws *WSClient) UnsubscribeTrades(marketId uint8) error {
	ws.RemoveHandler(ChannelTrades)
	return ws.Unsubscribe(ChannelTrades, fmt.Sprintf("%d", marketId))
}

// UnsubscribeMarkPrice unsubscribes from mark price updates
func (ws *WSClient) UnsubscribeMarkPrice(marketId uint8) error {
	ws.RemoveHandler(ChannelMarkPrice)
	return ws.Unsubscribe(ChannelMarkPrice, fmt.Sprintf("%d", marketId))
}

// UnsubscribeAccount unsubscribes from account updates
func (ws *WSClient) UnsubscribeAccount() error {
	ws.RemoveHandler(ChannelAccount)
	return ws.Unsubscribe(ChannelAccount, "")
}

// UnsubscribeOrders unsubscribes from order updates
func (ws *WSClient) UnsubscribeOrders() error {
	ws.RemoveHandler(ChannelOrders)
	return ws.Unsubscribe(ChannelOrders, "")
}

// Utility methods for data conversion

// ConvertWSOrderBookToKubi converts WebSocket order book data to kubi format
func ConvertWSOrderBookToKubi(wsData *WSOrderBookUpdate, marketId uint8) (*OrderBookResponse, error) {
	orderBooks := make([]OrderBook, 0, 1)
	
	bids := make([]PriceLevel, 0, len(wsData.Bids))
	for _, bid := range wsData.Bids {
		if len(bid) >= 2 {
			bids = append(bids, PriceLevel{
				Price:    bid[0],
				Quantity: bid[1],
			})
		}
	}
	
	asks := make([]PriceLevel, 0, len(wsData.Asks))
	for _, ask := range wsData.Asks {
		if len(ask) >= 2 {
			asks = append(asks, PriceLevel{
				Price:    ask[0],
				Quantity: ask[1],
			})
		}
	}
	
	orderBook := OrderBook{
		MarketId: marketId, // Use provided market ID directly
		Bids:     bids,
		Asks:     asks,
	}
	
	orderBooks = append(orderBooks, orderBook)
	
	return &OrderBookResponse{
		ResultCode: ResultCode{Code: CodeOK},
		OrderBooks: orderBooks,
	}, nil
}

// ConvertWSTickerToKubi converts WebSocket ticker data to kubi format
func ConvertWSTickerToKubi(wsData *WSTickerUpdate) (*WSTickerUpdate, error) {
	// For now, just return the same structure
	// This might need conversion based on the actual kubi PriceUpdate structure
	return wsData, nil
}

// Note: GetMarketIdFromSymbol function removed as WebSocket streams now use market_id directly