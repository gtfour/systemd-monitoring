package common

import "os"

func FileExists(filepath string)(bool){
    if _, err := os.Stat(filepath); err == nil {
        return true
    } else {
        return false
    }
}

