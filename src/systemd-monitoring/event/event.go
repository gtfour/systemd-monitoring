package event

//import "sync"
import "fmt"
import "time"
import "systemd-monitoring/common"

var EVENT_NEW      int = 8002
var EVENT_RUNNING  int = 8004
var EVENT_SLEEPING int = 8006

type Event struct {
   // sync.RWMutex
   Id            string        `json:"id"`
   ActionSet     ActionSet
   ConditionSet  ConditionSet
   CheckTimeout  time.Duration `json:"check-timeout"`
   FlushAfter    time.Duration `json:"flush-after"` // setting all conditions to unsatisfied after that's  minutes 
   DisableOn     time.Duration `json:"disable-on"`  // disable event on this period of time (usually used after all conditions are marked as satisfied)
   events        chan *Event
   actionsIn     chan Action
   conditionsOut chan Condition
   state         int
   //
}

func(e *Event)Handle()(){
    for {
        select {
            case c:=<-e.conditionsOut:
                fmt.Printf("\t--event %v: recieved condition : %v\n",e.Id,c)
                if e.Id == c.EventId && e.state != EVENT_SLEEPING  {
                    e.ConditionSet.setConditionById(c.Id, c.Satisfied)
                    //fmt.Printf("Setting condition err: %v\n",err)
                }
            default:
                if e.state != EVENT_SLEEPING {
                    conditions_are_satisfied := e.ConditionSet.IsSatisfied()
                    fmt.Printf("\t-- event_id: %v conditions are satisfied: %v\n",e.Id,conditions_are_satisfied)
                    if conditions_are_satisfied {
                        fmt.Printf("\t++ trying to execute following actions: %v\n",e.ActionSet.actions)
                        for i:= range e.ActionSet.actions {
                            action := e.ActionSet.actions[i]
                            fmt.Printf("\t++ Sending action: %v <<\n",action)
                            e.actionsIn <- action
                        }
                        e.state = EVENT_SLEEPING
                        e.ConditionSet.SetAll(false)
                        go func() {
                            time.Sleep(time.Second * e.DisableOn)
                            e.state = EVENT_RUNNING
                            //e.ConditionSet.SetAll(false)
                        }()
                    } else { time.Sleep(time.Second * e.CheckTimeout) }
                } else {
                    time.Sleep(time.Second * e.CheckTimeout)
                }
        }
    }
}

func(e *Event)SatisfyCondition(condition_id string)(error){
    err := e.ConditionSet.setConditionById(condition_id,true)
    return err
}


func(b *Bus)NewEvent(props ...string)(e_ptr *Event,err error){
    if b.ready {
        var e Event
        e.CheckTimeout  = 10
        e.FlushAfter    = 10
        e.DisableOn     = 10
        e.events        = b.events
        e.conditionsOut = make(chan Condition, 10)
        e.actionsIn,err = b.GetActionsWritePipe()
        if err != nil {
            return nil,err
        }
        e.state         =  EVENT_NEW
        if err != nil { return nil, err }
        if len(props)>0 {
            name := props[0]
            e.Id =  name
            return &e,nil
        }
        event_id,err := common.GenId()
        if err == nil {
            e.Id = event_id
            return &e,nil
        } else {
            return nil,err
        }
    } else {
        return nil, busNotReady
    }
}
