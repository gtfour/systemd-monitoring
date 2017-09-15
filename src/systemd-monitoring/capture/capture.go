package capture

import "io"
import "fmt"
import "time"
import "bufio"
import "os/exec"
import "errors"
import "systemd-monitoring/common"

var fileNotExists        = errors.New("File doesn't exist")
var unableToCreateTarget = errors.New("Unable to create target")
var tail_path     string = "/usr/bin/tail"

type Update struct {
    path string
    text string
}

type Target struct {
    path     string
    cmd      *exec.Cmd
    stdout   io.ReadCloser
    quit     chan bool
    updates  chan Update
    active   bool
}

type Runner struct {
    updates     chan Update
    quit        chan bool
    MainQuit   chan bool
    QuitDone    chan bool
    targets     []*Target
    timeout_sec time.Duration
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

    if !common.FileExists(path) { return nil, fileNotExists }
    var command = []string{tail_path,"-f",path}
    t,err := NewTarget(command)
    if err!=nil { return nil, unableToCreateTarget }
    t.path = path
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
    //
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
                    lineStr := string(line)
                    var update Update
                    update.path = t.path
                    update.text = deffered+lineStr
                    t.updates<-update
                    deffered = ""
                }
                if err!= nil { break }
            case <- t.quit:
                exit = true
        }
    }
    fmt.Printf("capture has been finished for :%v\n",t.path)
    t.active = false
    t.cmd.Process.Kill()
    //
}


func NewRunner(pathes []string)( *Runner , error){
    //
    var r Runner
    r.updates   = make(chan Update)
    r.quit      = make(chan bool)
    r.MainQuit  = make(chan bool)
    r.QuitDone  = make(chan bool)
    for i:= range pathes {
        path:=pathes[i]
        t,err := NewTailTarget(path)
        if err!=nil {continue}
        t.updates = r.updates
        t.quit    = r.quit
        r.targets  = append( r.targets, t )
    }
    //
    r.timeout_sec       = 2
    fmt.Printf("runner:\n%v\n",r)
    return &r, nil
}

func (r *Runner)Handle()(){
    //
    for i:= range r.targets {
        target := r.targets[i]
        target.run()
        go target.capture()
    }
    finish := false
    //
    for {
        select {
            case u, ok := <-r.updates:
                    if !ok {
                        break
                    }
                    fmt.Println(u)
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


