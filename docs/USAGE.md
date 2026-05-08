# 4G LTE S1AP Proxy - Usage Guide

This guide will teach you how to use the 4G LTE S1AP Proxy to intercept, inspect, delay, and drop S1AP messages between your eNB and MME.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Network Setup](#network-setup)
3. [Running the Proxy](#running-the-proxy)
4. [Viewing Messages](#viewing-messages)
5. [Dropping Messages](#dropping-messages)
6. [Delaying Messages](#delaying-messages)
7. [Using the HTTP API](#using-the-http-api)
8. [Kubernetes Deployment](#kubernetes-deployment)
9. [Common Use Cases](#common-use-cases)
10. [Troubleshooting](#troubleshooting)

---

## Quick Start

### 1. Build the binaries

```bash
cd /home/sklee/4g-proxy-private

# Build the proxy
go build -o 4g-proxy ./cmd/app.go

# Build the CLI tool
go build -o dropctl ./cmd/dropctl/main.go
```

### 2. Run the proxy

```bash
# Basic usage: proxy listens on port 36412, forwards to MME at 192.168.1.100:36412
./4g-proxy --mme 192.168.1.100 --mme-port 36412
```

### 3. Point your eNB to the proxy

Configure your eNB to connect to the proxy's IP address instead of the MME directly.

---

## Network Setup

### Before (Direct Connection)
```
┌─────────┐                          ┌─────────┐
│   eNB   │ ──────── SCTP ─────────► │   MME   │
│         │        port 36412        │         │
└─────────┘                          └─────────┘
```

### After (With Proxy)
```
┌─────────┐         ┌─────────────┐         ┌─────────┐
│   eNB   │ ──────► │  4G Proxy   │ ──────► │   MME   │
│         │  SCTP   │ (this host) │  SCTP   │         │
└─────────┘ :36412  └─────────────┘ :36412  └─────────┘
                          │
                          │ HTTP :8080
                          ▼
                    ┌─────────────┐
                    │   dropctl   │
                    └─────────────┘
```

### IP Address Example

| Component | IP Address | Port |
|-----------|------------|------|
| eNB | 192.168.1.10 | - |
| Proxy | 192.168.1.50 | 36412 (SCTP), 8080 (HTTP) |
| MME | 192.168.1.100 | 36412 |

**eNB Configuration**: Change MME address from `192.168.1.100` to `192.168.1.50`

---

## Running the Proxy

### Basic Usage

```bash
# Minimal: specify MME address
./4g-proxy --mme 192.168.1.100 --mme-port 36412
```

### With Custom Listen Port

```bash
# Listen on a different port (e.g., 36413)
./4g-proxy --port 36413 --mme 192.168.1.100 --mme-port 36412
```

### With Verbose Logging

```bash
# See detailed message logs
./4g-proxy --mme 192.168.1.100 --mme-port 36412 -v
```

### Using a Config File

```bash
# Use configuration file
./4g-proxy -c config/config.yaml
```

### With Pre-configured Delays

```bash
# Add 1 second delay to Attach, 500ms to TAU
DELAY_ATTACH_MS=1000 DELAY_TAU_MS=500 ./4g-proxy --mme 192.168.1.100 --mme-port 36412
```

---

## Viewing Messages

When the proxy is running, you'll see S1AP messages logged in real-time:

```
2024/01/15 10:30:45 S1AP [UL] InitialUEMessage (EMM: AttachRequest) UE: eNB-UE-ID=1 [Signal: Attach]
2024/01/15 10:30:45 S1AP [DL] DownlinkNASTransport (EMM: AuthenticationRequest) UE: MME-UE-ID=1 eNB-UE-ID=1
2024/01/15 10:30:46 S1AP [UL] UplinkNASTransport (EMM: AuthenticationResponse) UE: MME-UE-ID=1 eNB-UE-ID=1
2024/01/15 10:30:46 S1AP [DL] DownlinkNASTransport (EMM: SecurityModeCommand) UE: MME-UE-ID=1 eNB-UE-ID=1
2024/01/15 10:30:46 S1AP [UL] UplinkNASTransport (EMM: SecurityModeComplete) UE: MME-UE-ID=1 eNB-UE-ID=1
2024/01/15 10:30:47 S1AP [DL] InitialContextSetupRequest (EMM: AttachAccept) UE: MME-UE-ID=1 eNB-UE-ID=1 [Signal: Attach]
```

### Understanding the Log Format

```
[UL] = Uplink (eNB → MME)
[DL] = Downlink (MME → eNB)

Signal types:
- Attach      = Attach procedure
- Detach      = Detach procedure
- TAU         = Tracking Area Update
- ServiceRequest = Service Request
- UEContextRelease = UE Context Release
```

---

## Dropping Messages

### Using the CLI Tool

```bash
# Drop all Attach messages (UE won't be able to attach)
./dropctl drop attach

# Check current drop settings
./dropctl list

# Allow Attach messages again
./dropctl allow attach

# Drop multiple types
./dropctl drop attach
./dropctl drop tau

# Reset all drop flags
./dropctl reset
```

### What Happens When You Drop

When a message is dropped, the proxy logs it:

```
2024/01/15 10:35:12 DROPPED: [UL] InitialUEMessage (EMM: AttachRequest) UE: eNB-UE-ID=2 [Signal: Attach]
```

The message is not forwarded to the destination, simulating network loss or filtering.

---

## Delaying Messages

### Using the CLI Tool

```bash
# Add 1 second (1000ms) delay to Attach messages
./dropctl delay attach 1000

# Add 500ms delay to TAU messages
./dropctl delay tau 500

# Check current delays
./dropctl delays

# Remove delay (set to 0)
./dropctl delay attach 0

# Reset all delays
./dropctl reset-delays
```

### Using Environment Variables

```bash
# Set delays before starting the proxy
export DELAY_ATTACH_MS=1000
export DELAY_TAU_MS=500
./4g-proxy --mme 192.168.1.100 --mme-port 36412
```

### What Happens When You Delay

When a message is delayed, the proxy logs it:

```
2024/01/15 10:40:22 DELAYING 1s: [UL] InitialUEMessage (EMM: AttachRequest) UE: eNB-UE-ID=3 [Signal: Attach]
```

The message is held for the specified duration before being forwarded.

---

## Using the HTTP API

The proxy exposes a REST API on port 8080 (default).

### Check Status

```bash
curl http://localhost:8080/api/v1/status
```

Response:
```json
{
  "status": "running",
  "stats": {
    "totalUplinkPackets": 42,
    "totalDownlinkPackets": 38,
    "droppedPackets": 2,
    "delayedPackets": 5
  },
  "dropFlags": {
    "attach": false,
    "detach": false,
    "tau": false
  },
  "delayConfig": {
    "attach": 1000,
    "detach": 0,
    "tau": 500
  }
}
```

### Get Statistics

```bash
curl http://localhost:8080/api/v1/stats
```

### Drop Messages via API

```bash
# Enable dropping Attach
curl -X POST http://localhost:8080/api/v1/drop/attach

# Disable dropping Attach
curl -X DELETE http://localhost:8080/api/v1/drop/attach

# Set multiple drop flags
curl -X PUT http://localhost:8080/api/v1/drop \
  -H "Content-Type: application/json" \
  -d '{"attach": true, "tau": true}'
```

### Set Delays via API

```bash
# Set 1000ms delay for Attach
curl -X PUT http://localhost:8080/api/v1/delay/attach \
  -H "Content-Type: application/json" \
  -d '{"delayMs": 1000}'

# Set multiple delays
curl -X PUT http://localhost:8080/api/v1/delay \
  -H "Content-Type: application/json" \
  -d '{"attach": 1000, "tau": 500, "serviceRequest": 200}'

# Get current delays
curl http://localhost:8080/api/v1/delay

# Reset all delays
curl -X DELETE http://localhost:8080/api/v1/delay
```

---

## Kubernetes Deployment

### Deployment YAML

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: s1ap-proxy
  namespace: lte-core
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
        - name: s1ap
          containerPort: 36412
          protocol: SCTP
        - name: api
          containerPort: 8080
          protocol: TCP
        env:
        - name: MME_ADDRESS
          value: "mme-service.lte-core.svc.cluster.local"
        - name: MME_PORT
          value: "36412"
        # Optional: Set initial delays
        - name: DELAY_ATTACH_MS
          value: "0"
        - name: DELAY_TAU_MS
          value: "0"
---
apiVersion: v1
kind: Service
metadata:
  name: s1ap-proxy
  namespace: lte-core
spec:
  selector:
    app: s1ap-proxy
  ports:
  - name: s1ap
    port: 36412
    protocol: SCTP
  - name: api
    port: 8080
    protocol: TCP
```

### Changing Delays at Runtime

Once deployed, use the API to change delays without restarting:

```bash
# From within the cluster or via port-forward
kubectl port-forward svc/s1ap-proxy 8080:8080 -n lte-core

# Then set delays
curl -X PUT http://localhost:8080/api/v1/delay \
  -H "Content-Type: application/json" \
  -d '{"attach": 2000}'
```

---

## Common Use Cases

### 1. Testing Attach Timeout

Simulate slow network causing attach timeout:

```bash
# Set 5 second delay on Attach Accept
./dropctl delay attach 5000
```

### 2. Blocking New Attachments

Prevent new UEs from attaching while allowing existing UEs:

```bash
./dropctl drop attach
```

### 3. Testing TAU Failure

Simulate TAU failures:

```bash
./dropctl drop tau
```

### 4. Simulating Network Congestion

Add delays to all message types:

```bash
./dropctl delay attach 1000
./dropctl delay tau 1000
./dropctl delay serviceRequest 500
./dropctl delay pdnConnectivity 1000
```

### 5. Testing Handover Scenarios

Add delay to handover messages:

```bash
# Delay all handover messages
./dropctl delay handover 2000

# Or delay specific handover messages to simulate core-side lag:

# Delay HandoverRequired (Source eNB -> MME, initiates S1 handover)
./dropctl delay handover-required 2000

# Delay HandoverNotify (Target eNB -> MME, UE has arrived at target cell)
./dropctl delay handover-notify 1500
```

Using environment variables for K8s:
```bash
DELAY_HANDOVER_REQUIRED_MS=2000 DELAY_HANDOVER_NOTIFY_MS=1500 ./4g-proxy --mme 192.168.1.100 --mme-port 36412
```

### 6. Observing Message Flow

Just run the proxy without any drops/delays to observe messages:

```bash
./4g-proxy --mme 192.168.1.100 --mme-port 36412 -v
```

---

## Troubleshooting

### Proxy Won't Start

**Error**: `failed to listen on SCTP`

- Make sure no other process is using port 36412
- Check if SCTP kernel module is loaded: `lsmod | grep sctp`
- Load SCTP if needed: `sudo modprobe sctp`

### eNB Can't Connect

- Verify the proxy is listening: `ss -lnp | grep 36412`
- Check firewall rules allow SCTP traffic
- Verify eNB is configured with proxy's IP address

### No Messages Appearing

- Verify eNB is actually connecting (check for "New session" log)
- Make sure the MME address is correct
- Check if MME is reachable from proxy: `nc -z <mme-ip> 36412`

### dropctl Can't Connect

**Error**: `Error: connection refused`

- Make sure proxy is running
- Check API port (default 8080)
- Verify: `curl http://localhost:8080/health`

### Messages Not Being Dropped/Delayed

- Verify the signal type name is correct
- Check current settings: `./dropctl list` and `./dropctl delays`
- Remember: signal types are case-sensitive in API

---

## Signal Type Reference

| Signal Type | CLI Name | Affects |
|-------------|----------|---------|
| Attach | `attach` | AttachRequest, AttachAccept, AttachComplete, AttachReject |
| Detach | `detach` | DetachRequest, DetachAccept |
| TAU | `tau` | TAURequest, TAUAccept, TAUComplete, TAUReject |
| Service Request | `service-request` | ServiceRequest, ExtendedServiceRequest |
| UE Context Release | `ue-context-release` | UEContextReleaseRequest, UEContextReleaseCommand |
| PDN Connectivity | `pdn-connectivity` | PDNConnectivityRequest, PDNConnectivityReject |
| Handover | `handover` | All handover messages (general) |
| Handover Required | `handover-required` | HandoverRequired (Source eNB -> MME) |
| Handover Notify | `handover-notify` | HandoverNotify (Target eNB -> MME) |
| Reset | `reset` | S1 Reset |
| Paging | `paging` | Paging |
| Default | `default` | All other unspecified messages |

**Note**: Handover-specific delays (`handover-required`, `handover-notify`) take precedence over the general `handover` delay when set.

---

## Getting Help

```bash
# Proxy help
./4g-proxy --help

# CLI tool help
./dropctl --help
```

For issues, check the logs and verify your network setup matches the expected topology.
