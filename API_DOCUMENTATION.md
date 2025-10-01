# GPU Cloud Manager API Documentation

A REST API for managing GPU instances across multiple cloud providers like Vast.ai.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Currently, authentication is handled via API keys. Include your API key in the Authorization header:
```
Authorization: Bearer your_api_key_here
```

## Endpoints

### Health Check
```http
GET /health
```
Returns the health status of the API.

**Response:**
```json
{
  "status": "healthy",
  "service": "gpu-cloud-manager",
  "version": "1.0.0"
}
```

---

### Get Supported Providers
```http
GET /api/v1/providers
```
Returns a list of supported GPU cloud providers.

**Response:**
```json
{
  "success": true,
  "message": "Providers retrieved successfully",
  "data": ["vast_ai"]
}
```

---

### Search GPU Offers
```http
GET /api/v1/offers/search
```
Search for available GPU offers across different providers.

**Query Parameters:**
- `provider` (string, optional): Filter by provider (e.g., "vast_ai")
- `gpu_model` (string, optional): Filter by GPU model (e.g., "RTX 4090")
- `min_gpu_count` (int, optional): Minimum number of GPUs
- `max_price` (float, optional): Maximum price per hour
- `region` (string, optional): Filter by region/datacenter
- `available` (bool, optional): Show only available instances

**Example Request:**
```http
GET /api/v1/offers/search?provider=vast_ai&gpu_model=RTX 4090&max_price=2.0&available=true
```

**Response:**
```json
{
  "success": true,
  "message": "Offers retrieved successfully",
  "data": [
    {
      "id": "vast_12345",
      "provider": "vast_ai",
      "provider_id": "12345",
      "name": "Vast.ai Machine 67890",
      "status": "offline",
      "gpu_model": "RTX 4090",
      "gpu_count": 1,
      "cpu_count": 8,
      "ram_gb": 32,
      "storage_gb": 100,
      "price_per_hour": 1.50,
      "region": "US-East",
      "provider_data": {
        "machine_id": 67890,
        "compute_cap": 89,
        "reliability": 0.95,
        "score": 8.5
      }
    }
  ]
}
```

---

### Get User Instances
```http
GET /api/v1/instances
```
Retrieve all GPU instances owned by the authenticated user.

**Response:**
```json
{
  "success": true,
  "message": "Instances retrieved successfully",
  "data": [
    {
      "id": "vast_12345",
      "provider": "vast_ai",
      "provider_id": "12345",
      "name": "My Training Instance",
      "status": "running",
      "gpu_model": "RTX 4090",
      "gpu_count": 1,
      "cpu_count": 8,
      "ram_gb": 32,
      "storage_gb": 100,
      "price_per_hour": 1.50,
      "region": "US-East",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z",
      "provider_data": {
        "ssh_host": "ssh5.vast.ai",
        "ssh_port": 12345,
        "public_ipaddr": "1.2.3.4"
      }
    }
  ]
}
```

---

### Get Instance Details
```http
GET /api/v1/instances/{id}
```
Retrieve details of a specific GPU instance.

**Path Parameters:**
- `id` (string): Instance ID (e.g., "vast_12345")

**Response:**
```json
{
  "success": true,
  "message": "Instance retrieved successfully",
  "data": {
    "id": "vast_12345",
    "provider": "vast_ai",
    "provider_id": "12345",
    "name": "My Training Instance",
    "status": "running",
    "gpu_model": "RTX 4090",
    "gpu_count": 1,
    "cpu_count": 8,
    "ram_gb": 32,
    "storage_gb": 100,
    "price_per_hour": 1.50,
    "region": "US-East",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "provider_data": {
      "ssh_host": "ssh5.vast.ai",
      "ssh_port": 12345,
      "public_ipaddr": "1.2.3.4",
      "image": "pytorch/pytorch:latest"
    }
  }
}
```

---

### Create Instance
```http
POST /api/v1/instances
```
Create a new GPU instance from an available offer.

**Request Body:**
```json
{
  "provider": "vast_ai",
  "offer_id": "12345",
  "image": "pytorch/pytorch:latest",
  "label": "My Training Instance",
  "onstart_script": "#!/bin/bash\necho 'Instance started'\npip install -r requirements.txt",
  "ssh_key": "ssh-rsa AAAAB3NzaC1yc2E..."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Instance created successfully",
  "data": {
    "id": "vast_12345",
    "provider": "vast_ai",
    "provider_id": "12345",
    "name": "My Training Instance",
    "status": "loading",
    "price_per_hour": 1.50,
    "created_at": "2024-01-15T10:30:00Z",
    "provider_data": {
      "machine_id": 67890,
      "image": "pytorch/pytorch:latest",
      "onstart_script": "#!/bin/bash\necho 'Instance started'\npip install -r requirements.txt"
    }
  }
}
```

---

### Start Instance
```http
POST /api/v1/instances/{id}/start
```
Start a stopped GPU instance.

**Path Parameters:**
- `id` (string): Instance ID

**Response:**
```json
{
  "success": true,
  "message": "Instance started successfully"
}
```

---

### Stop Instance
```http
POST /api/v1/instances/{id}/stop
```
Stop a running GPU instance.

**Path Parameters:**
- `id` (string): Instance ID

**Response:**
```json
{
  "success": true,
  "message": "Instance stopped successfully"
}
```

---

### Destroy Instance
```http
DELETE /api/v1/instances/{id}
```
Permanently terminate a GPU instance.

**Path Parameters:**
- `id` (string): Instance ID

**Response:**
```json
{
  "success": true,
  "message": "Instance destroyed successfully"
}
```

---

## Error Responses

All error responses follow this format:
```json
{
  "success": false,
  "error": "Error message describing what went wrong"
}
```

### Common HTTP Status Codes
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Rate Limiting

The API implements rate limiting to prevent abuse:
- Default: 100 requests per minute per IP address
- Rate limit headers are included in responses:
  - `X-RateLimit-Limit`: Request limit per window
  - `X-RateLimit-Remaining`: Remaining requests in current window
  - `X-RateLimit-Reset`: Window reset time (Unix timestamp)

## Provider-Specific Notes

### Vast.ai
- Instance IDs are prefixed with "vast_"
- Supports various Docker images from Docker Hub
- SSH access is provided via `ssh_host` and `ssh_port` in provider_data
- OnStart scripts can be used for automatic setup

## Examples

### Python Client Example
```python
import requests

BASE_URL = "http://localhost:8080/api/v1"
API_KEY = "your_api_key_here"

headers = {
    "Authorization": f"Bearer {API_KEY}",
    "Content-Type": "application/json"
}

# Search for GPU offers
response = requests.get(
    f"{BASE_URL}/offers/search",
    params={
        "provider": "vast_ai",
        "gpu_model": "RTX 4090",
        "max_price": 2.0,
        "available": True
    },
    headers=headers
)
offers = response.json()["data"]

# Create an instance from the first offer
if offers:
    offer_id = offers[0]["provider_id"]
    create_response = requests.post(
        f"{BASE_URL}/instances",
        json={
            "provider": "vast_ai",
            "offer_id": offer_id,
            "image": "pytorch/pytorch:latest",
            "label": "My ML Training Instance"
        },
        headers=headers
    )
    instance = create_response.json()["data"]
    print(f"Created instance: {instance['id']}")
```

### cURL Examples
```bash
# Search offers
curl -H "Authorization: Bearer your_api_key" \
     "http://localhost:8080/api/v1/offers/search?provider=vast_ai&available=true"

# Create instance
curl -X POST \
     -H "Authorization: Bearer your_api_key" \
     -H "Content-Type: application/json" \
     -d '{
       "provider": "vast_ai",
       "offer_id": "12345",
       "image": "pytorch/pytorch:latest",
       "label": "Training Instance"
     }' \
     "http://localhost:8080/api/v1/instances"

# Get instances
curl -H "Authorization: Bearer your_api_key" \
     "http://localhost:8080/api/v1/instances"
```
