# imposter-db

## Description

This program acts as a proxy database for a remote MySQL server, in which a single table can be spoofed, allowing developers the ability to work without conflicting with other developers or testers. The proxy allows the developer to customize the "spoofed" table locally, while being connected to the remote database. The power of this can be seen when connecting an API or DBMS to the proxy. Physically, there are two different databases (the local and the remote), but the API or DBMS will only be aware of the proxy database.

Once you specify the table you want to spoof, it will be replicated inside the running Docker container. When you connect to the proxy, you will see all the tables from the remote database, but the spoofed table will be coming from the Docker container. This means you can change this table without affecting other people who are using the remote database.

> [!CAUTION]
> This project is to be best used on a v8 server. This should work on anything down to v5 server or client though, please open an issue if you encounter any problems.

> [!NOTE]
> I am working on coding the `caching_sha2_password` and `sha256_password` authentication methods. Currently, only the fast_auth path and `mysql_native_password` work. 

## Motivation

Have you ever been in development where the needed data is in the test environment, but working with it is almost impossible because application gates are constantly being opened and closed? With this program, you can spoof the application gate table and connect to the proxy. This allows you to easily modify the application gates locally without affecting the real gates in the test environment.

## Getting started

To begin, you must have Go and Docker installed.

### Option 1: Getting Started with a Dummy Remote (For Testing)

If you want to try out the proxy with a dummy remote database first, start the Docker service to launch a local MySQL instance that simulates a remote database:
```
docker compose up
```

This sets up two local MySQL instances (one acting as "remote" and one as "local") that you can use to familiarize yourself with the proxy functionality without connecting to a real remote database.

### Prepare the Proxy

In a different terminal, install dependencies and build the binary:
```
go mod download  # to install all dependencies

go build  # creates binary
```

Continue by selecting the schema and table to be spoofed, as the program is interactive. After this, the proxy will begin running. The idea is that you connect your DBMS and your locally running APIs to this proxy, so that you can modify the locally spoofed table, without changing configurations that impact others, and such that others cannot impact you.

The program will prompt you to choose a database:
```
Choose a database:
> TEST-DB
```
Use the **spacebar** to select and **Enter** to continue. Then select a table:
```
Select a table:
> application_gates
> application
> user
> user_types
```
After making your selection, **a replica table will be created inside a replica database in the running Docker container**.

For example, if you replicate the `application_gates` table, this would allow you to locally specify your own timeline for these gates, which is very powerful when trying to develop with data that may only exist in a test environment.

The proxy is now running! You can connect to it using the credentials needed to access the remote database:
```
host:     127.0.0.1
port:     3308
username: USER  (the user of the remote database)
password: PASS  (the password of the remote database user)
```

If the interactive selection process is too cumbersome, you can also use optional command-line flags:

- `-fk` - indicates whether tables with a foreign key reference to the identified table should be replicated
- `-schema=NAME` - name of the schema/database to use
- `-table=NAME` - name of the table to spoof

You can connect to the proxy from a DBMS (like MySQL Workbench or DBeaver) or configure a locally running API to use it as the database connection.

### Option 2: Running with a Real Remote Database

When you're ready to work with an actual remote database (instead of the dummy setup above), follow these steps:

1. **Configure your connection details** in `.env.local`, which specifies all the required variables. You will most likely only need to modify the remote database variables (REMOTE_DB_HOST, REMOTE_DB_PORT, REMOTE_DB_USER, REMOTE_DB_PASS, REMOTE_DB_NAME)

2. **Start the local database container** for the proxy to use:
```
docker compose up localdb
```

This starts only the `localdb` container, which will store the spoofed tables locally while the proxy connects to your actual remote database.

3. **Source the environment file and run the proxy**:
```bash
source .env.local
./imposter-db [-fk] [-schema=NAME] [-table=NAME]
```

The proxy will then interactively prompt you to select which schema and table to spoof (unless you provided the `-schema` and `-table` flags), replicate the table to the local database, and start listening on the configured PROXY_PORT (default: 3308).

## Architecture

Here's a flow chart depicting the architecture of what the proxy does:

<img width="512" height="512" alt="image" src="https://github.com/user-attachments/assets/11d7c52e-93cb-48f5-ad02-15a3fcce05dc" />

## Authors

- Matt Maloney : matttm
