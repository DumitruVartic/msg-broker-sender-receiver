run-broker:
	go run broker/main.go

run-sender:
	cd sender && dotnet run sender.cs

run-receiver:
	python3 receiver/receiver.py
