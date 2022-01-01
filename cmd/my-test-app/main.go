package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/larwef/rpi-docker-test/internal/server"
	"github.com/larwef/rpi-docker-test/internal/storage"
	"github.com/larwef/rpi-docker-test/pkg/enemy"
	"google.golang.org/grpc"
)

// Version injected at compile time.
var version = "No version provided"

func main() {
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := realMain(ctx)
	done()
	if err != nil {
		log.Fatal(err)
	}
}

func realMain(ctx context.Context) error {
	log.Printf("Starting my-test-app %s\n", version)

	// Set up listener.
	port := os.Getenv("PORT")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	defer listener.Close()

	// Set up db connection.
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbURI := fmt.Sprintf("host=%s user=%s password=%s port=%d database=%s sslmode=disable", dbHost, dbUser, dbPass, 5432, dbName)
	db, err := sql.Open("pgx", dbURI)
	if err != nil {
		return fmt.Errorf("unable to open db connection: %v", err)
	}
	defer db.Close()

	// Make sure database is ready before trying to initialize storage, or else
	// it will fail.
	if err := PingRetry(ctx, db, 5*time.Second, 60*time.Second); err != nil {
		return err
	}
	store, err := storage.NewEnemyStore(db)
	if err != nil {
		return fmt.Errorf("unable to initialize storage: %v", err)
	}

	opts := []grpc.ServerOption{}
	srv := grpc.NewServer(opts...)
	enemy.RegisterEnemyServiceServer(srv, server.New(store))

	errCh := make(chan error)
	go func() {
		if err := srv.Serve(listener); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		srv.GracefulStop()
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func PingRetry(ctx context.Context, db *sql.DB, pingInterval, timeout time.Duration) error {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()
	timeoutExceeded := time.After(timeout)
	for {
		if err := db.PingContext(ctx); err == nil {
			return nil
		} else {
			log.Printf("connecting to db failed: %v. Retrying in %s", err, pingInterval)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeoutExceeded:
			return fmt.Errorf("db connection timed out after %s", timeout)
		case <-ticker.C:
			continue
		}
	}
}
