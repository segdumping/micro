package signal

import (
	"os"
	"os/signal"
	"syscall"
)

var signals = []os.Signal {
	syscall.SIGINT,
	syscall.SIGQUIT,
	syscall.SIGTERM,
}

func Notify(ch chan os.Signal) {
	signal.Notify(ch, signals...)
}

func Reset() {
	signal.Reset(signals...)
}
