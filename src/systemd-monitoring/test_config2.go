package main

import "fmt"
import "systemd-monitoring/config"
import "systemd-monitoring/event"

func main(){
    //
    //
    sampleConfig  := `[{"id":"python-traceback-file1.txt","check-timeout":"60","flush-after":"300","disable-on":"30","conditions":[{"id":"catched-python-traceback"}],"actions":[{"id":"service-cron-restart","area":"service","args":["cron","restart"]}]}]`
    eventConfigList,err := config.ParseEvents(sampleConfig)
    events,_ := event.InputEvents(eventConfigList)
    fmt.Printf(">>\nEvents:\n%v\n%v\n",events,err)
    for i := range events {
        e := events[i]
        fmt.Printf("\tId: %v\n",e.Id)
        fmt.Printf("\tActionSet: %v\n",e.ActionSet)
        fmt.Printf("\tConditionSet: %v\n",e.ConditionSet)
    }
    //
    //
}
