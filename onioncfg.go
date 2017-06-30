//*ClientTransportPlugin transport socks4|socks5 IP:PORT
//*HTTPProxy host[:port]
//*HTTPProxyAuthenticator username:password
//*HTTPSProxy host[:port]
//*HTTPSProxyAuthenticator username:password
//*Socks4Proxy host[:port]
//*Socks5Proxy host[:port]
//*Socks5ProxyUsername username
//*Socks5ProxyPassword password




package main

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
	"fmt"
	"strings"
	"github.com/yawning/bulb"
)


var mainWin *gtk.Window
var Notebook *gtk.Notebook


func promptInfo(msg string) {
	dialog := gtk.MessageDialogNew(mainWin, 0, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "Displaying full log info:")
//	dialog.SetDefaultGeometry(500, 200)

	tv, err := gtk.TextViewNew()

	if err != nil {
		log.Fatal("Unable to create TextView:", err)
	}

	tvbuf, err := tv.GetBuffer()

	if err != nil {
		log.Fatal("Unable to get buffer:", err)
	}

	tvbuf.SetText(msg)
	tv.SetEditable(false)
	tv.SetWrapMode(gtk.WRAP_WORD)

	scrollbox, err := gtk.ScrolledWindowNew(nil, nil)

	if err != nil {
		log.Fatal("Unable to create scrolled window:", err)
	}

	scrollbox.Add(tv)
	scrollbox.SetSizeRequest(600, 100)

	box, err := dialog.GetContentArea()

	if err != nil {
		log.Fatal("Unable to get content area of dialog:", err)
	}

	box.Add(scrollbox)
	dialog.ShowAll()
	dialog.Run()
	dialog.Destroy()
//self.set_default_size(150, 100)
}

func promptChoice(msg string) int {
	dialog := gtk.MessageDialogNew(mainWin, 0, gtk.MESSAGE_ERROR, gtk.BUTTONS_YES_NO, msg)
	result := dialog.Run()
	dialog.Destroy()
	return result
}

func promptError(msg string) {
	dialog := gtk.MessageDialogNew(mainWin, 0, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE, "Error: %s", msg)
	dialog.Run()
	dialog.Destroy()
}

func get_scrollbox() *gtk.ScrolledWindow {
	sbox, err := gtk.ScrolledWindowNew(nil, nil)

	if err != nil {
		log.Fatal("Unable to create scrolled window:", err)
	}

	return sbox
}

func get_frame(text string) *gtk.Frame {
	frame, err := gtk.FrameNew(text)

	if err != nil {
		log.Fatal("Unable to create frame:", err)
	}

	return frame
}

func get_hbox() *gtk.Box {
        hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)

        if err != nil {
                log.Fatal("Unable to create horizontal box:", err)
        }

        return hbox
}

func get_vbox() *gtk.Box {
        vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

        if err != nil {
                log.Fatal("Unable to create vertical box:", err)
        }

        return vbox
}

func get_entry(text string) *gtk.Entry {
        entry, err := gtk.EntryNew()

        if err != nil {
                log.Fatal("Unable to create text entry:", err)
        }

        entry.SetText(text)
        return entry
}

func get_textview(text string) *gtk.TextView {
	buf, err := gtk.TextBufferNew(nil)

	if err != nil {
		log.Fatal("Unable to create text buffer:", err)
	}

	buf.SetText(text);

	tv, err := gtk.TextViewNewWithBuffer(buf)

	if err != nil {
		log.Fatal("Unable to create textview:", err)
	}

	return tv
}

func get_label(text string) *gtk.Label {
	label, err := gtk.LabelNew(text)

	if err != nil {
		log.Fatal("Unable to create label in GUI:", err)
		return nil
	}

	return label
}

func get_radiobutton(group *gtk.RadioButton, label string, activated bool) *gtk.RadioButton {

        if group == nil {
                radiobutton, err := gtk.RadioButtonNewWithLabel(nil, label)

                if err != nil {
                        log.Fatal("Unable to create radio button:", err)
                }

                radiobutton.SetActive(activated)
                return radiobutton
        }

        radiobutton, err := gtk.RadioButtonNewWithLabelFromWidget(group, label)

        if err != nil {
                log.Fatal("Unable to create radio button in group:", err)
        }

        radiobutton.SetActive(activated)
        return radiobutton
}

func get_checkbox(text string, activated bool) *gtk.CheckButton {
        cb, err := gtk.CheckButtonNewWithLabel(text)

        if err != nil {
                log.Fatal("Unable to create new checkbox:", err)
        }

        cb.SetActive(activated)
        return cb
}

func get_button(label string) *gtk.Button {
        button, err := gtk.ButtonNewWithLabel(label)

        if err != nil {
                log.Fatal("Unable to create new button:", err)
        }

        return button
}

func get_edit_text(edit *gtk.Entry) string {
	text, err := edit.GetText()

	if err != nil {
		log.Fatal("Unable to read entry text:", err)
	}

	return text
}

func emitConfigFree(c *bulb.Conn, val string) {
	fmt.Println("SETCONF / " + val)
	resp, err := c.Request("SETCONF " + val)
	if err != nil {
		promptError("SETCONF on " + val + " failed:" + err.Error())
		log.Fatal("SETCONF on " + val + " failed:", err)
	}
	log.Println("SETCONF response: ", resp)
}

func emitConfig(c *bulb.Conn, key, val string) {
	fmt.Println("SETCONF / " + key + " = "  + val)
	resp, err := c.Request("SETCONF " + key + "=\"" + val + "\"")
	if err != nil {
		promptError("SETCONF on " + key + " failed:" + err.Error())
		log.Fatal("SETCONF on " + key + " failed:", err)
	}
	log.Println("SETCONF response: ", resp)
}

func main() {
	gtk.Init(nil)

	c, err := bulb.Dial("unix", "/var/run/tor/control")
	if err != nil {
		log.Fatal("failed to connect to control port:", err)
	}
	defer c.Close()
	c.Debug(true)

	if err := c.Authenticate(""); err != nil {
		log.Fatal("Tor control port authentication failed:", err)
	}

	resp, err := c.Request("GETINFO version")
	if err != nil {
		log.Fatalf("GETINFO version failed: %v", err)
	}
	log.Printf("GETINFO version: %v", resp)

	resp, err = c.Request("GETCONF bridge")
	if err != nil {
		log.Fatalf("GETINFO bridges failed: %v", err)
	}
	log.Printf("GETINFO bridges: %v", resp)

	// Create a new toplevel window, set its title, and connect it to the "destroy" signal to exit the GTK main loop when it is destroyed.
	mainWin, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)

	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	mainWin.SetTitle("Subgraph Onion Bridge Configuration")

	mainWin.Connect("destroy", func() {
		fmt.Println("Shutting down...")
	        gtk.MainQuit()
	})

	mainWin.SetPosition(gtk.WIN_POS_CENTER)

	Notebook, err = gtk.NotebookNew()

	if err != nil {
		log.Fatal("Unable to create new notebook:", err)
	}

	nbLabel, err := gtk.LabelNew("Tor Network Setup")

	if err != nil {
		log.Fatal("Unable to create notebook label:", err)
	}

	box := get_vbox()
	scrollbox := get_scrollbox()
	scrollbox.Add(box)
	hlb := get_vbox()
//	lb := get_label("Tor setup options:")
//	hlb.PackStart(lb, false, false, 0)
	rb := get_radiobutton(nil, "Direct connection (default)", true)
	rb2 := get_radiobutton(rb, "Configure bridge/proxy (restricted Internet access)", false)


	box.Add(hlb)
	box.Add(rb)
	box.Add(rb2)

//	lb = get_label("Input bridges (one per line):")
//		box.PackStart(lb, false, true, 10)
//	box.Add(lb)

	bframe := get_frame("Input bridges (one per line):")
	box.PackStart(bframe, false, false, 10)

	tvsb := get_scrollbox()
	tvsb.SetSizeRequest(0, 55)
	tv := get_textview("")
	tv.SetBorderWidth(5)
	tvsb.Add(tv)
	bframe.Add(tvsb)
	bframe.SetSensitive(false)

	rb.Connect("toggled", func() {
		bframe.SetSensitive(!rb.GetActive())
	})

	npcheck := get_checkbox("Use network proxy", false)
	box.PackStart(npcheck, false, false, 10)


	pframe := get_frame("Network proxy settings")
	fbox := get_vbox()
	hbox := get_hbox()
	s4radio := get_radiobutton(nil, "SOCKS4", false)
	s5radio := get_radiobutton(s4radio, "SOCKS5", true)
	hproxy := get_radiobutton(s4radio, "HTTP proxy", false)
	hsproxy := get_radiobutton(s4radio, "HTTPS proxy", false)
	hbox.PackStart(s4radio, false, true, 0)
	hbox.PackStart(s5radio, false, true, 20)
	hbox.PackStart(hproxy, false, true, 10)
	hbox.PackStart(hsproxy, false, true, 10)
	fbox.Add(hbox)
	pframe.Add(fbox)

	hbox = get_hbox()
	slabel := get_label("SOCKS/Proxy address:")
	hbox.PackStart(slabel, false, true, 5)
	addredit := get_entry("")
	addredit.SetPlaceholderText("Host:port")
	hbox.PackStart(addredit, false, true, 10)
	fbox.Add(hbox)

	hbox = get_hbox()
	ulabel := get_label("Username:")
	hbox.PackStart(ulabel, false, true, 5)
	fbox.Add(hbox)
	uedit := get_entry("")
	hbox.PackStart(uedit, false, true, 10)
	plabel := get_label("Password:")
	hbox.PackStart(plabel, false, true, 5)
	pedit := get_entry("")
	hbox.PackStart(pedit, false, true, 10)
	fbox.Add(get_label(""))

	box.Add(pframe)
	pframe.SetSensitive(false)


	npcheck.Connect("toggled", func() {
		pframe.SetSensitive(npcheck.GetActive())
	})


	btnbox := get_hbox()
	btn := get_button("Go configure")
	btnbox.PackStart(btn, false, false, 10)
	box.PackStart(btnbox, false, false, 10)

	btn.Connect("clicked", func() {
		addrempty := strings.TrimSpace(get_edit_text(addredit)) == ""
		uempty := strings.TrimSpace(get_edit_text(uedit)) == ""
		pempty := strings.TrimSpace(get_edit_text(pedit)) == ""

		if rb2.GetActive() {
			tvb, err := tv.GetBuffer()
			if err != nil {
				promptError("Error reading user-supplied bridge data")
				return
			}

			text, err := tvb.GetText(tvb.GetStartIter(), tvb.GetEndIter(), false)
			if err != nil {
				promptError("Error reading user-supplied bridge data")
				return
			}

			found := false
			lines := strings.Split(text, "\n")

			for _, line := range lines {
				line  = strings.TrimSpace(line)
				if line == "" {
					continue
				}

				fmt.Println("Adding bridge: ", line)
				emitConfig(c, "Bridge", line)
				found = true
			}

			if !found {
				promptError("No bridges supplied by user")
				return
			} else {
				fmt.Println("Enabling UseBridges")
				emitConfig(c, "UseBridges", "1")
			}

		}

		if npcheck.GetActive() {

			if addrempty {
				promptError("SOCKS/proxy address field cannot be left blank!")
				return
			}

			if (!uempty || !pempty) && (uempty != pempty) {
				promptError("Use of credentials requires both a username and a password!")
				return
			}

			if s4radio.GetActive() {
				emitConfig(c, "Socks4Proxy", get_edit_text(addredit))
			} else if s5radio.GetActive() {
				emitConfig(c, "Socks5Proxy", get_edit_text(addredit))
			} else if hproxy.GetActive() {
				emitConfig(c, "HTTPProxy", get_edit_text(addredit))

				if !uempty && !pempty {
					emitConfig(c, "HTTPProxyAuthenticator", get_edit_text(uedit) + ":" + get_edit_text(pedit))
				}
			} else if hsproxy.GetActive() {
				emitConfig(c, "HTTPSProxy", get_edit_text(addredit))

				if !uempty && !pempty {
					emitConfig(c, "HTTPSProxyAuthenticator", get_edit_text(uedit) + ":" + get_edit_text(pedit))
				}
			}

/*			if rb2.GetActive() {
				if s4radio.GetActive() {
					emitConfigFree(c, "ClientTransportPlugin ... socks4 " + get_edit_text(addredit))
				} else if s5radio.GetActive() {
					if !uempty && !pempty {
						emitConfigFree(c, "ClientTransportPlugin ... socks5 " + get_edit_text(addredit))
					} else {
						cfgline := fmt.Sprintf("ClientTransportPlugin ... socks5 %s username=%s password=%s",
							get_edit_text(addredit), get_edit_text(uedit), get_edit_text(pedit))
						emitConfigFree(c, cfgline)
					}
				}

			} */

			if s5radio.GetActive() && !uempty && !pempty {
				emitConfig(c, "Socks5ProxyUsername", get_edit_text(uedit))
				emitConfig(c, "Socks5ProxyPassword", get_edit_text(pedit))
			}

		}

	})


	scrollbox.SetSizeRequest(550, 350)
	Notebook.AppendPage(scrollbox, nbLabel)

	mainWin.Add(Notebook)
	mainWin.SetResizable(true)
//	mainWin.SetDefaultSize(500, 300)

	mainWin.ShowAll()
	gtk.Main()      // GTK main loop; blocks until gtk.MainQuit() is run. 
}
