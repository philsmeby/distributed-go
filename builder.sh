#!/bin/bash

install_deps(){
    go get -u github.com/jinzhu/gorm
    go get -u github.com/streadway/amqp
    go get -u github.com/mattn/go-sqlite3
}

build_sensors() {
    cd $(pwd)/sensors
    GOOS=linux GOARCH=amd64 go build
    cd -
}

build_coordinator() {
    cd $(pwd)/coordinator/exec
    GOOS=linux GOARCH=amd64 go build -o coordinator
    cd -
}

build_metric_manager(){
    cd $(pwd)/monitoring/exec
    GOOS=linux GOARCH=amd64 go build -o metricmanager
    cd -
}

docker_compose(){
    docker-compose up -d --build
}

main (){
    install_deps
    build_sensors
    build_coordinator
    build_metric_manager

    docker_compose
}

main