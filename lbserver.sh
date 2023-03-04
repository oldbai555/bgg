docker stop web
docker rm web
docker rmi bgg_web
docker-compose -f docker-compose.yml up -d
