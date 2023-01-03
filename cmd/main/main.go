package main

import (
	"net"
	"net/http"
	"rest-api/cmd/main/internal/user"
	"rest-api/cmd/main/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	timeOut = time.Second * 15
	baseURL = ""
	host    = ":1234"
	network = "tcp"
)

func main() {
	// logging.Init()
	logger := logging.GetLogger()

	logger.Info("Create router")
	router := httprouter.New()

	logger.Info("Register user handler")
	userHandler := user.NewHandler(logger)
	userHandler.Register(router)

	start(router)
}

func start(router *httprouter.Router) {
	logger := logging.GetLogger()

	logger.Info("Start application")
	listener, err := net.Listen(network, host)

	if err != nil {
		panic(err)
	}

	server := http.Server{
		Handler:      router,
		WriteTimeout: timeOut,
		ReadTimeout:  timeOut,
	}

	logger.Infof("Server is listening port %s, start on connection %s", host, network)
	logger.Fatal(server.Serve(listener))
}
