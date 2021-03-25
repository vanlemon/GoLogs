#!/usr/bin/env bash

cd `dirname $0`



#sudo systemctl start docker

##https://docs.docker.com/compose/install/
#sudo curl -L "https://github.com/docker/compose/releases/download/1.28.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
#sudo chmod +x /usr/local/bin/docker-compose
#sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose

sudo docker-compose up -d

sudo docker-compose ps

#docker-compose stop
#docker-compose rm

# http://127.0.0.1:9000/

