#!/bin/bash
docker exec chaincode /bin/bash  -c " cd lbh_tuna_demo ; echo prepare to build ; go build ; "
docker exec chaincode /bin/bash -c " cd lbh_tuna_demo ; CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=lbh_tuna_demo:0 ./lbh_tuna_demo "
