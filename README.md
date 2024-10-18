# msg-broker-sender-receiver

An implementation of a message broker with sender and receiver, developed as part of the Distributed Applications Programming course.

## Prerequisites

- Go (version 1.23.1 or later)
- .NET 8.0
- Python 3.x

## Makefile Targets

The project includes a Makefile for convenient operations. Here are the available targets:

- **run-broker**: Start the message broker.
- **run-sender**: Start the sender application.
- **run-receiver**: Start the receiver application.
- **clean**: Clean build artifacts (specifically for C#).
- **requirements**: Install all required dependencies.
- **proto**: Compile the protobuf files.

## Installing Requirements

To install all necessary requirements for each component, run (linux specific):

```bash
make requirements
```

Compiling Protobuf Files
To compile the protobuf files, you can use:

```bash
make proto
```

This target compiles the protobuf files for Go, C#, and Python.
