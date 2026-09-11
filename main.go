package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"entrevista/core/adapter/handler"
	"entrevista/core/adapter/repository"
	"entrevista/core/service"
	"entrevista/platform"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	db := platform.NewPostgresSQL()
	if err := platform.Migrate(db); err != nil {
		log.Fatalf("error al migrar el esquema: %v", err)
	}
	employeeRepo := repository.NewEmployeeRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	employeeSvc := service.NewEmployeeService(employeeRepo)
	taskSvc := service.NewTaskService(taskRepo, employeeRepo)
	employeeHandler := handler.NewEmployeeHandler(employeeSvc, taskSvc)
	taskHandler := handler.NewTaskHandler(taskSvc)
	router := handler.NewRouter(employeeHandler, taskHandler)
	srv := &http.Server{
		Addr:              ":8000",
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() {
		log.Printf("servidor escuchando en %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error al iniciar el servidor: %v", err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("apagando el servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("apagado forzado: %v", err)
	}
	log.Println("servidor detenido")
}
