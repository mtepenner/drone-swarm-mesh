.PHONY: build validate up scale

build:
	docker compose -f orchestration/docker-compose.yaml build

validate:
	python -m compileall swarm_simulator/app
	cd drone_agent && go test ./...
	cd swarm_visualizer && npm install
	cd swarm_visualizer && npm run build
	if command -v bash >/dev/null 2>&1; then bash -n orchestration/scale_swarm.sh; fi

up:
	docker compose -f orchestration/docker-compose.yaml up -d --build

scale:
	./orchestration/scale_swarm.sh 25
