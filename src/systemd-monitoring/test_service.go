package main

import "fmt"
import "time"
import "systemd-monitoring/common"
import "systemd-monitoring/service"
import "systemd-monitoring/notifier"

func main() {
    // s,err:=service.CheckSystemdService("cron")
    // fmt.Printf("service: %v err %v\n",s,err)
    var notifier notifier.NativeNotifier
    notifier.Address      = "127.0.0.1:8080"
    notifier.SecretPhrase = "1234567890123456"
    var services =  []string {"cron","pcap-log"}
    updates      := make(chan common.DataUpdate)
    chain,err:=service.NewServiceChain(services, updates)
    go chain.Proceed()
    if err == nil {
        for {
            select {
                case u:=<-updates:
                    err = notifier.Notify(u)
                    fmt.Printf("Update:%v\nnotifier_err:%v\n",u,err)
                default:
                    time.Sleep(time.Second * 2)
            }

        }
    }
}
