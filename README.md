# Alarm Service

An in-memory alarm management service built in Golang that provides RESTful APIs to create, update, and manage alarms with notification capabilities.

## Project Structure

```
alarm-service
├─ .http
│   └─ alarm_service.http
├─ cmd
│   └─ main.go
├─ internal
│   ├─ handlers
│   │   ├─ handlers_test.go
│   │   └─ handlers.go
│   ├─ models
│   │   ├─ alarm_test.go
│   │   └─ alarm.go
│   └─ services
│       ├─ alarm_service_test.go
│       └─ alarm_service.go
├─ testdata
│   └─ sample_alarms.json
├─ go.mod
├─ go.sum
└─ README.md
```

---

## Important Notes

Currently, this repository is **private** and access is restricted. Since the repo is private, the code has been downloaded as a ZIP file and shared via email.

If you require repository access, please provide your email ID so that I can add you as a contributor for direct cloning and evaluation.

---

## Prerequisites

- **Golang 1.23.4** (Recommended)
- **Lower versions of Golang** may also be compatible.

---

## Installation

1. Download the provided `.zip` file and extract it.
2. Navigate to the project directory:

```sh
cd alarm-service
```

3. Install dependencies:

```sh
go mod tidy
```

---

## Usage

### Run the Service

```sh
go run cmd/main.go
```

The service will start and listen on `http://localhost:8080`

### Sample HTTP Requests

Using `.http` file (Recommended for VSCode REST Client Plugin):

1. Open the file `alarm_service.http`.
2. Send the predefined requests to test the endpoints.

Alternatively, use `curl`:

**Create Alarm:**

```sh
curl -X POST -H "Content-Type: application/json" -d '{
  "name": "CPU Overload",
  "state": "Triggered"
}' http://localhost:8080/alarm
```

**Get All Alarms:**

```sh
curl -X GET http://localhost:8080/alarms
```

**Get Alarm By ID:**

```sh
curl -X GET http://localhost:8080/alarm?id={alarm_id}
```

**Update Alarm State:**

```sh
curl -X PUT -H "Content-Type: application/json" -d '{"state": "ACKed"}' http://localhost:8080/alarm?id={alarm_id}
```

**Delete Alarm:**

```sh
curl -X DELETE http://localhost:8080/alarm?id={alarm_id}
```

---

## Testing

1. Run tests with coverage:

```sh
go test ./... -coverprofile=coverage.out
```

2. Generate coverage report and display in the browser:

```sh
go tool cover -html=coverage.out
```

This command will open a detailed coverage report in your default browser.

3. To run specific tests:

```sh
go test ./internal/handlers -v
```

---

## Sample Data for Evaluation

Sample alarms data is available in `testdata/sample_alarms.json` for testing bulk creation and other operations.

---

## Key Features

- **Thread-safe Alarm Management:** Ensures concurrency safety using `sync.RWMutex`.
- **In-memory Storage:** Alarms are stored in-memory for simplicity and faster operations.
- **Notification Support:** Automatically sends notifications based on state transitions.
- **Bulk Creation Support:** Efficiently creates multiple alarms in one request.
- **Flexible REST API Design:** Easy integration with third-party services.

---

## Contribution

Contributions are welcome! Follow the steps below to contribute:

1. Fork the repository.
2. Create a feature branch.
3. Commit your changes with meaningful messages.
4. Raise a pull request for review.

---

## License

This project is licensed under the [MIT License](LICENSE).
