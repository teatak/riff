package common

const (
	Name            = "riff"
	Reset           = "\033[0m"
	Red             = "\033[31;1m"
	Green           = "\033[32;1m"
	Success         = Green + "[  OK  ]" + Reset
	Failed          = Red + "[FAILED]" + Reset
	DefaultHttpPort = 8610
	DefaultDnsPort  = 8620
	DefaultRpcPort  = 8630
)
