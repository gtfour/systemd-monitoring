package event

//import "sync"
import "systemd-monitoring/common"

type Event struct {
   //
   // sync.RWMutex
   id           string
   actionSet    ActionSet
   conditionSet ConditionSet
   CheckTimeout int
   FlushAfter   int           // setting all conditions to unsatisfied after that's  minutes 
   DisableOn    int           // disable event on this period of time
   events       chan *Event
   state        int
   //
   //
}

func(e *Event)Handle()(){
    //
    for {
        conditions_is_satisfied := e.conditionSet.IsSatisfied()
        if conditions_is_satisfied {

        } else {

        }
    }
    //
}

func(e *Event)SatisfyCondition(condition_id string)(error){
    err := e.conditionSet.satisfy(condition_id)
    return err
}


func(b *Bus)NewEvent(props ...string)(*Event,error){
    if b.ready {
        var e Event
        e.events = b.events
        if len(props)>0 {
            name := props[0]
            e.id =  name
            return &e,nil
        }
        event_id,err := common.GenId()
        if err == nil {
            e.id = event_id
            return &e,nil
        } else {
            return nil,err
        }
    } else {
        return nil, busNotReady
    }
}
