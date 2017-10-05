package common

type AgentConfig struct {
    SecretPhrase     string
    MasterAddress    string
    FilesList        []string
    NginxLogs        []string
    ServiceList      []string
    PythonTracebacks []string
    DockerEvents     bool
}

