package repository

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/mymmrac/project-glynn/pkg/data/message"
	"github.com/mymmrac/project-glynn/pkg/data/user"
	"github.com/mymmrac/project-glynn/pkg/uuid"
	"github.com/sirupsen/logrus"
)

const (
	keyspace            = "glynn"
	createKeyspaceQuery = "CREATE KEYSPACE IF NOT EXISTS " +
		keyspace + " WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };"

	createUsersTable    = "CREATE TABLE IF NOT EXISTS " + keyspace + ".users (id uuid PRIMARY KEY, username text);"
	createRoomsTable    = "CREATE TABLE IF NOT EXISTS " + keyspace + ".rooms (id uuid PRIMARY KEY);"
	createMessagesTable = "CREATE TABLE IF NOT EXISTS " +
		keyspace + ".messages (id uuid, roomID uuid, userID uuid, text text, time timestamp, PRIMARY KEY (roomID, time));"

	selectTimeOfMessage = "SELECT time FROM messages WHERE id = ? LIMIT 1 ALLOW FILTERING;"
	selectMessages      = "SELECT id, roomID, userID, text, time FROM messages " +
		"WHERE roomID = ? AND time > ? ORDER BY time DESC LIMIT ? ALLOW FILTERING;"
	selectUsersByIDs  = "SELECT id, username FROM users WHERE id IN ?"
	selectIfRoomExist = "SELECT count(*) FROM rooms WHERE id = ?;"

	insertMessage = "INSERT INTO messages (id, userID, roomID, text, time) VALUES (?, ?, ?, ?, ?);"
)

// Cassandra implementation of Repository
type Cassandra struct {
	session *gocql.Session
	log     *logrus.Logger
}

// NewCassandraRepository creates new Cassandra Repository
func NewCassandraRepository(log *logrus.Logger) *Cassandra {
	return &Cassandra{
		log: log,
	}
}

// Connect starts new connection and initializes database if needed
func (c *Cassandra) Connect(cassandraURL, cassandraUser, cassandraPass string, initDB bool) error {
	cluster := gocql.NewCluster(cassandraURL)
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: cassandraUser, Password: cassandraPass}
	// TODO use const ?
	cluster.ProtoVersion = 4
	cluster.Consistency = gocql.One
	cluster.Logger = c.log

	if initDB {
		if err := c.createKeyspace(cluster); err != nil {
			return fmt.Errorf("inti db: %w", err)
		}
	}
	cluster.Keyspace = keyspace

	session, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("create cluster: %w", err)
	}
	c.session = session

	if initDB {
		if err = c.createTables(); err != nil {
			return fmt.Errorf("inti db: %w", err)
		}
	}

	return nil
}

// Close connection
func (c *Cassandra) Close() {
	if c.session != nil {
		c.session.Close()
		//	TODO check error
	}
}

func (c *Cassandra) createKeyspace(cluster *gocql.ClusterConfig) error {
	session, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("create cluster: %w", err)
	}
	defer session.Close()

	if err := session.Query(createKeyspaceQuery).Exec(); err != nil {
		return fmt.Errorf("create keyspace %q: %w", keyspace, err)
	}
	return nil
}

func (c *Cassandra) createTables() error {
	if err := c.session.Query(createUsersTable).Exec(); err != nil {
		return fmt.Errorf("create users table: %w", err)
	}

	if err := c.session.Query(createRoomsTable).Exec(); err != nil {
		return fmt.Errorf("create rooms table: %w", err)
	}

	if err := c.session.Query(createMessagesTable).Exec(); err != nil {
		return fmt.Errorf("create messages table: %w", err)
	}
	return nil
}

func (c *Cassandra) GetMessageTime(messageID uuid.UUID) (time.Time, error) {
	var t time.Time
	if err := c.session.Query(selectTimeOfMessage, messageID.String()).Scan(&t); err != nil {
		return time.Time{}, fmt.Errorf("get time of message %s: %w", messageID, err)
	}
	return t, nil
}

func (c *Cassandra) GetMessages(roomID uuid.UUID, afterTime time.Time, limit uint) ([]message.Message, error) {
	it := c.session.Query(selectMessages, roomID.String(), afterTime, limit).Iter()
	scanner := it.Scanner()

	messages := make([]message.Message, it.NumRows())
	i := it.NumRows() - 1
	for scanner.Next() {
		var messageIDStr, userIDStr, roomIDStr string
		var msg message.Message
		err := scanner.Scan(&messageIDStr, &roomIDStr, &userIDStr, &msg.Text, &msg.Time)
		if err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}

		msg.ID, err = uuid.Parse(messageIDStr)
		if err != nil {
			return nil, fmt.Errorf("message id: %w", err)
		}
		msg.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("user id: %w", err)
		}
		msg.RoomID, err = uuid.Parse(roomIDStr)
		if err != nil {
			return nil, fmt.Errorf("room id: %w", err)
		}

		messages[i] = msg
		i--
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan messages: %w", err)
	}
	return messages, nil
}

func (c *Cassandra) GetUsersFromIDs(uuids []uuid.UUID) ([]user.User, error) {
	uuidsStr := uuid.ToStrings(uuids)
	scanner := c.session.Query(selectUsersByIDs, uuidsStr).Iter().Scanner()

	var users []user.User
	for scanner.Next() {
		var usr user.User
		var idSrt string
		if err := scanner.Scan(&idSrt, &usr.Username); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}

		var err error
		usr.ID, err = uuid.Parse(idSrt)
		if err != nil {
			return nil, fmt.Errorf("user id: %w", err)
		}
		users = append(users, usr)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan users: %w", err)
	}
	return users, nil
}

func (c *Cassandra) SaveMessage(msg *message.Message) error {
	err := c.session.Query(insertMessage,
		msg.ID.String(), msg.UserID.String(), msg.RoomID.String(), msg.Text, msg.Time).Exec()
	if err != nil {
		return fmt.Errorf("save message: %w", err)
	}
	return nil
}

func (c *Cassandra) IsRoomExist(roomID uuid.UUID) (bool, error) {
	var exist int
	if err := c.session.Query(selectIfRoomExist, roomID.String()).Scan(&exist); err != nil {
		return false, fmt.Errorf("check room: %w", err)
	}
	return exist >= 1, nil
}
