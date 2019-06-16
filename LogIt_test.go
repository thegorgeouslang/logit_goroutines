// Author: James Mallon <jamesmallondev@gmail.com>
// logit package - lib created to print and write logs
package logit

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"
)

// init function - data and process initialization
func init() {
	Syslog.Filepath = "log/"
}

// Test function TestGetLogDate to evaluate getLogDate
func TestGetLogDate(t *testing.T) {

	logDate := Syslog.getLogDate()
	currDate := time.Now().Format("2006/01/02 15:04:05")

	if logDate != currDate {
		t.Errorf("Expected return from getLogDate to be the current date %s, but got %s ", currDate, logDate)
	}
}

// Test function TestCreateDir to evaluate the createDir method
func TestCreateDir(t *testing.T) {
	Syslog.createDir() // creates the folder
	_, e := os.Stat(Syslog.Filepath)
	if e != nil { // check for non existent dir
		t.Errorf("Expected the directory to exists.")
	}
	os.Remove(Syslog.Filepath) // remove the dir
}

// Test function TestCheckPath to evaluate the checkPath method
func TestCheckPath(t *testing.T) {
	e := Syslog.checkPath()
	if e == nil { // check for non existent dir
		t.Errorf("Expected the directory to not exists.")
	}
}

// Test function TestStartLog to evaluate startLog method
func TestStartLog(t *testing.T) {
	Syslog.startLog() // check the existence of the folder and create it
	_, e := os.Stat(Syslog.Filepath)
	if e != nil { // check for non existent dir
		t.Errorf("Expected the directory to exists.")
	}
	os.Remove(Syslog.Filepath) // remove the dir
}

// Test function TestLoadCategories to evaluate loadCategories method
func TestLoadCategories(t *testing.T) {
	Syslog.loadCategories()
	if Syslog.categories["alert"][0] != "Alert:" {
		t.Errorf("Expected Syslog.categories[\"alert\"][0] == \"Alert\", but got %s", Syslog.categories["alert"][0])
	}
}

// Test function TestAppendCategories to evaluate AppendCategories
func TestAppendCategories(t *testing.T) {
	newCategory := map[string][]string{
		"checkpoint": {"Checkpoint:", "150.000.000,00"},
	}
	Syslog.AppendCategories(newCategory)
	if Syslog.categories["checkpoint"][0] != "Checkpoint:" {
		t.Errorf("Expected Checkpoint:, but got %s ", Syslog.categories["checkpoint"][0])
	}
}

// Test function TestWriteLog to evaluate WriteLog method
func TestWriteLog(t *testing.T) {
	Syslog.Filepath = fmt.Sprintf("%s%s.log", "logs/", time.Now().Format("2006_01_02"))
	Syslog.WriteLog("debug", "Testing...", Syslog.GetTraceMsg())

	// open and read the first line of the log file
	file, _ := os.Open(Syslog.Filepath)
	fs := bufio.NewScanner(file)
	fs.Scan()
	fline := fs.Text()

	// check for must have text
	match, _ := regexp.MatchString(".*Debug:.*Testing", fline)
	if !match {
		t.Errorf("Expected to find Debug: in the file")
	}

	os.Remove(Syslog.Filepath) // remove the file
	os.Remove("logs/")         // remove the dir
}

// Test function BenchmarkWriteLog to evaluate the WriteLog method
func BenchmarkWriteLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Syslog.Filepath = fmt.Sprintf("%s%s.log", "logs/", time.Now().Format("2006_01_02"))
		Syslog.WriteLog("debug", "Testing...", Syslog.GetTraceMsg())

		os.Remove(Syslog.Filepath) // remove the file
		os.Remove("logs/")         // remove the dir
	}
}

// Test function TestGetTraceMsg to evaluate GetTraceMsg method
func TestGetTraceMsg(t *testing.T) {
	pattern := fmt.Sprintf(".*PID: %d", os.Getpid())
	match, _ := regexp.MatchString(pattern, Syslog.GetTraceMsg())
	if !match {
		t.Errorf("Expected to match the PID")
	}
}
