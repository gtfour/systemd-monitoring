package main

import "fmt"
import "time"
import "flag"
import "errors"
import "gopkg.in/telegram-bot-api.v4"
import "systemd-monitoring/bot"
import "systemd-monitoring/user"
import "systemd-monitoring/common"
import "systemd-monitoring/reciever"

var allowedUsersListIsEmpty = errors.New("allowed users list is empty")
var parseError              = errors.New("cmd line parse error")
var listenAddrIsEmpty       = errors.New("Listen address is empty")

func main() {


    listenAddress,secretPhrase,token,users,err:=parseInput()
    if err!=nil {
        if err!=nil { fmt.Printf("Exiting:err:%v\n",err) ; return }
    }
    updates     := make(chan common.DataUpdate,100)
    updatesSend := make(chan common.DataUpdate,100)
    rec,err   := reciever.NewReciever(listenAddress,secretPhrase,updates)
    if err!=nil { fmt.Printf("Exiting:err:%v\n",err) ; return }

    runner,err := bot.NewRunner(token,users,updatesSend)
    if err != nil {
        fmt.Printf("error:\tCouldn't init runner: %v\n",err)
        return
    }
    go runner.Handle()
    go rec.Run()
    if err == nil {
        for {
            select {
                case u:=<-updates:
                   updatesSend<-u
                default:
                    time.Sleep(time.Second * 2)
            }
        }
    }





}





func parseInput()(listenAddress string,secretPhrase string,token string, users []tgbotapi.User, err error){
    //
    var allowed_users string
    //
    listenAddressPtr    := flag.String("listen"        ,"","Listen address")
    secretPhrasePtr     := flag.String("secret-phrase","","Phrase to crypt messages.The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.")
    tokenPtr            := flag.String("token"         ,"","Telegram bot token")
    allowedUsersListPtr := flag.String("allowed-users" ,"",`Allowed users.\nFormat: [{"id":"123456","first_name:"Ivan","last_name":"Ivanov"}]`)
    flag.Parse()
    //
    //
    if listenAddressPtr    != nil {  listenAddress  = *listenAddressPtr    } else { err = parseError ; return }
    if secretPhrasePtr     != nil { secretPhrase  = *secretPhrasePtr       } else { err = parseError ; return }
    if tokenPtr            != nil {  token          = *tokenPtr            } else { err = parseError ; return }
    if allowedUsersListPtr != nil {  allowed_users  = *allowedUsersListPtr } else { err = parseError ; return }
    //
    if allowed_users == "" { err = allowedUsersListIsEmpty ; return }
    users,err = user.ParseUsers(allowed_users)
    if err != nil {
        return
    }
    //
    if listenAddress == "" { err = listenAddrIsEmpty ; return  }
    err              = common.ValidateSecretPhrase(secretPhrase)
    if err!=nil {return}

    //
    return
    //
}

