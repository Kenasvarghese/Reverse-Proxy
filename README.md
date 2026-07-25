# Reverse Proxy

A high-performance HTTP reverse proxy built in Go featuring configurable rate limiting, request logging, and optimized connection management.

## Features

- 🚀 **Reverse Proxy** - Forward requests to configurable origin servers
- 🛡️ **Rate Limiting** - Token bucket algorithm (mutex and channel implementations)
- 📊 **Request Logging** - Method, path, duration, and status tracking
- ⚙️ **Configurable** - Environment variable based configuration
- 🔧 **Production-Ready** - Connection pooling, configurable timeouts

## Quick Start
```bash
# Clone and build
git clone https://github.com/Kenasvarghese/Reverse-Proxy.git
cd Reverse-Proxy
go build -o proxy ./cmd/proxy

# Run with defaults
./proxy
```

## Configuration
```bash
# Server
export PORT=8080
export ORIGIN="http://localhost:3000"

# Rate Limiting
export MAX_ALLOWED_REQUESTS=150        # Burst capacity
export REQUEST_RATE_PER_SEC=50         # Sustained rate
export RATE_LIMITER_TYPE="TokenBucket" # or "TokenBucketWithChannel"

# Transport (optional)
export DIAL_TIMEOUT="5s"
export MAX_IDLE_CONNS=100
```

## Usage
```bash
# Start proxy
PORT=8080 ORIGIN="http://api.example.com" ./proxy

# Make requests
curl http://localhost:8080/api/users
```

## Rate Limiting

Token bucket algorithm with configurable burst and sustained rates:

- **Burst Capacity**: Up to `MAX_ALLOWED_REQUESTS` simultaneous requests
- **Sustained Rate**: `REQUEST_RATE_PER_SEC` requests per second
- **Response**: 429 Too Many Requests when limit exceeded

### Implementations

- **TokenBucket** (default) - Mutex-based, ~50ns per operation
- **TokenBucketWithChannel** - Channel-based with goroutine refill

## Project Structure
```
├── cmd/proxy/main.go           # Entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── middlewares/            # Logging, rate limiting
│   ├── proxy/                  # Reverse proxy, transport
│   └── rate_limiter/           # Rate Limiter implementations
```

## Technical Highlights

- **Middleware Pattern** - Composable request processing
- **Token Bucket Algorithm** - Two implementations (mutex/channel)
- **Connection Pooling** - Optimized HTTP transport
- **Thread-Safe** - Concurrent request handling

## Building
```bash
go build -o proxy ./cmd/proxy                    # Current platform
GOOS=linux go build -o proxy-linux ./cmd/proxy   # Linux
GOOS=darwin go build -o proxy-mac ./cmd/proxy    # macOS
```

## Author

**Kenas Varghese**  
GitHub: [@Kenasvarghese](https://github.com/Kenasvarghese)