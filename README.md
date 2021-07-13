## Runnable
Runnable is a simple Linux process runner, implemented as a REST API.

## Build
```
make certs
```
This will generate keys and certs for the CA, 2 clients (Alice and Bob) and a bad client for testing.

```
$ make binaries
```

This will create the `runnable` (CLI client) and `runnable-client` (Server) binaries.

## Test
```
make test
```

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

### Starting a job
```
$ ./runnable-client --ca $path-to-ca-cert --cert $path-to-client-cert --key $path-to-client-key start echo hello world`


Started job : {"jobID":"eabc5579-7f8e-48d5-ba57-dc6e17f7a3ad"}
```

### Stopping a job
```
./runnable-client --ca $path-to-ca-cert --cert $path-to-client-cert --key $path-to-client-key stop eabc5579-7f8e-48d5-ba57-dc6e17f7a3ad

```

### Getting a job's status
```
$ ./runnable-client --ca certs/ca-cert.pem --cert certs/alice-cert.pem --key certs/alice-key.pem get eabc5579-7f8e-48d5-ba57-dc6e17f7a3ad

Job : {"state":"Completed","exitCode":0,"startTime":"2021-07-13T06:57:49.401392804-04:00","endTime":"2021-07-13T06:57:49.401505369-04:00"}
```

### Getting a job's logs
```
$ ./runnable-client --ca certs/ca-cert.pem --cert certs/alice-cert.pem --key certs/alice-key.pem logs eabc5579-7f8e-48d5-ba57-dc6e17f7a3ad


logs : hello world
```
