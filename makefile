#=============================================================================================================================================================
#Define the dependencies for local workstation (Linux (Debian or Ubuntu)only, must be AMD64/x_86 or ARM64), Windows & Macs are not supported at the moment)
GOLANG            := golang:1.23
ALPINE            := alpine:3.20
KIND              := kindest/node:v1.31.0
POSTGRES          := postgres:16.4
GRAFANA           := grafana/grafana:11.1.0
PROMETHEUS        := prom/prometheus:v2.54.0
KEYCLOAK          := quay.io/keycloak/keycloak:25.0.4
TEMPO             := grafana/tempo:2.5.0
LOKI              := grafana/loki:3.1.0
PROMTAIL          := grafana/promtail:3.1.0
TRAEFIK           := traefik:v3.1.3
KIND_CLUSTER      := local-cluster
NAMESPACE         := adoption-system
TRAEFIK_NAMESPACE := traefik-system
KEYCLOAK_NAMESPACE:= keycloak-system
AUTH_NAMESPACE	  := auth-system
ADOPT_APP         := adoptadog
ADOPT_DEPLOY      := adoptions
TRAEFIK_APP	      := traefik
TRAEFIK_DEPLOY    := traefik-proxy
KEYCLOAK_APP      := keycloak
AUTH_APP          := auth-server
BASE_IMAGE_NAME   := localhost/rmishgoog
VERSION           := 0.0.1
ADOPT_IMAGE       := $(BASE_IMAGE_NAME)/$(ADOPT_APP):$(VERSION)
METRICS_IMAGE     := $(BASE_IMAGE_NAME)/metrics:$(VERSION)
AUTH_IMAGE        := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)
CILIUM_CLI        := v0.16.15
CILIUM_NS	      := kube-system
CILIUM_VERSION    := v1.16.0
GOOS              := $(shell go env GOOS)
GOARCH            := $(shell go env GOARCH)
HUBBLE_ARCH       := amd64
HUBBLE_VERSION    := $(shell curl -s https://raw.githubusercontent.com/cilium/hubble/master/stable.txt)
#=============================================================================================================================================================
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
	docker pull $(TEMPO)& \
	docker pull $(KEYCLOAK)& \
	docker pull $(TRAEFIK)& \
	wait;

#=============================================================================================================================================================
#Prepare the Kubernetes environment & manage the cluster (AMD64/x_86 only)

dev-bootstrap-kind:	dev-cluster-up	dev-docker-loads

dev-cluster-up:
	kind create cluster --name $(KIND_CLUSTER) --image $(KIND) --config zarf/k8s/dev/kind-config.yaml

dev-docker-loads:
	kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER) & \
	kind load docker-image $(GRAFANA) --name $(KIND_CLUSTER) & \
	kind load docker-image $(PROMETHEUS) --name $(KIND_CLUSTER) & \
	kind load docker-image $(TEMPO) --name $(KIND_CLUSTER) & \
	kind load docker-image $(LOKI) --name $(KIND_CLUSTER) & \
	kind load docker-image $(PROMTAIL) --name $(KIND_CLUSTER) & \
	kind load docker-image $(KEYCLOAK) --name $(KIND_CLUSTER) & \
	kind load docker-image $(TRAEFIK) --name $(KIND_CLUSTER) & \
	wait;

#=============================================================================================================================================================
#Install Cilim and Hubble with Helm, this will be the approach moving forward
dev-install-cilium: dev-cluster-helm-cilium-repo-update	dev-cluster-helm-cilium-repo-add	dev-cluster-helm-cilium-install

dev-cluster-helm-cilium-repo-update:
	helm repo update

dev-cluster-helm-cilium-repo-add:
	helm repo add cilium https://helm.cilium.io/

dev-cluster-helm-cilium-install:
	helm install cilium cilium/cilium --version 1.16.3 \
  --namespace kube-system \
  --set prometheus.enabled=true \
  --set l7Proxy=true \
  --set operator.prometheus.enabled=true \
  --set hubble.enabled=true \
  --set hubble.relay.enabled=true \
  --set hubble.metrics.enabled="{dns,drop,tcp,flow,port-distribution,icmp,httpV2:exemplars=true;labelsContext=source_ip\,source_namespace\,source_workload\,destination_ip\,destination_namespace\,destination_workload\,traffic_direction}"

dev-cilium-encryption-enable: dev-cilium-enable-wireguard	dev-cilium-restart-ds

dev-cilium-enable-wireguard:
	helm upgrade cilium cilium/cilium --namespace kube-system \
  --reuse-values \
  --set l7Proxy=false \
  --set encryption.enabled=true \
  --set encryption.type=wireguard

dev-cilium-restart-ds:
	kubectl rollout restart daemonset/cilium -n kube-system

#=============================================================================================================================================================
#Install the Cilium & Hubble for the Kubernetes cluster!

dev-cluster-cilium-install:
	curl -L --remote-name-all https://github.com/cilium/cilium-cli/releases/download/$(CILIUM_CLI)/cilium-$(GOOS)-$(GOARCH).tar.gz{,.sha256sum} & \
	wait;


	sudo tar -C /usr/local/bin -xzvf cilium-$(GOOS)-$(GOARCH).tar.gz & \
	wait;

	rm cilium-$(GOOS)-$(GOARCH).tar.gz && rm cilium-linux-amd64.tar.gz.sha256sum
	
	cilium install --version $(CILIUM_VERSION) --namespace $(CILIUM_NS)   --set encryption.enabled=true   --set encryption.type=wireguard   --set encryption.nodeEncryption=true& \
	cilium status --wait

dev-cluster-hubble-cli-install:
	curl -L --fail --remote-name-all https://github.com/cilium/hubble/releases/download/$(HUBBLE_VERSION)/hubble-linux-$(HUBBLE_ARCH).tar.gz{,.sha256sum} & \
	wait;

	sudo tar xzvfC hubble-linux-$(HUBBLE_ARCH).tar.gz /usr/local/bin & \
	wait;

	rm hubble-linux-$(HUBBLE_ARCH).tar.gz && rm hubble-linux-$(HUBBLE_ARCH).tar.gz.sha256sum & \
	hubble version

	cilium hubble enable --ui & \
	cilium status --wait

dev-cluster-hubble-show: dev-cluster-hubble-port-forward dev-cluster-hubble-status

dev-cluster-hubble-port-forward:
	cilium hubble port-forward &

dev-cluster-hubble-status:
	hubble status

#=============================================================================================================================================================
#Tear down the Kubernetes cluster
dev-cluster-down:
	kind delete cluster --name $(KIND_CLUSTER)
#=============================================================================================================================================================
dev-cluster-status:
	kubectl get nodes -o wide
	kubectl get po -o wide --all-namespaces --watch

dev-pods-status:
	watch -n 5 kubectl get pods -o wide --all-namespaces
#=============================================================================================================================================================
#Local execution from command line & go moduling related commands

run:
	go run apis/services/adoptions/cmd/main.go

run-logfmt:
	go run apis/services/adoptions/cmd/main.go | go run apis/tooling/logfmt/main.go

tidy:
	go mod tidy
	go mod vendor
#=============================================================================================================================================================
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

build-auth: auth-server-image	auth-server-image-upload

auth-server-image:
	docker build \
		-t $(AUTH_IMAGE) \
		-f zarf/docker/dockerfile.auth \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

auth-server-image-upload:
	kind load docker-image $(AUTH_IMAGE) --name $(KIND_CLUSTER) & \
	wait;
#=============================================================================================================================================================
#Deploy/un-deploy the applications to the Kubernetes cluster & install traefik proxy, keycloak & other services

dev-apply:
	kustomize build zarf/k8s/dev/adoptions | kubectl apply -f -
	kubectl wait pods --for=condition=Ready --timeout=120s -n $(NAMESPACE) --selector app=$(ADOPT_DEPLOY)

dev-unapply:
	kustomize build zarf/k8s/dev/adoptions | kubectl delete -f -

dev-apply-traefik:
	kustomize build zarf/k8s/dev/traefik | kubectl apply -f -
	kubectl wait pods --for=condition=Ready --timeout=120s -n $(TRAEFIK_NAMESPACE) --selector app=$(TRAEFIK_APP)

dev-unapply-traefik:
	kustomize build zarf/k8s/dev/traefik | kubectl delete -f -

dev-apply-keycloak:
	kustomize build zarf/k8s/dev/keycloak | kubectl apply -f -
	kubectl wait pods --for=condition=Ready --timeout=900s -n $(KEYCLOAK_NAMESPACE) --selector app=$(KEYCLOAK_APP)

dev-unapply-keycloak:
	kustomize build zarf/k8s/dev/keycloak | kubectl delete -f -

dev-apply-auth:
	kustomize build zarf/k8s/dev/auth | kubectl apply -f -
	kubectl wait pods --for=condition=Ready --timeout=30s -n $(AUTH_NAMESPACE) --selector app=$(AUTH_APP)

dev-unapply-auth:ustomize build za
	kustomize build zarf/k8s/dev/auth | kubectl delete -f -

#=============================================================================================================================================================
#Restart the kubernetes deployments & get status
dev-restart:
	kubectl rollout restart deployment $(ADOPT_DEPLOY)  -n $(NAMESPACE)

dev-restart-keycloak:
	kubectl rollout restart deployment $(KEYCLOAK_APP)  -n $(KEYCLOAK_NAMESPACE)

dev-restart-traefik:
	kubectl rollout restart deployment $(TRAEFIK_DEPLOY)  -n $(TRAEFIK_NAMESPACE)

dev-restart-auth:
	kubectl rollout restart deployment $(AUTH_APP)  -n $(AUTH_NAMESPACE)

dev-deploy-status:
	kubectl rollout status deployment $(ADOPT_DEPLOY) -n $(NAMESPACE)

dev-deploy-traefik-status:
	kubectl rollout status deployment $(TRAEFIK_DEPLOY) -n $(TRAEFIK_NAMESPACE)

dev-deploy-keycloak-status:
	kubectl rollout status deployment $(KEYCLOAK_APP) -n $(KEYCLOAK_NAMESPACE)

dev-deploy-auth-status:
	kubectl rollout status deployment $(AUTH_APP) -n $(AUTH_NAMESPACE)

#=============================================================================================================================================================
#Get the logs from the adoption application pods
dev-logs-fmtd:
	kubectl logs -f -l app=$(ADOPT_DEPLOY) -n $(NAMESPACE) --tail=100 --max-log-requests=6 | go run apis/tooling/logfmt/main.go

dev-logs-raw:
	kubectl logs -f -l app=adoptions -n $(NAMESPACE)
#=============================================================================================================================================================
#Build & deploy the application from scratch

dev-build-deploy: build dev-apply dev-deploy-status
#=============================================================================================================================================================
#Build, upload image & restrat the deployment, no KRM changes

dev-build-restrat: build dev-restart
#=============================================================================================================================================================
#Describe the application deployment & pods

dev-describe-deployment:
	kubectl describe deployment $(ADOPT_DEPLOY) -n $(NAMESPACE)

dev-describe-pods:
	kubectl describe pods -n $(NAMESPACE) -l app=$(ADOPT_DEPLOY)
#=============================================================================================================================================================
# Operations for Cilium

dev-cilium-status:
	cilium status --wait

dev-reboot-cilium:	dev-cluster-cilium-reinstall dev-cluster-cilium-hubble-enable

dev-cluster-cilium-reinstall:
	cilium uninstall
	wait;

	cilium install --version $(CILIUM_VERSION) --namespace $(CILIUM_NS)   --set encryption.enabled=true   --set encryption.type=wireguard   --set encryption.nodeEncryption=true& \
	cilium status --wait

dev-cluster-cilium-hubble-enable:
	cilium hubble enable --ui & \
	cilium status --wait

dev-cluster-cilium-client:
	cilium version --client

dev-cluster-cilium-server:
	cilium version

dev-cluster-cilium-uninstall:
	cilium uninstall

dev-cluster-cilium-cep:
	kubectl get cep --all-namespaces

dev-cluster-cilium-config-view:
	cilium config view
#=============================================================================================================================================================
#Basic local service testing

dev-kubectl-forward:
	kubectl port-forward svc/adoptions 3000:3000 -n $(NAMESPACE)& >> /dev/null

dev-adoptadog-liveness:
	curl -X GET http://localhost:3000/liveness

dev-adoptadog-readiness:
	curl -X GET http://localhost:3000/readiness

dev-adoptadog-endpoint-load:
	hey -n 1000 -c 10 http://localhost:3000/liveness
#=============================================================================================================================================================
#Setting up kind cluster for LoadBalancer services

#This section is commented out in favor of MetaLB LoadBalancer IPAM & future versions of this application
#shall continue to Cilium as the CNI & MetaLB LoadBalancer service IPAM (non-cloud provider environments)

# kind-configure-cloud-provider-lb:	kind-remove-label-lb-access kind-install-cloud-provider-lb kind-enable-cloud-provider-lb

# kind-remove-label-lb-access:
# 	kubectl label node local-cluster-control-plane node.kubernetes.io/exclude-from-external-load-balancers-

# kind-install-cloud-provider-lb:
# 	go install sigs.k8s.io/cloud-provider-kind@latest
# 	sudo install ~/go/bin/cloud-provider-kind /usr/local/bin

# kind-enable-cloud-provider-lb:
# 	cloud-provider-kind > /dev/null 2>&1 &

# MetalLB LoadBalancer service configuration for the local kind cluster

dev-metallb-install: dev-metallb-apply	dev-configure-metallb

dev-metallb-apply:
	kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.14.8/config/manifests/metallb-native.yaml
	wait;

dev-configure-metallb:
	kustomize build zarf/k8s/dev/metal-lb | kubectl apply -f -
#=============================================================================================================================================================
#Setting up the private key & certificate for the keycloak service
#The following commands will use openssl to create a self-signed certificate for the keycloak service
# openssl genpkey -algorithm RSA -out key.pem
# openssl req -new -x509 -key key.pem -out cert.pem -days 365
# openssl x509 -text -noout -in cert.pem
# After creating the certificate, the following command will create a secret in the Kubernetes cluster
# kubectl create secret tls keycloak-tls-secret --key=key.pem --cert=cert.pem -n keycloak-system
#=============================================================================================================================================================
#Generate an OAuth2 access token for the API
#You need to do the set up before you can use this command, this app has keycloak as the OAuth2 provider
#and the self-signed certificate for the keycloak service use local.auth.adoptadog.com as the host, keycloak
#service is running on port 443 (LoadBalancer service) and the realm is apiclients
# curl -kv POST https://local.auth.adoptadog.com/realms/adoptadog/protocol/openid-connect/token \
# -H 'content-type: application/x-www-form-urlencoded' \
# -d 'client_id=local-test-harness' \
# -d 'username=api-developer&password=api-developer&grant_type=password' | jq --raw-output '.access_token'
dev-get-access-token:
	curl -kv POST https://local.auth.adoptadog.com/realms/adoptadog/protocol/openid-connect/token \
	-H 'content-type: application/x-www-form-urlencoded' \
	-d 'client_id=local-test-harness' \
	-d 'username=api-developer&password=api-developer&grant_type=password' | jq --raw-output '.access_token'
#=============================================================================================================================================================
#Discovery URL to get the JWKS for the realm
#https://local.auth.adoptadog.com/realms/adoptadog/.well-known/openid-configuration
#Tooling for local development
dev-get-jwks:
	curl -kv https://local.auth.adoptadog.com/realms/adoptadog/.well-known/openid-configuration | jq
dev-get-token-verify:
	go run apis/tooling/keys/keys.go