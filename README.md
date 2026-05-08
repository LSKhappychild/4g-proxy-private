# 4G LTE S1AP Proxy

A transparent SCTP proxy for S1AP protocol between eNB and MME with message inspection, logging, selective message dropping, and delay capabilities.

## Features

- **S1AP Message Inspection**: Decodes and logs S1AP messages in real-time
- **NAS EPS Decoding**: Extracts and parses EMM/ESM messages from NAS PDUs
- **Message Dropping**: Selective dropping of specific signal types (Attach, Detach, TAU, etc.)
- **Message Delay**: Configurable delay for each signal type (useful for testing timeout scenarios)
- **HTTP API**: REST API for controlling drop flags, delays, and viewing statistics
- **CLI Tool**: Command-line tool for managing drop controls and delays
- **Kubernetes Ready**: Environment variable configuration for easy K8s deployment

## Architecture

```
┌─────────┐         ┌─────────────┐         ┌─────────┐
│   eNB   │◄───────►│  4G Proxy   │◄───────►│   MME   │
└─────────┘  SCTP   │             │  SCTP   └─────────┘
                    │  - Inspect  │
                    │  - Log      │
                    │  - Drop     │
                    │  - Delay    │
                    └─────────────┘
                          │
                          │ HTTP API
                          ▼
                    ┌─────────────┐
                    │   dropctl   │
                    └─────────────┘
```

## Building

```bash
# Build the proxy
go build -o 4g-proxy ./cmd/app.go

# Build the CLI tool
go build -o dropctl ./cmd/dropctl/main.go
```

## Usage

### Running the Proxy

```bash
# Run with default settings (listen on :36412, forward to localhost:36412)
./4g-proxy

# Run with custom MME address
./4g-proxy --mme 192.168.1.100 --mme-port 36412

# Run with configuration file
./4g-proxy -c config/config.yaml

# Run with delay via environment variables
DELAY_ATTACH_MS=1000 DELAY_TAU_MS=500 ./4g-proxy

# Run with verbose logging
./4g-proxy -v
```

### Command Line Options

```
Options:
  -c, --config FILE      Path to configuration file
  --listen ADDR          Listen address for eNB connections (default: 0.0.0.0)
  --port PORT            Listen port (default: 36412)
  --mme ADDR             MME address
  --mme-port PORT        MME port
  --api-port PORT        HTTP API port (default: 8080)
  -v, --verbose          Enable verbose logging
  -h, --help             Show help message
  --version              Show version
```

### Environment Variables (for Kubernetes)

| Variable | Description |
|----------|-------------|
| `PROXY_LISTEN_ADDRESS` | Listen address |
| `PROXY_LISTEN_PORT` | Listen port |
| `MME_ADDRESS` | MME address |
| `MME_PORT` | MME port |
| `API_ENABLED` | Enable HTTP API (true/false) |
| `API_PORT` | HTTP API port |
| `DELAY_ATTACH_MS` | Delay for Attach messages (ms) |
| `DELAY_DETACH_MS` | Delay for Detach messages (ms) |
| `DELAY_TAU_MS` | Delay for TAU messages (ms) |
| `DELAY_SERVICE_REQUEST_MS` | Delay for Service Request (ms) |
| `DELAY_UE_CONTEXT_RELEASE_MS` | Delay for UE Context Release (ms) |
| `DELAY_PDN_CONNECTIVITY_MS` | Delay for PDN Connectivity (ms) |
| `DELAY_HANDOVER_MS` | Delay for Handover (ms) |
| `DELAY_RESET_MS` | Delay for Reset (ms) |
| `DELAY_PAGING_MS` | Delay for Paging (ms) |
| `DELAY_DEFAULT_MS` | Default delay for other messages (ms) |

### Using the CLI Tool

```bash
# Show proxy status
./dropctl status

# Show statistics
./dropctl stats

# List current drop flags
./dropctl list

# List current delay settings
./dropctl delays

# Drop Attach messages
./dropctl drop attach

# Allow Attach messages again
./dropctl allow attach

# Add 1000ms delay to Attach messages
./dropctl delay attach 1000

# Add 500ms delay to TAU messages
./dropctl delay tau 500

# Reset all drop flags
./dropctl reset

# Reset all delays
./dropctl reset-delays
```

### Signal Types

| Signal Type | Description |
|------------|-------------|
| `attach` | Attach Request/Accept/Complete/Reject |
| `detach` | Detach Request/Accept |
| `tau` | Tracking Area Update |
| `service-request` | Service Request |
| `ue-context-release` | UE Context Release |
| `pdn-connectivity` | PDN Connectivity |
| `handover` | Handover procedures |
| `reset` | S1 Reset |
| `paging` | Paging |
| `default` | Default (for unspecified types) |

## HTTP API

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/status` | Get proxy status (includes delays) |
| GET | `/api/v1/stats` | Get statistics |
| GET | `/api/v1/drop` | Get all drop flags |
| PUT | `/api/v1/drop` | Set multiple drop flags |
| DELETE | `/api/v1/drop` | Reset all drop flags |
| PUT | `/api/v1/drop/:signalType` | Set individual drop flag |
| GET | `/api/v1/delay` | Get all delay settings |
| PUT | `/api/v1/delay` | Set multiple delays |
| DELETE | `/api/v1/delay` | Reset all delays |
| PUT | `/api/v1/delay/:signalType` | Set individual delay |

### Examples

```bash
# Get status (includes drop flags and delays)
curl http://localhost:8080/api/v1/status

# Drop Attach messages
curl -X POST http://localhost:8080/api/v1/drop/attach

# Set multiple drop flags
curl -X PUT http://localhost:8080/api/v1/drop \
  -H "Content-Type: application/json" \
  -d '{"attach": true, "detach": true}'

# Set delay for Attach messages (1000ms)
curl -X PUT http://localhost:8080/api/v1/delay/attach \
  -H "Content-Type: application/json" \
  -d '{"delayMs": 1000}'

# Set multiple delays
curl -X PUT http://localhost:8080/api/v1/delay \
  -H "Content-Type: application/json" \
  -d '{"attach": 1000, "tau": 500, "serviceRequest": 200}'

# Get current delay settings
curl http://localhost:8080/api/v1/delay

# Reset all delays
curl -X DELETE http://localhost:8080/api/v1/delay
```

## Configuration

Create a `config.yaml` file:

```yaml
proxy:
  listenAddress: "0.0.0.0"
  listenPort: 36412

mme:
  address: "192.168.1.100"
  port: 36412

api:
  enabled: true
  address: "0.0.0.0"
  port: 8080

logging:
  level: "info"
  logS1AP: true
  logNAS: true
  verbose: false

# Delay settings (in milliseconds)
delay:
  attach: 0
  detach: 0
  tau: 0
  serviceRequest: 0
  ueContextRelease: 0
  pdnConnectivity: 0
  handover: 0
  reset: 0
  paging: 0
  default: 0
```

## Kubernetes Deployment Example

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: s1ap-proxy
spec:
  replicas: 1
  selector:
    matchLabels:
      app: s1ap-proxy
  template:
    metadata:
      labels:
        app: s1ap-proxy
    spec:
      containers:
      - name: s1ap-proxy
        image: your-registry/4g-proxy:latest
        ports:
        - containerPort: 36412
          protocol: SCTP
        - containerPort: 8080
          protocol: TCP
        env:
        - name: MME_ADDRESS
          value: "mme-service.core.svc.cluster.local"
        - name: MME_PORT
          value: "36412"
        - name: DELAY_ATTACH_MS
          value: "1000"
        - name: DELAY_TAU_MS
          value: "500"
```

## S1AP Messages Supported

### Elementary Procedures

- InitialUEMessage
- UplinkNASTransport
- DownlinkNASTransport
- InitialContextSetup
- UEContextRelease
- UEContextReleaseRequest
- E-RABSetup
- S1Setup
- HandoverRequired
- HandoverRequest
- PathSwitchRequest
- Reset
- Paging

### NAS EMM Messages

- Attach Request/Accept/Complete/Reject
- Detach Request/Accept
- Tracking Area Update Request/Accept/Complete/Reject
- Extended Service Request
- Authentication Request/Response/Reject/Failure
- Security Mode Command/Complete/Reject
- Identity Request/Response

### NAS ESM Messages

- PDN Connectivity Request/Reject
- Activate Default EPS Bearer Context Request/Accept/Reject
- Activate Dedicated EPS Bearer Context Request/Accept/Reject
- Deactivate EPS Bearer Context Request/Accept

## Testing

1. Set up the proxy to listen on the S1AP port (36412)
2. Configure the eNB to connect to the proxy instead of the MME
3. Configure the proxy to forward to the actual MME
4. Observe message logs in the proxy output
5. Use the HTTP API or CLI to control message dropping and delays

## Dependencies

- [github.com/free5gc/aper](https://github.com/free5gc/aper) - ASN.1 PER encoding
- [github.com/ishidawataru/sctp](https://github.com/ishidawataru/sctp) - SCTP support for Go
- [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [gopkg.in/yaml.v2](https://gopkg.in/yaml.v2) - YAML support

## License

Private - All rights reserved
