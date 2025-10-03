# 🔥 GPU Cloud Manager - Dynamic Multi-Provider Support

A comprehensive example of the enhanced GPU Cloud Manager with dynamic model support and multiple providers.

## 🚀 New Features

### ✅ **Multiple GPU Providers**
- **Vast.ai** - Community GPU marketplace
- **RunPod** - Professional GPU cloud platform
- **Lambda Labs** - Ready to integrate
- **Paperspace** - Ready to integrate

### ✅ **Dynamic GPU Model Detection**
- **Consumer GPUs**: RTX 4090, RTX 4080, RTX 3090, RTX 3080, etc.
- **Professional GPUs**: A6000, A40, Quadro RTX 8000
- **Datacenter GPUs**: H100, A100, V100
- **AMD GPUs**: RX 7900 XTX, RX 6900 XT

### ✅ **Advanced Search & Filtering**
- Filter by GPU category (consumer, professional, datacenter)
- Performance-based filtering
- Price range filtering
- Multi-region support
- Reliability scoring

## 🎯 API Examples

### 1. Search All Available GPUs
```bash
curl "http://localhost:8080/api/v1/offers/search?available=true"
```

### 2. Find RTX 4090s Under $2/hour
```bash
curl "http://localhost:8080/api/v1/offers/search?gpu_model=RTX%204090&max_price=2.0&available=true"
```

### 3. Advanced Search with Multiple Filters
```bash
curl -X POST http://localhost:8080/api/v1/offers/search/advanced \
  -H "Content-Type: application/json" \
  -d '{
    "gpu_models": ["RTX 4090", "RTX 4080", "A100"],
    "gpu_category": "consumer",
    "min_gpu_count": 1,
    "max_price": 3.0,
    "min_ram": 16,
    "regions": ["US-East", "US-West"],
    "min_performance": 80,
    "sort_by": "price",
    "sort_order": "asc",
    "available": true
  }'
```

### 4. Get Available GPU Models and Specs
```bash
curl "http://localhost:8080/api/v1/gpu-models"
```
**Response:**
```json
{
  "success": true,
  "message": "GPU models retrieved successfully",
  "data": {
    "RTX 4090": {
      "name": "RTX 4090",
      "memory_gb": 24,
      "compute_capability": 8.9,
      "architecture": "Ada Lovelace",
      "category": "consumer",
      "performance_score": 100
    },
    "H100": {
      "name": "H100",
      "memory_gb": 80,
      "compute_capability": 9.0,
      "architecture": "Hopper",
      "category": "datacenter",
      "performance_score": 150
    }
  }
}
```

### 5. Get Provider Information
```bash
curl "http://localhost:8080/api/v1/providers"
```
**Response:**
```json
{
  "success": true,
  "data": [
    {
      "name": "vast_ai",
      "display_name": "Vast.ai",
      "website": "https://vast.ai",
      "regions": ["US-East", "US-West", "Europe", "Asia"],
      "gpu_models": ["RTX 4090", "RTX 3090", "A100", "V100"],
      "features": ["SSH Access", "Docker Support", "Jupyter Notebooks"],
      "is_configured": true
    },
    {
      "name": "runpod",
      "display_name": "RunPod",
      "website": "https://runpod.io",
      "regions": ["Global", "US", "Europe", "Asia"],
      "gpu_models": ["RTX 4090", "A100", "H100"],
      "features": ["GraphQL API", "Jupyter Support", "Community & Secure Cloud"],
      "is_configured": true
    }
  ]
}
```

### 6. Get Marketplace Statistics
```bash
curl "http://localhost:8080/api/v1/marketplace/stats"
```
**Response:**
```json
{
  "success": true,
  "data": {
    "total_instances": 1247,
    "available_instances": 892,
    "model_stats": {
      "RTX 4090": {
        "model": {
          "name": "RTX 4090",
          "memory_gb": 24,
          "performance_score": 100
        },
        "available_count": 156,
        "min_price": 0.89,
        "max_price": 2.45,
        "avg_price": 1.67,
        "providers": ["vast_ai", "runpod"]
      }
    },
    "provider_stats": {
      "vast_ai": 567,
      "runpod": 325
    },
    "average_price": 1.34,
    "price_range": {
      "min": 0.23,
      "max": 8.99,
      "avg": 1.34
    }
  }
}
```

### 7. Create Instance with Advanced Options
```bash
curl -X POST http://localhost:8080/api/v1/instances \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "runpod",
    "offer_id": "gpu_type_id_here",
    "image": "pytorch/pytorch:2.0.1-cuda11.7-cudnn8-devel",
    "label": "ML Training Workstation",
    "environment": {
      "WANDB_API_KEY": "your_wandb_key",
      "HUGGINGFACE_TOKEN": "your_hf_token"
    },
    "ports": [
      {"container_port": 8888, "protocol": "tcp"},
      {"container_port": 6006, "protocol": "tcp"}
    ],
    "resources": {
      "min_ram_gb": 32,
      "min_storage_gb": 100,
      "min_gpus": 1
    }
  }'
```

## 🎛️ Configuration

### Environment Variables
```bash
# GPU Provider API Keys
VAST_AI_API_KEY=your_vast_ai_key_here
RUNPOD_API_KEY=your_runpod_key_here
LAMBDA_API_KEY=your_lambda_key_here
PAPERSPACE_API_KEY=your_paperspace_key_here

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/gpu_cloud_manager

# Server Settings
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info

# Features
ENABLE_METRICS=true
ENABLE_CORS=true
RATE_LIMIT_ENABLED=true
RATE_LIMIT_RPM=100
```

## 🧠 Dynamic GPU Model System

The system now automatically enriches GPU instances with detailed model information:

```json
{
  "id": "vast_12345",
  "provider": "vast_ai",
  "gpu_model": "RTX 4090",
  "gpu_info": {
    "name": "RTX 4090",
    "memory_gb": 24,
    "compute_capability": 8.9,
    "architecture": "Ada Lovelace",
    "category": "consumer",
    "performance_score": 100
  },
  "performance": 100,
  "reliability": 0.95,
  "network_info": {
    "download_mbps": 1000,
    "upload_mbps": 100
  }
}
```

## 🔍 Advanced Filtering Options

### Filter by GPU Category
```bash
# Consumer GPUs (RTX series, gaming cards)
curl -X POST http://localhost:8080/api/v1/offers/search/advanced \
  -d '{"gpu_category": "consumer", "available": true}'

# Professional GPUs (Quadro, A6000)
curl -X POST http://localhost:8080/api/v1/offers/search/advanced \
  -d '{"gpu_category": "professional", "available": true}'

# Datacenter GPUs (H100, A100, V100)
curl -X POST http://localhost:8080/api/v1/offers/search/advanced \
  -d '{"gpu_category": "datacenter", "available": true}'
```

### Filter by Performance Score
```bash
# High-performance GPUs (score > 90)
curl -X POST http://localhost:8080/api/v1/offers/search/advanced \
  -d '{"min_performance": 90, "sort_by": "performance", "sort_order": "desc"}'
```

### Filter by Multiple Models
```bash
# Find any NVIDIA RTX 40-series
curl -X POST http://localhost:8080/api/v1/offers/search/advanced \
  -d '{"gpu_models": ["RTX 4090", "RTX 4080", "RTX 4070"], "available": true}'
```

### Complex Multi-Filter Search
```bash
curl -X POST http://localhost:8080/api/v1/offers/search/advanced \
  -d '{
    "gpu_models": ["RTX 4090", "A100"],
    "min_price": 1.0,
    "max_price": 5.0,
    "min_ram": 24,
    "regions": ["US-East", "Europe"],
    "min_reliability": 0.8,
    "min_performance": 85,
    "sort_by": "price",
    "sort_order": "asc"
  }'
```

## 🚀 Usage Examples

### Python Client Example
```python
import requests
import json

BASE_URL = "http://localhost:8080/api/v1"

# Search for high-performance GPUs
def find_best_gpus():
    response = requests.post(f"{BASE_URL}/offers/search/advanced", 
        json={
            "gpu_category": "datacenter",
            "min_performance": 100,
            "max_price": 10.0,
            "sort_by": "performance",
            "sort_order": "desc",
            "available": True
        }
    )
    return response.json()["data"]

# Get marketplace insights
def get_market_insights():
    response = requests.get(f"{BASE_URL}/marketplace/stats")
    stats = response.json()["data"]
    
    print(f"Total available instances: {stats['available_instances']}")
    print(f"Average price: ${stats['average_price']:.2f}/hour")
    
    # Show top 3 models by availability
    models = sorted(stats["model_stats"].items(), 
                   key=lambda x: x[1]["available_count"], 
                   reverse=True)[:3]
    
    print("\nTop 3 Available Models:")
    for model, info in models:
        print(f"  {model}: {info['available_count']} available, "
              f"avg ${info['avg_price']:.2f}/hour")

# Create optimized instance
def create_ml_instance(gpu_model="RTX 4090"):
    # First, find the best offer
    offers = requests.post(f"{BASE_URL}/offers/search/advanced",
        json={
            "gpu_model": gpu_model,
            "max_price": 3.0,
            "sort_by": "price",
            "sort_order": "asc",
            "available": True
        }
    ).json()["data"]
    
    if not offers:
        print(f"No {gpu_model} available under $3/hour")
        return None
    
    best_offer = offers[0]
    
    # Create instance
    instance_data = {
        "provider": best_offer["provider"],
        "offer_id": best_offer["provider_id"],
        "image": "pytorch/pytorch:2.0.1-cuda11.7-cudnn8-devel",
        "label": f"ML Training - {gpu_model}",
        "environment": {
            "CUDA_VISIBLE_DEVICES": "0"
        }
    }
    
    response = requests.post(f"{BASE_URL}/instances", json=instance_data)
    return response.json()["data"]

if __name__ == "__main__":
    # Get market insights
    get_market_insights()
    
    # Find and create instance
    instance = create_ml_instance("RTX 4090")
    if instance:
        print(f"\nCreated instance: {instance['id']}")
        print(f"Provider: {instance['provider']}")
        print(f"GPU: {instance['gpu_model']}")
        print(f"Price: ${instance['price_per_hour']:.2f}/hour")
```

## 🔥 What's New

### 🎯 **Smart Model Recognition**
The system now automatically recognizes and categorizes GPU models, providing:
- Performance scoring (0-150 scale)
- Memory specifications
- Architecture information
- Category classification

### 🚀 **Multi-Provider Support**
- **Unified API** across all providers
- **Provider-specific features** handled transparently
- **Automatic failover** when providers are unavailable
- **Cost optimization** across providers

### 📊 **Advanced Analytics**
- **Real-time marketplace statistics**
- **Price trend analysis**
- **Availability monitoring**
- **Performance benchmarking**

### 🎛️ **Enhanced Filtering**
- **Multi-dimensional search** (price, performance, reliability)
- **Geographic filtering** with multiple regions
- **Resource-based filtering** (RAM, storage, GPU count)
- **Category-based search** (consumer vs professional vs datacenter)

### 🔄 **Dynamic Configuration**
- **Hot-swappable providers** via environment variables
- **Runtime model updates** without restarts
- **Flexible resource requirements**
- **Custom environment variables and port mappings**

## 🎉 Migration from v1

The new system is backward compatible. Your existing API calls will continue to work, but you can now take advantage of:

1. **Enhanced responses** with GPU model information
2. **Better performance** with optimized filtering
3. **More providers** for better availability and pricing
4. **Advanced search** for precise requirements

Start using the new features by upgrading your API calls to include the new filtering options!
