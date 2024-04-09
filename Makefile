up:
	docker-compose up

down:
	@docker ps -aq | xargs docker rm -f || true
	if [ -d "data" ]; then \
		sudo rm -r data; \
	fi

run:
	go run cmd/currency/main.go cmd/currency/init.go