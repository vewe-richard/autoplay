package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
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
	mygtk()
	return
}

func mygtk() {
	runtime.GOMAXPROCS(10)
	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gdk.ThreadsEnter()
	gtk.Init(nil)

	files, err := ioutil.ReadDir("./images")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}

	fmt.Print(rand.Intn(100), ",")
	fmt.Print(rand.Intn(100))
	fmt.Println()

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
		cnt += 1
		time.Sleep(5 * 100 * time.Millisecond)
		if cnt > 10 && cnt%5 == 0 {
			x, y := robotgo.GetMousePos()
			robotgo.Move(x+5, y+5)
			if x > 1900 {
				robotgo.Move(0, y+5)
			}
			if y > 900 {
				robotgo.Move(x+5, 0)
			}
			if x > 200 && y > 200 {
				robotgo.Click()
			}
		}
	}
}

func quit() {
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press ESC to quit")
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		fmt.Printf("You pressed: rune %q, key %X\r\n", event.Rune, event.Key)
		if event.Key == keyboard.KeyEnd || event.Key == keyboard.KeyEsc {
			break
		}
	}
}
