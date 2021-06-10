bin:
	GOARCH=amd64 GOOS=linux go build -o build/server main.go

start: bin
	docker compose down -v
	docker rmi $IMAGE
	docker compose up