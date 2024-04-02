# ACME Sky - Airline Service

This repo refers to the Airline service used by ACME Sky.

It is a REST API backend used by ACME Sky to find new flights and receive new
offerts.

## Build

> [!TIP]
> You can use `./build.sh` for a step-by-step setup guide for deploying.

You need to set up

```
POSTGRES_USER=user
POSTGRES_PASSWORD=pass
POSTGRES_DB=db
DATABASE_DSN="host=localhost user=user password=pass dbname=db port=5432"
JWT_TOKEN=t0k3n
DEBUG=0
```

and build

```
docker build -t acmesky-airlineservice-api .
```

after that you can put everything up

```
docker compose up
```

Now you have to create a new user:

```
$ docker ps                                                                                                                                                         1 ↵
CONTAINER ID   IMAGE                        COMMAND                  CREATED          STATUS          PORTS                                       NAMES
646fc7c342fd   acmesky-airlineservice-api   "./main"                 57 seconds ago   Up 56 seconds   0.0.0.0:8080->8080/tcp, :::8080->8080/tcp   airlineservice-api
c8f8f8782838   postgres:16-alpine           "docker-entrypoint.s…"   57 seconds ago   Up 56 seconds   5432/tcp                                    airlineservice-postgres
$ docker exec -it c8 psql -U acme -d db -W
Password:
psql (16.2)
Type "help" for help.

db=# create user
user              user mapping for
db=# insert into users (username, password) values ('sa', '6ea044c786f237c955b497b04b9247f2a663c5038e54175e62308c8b8457e23e');
INSERT 0 1
```

so to log in:

```
$ curl -X POST http://localhost:8080/v1/login/ -H 'content-type: application/json' -H 'accept: application/json, */*;q=0.5' -d '{"username":"sa","password":"*****"}'
HTTP/1.1 200 OK
Content-Length: 147
Content-Type: application/json; charset=utf-8
Date: Tue, 02 Apr 2024 10:14:41 GMT

{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXJuYW1lIiwiZXhwIjoxNzEyMDU2NDgxfQ.7R87BuuVkvOwojBpLmJ8OKtKC0B9Iq-wWSA_pqGBVXE"
}
```

now you have to use that token as `Authorization` header if you want to create
new data by API.
