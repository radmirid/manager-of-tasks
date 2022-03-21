FROM golang:1.17-buster AS build

ENV GOPATH=/
WORKDIR /src/
COPY ./ /src/

RUN go mod download; CGO_ENABLED=0 go build -o /manager-of-tasks ./cmd/main.go


FROM alpine:latest

COPY --from=build /manager-of-tasks /manager-of-tasks
COPY ./configs/ /configs/
COPY ./psql.sh ./

RUN apk --no-cache add postgresql-client && chmod +x psql.sh

CMD ["/manager-of-tasks"]