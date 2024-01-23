package datamodels

type Result struct {
	File     string
	Line     int
	Message  string
	Severity string // "error" or "warning"
}
