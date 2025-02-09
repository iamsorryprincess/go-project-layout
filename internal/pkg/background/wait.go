package background

import (
	"os"
	"os/signal"
	"syscall"
)

func Wait() os.Signal {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
	s := <-exit
	signal.Stop(exit)
	close(exit)
	return s
}
