#!/bin/bash

# sudo docker-compose -f docker-compose.yml up      # Alternative using docker-compose

sudo docker pull postgres
sudo docker container run -d --name=pg-forum -p 5400:5432 -e POSTGRES_PASSWORD=postgres -e PGDATA=/pgdata_docker -v /pgdata_docker:/pgdata_docker postgres # This will map port 5400 host to postgres port (5432)
PGPASSWORD='postgres' psql -h localhost -p 5400 -U postgres -c "create database forum_db"

# sudo docker exec -it pg-forum bash
# sudo docker exec -it pg-forum psql -U postgres
