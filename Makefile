.PHONY: run-broker run-sender run-receiver all clean requirements venv go-requirements cs-requirements python-requirements proto

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

# Install all requirements
requirements: venv go-requirements cs-requirements python-requirements

# Create a virtual environment for Python
venv:
	@echo "Creating Python virtual environment..."
	@python3 -m venv .venv
	@echo "Python virtual environment created at .venv."

# Install Go dependencies
go-requirements:
	@echo "Installing Go dependencies..."
	@cd broker && go mod tidy
	@cd broker && go get google.golang.org/grpc
	@cd broker && go get google.golang.org/protobuf
	@cd broker && go get google.golang.org/protobuf/cmd/protoc-gen-go
	@cd broker && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@cd broker && go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@cd broker && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Go dependencies installed."

# Install C# dependencies
cs-requirements:
	@echo "Restoring C# dependencies..."
	@cd sender && dotnet restore && dotnet build
	@echo "C# dependencies restored."

# Install Python dependencies within the virtual environment
# Possible platform problem with the source (only linux)
python-requirements: venv
	@echo "Installing Python dependencies in the virtual environment..."
	source .venv/bin/activate && pip install -r receiver/requirements.txt
	@echo "Python dependencies installed."

# Compile protobuf files (the go part(after export) is linux only)
proto:
	@echo "Compiling protobuf files..."
	protoc --csharp_out=sender proto/message_broker.proto
	protoc --python_out=receiver proto/message_broker.proto
	@export PATH=$PATH:$(go env GOPATH)/bin
	@source ~/.bashrc
	protoc --go_out=broker --go-grpc_out=broker proto/message_broker.proto
	@echo "Protobuf files compiled."
