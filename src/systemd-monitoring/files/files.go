package files

import "strings"
import "systemd-monitoring/common"

func ParseFileList(filesList string)(files []string){

    filesSlice := strings.Split(filesList," ")
    for i:=range filesSlice {
        filename := filesSlice[i]
        if common.FileExists(filename){
            files = append(files,filename)
        }
    }
    return files
}
