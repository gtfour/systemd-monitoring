package common

import "systemd-monitoring/config"

type AgentConfig struct {
    SecretPhrase     string
    MasterAddress    string
    FilesList        []string
    NginxLogs        []string
    ServiceList      []string
    PythonTracebacks []string
    DockerEvents     bool
    Monitors         config.Monitors
}

