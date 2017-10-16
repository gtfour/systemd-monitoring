package event

import "systemd-monitoring/common"

var ACTION_STATE_ACTIVATED int = 7002
var ACTION_STATE_RUNNING   int = 7004
var ACTION_STATE_PENDING   int = 7006

type ActionSet struct {
    //
    id      string
    actions []Action
    //
}

type Action struct {
    //
    Id        string `json:"id"`
    EventId   string `json:"event-id"`
    state     int
    //
}

func(a *Action)activate()(){
    a.state = ACTION_STATE_ACTIVATED
}

func(a *Action)run()(){
    a.state = ACTION_STATE_RUNNING
}

func(a *Action)finish()(){
    a.state = ACTION_STATE_PENDING
}

func(e *Event)NewAction(props ...string)(a Action,err error){
    if e.state ==  EVENT_NEW || e.state == EVENT_RUNNING {
        if len(props)>0 {
            name      := props[0]
            a.Id      =  name
            a.EventId =  e.Id
            e.ActionSet.actions = append(e.ActionSet.actions, a)
            return a,nil
        }
        action_id, err := common.GenId()
        if err == nil {
            a.Id      = action_id
            a.EventId = e.Id
            e.ActionSet.actions = append(e.ActionSet.actions, a)
            return a, nil
        } else {
            return a, err
        }
    } else {
        return a, eventStateUndefined
    }
}
