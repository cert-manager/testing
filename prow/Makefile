# Copyright 2021 The cert-manager Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# based on: https://github.com/kubernetes/test-infra/blob/e4d1738d6eb8c2c00f9d90ed9e694e48f14156c5/config/prow/Makefile

# GKE cluster configuration
TRUSTED_PROJECT := cert-manager-tests-trusted
UNTRUSTED_PROJECT := cert-manager-tests-untrusted
ZONE := europe-west1-b
TRUSTED_CLUSTER := prow-trusted
UNTRUSTED_CLUSTER := prow-untrusted
TRUSTED_CTX := gke_$(TRUSTED_PROJECT)_$(ZONE)_$(TRUSTED_CLUSTER)
UNTRUSTED_CTX := gke_$(UNTRUSTED_PROJECT)_$(ZONE)_$(UNTRUSTED_CLUSTER)

# Helper variables
KUBECTL_TRUSTED := kubectl --context=$(TRUSTED_CTX)
KUBECTL_UNTRUSTED := kubectl --context=$(UNTRUSTED_CTX)

.PHONY: help
help:
	@echo "Main targets:"
	@echo "  get-kubeconfig       - fetch kubeconfig for both clusters"
	@echo "  create-cookie-secret - create cookie secret with random data"
	@echo "  verify               - verify required secrets and resources exist"
	@echo "  diff                 - diff all resources"
	@echo "  deploy               - verify, deploy all, and restart deployments"
	@echo "  bootstrap-config     - update config/plugins/job-config configmaps"

.PHONY: get-kubeconfig
get-kubeconfig:
	@echo "Fetching kubeconfig..."
	@gcloud container clusters get-credentials $(UNTRUSTED_CLUSTER) --zone=$(ZONE) --project=$(UNTRUSTED_PROJECT)
	@gcloud container clusters get-credentials $(TRUSTED_CLUSTER) --zone=$(ZONE) --project=$(TRUSTED_PROJECT)

.PHONY: create-cookie-secret
create-cookie-secret:
	@echo "Creating cookie secret..."
	@$(KUBECTL_TRUSTED) create secret generic cookie --from-literal=secret=$$(openssl rand -base64 32) --dry-run=client -o yaml | $(KUBECTL_TRUSTED) apply -f -

.PHONY: bootstrap-config
bootstrap-config:
	@echo "Bootstrapping config to $(TRUSTED_CLUSTER)..."
	@kubectl config use-context $(TRUSTED_CTX) >/dev/null
	@cd ../ && go run sigs.k8s.io/prow/cmd/config-bootstrapper@v0.0.0-20260119105717-643bbf68a1ed \
		--dry-run=false --source-path=. --config-path=config/config.yaml \
		--plugin-config=config/plugins.yaml --job-config-path=config/jobs

.PHONY: verify
verify:
	@echo "Verifying prerequisites..."
	@$(KUBECTL_TRUSTED) get ns default >/dev/null || (echo "ERROR: Cannot connect to trusted cluster" && exit 1)
	@$(KUBECTL_UNTRUSTED) get ns default >/dev/null || (echo "ERROR: Cannot connect to untrusted cluster" && exit 1)
	
	@for secret in cert-manager-bot-github-token hmac-token github-app-token github-oauth-config; do \
		$(KUBECTL_TRUSTED) get secret $$secret -n default >/dev/null || (echo "ERROR: $$secret not found" && exit 1); \
	done
	@for cm in config plugins job-config; do \
		$(KUBECTL_TRUSTED) get cm $$cm -n default >/dev/null || (echo "ERROR: $$cm configmap not found" && exit 1); \
	done
	@for secret in cert-manager-bot-ssh-keys cert-manager-bot-gcp-service-account cert-manager-bot-github-token github-app-token; do \
		$(KUBECTL_TRUSTED) get secret $$secret -n test-pods >/dev/null 2>&1 || \
			echo "WARNING: $$secret not found"; \
	done

	@for secret in cloudflare-api-key venafi-tpp venafi-cloud venafi-ngts; do \
		$(KUBECTL_UNTRUSTED) get secret $$secret -n test-pods >/dev/null 2>&1 || \
			echo "WARNING: $$secret not found"; \
	done
	
	@echo "✓ Prerequisites verified"

.PHONY: diff
diff: diff-prow diff-test-pods diff-trusted-cluster

.PHONY: diff-prow
diff-prow:
	@$(KUBECTL_TRUSTED) diff -f ./cluster/ 2>/dev/null || true

.PHONY: diff-test-pods
diff-test-pods:
	@$(KUBECTL_TRUSTED) diff -f ./test-pods/ 2>/dev/null || true
	@$(KUBECTL_UNTRUSTED) diff -f ./test-pods/ 2>/dev/null || true

	@$(KUBECTL_TRUSTED) diff -f ./test-pods/trusted_cluster/ 2>/dev/null || true
	@$(KUBECTL_UNTRUSTED) diff -f ./test-pods/untrusted_cluster/ 2>/dev/null || true

.PHONY: deploy
deploy: verify deploy-prow deploy-test-pods deploy-trusted-cluster restart-deployments

.PHONY: deploy-prow
deploy-prow:
	@echo "Deploying Prow..."
	@$(KUBECTL_TRUSTED) apply --server-side -f ./cluster/prowjob-crd/prowjob_customresourcedefinition.yaml
	@$(KUBECTL_TRUSTED) apply --server-side -f ./cluster/

.PHONY: deploy-test-pods
deploy-test-pods:
	@echo "Deploying test-pods RBAC..."
	@$(KUBECTL_TRUSTED) apply --server-side -f ./test-pods/
	@$(KUBECTL_UNTRUSTED) apply --server-side -f ./test-pods/

	@$(KUBECTL_TRUSTED) apply --server-side -f ./test-pods/trusted_cluster/
	@$(KUBECTL_UNTRUSTED) apply --server-side -f ./test-pods/untrusted_cluster/

.PHONY: restart-deployments
restart-deployments:
	@echo "Restarting deployments..."
	@$(KUBECTL_TRUSTED) rollout restart deployment -n default
