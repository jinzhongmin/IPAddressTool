package main

import (
	"net"
	"strconv"

	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

type widgets struct {
	win       *gtk.Window
	IP        *gtk.Entry
	maskL     *gtk.Entry
	maskS     *gtk.Entry
	IPA       *gtk.Entry
	IPZ       *gtk.Entry
	net       *gtk.Entry
	broadcast *gtk.Entry

	view *gtk.TextView
}

func itoa(i int) string {
	return strconv.Itoa(i)
}
func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
func newWidgets() *widgets {
	w := new(widgets)
	w.win = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)

	w.win.SetPosition(gtk.WIN_POS_CENTER)
	w.win.SetTitle("IPAddressTool")
	w.win.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	})
	w.win.SetSizeRequest(250, 300)
	w.win.SetResizable(false)

	w.win.SetSizeRequest(250, 300)
	w.win.SetResizable(false)

	vbox := gtk.NewVBox(false, 0)
	w.win.Add(vbox)

	line1 := gtk.NewHBox(false, 5)
	IPLabel := gtk.NewLabel("IP地址")
	IPLabel.SetSizeRequest(60, -1)
	w.IP = gtk.NewEntry()
	w.maskL = gtk.NewEntry()
	w.maskL.SetSizeRequest(25, -1)
	line1.PackStart(IPLabel, false, false, 0)
	line1.PackStart(w.IP, true, true, 5)
	line1.PackStart(gtk.NewLabel("/"), false, false, 0)
	line1.PackEnd(w.maskL, false, false, 5)
	vbox.PackStart(line1, false, true, 5)

	line2 := gtk.NewHBox(false, 5)
	maskLabel := gtk.NewLabel("IP掩码")
	maskLabel.SetSizeRequest(60, -1)
	w.maskS = gtk.NewEntry()
	line2.PackStart(maskLabel, false, false, 0)
	line2.PackEnd(w.maskS, true, true, 5)
	vbox.PackStart(line2, false, true, 5)

	vbox.PackStart(gtk.NewHSeparator(), false, true, 5)

	line3 := gtk.NewHBox(false, 5)
	IPALabel := gtk.NewLabel("起始IP")
	IPALabel.SetSizeRequest(60, -1)
	w.IPA = gtk.NewEntry()
	line3.PackStart(IPALabel, false, false, 0)
	line3.PackEnd(w.IPA, true, true, 5)
	vbox.PackStart(line3, false, true, 5)

	line4 := gtk.NewHBox(false, 5)
	IPZLabel := gtk.NewLabel("终止IP")
	IPZLabel.SetSizeRequest(60, -1)
	w.IPZ = gtk.NewEntry()
	line4.PackStart(IPZLabel, false, false, 0)
	line4.PackEnd(w.IPZ, true, true, 5)
	vbox.PackStart(line4, false, true, 5)

	line5 := gtk.NewHBox(false, 5)
	netLabel := gtk.NewLabel("网络号")
	netLabel.SetSizeRequest(60, -1)
	w.net = gtk.NewEntry()
	line5.PackStart(netLabel, false, false, 0)
	line5.PackEnd(w.net, true, true, 5)
	vbox.PackStart(line5, false, true, 5)

	line6 := gtk.NewHBox(false, 5)
	broadcastLabel := gtk.NewLabel("广播地址")
	broadcastLabel.SetSizeRequest(60, -1)
	w.broadcast = gtk.NewEntry()
	line6.PackStart(broadcastLabel, false, false, 0)
	line6.PackEnd(w.broadcast, true, true, 5)
	vbox.PackStart(line6, false, true, 5)

	line7 := gtk.NewHBox(false, 5)
	w.view = gtk.NewTextView()
	w.view.SetEditable(false)
	line7.PackStart(w.view, true, true, 5)
	vbox.PackStart(line7, true, true, 5)

	return w
}
func do(w *widgets) {

	w.IP.Connect("activate", func() {
		if net.ParseIP(w.IP.GetText()) != nil {
			w.maskL.GrabFocus()
		} else {
			w.IPA.SetText("")
			w.IPZ.SetText("")
			w.net.SetText("")
			w.broadcast.SetText("")
			w.view.GetBuffer().SetText("不是有效的IP")
		}
	})

	maskLChanged := 0
	maskLChanged = w.maskL.Connect("changed", func() {
		w.maskL.HandlerBlock(maskLChanged)
		maskLI, err := strconv.Atoi(w.maskL.GetText())
		if err != nil || maskLI > 32 || maskLI < 0 {
			w.maskL.SetText("")
			w.view.GetBuffer().SetText("掩码长度在 0 - 32 之间")
		}
		w.maskL.HandlerUnblock(maskLChanged)
	})

	w.maskL.Connect("activate", func() {

		l := atoi(w.maskL.GetText())
		mask, _ := NewMaskByLen(l)
		w.maskS.SetText(mask.String())

		w.view.GetBuffer().SetText(w.maskL.GetText() + " 位的掩码是 " + mask.String())
		if w.IP.GetText() != "" {
			if net.ParseIP(w.IP.GetText()) != nil {
				ip, _ := NewIPByStr(w.IP.GetText())
				ip.SetMask(mask)

				w.IPA.SetText(ip.NetIP().NextIP().String())
				w.IPZ.SetText(ip.BroadcastIP().PrevIP().String())
				w.net.SetText(ip.NetIP().String())
				w.broadcast.SetText(ip.BroadcastIP().String())
			} else {
				w.IPA.SetText("")
				w.IPZ.SetText("")
				w.net.SetText("")
				w.broadcast.SetText("")
				w.view.GetBuffer().SetText("不是有效的IP")
			}

		}

	})

	w.maskS.Connect("activate", func() {
		mask, err := NewMaskByStr(w.maskS.GetText())

		if err != nil {
			w.IPA.SetText("")
			w.IPZ.SetText("")
			w.net.SetText("")
			w.broadcast.SetText("")
			w.view.GetBuffer().SetText("错误的掩码 " + w.maskS.GetText())
		} else {
			w.maskL.SetText(mask.StringLen())
			w.view.GetBuffer().SetText("掩码 " + w.maskS.GetText() + " 的长度为 " + mask.StringLen() + " 位")

			if w.IP.GetText() != "" {
				if net.ParseIP(w.IP.GetText()) != nil {
					ip, _ := NewIPByStr(w.IP.GetText())
					ip.SetMask(mask)

					w.IPA.SetText(ip.NetIP().NextIP().String())
					w.IPZ.SetText(ip.BroadcastIP().PrevIP().String())
					w.net.SetText(ip.NetIP().String())
					w.broadcast.SetText(ip.BroadcastIP().String())
				} else {
					w.IPA.SetText("")
					w.IPZ.SetText("")
					w.net.SetText("")
					w.broadcast.SetText("")
					w.view.GetBuffer().SetText("不是有效的IP")
				}
			}
		}

	})
}
func main() {
	gtk.Init(nil)

	widgets := newWidgets()
	do(widgets)
	widgets.win.ShowAll()

	gtk.Main()
}
