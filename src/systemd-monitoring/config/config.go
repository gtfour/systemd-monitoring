package config

import "strings"
import "encoding/json"

type Monitors []Monitor

var NGINX_LOG_MONITOR_TYPE string = "nginx-log-monitor"
var FILE_MONITOR_TYPE      string = "file-monitor"

/*
type Monitors struct {
    Monitors []FileMonitor `json:"monitors"`
}
*/

type Monitor struct {
    Type    string `json:"type"`
    Monitor string `json:"monitor"`
}

type FileMonitor struct {
    Path         string   `json:"path"`
    IgnoreString []string `json:"ignore-string"`
    MatchString  []string `json:"match-string"`
    Timeout      string   `json:"timeout"`
}

type NginxLogMonitor struct {
    Path         string   `json:"path"`
    IgnoreStatus []string `json:"ignore-status"`
    MatchStatus  []string `json:"match-status"`
    Timeout      string   `json:"timeout"`
}


func ParseMonitors(monitors_string string)(ml Monitors,err error){
    err=json.Unmarshal([]byte(monitors_string),&ml)
    return
}

func ParseFileMonitor(file_monitor_string string)(fm FileMonitor, err error){
    file_monitor_string = strings.Replace(file_monitor_string, "'", `"`, -1)
    err=json.Unmarshal([]byte(file_monitor_string),&fm)
    return
}

func ParseNginxLogMonitor(nginx_log_monitor_string string)(nm NginxLogMonitor, err error){
    nginx_log_monitor_string = strings.Replace(nginx_log_monitor_string, "'", `"`, -1)
    err=json.Unmarshal([]byte(nginx_log_monitor_string),&nm)
    return
}
