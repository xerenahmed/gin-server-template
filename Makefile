bin:
	GOARCH=amd64 GOOS=linux go build -o build/server main.go

clean:
	docker compose rm -fv

down:
	docker compose down -v --rmi 'local'

start: bin clean down
	docker compose up
