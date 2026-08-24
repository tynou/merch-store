package domain

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int32  `json:"quantity"`
}

type ReceivedTransfer struct {
	FromUser string `json:"fromUser"`
	Amount   int32  `json:"amount"`
}

type SentTransfer struct {
	ToUser string `json:"toUser"`
	Amount int32  `json:"amount"`
}

type CoinHistory struct {
	Received []ReceivedTransfer `json:"received"`
	Sent     []SentTransfer     `json:"sent"`
}
