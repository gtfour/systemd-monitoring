package filter

import "fmt"
import "errors"
import "systemd-monitoring/common"
import "jumper/cuda/analyze"
import "jumper/cuda/handling"

var emptyKeywordSlice   = errors.New("Keyword slice is empty")
var fileDoesntExist     = errors.New("File doesn't exist")
var nothingToSendYet    = errors.New("Nothing to send yet")
var wrongArea           = errors.New("Wrong area")
var pythonTracebackArea = "python-traceback"

type PythonTracebackHandlerSet struct {
    handlers []*PythonTracebackHandler
}





func(pths *PythonTracebackHandlerSet)Processing(messageInput common.DataUpdate)(messageOutput common.DataUpdate, err error){

    //
    messagePath     := messageInput.Path
    messageArea     := messageInput.Area
    messageHostname := messageInput.Hostname
    //
    if messageArea != pythonTracebackArea { return messageOutput,wrongArea }
    handlers    := pths.handlers
    for i := range handlers {
        handler := handlers[i]
        if messagePath == handler.Path {
            messageText := messageInput.Text
            entry,err   := handler.strProcessing( messageText )
            if err == nil {
                messageOutput.Path      = messagePath
                messageOutput.Area      = messageArea
                messageOutput.Hostname  = messageHostname
                messageOutput.Timestamp = common.GetTime()
                messageOutput.Text      = entry.String()
                return messageOutput, nil

            }
        }
    }
    //
    return messageOutput, nothingToSendYet
}


func(pths *PythonTracebackHandlerSet)AppendHandler(path string, header_keyword []string)(err error){
    //
    var handler PythonTracebackHandler
    if len(header_keyword)==0   { return emptyKeywordSlice }
    if !common.FileExists(path) { return fileDoesntExist   }
    matcher := func(phrase string)(bool){
        _, _, section_type := analyze.EscapeIndentSection(phrase, header_keyword)
        if section_type == analyze.INDENT_SECTION {
            return true
        } else {
            return false
        }
    }
    breaker  := handling.GetSectionBreaker("", [2]int{}, [2]int{}, analyze.INDENT_SECTION)
    new_id,_ := common.GenId()

    handler.Id      = new_id
    handler.Path    = path
    handler.Matcher = matcher
    handler.Breaker = breaker

    pths.handlers = append(pths.handlers, &handler)
    return nil

}

type PythonTracebackHandler struct {
    //
    Id           string
    Path         string
    Matcher      func(string)(bool)
    Breaker      func(string)(bool)
    CurrentEntry *pythonTracebackEntry
    CreateNew    bool
    //
}

func(pth *PythonTracebackHandler)strProcessing(inputString string)(pte *pythonTracebackEntry,err error){
    //
    //
    //
    fmt.Printf("\n=== Testing matcher: %v Input String: %v ===\n",inputString,pth.Matcher(inputString))
    //
    //
    //
    if pth.CurrentEntry == nil || pth.CreateNew {
        var NewEntry pythonTracebackEntry
        pth.CurrentEntry = &NewEntry
        pth.CreateNew = false
    }
    if pth.Matcher(inputString) && !pth.CurrentEntry.InProgress {
        pth.CurrentEntry.Header     = inputString
        pth.CurrentEntry.InProgress = true
        return nil, nothingToSendYet
    }
    if pth.Breaker(inputString) && pth.CurrentEntry.InProgress {
        pth.CurrentEntry.Footer     = inputString
        pth.CurrentEntry.InProgress = false
        pth.CreateNew               = true
        return pth.CurrentEntry, nil
    }
    if pth.CurrentEntry.InProgress && !pth.Breaker(inputString) {
        pth.CurrentEntry.Lines = append(pth.CurrentEntry.Lines, inputString)
    }
    return nil, nothingToSendYet
    //
    //
    //
}



type pythonTracebackEntry struct {
    //
    //Path     string
    Header     string
    Lines      []string
    InProgress bool
    Footer     string
    //
}

func(pte *pythonTracebackEntry)String()(res string){
    res+=pte.Header+"\n"
    for i := range pte.Lines {
        line := pte.Lines[i]
        res += line+"\n"
    }
    res+=pte.Footer
    return res
}

