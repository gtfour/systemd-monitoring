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

var parseError        = errors.New("cmd line parse error")
var masterAddrIsEmpty = errors.New("Master address is empty")
var nothingToDo       = errors.New("There are no any files to capture or services to monitor or task to monitor docker-events")

func main(){

    config, err := parseInput()
    if err!=nil { fmt.Printf("Exiting:err:%v\n",err) ; return }
    updates      := make(chan common.DataUpdate)
    //
    var notifier notifier.NativeNotifier
    notifier.Address      = config.MasterAddress
    notifier.SecretPhrase = config.SecretPhrase
    //
    chain,err := service.NewServiceChain(config.ServiceList, updates)
    if err !=nil {
        fmt.Printf("Warning:%v\n",err)
    }

    runner,err:= capture.NewRunner(config.FilesList, updates)
    if err !=nil {
        fmt.Printf("Warning:%v\n",err)
    }
    for i:= range config.NginxLogs {
        nginxLogPath        := config.NginxLogs[i]
        nginxLogsTarget,err := capture.NewNginxLogTarget(nginxLogPath)
        if err == nil {
            runner.AppendTarget(nginxLogsTarget)
        }
    }
    if config.DockerEvents {
        dockerEventsTarget,err := capture.NewDockerEventsTarget()
        if err != nil { fmt.Printf("Unable to create docker target\nerr:%v",err) }
        err=runner.AppendTarget(dockerEventsTarget)
        if err != nil { fmt.Printf("Unable to append new target\nerr:%v",err) }
    }

    go runner.Handle()
    go chain.Proceed()

    for {
        select {
            case u:=<-updates:
                err = notifier.Notify(u)
                fmt.Printf("Update:%v\nnotifier_err:%v\n",u,err)
            default:
                time.Sleep(time.Second * 2)
        }

    }
    // IsStringIn
}

func parseInput()(blank *common.AgentConfig,err error){
    //
    var agentConfig       common.AgentConfig
    var secretPhrase      string
    var masterAddress     string
    var files_list        string
    var nginx_logs        string
    var python_tracebacks string
    var service_list      string
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
    // 
    filesList        := files.ParseFileList(files_list)
    nginxLogs        := files.ParseFileList(nginx_logs)
    pythonTracebacks := files.ParseFileList(python_tracebacks)
    serviceList      := strings.Split(service_list," ")
    if len(serviceList) == 1 && serviceList[0] == "" {
        serviceList = []string{}
    }
    if (len(filesList)<1)&&(len(serviceList)<1)&&(docker_events == false)&&(len(nginxLogs)<1) {
        err = nothingToDo
        return
    }
    if masterAddress == "" { err =  masterAddrIsEmpty ; return  }
    err          = common.ValidateSecretPhrase(secretPhrase)
    if err!=nil {return}
    //
    agentConfig.SecretPhrase     = secretPhrase
    agentConfig.MasterAddress    = masterAddress
    agentConfig.FilesList        = filesList
    agentConfig.NginxLogs        = nginxLogs
    agentConfig.ServiceList      = serviceList
    agentConfig.PythonTracebacks = pythonTracebacks
    agentConfig.DockerEvents     = docker_events
    //
    return &agentConfig, nil
    //
}

