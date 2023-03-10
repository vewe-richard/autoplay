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
	"os/exec"
	"runtime"
	"strings"
	"time"
	"unsafe"
)

//install robotgo
//https://github.com/go-vgo/robotgo
var _files []string
var _fileno int

func images() []string {
	var results []string

	home, _ := os.UserHomeDir()
	dirs := []string{home + "/Pictures"}

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
	if _fileno >= len(_files) {
		_fileno = 0
	}
	p := _files[_fileno]
	_fileno += 1
	fmt.Println(p)
	return p
}

func main() {
	_files = images()
	if len(_files) < 1 {
		log.Fatal("Please provide enough files")
	}
	_fileno = 0

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
	//fixed.Put(image, -56, -51)
	fixed.Put(image, 0, 0)

	plates := gtk.NewTextView()
	plates.SetEditable(false)
	plates.ModifyText(gtk.STATE_NORMAL, gdk.NewColor("white"))
	plates.ModifyBase(0, gdk.NewColorRGB(0x4f4f, 0x4d4d, 0x4646))
	plates.ModifyFontEasy("11")
	plates.GetBuffer().SetText(timestr())

	fixed.Put(plates, 958, (2))
	window.Fullscreen()
	go func() {
		for {
			time.Sleep(time.Second * 60 * 8)

			gdk.ThreadsEnter()
			image.SetFromFile(getFile())
			plates.GetBuffer().SetText(timestr())
			gdk.ThreadsLeave()
			if _fileno >= len(_files) {
				break
			}
		}
		fmt.Println("Exit from show picture")
		cmd := exec.Command("gnome-screensaver-command", "-l")
		cmd.Run()
	}()

	//event := make(chan interface{})

	window.Connect("key-press-event", func(ctx *glib.CallbackContext) {
		arg := ctx.Args(0)
		event := *(**gdk.EventKey)(unsafe.Pointer(&arg))
		fmt.Println(event.Keyval)
		if event.Keyval == gdk.KEY_Escape {
			window.Unfullscreen()
		} else if event.Keyval == gdk.KEY_f {
			window.Fullscreen()
		} else {
			fmt.Println("key")
		}
	})
	window.ShowAll()
	go keyevent()

	gtk.Main()
}

func timestr() string {
	day := []string{"", "一", "二", "三", "四", "五", "六", "日"}
	now := time.Now()
	return fmt.Sprintf(" %s    %02d:%02d  ", day[now.Weekday()], now.Hour(), now.Minute())
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
	robotgo.Move(105, 105)
	for {
		if _fileno >= len(_files) {
			break
		}

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

		delay := (time.Duration)(10 + rand.Intn(100)*5)
		time.Sleep(100 * time.Millisecond * delay) //delay
		fmt.Println("sleep ", delay)

		x, y := robotgo.GetMousePos()
		if x < 100 || y < 100 {
			skipmouse = true
		} else {
			skipmouse = false
		}

		if !skipmouse {
			if cnt%(1+rand.Intn(10)) == 0 {
				robotgo.Move(x+5, y+5)
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
	fmt.Println("Exit from keyevent")
	cmd := exec.Command("gnome-screensaver-command", "-l")
	cmd.Run()
}
