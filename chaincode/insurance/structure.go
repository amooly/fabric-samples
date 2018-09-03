package main

//=============================常量===========================
const (
	DataTypePolicy         = "POLICY"
	DataTypeReport         = "REPORT"
	DataTypeAudit          = "AUDIT"
	DataTypeEvidence       = "EVIDENCE"
	DataTypeEvidenceRecord = "EVIDENCE_RECORD"
	DataTypeHolder         = "HOLDER"
	DataTypeInsurer        = "INSURER"
	DataTypeReporter       = "REPORTER"
	DataTypeBeneficiary    = "BENEFICIARY"
	StatusInit             = 0
	StatusPass             = 1
	StatusReject           = 2
	ResultInit             = 0
	ResultPass             = 1
	ResultReject           = 2
)

//=============================保险数据：私有数据==============
// 干系人：投保人、被保人、受益人
type Person struct {
	DataType        string
	Name            string
	Age             int
	IdCardNo        string
	SpNo            string
	RelatedNo       string
	RelatedDataType string
}

// 凭证（例如发票、交通事故认定责任书等）：报案凭证，审核凭证
// 单据号+凭证号
type Evidence struct {
	DataType        string
	EvidenceNo      string
	Name            string
	Url             string
	RelatedNo       string
	RelatedDataType string
}

// 保单
// 使用保单号为key
type Policy struct {
	DataType       string
	PolicyNo       string
	GmtInsured     Time
	GmtEffectStart Time
	SpNo           string
	ComId          string
	SumInsured     int
	ClaimedFee     int
}

// 报案记录
// 使用保单号+报案单号作为compositeKey
type Report struct {
	DateType   string
	ReportNo   string
	GmtReport  Time
	Status     int
	PolicyNo   string
	ReportFee  int
	Desc       string
	GmtResolve Time
	ResolveFee int
}

// 审核记录
// 使用报案单号+审核单号作为compositeKey
type Audit struct {
	DateType string
	AuditNo  string
	ReportNo string
	GmtAudit Time
	Auditor  string
	Result   int
	ClaimFee int
	Desc     string
}

//=============================共享数据=======================
// 凭证理赔记录
// 用于共享同一个凭证的使用情况
type EventRecord struct {
	DataType   string
	EvidenceNo string
	SpNo       string
	ComId      string
	PolicyNo   string
	ReportNo   string
	GmtReport  Time
	GmtResolve Time
	ResolveFee int
}
