package main

import (
	"github.com/ITA-Dnipro/Dp-210_Go/email_sender"
	"log"
)

func main() {
	ges, err := email_sender.NewGmailEmailSender("config.json", "token.json")
	if err != nil {
		log.Fatal(err)
	}

	if err = ges.Send("nicknema13@gmail.com", "email test", "gotcha"); err != nil {
		log.Fatal(err)
	}
}
