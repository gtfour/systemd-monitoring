package main

import "fmt"
import "systemd-monitoring/service"

func main() {
    s,err:=service.CheckSystemdService("cron")
    fmt.Printf("service: %v err %v\n",s,err)
}
