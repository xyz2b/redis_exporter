package main

import (
	"github.com/tkanos/gonfig"
)

var (
	config        redisExporterConfig
	defaultOption = redisExporterConfig{
		RedisAddr: "127.0.0.1:6379",
		RedisUser:	"",
		RedisPwd:	"",
		SentinelClusterName: "",
		ListenAddress:	":9121",
		Namespace:	"redis",
		CheckKeys:	"",
		CheckSingleKeys: "",
		CheckStreams: "",
		CheckSingleStreams: "",
		CountKeys: "",
		TlsClientKeyFile: "",
		TlsClientCertFile: "",
		TlsCaCertFile: "",
		ScriptPath: "",
		ConnectionTimeout: "15s",
		MetricPath: "/metrics",
		LogFormat: "txt",
		ConfigCommand: "CONFIG",
		TlsServerKeyFile: "",
		TlsServerCertFile: "",
		IsDebug: false,
		SetClientName: true,
		IsTile38: false,
		ExportClientList: false,
		ShowVersion: false,
		RedisMetricsOnly: false,
		PingOnConnect: false,
		InclSystemMetrics: false,
		SkipTLSVerification: false,
		SubSystemName:      "",
		SubSystemID:        "",
		ClusterName:		"",
	}
)

type redisExporterConfig struct {
	RedisAddr               string	  		        `json:"redis_addr"`
	RedisUser				string					`json:"redis_user"`
	RedisPwd				string					`json:"redis_pwd"`
	SentinelClusterName		string					`json:"sentinel_cluster_name"`
	ListenAddress			string					`json:"listen_address"`
	Namespace				string					`json:"namespace"`
	CheckKeys				string					`json:"check_keys"`
	CheckSingleKeys			string					`json:"check_single_keys"`
	CheckStreams			string					`json:"check_streams"`
	CheckSingleStreams		string					`json:"check_single_streams"`
	CountKeys				string					`json:"count_keys"`
	ScriptPath				string					`json:"script_path"`
	MetricPath				string					`json:"metric_path"`
	LogFormat				string					`json:"log_format"`
	IsDebug					bool					`json:"is_debug"`
	ConfigCommand			string					`json:"config_command"`
	ConnectionTimeout		string					`json:"connection_timeout"`
	TlsClientKeyFile		string					`json:"tls_client_key_file"`
	TlsClientCertFile		string					`json:"tls_client_cert_file"`
	TlsCaCertFile			string					`json:"tls_ca_cert_file"`
	TlsServerKeyFile		string					`json:"tls_server_key_file"`
	TlsServerCertFile		string					`json:"tls_server_cert_file"`
	SetClientName			bool					`json:"set_client_name"`
	IsTile38				bool					`json:"is_tile_38"`
	ExportClientList		bool					`json:"export_client_list"`
	ShowVersion				bool					`json:"show_version"`
	RedisMetricsOnly		bool					`json:"redis_metrics_only"`
	PingOnConnect			bool					`json:"ping_on_connect"`
	InclSystemMetrics		bool					`json:"incl_system_metrics"`
	SkipTLSVerification		bool					`json:"skip_tls_verification"`
	SubSystemName			string					`json:"sub_system_name"`
	SubSystemID				string					`json:"sub_system_id"`
	ClusterName				string					`json:"cluster_name"`
}

func initConfigFromFile(config_file string) error {
	config = defaultOption
	err := gonfig.GetConf(config_file, &config)
	if err != nil {
		return err
	}

	return nil
}
