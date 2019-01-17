package main

import (
	"fmt"
	"log"
	"os"
)

type FileType int
type PathsMap map[FileType]string

const (
	Service FileType = iota
	Timer
	Executable
)

func checkFileReadable(path string) bool {
	if fd, err := os.Open(path); err != nil {
		log.Println("Failed to open file:", err.Error())
		return false
	} else {
		if fd.Close() == nil {
			return true
		} else {
			log.Println("Failed to close file:", err.Error())
			return false
		}
	}
}

func (t FileType) sourcePath(serviceName string) string {
	var m PathsMap = map[FileType]string{
		Service:    fmt.Sprintf("systemd/%s/%s.service", serviceName, serviceName),
		Timer:      fmt.Sprintf("systemd/%s/%s.timer", serviceName, serviceName),
		Executable: fmt.Sprintf("systemd/%s/%s-service", serviceName, serviceName),
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
		Service:    fmt.Sprintf("systemd/%s/%s.service", serviceName, serviceName),
		Timer:      fmt.Sprintf("systemd/%s/%s.timer", serviceName, serviceName),
		Executable: fmt.Sprintf("systemd/%s/%s-service", serviceName, serviceName),
	}

	return m.getOrPanic(t)
}

func (t FileType) copyPath(serviceName string) (from string, to string) {
	return t.sourcePath(serviceName), t.installPath(serviceName)
}

func help() {
	fmt.Println("Available commands: list, install, remove")
	os.Exit(0)
}

func secondArgOrPanic() string {
	if len(os.Args) < 3 {
		log.Panicln("command argument required")
	}
	return os.Args[2]
}

func install(serviceName string) {
	fmt.Println(Executable.copyPath(serviceName))
	fmt.Println(Service.copyPath(serviceName))
	fmt.Println(Timer.copyPath(serviceName))
}

func list() {

}

func remove(serviceName string) {

}

func main() {
	if len(os.Args) < 2 {
		help()
	}
	switch os.Args[1] {
	case "list":
		list()
		break
	case "install":
		install(secondArgOrPanic())
		break
	case "remove":
		remove(secondArgOrPanic())
		break
	default:
		help()
	}
}
