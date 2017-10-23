package config

import "strings"
import "encoding/json"

// Loading Events from json

type ActionConfig struct {
    Id        string   `json:"id"`
    //EventId   string   `json:"event-id"`
    Area      string   `json:"area"`
    Type      string   `json:"type"`
    Args      []string `json:"args"`
}

type ConditionConfig struct {
    Id        string   `json:"id"`
    //EventId   string   `json:"event-id"`
    Satisfied string   `json:"satisfied"`
    Count     string   `json:"count"`
    Area      string   `json:"area"`
    Args      []string `json:"args"`
}

type EventConfig struct {
   Id            string            `json:"id"`
   Actions       []ActionConfig    `json:"actions"`
   Conditions    []ConditionConfig `json:"conditions"`
   CheckTimeout  string            `json:"check-timeout"`
   FlushAfter    string            `json:"flush-after"`
   DisableOn     string            `json:"disable-on"`
}

func ParseEvents(events_string string)(event_list []EventConfig,err error){
    events_string = strings.Replace(events_string, "'", `"`, -1)
    err           = json.Unmarshal([]byte(events_string),&event_list)
    return
}
