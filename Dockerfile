# syntax=docker/dockerfile:1
FROM golang:latest as build

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . /app
RUN go build ./cmd/main.go

FROM ubuntu:latest

RUN apt-get -y update && apt-get install -y tzdata
ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER
USER postgres


RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER just4n4cc WITH SUPERUSER PASSWORD 'password';" &&\
    createdb -O just4n4cc dbproject &&\
    /etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

WORKDIR /app
COPY . .
COPY --from=build /app/main .

USER root
ENV PGPASSWORD password
CMD service postgresql start && psql -h localhost -d dbproject -U just4n4cc -p 5432 -a -q -f ./db/db.sql && ./main
