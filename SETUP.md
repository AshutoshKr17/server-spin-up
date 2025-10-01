# GPU Cloud Manager Setup Guide

This guide will help you set up and run the GPU Cloud Manager locally or in production.

## Prerequisites

- Go 1.21 or later
- PostgreSQL 12 or later
- Redis (optional, for caching)
- Docker and Docker Compose (optional)

## Quick Start with Docker Compose

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd gpu-cloud-manager
   ```

2. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your API keys and configuration
   ```

3. **Start the services:**
   ```bash
   docker-compose up -d
   ```

4. **Check if everything is running:**
   ```bash
   curl http://localhost:8080/health
   ```

## Manual Setup

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Database Setup

Create a PostgreSQL database:
```sql
CREATE DATABASE gpu_cloud_manager;
CREATE USER gpu_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE gpu_cloud_manager TO gpu_user;
```

### 3. Environment Configuration

Create a `.env` file in the project root:
```env
# Server Configuration
PORT=8080
ENVIRONMENT=development

# Database Configuration
DATABASE_URL=postgres://gpu_user:your_password@localhost:5432/gpu_cloud_manager?sslmode=disable

# GPU Provider API Keys
VAST_AI_API_KEY=your_vast_ai_api_key_here

# Feature Flags
ENABLE_METRICS=true
ENABLE_CORS=true

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_RPM=100

# Logging
LOG_LEVEL=info
```

### 4. Run Database Migrations

The application will automatically run migrations on startup, but you can also run them manually:
```bash
go run cmd/api/main.go
```

### 5. Start the Application

```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`.

## Getting Vast.ai API Key

1. Sign up at [Vast.ai](https://vast.ai/)
2. Go to Account Settings
3. Generate a new API key
4. Copy the key to your `.env` file

## Testing the Setup

### 1. Health Check
```bash
curl http://localhost:8080/health
```

### 2. Test Vast.ai Integration
```bash
# Search for available GPU offers
curl "http://localhost:8080/api/v1/offers/search?provider=vast_ai&available=true"

# Get supported providers
curl "http://localhost:8080/api/v1/providers"
```

### 3. Create and Manage Instances
```bash
# Create an instance (replace offer_id with actual offer ID from search)
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "vast_ai",
    "offer_id": "12345",
    "image": "pytorch/pytorch:latest",
    "label": "Test Instance"
  }' \
  "http://localhost:8080/api/v1/instances"

# List your instances
curl "http://localhost:8080/api/v1/instances"
```

## Development Setup

### 1. Install Development Tools

```bash
# Install air for hot reloading
go install github.com/cosmtrek/air@latest

# Install golangci-lint for linting
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 2. Development Commands

```bash
# Run with hot reloading
air

# Run linting
golangci-lint run

# Run tests
go test ./...

# Format code
go fmt ./...
```

### 3. Development Environment Variables

For development, you might want to use these additional settings:
```env
# Development specific settings
ENVIRONMENT=development
LOG_LEVEL=debug
GIN_MODE=debug
ENABLE_METRICS=false
```

## Production Deployment

### 1. Using Docker

Build and run the production image:
```bash
# Build the image
docker build -t gpu-cloud-manager .

# Run the container
docker run -d \
  --name gpu-cloud-manager \
  -p 8080:8080 \
  -e DATABASE_URL="your_production_database_url" \
  -e VAST_AI_API_KEY="your_api_key" \
  -e ENVIRONMENT=production \
  gpu-cloud-manager
```

### 2. Environment Variables for Production

```env
# Production settings
ENVIRONMENT=production
PORT=8080
LOG_LEVEL=warn
GIN_MODE=release

# Database (use connection pooling)
DATABASE_URL=postgres://user:pass@localhost:5432/gpu_cloud_manager?sslmode=require

# Security
ENABLE_CORS=false  # Configure properly for your domain

# Rate limiting (adjust based on your needs)
RATE_LIMIT_ENABLED=true
RATE_LIMIT_RPM=1000
```

### 3. Reverse Proxy Setup (Nginx)

Create `/etc/nginx/sites-available/gpu-cloud-manager`:
```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 4. SSL/TLS Setup

Use Let's Encrypt for free SSL certificates:
```bash
# Install certbot
sudo apt-get install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d your-domain.com

# Auto-renewal is set up automatically
```

## Monitoring and Logging

### 1. Application Logs

Logs are written to stdout in JSON format for production:
```bash
# View logs in Docker
docker logs -f gpu-cloud-manager

# View logs in development
tail -f /var/log/gpu-cloud-manager.log
```

### 2. Health Monitoring

Set up health checks:
```bash
# Simple health check
curl -f http://localhost:8080/health || exit 1

# More detailed monitoring with external tools like Prometheus
# Metrics endpoint (if enabled)
curl http://localhost:8080/metrics
```

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   ```bash
   # Check if PostgreSQL is running
   sudo systemctl status postgresql
   
   # Test connection
   psql "postgres://user:pass@localhost:5432/gpu_cloud_manager"
   ```

2. **Vast.ai API Errors**
   ```bash
   # Test API key validity
   curl -H "Authorization: Bearer YOUR_API_KEY" \
        "https://console.vast.ai/api/v0/instances"
   ```

3. **Port Already in Use**
   ```bash
   # Find process using port 8080
   lsof -i :8080
   
   # Kill the process or use a different port
   export PORT=8081
   ```

### Performance Tuning

1. **Database Optimization**
   ```sql
   -- Check slow queries
   SELECT query, mean_time, calls 
   FROM pg_stat_statements 
   ORDER BY mean_time DESC LIMIT 10;
   
   -- Optimize connection pool
   ALTER SYSTEM SET max_connections = 200;
   ALTER SYSTEM SET shared_buffers = '256MB';
   ```

2. **Application Optimization**
   ```env
   # Increase connection pool size
   DATABASE_MAX_OPEN_CONNS=100
   DATABASE_MAX_IDLE_CONNS=10
   
   # Enable caching
   REDIS_URL=redis://localhost:6379
   CACHE_TTL=3600
   ```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines and how to contribute to the project.

## Security

- Never commit API keys or sensitive data to version control
- Use environment variables for all configuration
- Enable CORS only for trusted domains in production
- Implement proper authentication and authorization
- Keep dependencies updated
- Use HTTPS in production

## Support

- Check the [API Documentation](API_DOCUMENTATION.md)
- Review [Common Issues](https://github.com/your-repo/issues)
- Create a new issue if you need help
