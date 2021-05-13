# Project Glynn

CLI chat with server & client for communication in different chat rooms with any number of users.

For simplicity as "chat room" will be used just "room". 

## Server

* [ ] Configs:
  * [ ] Read configs
  * [ ] Parse CLI args
* [ ] Types:
  * [ ] User, Room, Connection, Message
  * [ ] Data repository
  * [ ] Server (HTTP & gRPC)
* [ ] Docker:
  * [ ] App build container
  * [ ] App run container
  * [ ] Cassandra container
* [ ] Cassandra:
  * [ ] Connect to Cassandra
  * [ ] Init Cassandra's keyspace & tables
* [ ] Basic info:
  * [ ] Start server (display initial server info)
  * [ ] Logging
  * [ ] Swagger UI
* [ ] HTTP:
  * [ ] Handle admin authentication middleware
  * [ ] Handle room creation
  * [ ] Handle user connection to room
  * [ ] Handle if user is new
  * [ ] Handle user connection status
  * [ ] Handle user disconnection from room
  * [ ] Handle new messages from users
  * [ ] Handle user reconnection (send unreaded messages)
  * [ ] Handle room deletion
  * [ ] Handle server info
* [ ] gRPC:
  * [ ] *Future plans*
* [ ] Encription
  * [ ] *Future plans*

## Client

* [ ] Configs:
  * [ ] Read configs
  * [ ] Parse CLI args
* [ ] Types:
  * [ ] Service
* [ ] HTTP:
  * [ ] Create connection to room
  * [ ] Read messages
  * [ ] Connection status (for server)
  * [ ] Get unreaded messages
  * [ ] Format massages
  * [ ] Create room
  * [ ] Delete room
  * [ ] Server info
* [ ] gRPC:
  * [ ] *Future plans*
* [ ] Encription
  * [ ] *Future plans*

