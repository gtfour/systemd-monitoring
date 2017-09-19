package bot

//
// use non-root account to run this daemon
// 

import "time"
import "errors"
import "gopkg.in/telegram-bot-api.v4"
//
import "systemd-monitoring/user"
import "systemd-monitoring/common"
//
var response_status_success           int = 10
var response_status_failed            int = 11
var response_status_permission_denied int = 15
var call_doesnt_exist                 int = 21

var tokenIsEmpty            = errors.New("token is empty")
var serviceListIsEmpty      = errors.New("service list is empty")
var allowedUsersListIsEmpty = errors.New("allowed users list is empty")
var unableToParseUsersList  = errors.New("unable to parse users list")
var userIsNotAllowed        = errors.New("this user is not allowed to run commands")
var updateIsNil             = errors.New("update is nil")
//


type Runner struct {
    token                          string
    bot                            *tgbotapi.BotAPI
    users                          user.Users
    updates                        <-chan tgbotapi.Update
    dataUpdates                    chan common.DataUpdate
    quitHandle                     chan bool
    quit                           chan bool
    timeout_sec                    time.Duration
}

func(r *Runner)Handle()(error){
    for {
        select {
            case update:=<-r.updates:
                message := update.Message
                if message == nil {
                    continue
                }
                //go r.handleMessage(message)
            case update:=<-r.dataUpdates:
                //fmt.Printf("<< response has been recieved >>\n")
                all_users := r.users.AllUsersIds()
                for i:= range all_users {
                    id  := all_users[i]
                    msg := tgbotapi.NewMessage(int64(id), update.String())
                    //msg.ReplyToMessageID =  response.message_id
                    r.bot.Send(msg)
                }
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


func NewRunner( token string, users []tgbotapi.User, dataUpdates chan common.DataUpdate )( *Runner , error){
    var r Runner
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        return nil, err
    }
    r.token   = token
    bot.Debug = true
    r.bot     = bot
    //fmt.Printf("Authorized on account %s",r.bot.Self.UserName)
    u         := tgbotapi.NewUpdate(0)
    u.Timeout =  60

    r.updates, err = r.bot.GetUpdatesChan(u)
    if err != nil {
        return nil,err
    }
    //
    //
    r.quitHandle           = make(chan bool)
    r.quit                 = make(chan bool)
    r.dataUpdates          = dataUpdates

    r.users                = users
    r.timeout_sec = 2
    return &r, nil
}
