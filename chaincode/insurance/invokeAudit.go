package main

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// 审核
// 第一个参数：保单号
// 第二个参数：报案单号
// 第三个参数：审核记录
// 第四个参数：审核凭证列表（可选）
func Audits(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 3 {
		return shim.Error("Incorrect size of parameters")
	}

	// 保存报案记录
	policyNo := args[0]
	reportNo := args[1]
	var audit = &Audit{DateType: DataTypeAudit}
	if err := json.Unmarshal([]byte(args[2]), audit); err != nil {
		return shim.Error("Failed to assemble report:" + err.Error())
	}

	if err := PutAudit(stub, audit); err != nil {
		return shim.Error(err.Error())
	}

	if err := saveAuditResult(stub, policyNo, reportNo, audit); err != nil {
		return shim.Error(err.Error())
	}

	if len(args) == 4 {
		evidenceStr := args[3]
		var evidences []Evidence
		if err := json.Unmarshal([]byte(evidenceStr), &evidences); err != nil {
			return shim.Error("Failed to assemble evidence:" + err.Error())
		}
		for _, evidence := range evidences {
			evidence.DataType = DataTypeEvidence
			evidence.RelatedNo = audit.AuditNo
			evidence.RelatedDataType = DataTypeAudit

			if err := PutAuditEvidence(stub, evidence); err != nil {
				return shim.Error(err.Error())
			}
		}
	}

	return shim.Success([]byte(reportNo))
}

func saveAuditResult(stub shim.ChaincodeStubInterface, policyNo string, reportNo string, audit *Audit) error {
	if audit.Result == ResultInit {
		return nil
	}

	report, err := GetReport(stub, policyNo, reportNo)
	if err != nil {
		return nil
	} else if report == nil {
		return errors.New("report find not found")
	}

	if audit.Result == ResultReject {
		report.Status = StatusReject
		if err := PutReport(stub, report); err != nil {
			return err
		}
		return nil
	} else if audit.Result == ResultPass {
		// 更新报案单状态
		report.Status = StatusPass
		report.GmtResolve = audit.GmtAudit
		report.ResolveFee = audit.ClaimFee
		if err := PutReport(stub, report); err != nil {
			return err
		}

		// 更新保单状态
		policy, err := GetPolicy(stub, policyNo)
		if err != nil {
			return err
		} else if policy == nil {
			return errors.New("policy find not found")
		}

		policy.ClaimedFee = policy.ClaimedFee + audit.ClaimFee
		if policy.SumInsured < policy.ClaimedFee {
			return errors.New("Claim fee is not allow to surpass sumInsured fee ")
		}
		if err := PutPolicy(stub, policy); err != nil {
			return err
		}

		// 创建凭证使用记录
		evidences, err := GetReportEvidences(stub, reportNo)
		if err != nil {
			return err
		}
		for i := 0; i < len(evidences); i++ {
			evidence := evidences[i]
			eventRecord := &EventRecord{
				DataType:   DataTypeEvidenceRecord,
				EvidenceNo: evidence.EvidenceNo,
				SpNo:       policy.SpNo,
				ComId:      policy.ComId,
				PolicyNo:   policyNo,
				ReportNo:   reportNo,
				GmtReport:  report.GmtReport,
				GmtResolve: report.GmtResolve,
				ResolveFee: audit.ClaimFee,
			}
			if err := PutEventRecord(stub, eventRecord); err != nil {
				return err
			}
		}

		return nil
	}
	return errors.New("Unexpected result")
}

func QueryAllAudits(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect size of parameters")
	}

	reportNo := args[0]
	audits, err := GetAllAudits(stub, reportNo)
	if err != nil {
		return shim.Error(err.Error())
	}

	if auditsByte, err := json.Marshal(audits); err != nil {
		return shim.Error(err.Error())
	} else {
		return shim.Success(auditsByte)
	}
}

func QueryAuditEvidences(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect size of parameters")
	}

	auditNo := args[0]
	evidences, err := GetAuditEvidences(stub, auditNo)
	if err != nil {
		return shim.Error(err.Error())
	}

	if evidencesByte, err := json.Marshal(evidences); err != nil {
		return shim.Error(err.Error())
	} else {
		return shim.Success(evidencesByte)
	}
}
