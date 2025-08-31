package client

const (
	CodeOK = 200
)

type ResultCode struct {
	Code    int32  `json:"code,example=200"`
	Message string `json:"message,omitempty"`
}

type NextNonce struct {
	ResultCode
	Nonce int64 `json:"nonce,example=722"`
}

type ApiKey struct {
	AccountIndex int64  `json:"account_index,example=3"`
	ApiKeyIndex  uint8  `json:"api_key_index,example=0"`
	Nonce        int64  `json:"nonce,example=722"`
	PublicKey    string `json:"public_key"`
}

type AccountApiKeys struct {
	ResultCode
	ApiKeys []*ApiKey `json:"api_keys"`
}

type TxHash struct {
	ResultCode
	TxHash string `json:"tx_hash,example=0x70997970C51812dc3A010C7d01b50e0d17dc79C8"`
}

type TransferFeeInfo struct {
	ResultCode
	TransferFee int64 `json:"transfer_fee_usdc"`
}

type AccountInfo struct {
	ResultCode
	AccountType              uint8  `json:"account_type,omitempty"`
	Index                    int64  `json:"index,omitempty"`
	L1Address               string `json:"l1_address,omitempty"`
	CancelAllTime           int64  `json:"cancel_all_time,omitempty"`
	TotalOrderCount         int64  `json:"total_order_count,omitempty"`
	TotalIsolatedOrderCount int64  `json:"total_isolated_order_count,omitempty"`
	PendingOrderCount       int64  `json:"pending_order_count,omitempty"`
	AvailableBalance        string `json:"available_balance,omitempty"`
	Status                  uint8  `json:"status,omitempty"`
	CreatedAt               int64  `json:"created_at,omitempty"`
	LastActiveAt            int64  `json:"last_active_at,omitempty"`
	MarketStats             []AccountMarketStats `json:"market_stats,omitempty"`
}

type AccountMarketStats struct {
	MarketId              uint8  `json:"market_id,omitempty"`
	OpenOrderCount        int64  `json:"open_order_count,omitempty"`
	Sign                  int8   `json:"sign,omitempty"`
	Position              string `json:"position,omitempty"`
	AvgEntryPrice         string `json:"avg_entry_price,omitempty"`
	PositionValue         string `json:"position_value,omitempty"`
	UnrealizedPnl         string `json:"unrealized_pnl,omitempty"`
	RealizedPnl           string `json:"realized_pnl,omitempty"`
}

type DetailedAccountsResponse struct {
	ResultCode
	DetailedAccounts []AccountInfo `json:"detailed_accounts,omitempty"`
}

type OrderBookResponse struct {
	ResultCode
	OrderBooks []OrderBook `json:"order_books,omitempty"`
}

type OrderBook struct {
	MarketId             uint8       `json:"market_id,omitempty"`
	Bids                 []PriceLevel `json:"bids,omitempty"`
	Asks                 []PriceLevel `json:"asks,omitempty"`
	// Market configuration fields from API
	SupportedSizeDecimals  int32  `json:"supported_size_decimals,omitempty"`
	SupportedPriceDecimals int32  `json:"supported_price_decimals,omitempty"`  
	SupportedQuoteDecimals int32  `json:"supported_quote_decimals,omitempty"`
	MinBaseAmount          string `json:"min_base_amount,omitempty"`
	MinQuoteAmount         string `json:"min_quote_amount,omitempty"`
	TakerFee              string `json:"taker_fee,omitempty"`
	MakerFee              string `json:"maker_fee,omitempty"`
	LiquidationFee        string `json:"liquidation_fee,omitempty"`
	Status                string `json:"status,omitempty"`
}

type PriceLevel struct {
	Price    string `json:"price,omitempty"`
	Quantity string `json:"quantity,omitempty"`
}

type OrdersResponse struct {
	ResultCode
	Orders []Order `json:"orders,omitempty"`
}

type Order struct {
	Id               string `json:"id,omitempty"`
	AccountIndex     int64  `json:"account_index,omitempty"`
	MarketId         uint8  `json:"market_id,omitempty"`
	ClientOrderIndex int64  `json:"client_order_index,omitempty"`
	IsAsk            uint8  `json:"is_ask,omitempty"`
	BaseQuantity     string `json:"base_quantity,omitempty"`
	Price            string `json:"price,omitempty"`
	OrderType        uint8  `json:"order_type,omitempty"`
	TimeInForce      uint8  `json:"time_in_force,omitempty"`
	ReduceOnly       uint8  `json:"reduce_only,omitempty"`
	TriggerPrice     string `json:"trigger_price,omitempty"`
	OrderExpiry      int64  `json:"order_expiry,omitempty"`
	CreatedAt        int64  `json:"created_at,omitempty"`
	Status           string `json:"status,omitempty"`
	FilledQuantity   string `json:"filled_quantity,omitempty"`
	RemainingQuantity string `json:"remaining_quantity,omitempty"`
}
