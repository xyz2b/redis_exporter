package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/oliver006/redis_exporter/exporter"
)

var (
	/*
		BuildVersion, BuildDate, BuildCommitSha are filled in by the build script
	*/
	BuildVersion   = "<<< filled in by build >>>"
	BuildDate      = "<<< filled in by build >>>"
	BuildCommitSha = "<<< filled in by build >>>"
)

func getEnv(key string, defaultVal string) string {
	if envVal, ok := os.LookupEnv(key); ok {
		return envVal
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if envVal, ok := os.LookupEnv(key); ok {
		envBool, err := strconv.ParseBool(envVal)
		if err == nil {
			return envBool
		}
	}
	return defaultVal
}

func main() {
	var configFile = flag.String("config-file", "./redis_exporter.conf", "path to json config")
	flag.Parse()

	err := initConfigFromFile(*configFile)
	if err!= nil {
		panic(err)
	}

	switch config.LogFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}
	log.Printf("Redis Metrics Exporter %s    build date: %s    sha1: %s    Go: %s    GOOS: %s    GOARCH: %s\nconfig: %s",
		BuildVersion, BuildDate, BuildCommitSha,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		config,
	)
	if config.IsDebug {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Enabling debug output")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if config.ShowVersion {
		return
	}

	to, err := time.ParseDuration(config.ConnectionTimeout)
	if err != nil {
		log.Fatalf("Couldn't parse connection timeout duration, err: %s", err)
	}

	var tlsClientCertificates []tls.Certificate
	if (config.TlsClientKeyFile != "") != (config.TlsClientCertFile != "") {
		log.Fatal("TLS client key file and cert file should both be present")
	}
	if config.TlsClientKeyFile != "" && config.TlsClientCertFile != "" {
		cert, err := tls.LoadX509KeyPair(config.TlsClientCertFile, config.TlsClientKeyFile)
		if err != nil {
			log.Fatalf("Couldn't load TLS client key pair, err: %s", err)
		}
		tlsClientCertificates = append(tlsClientCertificates, cert)
	}

	var tlsCaCertificates *x509.CertPool
	if config.TlsCaCertFile != "" {
		caCert, err := ioutil.ReadFile(config.TlsCaCertFile)
		if err != nil {
			log.Fatalf("Couldn't load TLS Ca certificate, err: %s", err)
		}
		tlsCaCertificates = x509.NewCertPool()
		tlsCaCertificates.AppendCertsFromPEM(caCert)
	}

	var ls []byte
	if config.ScriptPath != "" {
		if ls, err = ioutil.ReadFile(config.ScriptPath); err != nil {
			log.Fatalf("Error loading script file %s    err: %s", config.ScriptPath, err)
		}
	}

	registry := prometheus.NewRegistry()
	if !config.RedisMetricsOnly {
		registry = prometheus.DefaultRegisterer.(*prometheus.Registry)
	}

	exp, err := exporter.NewRedisExporter(
		config.RedisAddr,
		exporter.Options{
			User:                config.RedisUser,
			Password:            config.RedisPwd,
			Namespace:           config.Namespace,
			ConfigCommandName:   config.ConfigCommand,
			CheckKeys:           config.CheckKeys,
			CheckSingleKeys:     config.CheckSingleKeys,
			CheckStreams:        config.CheckStreams,
			CheckSingleStreams:  config.CheckSingleStreams,
			CountKeys:           config.CountKeys,
			LuaScript:           ls,
			InclSystemMetrics:   config.InclSystemMetrics,
			SetClientName:       config.SetClientName,
			IsTile38:            config.IsTile38,
			ExportClientList:    config.ExportClientList,
			SkipTLSVerification: config.SkipTLSVerification,
			ClientCertificates:  tlsClientCertificates,
			CaCertificates:      tlsCaCertificates,
			ConnectionTimeouts:  to,
			MetricsPath:         config.MetricPath,
			RedisMetricsOnly:    config.RedisMetricsOnly,
			PingOnConnect:       config.PingOnConnect,
			Registry:            registry,
			BuildInfo: exporter.BuildInfo{
				Version:   BuildVersion,
				CommitSha: BuildCommitSha,
				Date:      BuildDate,
			},
			SubSystemID: config.SubSystemID,
			SubSystemName: config.SubSystemName,
			ClusterName: config.ClusterName,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Providing metrics at %s%s", config.ListenAddress, config.MetricPath)
	log.Debugf("Configured redis addr: %#v", config.RedisAddr)
	if config.TlsServerCertFile != "" && config.TlsServerKeyFile != "" {
		log.Debugf("Bind as TLS using cert %s and key %s", config.TlsServerCertFile, config.TlsServerKeyFile)
		log.Fatal(http.ListenAndServeTLS(config.ListenAddress, config.TlsServerCertFile, config.TlsServerKeyFile, exp))
	} else {
		log.Fatal(http.ListenAndServe(config.ListenAddress, exp))
	}
}
