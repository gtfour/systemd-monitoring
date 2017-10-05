package capture

import "io"
import "os"
import "fmt"
import "time"
import "bufio"
import "os/exec"
import "errors"
import "systemd-monitoring/common"
import "systemd-monitoring/filter"

var fileNotExists        = errors.New("File doesn't exist")
var unableToCreateTarget = errors.New("Unable to create target")
var listOfTargetsIsNil   = errors.New("List of targets is nil")
var targetIsNil          = errors.New("New target is nil")
var updatesChanIsNil     = errors.New("Updates chan is nil")
var tail_path     string = "/usr/bin/tail"
var docker_path   string = "/usr/bin/docker"

type Update struct {
    path string
    text string
}

type Target struct {
    path     string
    cmd      *exec.Cmd
    stdout   io.ReadCloser
    quit     chan bool
    updates  chan common.DataUpdate
    active   bool
    area     string
}

type Runner struct {
    // updates     chan Update
    updates       chan common.DataUpdate
    globalUpdates chan common.DataUpdate
    quit          chan bool
    MainQuit      chan bool
    QuitDone      chan bool
    targets       []*Target
    timeout_sec   time.Duration
    running       bool
    logHandler    *filter.LogHandler
}

func NewTarget(cmd_line []string)(*Target,error){
    var t Target
    cmd,err := common.Command(cmd_line)
    if err!= nil {
        return nil,err
    }
    t.cmd = cmd
    return &t,err
}

func NewTailTarget(path string)(*Target,error){
    //
    if !common.FileExists(path) { return nil, fileNotExists }
    var command = []string{tail_path,"-F",path}
    t,err := NewTarget(command)
    if err!=nil { return nil, unableToCreateTarget }
    t.path = path
    t.area = "file"
    return t,nil
    //
}

func NewNginxLogTarget(path string)(*Target,error){
    //
    if !common.FileExists(path) { return nil, fileNotExists }
    var command = []string{tail_path,"-F",path}
    t,err := NewTarget(command)
    if err!=nil { return nil, unableToCreateTarget }
    t.path = path
    t.area = "nginx-log"
    return t,nil
    //
}

func NewDockerEventsTarget()(*Target,error){

    var command = []string{docker_path,"events"}
    t,err := NewTarget(command)
    if err!=nil { return nil, unableToCreateTarget }
    t.area = "docker_events"
    return t,nil

}


func (t *Target)run()(error){

    stdout, err := t.cmd.StdoutPipe()
    if err != nil { return err }
    t.stdout = stdout
    err = t.cmd.Start()
    if err != nil { return err }
    t.active = true
    return nil

}

func(t *Target)capture()(){
    fmt.Printf("\n<<< Capturing changes in area:%v  path:%v\n",t.area,t.path)
    exit       := false
    lineReader := bufio.NewReader(t.stdout)
    var deffered string
    for {
        select {
            default:
                if exit { break }
                line,isPrefix,err := lineReader.ReadLine()
                if isPrefix && err==nil {
                    deffered+=string(line)
                    continue
                }
                if err == nil && !isPrefix {
                    //
                    lineStr := string(line)
                    var update common.DataUpdate
                    hostname,err     := os.Hostname()
                    if err!=nil { hostname = "undefined" }
                    update.Hostname  = hostname
                    update.Area      = t.area
                    update.Text      = deffered+lineStr
                    update.Timestamp = common.GetTime()
                    t.updates<-update
                    deffered = ""
                    //
                }
                if err!= nil { break }
            case <- t.quit:
                exit = true
        }
    }
    fmt.Printf("capture has been finished for :%v\n",t.path)
    t.active = false
    t.cmd.Process.Kill()
}


func NewRunner(pathes []string, globalUpdates chan common.DataUpdate)( *Runner , error){
    //
    if globalUpdates == nil  { return nil, updatesChanIsNil }
    var r Runner
    r.updates       = make(chan common.DataUpdate)
    r.quit          = make(chan bool)
    r.MainQuit      = make(chan bool)
    r.QuitDone      = make(chan bool)
    r.globalUpdates = globalUpdates
    for i:= range pathes {
        path:=pathes[i]
        t,err := NewTailTarget(path)
        if err!=nil {continue}
        t.updates = r.updates
        t.quit    = r.quit
        r.targets  = append( r.targets, t )
    }
    //
    r.timeout_sec = 2
    r.logHandler  = filter.NewNginxLogHandler()
    // fmt.Printf("runner:\n%v\n",r)
    return &r, nil
}

func (r *Runner)Handle()(){
    //
    for i:= range r.targets {
        target := r.targets[i]
        err:=target.run()
        if err == nil {
            go target.capture()
        }
    }
    r.running = true
    finish := false
    //
    for {
        select {
            case u, ok := <-r.updates:
                    if !ok {
                        break
                    }
                    fmt.Println(u)
                    if u.Area == "nginx-log" {
                        r.handleNginxLogs(u,r.globalUpdates)
                    }else {
                        r.globalUpdates <- u
                    }
            case <-r.MainQuit:
                for i:= range r.targets {
                    target := r.targets[i]
                    if target.active {
                        target.quit <- true
                    }
                }
                r.QuitDone <- true
                finish = true
            default:
                if finish { break }
                time.Sleep(time.Second * r.timeout_sec)
        }
    }
}

func (r *Runner)AppendTarget(t *Target)(error){
    if t == nil { return targetIsNil  }
    t.updates = r.updates
    t.quit    = r.quit
    if r.running {
        t.run()
        go t.capture()
    }
    if r.targets==nil { r.targets = make([]*Target,0) }
    r.targets  = append(r.targets, t)
    return nil
}

func(r *Runner)handleNginxLogs(u common.DataUpdate, globalUpdates chan common.DataUpdate)(){
    _,status,beauty_message := r.logHandler.Handle(u.Text)
    fmt.Printf("\nBeauty Message:\n%v\n",beauty_message)
    if status == "500" || status  == "502" {
        u.Text = beauty_message
        r.globalUpdates <- u
    }
}
