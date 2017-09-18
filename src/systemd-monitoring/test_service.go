package main

import "fmt"
import "time"
import "systemd-monitoring/common"
import "systemd-monitoring/service"

func main() {
    // s,err:=service.CheckSystemdService("cron")
    // fmt.Printf("service: %v err %v\n",s,err)
    var services =  []string {"cron","pcap-log"}
    updates      := make(chan common.DataUpdate)
    chain,err:=service.NewServiceChain(services, updates)
    go chain.Proceed()
    if err == nil {
        for {
            select {
                case u:=<-updates:
                    fmt.Printf("Update:%v\n",u)
                default:
                    time.Sleep(time.Second * 2)
            }

        }
    }
}
