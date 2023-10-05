FROM golang:1.21.1 AS base
FROM base AS dev

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

ENV CGO_ENABLED=0 GOOS=linux GOPROXY=direct GOFLAGS=-buildvcs=false POSTGRES_HOST=postgres_db POSTGRES_USER=todoapp POSTGRES_PASSWORD=todoapp123 POSTGRES_DB=todo POSTGRES_PORT=5432

WORKDIR /opt/app/server
ENTRYPOINT [ "air" ]

#WORKDIR /go/src/app
#COPY . .

#RUN go build -v -o app .

#FROM scratch
#COPY --from=build /go/src/app/app /go/bin/app
#ENTRYPOINT ["/go/bin/app"]
