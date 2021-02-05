# Token (Signed JWT) Authorization Service

## Description

Authentication Service written in GO 1.15.x that supports 384Bit signed JWT for login , token refresh and logout use cases.
Login and Refresh endpoints sets an http-only AES encrypted cookie with the jwt and refresh token and returns the JWT (10 min expiration) to the client.
The http-only encrypted cookie expires after 48 hours.

## Primary Use Case

  Token based (signed JWT), secure authentication support for web applications

## Provides

```console
POST   /api/v1/login             --> local/auth-svc/handler.(*profileHandler).SendLoginCookie-fm (4 handlers)
POST   /api/v1/logout            --> local/auth-svc/handler.(*profileHandler).LogoutSession-fm (4 handlers)
POST   /api/v1/refresh           --> local/auth-svc/handler.(*profileHandler).RefreshSession-fm (4 handlers)
POST   /api/v1/account/new       --> local/auth-svc/services.AddAccount (5 handlers)
POST   /api/v1/account/update    --> local/auth-svc/services.ModifyAccount (5 handlers)
POST   /api/v1/account/remove    --> local/auth-svc/services.RemoveAccount (5 handlers)
POST   /api/v1/account/list      --> local/auth-svc/services.ListAccounts (5 handlers)
POST   /api/v1/account/find      --> local/auth-svc/services.FindUser (5 handlers)
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
      JSON
        {
          "display_name": "John D Smith",
          "role": "admin"
        }  
```

- Refresh
  - /api/v1/refresh

```json
 POST COOKIE {ENCRYPTED HTTP-Only cookie}
 200 RESPONSE
     COOKIE {ENCRYPTED HTTP-Only cookie}
     JSON 
       {
        "display_name": "John D Smith",
        "role": "admin"
       }
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

```console
  POST
  {   
    "id": 3,
    "username": "testuser",
    "password": "",
    "email": "test.user@testingcompany.com",
    "displayname": "test s. user",
    "role": "report-user"
  }
```

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
