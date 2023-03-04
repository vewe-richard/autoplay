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
	"runtime"
	"time"
)

//install robotgo
//https://github.com/go-vgo/robotgo

func main() {
	runtime.GOMAXPROCS(10)
	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gdk.ThreadsEnter()
	gtk.Init(nil)

	files, err := ioutil.ReadDir("./images")
	if err != nil {
		log.Fatal(err)
	}

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL) //WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)
	window.SetTitle("")
	window.Connect("destroy", gtk.MainQuit)
	fixed := gtk.NewFixed()

	window.Container.Add(fixed)
	window.SetSizeRequest(1880, 1030)

	image := gtk.NewImageFromFile("./images/" + files[0].Name())
	fixed.Put(image, -56, -51)
	go func() {
		for {
			time.Sleep(time.Second * 60 * 3)

			gdk.ThreadsEnter()
			r := rand.Intn(len(files))
			filepath := "./images/" + files[r].Name()
			fmt.Println(filepath)
			image.SetFromFile(filepath)
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
