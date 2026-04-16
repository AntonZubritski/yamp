//go:build windows

package mpris

import (
	"syscall"
	"unsafe"
)

var (
	user32              = syscall.NewLazyDLL("user32.dll")
	procRegisterHotKey  = user32.NewProc("RegisterHotKey")
	procGetMessage      = user32.NewProc("GetMessageW")
)

const (
	modNoRepeat = 0x4000
	wmHotkey    = 0x0312

	// Virtual key codes for media keys.
	vkMediaPlayPause = 0xB3
	vkMediaNextTrack = 0xB0
	vkMediaPrevTrack = 0xB1
	vkMediaStop      = 0xB2
)

type msg struct {
	hwnd    uintptr
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      [2]int32
}

// StartMediaKeyListener registers global media hotkeys and listens for them.
// It calls send() with the appropriate message type when a media key is pressed.
// This function blocks and should be called in a goroutine.
func StartMediaKeyListener(send func(interface{})) {
	// Register media keys as global hotkeys.
	// ID 1=PlayPause, 2=Next, 3=Prev, 4=Stop
	procRegisterHotKey.Call(0, 1, uintptr(modNoRepeat), uintptr(vkMediaPlayPause))
	procRegisterHotKey.Call(0, 2, uintptr(modNoRepeat), uintptr(vkMediaNextTrack))
	procRegisterHotKey.Call(0, 3, uintptr(modNoRepeat), uintptr(vkMediaPrevTrack))
	procRegisterHotKey.Call(0, 4, uintptr(modNoRepeat), uintptr(vkMediaStop))

	var m msg
	for {
		ret, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if ret == 0 {
			break
		}
		if m.message == wmHotkey {
			switch m.wParam {
			case 1:
				send(PlayPauseMsg{})
			case 2:
				send(NextMsg{})
			case 3:
				send(PrevMsg{})
			case 4:
				send(StopMsg{})
			}
		}
	}
}
