# build stage
FROM golang:alpine AS build
WORKDIR /go/src/app

# extra stage
# RUN go mod init eclaim-api
RUN go mod tidy
# TODO: Why comment this, please find out
# RUN go get -u github.com/swaggo/swag/cmd/swag@v1.7.0
# RUN swag init
# RUN apk add --no-cache git
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/api ./main.go

# final stage
FROM scratch
# Import from builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# Use an unprivileged user.
USER appuser:appuser

WORKDIR /usr/app

COPY --from=build /go/src/app/bin /go/bin
COPY --from=build /go/src/app/ ./
EXPOSE 3000
ENTRYPOINT /go/bin/api