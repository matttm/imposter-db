# imposter-db

## Description

This program acts as a proxy database, in which, any one table can be spoofed, allowing developers the ability to work without conflicting with other developers or testers.

> [!NOTE]
> This proxy has been written, with consideration of most cient versiions in mind. With that being said though, most testing has been with clients supporting at least the newer 4.1 version protocol.
>
> if you do find an issue, though, please document it well and raise an issue.

## Motivation.

Have you ever been in development and the needed data is in the test environment, but working with it is almost impossible because application gates are always being opened and closed? With this program, just spoof the application gate table and connect to the proxy. You can easily modify the application gates without affecting the real gates in the test environment.

## Getting started

To begin, you must have Go and Docker installed. In case you are on mac, docker won't work by itself, unless you install docker desktop. If you're like me and PREFER CLIs, then install, `colima` which is a container runtime, and can serve as a "docker daemon".

First, you need to run the docker daemon, which you can do if using Colima by running:
```
colima start
```
Then, you'll need to build the docker image and then start a container:
```
docker build -t imposter-img .

docker run -d --name imposter-cont -p 3306:3306 imposter-img
```

Then run:
```
go mod tidy  # to install all dependencies

go build  # creates binary
```
Before running the program, you need to export the following variables into the environment of the terminal that will be running the proxy:
```
export DB_HOST=""
export DB_USER=""
export DB_PASS=""
export DB_PORT=""
export DB_NAME=""
```
These variables should be the information you normally use to correct directly to the remote

Then run it with:
```
./imposter-db
```

Continue by selecting the schema and table to be spoofed, as the program is interactive. After this, the proxy will begin running. The idea is that you connect your DBMS and your locally running APIs to this proxy, so that you can modify the locally spoofed table, without changing configurations that impact others, and such that others cannot impact you.

Finally, the proxy is running! Now we want it to do some tom-foolery. We can connect to it using the following credentials:
```
host -  127.0.0.1
port - 3307
username - USER -- where USER is the user of the remote database?
password - PASS -- where PASS is the password of the above user in the remote database
```


## Authors

- Matt Maloney : matttm
