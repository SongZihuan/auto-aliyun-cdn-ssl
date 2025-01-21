package config

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func initSignal(sigChan chan os.Signal) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("init signal error: %v", r)
		}
	}()

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	return nil
}
