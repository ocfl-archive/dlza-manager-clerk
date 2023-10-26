package models

type LoggingConfig struct {
	Namespace       string
	StashHost       string
	StashPortNb     int
	TraceLevel      string
	StashTraceLevel string
	Filename        string
}
