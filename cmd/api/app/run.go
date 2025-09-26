package app

import (
	"cv-platform/pkg/thread"
	"os"
	"os/signal"
	"syscall"
)

type Application interface {
	Start()
	Shutdown()
}

func Run() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	rg := thread.NewRoutineGroup()

	httpServer := NewHTTPServer()
	rg.Run(httpServer.Shutdown)

	rg.Wait()
}
