package main

import (
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/usercardjob/user"
	"os"
)

func main() {
	// Начать диалог с человеком, который еще не заведен в системе как пациент.
	if err := user.StartDialog(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
