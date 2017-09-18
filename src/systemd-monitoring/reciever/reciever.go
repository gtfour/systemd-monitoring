package reciever

import "fmt"
import "errors"
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
        // temporary handler just for fun :)
        //
        var m common.Message
        //messageType := c.PostForm("type")
        //messageData := c.PostForm("data")
        c.BindJSON(&m)
        fmt.Printf("recieveUpdate:: Type:%v Data:%v\nsecret_phrase:%v\n",m.Type,m.Data,data["secret_phrase"])
        //
        // decrypt message
        //
        keyByte:=[]byte(r.secretPhrase)
        update,err:=crypted.Decrypt(keyByte,m.Data)
        if err == nil {
            r.updates <- common.DataUpdate{"zombie",string(update),common.GetTime()}
        }
        //
        //if err == nil {
        c.JSON( 200, gin.H{"status": "ok"} )
        //} else {
        //    c.JSON( 500, gin.H{"status": "error"} )
        //}
        //
        //
    }
}



