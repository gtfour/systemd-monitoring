package common

var TypeDataUpdate int = 2002

type DataUpdate struct {
    Hostname  string `json:"hostname"`
    Area      string `json:"area"`
    Text      string `json:"text"`
    Timestamp string `json:"timestamp"`
}

func(d *DataUpdate)String() string {
    return "Hostname:"+d.Hostname+"\n"+"Area:"+d.Area+"\n"+d.Text+"\nTimestamp:"+d.Timestamp+"\n"
}


type Message struct {
    Type int    `json:"type"`
    Data []byte `json:"data"`
}