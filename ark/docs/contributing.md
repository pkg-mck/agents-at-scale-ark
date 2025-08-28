# Contributing

Standard Kubernetes operator development practices:

1. **Development**: Use `make run` to run the controller locally against your kubeconfig cluster
2. **Testing**: Run `make test` for unit tests and `make test-e2e` for end-to-end tests  
3. **Code Quality**: Use `make lint` to check code style and `make fmt vet` to format code
4. **Documentation**: Update both README.md and CLAUDE.md when making significant changes

## Development Workflow

Generate manifests and code after API changes:
```bash
make manifests generate
```

Run locally for development:
```bash
make run
```

Build and test:
```bash
make build test lint
```

## Make Targets

Run `make help` for more information on all potential `make` targets.

## Testing Framework

- Uses Ginkgo/Gomega for testing framework
- E2E tests configured with Kind clusters
- Test suite setup in `internal/controller/suite_test.go`
- Environment test setup via controller-runtime's envtest

## Required Tools

The Makefile manages these tool dependencies:
- `controller-gen`: Generate CRDs and RBAC
- `kustomize`: Build Kubernetes manifests  
- `setup-envtest`: Set up test environment
- `golangci-lint`: Code linting
- `kind`: E2E test clusters

## Additional Resources

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)