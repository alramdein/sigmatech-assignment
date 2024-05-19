
build:
	docker-compose up --build     

run: 
	go run main.go

tidy:
	go mod tidy

test:
	go test ./...

mockgen:
	mockgen -source=repository/limit.go -destination=mocks/mock_limit_repository.go -package=mocks /
	mockgen -source=repository/transaction.go -destination=mocks/mock_transaction_repository.go -package=mocks
	mockgen -source=repository/customer.go -destination=mocks/mock_customer_repository.go -package=mocks
