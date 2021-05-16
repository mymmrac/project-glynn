package main

import (
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/mymmrac/project-glynn/pkg/repository"
	"github.com/mymmrac/project-glynn/pkg/server"
	"github.com/mymmrac/project-glynn/pkg/server/httpapi"
	"github.com/sirupsen/logrus"
)

var cli struct {
	Port string `kong:"default='8080',help='Server port'"`

	CassandraURL  string `kong:"default='localhost',help='Cassandra URL'"`
	CassandraUser string `kong:"default='',help='Cassandra URL'"`
	CassandraPass string `kong:"default='',help='Cassandra URL'"`
}

func main() {
	log := logrus.StandardLogger()
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "02.01.2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	log.SetLevel(logrus.DebugLevel)

	ctx := kong.Parse(&cli)

	switch ctx.Command() {
	case "":
		log.Infof("Starting server on port: %s", cli.Port)

		log.Infof("Connecting to cassandra on url: '%s', user: '%s'", cli.CassandraURL, cli.CassandraUser)
		cassandra := repository.NewCassandraRepository(log)
		err := cassandra.Connect(cli.CassandraURL, cli.CassandraUser, cli.CassandraPass)
		defer cassandra.Close()
		if err != nil {
			log.Error("Failed to connect to cassandra: ", err)
			return
		}
		log.Info("Connected")

		service := server.NewService(cassandra, log)

		httpServer := httpapi.NewServer(service, log)
		log.Info("Listening on port:", cli.Port)
		if err := http.ListenAndServe(":"+cli.Port, httpServer); err != nil {
			log.Error("Failed to server: ", err)
			return
		}
	default:
		log.Error("Unknown command: ", ctx.Command())
		return
	}
}
