package main

import "os"
import "os/signal"
import "fmt"
import "time"
import "systemd-monitoring/event"
import "systemd-monitoring/common"

func main() {
    //
    go event.EventBus.Handle()
    eventsWriteChan,err_events         := event.EventBus.GetEventsWritePipe()
    conditionsWriteChan,err_conditions := event.EventBus.GetConditionsWritePipe()
    actionsOutChan                     := make(chan event.Action,100)
    err_actions_out                    := event.EventBus.SubscribeActions(actionsOutChan)
    if err_events!=nil || err_conditions!=nil || err_actions_out!=nil {
        fmt.Printf("Exiting:\nerr:%v\t%v\t%v\n",err_events,err_conditions,err_actions_out)
        return
    }
    newEvent1,_        := event.EventBus.NewEvent()
    newEvent2,_        := event.EventBus.NewEvent()

    con1,_ := newEvent1.NewCondition("python-traceback-first")
    con2,_ := newEvent1.NewCondition("python-traceback-second")
    con3,_ := newEvent2.NewCondition("error-entry#1")
    con4,_ := newEvent2.NewCondition("error-entry#2")

    var conditionsUpdate =  []event.Condition {con1, con2, con3, con4}


    _,_=newEvent1.NewAction("pcap-log-service-restart")
    _,_=newEvent2.NewAction("crontab-service-restart")

    customAction1 :=  event.Action{Id:"service-mongodb-restart",   Args:[]string{"mongodb","restart"}}
    customAction2 :=  event.Action{Id:"service-postgresql-restart",Args:[]string{"postgresql","restart"}}

    newEvent1.AppendAction(customAction1)
    newEvent2.AppendAction(customAction2)

    eventsWriteChan<-newEvent1
    eventsWriteChan<-newEvent2

    /*
    for i := range conditionsUpdate {
        con          := conditionsUpdate[i]
        con.Satisfied =  true
        conditionsWriteChan<-con
    }
    */

    for {
        select {
            case a:=<-actionsOutChan:
                fmt.Printf("++++ Action %v activated at: %v++++\n",a,common.GetTime())
            default:
                fmt.Printf(":'default_stage':\n")
                //
                for i := range conditionsUpdate {
                    con          := conditionsUpdate[i]
                    con.Satisfied =  true
                    conditionsWriteChan<-con
                }
                //
                time.Sleep(time.Second * 2)
        }
    }
    catchExit()
    //
}

func catchExit()(){
    //
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
    //
}
