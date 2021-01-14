# Golang Forum API

## Description

-------
A simple API that already provide basic features of the forum web application and implemented in Golang language. This project is sample of production ready API that already has JWT auth system, logger, test, and also openAPI/swagger docs.

## Project Structure

-------
Here is summary general structure (folder and important files) in the project and It's short description, where structure separation itself is mainly based on It's context and functionality.
- **api**: Stuff that related to source code of the API itself
  - **auth**: Stuff that specifically related to auth system
  - **controllers**: Functions that handling request as the end/terminal handler, so function here is main judge to giving response of the request.
  - **database**: Stuff that related to interact with DB 
  - **logger**: Stuff that related to writing log of the API
  - **middlewares**: Functions that placed as an intermediate component that passed by the request and response
  - **models**: Structure that used as the DB table's schema representation in Go, and also being model for request-response format and It's corresponding casting function/method
  - **routers**: Function that define endpoint path and It's pipeline or route or function chain that become request-response way
  - **utils**: Helper function to help implement other functionalities faster, easier, and less redundant
  - <b>*server.go*</b>: Go code that become gateway (startup) of the API 
- **docs**: Related to the documentation of the project, especially of the API
  - **assets**: Static file that being served as the content of the docs (such as images, etc)
  - **swaggerui**: Satic file that specifically for the file server of the swagger docs
- **logs**: Log file(s) of the API
- **migrations**: Go source code that needed to migrate model/schema to the DB
- **seeders**: Go source code that needed to seed/insert data/record to the DB
- <b>*main.go*</b>: Go code that become gateway of the project application
- <b>*.env*</b>: Environment variable that used to run this app
- <b>* *.sh*</b>: Shell console code that highly needed and packed to help operation of the project
  
<!-- ```
root
├── config
│   └── routes.js
├── screens
│   └── App
│       ├── screens
│       │   ├── Admin
│       │   │   ├── screens
│       │   │   │   ├── Reports
│       │   │   │   │   └── index.js
│       │   │   │   └── Users
│       │   │   │       └── index.js
│       │   │   └── index.js
│       │   └── Course
│       │       ├── screens
│       │       │   └── Assignments
│       │       │       └── index.js
│       │       └── index.js
│       └── index.js
└── index.js
``` -->

## Models

-------

![db_schema](docs/assets/db_schema.png)

Here is the details of model's schema that used in the DB 

1. User -> users
   - ID: uint32 -> id: bigint (bigserial)
   - Username: string -> username: varchar(255) unique not null
   - Email: string -> email: varchar(255) unique not null
   - Password: string -> password: varchar(255) not null
   - CreatedAt: time.Time -> created_at: timestamp
   - UpdatedAt: time.Time -> updated_at: timestamp
2. Thread -> threads
   - ID: uint64 -> id: bigint (bigserial)
   - Title: string -> title: varchar(255) not null
   - Topic: string -> topic: varchar(255) not null
   - CreatorID: uint32 -> creator_id: bigint
   - Creator: *User
   - Posts: []*Post
   - CreatedAt: time.Time -> created_at: timestamp
   - UpdatedAt: time.Time -> updated_at: timestamp
3. Post -> posts
   - ID: serial, uint64
   - AuthorID: uint64 -> author_id: bigint
   - Author: *User
   - ThreadID: uint64 -> thread_id: bigint
   - Thread: *Thread
   - Content: string -> topic: teet not null
   - CreatedAt: time.Time -> created_at: timestamp
   - UpdatedAt: time.Time -> updated_at: timestamp

And the relations basically are

1. One-to-Many between User and Thread
   -> threads.creator_id referencing users.id
2. One-to-Many between User and Post
   -> posts.author_id referencing users.id
3. One-to-Many between Thread and Thread
   -> posts.thread_id referencing threads.id

## API Endpoints

-------
More details of the API endpoint you could see in the swagger docs, but here are summary of the available endpoints:
1. SwaggerUI Docs in `/docs/` [get]
2. User
   - `/signin` [post]
   - `/signup` [post]
   - `/users` [get]
   - `/users` [patch]
   - `/users/{id}` [get]
   - `/users/{id}` [delete]
3. Thread
   - `/threads` [get]
   - `/threads` [post]
   - `/threads` [patch]
   - `/threads/{id}` [get]
   - `/threads/{id}` [delete]
4. Post
   - `/posts` [get]
   - `/posts` [post]
   - `/posts` [patch]
   - `/posts/{id}` [delete]


## Application Dependency

-------

1. Go, here I'm using v1.13
2. PostgreSQL
3. Linux OS (optional), recommended using Ubuntu-based distros since I'm using it (Ubuntu 20.04).
4. Docker (optional), if you want to containerize the app and the db 
5. Postman (optional), to use or test the API with forging requests  directly, and could be used to open [postman_collection.json](docs/forum_api.postman_collection.json) that include request sample of every endpoints.
6. Browser (optional), such as Firefox or Chrome, which is needed to see swaggerui static file from the fileserver.

## Go Library/Module Used

-------

1. `gorm.io/gorm` as Object Relational Mapping (ORM) library that I used in this project to access database, and `gorm.io/driver/postgres` as the mandatory driver that ORM need to access Postgre DB.
2. `github.com/gorilla/mux` as routing helper, and `github.com/gorilla/context` as helper to passing data through the middleware.
3. `github.com/sirupsen/logrus` as main library to help logging the API and `github.com/rifflock/lfshook` as hook library to write file to the file.
4. `github.com/asaskevich/govalidator` as helper library to help validate input.
5. `github.com/mikunalpha/goas` as helper to make swagger spec from the comments in the source code.
6. `github.com/joho/godotenv` to help fetch environment var from a file.
7. `github.com/dgrijalva/jwt-go` library that needed to make JWT-based auth.
8. `github.com/go-errors/errors` helps to checking or handling things that related to error type.


## How to Install 

-------
For run the API for the first time, make sure to do these steps:

1. Prepare application that needed, minimum Go and Postgres. You could use dockerized version API from create the docker image with the `Dockerfile`, and also create PostgreSQL instance in docker container with command in `init_db.sh`, or simply you could just run in terminal
   ```
      sh init_db.sh
   ```

2. Install every Go library listed in `go.mod` with below command
   ```
      go mod tidy
   ```
   or from `glide.yaml` file if you're using and familiar with glide.
3. Set the environment variable in the `.env` file to the value that approriate.
4. Migrate the models to the DB and try to seed the data if you think It's needed.

## How to Run

-------
Running the application is just run the `main.go` in the root of the project based on main function of the app that you want to use

1. API
   ```
      go run main.go
   ```
2. Migrate model
   ```
      go run main.go -- migrate
   ```
3. Seed Data 
   ```
      go run main.go -- seed
   ```

Or you could try to build the projct first and then try to execute it with corresponding schema that written above. For example you could run this command below
```
   go build -o main
   ./main   # or ./main -- test or ./main -- migrate
```

## How to Test

-------
To test the application just run code in the `run_tests.sh`, or simply you could just run in terminal
```
   sh run_tests.sh
```
And see if there is any fail on it.

## Issues

-------

1. Still can't log GORM query into a file.
2. Log message into terminal still has bad format (spacing).