package common

var TypeDataUpdate int = 2002

type DataUpdate struct {
    Hostname  string `json:"hostname"`
    Area      string `json:"area"`
    Path      string `json:"path"`
    Text      string `json:"text"`
    Timestamp string `json:"timestamp"`
}

func(d *DataUpdate)String() string {
    if d.Area == "docker_events" || d.Area == "service" {
        return "Hostname:"+d.Hostname+"\n"+"Area:"+d.Area+"\n"+d.Text+"\nTimestamp:"+d.Timestamp+"\n"
    } else {
        return "Hostname:"+d.Hostname+"\n"+"Area:"+d.Area+"\n"+"Path:"+d.Path+"\n"+d.Text+"\nTimestamp:"+d.Timestamp+"\n"
    }
}


type Message struct {
    Type int    `json:"type"`
    Data []byte `json:"data"`
}
