SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# ==============================================================================
# Define dependencies

GOLANG          := golang:1.24
ALPINE          := alpine:3.21
KIND            := kindest/node:v1.32.0
POSTGRES        := postgres:17.2
GRAFANA         := grafana/grafana:11.4.0
PROMETHEUS      := prom/prometheus:v3.0.0
TEMPO           := grafana/tempo:2.6.0
LOKI            := grafana/loki:3.3.0
PROMTAIL        := grafana/promtail:3.3.0

KIND_CLUSTER    := ardan-starter-cluster
NAMESPACE       := sales-system
SALES_APP       := sales
AUTH_APP        := auth
BASE_IMAGE_NAME := localhost/ardanlabs
VERSION         := 0.0.1
SALES_IMAGE     := $(BASE_IMAGE_NAME)/$(SALES_APP):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/metrics:$(VERSION)
AUTH_IMAGE      := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)

# ==============================================================================
# Class Stuff

run:
	go run api/services/sales/main.go | go run api/tooling/logfmt/main.go

help:
	go run api/services/sales/main.go --help

version:
	go run api/services/sales/main.go --version

curl-test:
	curl -il -X GET http://localhost:3000/test

curl-live:
	curl -il -X GET http://localhost:3000/liveness

curl-ready:
	curl -il -X GET http://localhost:3000/readiness

curl-test-error:
	curl -il -X GET http://localhost:3000/test-error

curl-test-panic:
	curl -il -X GET http://localhost:3000/test-panic

admin:
	go run api/tooling/admin/genkey.go

#admin token
#export TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6IjU0YmIyMTY1LTcxZTEtNDFhNi1hZjNlLTdkYTRhMGUxZTJjMSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzZXJ2aWNlIHByb2plY3QiLCJzdWIiOiIyMzFjMGExZS02NjRmLTRiMmMtYThhOS1mM2ZiYWE5YThiZDEiLCJleHAiOjE3NzQxOTg0NTUsImlhdCI6MTc0MjY2MjQ1NSwiUm9sZXMiOlsiQURNSU4iXX0.AXZVVM329HJXdN2bNbaUNG2bynRQNYhTSedzmv7bAgeNDEusBJwDwZI4ZpIEZ0ge6--DNkqsuHhZaQQ5GN2FytJduyogK3nVbVGuLhXh9ALNhPfnyYfcByWTp2tVWKS0uEunpCERuDTK6hLOxSZQ2akQK4_RAADENi90daR49fj8cgJQR-HgcODhLe1Z4fJrc6iyHxXX9jjI2SLaBjocV1DWeZ2YYu5WqlSSc9vI7gRXh3RFUc53Jm3rXVTPKupWxiZTOL1YbDCpP26utFPrW7opbukHIN-BN-hXrnaiESUAvdtL90hqUA73wIg0WVPiG6d39yBLcPVzrrNKvJnX2g
#user token
#export TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6IjU0YmIyMTY1LTcxZTEtNDFhNi1hZjNlLTdkYTRhMGUxZTJjMSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzZXJ2aWNlIHByb2plY3QiLCJzdWIiOiIyMzFjMGExZS02NjRmLTRiMmMtYThhOS1mM2ZiYWE5YThiZDEiLCJleHAiOjE3NzQyMDMwNjksImlhdCI6MTc0MjY2NzA2OSwiUm9sZXMiOlsiVVNFUiJdfQ.mVLlZV3iLmTAb3OTWCjX8e3QTDgU5CiTo_LItVykEdlwEOMWZYUkZfnYxeEIzn4dOUrgMdHuGsSQ7VgbiLV5G8O94U-8s6xI1dc2_HuCg8KzZN2J2XAS_t2sSiOg4orMIkkw1Fxjp3aaAhcKCOSVd8LG-cKvSbbbFyRxGcQ6l169njva0-sv3lFQ9BYk4L6FzFaFd44c0I2-MCYjoT5-1AoLFacb5NqPPSRWsbEUCPAV6DifhIl_TpduZkiPMioI5q5qKHYj-DhMJFf-W6B3YlVrf5rbw1e0WXvTqlPkHmX6csD9JyvJKMAsoekQJxZPZEEA5rOLdt-fIkyCDOLGxA
curl-test-auth:
	curl -il \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:3000/test-auth"

token:
	curl -il \
	--user "admin@example.com:gophers" http://localhost:6000/auth/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1

curl-test-auth-service:
	curl -il \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:6000/auth/authenticate"

# ==============================================================================
# Running from within k8s/kind

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-status-all:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-status:
	watch -n 2 kubectl get pods -o wide --all-namespaces

# ------------------------------------------------------------------------------

dev-load:
	kind load image-archive --name $(KIND_CLUSTER) <(podman save $(SALES_IMAGE))
	kind load image-archive --name $(KIND_CLUSTER) <(podman save $(AUTH_IMAGE))

dev-apply:
	kustomize build zarf/k8s/dev/auth | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(AUTH_APP) --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(SALES_APP) --timeout=120s --for=condition=Ready

dev-restart:
	kubectl rollout restart deployment $(SALES_APP) --namespace=$(NAMESPACE)

dev-restart-auth:
	kubectl rollout restart deployment $(SALES_APP) --namespace=$(NAMESPACE)

dev-run: build dev-up dev-load dev-apply

dev-update: build dev-load dev-restart dev-restart-auth

dev-update-apply: build dev-load dev-apply

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(SALES_APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run api/tooling/logfmt/main.go -service=$(SALES_APP)

dev-logs-auth:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(AUTH_APP) --all-containers=true -f --tail=100 | go run api/cmd/tooling/logfmt/main.go

# ------------------------------------------------------------------------------

dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(SALES_APP)

dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(SALES_APP)

dev-describe-auth:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(AUTH_APP)

# ==============================================================================
# Metrics and Tracing

metrics-view:
	expvarmon -ports="localhost:3010" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

statsviz:
	open http://localhost:3010/debug/statsviz

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor

# ==============================================================================
# Building containers

build: sales auth

sales:
	podman build \
		-f zarf/docker/dockerfile.sales \
		-t $(SALES_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

auth:
	podman build \
		-f zarf/docker/dockerfile.auth \
		-t $(AUTH_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

# ==============================================================================
# Running tests within the local computer

test-r:
	CGO_ENABLED=1 go test -race -count=1 ./...

test-only:
	CGO_ENABLED=0 go test -count=1 ./...

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...

vuln-check:
	govulncheck ./...

test: test-only lint vuln-check

test-race: test-r lint vuln-check
