package main

//
// use non-root account to run this daemon
// 

import "os"
import "fmt"
import "flag"
import "time"
import "os/signal"
import "errors"
import "strings"
import "gopkg.in/telegram-bot-api.v4"
//
import "systemd-monitoring/user"
import "systemd-monitoring/service"
import "systemd-monitoring/common"
//
var response_status_success           int = 10
var response_status_failed            int = 11
var response_status_permission_denied int = 15
var call_doesnt_exist                 int = 21

var parseError              = errors.New("cmd line parse error")
var tokenIsEmpty            = errors.New("token is empty")
var serviceListIsEmpty      = errors.New("service list is empty")
var allowedUsersListIsEmpty = errors.New("allowed users list is empty")
var unableToParseUsersList  = errors.New("unable to parse users list")
var userIsNotAllowed        = errors.New("this user is not allowed to run commands")
var updateIsNil             = errors.New("update is nil")
var unableToCheckService    = errors.New("unable to check service")
//

type Response struct {
    chat_id    int64
    message_id int
    hostname   string
    text       string
    status     int
}

func(r *Response)String()string{
    if r.status == response_status_permission_denied {
        return "permissions denied"
    }
    return "hostname:"+r.hostname+"\n########\n"+r.text+"\n########\nstatus:"+string(r.status)
}


type Message struct {
    //
    chat_id        int64 // user's id who requests service information
    hostname       string
    date           string
    service_names  []string
    function_names []string
    response       string
    //
}

type Runner struct {
    //
    token                          string
    systemctl_path                 string
    service_names                  string
    downServiceNotificationPeriod  int
    bot                            *tgbotapi.BotAPI
    responses                      chan *Response
    users                          user.Users
    updates                        <-chan tgbotapi.Update
    // serviceCheckRequests           chan Message
    // serviceChecksResults           chan Message
    services                       []service.Service
    servicesChain                  chan *service.Service
    fileChecks                     chan string
    quitHandle                     chan bool
    quit                           chan bool
    timeout_sec                    time.Duration
    //
}




func main() {
    token,systemctl_path,service_names,users,err := parseInput()
    //fmt.Printf("token:%v\tsystemctl_path:%v\nservices:%v\tusers:%v\nerr:%v",token,systemctl_path,services,users,err)
    if err != nil {
        fmt.Printf("error:\tCouldn't parse command line arguments: %v\n",err)
        return
    }
    runner,err := NewRunner( token ,systemctl_path ,service_names ,users )
    if err != nil {
        fmt.Printf("error:\tCouldn't init runner: %v\n",err)
        return
    }
    runner.run()
    //
}


func parseInput()(token string,systemctl_path string,services []string, users []tgbotapi.User, err error){
    //
    var service_list  string
    var allowed_users string
    //
    tokenPtr            := flag.String("token"         ,"","Telegram bot token")
    systemctlPathPtr    := flag.String("systemctl-path","","Path to systemctl binary")
    serviceListPtr      := flag.String("services","","Service whose are going to be monitored")
    allowedUsersListPtr := flag.String("allowed-users","",`Users whose able to run commands.\nFormat: [{"id":"123456","first_name:"Ivan","last_name":"Ivanov"}]`)
    flag.Parse()
    //
    //
    if tokenPtr            != nil {  token          = *tokenPtr            } else { err = parseError ; return }
    if systemctlPathPtr    != nil {  systemctl_path = *systemctlPathPtr    } else { err = parseError ; return }
    if serviceListPtr      != nil {  service_list   = *serviceListPtr      } else { err = parseError ; return }
    if allowedUsersListPtr != nil {  allowed_users  = *allowedUsersListPtr } else { err = parseError ; return }
    //
    if service_list  == "" { err = serviceListIsEmpty      ; return }
    if allowed_users == "" { err = allowedUsersListIsEmpty ; return }
    fmt.Printf("\nUsers: %v\n",allowed_users)
    //
    services  = strings.Split(service_list," ")
    users,err = user.ParseUsers(allowed_users)
    //
    return
    //
}

func(r *Runner)run()(error){

    go r.handle()
    r.catchExit()
    return nil

}

func(r *Runner)catchExit()(){

    signalChan  := make(chan os.Signal, 1)
    cleanupDone := make(chan bool)
    signal.Notify(signalChan, os.Interrupt)
    signal.Notify(signalChan, os.Kill)
    go func() {
        for _ = range signalChan {
            r.quitHandle <- true
            <-r.quit
            cleanupDone <- true
            break
        }
    }()
    <-cleanupDone
    return

}


func(r *Runner)handle()(error){
    //
    //
    /*for update := range r.updates {

        if update.Message == nil {
            continue
        }

        fmt.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
        msg.ReplyToMessageID = update.Message.MessageID

        r.bot.Send(msg)

    }*/
    //
    //
    for {
        select {
            case update:=<-r.updates:
                message := update.Message
                if message == nil {
                    continue
                }
                go r.handleMessage(message)
            case response:=<-r.responses:
                fmt.Printf("<< response has been recieved >>\n")
                if response == nil {
                    continue
                }
                msg                  := tgbotapi.NewMessage(response.chat_id, response.String())
                msg.ReplyToMessageID =  response.message_id
                r.bot.Send(msg)
            case <-r.quitHandle:
                break
            default:
                time.Sleep(time.Second * r.timeout_sec)

        }
    }
    //
    r.quit<-true
    return nil
}

func(r *Runner)handleMessage(message *tgbotapi.Message)(){
    //
    fmt.Printf("===HandleMessage===\n")
    var resp Response
    resp.chat_id    = message.Chat.ID
    resp.message_id = message.MessageID
    if !r.users.UserIsAllowed(message) {
        resp.status = response_status_permission_denied
        r.responses     <- &resp
        return
    }
    id          := message.From.ID
    first_name  := message.From.FirstName
    text        := message.Text
    service,err := service.CheckSystemdService(text)
    fmt.Printf("service_name:%v\tservice:%v\terr:%v\n",text,service,err)
    if err == nil {
        resp.text = service.String()
    } else {
        resp.text = ""
    }
    fmt.Printf("func:handleUpdates:from:\tid:%v\tfirst_name:%v\ntext:%v\n",id,first_name,text)
    //
    resp.hostname   = "zombie"
    resp.status     = response_status_success
    r.responses     <- &resp
    return
    //
    //
}


func(r *Runner)checkServices()(){
    servicesLen := len(r.services)
    _           =  servicesLen
    for {


    }
}

func(r *Runner)initServices()(){


}


func NewRunner( token string, systemctl_path string , service_names []string, users []tgbotapi.User )( *Runner , error){
    var r Runner
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        return nil, err
    }
    r.token   = token
    bot.Debug = true
    r.bot     = bot
    fmt.Printf("Authorized on account %s",r.bot.Self.UserName)
    u         := tgbotapi.NewUpdate(0)
    u.Timeout =  60

    r.updates, err = r.bot.GetUpdatesChan(u)
    if err != nil {
        return nil,err
    }
    //
    //
    if systemctl_path == "" {
        r.systemctl_path = "/bin/systemctl"
    } else {
        if common.FileExists(systemctl_path) {
            r.systemctl_path = systemctl_path
        }
    }
    r.quitHandle           = make(chan bool)
    r.quit                 = make(chan bool)
    r.responses            = make(chan *Response)
    r.servicesChain        = make(chan *service.Service)
    // r.serviceCheckRequests = make(chan Message)
    // r.serviceChecksResults = make(chan Message)

    r.users                = users
    //
    // r.ch                = make(chan string,100)
    // r.quitCapture       = make(chan bool)
    // r.quitHandle        = make(chan bool)
    // r.quit              = make(chan bool)
    //
    // r.timeout_sec       = 2
    // fmt.Printf("runner:\n")
    // fmt.Printf("\n\tcmd_line:%v",cmd_line)
    // fmt.Printf("\n")
    //
    r.timeout_sec = 2
    return &r, nil
}
