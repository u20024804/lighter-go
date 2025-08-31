package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WSClient struct {
	config      *WSConfig
	conn        *websocket.Conn
	mu          sync.RWMutex
	handlers    map[string][]WSHandler
	isConnected bool
	stopCh      chan struct{}
	reconnectCh chan struct{}
	authToken   string
	
	// For managing subscriptions
	subscriptions map[string]bool
}

type WSHandler func(data []byte) error

// NewWSClient creates a new WebSocket client
func NewWSClient(config *WSConfig) *WSClient {
	if config == nil {
		config = DefaultWSConfig()
	}
	
	return &WSClient{
		config:        config,
		handlers:      make(map[string][]WSHandler),
		subscriptions: make(map[string]bool),
		stopCh:        make(chan struct{}),
		reconnectCh:   make(chan struct{}, 1),
	}
}

// SetAuthToken sets the authentication token for authenticated streams
func (ws *WSClient) SetAuthToken(token string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.authToken = token
}

// Connect establishes WebSocket connection
func (ws *WSClient) Connect(ctx context.Context) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	
	if ws.isConnected {
		return nil
	}
	
	u, err := url.Parse(ws.config.URL)
	if err != nil {
		return fmt.Errorf("invalid WebSocket URL: %v", err)
	}
	
	dialer := websocket.Dialer{
		HandshakeTimeout: 45 * time.Second,
	}
	
	headers := http.Header{}
	if ws.authToken != "" {
		headers.Set("Authorization", "Bearer "+ws.authToken)
	}
	
	conn, _, err := dialer.Dial(u.String(), headers)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %v", err)
	}
	
	ws.conn = conn
	ws.isConnected = true
	
	// Start message handler goroutines
	go ws.readMessages(ctx)
	go ws.ping(ctx)
	
	// Resubscribe to existing channels
	go ws.resubscribe()
	
	log.Println("[WSClient] Connected to Lighter WebSocket")
	return nil
}

// Disconnect closes the WebSocket connection
func (ws *WSClient) Disconnect() error {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	
	if !ws.isConnected {
		return nil
	}
	
	close(ws.stopCh)
	ws.isConnected = false
	
	if ws.conn != nil {
		err := ws.conn.Close()
		ws.conn = nil
		return err
	}
	
	log.Println("[WSClient] Disconnected from Lighter WebSocket")
	return nil
}

// IsConnected returns connection status
func (ws *WSClient) IsConnected() bool {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.isConnected
}

// Subscribe subscribes to a WebSocket channel
func (ws *WSClient) Subscribe(channel, symbol string) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	
	subscriptionKey := channel
	if symbol != "" {
		subscriptionKey = fmt.Sprintf("%s:%s", channel, symbol)
	}
	
	ws.subscriptions[subscriptionKey] = true
	
	if !ws.isConnected {
		return fmt.Errorf("WebSocket not connected")
	}
	
	msg := WSSubscribeMessage{
		Type:    MessageTypeSubscribe,
		Channel: channel,
		Symbol:  symbol,
	}
	
	return ws.sendMessage(msg)
}

// Unsubscribe unsubscribes from a WebSocket channel
func (ws *WSClient) Unsubscribe(channel, symbol string) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	
	subscriptionKey := channel
	if symbol != "" {
		subscriptionKey = fmt.Sprintf("%s:%s", channel, symbol)
	}
	
	delete(ws.subscriptions, subscriptionKey)
	
	if !ws.isConnected {
		return nil // Already disconnected
	}
	
	msg := WSUnsubscribeMessage{
		Type:    MessageTypeUnsubscribe,
		Channel: channel,
		Symbol:  symbol,
	}
	
	return ws.sendMessage(msg)
}

// AddHandler adds a message handler for a specific channel
func (ws *WSClient) AddHandler(channel string, handler WSHandler) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	
	if ws.handlers[channel] == nil {
		ws.handlers[channel] = make([]WSHandler, 0)
	}
	ws.handlers[channel] = append(ws.handlers[channel], handler)
}

// RemoveHandler removes all handlers for a channel
func (ws *WSClient) RemoveHandler(channel string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	
	delete(ws.handlers, channel)
}

func (ws *WSClient) sendMessage(msg interface{}) error {
	if ws.conn == nil {
		return fmt.Errorf("WebSocket not connected")
	}
	
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}
	
	ws.conn.SetWriteDeadline(time.Now().Add(ws.config.WriteTimeout))
	return ws.conn.WriteMessage(websocket.TextMessage, data)
}

func (ws *WSClient) readMessages(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[WSClient] Panic in readMessages: %v", r)
		}
	}()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ws.stopCh:
			return
		default:
			if ws.conn == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			
			ws.conn.SetReadDeadline(time.Now().Add(ws.config.ReadTimeout))
			_, data, err := ws.conn.ReadMessage()
			if err != nil {
				log.Printf("[WSClient] Read error: %v", err)
				ws.handleDisconnect(ctx)
				return
			}
			
			ws.handleMessage(data)
		}
	}
}

func (ws *WSClient) handleMessage(data []byte) {
	var msg WSMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("[WSClient] Failed to unmarshal message: %v", err)
		return
	}
	
	// Handle pong messages
	if msg.Type == MessageTypePong {
		return
	}
	
	// Route message to appropriate handlers
	ws.mu.RLock()
	handlers := ws.handlers[msg.Type]
	ws.mu.RUnlock()
	
	for _, handler := range handlers {
		go func(h WSHandler) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[WSClient] Handler panic: %v", r)
				}
			}()
			
			if err := h(data); err != nil {
				log.Printf("[WSClient] Handler error: %v", err)
			}
		}(handler)
	}
}

func (ws *WSClient) ping(ctx context.Context) {
	ticker := time.NewTicker(ws.config.PingInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ws.stopCh:
			return
		case <-ticker.C:
			if ws.isConnected && ws.conn != nil {
				ping := WSMessage{Type: MessageTypePing}
				if err := ws.sendMessage(ping); err != nil {
					log.Printf("[WSClient] Failed to send ping: %v", err)
				}
			}
		}
	}
}

func (ws *WSClient) handleDisconnect(ctx context.Context) {
	ws.mu.Lock()
	ws.isConnected = false
	if ws.conn != nil {
		ws.conn.Close()
		ws.conn = nil
	}
	ws.mu.Unlock()
	
	log.Println("[WSClient] Connection lost, attempting to reconnect...")
	
	select {
	case ws.reconnectCh <- struct{}{}:
	default:
	}
	
	go ws.reconnect(ctx)
}

func (ws *WSClient) reconnect(ctx context.Context) {
	attempt := 0
	
	for attempt < ws.config.MaxReconnects {
		select {
		case <-ctx.Done():
			return
		default:
		}
		
		attempt++
		log.Printf("[WSClient] Reconnect attempt %d/%d", attempt, ws.config.MaxReconnects)
		
		time.Sleep(ws.config.ReconnectDelay)
		
		if err := ws.Connect(ctx); err != nil {
			log.Printf("[WSClient] Reconnect failed: %v", err)
			continue
		}
		
		log.Println("[WSClient] Reconnected successfully")
		return
	}
	
	log.Printf("[WSClient] Max reconnect attempts reached, giving up")
}

func (ws *WSClient) resubscribe() {
	ws.mu.RLock()
	subscriptions := make(map[string]bool)
	for k, v := range ws.subscriptions {
		subscriptions[k] = v
	}
	ws.mu.RUnlock()
	
	for subscriptionKey := range subscriptions {
		channel := subscriptionKey
		symbol := ""
		
		// Parse channel:symbol format
		if idx := len(subscriptionKey); idx > 0 {
			if colonIdx := 0; colonIdx < len(subscriptionKey) {
				for i, c := range subscriptionKey {
					if c == ':' {
						colonIdx = i
						break
					}
				}
				if colonIdx > 0 && colonIdx < len(subscriptionKey)-1 {
					channel = subscriptionKey[:colonIdx]
					symbol = subscriptionKey[colonIdx+1:]
				}
			}
		}
		
		msg := WSSubscribeMessage{
			Type:    MessageTypeSubscribe,
			Channel: channel,
			Symbol:  symbol,
		}
		
		if err := ws.sendMessage(msg); err != nil {
			log.Printf("[WSClient] Failed to resubscribe to %s: %v", subscriptionKey, err)
		} else {
			log.Printf("[WSClient] Resubscribed to %s", subscriptionKey)
		}
		
		time.Sleep(100 * time.Millisecond) // Rate limit resubscriptions
	}
}