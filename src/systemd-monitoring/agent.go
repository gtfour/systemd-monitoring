package main

import "fmt"
import "flag"
import "time"
import "errors"
import "strings"
import "systemd-monitoring/files"
import "systemd-monitoring/common"
import "systemd-monitoring/capture"
import "systemd-monitoring/notifier"
import "systemd-monitoring/service"
import "systemd-monitoring/config"
import "systemd-monitoring/filter"


var parseError        = errors.New("cmd line parse error")
var masterAddrIsEmpty = errors.New("Master address is empty")
var nothingToDo       = errors.New("There are no any files to capture or services to monitor or task to monitor docker-events")

func main(){

    agentConfig, err := parseInput()
    if err!=nil { fmt.Printf("Exiting:err:%v\n",err) ; return }
    fmt.Printf("Monitors: %v\n",agentConfig.Monitors)
    updates      := make(chan common.DataUpdate)
    //
    var notifier notifier.NativeNotifier
    notifier.Address      = agentConfig.MasterAddress
    notifier.SecretPhrase = agentConfig.SecretPhrase
    //
    chain,err := service.NewServiceChain(agentConfig.ServiceList, updates)
    if err !=nil {
        fmt.Printf("Warning:%v\n",err)
    }

    //
    // follow simple file changes
    // every data-change will be marked as "file"
    //
    runner,err:= capture.NewRunner(agentConfig.FilesList, updates)
    if err !=nil {
        fmt.Printf("Warning:%v\n",err)
    }
    //
    // follow nginx log files
    // every data-change will be marked as "nginx-log"
    //
    for i:= range agentConfig.NginxLogs {
        nginxLogPath        := agentConfig.NginxLogs[i]
        nginxLogsTarget,err := capture.NewNginxLogTarget(nginxLogPath)
        if err == nil {
            runner.AppendTarget(nginxLogsTarget)
        }
    }
    //
    // follow to output of docker-events command
    // every data-change will be marked as "docker-events" 
    //
    if agentConfig.DockerEvents {
        dockerEventsTarget,err := capture.NewDockerEventsTarget()
        if err != nil { fmt.Printf("Unable to create docker target\nerr:%v",err) }
        err=runner.AppendTarget(dockerEventsTarget)
        if err != nil { fmt.Printf("Unable to append new target\nerr:%v",err) }
    }

    //
    // follow to file with python tracebacks
    // every data-change will be marked as "python-traceback"
    //
    for i:= range agentConfig.PythonTracebacks {
        ptPath       := agentConfig.PythonTracebacks[i].Path
        ptTarget,err := capture.NewPythonTracebackTarget(ptPath)
        if err == nil {
            runner.AppendTarget(ptTarget)
        }
    }
    //
    //
    //
    //

    go runner.Handle()
    go chain.Proceed()

    // define filter relay
    filteredUpdates   := make(chan common.DataUpdate)
    relay,_           := filter.NewRelay(updates, filteredUpdates)
    //
    for i := range agentConfig.Monitors {
        //
        m := agentConfig.Monitors[i]
        switch m.Type {
            case config.FILE_MONITOR_TYPE:
                fm,err := config.ParseFileMonitor(m.Monitor)
                if err == nil {
                    relay.FileMonitors=append(relay.FileMonitors, fm)
                }
            case config.NGINX_LOG_MONITOR_TYPE:
                nm,err := config.ParseNginxLogMonitor(m.Monitor)
                if err == nil {
                    relay.NginxLogMonitors = append(relay.NginxLogMonitors, nm)
                }
            default:
                fmt.Printf("Unrecognized type:\n")
        }
        //
    }
    if len(agentConfig.PythonTracebacks)>0 {
        var pythonTracebackHandlerSet filter.PythonTracebackHandlerSet
        for i := range agentConfig.PythonTracebacks {
            pTracebackHandlerConfig            := agentConfig.PythonTracebacks[i]
            destFilePath                       := pTracebackHandlerConfig.Path
            matchSectionContainingThisKeywords := pTracebackHandlerConfig.Keywords
            pythonTracebackHandlerSet.AppendHandler( destFilePath, matchSectionContainingThisKeywords )
        }
        relay.PythonTracebackHandlerSet = pythonTracebackHandlerSet
    }
    //
    go relay.Handle()
    //
    for {
        select {
            case u:=<-filteredUpdates:
                err = notifier.Notify(u)
            default:
                time.Sleep(time.Second * 2)
        }

    }
    // IsStringIn
}

func parseInput()(blank *common.AgentConfig,err error){
    //
    //
    var agentConfig       common.AgentConfig
    var secretPhrase      string
    var masterAddress     string
    var files_list        string
    var nginx_logs        string
    var python_tracebacks string
    var service_list      string
    var monitors          string
    var docker_events     bool
    //
    secretPhrasePtr     := flag.String("secret-phrase","","Phrase to crypt messages.The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.")
    masterAddressPtr    := flag.String("master-address"    ,"","Remote master server ip-address")
    //masterPortPtr     := flag.String("master-port"  ,"","Remote master server port")
    filesListPtr        := flag.String("files"        ,"","Log files that's will be captured")
    nginxLogsPtr        := flag.String("nginx-logs"        ,"","Nginx log files. Files will be handled in special filter chain")
    pythonTracebacksPtr := flag.String("python-tracebacks" ,"","Handle log-files containing python tracebacks in special filter chain")
    serviceListPtr      := flag.String("services"     ,"","Service whose are going to be monitored")
    dockerEventsPtr     := flag.Bool("docker-events"     ,false,"Capture output of 'docker events' command")
    monitorsPtr         := flag.String("monitors"     ,"","Capture file changes")
    flag.Parse()
    //
    //
    if secretPhrasePtr     != nil { secretPhrase      = *secretPhrasePtr     } else { err = parseError ; return }
    if masterAddressPtr    != nil { masterAddress     = *masterAddressPtr    } else { err = parseError ; return }
    //if masterPortPtr     != nil { masterPort        = *masterPortPtr       } else { err = parseError ; return }
    if filesListPtr        != nil { files_list        = *filesListPtr        } else { err = parseError ; return }
    if nginxLogsPtr        != nil { nginx_logs        = *nginxLogsPtr        } else { err = parseError ; return }
    if serviceListPtr      != nil { service_list      = *serviceListPtr      } else { err = parseError ; return }
    if pythonTracebacksPtr != nil { python_tracebacks = *pythonTracebacksPtr } else { err = parseError ; return }
    if dockerEventsPtr     != nil { docker_events     = *dockerEventsPtr     } else { err = parseError ; return }
    if monitorsPtr         != nil { monitors          = *monitorsPtr         } else { err = parseError ; return }
    // 
    filesList        := files.ParseFileList(files_list)
    nginxLogs        := files.ParseFileList(nginx_logs)
    //pythonTracebacks := files.ParseFileList(python_tracebacks)
    serviceList      := strings.Split(service_list," ")
    if len(serviceList) == 1 && serviceList[0] == "" {
        serviceList = []string{}
    }
    monitors_list,err_mon               := config.ParseMonitors(monitors)
    ptraceback_handler_config_list, err := config.ParsePythonTracebackHandlerConfig(python_tracebacks)
    if err !=nil && python_tracebacks != "" { return }
    if (len(filesList)<1)&&(len(serviceList)<1)&&(docker_events == false)&&(len(nginxLogs)<1)&&( monitors=="" )&&( python_tracebacks=="" ) {
        err = nothingToDo
        return
    }
    if masterAddress == "" { err =  masterAddrIsEmpty ; return  }
    err          = common.ValidateSecretPhrase(secretPhrase)
    if err !=nil { return }
    //
    agentConfig.SecretPhrase     = secretPhrase
    agentConfig.MasterAddress    = masterAddress
    agentConfig.FilesList        = filesList
    agentConfig.NginxLogs        = nginxLogs
    agentConfig.ServiceList      = serviceList
    agentConfig.PythonTracebacks = ptraceback_handler_config_list
    agentConfig.DockerEvents     = docker_events
    //
    if err_mon == nil {
        agentConfig.Monitors         = monitors_list
    }
    //
    return &agentConfig, nil
    //
    //
}

