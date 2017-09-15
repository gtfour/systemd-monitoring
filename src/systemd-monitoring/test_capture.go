package main

import "os"
import "fmt"
import "os/signal"
import "systemd-monitoring/capture"

func main(){

    r,err:=capture.NewRunner([]string{"/tmp/my1.log","/tmp/my2.log","/tmp/my3.log"})
    if err == nil {
        go r.Handle()
    }
    catchExit(r.MainQuit, r.QuitDone)
}

func catchExit(quit chan bool,quitDone chan bool)(){
    signalChan  := make(chan os.Signal, 1)
    cleanupDone := make(chan bool)
    signal.Notify(signalChan, os.Interrupt)
    //signal.Notify(signalChan, os.Kill)
    go func() {
        for _ = range signalChan {
            fmt.Printf("Catched signal")
            quit <- true
            <-quitDone
            cleanupDone <- true
            break
        }
    }()
    <-cleanupDone
    return
}

