package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/email/internal/usesase"
)

type Email struct {
	Title    string `json:"title"`
	Text     string `json:"text"`
	Receiver string `json:"receiver"`
}

type EventHandler struct {
	Sender *usesase.GmailEmailSender
}

func (h *EventHandler) EmailFromEvent(payload []byte) error {
	var e Email
	if err := json.Unmarshal(payload, &e); err != nil {
		return fmt.Errorf("unmarshaling email:%w", err)
	}

	return h.Sender.Send(e.Receiver, e.Title, e.Text)
}
