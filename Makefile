dc_up:
	docker-compose  -f docker-compose.redis.yaml -f docker-compose.pg.yaml up -d

dc_down:
	docker-compose -f docker-compose.redis.yaml -f docker-compose.pg.yaml down




