package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/dreamsofcode-io/orders-api/model"
	"github.com/dreamsofcode-io/orders-api/repository/order"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New() *App {
	app := &App{
		router: loadRoutes(),
		rdb:    redis.NewClient(&redis.Options{}),
	}

	repo := order.RedisRepo{
		Client: app.rdb,
	}

	now := time.Now()
	o := model.Order{
		OrderID:    2,
		CustomerID: uuid.New(),
		LineItems:  []model.LineItem{},
		CreatedAt:  &now,
	}

	err := repo.Insert(context.Background(), o)

	res, err := repo.FindAll(context.Background(), order.FindAllPage{})
	fmt.Println(res, err)

	println(repo.Client)

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	fmt.Println("Starting server")

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}

	return nil
}
