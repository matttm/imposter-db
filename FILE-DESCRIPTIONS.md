

├── Dockerfile -- Defines the steps to build a Docker image for the empty spoofed db
├── README.md
├── database.go -- Contains the core logic for interacting with the databases during init
├── dbeaver-init-success-no-proxy.pcapng
├── failed-codings
│   ├── ...
├── go.mod -- Defines the Go module path and lists project dependencies.
├── go.sum -- Contains cryptographic hashes of the project's dependencies for verification.
├── imposter-db -- The compiled executable for the imposter-db application.
├── login-dbeaver-mysql.pcapng -- Records network traffic of a successful login to a MySQL database using DBeaver.
├── main.go -- The entry point of the imposter-db application.
├── protocol
│   ├── auth-switch-request.go -- Handles the structure and logic for authentication switch requests.
│   ├── auth-switch-response.go -- Handles the structure and logic for authentication switch responses.
│   ├── auth.go -- Contains general logic and interfaces related to authentication.
│   ├── commands.go -- Defines constants and structures for database protocol commands.
│   ├── flags.go -- Defines constants for flags and options used in protocol messages.
│   ├── handshake-request.go -- Handles the structure and logic for the initial client handshake request.
│   ├── handshake-request_test.go -- Contains unit tests for the handshake request handling.
│   ├── handshake-response.go -- Handles the structure and logic for the server's handshake response.
│   ├── handshake-response_test.go -- Contains unit tests for the handshake response handling.
│   ├── ok-packet.go -- Handles the structure and logic for "OK" packets indicating success.
│   ├── ok-packet_test.go
│   ├── packet.go -- Defines the base structure and common functions for protocol packets.
│   ├── packet_test.go
│   ├── proxy.go -- contains logic for proxying database connections and traffic.
│   ├── query.go -- Handles the structure and logic for database query commands.
│   ├── resultset.go -- contains unused struct, for now...
│   ├── sql.go -- contains main logic for connection/command phase
│   └── utilities.go -- Contains general utility functions for the protocol package.
├── query.go -- contains functions returning strings of query
├── sql-manipulator.go -- Defines functionality for creating the INSERT SQL query
├── sql-manipulator_test.go.ignore
├── update-dbeaver-utan5.pcapng
└── view.go -- Contains logix for running the TUI
