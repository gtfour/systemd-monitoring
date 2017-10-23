package event

import "errors"
import "systemd-monitoring/common"

var ACTION_STATE_ACTIVATED int = 7002
var ACTION_STATE_RUNNING   int = 7004
var ACTION_STATE_PENDING   int = 7006
var ACTION_STATE_FAILED    int = 7008

var actionAlreadyExists = errors.New("Action with such id is already exist in action set")

type ActionSet struct {
    //
    id      string
    actions []Action
    //
}

type Action struct {
    //
    Id        string   `json:"id"`
    EventId   string   `json:"event-id"`
    Area      string   `json:"area"`
    Type      string   `json:"type"`
    Args      []string `json:"args"`
    Count     int
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

func(e *Event)AppendAction(a Action)(error){
    new_action_id := a.Id
    for i := range e.ActionSet.actions {
        ac := e.ActionSet.actions[i]
        if new_action_id == ac.Id {
            return actionAlreadyExists
        }
    }
    a.EventId           = e.Id
    e.ActionSet.actions = append(e.ActionSet.actions, a)
    return nil
}

