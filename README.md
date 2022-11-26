# Jumpnet

Jumpnet is a kafka clone written in Golang and Python.

## Project Structure

```
├── bootstrap.sh        | Script to start the entire cluster
├── examples            | A full-fledged example of both a producer and a consumer talking to each other via the jumpnet using our interface libraries
├── gates               | Jumpnet Gates (Python) is a library to interact with the jumpnet 
└── mux                 | Mux is what makes Jumpnet work. It's the broker set and the broker manager.
```

## Broker Interface
This is the web interface of the project.

We assume that every broker has an independent IP and port, so this interface is designed to work on each and every broker, regardless of who the leader is.

- These interfaces work for _every_ node in the broker cluster. 
- Every node supports publish and subscribe operations

### Publishing

- PUT `localhost:9090X/`

  Body - JSON Object with the messages, compulsory publisher ID token
    ```json
    {
      "Key": <topic name: str>
      "Value": <Message: str>
    }
    ```

### Consuming

- GET `localhost:909X/<topic: str>`
  Returns the newest messages on that topic

- GET `localhost:909X/history/<topic: str>`
  Returns all the known messages on that topic

### Example

Here `9091` is the current broker. These commands will work with _any_ broker instance, so examples will need to change the port number accordingly.

#### Publishing

```sh
$ curl -L http://127.0.0.1:9091/ -X PUT -d '{"Key":"<KEY>", "Value":"<VALUE>"}'
```

#### Subscribing

```sh
$ curl -L http://127.0.0.1:9091/<KEY>
```


### Running the Brokers

```sh
cd mux
go build main.go
```

You'll have a `main` binary in the `mux` folder.

1. Start the first node in the Raft Cluster

  ```
  ./main -state_dir=./tempdata/1 -raft :8080 -api :9090
  ```
  
2. Start Additional Nodes to join the cluster, **_in separate shells_**

  ```
  ./main -state_dir ./tempdata/2 -raft :8081 -api :9091 -join :8080
  ./main -state_dir ./tempdata/3 -raft :8082 -api :9092 -join :8080
  ```


## Architecture
Jumpnet's `mux` uses Raft as the replication layer, implementing a simple key-List[Value] store on top of Raft. The replicated log manages inserts, deletes, topic creation, and multiple other requirements.

![](./.github/arch_diagram.png)

## Components
A slightly more detailed explanation of each of the components of the project.

### Gates
Gates is the python library that interacts with the Jumpnet Mux to send and consume messages.

### Mux
Mux interacts with the gates to transmit messages from producers to consumers.
