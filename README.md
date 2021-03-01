# Authentication &  Authorization Service

## Description

Authentication and Role handling Service written in GO 1.15.x that supports 384Bit signed JWT for login , token refresh and logout use cases.
Login and Refresh endpoints sets an http-only AES encrypted cookie with the jwt and refresh token and returns the JWT (10 min expiration) to the client.
The http-only encrypted cookie expires after 48 hours. A second NON http-only non encrtpyed cookie is sent with displayname and role to be consumed by a front end application. This cookie is not usable for authentication or accessing endpoints. This service also provides endpoints for managing application roles that could be assigned to accounts. Toggling of account status and setting passwords are also supported.

## Primary Use Case

### Token based (signed JWT), secure authentication support for web applications

### All endpoints account and role endpoints require a token

## Provides

```console
POST   /api/v1/login             --> create auth token
POST   /api/v1/logout            --> expire auth token
POST   /api/v1/refresh           --> refresh auth token
POST   /api/v1/account/new       --> create new account
POST   /api/v1/account/update    --> update account
POST   /api/v1/account/remove    --> remove account
POST   /api/v1/account/list      --> list all accounts
POST   /api/v1/account/find      --> find account
POST   /api/v1/account/toggle    --> toggle account status (active : inactive )
POST   /api/v1/account/set       --> set password for account
POST   /api/v1/role/new          --> create a new role
POST   /api/v1/role/update       --> update role
POST   /api/v1/role/remove       --> remove role
POST   /api/v1/role/list         --> list current roles
```

- Login
  - /api/v1/login

```json
POST
  {
    username: "someUser",
    password: "somePassword"
  }
    200 RESPONSE 
      COOKIE {ENCRYPTED HTTP-Only cookie}
      COOKIE {display_name, role}
    TEXT 
      "successful"
```

- Refresh
  - /api/v1/refresh

```json
 POST COOKIE {ENCRYPTED HTTP-Only cookie}
 200 RESPONSE
     COOKIE {ENCRYPTED HTTP-Only cookie}
     COOKIE {display_name, role}
     TEXT 
      "successful"
```

- Logout
  - /api/v1/logout

```console
  POST {HTTP-Only cookie from refresh or login (ENCRYPTED)}
 ```

## Accounts Management Endpoints

### New Account

- /api/v1/accounts/new

```console
  POST
  {
    "username": "testuser2",
    "displayname": "test s. user2",
    "email": "test.user2@testingcompany.com",
    "role": "report-user",
    "password": "supersecret01"
  }
```

### list, update, find (one), remove account

- /api/v1/accounts/list
- /api/v1/accounts/update
- /api/v1/accounts/find
- /api/v1/accounts/remove
- /api/v1/accounts/toggle
- /api/v1/accounts/set

```console
  POST
  {   
    "id": 3,
    "username": "testuser",
    "password": "somesecret",
    "email": "test.user@testingcompany.com",
    "displayname": "test s. user",
    "role": "report-user"
  }
```

## Role Management Endpoints

### list, update, new, remove role

- /api/v1/role/new
- /api/v1/role/update
- /api/v1/role/remove
- /api/v1/role/list

```console
  POST
  {   
    "id": 3,
    "name": "default-user",
    "displayname": "Default User Role",
    "description": "Default user Role for application"
  }
```

## Token & Claims

- example login / refresh token

```console
"eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjE4ZWYyM2Y5LTQ4YTYtNGE5My1hZWE4LWY1MDZlN2NlN2JhMCIsImRpc3BsYXlfbmFtZSI6IlRlbXBlci1TdXJlIEFkbWluIiwiZW1haWwiOiJhZG1pbkB0ZW1wZXItc3VyZS5jb20iLCJleHAiOjE2MTIyODg0MTYsInJvbGUiOiJhZG1pbiIsInVzZXJfaWQiOiIxIn0.j0nle36e2yFv5qvZMxJFewZ41d4zczE5UnHpC5s1T0PxTF5UK1FQT0zSsnZpwjCR"
```

## JWT Claims given

```json
{
  "access_uuid": "18ef23f9-48a6-4a93-aea8-f506e7ce7ba0",
  "display_name": "John Smith",
  "email": "jsmith@example.com",
  "exp": 1612288416,
  "role": "reports-user",
  "user_id": "24"
}
```

## Requirements

### Environmentals required for JWT Signature Validation

- ACCESS_SECRET  {see Dockerfile}
- REFRESH_SECRET {see Dockerfile}
- .env can be used for development & testing

```console
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
ACCESS_SECRET=1234567890abcdefghij
REFRESH_SECRET=9876543210abcdefghij
PORT=4000
```

### Redis with port 6379 exposed

- Used for storing the access_token and refresh_tokens until logout or token expiration
- Redis authentication supported

### Redis Data Structure

  ```go
  type TokenDetails struct {
      AccessToken  string
      RefreshToken string
      TokenUuid    string
      RefreshUuid  string
      AtExpires    int64
      RtExpires    int64
}
```

### Dockerfile example

```docker
FROM golang:latest
WORKDIR /go/src/app
ENV REDIS_HOST=localhost
ENV REDIS_PORT=6379
ENV REDIS_PASSWORD=
ENV ACCESS_SECRET=1234567890abcdefghij
ENV REFRESH_SECRET=9876543210abcdefghij
ENV PORT=4000
ENV GIN_MODE=release
COPY auth-svc .
EXPOSE 4000
CMD ["./auth-svc"]
```

## TODO

- Simple UI for managing accounts
- Since when you start the container / service, it will gen a new encryption key. it may be better for production deployments that scale to externalize this

## BUILDING / DEPLOYING

```console
go build -o my-auth-service -ldflags "-s -w" 
./my-auth-service
```

## DOCKER
- and Dockerfile and a docker-compose script have been provided for testing the auth_svc and the ui also in the repo

## REFERENCES and CODE INSPIRATION

- <https://github.com/victorsteven/jwt-best-practices>
- <https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes/>
