# Twitter Demo

A Twitter demo application built with Go, using clean architecture with layer separation (Repository, Usecase, Controller).

## Twitter System Design Implementation

This repository contains the implementation of a simplified microblogging platform, designed with a focus on **horizontal scalability** and **read optimization** (High Read Throughput).

The project solves the technical challenge by proposing a distributed architecture that separates read and write responsibilities, using asynchronous patterns to ensure performance under high demand.

## Architecture and Design

The system is designed under the principle that in a social network like Twitter, **read operations massively outnumber write operations**.

![Architecture Diagram](./docs/scalable_architecture.png)

### Scalability Strategy
To meet the requirement of "scaling to millions of users," the following design decisions were made:

1.  **Separation of Concerns (CQRS - Concept):**
    * **Write API:** Handles content creation. Its priority is availability and fast data capture.
    * **Read API:** Handles timeline queries. Its priority is low latency.
    * **Independent Scaling:** By decoupling these APIs, we can scale the reading infrastructure (which receives 90% of the traffic) without over-provisioning the writing infrastructure.

2.  **Read Optimization (Fan-Out on Write):**
    * Instead of calculating the timeline in real-time (which would require costly database `JOINs` for every request), the system pre-calculates and "pushes" tweets to followers' timelines.
    * An **Asynchronous Fan-Out** pattern is used via message queues (Kafka/RabbitMQ) and Workers.

3.  **Hybrid Cache Strategy:**
    * **Push (Fan-Out):** When a tweet is created, the Worker immediately updates the cache (Redis) of active followers.
    * **Pull (Cache-Aside):** If the cache fails or is empty (Cold Start), the Read API queries the persistent database, rebuilds the timeline, and stores it in the cache for future queries.
    * **Consistency:** The system guarantees **Eventual Consistency**. There may be a slight delay between posting a tweet and its appearance on a follower's timeline, an acceptable trade-off to ensure system availability and speed.

### Data Flow
1.  **Write:** The Load Balancer directs the request to the **Write API**. The tweet is persisted in the Master DB, and a `tweet.created` event is published to the Message Broker.
2.  **Processing:** A **Consumer (Worker)** reads the event and distributes the tweet ID to the Redis timeline lists of the followers.
3.  **Read:** The Load Balancer directs the request to the **Read API**. It first queries Redis (O(1) access). If there is a "Cache Miss," it falls back to the Read Replica DB.

## Tech Stack

* **Language:** Go (Golang) 1.25
* **Relational Database:** PostgreSQL (Simulated Master-Slave configuration).
* **Cache / Key-Value Store:** Redis (Lists for timelines).
* **Message Broker:** Kafka (To decouple writing from processing).
* **Infrastructure:** Docker & Docker Compose.

## Getting Started

### Prerequisites

- Docker and Docker Compose installed
- Go 1.25+ (for local development)
- Make (optional, for using Makefile commands)

### Running the Project

1. **Clone the repository:**
```bash
git clone https://github.com/DunittMonagas/twitter-demo.git
cd twitter-demo
```

2. **Start all services with Docker Compose:**
```bash
docker-compose up --build
```

This will start:
- **Write API** on port `8081`
- **Read API** on port `8080`
- **Worker** (background processor)
- **PostgreSQL** on port `5432`
- **Redis** on port `6379`
- **Kafka** on ports `9092` (internal) and `9093` (external)
- **Kafka UI** on port `8090` (for monitoring)

3. **Wait for all services to be healthy:**
The first startup may take a few minutes as Docker downloads images and initializes the database.

4. **Verify the services are running:**
```bash
# Check Write API
curl http://localhost:8081/health

# Check Read API
curl http://localhost:8080/health
```

5. **Access Kafka UI (optional):**
Open your browser and navigate to `http://localhost:8090` to monitor Kafka topics and messages.

### Running Locally (Development)

If you prefer to run the services locally without Docker:

```bash
# Install dependencies
go mod download

# Run the Write API
go run cmd/write-api/main.go

# Run the Read API (in another terminal)
go run cmd/read-api/main.go

# Run the Worker (in another terminal)
go run cmd/worker/main.go
```

**Note:** You'll need PostgreSQL, Redis, and Kafka running locally and update the environment variables accordingly.

### Stopping the Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down -v
```

## API Examples

Once the services are running, you can interact with the APIs using the following cURL commands:

### User Operations (Write API - Port 8081)

**Create a new user:**
```bash
curl -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

**Update a user:**
```bash
curl -X PUT http://localhost:8081/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_updated",
    "email": "john.new@example.com",
    "password": "newsecurepassword"
  }'
```

### User Queries (Read API - Port 8080)

**Get all users:**
```bash
curl http://localhost:8080/users
```

**Get user by ID:**
```bash
curl http://localhost:8080/users/1
```

### Tweet Operations (Write API - Port 8081)

**Create a tweet:**
```bash
curl -X POST http://localhost:8081/tweets \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "content": "Hello, Twitter! This is my first tweet."
  }'
```

**Delete a tweet:**
```bash
curl -X DELETE http://localhost:8081/tweets/1
```

### Timeline Queries (Read API - Port 8080)

**Get user timeline (tweets from followed users):**
```bash
# Get timeline for user ID 1
curl http://localhost:8080/timeline/1

# With pagination
curl "http://localhost:8080/timeline/1?limit=10&offset=0"
```

**Get user's own tweets:**
```bash
curl http://localhost:8080/tweets/user/1
```

### Follow Operations (Write API - Port 8081)

**Follow a user:**
```bash
curl -X POST http://localhost:8081/followers \
  -H "Content-Type: application/json" \
  -d '{
    "follower_id": 1,
    "followed_id": 2
  }'
```

**Unfollow a user:**
```bash
curl -X DELETE http://localhost:8081/followers/1/2
```

### Follow Queries (Read API - Port 8080)

**Get user's followers:**
```bash
curl http://localhost:8080/followers/1
```

**Get users that a user is following:**
```bash
curl http://localhost:8080/following/1
```

### Example Workflow

Here's a complete example to test the entire flow:

```bash
# 1. Create two users
curl -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "email": "alice@example.com", "password": "pass123"}'

curl -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '{"username": "bob", "email": "bob@example.com", "password": "pass456"}'

# 2. Alice (ID: 1) follows Bob (ID: 2)
curl -X POST http://localhost:8081/followers \
  -H "Content-Type: application/json" \
  -d '{"follower_id": 1, "followed_id": 2}'

# 3. Bob creates a tweet
curl -X POST http://localhost:8081/tweets \
  -H "Content-Type: application/json" \
  -d '{"user_id": 2, "content": "Hello from Bob!"}'

# 4. Wait a moment for the worker to process the event (fan-out)
sleep 2

# 5. Get Alice's timeline (should contain Bob's tweet)
curl http://localhost:8080/timeline/1
```

## Code Structure

The project follows the guidelines of the [Standard Go Project Layout](https://github.com/golang-standards/project-layout) and applies **Clean Architecture** principles to isolate business logic from infrastructure.

```text
.
├── cmd/
│   ├── read-api/      # Entrypoint for the Read API
│   ├── write-api/     # Entrypoint for the Write API
│   └── worker/        # Entrypoint for the asynchronous processor
├── internal/
│   ├── domain/        # Pure entities (Enterprise Business Rules)
│   ├── usecase/       # Business logic (Application Business Rules)
│   ├── infrastructure/# Repository implementations (DB, Redis, Kafka)
│   └── interfaces/    # HTTP Controllers and DTOs
├── pkg/               # Shared libraries (DB Drivers, Configs)
└── database/          # Migrations and Seeds
```

## Assumptions and Considerations

For the scope of this challenge, the following assumptions and simplifications have been made:

- **Testing:** Demonstrative unit tests have been included for the main use cases (user, tweet), but 100% coverage is not provided.

- **Fan-Out Scope:** The asynchronous distribution pattern was implemented solely for the tweet.created event as a prototype. In a production system, events like "delete tweet" or "unfollow" should also be emitted to maintain cache consistency.

- **Database Agnosticism**: Although a relational database (PostgreSQL) is used for persistence, the code is decoupled via interfaces, allowing for migration to NoSQL or other engines if data volume requires it.

- **Security**: Security implementations such as password hashing (bcrypt) and authentication (JWT) have been omitted to focus on architecture and scalability patterns. Additionally, sensitive credentials (database passwords, API keys) are written in plain text in the configuration files for demonstration purposes only. In a production environment, these should be managed using secure secret management solutions (e.g., HashiCorp Vault, AWS Secrets Manager, Kubernetes Secrets).

- **Content**: The design assumes purely text-based tweets. Multimedia handling would require the integration of Object Storage (S3) and a CDN, components not represented in this diagram.

## Testing

This project includes comprehensive unit tests with mocks for all layers:

### Testing Dependencies

- **testify**: More readable assertions
- **go-sqlmock**: SQL query mocking
- **gomock**: Interface mock generation

### Running Tests

```bash
# All tests
go test ./...

# Tests with verbose mode
go test -v ./...

# Tests with coverage
go test -cover ./...

# Tests for specific layer
go test -v ./internal/infrastructure/repository/
go test -v ./internal/usecase/
go test -v ./internal/interfaces/controller/
```

### Generating Mocks

Mocks are automatically generated using `mockgen`:

```bash
# Generate all mocks
make generate-mocks

# Clean generated mocks
make clean-mocks
```

### Test Coverage

Current tests cover:

- **Repository Layer**: Tests with sqlmock for database operations
  - SelectByID (success and not found cases)
  - SelectByEmail
  - Insert (success and error cases)
  
- **Usecase Layer**: Tests with mock repository for business logic
  - CreateUser (success, duplicate email, duplicate username)
  - UpdateUser (success, user not found)
  - GetUserByID (success, error)
  
- **Controller Layer**: Tests with mock usecase for HTTP endpoints
  - GetUserByID (200, 400, 500)
  - CreateUser (201, 400, 500)
  - UpdateUser (200)
  - GetAllUsers (200)

## Available Commands

```bash
# Testing
make test              # Run all tests
make test-coverage     # Run tests with coverage
make generate-mocks    # Generate mocks
make clean-mocks       # Clean generated mocks
```