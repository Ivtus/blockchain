package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type DrmChaincode struct{}

type Rights struct {
	Owner string `json:"owner"`
	Hash  string `json:"hash"`
}

func (t *DrmChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Chaincode Initialize")
	return shim.Success(nil)
}

func (t *DrmChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	// 获取调用的方法和参数
	fn, args := stub.GetFunctionAndParameters()

	// 判断分发 业务方法
	if fn == "initRights" {
		// 调用创建弹珠方法
		return t.initRights(stub, args)
	} else if fn == "readRights" {
		// 调用读取弹珠信息的方法
		return t.readRights(stub, args)
	} else if fn == "deleteRights" {
		// 调用删除弹珠信息
		return t.deleteRights(stub, args)
	} else if fn == "transferRights" {
		// 调用交易弹珠的方法
		return t.transferRights(stub, args)
		// } else if fn == "queryRightsByOwner" {
		// 	// 调用查询用户拥有的弹珠信息
		// 	return t.queryRightsByOwner(stub, args)
		// } else if fn == "queryHistoryForRights" {
		// 	// 查询弹珠的历史操作信息
		// 	return t.queryHistoryForRights(stub, args)
	}

	// 如果没有对应的方法 返回错误
	return shim.Error(fn + " 方法不存在!")
}

func (t *DrmChaincode) initRights(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	hash := args[1]
	// 查询账本中是否已经存在该版权的信息
	rightsBytes, err := stub.GetState(hash)

	if err != nil {
		return shim.Error(err.Error())
	}

	if rightsBytes != nil {
		return shim.Error("Rights Exist")
	}

	// 如果不存在 则写入到账本中
	owner := args[0]

	// 组装测结构体
	rights := &Rights{owner, hash}

	// 将rights 转成json字符串 存储到账本
	rightsJsonStr, err := json.Marshal(rights)
	if err != nil {
		return shim.Error(err.Error())
	}

	// PutState json信息写入账本
	err = stub.PutState(hash, rightsJsonStr)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Rights Initialize success")
	// 同时创建组合键用于查询
	indexName := "owner~record"
	indexKey, err := stub.CreateCompositeKey(indexName, []string{owner, string(rightsJsonStr)})
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(indexKey, []byte{0x00})
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println(indexKey)
	return shim.Success(nil)
}

func (t *DrmChaincode) readRights(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 获取参数 参数为hash的rights
	hash := args[0]

	// 根据hash 读取rights的数据
	rightsBytes, err := stub.GetState(hash)

	if err != nil {
		return shim.Error(err.Error())
	}

	if rightsBytes == nil {
		return shim.Error("Rights dosen't exist")
	}

	// 返回信息
	return shim.Success(rightsBytes)

}

func (t *DrmChaincode) deleteRights(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 从参数总获取到rights的hash
	hash := args[0]

	// 判断hash是否存在
	rightsBytes, err := stub.GetState(hash)
	if err != nil {
		return shim.Error(err.Error())
	}

	if rightsBytes == nil {
		return shim.Error("rights don't exist")
	}

	// 删除弹珠
	err = stub.DelState(hash)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("rights have deleted")
	return shim.Success(nil)
}

func (t *DrmChaincode) transferRights(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// 获取到参数
	newOwner := args[0]

	hash := args[1]

	// 检查弹珠是否存在
	rightsBytes, err := stub.GetState(hash)
	if err != nil {
		return shim.Error(err.Error())
	}
	if rightsBytes == nil {
		return shim.Error("Rights don't exist")
	}

	// 将账本中的信息转为 Rights 结构体
	rightsInfo := Rights{}
	err = json.Unmarshal(rightsBytes, &rightsInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	// 修改拥有者
	rightsInfo.Owner = newOwner

	// 转为json数据
	newRightsBytes, err := json.Marshal(rightsInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	// 写入账本
	err = stub.PutState(hash, newRightsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Now" + newOwner + "has the" + hash + "'rights.")
	return shim.Success(nil)

}

func main() {
	err := shim.Start(new(DrmChaincode))
	if err != nil {
		fmt.Printf("Error start DrmChaincode")
	}
}
