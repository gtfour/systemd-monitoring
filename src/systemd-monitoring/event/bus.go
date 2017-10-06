package event

import "fmt"
import "time"
import "errors"

var busNotReady = errors.New("Bus is not ready")
var chanIsNil   = errors.New("Chan is nil")
var EventBus    = NewBus()

type Bus struct {
    events           chan *Event
    eventsIn         chan *Event
    eventsOut        []chan *Event
    actionSets       chan ActionSet
    actionSetsIn     chan ActionSet
    actionSetsOut    []chan ActionSet
    conditionSet     chan ConditionSet
    conditionSetIn   chan ConditionSet
    conditionSetOut  []chan ConditionSet
    quitCh           chan bool
    timeout_sec      time.Duration
    ready bool
}

func NewBus()(*Bus) {
    var bus Bus
    bus.events          = make(chan   *Event,100)
    bus.eventsIn        = make(chan   *Event,100)
    bus.eventsOut       = make([]chan *Event,0)
    bus.actionSets      = make(chan   ActionSet)
    bus.actionSetsIn    = make(chan   ActionSet)
    bus.actionSetsOut   = make([]chan ActionSet,0)
    bus.conditionSet    = make(chan   ConditionSet)
    bus.conditionSetIn  = make(chan   ConditionSet)
    bus.conditionSetOut = make([]chan ConditionSet,0)
    bus.quitCh          = make(chan bool)
    bus.timeout_sec     = 2
    bus.ready           = true
    return &bus


}

func(b *Bus)SubscribeEvents(eventsOut chan *Event)(err error){
    if b.ready {
        if eventsOut == nil { return chanIsNil }
        b.eventsOut = append(b.eventsOut, eventsOut)
        return nil
    } else {
        return busNotReady
    }
}

func(b *Bus)GetEventsWritePipe()(chan *Event,error){
    if b.ready {
        return b.eventsIn, nil
    } else {
        return nil,busNotReady
    }
}

func(b *Bus)Handle()(error){
    if !b.ready { return busNotReady }
    for {
        select {
            case e:=<-b.events:
                fmt.Printf("Event: %v\n",e)
                for i := range b.eventsOut {
                    eventOut:=b.eventsOut[i]
                    eventOut<-e
                }
            case eIn:=<-b.eventsIn:
                b.events<-eIn
            case <-b.quitCh:
                fmt.Printf("\n--Exiting--\n")
                break
            default:
                time.Sleep(time.Second * b.timeout_sec)
        }
    }
}

func(b *Bus)Exit()(){
    b.quitCh<-true
}
