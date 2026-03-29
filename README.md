
# URL Shortener Service (Golang)

A scalable URL shortening service built with Go, Redis, and PostgreSQL supporting custom aliases and analytics.


## Features

- Create short URLs
    - User can create Short URL. 1:1 i.e One Short URL per ORiginal URL
- Custom aliases
- Expiry support
- Click tracking & analytics
- Rate limiting
- RESTful API

## Tech Stack

- Language: Go (Golang)
- Framework: Gin
- Database: PostgreSQL
- Cache: Redis
- Messaging (optional): Kafka
- Containerization: Docker

## Libraries


## Project Structure

.
├── cmd/            # Entry point
├── internal/
|   |__ configs/
│   ├── handlers/   # HTTP handlers
│   ├── services/   # Business logic
│   ├── repository/ # DB layer
│   └── models/     # Storage layer
|   |__DTO/         # API contract layer
├── pkg/            # Reusable utilities
├── Dockerfile
└── README.md

##signin/signup login page 

- User signup/login (JWT / OAuth)
- Each user manages their own links
- API keys per user

## Setup

1. Clone the repo
2. Create .env file
3. Run PostgreSQL & Redis

```bash
docker-compose up -d
go run cmd/main.go


👉 This shows you understand **developer experience (DX)**.

```


# 6. 🔑 Environment Variables


| Variable     | Description              |
|--------------|--------------------------|
| DB_URL       | PostgreSQL connection    |
| REDIS_URL    | Redis connection         |
| APP_PORT     | Server port              |


## API Endpoints

### Create Short URL
POST /shorten

### Get Original URL
GET /{shortCode}

### Analytics
GET /analytics/{shortCode}


## Design Decisions

- Used Redis for caching hot URLs to reduce DB load
- PostgreSQL for persistence
- Implemented rate limiting to prevent abuse
- Used base62 encoding for short URLs

## Contributing

Pull requests are welcome.


## Next step

- Use gorm
- redis based rate limiting
- User data save - dashboard
- Authentication/ Athorisation (observability for admins)
- Allow user to create many short URLs per Original URL

pay?
Increase quota

