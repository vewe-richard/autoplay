package main

/*
TODO:
	Check command exist "gnome-screensaver-command"
	Check Files > 3
TODO:
	disable key noise
	random key and random mouse
*/

import (
	"bufio"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/micmonay/keybd_event"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
	"unsafe"
)

//install robotgo
//https://github.com/go-vgo/robotgo
var _files [][2]string
var _fileno int

func images() error {
	home, _ := os.UserHomeDir()

	file, err := os.Open(home + "/Pictures/text.log")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line[:10], line[11:])
		_files = append(_files, [2]string{home + "/Pictures/" + line[:10] + ".png", line[11:]})
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil

}

func getFile() (string, string) {
	if len(_files) < 1 {
		return "", ""
	}
	if _fileno >= len(_files) {
		_fileno = 0
	}
	p := _files[_fileno]
	_fileno += 1
	fmt.Println(p)

	return p[0], p[1]
}

func main() {
	images()
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

	pf, pt := getFile()
	image := gtk.NewImageFromFile(pf)
	window.SetTitle(pt)
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
			pf, pt := getFile()
			window.SetTitle(pt)
			image.SetFromFile(pf)
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
	day := []string{"日", "一", "二", "三", "四", "五", "六"}
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

		delay := (time.Duration)(10 + rand.Intn(200))
		time.Sleep(40 * time.Millisecond * delay) //delay
		fmt.Println("sleep ", delay)

		x, y := robotgo.GetMousePos()
		if x < 100 || y < 100 {
			skipmouse = true
		} else {
			skipmouse = false
		}

		if !skipmouse {
			if cnt%(1+rand.Intn(20)) == 0 {
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
