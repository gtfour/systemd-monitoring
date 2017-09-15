package main

import "fmt"
import "systemd-monitoring/service"

func main() {
    pid,err:=service.GetServiceMainPid("cron")
    fmt.Printf("pid %d err %v\n",pid,err)
}
