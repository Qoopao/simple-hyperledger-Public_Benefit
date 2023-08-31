#! /bin/bash
if [ $# -eq 0 ];then
	echo "启动容器: up"
	echo "关闭容器: down"
fi

case $1 in
	up)

		echo "启动 orderer 组织节点..."
		docker-compose -f orderer.yaml up -d

		echo "启动 org1 组织节点..."
		docker-compose -f org1.yaml up -d

		echo "启动 org2 组织节点..."
		docker-compose -f org2.yaml up -d

		echo "启动 org3 组织节点..."
		docker-compose -f org3.yaml up -d

		echo "==========查看进程=========="
		docker-compose -f orderer.yaml ps
		docker-compose -f org1.yaml ps
		docker-compose -f org2.yaml ps
		docker-compose -f org3.yaml ps
		;;
	down)
		echo "========= 关闭 =========="
		docker-compose -f orderer.yaml down -v
		docker-compose -f org1.yaml down -v
		docker-compose -f org2.yaml down -v
		docker-compose -f org3.yaml down -v
		docker rm -f `docker ps -aq`
		docker rmi -f $(docker images | grep "dev-" | awk '{print $ 3}')
		;;
	*)
		echo "args error, do nothing..."
esac

