package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/micmonay/keybd_event"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
)

//install robotgo
//https://github.com/go-vgo/robotgo
var _files []string

func images() []string {
	var results []string

	home, _ := os.UserHomeDir()
	dirs := []string{"./images", home + "/Pictures"}

	for _, d := range dirs {
		files, err := ioutil.ReadDir(d)

		if err == nil {
			for _, f := range files {
				if strings.Contains(f.Name(), "Screenshot") && strings.Contains(f.Name(), "png") {
					results = append(results, d+"/"+f.Name())
				}
			}
		}
	}
	return results
}

func getFile() string {
	if len(_files) < 1 {
		return ""
	}
	fmt.Println(_files)
	r := rand.Intn(len(_files))
	p := _files[r]
	_files = append(_files[:r], _files[r+1:]...)
	fmt.Println(p)
	return p
}

func main() {
	_files = images()
	if len(_files) < 1 {
		log.Fatal("Please provide enough files")
	}

	runtime.GOMAXPROCS(10)
	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gdk.ThreadsEnter()
	gtk.Init(nil)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL) //WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)
	window.SetTitle("")
	window.Connect("destroy", gtk.MainQuit)
	fixed := gtk.NewFixed()

	window.Container.Add(fixed)
	window.SetSizeRequest(1880, 1030)

	image := gtk.NewImageFromFile(getFile())
	fixed.Put(image, -56, -51)
	go func() {
		for {
			time.Sleep(time.Second * 60 * 8)

			gdk.ThreadsEnter()
			image.SetFromFile(getFile())
			gdk.ThreadsLeave()
		}
	}()

	window.ShowAll()
	go keyevent()
	gtk.Main()
}

func keyevent() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	keys := []int{keybd_event.VK_UP, keybd_event.VK_DOWN, keybd_event.VK_LEFT, keybd_event.VK_RIGHT}
	skipmouse := false
	cnt := 0
	for {
		// Select keys to be pressed
		r := rand.Intn(len(keys))
		kb.SetKeys(keys[r])

		// Press the selected keys
		err = kb.Launching()
		if err != nil {
			panic(err)
		}

		// Or you can use Press and Release
		kb.Press()
		time.Sleep(10 * time.Millisecond)
		kb.Release()

		delay := (time.Duration)(10 + rand.Intn(10)*5)
		time.Sleep(100 * time.Millisecond * delay) //delay
		fmt.Println("sleep ", delay)

		x, y := robotgo.GetMousePos()
		if x < 100 || y < 100 {
			skipmouse = true
		} else {
			skipmouse = false
		}

		if !skipmouse {
			robotgo.Move(x+5, y+5)
			if cnt%3 == 0 {
				fmt.Println("click mouse")
				robotgo.Click()
			}
			if x > 1900 {
				robotgo.Move(110, y+5)
			}
			if y > 900 {
				robotgo.Move(x+5, 110)
			}
		}
		cnt += 1

	}
}
