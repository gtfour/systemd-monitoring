package reciever

import "fmt"
import "errors"
import "github.com/gin-gonic/gin"
import "systemd-monitoring/common"

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
    r.app.POST("/updates/recieve", r.recieve_update(gin.H{"secret_phrase":secret_phrase}))



    return &r,nil
}

func(r *Reciever)Run()(error){


    r.app.Run(r.Address)
    fmt.Printf("Exiting\n")
    return nil

}

func (r *Reciever)recieve_update(data  gin.H)(func (c *gin.Context)) {
    return  func( c *gin.Context ) {
        //
        // temporary handler just for fun :)
        //
        messageType := c.PostForm("type")
        messageData := c.PostForm("data")
        fmt.Printf("recieveUpdate:: Type:%v Data:%v\nsecret_phrase:%v\n",messageType,messageData,data["secret_phrase"])
        r.updates <- common.DataUpdate{"zombie",string(messageData),common.GetTime()}
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



