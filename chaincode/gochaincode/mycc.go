package main

/*----------------------------------------------------*/
/*--------------------引入所需依赖---------------------*/
/*----------------------------------------------------*/
import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	//这里是将后面这个包设置一个protocolbuffers的别名
	protocolbuffers "github.com/hyperledger/fabric/protos/peer"
)

//链码的接口,为空
type Chaincode struct {
}

//公益慈善项目的基本信息
type PB_Project struct {
	ProjectID          string `json:ProjectID`          //项目ID
	ProjectDESC        string `json:ProjectDESC`        //项目介绍
	ProjectForWho      string `json:ProjectForWho`      //项目受益人
	ProjectFundNeed    int    `json:ProjectFundNeed`    //项目所需资金
	ProjectFundCurrent int    `json:ProjectFundCurrent` //项目目前筹集到的资金
	Status             string `json:Status`             //项目的状态
	NumOfDonate        int    `json:NumOfDonate`        //项目当前的捐款次数
	HelpID             string `json:HelpID`             //绑定的求助信息
}

//求助信息
type Help_Info struct {
	HelpID      string `json:HelpID`      //求助信息的ID
	UserID      string `json:UserID`      //求助人的ID
	FundsNeeded int    `json:FundNeeded`  //所需金额
	Desc        string `json:Description` //描述
	Status      string `json:Status`      //该求助的状态
	ProjectID   string `json:ProjectID`   //绑定的项目ID
}

//捐款收据
type Receipt_Info struct {
	ReceiptID   string `json:Receipt`     //收据ID
	UserID      string `json:UserID`      //捐款人的ID
	FundDonated int    `json:FundDonated` //捐款金额
	Date        string `json:Date`        //日期
}

//举报信息
type Report_Info struct {
	ReportID   string `json:ReportID`   //举报信息的ID
	ReportDESC string `json:ReportDESC` //举报信息的描述
	Status     string `json:Status`     //举报状态
}

type User_Info struct {
	UserID          string `json:UserID`          //用户ID
	UserName        string `json:Username`        //用户名
	UserPhoneNumber string `json:UserPhoneNumber` //用户手机号码
	UserEmail       string `json:UserEmail`       //用户邮箱
	UserGender      string `json:UserGender`      //用户性别
	UserOrg         string `json:UserOrg`         //用户所属组织
	UserPem         string `json:UserPem`         //用户的证书
	UserRole        string `json:UserRole`        //用户的角色
	UserBalance     uint   `json:UserBalance`     //用户余额
}

func (receiver *Chaincode) Init(stub shim.ChaincodeStubInterface) protocolbuffers.Response {
	return shim.Success(nil)
}

func (receiver *Chaincode) Invoke(stub shim.ChaincodeStubInterface) protocolbuffers.Response {

	functionInvoked, args := stub.GetFunctionAndParameters()
	if functionInvoked == "addPBProject" {
		return receiver.addPBProject(stub, args)
	} else if functionInvoked == "queryPBProject" {
		return receiver.queryPBProject(stub, args)
	} else if functionInvoked == "addHelp" {
		return receiver.addHelp(stub, args)
	} else if functionInvoked == "queryHelp" {
		return receiver.queryHelp(stub, args)
	} else if functionInvoked == "updatePBProject" {
		return receiver.updatePBProject(stub, args)
	} else if functionInvoked == "updateHelp" {
		return receiver.updateHelp(stub, args)
	} else if functionInvoked == "addReceipt" {
		return receiver.addReceipt(stub, args)
	} else if functionInvoked == "queryReceipt" {
		return receiver.queryReceipt(stub, args)
	} else if functionInvoked == "addReport" {
		return receiver.addReport(stub, args)
	} else if functionInvoked == "queryReport" {
		return receiver.queryReport(stub, args)
	} else if functionInvoked == "updateReport" {
		return receiver.updateReport(stub, args)
	} else if functionInvoked == "addUser" {
		return receiver.addUser(stub, args)
	} else if functionInvoked == "queryUser" {
		return receiver.queryUser(stub, args)
	}

	//不能出现else否则报错
	return shim.Error("function called doesnt exist")

}

//增加项目的函数
func (receiver *Chaincode) addPBProject(stub shim.ChaincodeStubInterface, args []string) protocolbuffers.Response {

	if len(args) != 4 {
		return shim.Error("incorrect numbers of args")
	}

	var msg string
	var flag bool
	flag, msg = receiver.isexist(stub, args[0])

	if flag == true {
		return shim.Error("add failed, already existed")
	}
	fmt.Println(msg)

	var pbp PB_Project
	pbp.ProjectID = args[0]
	pbp.ProjectDESC = args[1]
	pbp.ProjectForWho = args[2]
	pbp.HelpID = ""
	tmp, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("project fund needed :convert to number wrong")
	}
	pbp.ProjectFundNeed = tmp
	pbp.ProjectFundCurrent = 0
	pbp.Status = "newly created"
	pbp.NumOfDonate = 0

	pbpJson, err1 := json.Marshal(pbp)
	if err1 != nil {
		return shim.Error(err1.Error())
	}

	err = stub.PutState(pbp.ProjectID, pbpJson)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)

}

//查询项目的函数
func (receiver *Chaincode) queryPBProject(stub shim.ChaincodeStubInterface, args []string) protocolbuffers.Response {

	if len(args) != 1 {
		return shim.Error("incorrect number of args")
	}

	var msg string
	var flag bool
	flag, msg = receiver.isexist(stub, args[0])

	if flag == false && msg == "not existed" {
		return shim.Error("query failed, not existed")
	}
	fmt.Println(msg)

	projectID := args[0]
	resultsIterator, err := stub.GetHistoryForKey(projectID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var pbp PB_Project

	for resultsIterator.HasNext() {
		response, err1 := resultsIterator.Next()
		if err1 != nil {
			return shim.Error(err1.Error())
		}
		err := json.Unmarshal(response.Value, &pbp)
		if err != nil {
			return shim.Error("Unmarshal failed")
		}
	}

	jsonsAsBytes, err2 := json.Marshal(pbp)
	if err2 != nil {
		return shim.Error(err2.Error())
	}

	return shim.Success(jsonsAsBytes)
}

//修改项目状态的函数
func (receiver *Chaincode) updatePBProject(stub shim.ChaincodeStubInterface, args []string) protocolbuffers.Response {
	//第一个参数为项目ID，第二个参数为状态，第三个参数为绑定的求助信息ID
	if len(args) != 3 {
		return shim.Error("incorrect number of args")
	}

	var msg string
	var flag bool
	flag, msg = receiver.isexist(stub, args[0])

	if flag == false && msg == "not existed" {
		return shim.Error("query failed, not existed")
	}
	fmt.Println(msg)

	projectID := args[0]
	resultsIterator, err := stub.GetHistoryForKey(projectID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var pbp PB_Project

	for resultsIterator.HasNext() {
		response, err1 := resultsIterator.Next()
		if err1 != nil {
			return shim.Error(err1.Error())
		}
		err := json.Unmarshal(response.Value, &pbp)
		if err != nil {
			return shim.Error("Unmarshal failed")
		}
	}

	pbp.Status = args[1]
	pbp.HelpID = args[2]

	pbpJson, err2 := json.Marshal(pbp)
	if err2 != nil {
		return shim.Error(err2.Error())
	}

	err = stub.PutState(pbp.ProjectID, pbpJson)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//增加求助信息的函数
func (receiver *Chaincode) addHelp(stub shim.ChaincodeStubInterface, args []string) protocolbuffers.Response {

	if len(args) != 4 {
		return shim.Error("incorrect number of args")
	}

	var msg string
	var flag bool
	flag, msg = receiver.isexist(stub, args[0])

	if flag == true {
		return shim.Error("add failed, already existed")
	}
	fmt.Println(msg)

	var hi Help_Info

	hi.HelpID = args[0]
	hi.UserID = args[1]
	hi.Desc = args[2]
	tmp, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("help info fund needed: convert to number wrong")
	}
	hi.FundsNeeded = tmp
	hi.Status = "waiting for auditing"
	hi.ProjectID = ""

	hiJson, err1 := json.Marshal(hi)
	if err1 != nil {
		return shim.Error(err1.Error())
	}

	err = stub.PutState(hi.HelpID, hiJson)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//查询求助信息的函数
func (receiver *Chaincode) queryHelp(stub shim.ChaincodeStubInterface, args []string) protocolbuffers.Response {
	if len(args) != 1 {
		return shim.Error("incorrect number of args")
	}

	var msg string
	var flag bool
	flag, msg = receiver.isexist(stub, args[0])

	if flag == false && msg == "not existed" {
		return shim.Error("query failed, not existed")
	}
	fmt.Println(msg)

	helpID := args[0]
	resultsIterator, err := stub.GetHistoryForKey(helpID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var hi Help_Info

	for resultsIterator.HasNext() {
		response, err1 := resultsIterator.Next()
		if err1 != nil {
			return shim.Error(err1.Error())
		}
		err := json.Unmarshal(response.Value, &hi)
		if err != nil {
			return shim.Error("Unmarshal failed")
		}
	}

	jsonsAsBytes, err2 := json.Marshal(hi)
	if err2 != nil {
		return shim.Error(err2.Error())
	}

	return shim.Success(jsonsAsBytes)
}

//修改求助信息状态的函数
func (receiver *Chaincode) updateHelp(stub shim.ChaincodeStubInterface, args []string) protocolbuffers.Response {
	//第一个参数为求助信息ID，第二个参数为状态,第三个参数为绑定的项目ID
	if len(args) != 3 {
		return shim.Error("incorrect number of args")
	}

	var msg string
	var flag bool
	flag, msg = receiver.isexist(stub, args[0])

	if flag == false && msg == "not existed" {
		return shim.Error("query failed, not existed")
	}
	fmt.Println(msg)

	helpID := args[0]
	resultsIterator, err := stub.GetHistoryForKey(helpID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var hi Help_Info

	for resultsIterator.HasNext() {
		response, err1 := resultsIterator.Next()
		if err1 != nil {
			return shim.Error(err1.Error())
		}
		err := json.Unmarshal(response.Value, &hi)
		if err != nil {
			return shim.Error("Unmarshal failed")
		}
	}

	hi.Status = args[1]
	hi.ProjectID = args[2]

	hiJson, err2 := json.Marshal(hi)
	if err2 != nil {
		return shim.Error(err2.Error())
	}

	err = stub.PutState(hi.HelpID, hiJson)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//增加收据
func (receiver *Chaincode) addReceipt(stub shim.ChaincodeStubInterface, args []string) protocolbuffers.Response {
	if len(args) != 4 {
		return shim.Error("incorrect number of args")
	}

	//第一个参数为收据ID，第二个参数为用户ID，第三个参数为捐款金额，第四个参数为日期

	var msg string
	var flag bool
	flag, msg = receiver.isexist(stub, args[0])

	if flag == true {
		return shim.Error("add failed, already existed")
	}
	fmt.Println(msg)

	var ri Receipt_Info

	ri.ReceiptID = args[0]
	ri.UserID = args[1]
	ri.Date = args[3]
	tmp, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("help info fund needed: convert to number wrong")
	}
	ri.FundDonated = tmp

	riJson, err1 := json.Marshal(ri)
	if err1 != nil {
		return shim.Error(err1.Error())
	}

	err = stub.PutState(ri.ReceiptID, riJson)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//查询收据
func (receiver *Chaincode) queryReceipt(stub shim.ChaincodeStubInterface, args []string) protocolbuffers.Response {
	if len(args) != 1 {
		return shim.Error("incorrect number of args")
	}

	var msg string
	var flag bool
	flag, msg = receiver.isexist(stub, args[0])

	if flag == false && msg == "not existed" {
		return shim.Error("query failed, not existed")
	}
	fmt.Println(msg)

	receiptID := args[0]
	resultsIterator, err := stub.GetHistoryForKey(receiptID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var ri Receipt_Info

	for resultsIterator.HasNext() {
		response, err1 := resultsIterator.Next()
		if err1 != nil {
			return shim.Error(err1.Error())
		}
		err := json.Unmarshal(response.Value, &ri)
		if err != nil {
			return shim.Error("Unmarshal failed")
		}
	}

	jsonsAsBytes, err2 := json.Marshal(ri)
	if err2 != nil {
		return shim.Error(err2.Error())
	}

	return shim.Success(jsonsAsBytes)
}

//增加举报信息
func (receiver *Chaincode) addReport(stub shim.ChaincodeStubInterface, args []string) protocolbuffers.Response {
	if len(args) != 2 {
		return shim.Error("incorrect number of args")
	}

	//第一个参数为举报ID，第二个参数为举报描述

	var msg string
	var flag bool
	flag, msg = receiver.isexist(stub, args[0])

	if flag == true {
		return shim.Error("add failed, already existed")
	}
	fmt.Println(msg)

	var ri Report_Info

	ri.ReportID = args[0]
	ri.ReportDESC = args[1]
	ri.Status = "waiting for auditing"

	riJson, err1 := json.Marshal(ri)
	if err1 != nil {
		return shim.Error(err1.Error())
	}

	err := stub.PutState(ri.ReportID, riJson)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//查询举报信息
func (receiver *Chaincode) queryReport(stub shim.ChaincodeStubInterface, args []string) protocolbuffers.Response {
	if len(args) != 1 {
		return shim.Error("incorrect number of args")
	}

	var msg string
	var flag bool
	flag, msg = receiver.isexist(stub, args[0])

	if flag == false && msg == "not existed" {
		return shim.Error("query failed, not existed")
	}
	fmt.Println(msg)

	reportID := args[0]
	resultsIterator, err := stub.GetHistoryForKey(reportID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var ri Report_Info

	for resultsIterator.HasNext() {
		response, err1 := resultsIterator.Next()
		if err1 != nil {
			return shim.Error(err1.Error())
		}
		err := json.Unmarshal(response.Value, &ri)
		if err != nil {
			return shim.Error("Unmarshal failed")
		}
	}

	jsonsAsBytes, err2 := json.Marshal(ri)
	if err2 != nil {
		return shim.Error(err2.Error())
	}

	return shim.Success(jsonsAsBytes)
}

//修改举报信息状态
func (receiver *Chaincode) updateReport(stub shim.ChaincodeStubInterface, args []string) protocolbuffers.Response {
	//第一个参数为举报ID，第二个参数为状态
	if len(args) != 2 {
		return shim.Error("incorrect number of args")
	}

	var msg string
	var flag bool
	flag, msg = receiver.isexist(stub, args[0])

	if flag == false && msg == "not existed" {
		return shim.Error("query failed, not existed")
	}
	fmt.Println(msg)

	reportID := args[0]
	resultsIterator, err := stub.GetHistoryForKey(reportID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var ri Report_Info

	for resultsIterator.HasNext() {
		response, err1 := resultsIterator.Next()
		if err1 != nil {
			return shim.Error(err1.Error())
		}
		err := json.Unmarshal(response.Value, &ri)
		if err != nil {
			return shim.Error("Unmarshal failed")
		}
	}

	ri.Status = args[1]

	riJson, err2 := json.Marshal(ri)
	if err2 != nil {
		return shim.Error(err2.Error())
	}

	err = stub.PutState(ri.ReportID, riJson)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//增加用户信息
func (receiver *Chaincode) addUser(stub shim.ChaincodeStubInterface, args []string) protocolbuffers.Response {
	if len(args) != 8 {
		shim.Error("incorrect number of args")
	}
	//参数分别为：ID,用户名，手机号，邮箱，性别，组织，证书，角色

	var msg string
	var flag bool
	flag, msg = receiver.isexist(stub, args[0])

	if flag == true {
		return shim.Error("add failed, already existed")
	}
	fmt.Println(msg)

	var user User_Info
	user.UserID = args[0]
	user.UserName = args[1]
	user.UserPhoneNumber = args[2]
	user.UserEmail = args[3]
	user.UserGender = args[4]
	user.UserOrg = args[5]
	user.UserPem = args[6]
	user.UserRole = args[7]
	user.UserBalance = 0

	userJson, err1 := json.Marshal(user)
	if err1 != nil {
		return shim.Error(err1.Error())
	}

	err := stub.PutState(user.UserID, userJson)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

//查询用户信息
func (receiver *Chaincode) queryUser(stub shim.ChaincodeStubInterface, args []string) protocolbuffers.Response {
	if len(args) != 1 {
		return shim.Error("incorrect number of args")
	}

	var msg string
	var flag bool
	flag, msg = receiver.isexist(stub, args[0])

	if flag == false && msg == "not existed" {
		return shim.Error("query failed, not existed")
	}
	fmt.Println(msg)

	userID := args[0]
	resultsIterator, err := stub.GetHistoryForKey(userID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var user User_Info

	for resultsIterator.HasNext() {
		response, err1 := resultsIterator.Next()
		if err1 != nil {
			return shim.Error(err1.Error())
		}
		err := json.Unmarshal(response.Value, &user)
		if err != nil {
			return shim.Error("Unmarshal failed")
		}
	}

	jsonsAsBytes, err2 := json.Marshal(user)
	if err2 != nil {
		return shim.Error(err2.Error())
	}

	return shim.Success(jsonsAsBytes)
}

//通用方法
func (receiver *Chaincode) isexist(stub shim.ChaincodeStubInterface, args string) (bool, string) {

	ID := args
	resultsIterator, err := stub.GetHistoryForKey(ID)
	if err != nil {
		return false, "gethistoryforkey failed"
	}
	defer resultsIterator.Close()

	if resultsIterator.HasNext() {
		return true, "existed"
	}

	return false, "not existed"
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s ", err)
	}
}
