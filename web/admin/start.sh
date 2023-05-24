docker stop vue
docker rm vue
docker rmi admin_vue
docker-compose -f docker-compose.yml up -d
