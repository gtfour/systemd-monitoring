package main

import "os"
import "os/signal"
import "fmt"
import "time"
import "systemd-monitoring/event"

func main() {

    go event.EventBus.Handle()
    eventWriteChan1,err1 := event.EventBus.GetEventsWritePipe()
    eventWriteChan2,err2 := event.EventBus.GetEventsWritePipe()
    eventOutChan1   := make(chan *event.Event)
    eventOutChan2   := make(chan *event.Event)
    err3:=event.EventBus.SubscribeEvents(eventOutChan1)
    err4:=event.EventBus.SubscribeEvents(eventOutChan2)
    if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
        return
    }
    var eventRecievedFrom int = 1
    for {
        select {
            case ev1:=<-eventOutChan1:
                fmt.Printf("event1: %v\n",ev1)
                eventRecievedFrom = 1
            case ev2:=<-eventOutChan2:
                fmt.Printf("event2: %v\n",ev2)
                eventRecievedFrom = 2
            default:
                if eventRecievedFrom == 1 {
                    newEvent,_:=event.NewEvent()
                    eventWriteChan2 <- newEvent
                } else {
                    newEvent,_:=event.NewEvent()
                    eventWriteChan1 <- newEvent
                }
                time.Sleep(time.Second * 1)
        }
    }
    catchExit()
}

func catchExit()(){
    signalChan  := make(chan os.Signal, 1)
    cleanupDone := make(chan bool)
    signal.Notify(signalChan, os.Interrupt)
    go func() {
        for _ = range signalChan {
            event.EventBus.Exit()
            cleanupDone <- true
            break
        }
    }()
    <-cleanupDone
    return
}
