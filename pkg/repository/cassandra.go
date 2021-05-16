package repository

import (
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/mymmrac/project-glynn/pkg/data/message"
	"github.com/mymmrac/project-glynn/pkg/data/user"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	keyspace = "glynn"
)

type Cassandra struct {
	log     *logrus.Logger
	session *gocql.Session
}

func NewCassandraRepository(log *logrus.Logger) *Cassandra {
	return &Cassandra{
		log: log,
	}
}

func (c *Cassandra) Connect(cassandraURL, cassandraUser, cassandraPass string) error {
	cluster := gocql.NewCluster(cassandraURL)
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: cassandraUser, Password: cassandraPass}
	cluster.ProtoVersion = 4
	cluster.Consistency = gocql.One

	c.log.Info("Using keyspace: ", keyspace)
	if err := c.createKeyspace(cluster); err != nil {
		c.log.Error("Failed to create keyspace: ", err)
		return err
	}
	cluster.Keyspace = keyspace

	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}
	c.session = session

	err = c.createTables()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cassandra) Close() {
	if c.session != nil {
		c.session.Close()
	}
}

func (c *Cassandra) createKeyspace(cluster *gocql.ClusterConfig) error {
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Query("CREATE KEYSPACE IF NOT EXISTS " +
		keyspace + " WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };").
		Exec()
}

func (c *Cassandra) createTables() error {
	if err := c.session.Query("CREATE TABLE IF NOT EXISTS " +
		keyspace + ".users (id uuid PRIMARY KEY, username text);").Exec(); err != nil {
		return err
	}

	if err := c.session.Query("CREATE TABLE IF NOT EXISTS " +
		keyspace + ".rooms (id uuid PRIMARY KEY);").Exec(); err != nil {
		return err
	}

	return c.session.Query("CREATE TABLE IF NOT EXISTS " +
		keyspace + ".messages (id uuid, roomID uuid, userID uuid, text text, time timestamp, PRIMARY KEY (id, time));").
		Exec()
}

func (c *Cassandra) GetMessageTime(messageID uuid.UUID) (time.Time, error) {
	var t time.Time
	err := c.session.Query("SELECT time FROM messages WHERE id = ?;", messageID.String()).Scan(&t)
	return t, err
}

func (c *Cassandra) GetMessages(roomID uuid.UUID, afterTime time.Time, limit uint) ([]message.Message, error) {
	scanner :=
		c.session.Query("SELECT id, roomID, userID, text, time FROM messages WHERE roomID = ? AND time > ? LIMIT ? ALLOW FILTERING;",
			roomID.String(), afterTime, limit).
			Iter().Scanner()

	var messages []message.Message
	for scanner.Next() {
		var messageIDStr, userIDStr, roomIDStr string
		var msg message.Message
		err := scanner.Scan(&messageIDStr, &roomIDStr, &userIDStr, &msg.Text, &msg.Time)
		if err != nil {
			return nil, err
		}
		msg.ID = uuid.MustParse(messageIDStr)
		msg.UserID = uuid.MustParse(userIDStr)
		msg.RoomID = uuid.MustParse(roomIDStr)

		messages = append(messages, msg)
	}

	return messages, scanner.Err()
}

func (c *Cassandra) GetUsersFromIDs(uuids []uuid.UUID) ([]user.User, error) {
	uuidsStr := make([]string, len(uuids))
	for i, id := range uuids {
		uuidsStr[i] = id.String()
	}
	scanner := c.session.Query("SELECT id, username FROM users WHERE id IN ?", uuidsStr).Iter().Scanner()

	var users []user.User
	for scanner.Next() {
		var usr user.User
		var idSrt string
		if err := scanner.Scan(&idSrt, &usr.Username); err != nil {
			return nil, err
		}
		usr.ID = uuid.MustParse(idSrt)
		users = append(users, usr)
	}

	return users, scanner.Err()
}

func (c *Cassandra) SaveMessage(msg *message.Message) error {
	return c.session.Query("INSERT INTO messages (id, userID, roomID, text, time) VALUES (?, ?, ?, ?, ?);",
		msg.ID.String(), msg.UserID.String(), msg.RoomID.String(), msg.Text, msg.Time).Exec()
}

func (c *Cassandra) IsRoomExist(roomID uuid.UUID) (bool, error) {
	var exist int
	if err := c.session.Query("SELECT count(*) FROM rooms WHERE id = ?;", roomID.String()).Scan(&exist); err != nil {
		return false, err
	}
	return exist == 1, nil
}
