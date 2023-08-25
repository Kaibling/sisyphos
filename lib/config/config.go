package config

import (
	"os"
	"strconv"
	"strings"
)

var Config *Configuration
var OS_PREFIX = "SISYPHOS"

type Configuration struct {
	GlobalVars          map[string]interface{}
	AdminPassword       string
	BindingPort         string
	BindingIP           string
	DBUser              string
	DBHost              string
	DBPort              string
	DBPassword          string
	DBDatabase          string
	DBDialect           string
	ClusterHeatBeatRate int
	ClusterEnabled      bool
}

func Init() {
	clusterHeartBeatRate, err := strconv.Atoi(getEnv("CLUSTER_HEARTBEATRATE", "2000"))
	if err != nil {
		clusterHeartBeatRate = 2000
	}
	clusterEnabled := parseStringBool(getEnv("CLUSTER_ENABLED", "false"))
	Config = &Configuration{
		GlobalVars: map[string]interface{}{
			"ssh_user":     "test",
			"ssh_password": "test",
		},
		AdminPassword:       getEnv("ADMIN_PASSWORD", "adminpassword"),
		BindingIP:           getEnv("BINDING_IP", "0.0.0.0"),
		BindingPort:         getEnv("BINDING_PORT", "7800"),
		DBUser:              getEnv("DB_USER", "db"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBPassword:          getEnv("DB_PASSWORD", "example"),
		DBHost:              getEnv("DB_HOST", "db"),
		DBDatabase:          getEnv("DB_DATABASE", "db"),
		DBDialect:           getEnv("DB_DIALECT", "postgres"),
		ClusterHeatBeatRate: clusterHeartBeatRate,
		ClusterEnabled:      clusterEnabled,
	}
}

func getEnv(key string, defaultValue string) string {
	fullKey := OS_PREFIX + "_" + key
	val := os.Getenv(OS_PREFIX + "_" + key)
	if val == "" {
		if defaultValue != "" {
			return defaultValue
		}
		panic(fullKey + " is not set")
	}
	return val
}

func parseStringBool(s string) bool {
	return strings.EqualFold(s, "true")
}
