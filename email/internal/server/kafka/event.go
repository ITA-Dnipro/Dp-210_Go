package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/email/internal/usesase"
)

type Email struct {
	Title    string `json:"title"`
	Receiver string `json:"receiver"`
}

type EventHandler struct {
	Sender *usesase.GmailEmailSender
}

func (h *EventHandler) EmailFromEvent(payload []byte) error {
	var e Email
	if err := json.Unmarshal(payload, &e); err != nil {
		return fmt.Errorf("marshaling appointment:%w", err)
	}
	return h.Sender.Send(e.Title, "subj", e.Receiver)
}
