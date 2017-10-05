package event

type ConditionSet struct {
    id         string
    conditions []Condition
}

type Condition struct {
    id        string
    satisfied bool
}

func(c *Condition)Satisfy()(){

}
