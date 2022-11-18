# Jumpnet

Jumpnet is a kafka clone written in Golang and Python.

## Project Structure

```sh
├── bootstrap.sh        | Script to start the entire cluster
├── examples            | A full-fledged example of both a producer and a consumer talking to each other via the jumpnet using our interface libraries
├── gates               | Jumpnet Gates (Python) is a library to interact with the jumpnet 
└── mux                 | Mux is what makes Jumpnet work. It's the broker set and the broker manager.
```

## Architecture
![](./.github/arch_diagram.png)

## Components
A slightly more detailed explanation of each of the components of the project.

### Gates
Gates is the python library that interacts with the Jumpnet Mux to send and consume messages.

### Mux
Mux interacts with the gates to transmit messages from producers to consumers.
