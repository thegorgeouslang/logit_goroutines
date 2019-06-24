// Author: James Mallon <jamesmallondev@gmail.com>
// logit package - lib created to print and write logs
package logit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Struct type syslog -
type syslog struct {
	file       *os.File
	Filepath   string
	log        *log.Logger
	categories map[string][]string
}

// to be used as an external pointer to the syslog struct type
var Syslog *syslog

// init function - initialize values and performs a pre instantiation to make this lib
// methods work as static methods and avoid external instantiation of the struct
func init() {
	lg := syslog{} // pre instantiation
	lg.Filepath = fmt.Sprintf("%s%s.log", "logs/", time.Now().Format("2006_01_02"))
	lg.loadCategories() // loads all categories
	Syslog = &lg        // exported variable receives the instance
}

// getLogDate method - returns a string with the log format date
func (lg *syslog) getLogDate() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

// createDir function - function attempts to create the log file dir in case it doesn't exists
func (lg *syslog) createDir() (err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err = os.MkdirAll(filepath.Dir(lg.Filepath), 0755)
		if err != nil {
			msg := fmt.Sprintf("Logit error: path %s doesn't exists or is not writable and cannot be created",
				lg.Filepath)
			fmt.Printf("%s %s on %s\n", lg.getLogDate(),
				msg, lg.GetTraceMsg())
		}
		defer wg.Done()
	}()
	wg.Wait()
	return
}

// checkPath method - verifies if the directory exists and is writable
func (lg *syslog) checkPath() bool {
	if _, err := os.Stat(filepath.Dir(lg.Filepath)); os.IsNotExist(err) {
		return false
	}
	return true
}

// startLog method - processes the dir. and open the log file
func (lg *syslog) startLog() (err error) {
	ex := lg.checkPath()
	if !ex {
		err = lg.createDir()
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if err == nil {
			lg.file, _ = os.OpenFile(lg.Filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 1444)
			lg.log = log.New(lg.file, "", log.Ldate|log.Ltime)
		}
		defer wg.Done()
	}()
	wg.Wait()
	return
}

// loadCategories method - loads all categories
func (lg *syslog) loadCategories() {
	lg.categories = map[string][]string{
		"emergency": {"Emergency:", "an emergency"},
		"alert":     {"Alert:", "an alert"},
		"critical":  {"Critical:", "a critical"},
		"error":     {"Error:", "an error"},
		"warning":   {"Warning:", "a warning"},
		"notice":    {"Notice:", "a notice"},
		"info":      {"Info:", "an info"},
		"debug":     {"Debug:", "a debug"},
	}
}

// AppendCategories method - it allow the user to append new categories
func (lg *syslog) AppendCategories(newCategories map[string][]string) {
	for k, v := range newCategories {
		lg.categories[k] = v
	}
}

// WriteLog method - writes the message to the log file
func (lg *syslog) WriteLog(category string, msg string, trace string) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := lg.startLog()
		if err == nil {
			val, res := lg.categories[category]
			if !res {
				fmt.Printf("%s %s The category %s does not exists on %s\n", lg.getLogDate(),
					lg.categories["warning"][0], category, lg.GetTraceMsg())
				lg.log.Printf("%s (non existent category) %s on %s", category, msg, trace)
			} else {
				lg.log.Printf("%s %s on %s", val[0], msg, trace)
			}
			defer lg.file.Close()
		}
		defer wg.Done()
	}()
	wg.Wait()
}

// GetTraceMsg method - get the full error stack trace
func (lg *syslog) GetTraceMsg() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Sprintf("%s:%d PID: %d", frame.File, frame.Line, os.Getpid())
}
