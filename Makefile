producer_img := producer
consumer_img := consumer
rabbit_host_name := rabbit-1

rabbit-network:
	docker network create rabbits

rabbit-instance:
# --hostname rabbit-1 ==> is the host name of our rabbit instance 
# --name rabbit-1 ==> is the name of our container 
# its important to know the hostname because rabbitmq uses <identifier>@<hostname> so rabbitmq instance can talk to each others 
	docker run -d --rm --net rabbits --hostname $(rabbit_host_name) -p 15672:15672 --name rabbit-1 rabbitmq:3-management

rabbit-logs:
	docker logs rabbit-1

rabbit-bash:
	docker exec -it rabbit-1 bash
# Then write this in the bash session 
# $ rabbitmq-plugins enable rabbitmq_management


rabbit-build-producer:
	docker build -t $(producer_img):1.0 -f producer.Dockerfile .

rabbit-producer:
	docker run -it --rm --net rabbits -e RABBIT_HOST=$(rabbit_host_name) -e RABBIT_PORT=5672 -e RABBIT_USERNAME=guest -e RABBIT_PASSWORD=guest --name producer_api -p 9876:9876 $(producer_img):1.0 

rabbit-build-consumer:
	docker build -t $(consumer_img):1.0 -f consumer.Dockerfile .

rabbit-consumer:
	docker run -it --rm --net rabbits -e RABBIT_HOST=$(rabbit_host_name) -e RABBIT_PORT=5672 -e RABBIT_USERNAME=guest -e RABBIT_PASSWORD=guest --name consumer_service $(consumer_img):1.0 

rabbit-instance-cookie:
	docker exec -it $(rabbit_host_name) cat /var/lib/rabbitmq/.erlang.cookie

## ======================================================================================
# Working with clustering [manuallllly] 
rabbit-cluster-node1:
	docker run -d --rm --net rabbits --hostname rabbit-1 --name rabbit-node-1 -p 8081:15672 rabbitmq:3-management

rabbit-node1-cookie:
	docker exec -it rabbit-1 cat /var/lib/rabbitmq/.erlang.cookie

#! Once i have the cookie of the node[1], i know will stop all the instances and run them again with the same cookie 

rabbit-cluster-node1-with-cookie:
	docker run -d --rm --net rabbits --hostname rabbit-1 --name rabbit-node-1 -p 8081:15672 -e RABBITMQ_ERLANG_COOKIE=TEHEYMDVXUBUZKFYWMAO rabbitmq:3-management

rabbit-cluster-node2-with-cookie:
	docker run -d --rm --net rabbits --hostname rabbit-2 --name rabbit-node-2 -p 8082:15672 -e RABBITMQ_ERLANG_COOKIE=TEHEYMDVXUBUZKFYWMAO rabbitmq:3-management

rabbit-cluster-node3-with-cookie:
	docker run -d --rm --net rabbits --hostname rabbit-3 --name rabbit-node-3 -p 8083:15672 -e RABBITMQ_ERLANG_COOKIE=TEHEYMDVXUBUZKFYWMAO rabbitmq:3-management

# to know the cluster of each node [just for testing] 
rabbit-cluster-node1-info:
	docker exec -it rabbit-node-1 rabbitmqctl cluster_status

rabbit-cluster-node2-info:
	docker exec -it rabbit-node-2 rabbitmqctl cluster_status

rabbit-cluster-node3-info:
	docker exec -it rabbit-node-3 rabbitmqctl cluster_status

rabbit-node2-join-node1-cluster:
# make node 2 joining node 1 cluster
	docker exec -it rabbit-node-2 rabbitmqctl stop_app
	docker exec -it rabbit-node-2 rabbitmqctl reset
	docker exec -it rabbit-node-2 rabbitmqctl join_cluster rabbit@rabbit-1
	docker exec -it rabbit-node-2 rabbitmqctl start_app
	docker exec -it rabbit-node-2 rabbitmqctl cluster_status

rabbit-node3-join-node1-cluster:
# make node 3 joining node 1 cluster
	docker exec -it rabbit-node-3 rabbitmqctl stop_app
	docker exec -it rabbit-node-3 rabbitmqctl reset
	docker exec -it rabbit-node-3 rabbitmqctl join_cluster rabbit@rabbit-1
	docker exec -it rabbit-node-3 rabbitmqctl start_app
	docker exec -it rabbit-node-3 rabbitmqctl cluster_status