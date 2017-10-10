package main

import "fmt"
import "systemd-monitoring/config"

func main() {
    //sampleConfig:=`[{"type":"file-monitor","monitor":"{'ignore-string':['No such file or directory']}"},{"type":"nginx-log-monitor","monitor":"{'match-status':['500','502']}"}]`
    sampleEmptyConfig:=""
    ml,err:=config.ParseMonitors(sampleEmptyConfig)
    for i := range ml {
        m := ml[i]
        fmt.Printf("Monitor: type: %s monitor_body: %v\n",m.Type,m.Monitor)
        switch m.Type {
            case config.FILE_MONITOR_TYPE:
                fm,err := config.ParseFileMonitor(m.Monitor)
                fmt.Printf("fm: %v err: %v\n",fm,err)
            case config.NGINX_LOG_MONITOR_TYPE:
                nm,err := config.ParseNginxLogMonitor(m.Monitor)
                fmt.Printf("nm: %v err: %v\n",nm,err)
            default:
                fmt.Printf("Unrecognized type:\n")
        }

    }
    fmt.Printf("Monitors:\n%v\n%v\n",ml,err)
}
