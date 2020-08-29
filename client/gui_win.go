// +build windows

package main

import (
	"io/ioutil"
	"os"

	"github.com/getlantern/systray"
	"github.com/nynicg/cake/lib/log"
	"github.com/skratchdot/open-golang/open"
)

func loadIcon(path string) ([]byte, error) {
	f, e := os.OpenFile(path, os.O_RDONLY, 0755)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func RunAsIcon() {
	go systray.Run(onReady, func() {
		unconfigure()
		os.Exit(0)
	})
}

func onReady() {
	icob, e := loadIcon("cake.ico")
	if e != nil {
		log.Error("load icon ", e)
		panic(e)
	}
	systray.SetTemplateIcon(icob, icob)
	systray.SetTitle("Cake")
	systray.SetTooltip("Love and Spanner")
	configure()

	prox := true
	go func() {
		stt := systray.AddMenuItem("Status: ON", "")
		stt.Disable()
		servs := systray.AddMenuItem(config.ServerPerfer, "")
		for _, v := range globCfg.ProxyCfg.ServerAddr {
			s := servs.AddSubMenuItem(v, "")
			go func(addr string) {
				for {
					select {
					case <-s.ClickedCh:
						config.ServerPerfer = addr
						servs.SetTitle(addr)
					}
				}
			}(v)
		}
		systray.AddSeparator()
		update := systray.AddMenuItem("Update", "")
		runStop := systray.AddMenuItem("Stop", "")
		systray.AddSeparator()
		mQuitOrig := systray.AddMenuItem("Quit", "")
		for {
			select {
			case <-update.ClickedCh:
				open.Run("https://github.com/nynicg/cake")
			case <-runStop.ClickedCh:
				if prox {
					stt.SetTitle("Status:OFF")
					runStop.SetTitle("Run")
					unconfigure()
				} else {
					stt.SetTitle("Status:ON")
					runStop.SetTitle("Stop")
					configure()
				}
				prox = !prox
			case <-mQuitOrig.ClickedCh:
				systray.Quit()
			}
		}
	}()
}