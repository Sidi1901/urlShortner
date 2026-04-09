
# URL Shortener Service (Golang)

A scalable URL shortening service built with Go, Redis, and PostgreSQL supporting custom aliases and analytics.


## Services
-Create Short URL
Purpose: Generate a short URL

-Redirect to Original URL
Purpose: Core redirect logic

-Get URL Info (Analytics-lite) 
Purpose:Fetch metadata (original URL, created_at, expiry, etc.)

-Delete Short URL
Purpose: Remove mapping

-Analytics
Purpose: Return: Click count, Last accessed, IP/device info (if tracked)

-Update URL
Purpose: Change expiry, Update original URL (optional)

-Auth
Purpose: Multiuse systemAuth

-Health Endpoint
Purpose: Used by deployment system to check health of instances


## Middlewawre
-Rate Limiting is applied on,
1) Create Short URL

-Logging
Globally


## Non Functional requiements
-Minimum redirect latency
-1B lifetime URLs

## Decision
- 62*62*62*62*62*62 unique short url
-Retry logic for collisions 


## Tech Stack

- Language: Go (Golang)
- Framework: Gin
- Database: PostgreSQL
- Cache: Redis
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
- email of type Email. And get email fro jw not from header

pay?
Increase quota



My steps

1) how to send cfg effectively following dependency injection.


Repository package is following DIP principle design.
For example sqlx concerte class is implementin user reposioty interface.
If later it's sql conconcrete class need, we can add one one concrete class implementing user repository interface. 

It's also follwing SIP and IP, as interface is segrated to allow needed and related functioanilities.


Analytics

1. Performance Metrics (System-level)

These tell you how fast and scalable your system is

Core metrics:
Latency (p50, p95, p99) :Redirect endpoint (GET /:shortcode) → most critical, URL creation endpoint  
Throughput : (RPS) Requests per second for redirects
Error rate : % of failed redirects or DB failures


2. Infrastructure Metrics

These tell you where the bottleneck is

DB query latency (Postgres)
Redis latency (once added)
CPU / Memory usage
Connection pool usage (very important in Go)


3. Business / Product Analytics (Core for your app)

This is what makes your project real-world ready

For each short URL:
Total clicks
Unique visitors
Clicks over time (time-series)
Expiry-based usage

Request metadata:
IP address (or hashed)
Location (GeoIP)
Device / browser (User-Agent parsing)
Referrer (who sent traffic)



1 & 2. Performance and Infrastructure Metrics (The Prometheus Way)
For Latency, Throughput, Error Rates, and Resource Usage (CPU/RAM/DB Pools), you should use Prometheus.

Why: These are "numerical counters" and "histograms." Prometheus is designed to store these efficiently and allows Grafana to calculate rates (like RPS) or percentiles (p95 latency) instantly.

Implementation in Go: You don't send this to Kafka. Instead, you use the Prometheus Go client to expose a /metrics endpoint in your Gin/Standard library app.

Database Metrics: For Postgres, you run a sidecar container called postgres_exporter. It scrapes Postgres and sends those metrics to Prometheus automatically.

3. Business & Product Analytics (The Kafka Way)
For Clicks, GeoIP, User-Agent, and Referrer data, you should use Kafka as your transport layer.

Why: This is "event-level" data. You need the raw details (like the IP and User-Agent) to perform high-cardinality analysis (e.g., "Which specific URLs are trending in Bangalore?").

The Flow:

Producer: Your Go app sends a JSON "Click Event" to a Kafka topic every time a redirect happens.

Enrichment: You can have a small Go consumer that reads the message, looks up the GeoIP based on the IP address, and parses the User-Agent string.

Storage: The enriched data is saved to a "Columnar Database" like ClickHouse or a standard PostgreSQL table.

Visualization: You connect Grafana to that database using a SQL datasource to build your "Clicks over time" or "Location" heatmaps.
