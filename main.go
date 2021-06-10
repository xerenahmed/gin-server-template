package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	server2 "server-app/server"
	"syscall"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("can't load env: %v", err))
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		DisableColors: false,
	})

	if os.Getenv("MODE") == "dev" {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetReportCaller(true)
	}
}

func main() {
	server, err := server2.NewServer()
	if err != nil {
		logrus.WithError(err).Fatal("Sunucu oluşturulamadı.")
	}

	server.Load()
	go server.Run()

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGKILL, syscall.SIGINT)

	select {
	case <-c:
		logrus.Info("Sinyal tarafından emredildi.")
		break
	case <-server.CloseRecv:
		logrus.Info("Sunucu tarafından emredildi.")
		break
	}

	logrus.Info("Sunucu kapatılıyor.")
	server.Shutdown()
	logrus.Info("Sunucu kapandı.")
}
