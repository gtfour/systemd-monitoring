package event

import "sync"
import "errors"

var conditionNotFound = errors.New("Condition not found")
var conditionSetIsNil = errors.New("ConditionSet is nil")

type ConditionSet struct {
    sync.RWMutex
    id         string
    conditions []Condition
}

type Condition struct {
    id        string
    event_id  string
    satisfied bool
}


func(c *ConditionSet)IsSatisfied()(yes bool){
    yes = true
    for i := range c.conditions {
        condition := c.conditions[i]
        if !condition.satisfied {
            yes = false
            break
        }
    }
    return
}

func(c *ConditionSet)setConditionById(condition_id string,state bool)(error){
    c.Lock()
    defer c.Unlock()
    found := false
    for i := range c.conditions {
        condition := c.conditions[i]
        if condition.id == condition_id {
            condition.satisfied = state
            found               = true
            break
        }
    }
    if found == false {  return conditionNotFound } else { return nil }
}


func(c *Condition)satisfy()(){
    c.satisfied = true
}
