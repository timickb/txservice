package domain

type ClientRequest struct {
	SenderId   string `json:"sender_id,omitempty"`
	ReceiverId string `json:"receiver_id,omitempty"`
	Amount     int64  `json:"amount,omitempty"`
}

type HttpResponse struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}
