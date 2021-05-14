# Project Glynn

CLI chats with a server & client for communication in different chat rooms with any number of users.

For simplicity as "chat room" will be used just "room". 

🕒 - delayed for an unknown amount of time

## Server

* [ ] Configs:
  * [ ] 🕒 Read configs
  * [ ] Parse CLI args
* [ ] Types:
  * [ ] User, Message, Room
  * [ ] Data repository
  * [ ] Server
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
* [ ] Server (HTTP):
  * [ ] Handle if user is new
  * [ ] Handle new messages from users
  * [ ] Handle room creation
  * [ ] Handle room deletion
  * [ ] Handle admin authentication middleware    
  * [ ] Handle server info
  * [ ] 🕒 Handle user connection to room
  * [ ] 🕒 Handle user disconnection from room
  * [ ] 🕒 Handle user connection status
* [ ] Server (gRPC):
  * [ ] *Future plans*
* [ ] Encryption
  * [ ] *Future plans*

## Client

* [ ] Configs:
  * [ ] 🕒 Read configs
  * [ ] Parse CLI args
* [ ] Types:
  * [ ] Service
* [ ] Service (HTTP):
  * [ ] User creation
  * [ ] Read messages
  * [ ] Format massages
  * [ ] Create room
  * [ ] Delete room
  * [ ] Server info
  * [ ] 🕒 Connection status (for the server)
  * [ ] 🕒 Create connection to room
* [ ] Service (gRPC):
  * [ ] *Future plans*
* [ ] Encryption
  * [ ] *Future plans*

