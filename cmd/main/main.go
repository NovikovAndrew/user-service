package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"rest-api/cmd/internal/config"
	"rest-api/cmd/internal/user"
	"rest-api/cmd/internal/user/db"
	"rest-api/pkg/client/mongodb"
	"rest-api/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	timeOut = time.Second * 15
	TCP     = "tcp"
	SOCKET  = "sock"
	UNIX    = "unix"
)

func main() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()

	mongoDBClient, err := mongodb.NewClient(
		context.Background(),
		cfg.MongoDB.Host,
		cfg.Listen.Port,
		cfg.MongoDB.Username,
		cfg.MongoDB.Password,
		cfg.MongoDB.Database,
		cfg.MongoDB.AuthDB,
	)

	if err != nil {
		panic(err)
	}

	storage := db.NewStorage(mongoDBClient, cfg.MongoDB.Collection, logger)
	u := user.User{
		ID:           "",
		Username:     "Test",
		PasswordHash: "12345",
		Email:        "test@example.com",
	}
	user1ID, err := storage.Create(context.Background(), u)

	if err != nil {
		panic(err)
	}

	logger.Info("USER CREATED", user1ID)

	logger.Info("Create router")
	router := httprouter.New()

	logger.Info("Register user handler")
	userHandler := user.NewHandler(logger)
	userHandler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("Start application")
	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == SOCKET {
		// path/to/binary
		// Dir{} Path/To
		// TODO: set output directory vscode vscode
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}

		// appDir = "/Users/andreynovikov/Desktop/all-project/Go/artdev-go-advanced/rest-api/build/app.sock"

		logger.Info("Create socket")
		socketPath := path.Join(appDir, "app.sock")
		logger.Info("Listen unix socket")

		listener, listenErr = net.Listen(UNIX, socketPath)
		logger.Infof("Server is listening unix socket: %s", socketPath)
	} else {
		logger.Info("Listen tcp socket")
		listener, listenErr = net.Listen(TCP, fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port))
		logger.Infof("Server is listening port %s:%s, start on connection %s", cfg.Listen.BindIp, cfg.Listen.Port, cfg.Listen.Type)
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := http.Server{
		Handler:      router,
		WriteTimeout: timeOut,
		ReadTimeout:  timeOut,
	}

	logger.Fatal(server.Serve(listener))
}
