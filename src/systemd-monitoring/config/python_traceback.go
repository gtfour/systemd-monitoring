package config

import "strings"
import "encoding/json"

type PythonTracebackHandlerConfig struct {
    Path     string   `json:"path"`
    Keywords []string `json:"keywords"`
}

func ParsePythonTracebackHandlerConfig(python_tracebacks_string string)(list []PythonTracebackHandlerConfig, err error){
    python_tracebacks_string = strings.Replace(python_tracebacks_string, "'", `"`, -1)
    err=json.Unmarshal([]byte(python_tracebacks_string), &list)
    return
}



