# imposter-db

## Description

This program acts as a proxy database for MySQL, in which, any one table can be spoofed, allowing developers the ability to work without conflicting with other developers or testers.

Once you specified the table you want to spoof, it will be replicated inside the running, docker container, so once you connect to the proxy, you will see all the tables from the remote, but this spoofed table will be coming from the docker container, meaning that you can change this table without affecting other people that are using the remote database.

> [!CAUTION]
> This project is to be best used on a v8 server. This should work on anything down to v5 server or client though, please open an issue if you encounter any problems.

> [!NOTE]
> I am working on coding the `caching_sha2_password` and `sha256_password` authentication methods. Currently, only the fast_auth path and `mysql_native_password` work. 

## Motivation.

Have you ever been in development and the needed data is in the test environment, but working with it is almost impossible because application gates are always being opened and closed? With this program, just spoof the application gate table and connect to the proxy. You can easily modify the application gates without affecting the real gates in the test environment.

## Getting started

To begin, you must have Go and Docker installed.

Then start the docker service:
```
docker compose up
```
and in a different terminal, run:
```
go mod download  # to install all dependencies

go build  # creates binary
```
These variables should be in the environment of the terminal running the proxy and should be the information you normally use to correct directly to the remote:
```
export DB_HOST=""
export DB_USER=""
export DB_PASS=""
export DB_PORT=""
export DB_NAME=""
```

Then run it with:
```
./imposter-db [-fk] [-schema=NAME] [-table=NAME]
```
Continue by selecting the schema and table to be spoofed, as the program is interactive. After this, the proxy will begin running. The idea is that you connect your DBMS and your locally running APIs to this proxy, so that you can modify the locally spoofed table, without changing configurations that impact others, and such that others cannot impact you.

Chooose DB:
```
Choose a database:
> A
> B
> C
```
Using space for select and enter to continue, then select a table.
```
Select a table:
> D
> E
> F
```
After choosing those, **a replica table `D` will be made inside a replica database `A` in the running docker container**.

Finally, the proxy is running! Now we want it to do some tom-foolery. We can connect to it using the credentials needed to access the remote.
```
host -  127.0.0.1
port - 3307
username - USER -- where USER is the user of the remote database
password - PASS -- where PASS is the password of the above user in the remote database
```

If this selection process is too cumbersome, you can also take advantage of the optional flags:

- fk - indicates whether, the tables with a foreign key reference to the identified table, should be replicated
- schema - name of the schema
- table - name of the table

You can connect to it from a DBMS or you can set a local running API to use it as the database.

# Architecture

Here's a flow chart depicting the architecture of what the proxy does:
<img width="256" height="256" alt="image" src="https://github.com/user-attachments/assets/11d7c52e-93cb-48f5-ad02-15a3fcce05dc" />

## Authors

- Matt Maloney : matttm
