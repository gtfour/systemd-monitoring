package main

import "os"
import "os/signal"
import "systemd-monitoring/capture"

func main(){

    r,err:=capture.NewRunner([]string{"/tmp/my1.log","/tmp/my2.log","/tmp/my3.log"})
    if err == nil {
        r.Handle()
    }
    catchExit(r.MainQuit, r.QuitDone)
}

func catchExit(quit chan bool,quitDone chan bool)(){

    signalChan  := make(chan os.Signal, 1)
    cleanupDone := make(chan bool)
    signal.Notify(signalChan, os.Interrupt)
    go func() {
        for _ = range signalChan {
            quit <- true
            <-quitDone
            cleanupDone <- true
            break
        }
    }()
    <-cleanupDone
    return

}

