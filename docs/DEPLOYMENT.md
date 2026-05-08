# Docker & Kubernetes Deployment Guide

This guide explains how to build the Docker image and deploy the S1AP proxy on Kubernetes.

## Table of Contents

1. [Building the Docker Image](#building-the-docker-image)
2. [Running with Docker](#running-with-docker)
3. [Pushing to Registry](#pushing-to-registry)
4. [Kubernetes Deployment](#kubernetes-deployment)
5. [Configuration](#configuration)
6. [Accessing the Proxy](#accessing-the-proxy)
7. [Updating Configuration](#updating-configuration)
8. [Monitoring](#monitoring)
9. [Troubleshooting](#troubleshooting)

---

## Building the Docker Image

### Prerequisites

- Docker installed
- Go 1.21+ (for local builds)

### Build the Image

```bash
cd /home/sklee/4g-proxy-private

# Build the image
docker build -t s1ap-proxy:latest .

# Verify the image
docker images | grep s1ap-proxy
```

### Build with a Specific Tag

```bash
# With version tag
docker build -t s1ap-proxy:v1.0.0 .

# With registry prefix
docker build -t myregistry.com/s1ap-proxy:v1.0.0 .
```

---

## Running with Docker

### Basic Run

```bash
docker run -d \
  --name s1ap-proxy \
  -p 36412:36412/sctp \
  -p 8080:8080 \
  -e MME_ADDRESS=192.168.1.100 \
  -e MME_PORT=36412 \
  s1ap-proxy:latest
```

### With Delay Configuration

```bash
docker run -d \
  --name s1ap-proxy \
  -p 36412:36412/sctp \
  -p 8080:8080 \
  -e MME_ADDRESS=192.168.1.100 \
  -e MME_PORT=36412 \
  -e DELAY_ATTACH_MS=1000 \
  -e DELAY_TAU_MS=500 \
  s1ap-proxy:latest
```

### View Logs

```bash
docker logs -f s1ap-proxy
```

### Stop and Remove

```bash
docker stop s1ap-proxy
docker rm s1ap-proxy
```

---

## Pushing to Registry

### Docker Hub

```bash
# Tag the image
docker tag s1ap-proxy:latest yourusername/s1ap-proxy:latest

# Login
docker login

# Push
docker push yourusername/s1ap-proxy:latest
```

### Private Registry

```bash
# Tag for private registry
docker tag s1ap-proxy:latest myregistry.com/lte/s1ap-proxy:v1.0.0

# Login to private registry
docker login myregistry.com

# Push
docker push myregistry.com/lte/s1ap-proxy:v1.0.0
```

### AWS ECR

```bash
# Get login token
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 123456789.dkr.ecr.us-east-1.amazonaws.com

# Tag
docker tag s1ap-proxy:latest 123456789.dkr.ecr.us-east-1.amazonaws.com/s1ap-proxy:latest

# Push
docker push 123456789.dkr.ecr.us-east-1.amazonaws.com/s1ap-proxy:latest
```

---

## Kubernetes Deployment

### Prerequisites

- Kubernetes cluster with SCTP support
- kubectl configured
- Image pushed to accessible registry

### Check SCTP Support

```bash
# SCTP must be enabled in your cluster
# Check if SCTP is available
kubectl get nodes -o jsonpath='{.items[*].status.conditions[?(@.type=="Ready")].status}'
```

> **Note**: SCTP support varies by Kubernetes distribution. Most cloud providers require additional configuration. On bare-metal, ensure the SCTP kernel module is loaded on all nodes.

### Deploy Using kubectl

```bash
cd /home/sklee/4g-proxy-private/deploy/kubernetes

# Create namespace and deploy all resources
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

### Deploy Using Kustomize

```bash
cd /home/sklee/4g-proxy-private/deploy/kubernetes

# Preview what will be deployed
kubectl kustomize .

# Apply
kubectl apply -k .
```

### Deploy with Custom Image

```bash
# Using kustomize to override image
cd /home/sklee/4g-proxy-private/deploy/kubernetes

kubectl kustomize . | \
  sed 's|s1ap-proxy:latest|myregistry.com/s1ap-proxy:v1.0.0|g' | \
  kubectl apply -f -
```

Or edit `kustomization.yaml`:

```yaml
images:
- name: s1ap-proxy
  newName: myregistry.com/s1ap-proxy
  newTag: v1.0.0
```

### Verify Deployment

```bash
# Check pods
kubectl get pods -n s1ap-proxy

# Check service
kubectl get svc -n s1ap-proxy

# View logs
kubectl logs -f deployment/s1ap-proxy -n s1ap-proxy
```

---

## Configuration

### Update MME Address

Edit the ConfigMap:

```bash
kubectl edit configmap s1ap-proxy-config -n s1ap-proxy
```

Or patch it:

```bash
kubectl patch configmap s1ap-proxy-config -n s1ap-proxy \
  --type merge \
  -p '{"data":{"MME_ADDRESS":"new-mme.example.com"}}'
```

Then restart the pod:

```bash
kubectl rollout restart deployment/s1ap-proxy -n s1ap-proxy
```

### Update Delays

Edit the ConfigMap:

```bash
kubectl patch configmap s1ap-proxy-config -n s1ap-proxy \
  --type merge \
  -p '{"data":{"DELAY_ATTACH_MS":"2000","DELAY_TAU_MS":"1000"}}'

# Restart to apply
kubectl rollout restart deployment/s1ap-proxy -n s1ap-proxy
```

### Runtime Configuration via API

You can also change delays at runtime without restarting:

```bash
# Port forward to access API
kubectl port-forward svc/s1ap-proxy 8080:8080 -n s1ap-proxy &

# Set delay
curl -X PUT http://localhost:8080/api/v1/delay/attach \
  -H "Content-Type: application/json" \
  -d '{"delayMs": 2000}'
```

---

## Accessing the Proxy

### From Inside the Cluster

Other pods can reach the proxy at:
- **S1AP**: `s1ap-proxy.s1ap-proxy.svc.cluster.local:36412`
- **API**: `s1ap-proxy.s1ap-proxy.svc.cluster.local:8080`

### From Outside the Cluster (NodePort)

Uncomment the NodePort service in `service.yaml` and apply:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: s1ap-proxy-external
  namespace: s1ap-proxy
spec:
  type: NodePort
  selector:
    app.kubernetes.io/name: s1ap-proxy
  ports:
  - name: s1ap
    port: 36412
    targetPort: s1ap
    protocol: SCTP
    nodePort: 36412
```

Then access via any node IP:
- **S1AP**: `<node-ip>:36412`

### Using LoadBalancer (Cloud)

If your cloud provider supports SCTP LoadBalancers:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: s1ap-proxy-lb
  namespace: s1ap-proxy
spec:
  type: LoadBalancer
  selector:
    app.kubernetes.io/name: s1ap-proxy
  ports:
  - name: s1ap
    port: 36412
    protocol: SCTP
```

### Port Forward for API Access

```bash
# Forward API port
kubectl port-forward svc/s1ap-proxy 8080:8080 -n s1ap-proxy

# Now access API locally
curl http://localhost:8080/api/v1/status
```

---

## Updating Configuration

### Method 1: Edit ConfigMap + Restart

```bash
# Edit
kubectl edit configmap s1ap-proxy-config -n s1ap-proxy

# Restart
kubectl rollout restart deployment/s1ap-proxy -n s1ap-proxy
```

### Method 2: Runtime API (No Restart)

```bash
# Port forward
kubectl port-forward svc/s1ap-proxy 8080:8080 -n s1ap-proxy &

# Set delays
curl -X PUT http://localhost:8080/api/v1/delay \
  -H "Content-Type: application/json" \
  -d '{"attach": 1000, "tau": 500}'

# Set drops
curl -X PUT http://localhost:8080/api/v1/drop \
  -H "Content-Type: application/json" \
  -d '{"attach": true}'
```

---

## Monitoring

### Check Pod Status

```bash
kubectl get pods -n s1ap-proxy -w
```

### View Logs

```bash
# Stream logs
kubectl logs -f deployment/s1ap-proxy -n s1ap-proxy

# Last 100 lines
kubectl logs --tail=100 deployment/s1ap-proxy -n s1ap-proxy
```

### Check Stats via API

```bash
kubectl port-forward svc/s1ap-proxy 8080:8080 -n s1ap-proxy &
curl http://localhost:8080/api/v1/stats
```

### Resource Usage

```bash
kubectl top pod -n s1ap-proxy
```

---

## Troubleshooting

### Pod Not Starting

```bash
# Check pod status
kubectl describe pod -l app.kubernetes.io/name=s1ap-proxy -n s1ap-proxy

# Check events
kubectl get events -n s1ap-proxy --sort-by='.lastTimestamp'
```

### SCTP Not Working

1. Check if SCTP is enabled on nodes:
```bash
# On each node
lsmod | grep sctp
# If not loaded: sudo modprobe sctp
```

2. Check network policy allows SCTP

3. Some CNI plugins don't support SCTP. Check your CNI documentation.

### Can't Connect to MME

```bash
# Exec into pod and test connectivity
kubectl exec -it deployment/s1ap-proxy -n s1ap-proxy -- sh

# Inside the pod
nc -z <mme-address> 36412
```

### Image Pull Errors

```bash
# Check if imagePullSecrets are configured
kubectl get deployment s1ap-proxy -n s1ap-proxy -o yaml | grep imagePullSecrets

# Create secret for private registry
kubectl create secret docker-registry regcred \
  --docker-server=myregistry.com \
  --docker-username=user \
  --docker-password=pass \
  -n s1ap-proxy
```

Then add to deployment:
```yaml
spec:
  template:
    spec:
      imagePullSecrets:
      - name: regcred
```

### Health Check Failing

```bash
# Check if API is responding
kubectl exec -it deployment/s1ap-proxy -n s1ap-proxy -- wget -qO- http://localhost:8080/health
```

---

## Complete Example

```bash
# 1. Build image
docker build -t myregistry.com/s1ap-proxy:v1.0.0 .

# 2. Push to registry
docker push myregistry.com/s1ap-proxy:v1.0.0

# 3. Update kustomization.yaml with your image
cd deploy/kubernetes
cat > kustomization.yaml << EOF
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: s1ap-proxy
resources:
- namespace.yaml
- configmap.yaml
- deployment.yaml
- service.yaml
images:
- name: s1ap-proxy
  newName: myregistry.com/s1ap-proxy
  newTag: v1.0.0
EOF

# 4. Update ConfigMap with your MME address
sed -i 's/mme-service.lte-core.svc.cluster.local/your-mme-address/' configmap.yaml

# 5. Deploy
kubectl apply -k .

# 6. Verify
kubectl get pods -n s1ap-proxy
kubectl logs -f deployment/s1ap-proxy -n s1ap-proxy
```
