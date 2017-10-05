package filter

import "fmt"
import "encoding/json"
import "jumper/cuda/targets"
import "jumper/cuda/handling"
import "jumper/cuda/filtering"

type Line struct {
    Data []string `json:"data"`
}

type PythonTraceback struct {
    Header string
    Lines  []string
    Footer string
}

type LogHandler struct {
    //
    handler  handling.Handler
    filters  filtering.FilterList
    target   *targets.Target
    //
}


func(l *LogHandler)HandleNginxLog(log_entry string)(request string, status string, beauty_log string){
    //
    //if l.target == nil {
    //    fmt.Printf("l.target is nil\n")
    //}
    //
    target_config          := make(map[string]string, 0)
    target_config["type"]  =  "SINGLE_LINE"
    tgt,_                  := targets.NewTarget(target_config)
    tgt.SetLine(log_entry)
    l.handler.AddTargetPtr(tgt)
    //
    res,err :=  l.handler.Handle()
    //
    if err == nil {
        var line Line
        result_js,err := res.GetJson()
        if err == nil {
            err_unmarshal:=json.Unmarshal(result_js,&line)
            if err_unmarshal == nil {
                fmt.Printf("result_line: %v\n",line)
                if len(line.Data)>=6 {
                    //for i := range line.Data {
                    //    fmt.Printf("i:%d string:%s\n",i,line.Data[i])
                    //}
                    request = line.Data[4]
                    status  = line.Data[5]
                }
                for i:= range line.Data {
                    entry := line.Data[i]
                    beauty_log+=entry+"\n"
                }
            }
        }
    }
    //
    return
}


func NewNginxLogHandler()(*LogHandler){
    //
    var log_handler LogHandler
    log_handler.handler = handling.NewHandler(nil)
    //
    var fl filtering.FilterList
    var sq_filter       = filtering.Filter{ Name:"square_brackets_filter", Call:filtering.SquareBracketsFilter, Enabled:true }
    var url_filter      = filtering.Filter{ Name:"url_filter",             Call:filtering.UrlFilter,            Enabled:true }
    var path_filter     = filtering.Filter{ Name:"path_filter",            Call:filtering.PathFilter,           Enabled:true }
    var quotes_filter   = filtering.Filter{ Name:"quotes_filter",          Call:filtering.QuotesFilter,         Enabled:true }
    var dot_filter      = filtering.Filter{ Name:"dot_filter",             Call:filtering.DotFilter,            Enabled:true }
    //
    fl.Append(sq_filter)
    fl.Append(url_filter)
    fl.Append(path_filter)
    fl.Append(quotes_filter)
    fl.Append(dot_filter)
    log_handler.handler.AddFilters(fl)
    //
    // target_config          := make(map[string]string, 0)
    // target_config["type"]  =  "SINGLE_LINE"
    // tgt,_                  := targets.NewTarget(target_config)
    // log_handler.handler.AddTargetPtr(tgt)
    // dangerous thing
    // log_handler.target = tgt
    //
    return &log_handler
    //
}
