package notifier

import "fmt"
import "bytes"
import "net/http"
import "io/ioutil"
import "encoding/json"
import "systemd-monitoring/common"
import "systemd-monitoring/crypted"


type Notifier interface {
    Notify(u common.DataUpdate)(error)
}

type NativeNotifier struct {
    Address      string
    SecretPhrase string
}

type SysLogNotifier struct {
    syslogServerAddr string
}

func(n *NativeNotifier)Notify(u common.DataUpdate)(error){
    //
    // TypeDataUpdate 
    //
    url := "http://"+n.Address+"/updates/recieve"
    secret_phrase   := n.SecretPhrase
    updateByte,err  := json.Marshal(u)
    if err != nil { return err }
    keyByte         := []byte( secret_phrase )
    cryptedByte,err := crypted.Encrypt(keyByte,updateByte)
    if err != nil { return err }
    var m common.Message
    m.Type = common.TypeDataUpdate
    m.Data = cryptedByte
    messageByte,err := json.Marshal(m)
    if err != nil { return err }
    //
    // sending json 
    //
    req, err  := http.NewRequest("POST", url, bytes.NewBuffer(messageByte))
    req.Header.Set("Content-Type", "application/json")
    client    := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    //
    // print response
    //
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
    //
    //
    //
    return nil
}

func(s *SysLogNotifier)Notify(u common.DataUpdate)(error){
    return nil
}

