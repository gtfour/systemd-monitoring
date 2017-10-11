package event

//import "sync"
import "systemd-monitoring/common"

type Event struct {
   //
   // sync.RWMutex
   id            string
   actionSet     ActionSet
   conditionSet  ConditionSet
   CheckTimeout  int
   FlushAfter    int           // setting all conditions to unsatisfied after that's  minutes 
   DisableOn     int           // disable event on this period of time
   events        chan *Event
   actionsIn     chan Action
   conditionsOut chan Condition
   state         int
   //
   //
}

func(e *Event)Handle()(){
    //
    for {
        select {
            case c:=<-e.conditionsOut:
                if e.id == c.event_id {
                    e.conditionSet.setConditionById(c.id,c.satisfied)
                }
            default:
                conditions_is_satisfied := e.conditionSet.IsSatisfied()
                if conditions_is_satisfied {

                } else {

                }

        }
    }
    //
}

func(e *Event)SatisfyCondition(condition_id string)(error){
    err := e.conditionSet.setConditionById(condition_id,true)
    return err
}


func(b *Bus)NewEvent(props ...string)(*Event,error){
    if b.ready {
        var e Event
        e.events        = b.events
        e.conditionsOut = make(chan Condition,10)
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
