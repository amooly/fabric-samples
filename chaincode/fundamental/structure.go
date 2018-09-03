package main

import "time"

const (
	standardFormat = "2006-01-02 15:04:05"
)

// 自定义时间类型，用于进行反序列化
type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+standardFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(standardFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, standardFormat)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(standardFormat)
}

// 事件
type Event struct {
	EventNo     string
	DataType    string
	GmtHappened Time
	Desc        string
	Responsible string
	Damaged     string
	Amount      int
}
