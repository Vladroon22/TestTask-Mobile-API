package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	golog "github.com/Vladroon22/GoLog"
	"github.com/Vladroon22/TestTask-Mobile-API/config"
	"github.com/Vladroon22/TestTask-Mobile-API/internal/database"
	"github.com/Vladroon22/TestTask-Mobile-API/internal/handlers"
	"github.com/Vladroon22/TestTask-Mobile-API/internal/service"
	"github.com/gorilla/mux"
)

func main() {
	logger := golog.New()
	cnf := config.CreateConfig()

	db := database.NewDB(cnf, logger)
	if err := db.Connect(); err != nil {
		logger.Fatalln(err)
	}

	router := mux.NewRouter()
	srv := service.NewService()
	repo := database.NewRepo(db)
	h := handlers.NewHandlers(logger, repo, srv)

	router.HandleFunc("/sign-up", h.SignUP).Methods("POST")
	router.HandleFunc("/login", h.Login).Methods("POST") // auth

	authRouter := router.PathPrefix("/user/").Subrouter()
	authRouter.Use(h.AuthMiddleWare)

	authRouter.HandleFunc("/post", h.Post).Methods("POST")
	authRouter.HandleFunc("/checkPost/{id:[0-9]+}", h.GetPost).Methods("GET")

	authRouter.HandleFunc("/leftComm", h.Comment).Methods("POST")
	authRouter.HandleFunc("/readComm/{id:[0-9]+}", h.ReadComm).Methods("GET")

	authRouter.HandleFunc("/post/{id:[0-9]+}", h.Like).Methods("POST")
	authRouter.HandleFunc("/comment/{id:[0-9]+}", h.Like).Methods("POST")
	authRouter.HandleFunc("/getLiker/{id:[0-9]+}", h.GetLiker).Methods("GET")

	logger.Infoln("Server is listening --> localhost:8888")
	go http.ListenAndServe(":8888", router)

	exitSig := make(chan os.Signal, 1)
	signal.Notify(exitSig, syscall.SIGINT, syscall.SIGTERM)
	<-exitSig

	go func() {
		if err := db.CloseDB(); err != nil {
			logger.Fatalln(err)
		}
	}()
}
