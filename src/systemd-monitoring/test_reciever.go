package main

import "fmt"
import "time"
import "systemd-monitoring/common"
import "systemd-monitoring/reciever"

func main() {

    updates      := make(chan common.DataUpdate)
    r,err:=reciever.NewReciever("0.0.0.0:8080","hello",updates)
    go r.Run()
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
