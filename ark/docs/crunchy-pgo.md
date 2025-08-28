# Crunchy PostgreSQL Operator Integration

## Installation

```bash
make install-pgo
```

## PostgreSQL Cluster

### Basic Configuration
```yaml
postgres:
  enabled: true
  version: 15
  replicas: 1
  storage: "1Gi"
```

### Production Configuration
```yaml
postgres:
  enabled: true
  version: 15
  replicas: 3
  storage: "50Gi"
  backupStorage: "20Gi"
  resources:
    limits:
      cpu: 2000m
      memory: 4Gi
    requests:
      cpu: 1000m
      memory: 2Gi
```

## Features

### High Availability
- Multi-replica PostgreSQL clusters
- Automatic failover
- Connection pooling

### Backup & Recovery
- pgBackRest integration
- Scheduled backups
- Point-in-time recovery

### Security
- TLS encryption
- User management
- Secret generation

## Generated Resources

### PostgresCluster
- Primary PostgreSQL cluster
- Backup repository configuration
- User and database setup

### Secrets
- Connection strings in `{name}-postgres-pguser-{username}`
- Automatic credential rotation
- TLS certificates

## Monitoring

### Metrics
- PostgreSQL metrics via pgMonitor
- Custom dashboards
- Alerting rules

### Health Checks
- Database connectivity
- Replication lag
- Backup status

## Commands

```bash
# Install operator
make install-pgo

# Remove operator
make uninstall-pgo

# Deploy with PGO
helm install postgres-memory ./services/postgres-memory/chart \
  --set postgres.enabled=true \
  --set postgres.storage=10Gi

# Scale replicas
kubectl patch postgrescluster postgres-memory-postgres \
  --type='merge' -p='{"spec":{"instances":[{"replicas":3}]}}'
```