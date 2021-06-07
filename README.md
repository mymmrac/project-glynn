# Project Glynn

[![Go Status](https://github.com/mymmrac/project-glynn/actions/workflows/go.yml/badge.svg)](https://github.com/mymmrac/project-glynn/actions/workflows/go.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=mymmrac_project-glynn&metric=alert_status)](https://sonarcloud.io/dashboard?id=mymmrac_project-glynn)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=mymmrac_project-glynn&metric=bugs)](https://sonarcloud.io/dashboard?id=mymmrac_project-glynn)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=mymmrac_project-glynn&metric=coverage)](https://sonarcloud.io/dashboard?id=mymmrac_project-glynn)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=mymmrac_project-glynn&metric=ncloc)](https://sonarcloud.io/dashboard?id=mymmrac_project-glynn)

CLI chats with a server & client for communication in different chat rooms with any number of users.

For simplicity as "chat room" will be used just "room". 

🕒 - delayed for an unknown amount of time

## Server

* [ ] Configs:
  * [ ] 🕒 Read configs
  * [X] Parse CLI args
* [X] Types:
  * [X] User, Message, Room
  * [X] Data repository
  * [X] Service
* [ ] Docker:
  * [ ] App build container
  * [ ] App run container
  * [X] Cassandra container
* [X] Cassandra:
  * [X] Connect to Cassandra
  * [X] Init Cassandra's keyspace & tables
* [ ] Basic info:
  * [ ] Start server (display initial server info)
  * [ ] Logging
  * [ ] Swagger UI
* [ ] Service:
  * [X] Get messages
  * [X] Send message
  * [ ] Create room
  * [ ] Delete room
  * [ ] Validate room
  * [ ] Validate user
  * [ ] Validate message
  * [ ] Get info
* [ ] Server (HTTP):
  * [ ] Handle if user is new
  * [X] Handle get messages
  * [X] Handle new messages
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
  * [X] Parse CLI args
* [X] Types:
  * [X] Client
  * [X] Chat related
* [ ] Service (HTTP):
  * [ ] User creation
  * [X] Read messages
  * [X] Format massages
  * [ ] Create room
  * [ ] Delete room
  * [ ] Server info
  * [ ] 🕒 Connection status (for the server)
  * [ ] 🕒 Create connection to room
* [ ] Service (gRPC):
  * [ ] *Future plans*
* [ ] Encryption
  * [ ] *Future plans*

