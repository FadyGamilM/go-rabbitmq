rabbit-network:
	docker network create rabbits

rabbit-instance:
# --hostname rabbit-1 ==> is the host name of our rabbit instance 
# --name rabbit-1 ==> is the name of our container 
# its important to know the hostname because rabbitmq uses <identifier>@<hostname> so rabbitmq instance can talk to each others 
	docker run -d --rm --net rabbits --hostname rabbit-1 -p 15672:15672 --name rabbit-1 rabbitmq:3-management

rabbit-logs:
	docker logs rabbit-1

rabbit-bash:
	docker exec -it rabbit-1 bash

	