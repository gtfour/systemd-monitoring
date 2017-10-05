package filter
//
import "fmt"
import "time"
import "systemd-monitoring/common"
//
type Relay struct {
    //
    updatesInput  chan common.DataUpdate
    updatesOutput chan common.DataUpdate
    quit          chan bool
    timeoutSec    time.Duration
    //
}

type Matcher struct {
    //

    //
}

func(r *Relay)Handle()(){

    for {
        select {
            case u:=<-r.updatesInput:
                fmt.Printf("Relay:updates:\n%v", u)
                r.updatesOutput<-u
            case <-r.quit:
                break
            default:
                time.Sleep(time.Second * r.timeoutSec)
        }
    }
}

//
func NewRelay(config common.AgentConfig)(*Relay,error){
    //
    var relay Relay
    relay.updatesInput  = make(chan common.DataUpdate)
    relay.updatesOutput = make(chan common.DataUpdate)
    relay.quit          = make(chan bool)
    //
    return &relay, nil
}
//
