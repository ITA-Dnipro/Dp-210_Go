package usecase

import (
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/client"
)

type PasswordCodeSender struct {
	//emailSender *GmailEmailSender
	kafka *client.Kafka
}

func NewPasswordCodeSender(kafka *client.Kafka) *PasswordCodeSender {
	return &PasswordCodeSender{kafka: kafka}
}

func (pcs *PasswordCodeSender) Send(to, code string) error {
	email := Email{
		Receiver: to,
		Text:     fmt.Sprintf("your restore code is %v, you should type it into our app", code),
		Title:    "Password restore",
	}

	if err := pcs.kafka.Send(client.PasswCodeTopic, email); err != nil {
		return fmt.Errorf("passw code sender: %w", err)
	}

	return nil
}

type Email struct {
	Title    string `json:"title"`
	Text     string `json:"text"`
	Receiver string `json:"receiver"`
}
