package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	now := time.Now().Unix()
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("homedir", err)
	}

	err = exec.Command("gnome-screenshot", "-f", fmt.Sprintf(dirname+"/Pictures/%d.png", now)).Run()
	if err != nil {
		log.Fatal("screenshot", err)
	}

	out, _ := exec.Command("xdotool", "getactivewindow", "getwindowname").Output()
	line := fmt.Sprintf("%d:%s", now, out)

	f, err := os.OpenFile(dirname+"/Pictures/text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("open log", err)
	}
	defer f.Close()

	if _, err := f.WriteString(line); err != nil {
		log.Fatal(err)
	}
}
