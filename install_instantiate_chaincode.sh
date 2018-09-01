#!/bin/bash
docker exec cli peer chaincode install -p chaincodedev/chaincode/lbh_tuna_demo -n lbh_tuna_demo -v 0
docker exec cli peer chaincode instantiate -n lbh_tuna_demo -v 0  -c '{"Args":[]}' -C myc
