package main

import "flag"
import "errors"
import "strings"
import "systemd-monitoring/files"

var parseError = errors.New("cmd line parse error")

func main(){


}

func parseInput()(secretPhrase string,masterIp string,masterPort string, filesList []string, serviceList []string , err error){
    //
    var files_list    string
    var service_list  string
    //
    secretPhrasePtr  := flag.String("secret-phrase","","Phrase to crypt messages.The key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.")
    masterIpPtr      := flag.String("master-ip"    ,"","Remote master server ip-address")
    masterPortPtr    := flag.String("master-port"  ,"","Remote master server port")
    filesListPtr     := flag.String("files"        ,"","Log files that's will be captured")
    serviceListPtr   := flag.String("services"     ,"","Service whose are going to be monitored")
    flag.Parse()
    //
    //
    if secretPhrasePtr != nil { secretPhrase = *secretPhrasePtr } else { err = parseError ; return }
    if masterIpPtr     != nil { masterIp     = *masterIpPtr     } else { err = parseError ; return }
    if masterPortPtr   != nil { masterPort   = *masterPortPtr   } else { err = parseError ; return }
    if filesListPtr    != nil { files_list   = *filesListPtr    } else { err = parseError ; return }
    if serviceListPtr  != nil { service_list = *serviceListPtr  } else { err = parseError ; return }
    // 
    filesList    = files.ParseFileList(files_list)
    serviceList  = strings.Split(service_list," ")
    return
    //
}

