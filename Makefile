.PHONY: run test fmt lint migrate-up migrate-down

run:
	go run cmd/api/main.go

test:
	go test ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run

migrate-up:
	# Add your migration tool command here
	echo "Applying migrations..."

migrate-down:
	# Add your migration tool command here
	echo "Reverting migrations..."
