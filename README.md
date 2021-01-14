# Golang Forum API

API that already provide basic features of the forum web application and implemented in Golang language. This project is sample of production ready API that already has JWT auth system, logger, test, and also openAPI/swagger docs.

## Project Structure

-------
Here is summary folder structure in the project and It's short description
- **api** : Related to source code of the API itself
  - **auth** :
  - **controllers** :
  - **database** :
  - **logger** :
  - **middlewares** :
  - **models** :
  - **routers** :
  - **utils** :
  - <b>*server.go*</b> : go code that become gateway (startup) of the API 
- **docs** : Related to the documentation of the project, especially of the API
  - **assets** :
  - **swaggerui** :
- **logs** : Contains the log of the 
- **migrations** : Contains file (go source code) that needed to migrate model to the DB
- **seeders** :
- <b>*main.go*</b> : go code that become gateway of the project application
- <b>*.env*</b> :
- <b>* *.sh*</b> :
  
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

So, thre relation basically

1. One-to-Many between User and Thread
   -> threads.creator_id referencing users.id
2. One-to-Many between User and Post
   -> posts.author_id referencing users.id
3. One-to-Many between Thread and Thread
   -> posts.thread_id referencing threads.id

## API Endpoints

-------
More details of the API endpoint you could see in the swagger docs, but here are summary of the available endpoints:

1. User
   - `/signin` [post]
   - `/signup` [post]
   - `/users` [get]
   - `/users` [patch]
   - `/users/{id}` [get]
   - `/users/{id}` [delete]
2. Thread
   - `/threads` [get]
   - `/threads` [post]
   - `/threads` [patch]
   - `/threads/{id}` [get]
   - `/threads/{id}` [delete]
3. Post
   - `/posts` [get]
   - `/posts` [post]
   - `/posts` [patch]
   - `/posts/{id}` [delete]


## Application Dependency

-------

1. Go, here I'm using v1.13
2. PostgreSQL
3. Linux OS (optional), recommended using Ubuntu-based distros since I'm using it.
4. Docker (optional), if you want to containerize the app and the db 
5. Postman (optional), to use or test the API with forging requests  directly, and could be used to open [postman_collection.json](docs/forum_api.postman_collection.json) that include request sample of every endpoints.
6. Browser (optional), This needed to see swaggerui static file from the fileserver.

## Go Library/Module Used

-------

1. 

## How to Install 

-------

## How to Run

-------

## Issues

-------