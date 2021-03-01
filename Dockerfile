FROM golang:alpine AS build-env
RUN apk --no-cache add build-base gcc
ADD . /src
RUN cd /src && go build -o auth-svc

# final stage
FROM alpine
WORKDIR /app
ENV ADMIN_PASS=admin
ENV REDIS_HOST=redis
ENV REDIS_PORT=6379
ENV REDIS_PASSWORD=
ENV ACCESS_SECRET=9ti4gj2dgfddrad3llr9
ENV REFRESH_SECRET=37fh79fjw955wdt321a9
ENV AESKEY=41a93e19086cbb496a6e1b728bd994a9d2070064de899b85e113868f50996de0
ENV ALLOWED=http://localhost:8888
ENV PORT=4000
ENV GIN_MODE=debug
COPY --from=build-env /src/auth-svc /app/
EXPOSE 4000
CMD ["./auth-svc"]
