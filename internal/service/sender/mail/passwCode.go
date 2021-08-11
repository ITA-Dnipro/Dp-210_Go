package mail

import "fmt"

type PasswordCodeSender struct {
	emailSender *GmailEmailSender
}

func (pcs *PasswordCodeSender) Send(to, code string) error {
	err := pcs.emailSender.Send(to, "Pasword restore",
		fmt.Sprintf("your restore code is %v, you should type it into our app", code))

	return fmt.Errorf("passw code sender: %w", err)
}
