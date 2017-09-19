package common

import "os"
import "time"
import "os/exec"
import "errors"
import "path/filepath"

var incorrectPhrase = errors.New("Length of secret phrase is incorrect.\nThe key argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.")

func FileExists(filepath string)(bool){
    //info.IsDir()
    info,err := os.Stat(filepath)
    if err == nil && info.IsDir() == false  {
        return true
    } else {
        return false
    }
}

func Command(args []string) (cmd *exec.Cmd,err error) {
    // overwriting existing exec.Command  function 
    var name string
    if len(args) > 0 { name = args[0] }
    cmd = &exec.Cmd{
        Path: name,
        Args: args,
    }
    if filepath.Base(name) == name {
        if lp, err := exec.LookPath(name); err != nil {
            return nil,err
        } else {
            cmd.Path = lp
        }
    }
    return cmd, nil
}

func GetTime()(time_now string) {
    t := time.Now()
    return t.Format(time.RFC3339Nano)
}

func ValidateSecretPhrase(phrase string)(error){
    phraseLen:=len(phrase)
    if phraseLen == 16 || phraseLen == 24 || phraseLen == 32 {
        return nil
    } else {
        return incorrectPhrase
    }
}



