package common

import "os"
import "os/exec"
import "path/filepath"

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


