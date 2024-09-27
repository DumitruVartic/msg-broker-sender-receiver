.PHONY: run-broker run-sender run-receiver all clean

# Default target
all: run-broker run-sender run-receiver

run-broker:
	go run broker/main.go

run-sender:
	cd sender && dotnet run

run-receiver:
	python3 receiver/receiver.py

# Clean build artifacts (if applicable)
clean:
	cd sender && dotnet clean