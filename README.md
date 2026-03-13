# urlShortner
## Project


urlShortener
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│
│   ├── database/
│   │   └── postgres.go
│
│   ├── models/
│   │   └── url.go
│
│   ├── repository/
│   │   └── url_repository.go
│
│   ├── service/
│   │   └── url_service.go
│
│   ├── handler/
│   │   └── url_handler.go
│
│   └── router/
│       └── router.go
│
├── web/
│   └── templates/
│
├── go.mod
└── go.sum