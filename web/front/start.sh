docker stop nuxt
docker rm nuxt
docker rmi front_nuxt
docker-compose -f docker-compose.yml up -d
