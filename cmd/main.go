package main

import (
	"avito/internal/factory"
	"avito/internal/handlers"
	"avito/internal/repository"
	"avito/internal/service"
	"avito/internal/usecase"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"os"
	"errors"
	"syscall"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"time"
	flag "github.com/spf13/pflag"
)

func main() {
	//загружаем env
	if err := godotenv.Load(); err != nil {
		fmt.Println("Не удалось загрузить .env:", err)
	}
	defport := os.Getenv("PORT")
	if defport == "" {
		defport = "8080"
	}
	//учли флаг --port
	port := flag.StringP("port", "p", defport, "порт для сервера")
	flag.Parse()
	//прочитаем переменные постгрес собрали адрес БД (DSN)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	//создаем пул соединений
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("не удалось подключиться к БД:", err)
	}
	defer pool.Close()
	//проверка подключились ли
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("БД не отвечает:", err)
	}
	//запускаем
	repo := repository.NewPostgresCourierRepository(pool)

	deliveryRepo:=repository.NewDeliveryRepository()
	timeFactory := factory.DeliveryTimeFactory{}
	deliveryUsecase:=usecase.NewDeliveryUsecase(pool, repo, deliveryRepo, timeFactory)
	deliveryHandler := handlers.NewDeliveryHandler(deliveryUsecase)
	svc := service.NewCourierService(repo)
	h:= handlers.NewCourierHandler(svc)

	r := chi.NewRouter()
	r.Get("/ping", h.Ping)
	r.Head("/healthcheck", h.Healthcheck)
	r.Post("/courier", h.Create)
	r.Get("/couriers", h.List)
	r.Get("/courier/{id}", h.GetById)
	r.Put("/courier/{id}", h.Update)

	r.Delete("/courier/{id}", h.Delete)
	
	r.Post("/delivery/assign", deliveryHandler.Assign)
	r.Post("/delivery/unassign", deliveryHandler.Unassign)

	addr := ":" + *port
	fmt.Println("Сервер запущен на", addr)



	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	srv := &http.Server{Addr: addr, Handler: r}
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()
	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("сервер упал: ", err)
		}
	case <-ctx.Done():
		fmt.Println("Shutting down service-courier")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			fmt.Println("ошибка при остановке сервера:", err)
		}
	}

}
