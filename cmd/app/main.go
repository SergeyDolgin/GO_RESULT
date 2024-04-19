package main

import (
	"context"
	"log"
	"my_fund/internal/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	file, err := os.OpenFile("/app/logs/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(file)

	srv, err := service.NewService(ctx)

	if err != nil {
		log.Println("main/NewService: ", err)
		return
	}

	go func() {
		<-ctx.Done()

		srv.Stop()

		if err = file.Close(); err != nil {
			log.Println("LogFile/Close: ", err)
		}
		os.Exit(1)

	}()

	srv.Run()
}
