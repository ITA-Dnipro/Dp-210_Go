package sender

import "fmt"

type PasswordCodeSender struct {
	emailSender *GmailEmailSender
}

func NewPasswordCodeSender(emailSender *GmailEmailSender) *PasswordCodeSender {
	return &PasswordCodeSender{emailSender: emailSender}
}

func (pcs *PasswordCodeSender) Send(to, code string) error {
	err := pcs.emailSender.Send(to, "Pasword restore",
		fmt.Sprintf("your restore code is %v, you should type it into our app", code))

	if err != nil {
		return fmt.Errorf("passw code sender: %w", err)
	}

	return err
}
