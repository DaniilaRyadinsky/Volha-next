#!/usr/bin/env bash
set -Eeuo pipefail

cd "${WORKSPACE:-$(git rev-parse --show-toplevel)}"
workspace="$PWD"
reports="$workspace/ci-reports"
cache="$workspace/.ci-cache"
mkdir -p "$reports" "$cache"

# The dind daemon mounts the same Jenkins volume at /var/jenkins_home.
container() {
    docker run --rm --user "$(id -u):$(id -g)" --network host \
        --volume "$workspace:$workspace" --workdir "$workspace" \
        --env HOME=/tmp \
        "$@"
}

go_container() {
    container \
        --env "GOPATH=$cache/go" \
        --env "GOCACHE=$cache/go-build" \
        --env "GOBIN=$cache/go/bin" \
        golang:1.25.4-bookworm bash -euo pipefail -c "$1"
}

case "${1:-}" in
    prepare)
        docker version
        ;;
    backend-build)
        go_container 'go -C backend mod verify && go -C backend build ./...'
        ;;
    frontend-build)
        container --env CI=true --env NEXT_TELEMETRY_DISABLED=1 node:24-bookworm \
            bash -euo pipefail -c 'cd frontend && npx --yes pnpm@11.3.0 install --frozen-lockfile && npx --yes pnpm@11.3.0 build'
        ;;
    unit-tests)
        go_container 'go install gotest.tools/gotestsum@v1.13.0 && cd backend && "${GOBIN}/gotestsum" --junitfile ../ci-reports/unit.xml -- -race -count=1 -covermode=atomic -coverprofile=../ci-reports/coverage.out ./...'
        ;;
    vet)
        go_container 'go -C backend vet ./...'
        ;;
    sonar)
        : "${SONAR_TOKEN:?Create Jenkins Secret text credential sonar-token}"
        container --env SONAR_TOKEN --env SONAR_HOST_URL=http://sonarqube:9000 \
            --env "SONAR_USER_HOME=$cache/sonar" \
            sonarsource/sonar-scanner-cli:11 \
            -Dsonar.qualitygate.wait=true -Dsonar.qualitygate.timeout=600
        ;;
    govulncheck)
        go_container 'go install golang.org/x/vuln/cmd/govulncheck@v1.7.0 && cd backend && "${GOBIN}/govulncheck" ./... 2>&1 | tee ../ci-reports/govulncheck.txt'
        ;;
    osv)
        status=0
        container ghcr.io/google/osv-scanner:v2.3.2 \
            scan source --lockfile=backend/go.mod \
            --format=table > "$reports/osv.txt" || status=$?
        if [[ -f "$reports/osv.txt" ]]; then cat "$reports/osv.txt"; fi
        exit "$status"
        ;;
    gosec)
        status=0
        go_container 'go install github.com/securego/gosec/v2/cmd/gosec@v2.29.0 && cd backend && "${GOBIN}/gosec" -exclude-generated -exclude-dir=api/proto/gen -exclude-dir=docs -severity=high -confidence=medium -fmt=text -out ../ci-reports/gosec.txt ./...' || status=$?
        if [[ -f "$reports/gosec.txt" ]]; then cat "$reports/gosec.txt"; fi
        exit "$status"
        ;;
    semgrep)
        container semgrep/semgrep:1.182.0 semgrep scan --config p/golang \
            --metrics off --disable-version-check --severity ERROR --error --strict \
            --exclude 'backend/api/proto/gen/**' --exclude 'backend/docs/**' \
            --sarif-output "$reports/semgrep.sarif" backend
        ;;
    gitleaks)
        container ghcr.io/gitleaks/gitleaks:v8.30.1 git . --log-opts=HEAD \
            --redact --verbose --report-format json --report-path "$reports/gitleaks.json"
        ;;
    docker-build)
        commit="$(git rev-parse --short=12 HEAD)"
        built_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
        for target in gateway auth-service product-service migrator dumper; do
            docker build --file backend/Dockerfile --target "$target" \
                --build-arg "VERSION=ci-${BUILD_NUMBER:-local}" \
                --build-arg "COMMIT=$commit" --build-arg "BUILD_TIME=$built_at" \
                --tag "volha/$target:${BUILD_NUMBER:-local}" backend
        done
        ;;
    trivy)
        for target in gateway auth-service product-service migrator dumper; do
            status=0
            container --volume /certs/client:/certs/client:ro \
                --env DOCKER_HOST --env DOCKER_CERT_PATH --env DOCKER_TLS_VERIFY \
                aquasec/trivy:0.76.1 image --image-src docker --scanners vuln \
                --severity HIGH,CRITICAL --exit-code 1 --no-progress \
                --format table --output "$reports/trivy-$target.txt" \
                "volha/$target:${BUILD_NUMBER:-local}" || status=$?
            if [[ -f "$reports/trivy-$target.txt" ]]; then cat "$reports/trivy-$target.txt"; fi
            if (( status != 0 )); then exit "$status"; fi
        done
        status=0
        container aquasec/trivy:0.76.1 config --severity HIGH,CRITICAL \
            --exit-code 1 --format table \
            --output "$reports/trivy-config.txt" backend || status=$?
        if [[ -f "$reports/trivy-config.txt" ]]; then cat "$reports/trivy-config.txt"; fi
        exit "$status"
        ;;
    publish)
        : "${NEXUS_USER:?Create Jenkins username/password credential nexus-docker}"
        : "${NEXUS_PASSWORD:?Create Jenkins username/password credential nexus-docker}"
        export DOCKER_CONFIG="$(mktemp -d)"
        trap 'rm -f "$DOCKER_CONFIG/config.json"; rmdir "$DOCKER_CONFIG"' EXIT
        printf '%s' "$NEXUS_PASSWORD" | docker login nexus:8082 \
            --username "$NEXUS_USER" --password-stdin
        for target in gateway auth-service product-service migrator dumper; do
            image="nexus:8082/volha/$target:${BUILD_NUMBER:-local}"
            docker tag "volha/$target:${BUILD_NUMBER:-local}" "$image"
            docker push "$image"
        done
        ;;
    *)
        echo "Usage: $0 {prepare|backend-build|frontend-build|unit-tests|vet|sonar|govulncheck|osv|gosec|semgrep|gitleaks|docker-build|trivy|publish}" >&2
        exit 2
        ;;
esac
