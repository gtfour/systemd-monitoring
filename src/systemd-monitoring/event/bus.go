package event

import "fmt"
import "time"
import "sync"
import "errors"

var busNotReady         = errors.New("Bus is not ready")
var chanIsNil           = errors.New("Chan is nil")
var eventNotFound       = errors.New("Event with such id is not found")
var eventIsNil          = errors.New("Event is nil")
var eventIsAlreadyExist = errors.New("Event with such id is already in list")
var EventBus            = NewBus()

type Bus struct {
    //
    sync.RWMutex
    eventsList       []*Event
    events           chan   *Event
    eventsIn         chan   *Event
    eventsOut        []chan *Event
    /*
    actionSets       chan ActionSet
    actionSetsIn     chan ActionSet
    actionSetsOut    []chan ActionSet
    conditionSet     chan ConditionSet
    conditionSetIn   chan ConditionSet
    conditionSetOut  []chan ConditionSet
    */
    conditionsIn     chan     Condition
    conditionsOut    []chan   Condition
    actionsIn        chan     Action
    actionsOut       []chan   Action
    //
    quitCh           chan bool
    timeout_sec      time.Duration
    ready            bool
    //
}

func NewBus()(*Bus) {
    //
    var bus Bus
    bus.eventsList      = make([]*Event,0)
    bus.events          = make(chan   *Event,100)
    bus.eventsIn        = make(chan   *Event,100)
    bus.eventsOut       = make([]chan *Event,0)
    /*
    bus.actionSets      = make(chan   ActionSet)
    bus.actionSetsIn    = make(chan   ActionSet)
    bus.actionSetsOut   = make([]chan ActionSet,0)
    bus.conditionSet    = make(chan   ConditionSet)
    bus.conditionSetIn  = make(chan   ConditionSet)
    bus.conditionSetOut = make([]chan ConditionSet,0)
    */
    //
    //
    bus.conditionsIn    = make(chan   Condition)
    bus.conditionsOut   = make([]chan Condition,100)
    bus.actionsIn       = make(chan   Action)
    bus.actionsOut      = make([]chan Action,100)
    //
    //
    bus.quitCh          = make(chan bool)
    bus.timeout_sec     = 2
    bus.ready           = true
    return &bus
    //
}

func(b *Bus)SubscribeEvents(eventsOutSingle chan *Event)(err error){
    if b.ready {
        if eventsOutSingle == nil { return chanIsNil }
        b.eventsOut = append(b.eventsOut, eventsOutSingle)
        return nil
    } else {
        return busNotReady
    }
}

func(b *Bus)SubscribeConditions(conditionsOutSingle chan Condition)(err error){
    if b.ready {
        if conditionsOutSingle == nil { return chanIsNil }
        b.conditionsOut = append(b.conditionsOut, conditionsOutSingle)
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
            //case e:=<-b.events:
            //    for i := range b.eventsOut {
            //        eventsOut := b.eventsOut[i]
            //        eventsOut <- e
            //    }
            case eIn:=<-b.eventsIn:
                err := b.AppendEvent(eIn)
                if err == nil {
                    go eIn.Handle()
                }
            case cIn:=<-b.conditionsIn:
                for i := range b.conditionsOut {
                    conditionsOutSingle := b.conditionsOut[i]
                    conditionsOutSingle<-cIn
                }
            case aIn:=<-b.actionsIn:
                for i := range b.actionsOut {
                    actionsOut := b.actionsOut[i]
                    actionsOut<-aIn
                }
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

func(b *Bus)GetEvent(event_id string)(event_copy Event, err error){
    //
    b.Lock()
    defer b.Unlock()
    for i := range b.eventsList {
        event_ptr := b.eventsList[i]
        fmt.Printf("event_ptr.id: %v event_id: %v\n",event_ptr.id,event_id)
        if event_ptr!=nil && event_ptr.id==event_id {
            event_copy =*event_ptr
            return event_copy,nil
        }
    }
    return event_copy,eventNotFound
    //
}

func(b *Bus)SetCondition(event_id string,condition_id string,state bool)(err error){
    b.Lock()
    defer b.Unlock()
    for i := range b.eventsList {
        ev := b.eventsList[i]
        if ev.id == event_id {
            err = ev.conditionSet.setConditionById(condition_id,state)
            return err
        }
    }
    return eventNotFound
}


func(b *Bus)AppendEvent(new_event *Event)(error){
    b.Lock()
    defer b.Unlock()
    if new_event == nil { return eventIsNil }
    new_event_id := new_event.id
    for i := range b.eventsList {
        ev := b.eventsList[i]
        if new_event_id == ev.id {
            return eventIsAlreadyExist
        }
    }
    b.eventsList = append(b.eventsList, new_event)
    b.SubscribeConditions(new_event.conditionsOut)
    return nil
}
