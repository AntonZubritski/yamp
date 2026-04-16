//go:build darwin

package mpris

/*
#cgo LDFLAGS: -framework CoreGraphics -framework Carbon
#include <CoreGraphics/CoreGraphics.h>
#include <Carbon/Carbon.h>

// Callback bridge — defined in Go, called from C event tap.
extern void mediaKeyCallback(int keyCode);

static CGEventRef tapCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon) {
    if (type == kCGEventTapDisabledByTimeout || type == kCGEventTapDisabledByUserInput) {
        return event;
    }
    int64_t keyCode = CGEventGetIntegerValueField(event, kCGKeyboardEventKeycode);
    // NX_KEYTYPE_PLAY=16, NX_KEYTYPE_NEXT=17, NX_KEYTYPE_PREVIOUS=18, NX_KEYTYPE_FAST=19, NX_KEYTYPE_REWIND=20
    if (type == NX_SYSDEFINED) {
        NSEvent *nsEvent = [NSEvent eventWithCGEvent:event];
        if ([nsEvent subtype] == 8) { // NX_SUBTYPE_AUX_CONTROL_BUTTONS
            int data = [nsEvent data1];
            int keyCode2 = (data & 0xFFFF0000) >> 16;
            int keyFlags = data & 0x0000FFFF;
            int keyState = (keyFlags & 0xFF00) >> 8; // 0xA = down, 0xB = up
            if (keyState == 0x0A) { // key down only
                mediaKeyCallback(keyCode2);
            }
            return NULL; // consume the event
        }
    }
    return event;
}

static void startTap() {
    CGEventMask mask = NX_SYSDEFINEDMASK;
    CFMachPortRef tap = CGEventTapCreate(kCGSessionEventTap, kCGHeadInsertEventTap, kCGEventTapOptionDefault, mask, tapCallback, NULL);
    if (!tap) return;
    CFRunLoopSourceRef source = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, tap, 0);
    CFRunLoopAddSource(CFRunLoopGetCurrent(), source, kCFRunLoopCommonModes);
    CGEventTapEnable(tap, true);
    CFRunLoopRun();
}
*/
import "C"

var darwinSend func(interface{})

//export mediaKeyCallback
func mediaKeyCallback(keyCode C.int) {
	if darwinSend == nil {
		return
	}
	switch keyCode {
	case 16: // NX_KEYTYPE_PLAY
		darwinSend(PlayPauseMsg{})
	case 17: // NX_KEYTYPE_NEXT
		darwinSend(NextMsg{})
	case 18: // NX_KEYTYPE_PREVIOUS
		darwinSend(PrevMsg{})
	}
}

// StartMediaKeyListener starts listening for macOS media keys.
// Blocks — call in a goroutine.
func StartMediaKeyListener(send func(interface{})) {
	darwinSend = send
	C.startTap()
}
