#=====================================================================================================
#Define the dependencies for local workstation (Linux (Debian or Ubuntu)only, must be AMD64/x_86 or ARM64), Windows & Macs are not supported at the moment)
GOLANG          := golang:1.23
ALPINE          := alpine:3.20
KIND            := kindest/node:v1.31.0
POSTGRES        := postgres:16.4
GRAFANA         := grafana/grafana:11.1.0
PROMETHEUS      := prom/prometheus:v2.54.0
TEMPO           := grafana/tempo:2.5.0
LOKI            := grafana/loki:3.1.0
PROMTAIL        := grafana/promtail:3.1.0
KIND_CLUSTER    := local-cluster
NAMESPACE       := adoption-system
ADOPT_APP       := adoptadog
AUTH_APP        := auth
BASE_IMAGE_NAME := localhost/rmishgoog
VERSION         := 0.0.1
ADOPT_IMAGE     := $(BASE_IMAGE_NAME)/$(ADOPT_APP):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/metrics:$(VERSION)
AUTH_IMAGE      := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)
CILIUM_CLI      := v0.16.15
CILIUM_VERSION  := v1.16.0
GOOS            := $(shell go env GOOS)
GOARCH          := $(shell go env GOARCH)
#=====================================================================================================
#Install environment dependencies

dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

dev-kubernetes-amd64:
	curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.24.0/kind-linux-amd64 && \
	chmod +x ./kind && \
	sudo mv ./kind /usr/local/bin/kind && \
	wait;

dev-kubernetes-arm64:
	curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.24.0/kind-linux-arm64 && \
	chmod +x ./kind && \
	sudo mv ./kind /usr/local/bin/kind && \
	wait;

dev-kubernetes-tooling: system-update	kubectl	kustomize

system-update:
	sudo apt-get update && \
	sudo apt-get install -y apt-transport-https ca-certificates curl gnupg && \
	curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.31/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
	sudo chmod 644 /etc/apt/keyrings/kubernetes-apt-keyring.gpg && \
	echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.31/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list && \
	sudo chmod 644 /etc/apt/sources.list.d/kubernetes.list

kubectl:
	sudo apt-get update && \
	sudo apt-get install -y kubectl

kustomize:
	curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"  | bash
	sudo mv kustomize /usr/local/bin/kustomize


dev-docker:
	docker pull $(GOLANG)& \
	docker pull $(ALPINE)& \
	docker pull $(KIND)& \
	docker pull $(POSTGRES)& \
	docker pull $(PROMETHEUS)& \
	docker pull $(GRAFANA)& \
	docker pull $(LOKI)& \
	docker pull $(PROMTAIL)& \
	docker pull $(TEMPO)
	wait;

#=====================================================================================================
#Prepare the Kubernetes environment & manage the cluster

dev-cluster-up:
	kind create cluster --name $(KIND_CLUSTER) --image $(KIND) --config zarf/k8s/dev/kind-config.yaml

	kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER) & \
	kind load docker-image $(GRAFANA) --name $(KIND_CLUSTER) & \
	kind load docker-image $(PROMETHEUS) --name $(KIND_CLUSTER) & \
	kind load docker-image $(TEMPO) --name $(KIND_CLUSTER) & \
	kind load docker-image $(LOKI) --name $(KIND_CLUSTER) & \
	kind load docker-image $(PROMTAIL) --name $(KIND_CLUSTER) & \
	wait;

dev-cluster-cilium-install:
	curl -L --remote-name-all https://github.com/cilium/cilium-cli/releases/download/$(CILIUM_CLI)/cilium-$(GOOS)-$(GOARCH).tar.gz{,.sha256sum} & \
	wait;


	sudo tar -C /usr/local/bin -xzvf cilium-$(GOOS)-$(GOARCH).tar.gz & \
	wait;

	rm cilium-$(GOOS)-$(GOARCH).tar.gz && rm cilium-linux-amd64.tar.gz.sha256sum
	
	cilium install --version $(CILIUM_VERSION)   --set encryption.enabled=true   --set encryption.type=wireguard   --set encryption.nodeEncryption=true& \
	cilium status --wait

dev-cluster-cilium-uninstall:
	cilium uninstall --all

dev-cluster-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-cluster-status:
	kubectl get nodes -o wide
	kubectl get po -o wide --all-namespaces --watch

dev-pods-status:
	watch -n 5 kubectl get pods -o wide --all-namespaces
#=====================================================================================================
#Local execution from command line & go moduling related commands

run:
	go run apis/services/adoptions/main.go

run-logfmt:
	go run apis/services/adoptions/main.go | go run apis/tooling/logfmt/main.go

tidy:
	go mod tidy
	go mod vendor
#=====================================================================================================
#Build the application images

build: adoptadog-image	adoptadog-image-upload

adoptadog-image:
	docker build \
		-t $(ADOPT_IMAGE) \
		-f zarf/docker/dockerfile.adopt \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

adoptadog-image-upload:
	kind load docker-image $(ADOPT_IMAGE) --name $(KIND_CLUSTER) & \
	wait;
#=====================================================================================================