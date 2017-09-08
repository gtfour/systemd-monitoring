package capture

import "io"
import "fmt"
import "time"
import "os/exec"
import "errors"
import "systemd-monitoring/common"

var fileNotExists        = errors.New("File doesn't exist")
var unableToCreateTarget = errors.New("Unable to create target")
var tail_path     string = "/usr/bin/tail"

type Update struct {
    path   string
    update string
}

type Target struct {
    cmd                *exec.Cmd
    stdout             io.ReadCloser
    quit               chan bool
    updates            chan Update
}

type Runner struct {
    updates     chan Update
    quit        chan bool
    main_quit   chan bool
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
    return t,nil

}

func NewRunner(pathes []string)( *Runner , error){
    //
    var r Runner
    r.updates   =  make(chan Update)
    r.quit      = make(chan bool)
    r.main_quit = make(chan bool)
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
