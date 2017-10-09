package event

type ActionSet struct {
    id      string
    actions []Action
}

type Action struct {
    id        string
    event_id  string
    activated bool
}
