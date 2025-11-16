package main

import (
    "backend-trainee-assignment/internal/app"
    memory "backend-trainee-assignment/internal/infrastructure/persistance/in_memory"
    pg "backend-trainee-assignment/internal/infrastructure/persistance/postgres"
    "backend-trainee-assignment/internal/transport/http"
    "database/sql"
    "log"
    "math/rand"
    "net/http"
    "os"
    "time"
)

func main() {
    logger := log.New(os.Stdout, "[pr-reviewer] ", log.LstdFlags|log.Lshortfile)

    var store app.Store

    switch os.Getenv("STORE") {
    case "postgres":
        dsn := os.Getenv("DB_DSN")
        if dsn == "" {
            logger.Fatal("DB_DSN is not set")
        }
        db, err := sql.Open("postgres", dsn)
        if err != nil {
            logger.Fatalf("failed to connect db: %v", err)
        }
        if err := db.Ping(); err != nil {
            logger.Fatalf("db ping error: %v", err)
        }
        logger.Println("Using PostgreSQL store")
        store = pg.NewPostgresStore(db)

    default:
        logger.Println("Using InMemory store")
        store = memory.NewInMemoryStore()
    }

    randSrc := rand.New(rand.NewSource(time.Now().UnixNano()))
    svc := app.NewService(store, randSrc)

    handler := httpapi.NewHandler(svc)

    server := &http.Server{
        Addr:         ":8080",
        Handler:      handler,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 5 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    logger.Println("starting HTTP server on :8080")
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        logger.Fatalf("server error: %v", err)
    }
}
