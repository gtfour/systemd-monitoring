package event

import "systemd-monitoring/common"

type Event struct {
   id           string
   actionSet    ActionSet
   conditionSet ConditionSet
}

func NewEvent()(*Event,error){
    event_id,err := common.GenId()
    if err == nil {
        var e Event
        e.id = event_id
        return &e,nil
    } else {
        return nil,err
    }
}
