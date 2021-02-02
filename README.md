# Token (Signed JWT) Authorization Service

## Description

Authentication Service written in GO 1.15.x that supports 384Bit signed JWT for login , token refresh and logout use cases.
Login and Refresh endpoints sets an http-only cookie with the jwt and refresh token and returns the JWT (10 min expiration) to the client.
The http-only cookie expires after 48 hours.

## Primary Use Case

  Token based (signed JWT), secure authentication support for web applications

## Provides

```javascript
POST   /api/v1/login             --> local/auth-svc/handler.(*profileHandler).SendLoginCookie-fm (4 handlers)
POST   /api/v1/logout            --> local/auth-svc/handler.(*profileHandler).LogoutSession-fm (4 handlers)
POST   /api/v1/refresh           --> local/auth-svc/handler.(*profileHandler).RefreshSession-fm (4 handlers)
```

- Login
  - /api/v1/login

   ```json
   POST
     {
      username: "someUser",
      password: "somePassword"
     }
     
   ```

- Refresh
  - /api/v1/refresh

  ```console
  POST {HTTP-Only cookie from login}
  ```

- Logout
  - /api/v1/logout

  ```console
  POST {HTTP-Only cookie from refresh or login}
  ```

- example login / refresh token

```console
"eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjE4ZWYyM2Y5LTQ4YTYtNGE5My1hZWE4LWY1MDZlN2NlN2JhMCIsImRpc3BsYXlfbmFtZSI6IlRlbXBlci1TdXJlIEFkbWluIiwiZW1haWwiOiJhZG1pbkB0ZW1wZXItc3VyZS5jb20iLCJleHAiOjE2MTIyODg0MTYsInJvbGUiOiJhZG1pbiIsInVzZXJfaWQiOiIxIn0.j0nle36e2yFv5qvZMxJFewZ41d4zczE5UnHpC5s1T0PxTF5UK1FQT0zSsnZpwjCR"
```

## JWT Claims given

```json

  "access_uuid": "18ef23f9-48a6-4a93-aea8-f506e7ce7ba0",
  "display_name": "John Smith",
  "email": "jsmith@example.com",
  "exp": 1612288416,
  "role": "reports-user",
  "user_id": "24"
}
```

## Requires

### Environmentals required for JWT Signature Validation

- ACCESS_SECRET  {see Dockerfile}
- REFRESH_SECRET {see Dockerfile}

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
ENV ACCESS_SECRET=9ti4gj2dgfddrad3llr9
ENV REFRESH_SECRET=37fh79fjw955wdt321a9
ENV PORT=4000
ENV GIN_MODE=release
COPY auth-svc .
EXPOSE 4000
CMD ["./auth-svc"]
```

## TODO

- Add backend accounts db with SQLlite as the storage engine
- Simple UI for managing accounts
