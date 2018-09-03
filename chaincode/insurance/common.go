package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//==================================通用==================================
func getData(stub shim.ChaincodeStubInterface, v interface{}, key string) error {
	dataByte, err := stub.GetState(key)
	if err != nil {
		return err
	}
	if len(dataByte) == 0 {
		return nil
	}
	if err := json.Unmarshal(dataByte, v); err != nil {
		return err
	}
	return nil
}

// 使用组合建获取满足的单个数据
// 适用于组合键中仅保存一个数据的情况，例如报案人
func getSingleDataByPartialCompositeKey(stub shim.ChaincodeStubInterface, v interface{}, objectType string, keys []string) error {
	iterator, err := stub.GetStateByPartialCompositeKey(objectType, keys)
	if err != nil {
		return err
	}
	defer iterator.Close()

	if iterator.HasNext() {
		responseRange, err := iterator.Next()
		if err != nil {
			return err
		}

		if err := json.Unmarshal(responseRange.Value, v); err != nil {
			return err
		}
		return nil
	}
	return nil
}

//==================================单据查询==================================
func GetPolicy(stub shim.ChaincodeStubInterface, policyNo string) (*Policy, error) {
	policy := &Policy{}

	return policy, getData(stub, policy, policyNo)
}

func GetReport(stub shim.ChaincodeStubInterface, policyNo string, reportNo string) (report *Report, err error) {
	reportKey, err := stub.CreateCompositeKey(DataTypePolicy+"~"+DataTypeReport, []string{policyNo, reportNo})

	report = &Report{}
	return report, getData(stub, report, reportKey)
}

func GetAllReports(stub shim.ChaincodeStubInterface, policyNo string) (report []Report, err error) {
	iterator, err := stub.GetStateByPartialCompositeKey(DataTypePolicy+"~"+DataTypeReport, []string{policyNo})
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	reports := []Report{}
	for iterator.HasNext() {
		responseRange, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		var report Report
		if err := json.Unmarshal(responseRange.Value, &report); err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}

func GetAllAudits(stub shim.ChaincodeStubInterface, reportNo string) (audits []Audit, err error) {
	iterator, err := stub.GetStateByPartialCompositeKey(DataTypeReport+"~"+DataTypeAudit, []string{reportNo})
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	audits = []Audit{}
	for iterator.HasNext() {
		responseRange, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		var audit Audit
		if err := json.Unmarshal(responseRange.Value, &audit); err != nil {
			return nil, err
		}
		audits = append(audits, audit)
	}
	return audits, nil
}

//==================================干系人查询==================================
func GetHolder(stub shim.ChaincodeStubInterface, policyNo string) (holder *Person, err error) {
	holder = &Person{}
	return holder, getSingleDataByPartialCompositeKey(stub, holder, DataTypePolicy+"~"+DataTypeHolder, []string{policyNo})
}

func GetInsurers(stub shim.ChaincodeStubInterface, policyNo string) (insurers []Person, err error) {
	iterator, err := stub.GetStateByPartialCompositeKey(DataTypePolicy+"~"+DataTypeInsurer, []string{policyNo})
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	insurers = []Person{}
	for iterator.HasNext() {
		responseRange, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		insurer := &Person{}
		if err := json.Unmarshal(responseRange.Value, insurer); err != nil {
			return nil, err
		}
		insurers = append(insurers, *insurer)
	}
	return insurers, nil
}

func GetReporter(stub shim.ChaincodeStubInterface, reportNo string) (reporter *Person, err error) {
	iterator, err := stub.GetStateByPartialCompositeKey(DataTypeReport+"~"+DataTypeReporter, []string{reportNo})
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	reporter = &Person{}
	for iterator.HasNext() {
		responseRange, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(responseRange.Value, reporter); err != nil {
			return nil, err
		}
		return reporter, nil
	}
	return nil, nil
}
func GetBeneficiary(stub shim.ChaincodeStubInterface, reportNo string) (beneficiary *Person, err error) {
	iterator, err := stub.GetStateByPartialCompositeKey(DataTypeReport+"~"+DataTypeBeneficiary, []string{reportNo})
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	beneficiary = &Person{}
	for iterator.HasNext() {
		responseRange, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(responseRange.Value, &beneficiary); err != nil {
			return nil, err
		}
		return beneficiary, nil
	}
	return nil, nil
}

//==================================凭证查询==================================
func GetReportEvidences(stub shim.ChaincodeStubInterface, reportNo string) (evidences []Evidence, err error) {
	iterator, err := stub.GetStateByPartialCompositeKey(DataTypeReport+"~"+DataTypeEvidence, []string{reportNo})
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	evidences = []Evidence{}
	for iterator.HasNext() {
		responseRange, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		var evidence Evidence
		if err := json.Unmarshal(responseRange.Value, &evidence); err != nil {
			return nil, err
		}
		evidences = append(evidences, evidence)
	}
	return evidences, nil
}

func GetAuditEvidences(stub shim.ChaincodeStubInterface, auditNo string) (evidences []Evidence, err error) {
	iterator, err := stub.GetStateByPartialCompositeKey(DataTypeAudit+"~"+DataTypeEvidence, []string{auditNo})
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	evidences = []Evidence{}
	for iterator.HasNext() {
		responseRange, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		var evidence Evidence
		if err := json.Unmarshal(responseRange.Value, &evidence); err != nil {
			return nil, err
		}
		evidences = append(evidences, evidence)
	}
	return evidences, nil
}

//==================================单据保存==================================
func PutPolicy(stub shim.ChaincodeStubInterface, policy *Policy) error {
	if policyByte, err := json.Marshal(policy); err != nil {
		return err
	} else {
		return stub.PutState(policy.PolicyNo, policyByte)
	}
}

func PutReport(stub shim.ChaincodeStubInterface, report *Report) error {
	if reportKey, err := stub.CreateCompositeKey(DataTypePolicy+"~"+DataTypeReport, []string{report.PolicyNo, report.ReportNo}); err != nil {
		return err
	} else if reportByte, err := json.Marshal(report); err != nil {
		return err
	} else {
		return stub.PutState(reportKey, reportByte)
	}
}

func PutAudit(stub shim.ChaincodeStubInterface, audit *Audit) error {
	if auditKey, err := stub.CreateCompositeKey(DataTypeReport+"~"+DataTypeAudit, []string{audit.ReportNo, audit.AuditNo}); err != nil {
		return err
	} else if auditByte, err := json.Marshal(audit); err != nil {
		return err
	} else {
		return stub.PutState(auditKey, auditByte)
	}
}

//==================================干系人保存==================================

func PutHolder(stub shim.ChaincodeStubInterface, holder *Person) error {
	if holderKey, err := stub.CreateCompositeKey(DataTypePolicy+"~"+DataTypeHolder, []string{holder.RelatedNo, holder.IdCardNo}); err != nil {
		return err
	} else if holderByte, err := json.Marshal(holder); err != nil {
		return err
	} else {
		return stub.PutState(holderKey, holderByte)
	}
}

func PutInsurer(stub shim.ChaincodeStubInterface, insurers *Person) error {
	if insurerKey, err := stub.CreateCompositeKey(DataTypePolicy+"~"+DataTypeInsurer, []string{insurers.RelatedNo, insurers.IdCardNo}); err != nil {
		return err
	} else if insurerByte, err := json.Marshal(insurers); err != nil {
		return err
	} else {
		return stub.PutState(insurerKey, insurerByte)
	}
}

func PutReporter(stub shim.ChaincodeStubInterface, reporter *Person) error {
	if reporterKey, err := stub.CreateCompositeKey(DataTypeReport+"~"+DataTypeReporter, []string{reporter.RelatedNo, reporter.IdCardNo}); err != nil {
		return err
	} else if reporterByte, err := json.Marshal(reporter); err != nil {
		return err
	} else {
		return stub.PutState(reporterKey, reporterByte)
	}
}

func PutBeneficiary(stub shim.ChaincodeStubInterface, beneficiary *Person) error {
	if beneficiaryKey, err := stub.CreateCompositeKey(DataTypeReport+"~"+DataTypeBeneficiary, []string{beneficiary.RelatedNo, beneficiary.IdCardNo}); err != nil {
		return err
	} else if beneficiaryByte, err := json.Marshal(beneficiary); err != nil {
		return err
	} else {
		return stub.PutState(beneficiaryKey, beneficiaryByte)
	}
}

//==================================凭证保存==================================
func PutReportEvidence(stub shim.ChaincodeStubInterface, evidence Evidence) error {
	if evidenceKey, err := stub.CreateCompositeKey(DataTypeReport+"~"+DataTypeEvidence, []string{evidence.RelatedNo, evidence.EvidenceNo}); err != nil {
		return err
	} else if evidenceByte, err := json.Marshal(evidence); err != nil {
		return err
	} else {
		return stub.PutState(evidenceKey, evidenceByte)
	}
}

func PutAuditEvidence(stub shim.ChaincodeStubInterface, evidence Evidence) error {
	if evidenceKey, err := stub.CreateCompositeKey(DataTypeAudit+"~"+DataTypeEvidence, []string{evidence.RelatedNo, evidence.EvidenceNo}); err != nil {
		return err
	} else if evidenceByte, err := json.Marshal(evidence); err != nil {
		return err
	} else {
		return stub.PutState(evidenceKey, evidenceByte)
	}
}

func PutEventRecord(stub shim.ChaincodeStubInterface, record *EventRecord) error {
	if recordKey, err := stub.CreateCompositeKey(DataTypeEvidenceRecord, []string{record.EvidenceNo, record.SpNo, record.ComId}); err != nil {
		return err
	} else if recordByte, err := json.Marshal(record); err != nil {
		return err
	} else {
		return stub.PutState(recordKey, recordByte)
	}
}
