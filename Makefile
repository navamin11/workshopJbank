docker_rebuild_app:
	docker image prune -f && docker-compose up --build api

docker_stop_app:
	docker-compose stop api && docker image prune -f

docker_up_all:
	docker-compose build --no-cache && docker-compose up

docker_down_all:
	docker-compose down && docker image prune -f
		
docker_ps:
	docker-compose ps

docker_rebuild_all:
	docker-compose down && docker image prune -f && docker-compose build --no-cache && docker-compose up