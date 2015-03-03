package navigator

import (
	"fmt"
	"github.com/nvsoft/cef"
	"github.com/nvsoft/win"
	"strings"
	"syscall"
)

// 模拟按键功能
// 参考: https://github.com/svn2github/cef3/blob/master/tests/unittests/os_rendering_unittest.cc

// 设置焦点
func InjectFocus(b *cef.Browser, focus bool) {
	h := b.GetHost()
	h.SetFocus(focus)
}

// 模拟鼠标点击
func InjectMouseClick(b *cef.Browser, x, y /*, buttonType int, mouseUp bool, clickCount*/ int) {
	var event cef.CefMouseEvent
	event.Modifiers = 0
	event.X = x
	event.Y = y
	h := b.GetHost()
	h.SendMouseClickEvent(&event, 0, false, 1)
	h.SendMouseClickEvent(&event, 0, true, 1)
}

// 模拟输入字符串
func InjectKeyPress(b *cef.Browser, text string) {
	h := b.GetHost()
	ss := strings.Split(text, "")
	for i := 0; i < len(ss); i++ {
		c := ss[i]
		InjectKey(h, c)
	}
}

// 模拟按键
func InjectKey(h *cef.BrowserHost, key string) {
	if len(key) != 1 {
		fmt.Printf("Argement error.")
		return
	}

	var event cef.CefKeyEvent
	var scanKey uint16
	var vkCode byte
	var scanCode uint32

	scanKey = *syscall.StringToUTF16Ptr(key)
	vkCode = win.LOBYTE(uint16(win.VkKeyScan(scanKey)))
	scanCode = win.MapVirtualKey(uint32(vkCode), win.MAPVK_VK_TO_VSC)
	//fmt.Printf("vkCode=%v, scanCode=%v\n", VkCode, scanCode)

	event.IsSystemKey = false
	event.Modifiers = 0
	event.NativeKeyCode = int((scanCode << 16) | // key scan code
		1) // key repeat count
	event.WindowsKeyCode = int(vkCode)

	event.Type = cef.KEYEVENT_RAWKEYDOWN
	h.SendKeyEvent(&event)

	event.WindowsKeyCode = int(scanKey)
	event.Type = cef.KEYEVENT_CHAR
	h.SendKeyEvent(&event)

	event.WindowsKeyCode = int(vkCode)
	// bits 30 and 31 should be always 1 for WM_KEYUP
	event.NativeKeyCode = int(uint32(event.NativeKeyCode) | 0xC0000000)
	event.Type = cef.KEYEVENT_KEYUP
	h.SendKeyEvent(&event)
}
