# simple-hyperledger-Public_Benefit
Fabric1.4.4毕设，三组织，solo排序，仅区块链部分
二进制工具以及docker镜像需要自行下载

//查询链码实例化情况
docker exec -it cli bash
peer chaincode list -C mychannel --instantiated

//关闭dockers环境
docker stop $(docker ps -aq)
docker rm $(docker ps -aq)

//调用链码，进到docker环境后执行

{
    增加项目：peer chaincode invoke -n mycc -c '{"Args":["addPBProject","ProjectID1","simple","boiy","100"] }' -C mychannel 
    查看项目：peer chaincode invoke -n mycc -c '{"Args":["queryPBProject","ProjectID1"] }' -C mychannel
    修改项目：peer chaincode invoke -n mycc -c '{"Args":["updatePBProject","ProjectID1","Help request binded","HelpID1"] }' -C mychannel

    增加求助信息：peer chaincode invoke -n mycc -c '{"Args":["addHelp","HelpID1","UserID1","helphelphelp","100"]}' -C mychannel
    查看求助信息：peer chaincode invoke -n mycc -c '{"Args":["queryHelp","HelpID1"]}' -C mychannel
    修改求助信息：peer chaincode invoke -n mycc -c '{"Args":["updateHelp","HelpID1","Project binded","ProjectID1"] }' -C mychannel

    增加收据：peer chaincode invoke -n mycc -c '{"Args":["addReceipt","ReceiptID1","UserID1","100","today"] }' -C mychannel
    查看收据：peer chaincode invoke -n mycc -c '{"Args":["queryReceipt","ReceiptID1"] }' -C mychannel

    增加举报信息：peer chaincode invoke -n mycc -c '{"Args":["addReport","ReportID2","I want to Report"]}' -C mychannel
    查看举报信息：peer chaincode invoke -n mycc -c '{"Args":["queryReport","ReportID2"]}' -C mychannel
    修改举报信息：peer chaincode invoke -n mycc -c '{"Args":["updateReport","ReportID2","done"] }' -C mychannel

    增加用户：peer chaincode invoke -n mycc -c '{"Args":["addUser","UserID1","lyh","1768","@qq.com","male","Org1","a.pem","admin"] }' -C mychannel
    查看用户：peer chaincode invoke -n mycc -c '{"Args":["queryUser","UserID1"] }' -C mychannel



1.结构体变量首字母一定要大写，这是golang变量类型决定的
2.每一次重新部署都要改链码名字
