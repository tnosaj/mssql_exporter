package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tnosaj/mssql_exporter/internal"
)

func evaluateInputs() (internal.Settings, error) {
	var s internal.Settings
	//	var err error

	flag.BoolVar(&s.Debug, "v", false, "Enable verbose debugging output")
	flag.StringVar(&s.Port, "p", "8080", "Starts server on this port")
	flag.IntVar(&s.Timeout, "t", 1, "Timeout in seconds for a backend answer")

	flag.StringVar(&s.MetricsPath, "u", "/metrics", "Url for the user service")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [flags] command [command args…]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	setupLogger(s.Debug)

	return s, nil
}

func setupLogger(debug bool) {
	logrus.SetReportCaller(false)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.Debug("Configured logger")
}
func getEnvVar(envVarName string) (string, error) {
	check := os.Getenv(envVarName)
	if check == "" {
		printRequiredEnvVars()
		return "", fmt.Errorf("Missing env var %q", envVarName)
	}
	return check, nil
}

func printRequiredEnvVars() {
	fmt.Println("\nRequired Environment variables:")
	fmt.Println("  XXXXX")
}
