# Quick Reference Card

## Starting the Proxy

```bash
# Basic
./4g-proxy --mme <MME_IP> --mme-port 36412

# With delays
DELAY_ATTACH_MS=1000 ./4g-proxy --mme <MME_IP> --mme-port 36412

# With config file
./4g-proxy -c config/config.yaml
```

## CLI Commands (dropctl)

| Command | Description |
|---------|-------------|
| `./dropctl status` | Show proxy status |
| `./dropctl stats` | Show statistics |
| `./dropctl list` | List drop flags |
| `./dropctl delays` | List delay settings |
| `./dropctl drop <type>` | Drop message type |
| `./dropctl allow <type>` | Allow message type |
| `./dropctl delay <type> <ms>` | Set delay in milliseconds |
| `./dropctl reset` | Reset all drop flags |
| `./dropctl reset-delays` | Reset all delays |

## Signal Types

| Type | CLI Name |
|------|----------|
| Attach | `attach` |
| Detach | `detach` |
| TAU | `tau` |
| Service Request | `service-request` |
| UE Context Release | `ue-context-release` |
| PDN Connectivity | `pdn-connectivity` |
| Handover | `handover` |
| Reset | `reset` |
| Paging | `paging` |
| Default | `default` |

## HTTP API Quick Reference

```bash
# Status
curl http://localhost:8080/api/v1/status

# Drop
curl -X POST http://localhost:8080/api/v1/drop/attach
curl -X DELETE http://localhost:8080/api/v1/drop/attach

# Delay
curl -X PUT http://localhost:8080/api/v1/delay/attach \
  -H "Content-Type: application/json" \
  -d '{"delayMs": 1000}'

# Reset
curl -X DELETE http://localhost:8080/api/v1/drop
curl -X DELETE http://localhost:8080/api/v1/delay
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `MME_ADDRESS` | MME IP address |
| `MME_PORT` | MME port |
| `DELAY_ATTACH_MS` | Attach delay (ms) |
| `DELAY_DETACH_MS` | Detach delay (ms) |
| `DELAY_TAU_MS` | TAU delay (ms) |
| `DELAY_SERVICE_REQUEST_MS` | Service Request delay (ms) |
| `DELAY_UE_CONTEXT_RELEASE_MS` | UE Context Release delay (ms) |
| `DELAY_PDN_CONNECTIVITY_MS` | PDN Connectivity delay (ms) |
| `DELAY_HANDOVER_MS` | Handover delay (ms) |
| `DELAY_DEFAULT_MS` | Default delay (ms) |

## Common Examples

```bash
# Block all new attachments
./dropctl drop attach

# Simulate slow attach (3 seconds)
./dropctl delay attach 3000

# Simulate slow network
./dropctl delay attach 1000
./dropctl delay tau 1000
./dropctl delay serviceRequest 500

# Reset everything
./dropctl reset
./dropctl reset-delays
```
