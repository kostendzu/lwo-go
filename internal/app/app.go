package app

import (
	"log"
	"net/http"
	"os"
	"time"
	"todo/internal/db"
	"todo/internal/handlers"
	"todo/pkg/config"
)

type App struct {
	handler *handlers.Handler
}

func NewApp() (*App, error) {
	a := &App{}
	err := a.initConfig()

	if err != nil {
		return nil, err
	}

	repository, err := db.TaskRepositoryInit()

	if err != nil {
		return nil, err
	}

	handler := handlers.NewHandler(repository)

	a.handler = handler
	return a, nil
}

func (a *App) initConfig() error {
	err := config.Load(".env")
	if err != nil {
		return err
	}
	log.Println("Configs are inited")
	return nil
}

func (a *App) StartBackgroundTask(stopChan <-chan os.Signal) {
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Println("Checking for overdue tasks...")
				count, err := a.handler.UpdateOverdueTasks()
				if err != nil {
					log.Println("Error checking overdue tasks:", err)
				}

				log.Println("Count of overdue tasks: ", count)
			case sig := <-stopChan:
				log.Printf("Received signal: %v. Stopping background task.", sig)
				return
			}
		}
	}()

}

func (a *App) StartServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", a.handler.HandleTasks)
	mux.HandleFunc("/tasks/", a.handler.HandleTaskByID)
	mux.HandleFunc("/tasks/complete/", a.handler.HandleCompleteTask)

	address := os.Getenv("SERVER_ADDRESS")

	server := &http.Server{
		Addr:    address,
		Handler: LoggerMiddleware(mux),
	}

	log.Printf("Starting server on: %v ", server.Addr)

	return server
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Логируем метод и путь запроса
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)

		// Логируем время выполнения
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}
