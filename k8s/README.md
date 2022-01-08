![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/hdb-renderer-syncsign)
[![Actions Status](https://github.com/tommzn/hdb-renderer-syncsign/actions/workflows/go.image.build.yml/badge.svg)](https://github.com/tommzn/hdb-renderer-syncsign/actions)

# HomeDashboard Rendering Server for SyncSign® eInk Displays
Provides a server wich listens for refresh requests send by SyncSign® eInk displays.

## K8s Deployment
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: hdb
  labels:
    name: hdb

---

kind: Deployment
apiVersion: apps/v1
metadata:
  name: hdb-renderer-syncsign
  namespace: hdb
  labels:
    app: hdb-renderer-syncsign
spec:
  replicas: 3
  selector:
    matchLabels:
      app: hdb-renderer-syncsign
  template:
    metadata:
      labels:
        app: hdb-renderer-syncsign
    spec:
      containers:
        - name: hdb-renderer-syncsign
          image: ghcr.io/tommzn/hdb-renderer-syncsign/arm64:latest
          volumeMounts:
            - name: secret-volume
              mountPath: /run/secrets/token
              readOnly: true
          imagePullPolicy: Always

---

kind: Service
apiVersion: v1
metadata:
  name: hdb-datasource-indoorclimate
spec:
  selector:
    app: hdb-datasource-indoorclimate
  ports:
    - port: 8080
```

## Endpoints
As describes in SyncSign render docs this server provides two endpoints.
### Node
Path: /renders/nodes/{nodeid}
Endpoint display will call to get updated content. Server will start rendering for passed node id (=display id) and return new content as JSON. Errors will be renderer as JSON as well and status code will always be 200 OK.
### Render
Path: /renders/{renderid}
Request to this endpoint will always be answered with a 204 status code.
### Health Check
Path: /health
If desired you can observe server health status with this endpoint.

# Links
- [SyncSign](https://sync-sign.com)
- [HomeDashboard Documentation](https://github.com/tommzn/hdb-docs/wiki)
