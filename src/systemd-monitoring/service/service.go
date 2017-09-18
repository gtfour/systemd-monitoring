package service

import "os"
import "time"
import "os/exec"
import "strconv"
import "errors"
import "strings"

var cantGetMainPid               = errors.New("can't get main pid")
var unableToRetrieveServiceInfo  = errors.New("unable to get service information")
var serviceChainIsEmpty          = errors.New("service chain is empty")
var systemctl_path string        = "/bin/systemctl"
var dead_pid                 int = -1
var unrecognized_service_pid int = -3


type Service struct {
    //
    name       string
    state      string
    pid        int
    changed    bool
    isActive   bool
    isEnabled  bool
    //
}

type Chain struct {
    //
    hostname    string
    services    []*Service
    processing  chan *Service
    updates     chan string
    timeout_sec time.Duration
    //
}




func(s *Service)String() string {
    return "service:"+s.name+"  "+"active:"+strconv.FormatBool(s.isActive)+"  enabled:"+strconv.FormatBool(s.isEnabled)+"  pid:"+strconv.Itoa(s.pid)+"\t"
}


func CheckSystemdService(service_name string)(*Service,error){
    //
    var service Service
    //
    out_active_byte,_ := exec.Command(systemctl_path,"is-active",service_name).Output()
    out_enabled_byte,_:= exec.Command(systemctl_path,"is-enabled",service_name).Output()
    out_active  := string(out_active_byte)
    out_enabled := string(out_enabled_byte)
    out_active  = strings.Replace(out_active, "\n", "", 1)
    out_enabled = strings.Replace(out_enabled, "\n","", 1)
    //
    //if err_active != nil || err_enabled != nil {
    //    fmt.Printf("checkSystemdService:err_active:%v\terr_enabled:%v\t\n",err_active,err_enabled)
    //    return nil, unableToCheckService
    //}

    if out_active == "active"{
        service.isActive = true
    } else {
        service.isActive = false
    }
    //
    if out_enabled == "enabled"{
        service.isEnabled = true
    } else {
        service.isEnabled = false
    }
    pid,err := GetMainPid(service_name)
    if err == nil {
        service.pid = pid
    } else {
        return nil,err
    }
    return &service, nil
    //
}

func GetMainPid(service_name string)(main_pid int, err error){
    //
    main_pid          =  -1
    out_status_byte,_ := exec.Command(systemctl_path,"status",service_name).Output()
    out_status        := string(out_status_byte)
    status            := strings.Split(out_status, "\n")
    main_pid_line     := ""
    //
    is_dead := false
    for i := range status {
        line := status[i]
        if strings.HasPrefix(line," Main PID:") {
            main_pid_line = line
            break
        } else if strings.HasPrefix(line, "   Active: inactive (dead)"){
            main_pid = dead_pid
            is_dead  = true
            break
        }
    }
    if is_dead {
        return main_pid, nil
    }
    if main_pid_line == "" {
        return unrecognized_service_pid, cantGetMainPid
    }
    //
    main_pid_slice := strings.Split(main_pid_line," ")
    for z := range main_pid_slice {
        word         := main_pid_slice[z]
        intWord, err := strconv.Atoi(word)
        if err == nil {
            main_pid = intWord
            break
        }
    }
    //
    return
}

func NewServiceChain(service_names []string, updates chan string)(c *Chain,err error){
    //
    //
    var chain Chain
    for i:= range service_names {
        service_name := service_names[i]
        s,err        := CheckSystemdService(service_name)
        if err == nil {
             chain.services = append(chain.services, s)
        }
    }
    // serviceChainIsEmpty
    if len(chain.services) == 0 { return nil, serviceChainIsEmpty }
    hostname,err := os.Hostname()
    if err == nil {
        c.hostname = hostname
        return &chain, nil
    } else {
        return nil, err
    }
    chain.timeout_sec    = 5
    chain.updates    = make(chan string)
    chain.processing = make(chan *Service)
    return &chain,nil
    //
    //
}

func (c *Chain)Proceed()(){
    //
    if c == nil { return }
    for i := range c.services {
        s := c.services[i]
        if s != nil {
           name                    := s.name
           currentServiceState,err := CheckSystemdService(name)
           if currentServiceState != nil && err == nil {
               changes,err := Compare(s, currentServiceState)
               if err == nil && len(changes)>0 {
                   for i:= range changes {
                       c.updates <- changes[i]
                   }
               }
           }
        }
    }
    //
}


func (c *Chain)Proceed2()(){
    //
    //
    if c == nil { return }
    for  {
        select {
            case s := <-c.processing:
                if s != nil {
                    name                    := s.name
                    currentServiceState,err := CheckSystemdService(name)
                    if currentServiceState != nil && err == nil {
                        changes,err := Compare(s, currentServiceState)
                        if err == nil && len(changes)>0 {
                            for i:= range changes {
                                c.updates <- "hostname:"+c.hostname+":\n"+changes[i]+"\n"
                            }
                        }
                    }
                }
            default:
                time.Sleep(time.Second * c.timeout_sec)
        }
    }
    //
    //
}



func Compare(oldS *Service, newS *Service)( changes []string ,err error) {
    if oldS == nil || newS == nil {
        return changes, unableToRetrieveServiceInfo
    }
    if oldS.name != newS.name {
        return changes, unableToRetrieveServiceInfo
    }
    service_name := "service '"+newS.name+"' "
    if oldS.pid != newS.pid {
        if oldS.pid == dead_pid {
            changes = append(changes, service_name+"has been started")
        } else if newS.pid == dead_pid {
            changes = append(changes, service_name+"has been stoped")
        } else {
            changes = append(changes, service_name+"has been restarted")
        }
    }
    if oldS.isEnabled != newS.isEnabled {
        if newS.isEnabled {
            changes = append(changes, service_name+"has been enabled")
        } else {
            changes = append(changes, service_name+"has been disabled")
        }
    }
    return changes,nil
}
