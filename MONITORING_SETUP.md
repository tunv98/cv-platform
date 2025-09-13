# CV Platform Monitoring Setup

Hướng dẫn setup Prometheus + Grafana để monitor CV Platform API.

## 🏗️ Architecture

```
CV Platform API (:8080)
    ↓ /metrics endpoint
Prometheus (:9090)
    ↓ scrape metrics
Grafana (:3000)
    ↓ visualize data
```

## 📋 Prerequisites

- Docker & Docker Compose
- CV Platform API đang chạy trên port 8080
- Port 3000, 9090, 9100 available

## 🚀 Quick Start

### 1. Start Monitoring Stack

```bash
# Start Prometheus + Grafana + Node Exporter
docker-compose -f docker-compose.monitoring.yml up -d

# Check status
docker-compose -f docker-compose.monitoring.yml ps
```

### 2. Start CV Platform API

```bash
# Build and run your API
go run cmd/api/main.go
```

### 3. Access Services

| Service         | URL                           | Credentials      |
| --------------- | ----------------------------- | ---------------- |
| **Grafana**     | http://localhost:3000         | admin / admin123 |
| **Prometheus**  | http://localhost:9090         | -                |
| **API Metrics** | http://localhost:8080/metrics | -                |

## 📊 Grafana Dashboard

Dashboard "**CV Platform API Metrics**" sẽ tự động được load với các panels:

### 🔴 RED Metrics (Request, Error, Duration)

- **HTTP Request Rate**: Rate của HTTP requests theo method/path/status
- **Total HTTP Requests**: Tổng số requests
- **HTTP Status Codes**: Phân bố status codes (pie chart)
- **HTTP Request Latency**: P95, P50 latency của các endpoints

### 📈 Business Metrics

- **CV Upload Operations Rate**: Rate của start/complete upload operations
- **Profile Requests Rate**: Rate của profile requests
- **CV Upload Duration**: Latency của CV upload operations

## 🎯 Key Metrics Explained

### HTTP Metrics

```promql
# Request rate per second
rate(http_requests_total[5m])

# 95th percentile latency
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))

# Error rate
rate(http_errors_total[5m]) / rate(http_requests_total[5m])
```

### Business Metrics

```promql
# CV upload success rate
rate(cv_uploads_total{status="success"}[5m]) / rate(cv_uploads_total[5m])

# Profile request rate
rate(profile_requests_total[5m])
```

## 🔧 Configuration

### Prometheus Target

Trong `monitoring/prometheus.yml`, target được config cho API:

```yaml
- job_name: "cv-platform-api"
  static_configs:
    - targets: ["host.docker.internal:8080"] # Adjust port if needed
```

**⚠️ Important**: Nếu API chạy trên port khác, update target trong `prometheus.yml`.

### Custom Dashboards

Để tạo dashboard mới:

1. Access Grafana → Create → Dashboard
2. Add queries với Prometheus datasource
3. Export JSON và save vào `monitoring/grafana/dashboards/`

## 🧪 Testing Metrics

### Generate Test Traffic

```bash
# Test CV upload endpoints
curl -X POST http://localhost:8080/api/v1/cvs/upload \
  -H "Content-Type: application/json" \
  -d '{"file_name":"test.pdf","mime_type":"application/pdf"}'

# Test profile endpoint
curl http://localhost:8080/api/v1/profiles/123456789

# View raw metrics
curl http://localhost:8080/metrics
```

### Expected Metrics Output

```
# HTTP metrics
http_requests_total{method="POST",path="/api/v1/cvs/upload",status="200"} 5
http_request_duration_seconds_bucket{method="POST",path="/api/v1/cvs/upload",status="200",le="0.025"} 3

# Business metrics
cv_uploads_total{operation="start",status="success"} 5
profile_requests_total{operation="get",status="success"} 10
```

## 🔍 Troubleshooting

### API Metrics Not Showing

1. **Check API is running**:

   ```bash
   curl http://localhost:8080/metrics
   ```

2. **Check Prometheus targets**:

   - Go to http://localhost:9090/targets
   - Ensure `cv-platform-api` target is UP

3. **Check Docker network**:
   ```bash
   docker network ls
   docker-compose -f docker-compose.monitoring.yml logs prometheus
   ```

### Grafana Dashboard Issues

1. **Dashboard not loading**:

   - Check datasource: Grafana → Configuration → Data Sources
   - Verify Prometheus URL: `http://prometheus:9090`

2. **No data in panels**:
   - Check time range (default: last 1 hour)
   - Verify metrics exist in Prometheus

### Performance Tuning

```yaml
# prometheus.yml - adjust scrape intervals
global:
  scrape_interval: 15s # Default scrape interval

scrape_configs:
  - job_name: "cv-platform-api"
    scrape_interval: 5s # More frequent for API
```

## 📚 Advanced Usage

### Alerting Rules

Create `monitoring/rules.yml`:

```yaml
groups:
  - name: cv-platform-alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_errors_total[5m]) / rate(http_requests_total[5m]) > 0.1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
```

### Custom Metrics

Add new business metrics trong `internal/metrics/business.go`:

```go
var NewMetric = promauto.NewCounterVec(
    prometheus.CounterOpts{
        Name: "custom_metric_total",
        Help: "Description of custom metric",
    },
    []string{"label1", "label2"},
)
```

## 🛑 Cleanup

```bash
# Stop monitoring stack
docker-compose -f docker-compose.monitoring.yml down

# Remove volumes (⚠️ deletes all data)
docker-compose -f docker-compose.monitoring.yml down -v
```

---

## 📞 Support

Nếu có issues:

1. Check logs: `docker-compose -f docker-compose.monitoring.yml logs`
2. Verify ports không bị conflict
3. Ensure API `/metrics` endpoint accessible
