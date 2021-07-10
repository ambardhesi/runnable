## Runnable
Runnable is a simple Linux process runner, implemented as a REST API.

## Build
```
$ make binaries
```

This will create the `runnable` (CLI client) and `runnable-client` (Server) binaries in the `bin` dir. 

## Test
make test

## Server
Start the server by following the build step, then running `./runnable`.

This will start up the server on localhost at port 8080. The server address is currently hardcoded inside runnable/main.go.

The server provides the following endpoints : 
* POST `/job` to start a job
* GET `/job/:id` to get a job
* POST `/job/:id/stop` to stop a job
* GET `/job/:id/logs` to get the logs for a job (stdout and stderr)

To interact with the server, use any tool to make REST calls (eg curl, Postman.

Sample call : 
```
curl -X POST -H "Content-Type: application/json" \
-d '{"command": "echo hello world"}'\ 
http://localhost:8080/job
```

## Client
The `./runnable-client` binary provides a CLI interface for making calls to the server.

Run `./runnable-client --help` for usage instructions.

Sample command : 
```
./runnable-client start echo hello world
```
