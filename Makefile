.PHONY: frontend-dev frontend-build run

frontend-dev:
	npm run start --prefix frontend

frontend-build:
	npm install --include dev --prefix frontend
	npm run build --prefix frontend

run: frontend-build
	go run main.go