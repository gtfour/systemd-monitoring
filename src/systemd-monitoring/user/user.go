package user

import "fmt"
import "encoding/json"
import "gopkg.in/telegram-bot-api.v4"

type Users []tgbotapi.User

func (users Users)UserIsAllowed(message *tgbotapi.Message)(found bool){
    //
    //
    fmt.Printf("\n::Users::\n%v\n",users)
    id         := message.From.ID
    first_name := message.From.FirstName
    found =  false
    //
    //
    for i := range users {
        user        := users[i]
        idI         := user.ID
        first_nameI := user.FirstName
        if idI == id && first_nameI == first_name {
            found = true
            break
        }
    }
    return
    //
    //
}





func ParseUsers(users_string string)(users []tgbotapi.User,err error){
    err = json.Unmarshal([]byte(users_string), &users)
    return
}

