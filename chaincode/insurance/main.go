package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("main")
var ccFunctions = map[string]func(shim.ChaincodeStubInterface, []string) peer.Response{
	// 保单相关
	"insure":        Insures,
	"queryPolicy":   QueryPolicy,
	"queryHolder":   QueryHolder,
	"queryInsurers": QueryInsurers,

	// 报案
	"report":               Reports,
	"queryAllReports":      QueryAllReports,
	"queryReporter":        QueryReporter,
	"queryBeneficiary":     QueryBeneficiary,
	"queryReportEvidences": QueryReportEvidences,

	// 审核
	"audit":               Audits,
	"queryAllAudits":      QueryAllAudits,
	"queryAuditEvidences": QueryAuditEvidences,
}

type Insurance struct {
}

func (t *Insurance) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}
func (t *Insurance) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
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

	err := shim.Start(new(Insurance))
	if err != nil {
		logger.Error("Error starting Simple chaincode: ", err)
	}
}
