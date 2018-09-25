package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/criminalist/gopress/helper"
	"github.com/criminalist/gopress/module/core"
	"github.com/criminalist/gopress/module/hook"

	"github.com/insionng/makross"

	"github.com/fsnotify/fsnotify"

	_ "qlang.io/lib/builtin"
)

var (
	wg   sync.WaitGroup
	port int = 9999
	app  *makross.Makross
)

func init() {
	defer core.VmString(`
		hook.DoActionHook("init")
	`)

	core.VmString(`
		hook.AddActionHook("init", fn {
			println("<Application init>")
		})
	`)
}

func quit() {
	fmt.Println("<Application quit>")
}

/*
func TheTitle() []byte {
	quote := "The bird is the word."
	return []byte(quote)
}

func changeTheQuote(quote []byte) []byte {
	quote = []byte(strings.Replace(string(quote), "bird", "nerd", -1))
	fmt.Println(string(quote))
	return quote
}

func bigTitle(quote []byte) []byte {
	quote = []byte("<<<" + string(quote) + ">>>")
	fmt.Println(string(quote))
	return quote
}
*/

func main() {
	defer func() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			core.AddFunc("quit", quit)
			core.VmString(`
			hook.AddActionHook("quit", quit)
			hook.DoActionHook("quit")
`)
		}()
		wg.Wait()
	}()

	core.VmString(`
		hook.AddActionHook("bootstrap", fn {
			println("<Application bootstrap>")
		})
	   	hook.DoActionHook("bootstrap")
	   `)

	//------------------------------------------------------//

	core.Plugins()
	hook.DoActionHook("plugin")

	//------------------------------------------------------//

	//------------------------------------------------------//

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	//Execute the theme logic
	var theme = "default"
	var filter, reload bool
	filter = true
	reload = true

	var okay bool
	app, okay = core.GetAppByTheme(theme, filter, reload)

	/*------------------------------------*/
	/*
		m.Use(toolbox.Toolboxer(m, toolbox.Options{
			HealthCheckFuncs: []*toolbox.HealthCheckFuncDesc{
				&toolbox.HealthCheckFuncDesc{
					Desc: "Database connection",
					Func: models.Ping,
				},
			},
		}))
	*/
	/*------------------------------------*/

	if okay {
		fmt.Printf("app.Listen(%v)\n", port)
		go app.Listen(port)
	} else {
		panic("cannot convert app to (*makross.Makross)")
	}

	// Process events
	go func() {
		for {
			select {
			case <-watcher.Events:
				// Loading App Logic
				fmt.Println("Reload Application")
				if app == nil {
					panic("app == nil")
				}
				app.Close()

				app, okay = core.GetAppByTheme(theme, filter, reload)
				if okay {
					fmt.Printf("app.Listen(%v)\n", port)
					go app.Listen(port)
				} else {
					panic("cannot convert app to (*makross.Makross)")
				}

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	fmt.Println(".........................................................")
	fmt.Printf("Application pid is %d\n", os.Getpid())
	fmt.Println(".........................................................")

	//Hot update monitoring directory
	var applicationDir = "content/application"
	var applicationRootDir = fmt.Sprintf("%s/root", applicationDir)

	var themeAppDir = fmt.Sprintf("content/theme/%s/handler", theme)
	var themeAppRootDir = fmt.Sprintf("%s/root", themeAppDir)

	err = watcher.Add(applicationDir)
	if err != nil {
		log.Fatal(fmt.Sprintf("watcher.Add applicationDir has error:%v", err))
	}

	err = watcher.Add(applicationRootDir)
	if err != nil {
		log.Fatal(fmt.Sprintf("watcher.Add applicationRootDir has error:%v", err))
	}

	err = watcher.Add(themeAppDir)
	if err != nil {
		log.Fatal(fmt.Sprintf("watcher.Add themeAppDir has error:%v", err))
	}

	if helper.IsExist(themeAppRootDir) {
		err = watcher.Add(themeAppRootDir)
		if err != nil {
			log.Fatal(fmt.Sprintf("watcher.Add themeAppRootDir has error:%v", err))
		}
	}

	// Hang so program doesn't exit
	<-done
	watcher.Close()
}
