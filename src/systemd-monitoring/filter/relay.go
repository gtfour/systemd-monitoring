package filter
//
import "fmt"
import "time"
import "errors"
import "strings"
import "systemd-monitoring/common"
import "systemd-monitoring/config"

var chanIsNil           = errors.New("Chan is nil")
//
type Relay struct {

    updatesInput      chan common.DataUpdate
    updatesOutput     chan common.DataUpdate
    quit              chan bool
    timeoutSec        time.Duration
    //
    NginxLogMonitors          []config.NginxLogMonitor
    FileMonitors              []config.FileMonitor
    PythonTracebackHandlerSet PythonTracebackHandlerSet
    //
    //
    nginxLogHandler           *LogHandler

}


func(r *Relay)Handle()(){
    for {
        select {
            case u:=<-r.updatesInput:
                fmt.Printf("Relay:updates:\n%v", u)
                //r.updatesOutput<-u
                r.passThroughMonitors(u)
            case <-r.quit:
                break
            default:
                time.Sleep(time.Second * r.timeoutSec)
        }
    }
}

func(r *Relay)passThroughMonitors(du common.DataUpdate)(){

    path := du.Path
    text := du.Text
    switch du.Area {
        case "file":
            skip := false
            //send := false
            for a := range r.FileMonitors {
                fm := r.FileMonitors[a]
                if path == fm.Path {
                    for aa := range fm.IgnoreString {
                        stringToIgnore := fm.IgnoreString[aa]
                        str_index      := strings.Index(text, stringToIgnore)
                        if str_index >= 0 {
                            skip = true
                            break
                        }
                    }
                }
                if skip == true { break } else { r.updatesOutput <- du ; break }
            }
        case "nginx-log":
            for b := range r.NginxLogMonitors {
                nm := r.NginxLogMonitors[b]
                if path == nm.Path {
                    _,status,beauty_message := r.nginxLogHandler.HandleNginxLog(text)
                    statusIsMatched         := common.IsStringIn(status, nm.MatchStatus)
                    if statusIsMatched {
                        du.Text = beauty_message
                        r.updatesOutput <- du
                        break
                    } else {
                        // then ignore message
                    }
                }
            }
        case "python-traceback":
            //fmt.Printf("Handling python tracebacks\n")
            duOut,err:=r.PythonTracebackHandlerSet.Processing(du)
            if err == nil {
                r.updatesOutput<-duOut
            }
        default:
            fmt.Printf("Unrecognized area type:\n")
    }
}

//
func NewRelay(updatesInput chan common.DataUpdate, updatesOutput chan common.DataUpdate)(*Relay,error){
    //
    if updatesInput == nil || updatesOutput == nil  {
        return nil,chanIsNil
    }
    // 
    var relay Relay
    relay.updatesInput    = updatesInput
    relay.updatesOutput   = updatesOutput
    relay.nginxLogHandler = NewNginxLogHandler()
    relay.quit            = make(chan bool)
    relay.timeoutSec      = 2
    //
    return &relay, nil
}
//
