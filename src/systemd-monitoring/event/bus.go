package event

import "errors"

var busNotReady = errors.New("Bus is not ready")

type Bus struct {
    //
    events          chan Event
    eventsIn        chan Event
    eventsOut       []chan Event
    //
    actionSets      chan ActionSet
    actionSetsIn    chan ActionSet
    actionSetsOut   []chan ActionSet
    //
    conditionSet    chan ConditionSet
    conditionSetIn  chan ConditionSet
    conditionSetOut []chan ConditionSet
    //
    ready bool
}

func NewBus()(*Bus) {
    var bus Bus
    bus.events          = make(chan   Event)
    bus.eventsIn        = make(chan   Event)
    bus.eventsOut       = make([]chan Event,0)
    bus.actionSets      = make(chan   ActionSet)
    bus.actionSetsIn    = make(chan   ActionSet)
    bus.actionSetsOut   = make([]chan ActionSet,0)
    bus.conditionSet    = make(chan   ConditionSet)
    bus.conditionSetIn  = make(chan   ConditionSet)
    bus.conditionSetOut = make([]chan ConditionSet,0)
    bus.ready           = true
    return &bus


}

func(b *Bus)SubscribeEvents(eventsOut chan Event)(err error){
    if b.ready {
        return nil
    } else {
        return busNotReady
    }

}

func(b *Bus)GetEventsWritePipe()(chan Event,error){
    if b.ready {
        return b.eventsIn, nil
    } else {
        return nil,busNotReady
    }
}
