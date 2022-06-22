package internal

type Settings struct {
	Debug   bool
	Port    string
	Timeout int

	MetricsPath string

	DBConnectionInfo DBConnectionInfo
}

type DBConnectionInfo struct {
	User     string
	Password string
	HostName string
	DBName   string
	Port     string
}
