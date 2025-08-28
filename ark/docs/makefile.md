# Makefile Commands

## Build
- `make build` - Build manager binary
- `make docker-build` - Build container image
- `make docker-push` - Push container image

## Test
- `make test` - Run unit tests
- `make test-e2e` - Run end-to-end tests
- `make lint` - Run linter
- `make lint-fix` - Fix linting issues

## Development
- `make run` - Run controller locally with webhooks
- `make dev` - Run controller locally without webhooks
- `make manifests` - Generate CRDs and RBAC
- `make generate` - Generate deepcopy code

## Deployment
- `make install` - Install CRDs
- `make deploy` - Deploy controller (or scale back up if ejected)
- `make eject-controller` - Scale controller to zero for out-of-cluster development
- `make uninstall` - Remove CRDs

## Dependencies
- `make install-cert-manager` - Install cert-manager
- `make uninstall-cert-manager` - Remove cert-manager
- `make install-pgo` - Install Crunchy PostgreSQL Operator
- `make uninstall-pgo` - Remove PostgreSQL Operator

## Development Environment
- `make deploy-dev` - Deploy without cert-manager
- `make admission-dev` - Run with webhooks in dev mode

## Testing
- `make setup-test-e2e` - Create Kind cluster
- `make cleanup-test-e2e` - Delete Kind cluster

## Tools
All tools are downloaded to `bin/` directory:
- controller-gen
- kustomize
- setup-envtest
- golangci-lint