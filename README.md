# Liquor OS Terminal

This is a backend service for Liquor OS.

## Project Structure

* `/cmd/api/main.go`: application entrypoint
* `/internal/app/app.go`: application lifecycle
* `/internal/config/`: configuration loading
* `/internal/server/`: HTTP server setup
* `/internal/router/`: chi routes
* `/internal/handlers/`: HTTP handlers
* `/internal/middleware/`: middleware
* `/internal/store/`: database layer
* `/internal/db/`: DB initialization
* `/internal/models/`: domain models
* `/internal/service/`: business logic
* `/internal/logger/`: logging setup
* `/internal/shutdown/`: graceful shutdown utilities
* `/pkg/optional`: reusable utilities
* `/migrations/`: database migrations
* `/scripts/`: helper scripts
* `.env.example`: example environment variables
* `Makefile`: development commands
* `README.md`: this file

## Startup Flow

1.  `main.go` calls `app.Run()`.
2.  `app.Run()` initializes the config, logger, and DB connection.
3.  An `App` struct is created to hold these dependencies.
4.  The HTTP server is created with the router and middleware.
5.  The server is started in a goroutine.
6.  Graceful shutdown is handled by listening for OS signals.

## Dependency Wiring

Dependencies are explicitly injected via constructors. The `App` struct acts as a container for shared dependencies.

## Request Flow

1.  A request hits the chi router.
2.  Middleware is executed (logging, request ID, etc.).
3.  The request is passed to the appropriate handler.
4.  The handler calls a service method.
5.  The service method calls a store (repository) method.
6.  The store method executes a SQL query.
7.  Data is returned up the call stack and encoded as JSON.
