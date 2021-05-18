package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/mymmrac/project-glynn/pkg/repository"
	"github.com/mymmrac/project-glynn/pkg/server"
	"github.com/mymmrac/project-glynn/pkg/server/httpapi"
	"github.com/sirupsen/logrus"
)

const timeoutThreshold = 5 * time.Second

var cli struct {
	Port string `kong:"default='8080',help='Server port'"`

	CassandraInit bool `kong:"default='false',help='Create keyspace & tables if not exist'"`
	// TODO move to env
	CassandraURL  string `kong:"default='localhost',help='Cassandra URL'"`
	CassandraUser string `kong:"default='',help='Cassandra User'"`
	CassandraPass string `kong:"default='',help='Cassandra Pass'"`
}

func main() {
	// TODO move to other files
	log := logrus.StandardLogger()
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = time.RFC1123
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	log.SetLevel(logrus.DebugLevel)

	ctx := kong.Parse(&cli)

	switch ctx.Command() {
	case "":
		log.Infof("Starting server on port: %s", cli.Port)

		log.Infof("Connecting to cassandra on url: '%s', user: '%s'", cli.CassandraURL, cli.CassandraUser)
		if cli.CassandraInit {
			log.Info("Creating keyspace & tables if not exist")
		}
		cassandra := repository.NewCassandraRepository(log)
		err := cassandra.Connect(cli.CassandraURL, cli.CassandraUser, cli.CassandraPass, cli.CassandraInit)
		defer cassandra.Close()
		if err != nil {
			log.Error("Failed to connect to cassandra: ", err)
			return
		}
		log.Info("Connected")

		service := server.NewService(cassandra, log)

		httpServer := httpapi.NewServer(service, log)

		srv := http.Server{
			Addr:    ":" + cli.Port,
			Handler: httpServer,
		}

		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error("Failed to server: ", err)
				return
			}
		}()
		log.Info("Listening on port:", cli.Port)

		<-done
		log.Info("Stopping server")

		ctx, cancel := context.WithTimeout(context.Background(), timeoutThreshold)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("Server shutdown failed: ", err)
			return
		}

		log.Info("Server stopped")
	default:
		log.Error("Unknown command: ", ctx.Command())
		return
	}
}
