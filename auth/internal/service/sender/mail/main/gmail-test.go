package main

import (
	"log"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/service/sender/mail"
)

func main() {
	ges, err := mail.NewGmailEmailSender("config.json", "token.json")
	if err != nil {
		log.Fatal(err)
	}

	if err = ges.Send("nicknema13@gmail.com", "email test", "gotcha"); err != nil {
		log.Fatal(err)
	}
}
