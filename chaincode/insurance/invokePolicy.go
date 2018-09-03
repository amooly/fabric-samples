package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// 投保
// 第一个参数：保单信息
// 第二个参数：投保人
// 第三个参数：被保人
func Insures(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect size of parameters")
	}

	// 添加保单
	policy := &Policy{DataType: DataTypePolicy}
	if err := json.Unmarshal([]byte(args[0]), policy); err != nil {
		return shim.Error("Failed to assemble policy:" + err.Error())
	}
	if err := PutPolicy(stub, policy); err != nil {
		return shim.Error(err.Error())
	}

	// 添加投保人
	holder := &Person{DataType: DataTypeHolder, SpNo: policy.SpNo, RelatedNo: policy.PolicyNo, RelatedDataType: DataTypePolicy}
	if err := json.Unmarshal([]byte(args[1]), holder); err != nil {
		return shim.Error("Failed to assemble holder:" + err.Error())
	}
	if err := PutHolder(stub, holder); err != nil {
		return shim.Error(err.Error())
	}

	// 添加被保人列表
	var insurers []Person
	if err := json.Unmarshal([]byte(args[2]), &insurers); err != nil {
		return shim.Error("Failed to assemble insurer:" + err.Error())
	}
	for _, insurer := range insurers {
		insurer.DataType = DataTypeInsurer
		insurer.SpNo = policy.SpNo
		insurer.RelatedNo = policy.PolicyNo
		insurer.RelatedDataType = DataTypePolicy

		if err := PutInsurer(stub, &insurer); err != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success([]byte(policy.PolicyNo))
}

func QueryPolicy(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect size of parameters")
	}

	policyNo := args[0]

	policy, err := GetPolicy(stub, policyNo)
	if err != nil {
		return shim.Error(err.Error())
	}

	if policyByte, err := json.Marshal(&policy); err != nil {
		return shim.Error(err.Error())
	} else {
		return shim.Success(policyByte)
	}
}

func QueryHolder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect size of parameters")
	}

	policyNo := args[0]

	holder, err := GetHolder(stub, policyNo)
	if err != nil {
		return shim.Error(err.Error())
	}

	if holderByte, err := json.Marshal(&holder); err != nil {
		return shim.Error(err.Error())
	} else {
		return shim.Success(holderByte)
	}
}

func QueryInsurers(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect size of parameters")
	}

	policyNo := args[0]

	insurers, err := GetInsurers(stub, policyNo)
	if err != nil {
		return shim.Error(err.Error())
	}

	if insurersByte, err := json.Marshal(&insurers); err != nil {
		return shim.Error(err.Error())
	} else {
		return shim.Success(insurersByte)
	}
}
