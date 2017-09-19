package reciever

import "fmt"
import "errors"
import "encoding/json"
import "github.com/gin-gonic/gin"
import "systemd-monitoring/common"
import "systemd-monitoring/crypted"

var updatesChanIsNil = errors.New("updates chan is nil")

type Reciever struct {
    Address      string
    secretPhrase string
    updates      chan common.DataUpdate
    app          *gin.Engine
}

func NewReciever(address string,secret_phrase string,updates chan common.DataUpdate) ( *Reciever,error ) {

    if updates == nil { return nil, updatesChanIsNil }
    var r Reciever
    r.updates      = updates
    r.Address      = address
    r.secretPhrase = secret_phrase
    r.app          = gin.Default()
    r.app.POST("/updates/recieve", r.recieveUpdate(gin.H{"secret_phrase":secret_phrase}))
    return &r,nil

}

func(r *Reciever)Run()(error){

    r.app.Run(r.Address)
    fmt.Printf("Exiting\n")
    return nil

}

func (r *Reciever)recieveUpdate(data  gin.H)(func (c *gin.Context)) {
    return  func( c *gin.Context ) {
        //
        var m common.Message
        c.BindJSON(&m)
        // checking data type
        if m.Type != common.TypeDataUpdate { c.JSON( 500, gin.H{"status": "wrong data type"}) }
        // decrypt message
        keyByte:=[]byte(r.secretPhrase)
        updateByte,err:=crypted.Decrypt(keyByte,m.Data)
        if err == nil {
            var update common.DataUpdate
            // parse decrypted json
            err_unmarshal := json.Unmarshal(updateByte, &update)
            if err_unmarshal == nil {
                r.updates <- update
                c.JSON( 200, gin.H{"status": "ok"})
            } else {
                c.JSON( 500, gin.H{"status": "encode json error"})
            }
        } else {
            c.JSON( 500, gin.H{"status": "encryption error"} )
        }
    }
}



