package main

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

type FileType int
type PathsMap map[FileType]string

const (
	Service FileType = iota
	Timer
	Executable
)

var fileTypes = [...]FileType{Service, Timer, Executable}

const serviceSourceDir = "systemd"

var excludeServices = map[string]bool{"template": true}
var legacyServices = [...]string{"off-on"}

func CloseOrPanic(c io.Closer) {
	if c.Close() != nil {
		log.Panicf("Failed to close resource")
	}
}

func checkFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func checkFileReadable(path string) bool {
	if !checkFileExists(path) {
		log.Printf("File %s doesnt exists", path)
		return false
	}
	fd, err := os.Open(path)
	defer CloseOrPanic(fd)
	if err != nil {
		log.Println("Failed to open file:", err.Error())
		return false
	}
	return true
}

func (t FileType) sourcePath(serviceName string) string {
	var m PathsMap = map[FileType]string{
		Service:    fmt.Sprintf("%s/%s/%s.service", serviceSourceDir, serviceName, serviceName),
		Timer:      fmt.Sprintf("%s/%s/%s.timer", serviceSourceDir, serviceName, serviceName),
		Executable: fmt.Sprintf("%s/%s/%s-service", serviceSourceDir, serviceName, serviceName),
	}
	fileName := m.getOrPanic(t)
	if !checkFileReadable(fileName) {
		log.Panicln("File is not readable")
	}
	return fileName
}

func (m PathsMap) getOrPanic(fType FileType) string {
	if result, found := m[fType]; found {
		return result
	} else {
		log.Panicf("Mapping for file type %d is not found\n", fType)
		return "unreachable"
	}
}

func (t FileType) installPath(serviceName string) string {
	var m PathsMap = map[FileType]string{
		Service:    fmt.Sprintf("/etc/systemd/system/%s.service", serviceName),
		Timer:      fmt.Sprintf("/etc/systemd/system/%s.timer", serviceName),
		Executable: fmt.Sprintf("/usr/local/bin/%s-service", serviceName),
	}

	return m.getOrPanic(t)
}

func (t FileType) copyPath(serviceName string) (from string, to string) {
	return t.sourcePath(serviceName), t.installPath(serviceName)
}

func help() {
	fmt.Printf("Available commands: { %s }\n", strings.Join(commandNames, " | "))
	os.Exit(0)
}

func secondArgOrPanic() string {
	if len(os.Args) < 3 {
		log.Panicln("command argument required")
	}
	return os.Args[2]
}

func install() {
	serviceName := secondArgOrPanic()
	if serviceName == "template" {
		log.Panicln("Installing template service is forbidden")
	}
	for _, fType := range fileTypes {
		fromPath, toPath := fType.copyPath(serviceName)
		from, err := os.Open(fromPath)
		if err != nil {
			log.Panicln("Failed to open source file.", err)
		}
		var mode os.FileMode
		if fType == Executable {
			mode = os.FileMode(0744)
		} else {
			mode = os.FileMode(0644)
		}
		to, err := os.OpenFile(toPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, mode)
		if err != nil {
			log.Panicln("Failed to open destination file.", err)
		}

		if _, err := io.Copy(to, from); err != nil {
			log.Panicln("Failed to copy file.", err)
		}
		fmt.Printf("[%s] %s -> %s\n", aurora.Blue(" Copied"), fromPath, toPath)
	}
	timerName := serviceName + ".timer"
	runCommand("systemctl enable " + timerName, true)
	fmt.Printf("[%s] %s\n", aurora.Green("Enabled"), timerName)
	runCommand("systemctl start " + timerName, true)
	fmt.Printf("[%s] %s\n", aurora.Green("Started"), timerName)
}

func runCommand(cmd string, panicOrError bool) {
	parts := strings.Split(cmd, " ")
	command := exec.Command(parts[0], parts[1:]...)
	if err := command.Run(); err != nil && panicOrError	{
		log.Panicln(err)
	}
}

type InstallStatus int
type colorFunc func(interface{}) aurora.Value

const (
	Installed InstallStatus = iota
	Missing
	Broken
)

func (s InstallStatus) String() string {
	return [...]string{"installed", "missing", "broken"}[s]
}

func (s InstallStatus) colorFunc(isLegacy bool) colorFunc {
	if s == Broken {
		return aurora.Brown
	}
	success := s == Installed && !isLegacy || s == Missing && isLegacy

	if success {
		return aurora.Green
	} else {
		return aurora.Red
	}
}

func serviceInstallStatus(serviceName string) InstallStatus {
	hasExecutable := checkFileExists(Executable.installPath(serviceName))
	hasTimer := checkFileExists(Timer.installPath(serviceName))
	hasServiceFile := checkFileExists(Service.installPath(serviceName))
	isFullyInstalled := hasExecutable && hasServiceFile && hasTimer
	if isFullyInstalled {
		return Installed
	}
	isPartiallyInstalled := hasExecutable || hasServiceFile || hasTimer
	if isPartiallyInstalled {
		return Broken
	}
	return Missing
}

func listServiceStatus(serviceName string, isLegacy bool) {
	state := serviceInstallStatus(serviceName)
	status := state.colorFunc(isLegacy)(state.String())
	fmt.Printf("[%9s] %s\n", status, serviceName)
}

func list() {
	fmt.Printf("[%s] %s\n", aurora.Green(" Syncing"), "git pull")
	runCommand("git pull", true)
	fmt.Printf("[%s] %s\n", aurora.Green("Building"), "go build")
	runCommand("go build", true)
	result, err := ioutil.ReadDir(serviceSourceDir)
	if err != nil {
		log.Panicf("Cant open dir %s", serviceSourceDir)
	}
	fmt.Println("Available services:")
	for _, file := range result {
		serviceName := file.Name()
		if excludeServices[serviceName] {
			continue
		}
		listServiceStatus(serviceName, false)
	}
	fmt.Println("\nLegacy services:")
	for _, serviceName := range legacyServices {
		listServiceStatus(serviceName, true)
	}
}

func remove() {
	serviceName := secondArgOrPanic()
	status := serviceInstallStatus(serviceName)
	isInstalled := status == Installed
	if status == Missing {
		fmt.Println("Already removed!")
		return
	}
	timerName := serviceName + ".timer"
	recover()
	runCommand("systemctl stop " + timerName, isInstalled)
	fmt.Printf("[%s] %s\n", aurora.Green(" Stopped"), timerName)
	runCommand("systemctl disable " + timerName, isInstalled)
	fmt.Printf("[%s] %s\n", aurora.Green("Disabled"), timerName)
	for _, fType := range fileTypes {
		path := fType.installPath(serviceName)
		if !checkFileExists(path) {
			continue
		}
		if err := os.Remove(path); err != nil {
			log.Panicln(err)
		}
		fmt.Printf("[%s] %s\n", aurora.Red(" Removed"), path)
	}
}

func template() {

}

func osTweaks() {
	/*
	systemctl disable apt-daily-upgrade.timer
systemctl stop apt-daily-upgrade.timer
systemctl disable apt-daily.timer
systemctl stop apt-daily.timer
	 */
}

var commands = map[string]func(){
	"list":     list,
	"install":  install,
	"remove":   remove,
	"template": template,
	"tweak":    osTweaks,
}
var commandNames = getCommandNames()

func getCommandNames() []string {
	var result []string
	for _, v := range reflect.ValueOf(commands).MapKeys() {
		result = append(result, v.String())
	}
	return result
}

func main() {
	if len(os.Args) < 2 {
		list()
		return
	}
	if f, found := commands[os.Args[1]]; found {
		f()
	} else {
		help()
	}
}
