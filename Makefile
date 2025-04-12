lint:
	~/go/bin/golangci-lint run ./... -c .golangci.yaml

# Dev
# Backend dev dir
backend_dev_dir = deploy/dev

dev-infrastructure-run:
	docker compose -f $(backend_dev_dir)/infrastructure/docker-compose.yaml -p go-project-layout-dev-infrastructure up -d

dev-infrastructure-down:
	docker compose -f $(backend_dev_dir)/infrastructure/docker-compose.yaml -p go-project-layout-dev-infrastructure down

dev-services-build:
	docker compose -f $(backend_dev_dir)/docker-compose.yaml -p go-project-layout-dev-services build

dev-services-run:
	docker compose -f $(backend_dev_dir)/docker-compose.yaml -p go-project-layout-dev-services up -d

dev-services-down:
	docker compose -f $(backend_dev_dir)/docker-compose.yaml -p go-project-layout-dev-services down