package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"retask/api/handler"
	"retask/config"
	"retask/internal/packing"
	"retask/server"
	"syscall"
)

func main() {
	conf, err := config.New()
	if err != nil {
		logrus.Fatal(err)
	}

	logger := logrus.WithField("method", "main")

	packingRepo := packing.New()
	h := handler.New(conf, packingRepo)
	serv, err := server.New(conf, h)
	if err != nil {
		logger.WithField("error", err).Fatal("failed to init server")
	}

	serv.ListenAndServe(logger)

	// This allows us to listen for interrupts (ctrl+c, shutting down the run in goland/vscode, etc)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	logger.Info("signal interrupt detected, shutting down ...")

	// shutdown the server
	if err := serv.Shutdown(); err != nil {
		logger.Fatalf("failed to shutdown server with error: %v", err)
	}

	// there was an interrupt so exit with code 1
	os.Exit(1)

}
