package main

import "fmt"
import "time"
import "systemd-monitoring/common"
import "systemd-monitoring/reciever"

func main() {

    updates := make(chan common.DataUpdate,100)
    r,err   := reciever.NewReciever("0.0.0.0:8080","1234567890123456",updates)
    go r.Run()
    if err == nil {
        for {
            select {
                case u:=<-updates:
                    fmt.Printf("Update:\n%s",u)
                default:
                    time.Sleep(time.Second * 2)
            }
        }
    }
}
