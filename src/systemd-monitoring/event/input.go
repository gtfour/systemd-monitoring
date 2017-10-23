package event

import "time"
import "strconv"
import "systemd-monitoring/config"

func InputEvents(eventConfigs []config.EventConfig)(events_list []*Event,err error){
    events_list = make([]*Event, 0)
    for i := range eventConfigs {
        //
        ec               := eventConfigs[i]
        checkTimeoutStr  := ec.CheckTimeout
        flushAfterStr    := ec.FlushAfter
        disableOnStr     := ec.DisableOn
        checkTimeout,err := strconv.Atoi(checkTimeoutStr)
        if err!=nil      { continue }
        flushAfter,err   := strconv.Atoi(flushAfterStr)
        if err!=nil      { continue }
        disableOn,err    := strconv.Atoi(disableOnStr)
        if err!=nil      { continue }
        id               := ec.Id
        if id == ""      { continue /* event with empty id is denied */ }
        var new_event Event
        new_event.Id           = id
        new_event.CheckTimeout = time.Duration(checkTimeout)
        new_event.FlushAfter   = time.Duration(flushAfter)
        new_event.DisableOn    = time.Duration(disableOn)
        if len(ec.Actions) > 0 {
            new_event.ActionSet.actions = make([]Action, 0)
            new_event.InputActions(ec.Actions)
        }
        if len(ec.Conditions) > 0 {
            new_event.ConditionSet.conditions = make([]Condition, 0)
            new_event.InputConditions(ec.Conditions)
        }
        events_list = append(events_list, &new_event)
        //
    }
    return
}

func(e *Event)InputActions(actionConfigs []config.ActionConfig)(err error){
    //
    if e == nil { return eventIsNil }
    //
    for i := range actionConfigs {

        ac         := actionConfigs[i]
        idStr      := ac.Id
        areaStr    := ac.Area
        args       := ac.Args
        if idStr == "" { continue }
        var new_action Action
        new_action.Id     = idStr
        new_action.Args   = args
        new_action.Area   = areaStr


        e.AppendAction(new_action)

    }
    return
    //
}

func(e *Event)InputConditions(conditionConfigs []config.ConditionConfig)(err error){
    //
    if e == nil { return eventIsNil }

    for i := range conditionConfigs {
        con     := conditionConfigs[i]
        idStr   := con.Id
        if idStr == "" { continue }
        areaStr    := con.Area
        countStr   := con.Count
        args       := con.Args
        count, err := strconv.Atoi(countStr)
        if err != nil && countStr != "" { continue }
        if err !=nil  { count = 0 }
        var new_condition Condition
        new_condition.Id    = idStr
        new_condition.Area  = areaStr
        new_condition.Args  = args
        new_condition.Count = count
        e.AppendCondition(new_condition)
    }
    return
    //
}
