package event

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
    id        string
    event_id  string
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


