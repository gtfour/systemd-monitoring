package service

import "fmt"
import "strconv"
import "errors"
import "os/exec"
import "strings"

var cantGetMainPid        = errors.New("can't get main pid")
var systemctl_path string = "/bin/systemctl"


type Service struct {
    //
    name       string
    state      string
    pid        string
    changed    bool
    isActive   bool
    isEnabled  bool
    //
}

func(s *Service)String() string {
    return "service:"+s.name+"  "+"active:"+strconv.FormatBool(s.isActive)+"  "+"enabled:"+strconv.FormatBool(s.isEnabled)
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
    //
    fmt.Printf("=== === === service:%v:active:%v\tenabled:%v\n",service_name,out_active,out_enabled)
    //
    //
    return &service, nil
    //
    //
}

func GetServiceMainPid(service_name string)(main_pid int, err error){
    main_pid          = -1
    out_status_byte,_ := exec.Command(systemctl_path,"status",service_name).Output()
    out_status        := string(out_status_byte)
    status            := strings.Split(out_status, "\n")
    main_pid_line     := ""
    for i:= range status {
        line:=status[i]
        if strings.HasPrefix(line," Main PID:") {
            main_pid_line = line
            break
        }
    }
    if main_pid_line == "" {
        return -1, cantGetMainPid
    }
    main_pid_slice := strings.Split(main_pid_line," ")
    for z := range main_pid_slice {
        word         := main_pid_slice[z]
        intWord, err := strconv.Atoi(word)
        if err == nil {
            main_pid = intWord
            break
        }
    }
    return
}

