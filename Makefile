.PHONY: build-frontend run

build-frontend:
	npm run build --prefix frontend

run: build-frontend
	go run main.go