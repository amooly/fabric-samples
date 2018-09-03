package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("main")
var ccFunctions = map[string]func(shim.ChaincodeStubInterface, []string) peer.Response{
	// 保单相关
	"saveEvent":           saveEvent,
	"queryEventByEventNo": queryEventByEventNo,
	"queryEventByDamaged": queryEventByDamaged,
}

type Fundamental struct {
}

func (t *Fundamental) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}
func (t *Fundamental) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("start to query")

	funcName, args := stub.GetFunctionAndParameters()

	ccFunc := ccFunctions[funcName]
	if ccFunc == nil {
		logger.Error("Failed to find corresponding function:" + funcName)
		return shim.Error("Invalid invoke funcName.")
	}
	return ccFunc(stub, args)
}

func main() {
	logger.SetLevel(shim.LogInfo)

	err := shim.Start(new(Fundamental))
	if err != nil {
		logger.Error("Error starting Fundamental chaincode: ", err)
	}
}
