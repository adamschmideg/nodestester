#!/bin/sh

server=$(docker ps -q -f name=light_server)
cmd='admin.nodeInfo.enode'
datadir=/root/.ethereum/goerli
enode=$(docker exec -it $server geth --datadir $datadir attach --exec $cmd)

les_cli=$(docker ps -q -f name=les_cli)
cmd="\'admin.addPeer(\"$enode\")\'"
echo $cmd
#x=$(docker exec -it $server geth --datadir $datadir attach --exec $cmd)
echo $x