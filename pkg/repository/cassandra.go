package repository

import (
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
	// TODO use const
	cluster.ProtoVersion = 4
	cluster.Consistency = gocql.One

	c.log.Info("Using keyspace: ", keyspace)
	if err := c.createKeyspace(cluster); err != nil {
		// TODO error handling
		// c.log.Error("Failed to create keyspace: ", err)
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
		//	TODO check error
	}
}

func (c *Cassandra) createKeyspace(cluster *gocql.ClusterConfig) error {
	// TODO ask if needed (flag)
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// TODO fmt.Errorf("create keyspace: %w")
	return session.Query(createKeyspaceQuery).
		Exec()
}

func (c *Cassandra) createTables() error {
	// TODO ask if needed (flag)
	if err := c.session.Query(createUsersTable).Exec(); err != nil {
		// TODO creates room table: %w
		return err
	}

	if err := c.session.Query(createRoomsTable).Exec(); err != nil {
		return err
	}

	return c.session.Query(createMessagesTable).
		Exec()
}

func (c *Cassandra) GetMessageTime(messageID uuid.UUID) (time.Time, error) {
	var t time.Time
	err := c.session.Query(selectTimeOfMessage, messageID.String()).Scan(&t)
	return t, err
}

func (c *Cassandra) GetMessages(roomID uuid.UUID, afterTime time.Time, limit uint) ([]message.Message, error) {
	it := c.session.Query(selectMessages, roomID.String(), afterTime, limit).
		Iter()
	scanner := it.Scanner()

	messages := make([]message.Message, it.NumRows())
	i := it.NumRows() - 1
	for scanner.Next() {
		var messageIDStr, userIDStr, roomIDStr string
		var msg message.Message
		err := scanner.Scan(&messageIDStr, &roomIDStr, &userIDStr, &msg.Text, &msg.Time)
		if err != nil {
			return nil, err
		}

		msg.ID, err = uuid.Parse(messageIDStr)
		if err != nil {
			return nil, err
		}
		msg.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}
		msg.RoomID, err = uuid.Parse(roomIDStr)
		if err != nil {
			return nil, err
		}

		messages[i] = msg
		i--
	}

	return messages, scanner.Err()
}

func (c *Cassandra) GetUsersFromIDs(uuids []uuid.UUID) ([]user.User, error) {
	uuidsStr := make([]string, len(uuids))
	// TODO move to func
	for i, id := range uuids {
		uuidsStr[i] = id.String()
	}
	scanner := c.session.Query(selectUsersByIDs, uuidsStr).Iter().Scanner()

	var users []user.User
	for scanner.Next() {
		var usr user.User
		var idSrt string
		if err := scanner.Scan(&idSrt, &usr.Username); err != nil {
			return nil, err
		}

		var err error
		usr.ID, err = uuid.Parse(idSrt)
		if err != nil {
			return nil, err
		}
		users = append(users, usr)
	}

	return users, scanner.Err()
}

func (c *Cassandra) SaveMessage(msg *message.Message) error {
	return c.session.Query(insertMessage,
		msg.ID.String(), msg.UserID.String(), msg.RoomID.String(), msg.Text, msg.Time).Exec()
}

func (c *Cassandra) IsRoomExist(roomID uuid.UUID) (bool, error) {
	var exist int
	if err := c.session.Query(selectIfRoomExist, roomID.String()).Scan(&exist); err != nil {
		return false, err
	}
	return exist >= 1, nil
}
