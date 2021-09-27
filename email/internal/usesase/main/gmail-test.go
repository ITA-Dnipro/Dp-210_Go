package main

import (
	"github.com/ITA-Dnipro/Dp-210_Go/email/internal/usesase"
	"log"
)

func main() {
	ges, err := usesase.NewGmailEmailSender("config.json", "token.json")
	if err != nil {
		log.Fatal(err)
	}

	if err = ges.Send("dp210go@gmail.com", "email test", "gotcha"); err != nil {
		log.Fatal(err)
	}
}
