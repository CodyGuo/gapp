package navigator

import (
	"fmt"
	"github.com/nvsoft/cef"
	"github.com/nvsoft/win"
	"strconv"
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

// http://stackoverflow.com/questions/3720968/win32-simulate-a-click-without-simulating-mouse-movement
func MouseClick(hWnd win.HWND, x, y int) {
	var pt win.POINT
	pt.X = int32(x) // This is your click coordinates
	pt.Y = int32(y)

	//hWnd = win.WindowFromPoint(pt)
	lParam := uintptr(win.MAKELONG(uint16(pt.X), uint16(pt.Y)))
	win.PostMessage(hWnd, win.WM_LBUTTONDOWN, win.MK_LBUTTON, lParam)
	win.PostMessage(hWnd, win.WM_LBUTTONUP, win.MK_LBUTTON, lParam)
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

// 获取Html原始坐标
func GetHtmlElementOffset(frame *cef.CefFrame, selector string) (left, top int) {
	js := `
    function getOffset( el ) {
        var _x = 0;
        var _y = 0;
        while( el && !isNaN( el.offsetLeft ) && !isNaN( el.offsetTop ) ) {
            _x += el.offsetLeft - el.scrollLeft;
            _y += el.offsetTop - el.scrollTop;
            // chrome/safari
            //if ($.browser.webkit) {
                el = el.parentNode;
            //} else {
                // firefox/IE
                //el = el.offsetParent;
            //}
        }
        return { left: _x, top: _y };
    }
    var e = document.querySelector("` + selector + `");
    var offset = getOffset(e);
    cef.setResult(offset.left + "," + offset.top);
    `
	strOffset := frame.ExecuteJavaScriptWithResult(js)

	ss := strings.Split(strOffset, ",")

	if len(ss) == 2 {
		left, _ = strconv.Atoi(ss[0])
		top, _ = strconv.Atoi(ss[1])
	}
	return
}
