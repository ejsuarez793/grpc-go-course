# gRPC GO COURSE

This is a repo created following and Udemy course called "gRPC [Golang] Master Class: Build Modern API & Microservices" by Stephane Maarek.

It have implementation of server and client for a Greet app, Calculator app and a Blog app

Only the server of blog app is Dockerized with a docker-compose file which you can run by running the command:

`docker-compose up`

This will start mongo and server

And you can run the client by running on the root of app:

`go run /blog/client/client.go`

To try other "apps" you have to run the server and client manually, for instace:


`go run/calculator/server/server.go`
`go run/calculator/client/client.go`