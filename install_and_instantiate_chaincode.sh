#用peer0.org1创建一个名为mychannel的channel
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.public.com/users/Admin@org1.public.com/msp" cli-org1 peer channel create -o orderer.public.com:7050 -c mychannel -t 30000ms --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/public.com/orderers/orderer.public.com/msp/tlscacerts/tlsca.public.com-cert.pem -f ./channel.tx

#peer0.org1加入频道
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.public.com/users/Admin@org1.public.com/msp" cli-org1 peer channel join -b mychannel.block

#peer0.org2加入频道
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.public.com/users/Admin@org2.public.com/msp" cli-org2 peer channel join -b mychannel.block

#peer0.org3加入频道
docker exec -e "CORE_PEER_LOCALMSPID=Org3MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.public.com/users/Admin@org3.public.com/msp" cli-org3 peer channel join -b mychannel.block

#peer0.org1安装链码
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.public.com/users/Admin@org1.public.com/msp" cli-org1 peer chaincode install -n mycc -v 1.0 -p github.com/chaincode/gochaincode

#peer0.org2安装链码
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.public.com/users/Admin@org2.public.com/msp" cli-org2 peer chaincode install -n mycc -v 1.0 -p github.com/chaincode/gochaincode

#peer0.org3安装链码
docker exec -e "CORE_PEER_LOCALMSPID=Org3MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.public.com/users/Admin@org3.public.com/msp" cli-org3 peer chaincode install -n mycc -v 1.0 -p github.com/chaincode/gochaincode

#peer0.org1实例化链码
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.public.com/users/Admin@org1.public.com/msp" cli-org1 peer chaincode instantiate -o orderer.public.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/public.com/orderers/orderer.public.com/msp/tlscacerts/tlsca.public.com-cert.pem -C mychannel -n mycc -v 1.0 -c '{"Args":[""]}' -P "OR ('Org1MSP.member')"


#注意一切需要与orderer通信的操作都需要加上
# --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/public.com/orderers/orderer.public.com/msp/tlscacerts/tlsca.public.com-cert.pem