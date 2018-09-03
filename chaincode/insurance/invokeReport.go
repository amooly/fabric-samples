package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// 报案
// 第一个参数：保单号
// 第二个参数：报案单信息
// 第三个参数：报案人信息
// 第四个参数：受益人
// 第五个参数：附件(可选)
func Reports(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 4 {
		return shim.Error("Incorrect size of parameters")
	}

	policyNo := args[0]

	// 保存报案记录
	report := &Report{DateType: DataTypeReport, PolicyNo: policyNo}
	if err := json.Unmarshal([]byte(args[1]), report); err != nil {
		return shim.Error(fmt.Sprintf("Failed to assemble report:%s,err:%s", args[1], err))
	}
	if err := PutReport(stub, report); err != nil {
		return shim.Error(err.Error())
	}

	// 报案人
	reporter := &Person{DataType: DataTypeReporter, RelatedNo: report.ReportNo, RelatedDataType: DataTypeReport}
	if err := json.Unmarshal([]byte(args[2]), reporter); err != nil {
		return shim.Error("Failed to assemble reporter:" + err.Error())
	}
	if err := PutReporter(stub, reporter); err != nil {
		return shim.Error(err.Error())
	}

	// 受益人
	var beneficiary = &Person{DataType: DataTypeBeneficiary, RelatedNo: report.ReportNo, RelatedDataType: DataTypeReport}
	if err := json.Unmarshal([]byte(args[3]), beneficiary); err != nil {
		return shim.Error("Failed to assemble beneficiary:" + err.Error())
	}
	if err := PutBeneficiary(stub, beneficiary); err != nil {
		return shim.Error(err.Error())
	}

	// 报案凭证
	if len(args) == 5 {
		evidenceStr := args[4]
		var evidences []Evidence
		if err := json.Unmarshal([]byte(evidenceStr), &evidences); err != nil {
			return shim.Error("Failed to assemble evidence:" + err.Error())
		}
		for _, evidence := range evidences {
			evidence.DataType = DataTypeEvidence
			evidence.RelatedNo = report.ReportNo
			evidence.RelatedDataType = DataTypeReport

			if err := PutReportEvidence(stub, evidence); err != nil {
				return shim.Error(err.Error())
			}
		}
	}
	return shim.Success([]byte(report.ReportNo))
}

func QueryAllReports(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect size of parameters")
	}

	policyNo := args[0]
	reports, err := GetAllReports(stub, policyNo)
	if err != nil {
		return shim.Error(err.Error())
	}

	if reportsByte, err := json.Marshal(reports); err != nil {
		return shim.Error(err.Error())
	} else {
		return shim.Success(reportsByte)
	}
}

func QueryReporter(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect size of parameters")
	}

	reportNo := args[0]
	reporter, err := GetReporter(stub, reportNo)
	if err != nil {
		return shim.Error(err.Error())
	}
	if reporter == nil {
		return shim.Error("reporter find not found")
	}

	if reportByte, err := json.Marshal(reporter); err != nil {
		return shim.Error(err.Error())
	} else {
		return shim.Success(reportByte)
	}
}

func QueryBeneficiary(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect size of parameters")
	}

	reportNo := args[0]
	beneficiary, err := GetBeneficiary(stub, reportNo)
	if err != nil {
		return shim.Error(err.Error())
	}
	if beneficiary == nil {
		return shim.Error("beneficiary find not found")
	}

	if beneficiaryByte, err := json.Marshal(beneficiary); err != nil {
		return shim.Error(err.Error())
	} else {
		return shim.Success(beneficiaryByte)
	}
}

func QueryReportEvidences(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect size of parameters")
	}

	reportNo := args[0]
	evidences, err := GetReportEvidences(stub, reportNo)
	if err != nil {
		return shim.Error(err.Error())
	}

	if evidencesByte, err := json.Marshal(evidences); err != nil {
		return shim.Error(err.Error())
	} else {
		return shim.Success(evidencesByte)
	}
}
