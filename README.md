# imposter-db

## Description

This program that acts as a proxy database, in which, any one table can be spoofed, allowing developers the ability to work without conflicting with other developers or testers.

## Motivation.

Have you ever been in development and the needed data is in the test environment, but working with it is almost impossible because application gates are always being opened and closed? With this program, just spoof the application gate table and connect to the proxy. You can easily modify the application gates without affecting the real gates in the test environment.

## Getting started

To begin, you must have Go and Docker installed. In case you are on mac, docker won't work by itself, unless you install docker desktop. If you're like me and PREFER CLIs, then install, `colima` which is a container runtime, and can serve as a "docker daemon".


First set these environment variables as to connect to your desired remote database
```
export DB_HOST=""
export DB_USER=""
export DB_PASS=""
export DB_NAME=""

```

Then you'll need to build the docker image and then start a container:
```
docker build -t imposter-img .

docker run -d --name imposter-cont -p 3306:3306 imposter-img
```

Then run:
```
go build
```

Continue by selecting the schema and table to be spoofed. After this, the proxy will begin running. The idea is that you connect your DBMS and your locally running APIs to this proxy, so that you can modify the locally spoofed table, without changing configurations that impact others, and such that others cannot impact you.

## Authors

- Matt Maloney : matttm
