package event

import "fmt"
import "sync"
import "errors"
import "systemd-monitoring/common"

var conditionNotFound       = errors.New("Condition not found")
var conditionSetIsNil       = errors.New("ConditionSet is nil")
var eventStateUndefined     = errors.New("Unable to append condition to event. Event state is undefined")
var conditionIsAlreadyExist = errors.New("Condition with such id is already in list")

type ConditionSet struct {
    //
    sync.RWMutex
    id         string
    conditions []Condition
    //
}

type Condition struct {
    //
    Id        string   `json:"id"`
    EventId   string   `json:"event-id"`
    Area      string   `json:"area"`
    Args      []string `json:"args"`
    Satisfied bool
    Count     int      `json:"count"`
    //
}


func(c *ConditionSet)IsSatisfied()(yes bool){
    c.Lock()
    defer c.Unlock()
    yes = true
    for i := range c.conditions {
        condition := c.conditions[i]
        if !condition.Satisfied {
            yes = false
            break
        }
    }
    return
}

func(c *ConditionSet)SetAll(state bool)(){
    c.Lock()
    defer c.Unlock()
    for i := range c.conditions {
        condition           := c.conditions[i]
        condition.Satisfied =  state
        c.conditions[i]     =  condition
    }
    return
}


func(c *ConditionSet)setConditionById(condition_id string,state bool)(error){
    //
    c.Lock()
    defer c.Unlock()
    fmt.Printf("\t== ConditionsSet before %v\n",c.conditions)
    var new_condition       Condition
    found := false
    var new_condition_index int
    for i := range c.conditions {
        condition := c.conditions[i]
        if condition.Id == condition_id {
            //condition.Satisfied = state
            new_condition       = condition
            new_condition_index = i
            found               = true
            //c.conditions[i] = condition
            break
        }
    }
    fmt.Printf("\t== setConditionById: updated conditions set  %v\n",c.conditions)
    if found == false {
        return conditionNotFound
    } else {
        new_condition.Satisfied           = true
        c.conditions[new_condition_index] = new_condition
        fmt.Printf("\t== ConditionsSet after update  %v\n",c.conditions)
        return nil
    }
    //
}


func(c *Condition)satisfy()(){
    c.Satisfied = true
}

func(e *Event)NewCondition(props ...string)(c Condition,err error){
    if e.state ==  EVENT_NEW || e.state == EVENT_RUNNING {
        if len(props)>0 {
            name      := props[0]
            c.Id      =  name
            c.EventId =  e.Id
            e.ConditionSet.conditions = append(e.ConditionSet.conditions, c)
            return c, nil
        }
        condition_id,err := common.GenId()
        if err == nil {
            c.Id      = condition_id
            c.EventId = e.Id
            e.ConditionSet.conditions = append(e.ConditionSet.conditions, c)
            return c, nil
        } else {
            return c, err
        }
    } else {
        return c, eventStateUndefined
    }
}

func(e *Event)AppendCondition(c Condition)(error){
    new_condition_id := c.Id
    for i := range e.ConditionSet.conditions {
        con := e.ConditionSet.conditions[i]
        if new_condition_id == con.Id {
            return conditionIsAlreadyExist
        }
    }
    c.EventId                 = e.Id
    e.ConditionSet.conditions = append( e.ConditionSet.conditions, c )
    return nil
}

