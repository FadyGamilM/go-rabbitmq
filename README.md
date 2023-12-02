# go-rabbitmq

# step.1 - run container
- create a new network so all instances can talk to each other
```shell
make rabbit-network
```
- spin rabbitmq docker container 
```shell
make rabbit-instance
```
- if we explored the container 
```shell
➜ go-rabbitmq git:(main) docker container ps
CONTAINER ID   IMAGE                  COMMAND                  CREATED          STATUS          PORTS                                                 NAMES
abbfe94ed74a   rabbitmq               "docker-entrypoint.s…"   18 minutes ago   Up 17 minutes   4369/tcp, 5671-5672/tcp, 15691-15692/tcp, 25672/tcp   rabbit-1
```
notice that :
 * if we run inside a cluster, the port is the 4369/tcp 
 * if we have another application needs to communicate to the queue and consume it, the port is 5671-5672/tcp

# step.2 - enable management ui
- open bash inside the container by running 
```shell
make rabbit-bash
```
- then run this command via the rabbitmq-plugin cli 
```bash
rabbitmq-plugins enable rabbitmq_management
```
- make a http request on the browser to `localhost:15672` <br>
- the port that I am pointing to is the port that we will use for the application that publishs msgs to the queue 
![ports](./app-ports.png) 

# Theory 
- Channel is a vritual connection to a specific queue 