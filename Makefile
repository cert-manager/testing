# Copyright 2023 The cert-manager Authors.
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

.PHONY: verify-boilerplate
verify-boilerplate:
	go run github.com/cert-manager/boilersuite@v0.1.0 .

.PHONY: prowgen
prowgen:
	cd ./config/prowgen/ && go run . --branch=* -o $(CURDIR)/config/jobs/cert-manager/cert-manager/

.PHONY: verify-prowgen
verify-prowgen:
	mkdir -p _temp
	cd ./config/prowgen/ && go run . --branch=* -o $(CURDIR)/_temp/cert-manager/
	diff -q -r $(CURDIR)/config/jobs/cert-manager/cert-manager/ $(CURDIR)/_temp/cert-manager/

.PHONY: verify
verify: verify-boilerplate verify-prowgen

.PHONY: test
test:
	cd ./config/prowgen/ && go test ./...

# Run checkconfig locally to verify the Prow configuration, CI runs this
# directly in the Prow cluster.
local-checkconfig:
	docker run --rm \
		-v $(CURDIR)/config:/config \
		us-docker.pkg.dev/k8s-infra-prow/images/checkconfig:v20260119-643bbf68a \
		--strict=true \
        --config-path=/config/config.yaml \
        --job-config-path=/config/jobs \
        --plugin-config=/config/plugins.yaml

	docker run --rm \
		-v $(CURDIR)/config:/config \
		gcr.io/k8s-staging-test-infra/configurator:v20251212-3f863f26d7 \
        --yaml=/config/testgrid/dashboards.yaml \
        --default=config/testgrid/default.yaml \
        --prow-config=/config/config.yaml \
        --prow-job-config=/config/jobs \
        --prowjob-url-prefix=https://github.com/cert-manager/testing/tree/master/config/jobs \
        --update-description \
        --validate-config-file \
        --oneshot
