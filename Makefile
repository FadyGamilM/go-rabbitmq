producer_img := producer
consumer_img := consumer

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
# Then write this in the bash session 
# $ rabbitmq-plugins enable rabbitmq_management


rabbit-build-producer:
	docker build -t $(producer_img):1.0 -f producer.Dockerfile .

rabbit-producer:
	docker run -it --rm --net rabbits -e RABBIT_HOST=rabbit-1 -e RABBIT_PORT=5672 -e RABBIT_USERNAME=guest -e RABBIT_PASSWORD=guest --name producer_api -p 9876:9876 $(producer_img):1.0 

rabbit-build-consumer:
	docker build -t $(consumer_img):1.0 -f consumer.Dockerfile .

rabbit-consumer:
	docker run -it --rm --net rabbits -e RABBIT_HOST=rabbit-1 -e RABBIT_PORT=5672 -e RABBIT_USERNAME=guest -e RABBIT_PASSWORD=guest --name consumer_service $(consumer_img):1.0 

	