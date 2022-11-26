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

- PUT `/`

  Body - JSON Object with the messages, compulsory publisher ID token
    ```json
    {
      "Key": <topic name: str>
      "Value": <Message: str>
    }
    ```

### Consuming

- GET `/<topic: str>`
  Returns the newest messages on that topic

- GET `/history/<topic: str>`
  Returns all the known messages on that topic

## Architecture
![](./.github/arch_diagram.png)

## Components
A slightly more detailed explanation of each of the components of the project.

### Gates
Gates is the python library that interacts with the Jumpnet Mux to send and consume messages.

### Mux
Mux interacts with the gates to transmit messages from producers to consumers.
