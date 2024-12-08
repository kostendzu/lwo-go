package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo/internal/app"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize TODO application: %v", err)
	}

	server := a.StartServer()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	a.StartBackgroundTask(sigs)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	log.Printf("Server started on %s", server.Addr)

	// Ждем сигнал для завершения программы
	sig := <-sigs
	log.Printf("Main function received signal: %v. Initiating shutdown...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Да, я использовал этот канал и для закрытия фоновой задачи
	//(Лучше бы использовал контекст, конечно)
	sigs <- sig

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped")
}
