.PHONY:

force: 
	migrate -path ./schema -database 'postgresql://postgres:11111@localhost:5320/postgres?sslmode=disable' force 1

mig-up:
	migrate -path ./schema -database 'postgresql://postgres:11111@localhost:5320/postgres?sslmode=disable' up 

mig-down:
	migrate -path ./schema -database 'postgresql://postgres:11111@localhost:5320/postgres?sslmode=disable' down 

build: 
	go build -o ./api cmd/main.go

run: build
	./api
