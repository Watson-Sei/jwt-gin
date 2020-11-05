# jwt-gin
JWT認証をGINに実装する

## setting
```shell script
# docker/db/.env

MYSQL_ROOT_PASS=xxxx
MYSQL_DATABASE=xxxx
MYSQL_USER=xxxx
MYSQL_PASSWORD=xxxx
```

## start
```shell script
# ~/home
docker-compose build --no-cache
docker-compose up -d
```

## logs check
```shell script
# ~/home
docker-compose logs -f
```

## User add check Database
```shell script
# ~/home
docker exec -it jwt-gin_db_1 bash
mysql -u xxxx -p
password: xxxx
use xxxx;
select * from user_model;
```