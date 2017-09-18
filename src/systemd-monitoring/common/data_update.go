package common

var typeDataUpdate int = 2002

type DataUpdate struct {
    Hostname  string `json:"hostname"`
    Text      string `json:"text"`
    Timestamp string `json:"timestamp"`
}

type Message struct {
    Type int    `json:"type"`
    Data []byte `json:"data"`
}
