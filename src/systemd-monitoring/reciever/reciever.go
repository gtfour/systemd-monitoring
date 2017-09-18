package reciever

import "github.com/gin-gonic/gin"
import "systemd-monitoring/common"

type Reciever struct {
    Address      string
    SecretPhrase string
    Updates      chan common.DataUpdate
    app          *gin.Engine
}

func NewReciever(address string,secret_phrase string,updates chan common.DataUpdate) ( *Reciever,error ) {

    var r Reciever
    r.app = gin.Default()
    r.app.POST("/updates/recieve", recieve_update(gin.H{"secret_phrase":secret_phrase}))



    return &r,nil
}

func(r *Reciever)Run()(error){



    return nil

}

func recieve_update(data  gin.H)(func (c *gin.Context)) {
    return  func( c *gin.Context ) {
        //
        // temporary handler just for fun :)
        //
        dashboardName   := c.PostForm("update")
        //
        if err == nil {
            c.JSON( 200, gin.H{"status": "ok", "data":wsInfoResponse} )
        } else {
            c.JSON( 500, gin.H{"status": "error"} )
        }
        //
        //
    }
}



