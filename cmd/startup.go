package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tnosaj/mssql_exporter/internal"
)

func evaluateInputs() (internal.Settings, error) {
	var s internal.Settings
	var includeMetrics string
	flag.BoolVar(&s.Debug, "v", false, "Enable verbose debugging output")
	flag.StringVar(&s.Port, "p", "8080", "Starts server on this port")
	flag.StringVar(&includeMetrics, "i", "exec,filespace,index,locks,memory,schedulers,tasks,waits", "Comma seperated list of metrics to gather. Possible values are [exec,filespace,index,locks,memory,performance,schedulers,tasks,waits]")
	flag.IntVar(&s.Timeout, "t", 10, "Timeout in seconds for a backend answer")

	flag.StringVar(&s.MetricsPath, "u", "/metrics", "Url for the user service")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [flags] command [command args…]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	setupLogger(s.Debug)

	s.EnabledMetrics = strings.Split(includeMetrics, ",")

	var err error
	s.DBConnectionInfo.User, err = getEnvVar("DBUSER")
	if err != nil {
		return internal.Settings{}, err
	}
	s.DBConnectionInfo.Password, err = getEnvVar("DBPASSWORD")
	if err != nil {
		return internal.Settings{}, err
	}
	s.DBConnectionInfo.HostName, err = getEnvVar("DBHOSTNAME")
	if err != nil {
		return internal.Settings{}, err
	}
	s.DBConnectionInfo.DBName, err = getEnvVar("DBNAME")
	if err != nil {
		return internal.Settings{}, err
	}
	s.DBConnectionInfo.Port, err = getEnvVar("DBPORT")
	if err != nil {
		return internal.Settings{}, err
	}

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
