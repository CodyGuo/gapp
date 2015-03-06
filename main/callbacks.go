package main

import (
	"fmt"
	"github.com/nvsoft/cef"
	"github.com/nvsoft/gapp/navigator"
	"github.com/nvsoft/gapp/sendinput"
	"github.com/nvsoft/win"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const spec_chars = "!@#$%^&*()"
const base_chars = "1234567890"

func init() {
	cef.RegisterV8CallbackHandler("start", win_start_browser)
	cef.RegisterV8CallbackHandler("emuClick", win_emuClick)
	cef.RegisterV8CallbackHandler("emuInput", win_emuInput)
}

func win_start_browser(browser *cef.Browser, message *cef.CefProcessMessage) interface{} {
	fmt.Printf("win_start_browser\n")
	handle := browser.GetWindowHandle()
	openerHandle := browser.GetOpenerWindowHandle()
	rootHandle := browser.GetRootWindowHandle()
	fmt.Printf("win_start_browser handle=%v openerHandle=%v rootHandle=%v\n", handle, openerHandle, rootHandle)
	win.ShowWindow(win.HWND(rootHandle), win.SW_MAXIMIZE)
	return 0
}

// 模拟点击
func win_emuClick(browser *cef.Browser, message *cef.CefProcessMessage) interface{} {
	fmt.Printf("win_emuClick\n")
	handle := browser.GetWindowHandle()
	openerHandle := browser.GetOpenerWindowHandle()
	rootHandle := browser.GetRootWindowHandle()
	fmt.Printf("win_start_browser handle=%v openerHandle=%v rootHandle=%v\n", handle, openerHandle, rootHandle)
	//win.ShowWindow(win.HWND(rootHandle), win.SW_MAXIMIZE)
	// 模拟点击
	// 查找窗口
	url := message.GetArgumentList().GetString(1)
	fmt.Printf("url=%v\n", url)
	w, ok := windowHolders[url]
	if ok {
		fmt.Printf("找到窗口\n")
		b, o := cef.BrowserByHandle(unsafe.Pointer(w))
		if o {
			x := message.GetArgumentList().GetInt(2)
			y := message.GetArgumentList().GetInt(3)
			buttonType := message.GetArgumentList().GetBool(4)
			fmt.Printf("X=%v Y=%v ButtonType=%v\n", x, y, buttonType)
			rootHandle = b.GetRootWindowHandle()
			//win.ShowWindow(win.HWND(rootHandle), win.SW_MAXIMIZE)
			// 模拟鼠标
			var pt win.POINT
			pt.X = int32(x) // This is your click coordinates
			pt.Y = int32(y)

			hWnd := win.HWND(rootHandle)

			go func() {
				win.SetForegroundWindow(hWnd)
				time.Sleep(3 * time.Second)

				win.ClientToScreen(hWnd, &pt)

				fmt.Printf("ClientToScreen X=%v Y=%v\n", pt.X, pt.Y)

				cx_screen := win.GetSystemMetrics(win.SM_CXSCREEN) //屏幕 宽
				cy_screen := win.GetSystemMetrics(win.SM_CYSCREEN) //     高

				real_x := 65535 * pt.X / cx_screen //转换后的 x
				real_y := 65535 * pt.Y / cy_screen //         y

				var input win.MOUSE_INPUT
				input.Type = win.INPUT_MOUSE
				input.Mi.Dx = real_x
				input.Mi.Dy = real_y
				if buttonType {
					input.Mi.DwFlags = (win.MOUSEEVENTF_ABSOLUTE | win.MOUSEEVENTF_MOVE | win.MOUSEEVENTF_RIGHTDOWN | win.MOUSEEVENTF_RIGHTUP)
				} else {
					input.Mi.DwFlags = (win.MOUSEEVENTF_ABSOLUTE | win.MOUSEEVENTF_MOVE | win.MOUSEEVENTF_LEFTDOWN | win.MOUSEEVENTF_LEFTUP)
				}
				input.Mi.MouseData = 0
				input.Mi.DwExtraInfo = 0
				input.Mi.Time = 0
				win.SendInput(2, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))

				fmt.Printf("点击")
			}()
			//hWnd = win.WindowFromPoint(pt)
		}
	} else {
		fmt.Printf("找不到窗口\n")
	}
	return 0
}

// 模拟输入
func win_emuInput(browser *cef.Browser, message *cef.CefProcessMessage) interface{} {
	rootHandle := browser.GetRootWindowHandle()
	url := message.GetArgumentList().GetString(1)
	fmt.Printf("url=%v\n", url)
	w, ok := windowHolders[url]
	if ok {
		fmt.Printf("找到窗口\n")
		b, o := cef.BrowserByHandle(unsafe.Pointer(w))
		if o {
			inputText := message.GetArgumentList().GetString(2)
			fmt.Printf("inputText=%v\n", inputText)
			rootHandle = b.GetRootWindowHandle()
			//win.ShowWindow(win.HWND(rootHandle), win.SW_MAXIMIZE)
			// 模拟鼠标
			var pt win.POINT
			pt.X = int32(0) // This is your click coordinates
			pt.Y = int32(0)

			hWnd := win.HWND(rootHandle)

			go func() {
				win.SetForegroundWindow(hWnd)
				time.Sleep(5 * time.Second)

				win.ClientToScreen(hWnd, &pt)

				fmt.Printf("ClientToScreen X=%v Y=%v\n", pt.X, pt.Y)

				//sendinput.SendString(inputText)

				ss := strings.Split(inputText, "")
				for i := 0; i < len(ss); i++ {
					c := ss[i]
					cc := []rune(c)[0]
					keyCode := *syscall.StringToUTF16Ptr(c)
					shift := false
					if cc >= 'a' && cc <= 'z' {
						fmt.Printf("小写\n")
						c = strings.ToUpper(c)
						keyCode = *syscall.StringToUTF16Ptr(c)
					}
					if cc >= 'A' && cc <= 'Z' {
						fmt.Printf("大写\n")
						shift = true
					}
					index := strings.Index(spec_chars, c)
					if index >= 0 {
						fmt.Printf("特殊字符1 c=%v index=%v\n", c, index)
						ss := strings.Split(base_chars, "")
						fmt.Printf("ss=%v\n", ss)
						c = ss[index]
						fmt.Printf("特殊字符2 c=%v\n", c)
						keyCode = *syscall.StringToUTF16Ptr(c)
						shift = true
					}
					switch {
					case c == ";":
						keyCode = win.VK_OEM_1
					case c == ":":
						keyCode = win.VK_OEM_1
						shift = true
					case c == "=":
						keyCode = win.VK_OEM_PLUS
					case c == "+":
						keyCode = win.VK_OEM_PLUS
						shift = true
					case c == "-":
						keyCode = win.VK_OEM_MINUS
					case c == "_":
						keyCode = win.VK_OEM_MINUS
						shift = true
					case c == ",":
						keyCode = win.VK_OEM_COMMA
					case c == "<":
						keyCode = win.VK_OEM_COMMA
						shift = true
					case c == ".":
						keyCode = win.VK_OEM_PERIOD
					case c == ">":
						keyCode = win.VK_OEM_PERIOD
						shift = true
					case c == "/":
						keyCode = win.VK_OEM_2
					case c == "?":
						keyCode = win.VK_OEM_2
						shift = true
					case c == "`":
						keyCode = win.VK_OEM_3
					case c == "~":
						keyCode = win.VK_OEM_3
						shift = true
					case c == "[":
						keyCode = win.VK_OEM_4
					case c == "{":
						keyCode = win.VK_OEM_4
						shift = true
					case c == `\`:
						keyCode = win.VK_OEM_5
					case c == "|":
						keyCode = win.VK_OEM_5
						shift = true
					case c == "]":
						keyCode = win.VK_OEM_6
					case c == "}":
						keyCode = win.VK_OEM_6
						shift = true
					case c == "'":
						keyCode = win.VK_OEM_7
					case c == `"`:
						keyCode = win.VK_OEM_7
						shift = true
					}
					fmt.Printf("输入:%v %v\n", keyCode, shift)
					sendinput.KeyPress(keyCode, shift)
				}

				//hWnd = win.WindowFromPoint(pt)

				//PostMessage(edit, WM_KEYDOWN, ...);
				//PostMessage(edit, WM_CHAR, ...);
				//PostMessage(edit, WM_KEYUP, ...);

				//win.SendMessage(hWnd, win.WM_KEYDOWN, 0, 0)

				/*
				   lParam := uintptr(win.MAKELONG(uint16(pt.X), uint16(pt.Y)))
				   win.SendMessage(hWnd, win.WM_LBUTTONDOWN, 0, lParam)    //WM_LBUTTONDOWN
				   //time.Sleep(50 * time.Millisecond)
				   win.SendMessage(hWnd, win.WM_LBUTTONUP, 0, lParam)
				*/
				fmt.Printf("输入")
			}()
			//hWnd = win.WindowFromPoint(pt)
		}
	} else {
		fmt.Printf("找不到窗口\n")
	}
	return 0
}

// http://blog.sina.com.cn/s/blog_648d306d0101gjxh.html
func emuInputText(c string, shift bool) {
	scanKey := *syscall.StringToUTF16Ptr(c)
	var input win.KEYBD_INPUT
	if shift {
		input.Type = win.INPUT_KEYBOARD
		input.Ki.WVk = win.VK_SHIFT
		win.SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))
	}

	input.Type = win.INPUT_KEYBOARD
	input.Ki.WVk = scanKey
	win.SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))

	input.Type = win.INPUT_KEYBOARD
	input.Ki.WVk = scanKey
	input.Ki.DwFlags = win.KEYEVENTF_KEYUP
	win.SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))

	if shift {
		input.Type = win.INPUT_KEYBOARD
		input.Ki.WVk = win.VK_SHIFT
		input.Ki.DwFlags = win.KEYEVENTF_KEYUP
		win.SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))
	}

	win.SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))
}

func emuInput(browser *cef.Browser) {
	fmt.Printf("emuInput\n")
	navigator.InjectKeyPress(browser, "k")
}

func emuInput2(browser *cef.Browser) {
	fmt.Printf("emuInput\n")
	// 查找

	frame := browser.GetMainFrame()

	//frameCount := browser.GetFrameCount()
	//fmt.Printf("frameCount=%v\n", frameCount)

	//frameIdentifiers := browser.GetFrameIdentifiers()
	//fmt.Printf("frameIdentifiers=%v\n", frameIdentifiers)

	// 计算frame坐标
	//mainFrame := browser.GetMainFrame()
	//login_frame_left, login_frame_top := navigator.GetHtmlElementOffset(mainFrame, "#loginIframe")
	//fmt.Printf("#loginIframe坐标:left=%v,top=%v\n", login_frame_left, login_frame_top)
	// 获取frameElement id
	//for i := 0; i < len(frameIdentifiers); i++ {
	//identifier := frameIdentifiers[i]
	//fmt.Printf("frame-identifier=%v\n", identifier)
	// id=loginIframe
	//frame := browser.GetFrameByIdent(identifier)
	//if !frame.IsValid() {
	//	fmt.Printf("IsValid fail.\n")
	//	continue
	//}

	/*c := `
	  var id_ = "";
	  var frame = window.frameElement;  //Get <iframe> element of the window
	  if (frame) {
	  //if (typeof frameElement_.id !== "undefined" && frameElement_.id !== null) {
	      // some code here
	      if (typeof frame.id !== "undefined" && frame.id !== null) {
	          id_ = frame.id;
	      }
	  }
	  app.cefResult(id_);
	  `*/
	//frame_id := frame.ExecuteJavaScriptWithResult(c)
	//fmt.Printf("frame_id=%v\n", frame_id)
	//if frame_id == "loginIframe" {
	fmt.Printf("找到登录界面")

	navigator.InjectFocus(browser, true)
	fmt.Printf("开始登录支付宝...\n")
	fmt.Printf("获取账号输入框按钮坐标\n")

	left, top := navigator.GetHtmlElementOffset(frame, "#J-input-user")
	fmt.Printf("输入框坐标:left=%v,top=%v\n", left, top)

	x := left + 10
	y := top + 10
	fmt.Printf("点击账号输入框 x=%v,y=%v\n", x, y)

	hWnd := win.HWND(browser.GetWindowHandle())
	fmt.Printf("hWnd=%v\n", hWnd)
	navigator.MouseClick(hWnd, x, y)
	fmt.Printf("点击账号输入框\n")
	time.Sleep(2 * time.Second)
	navigator.InjectKeyPress(browser, "kevin.l.zhou@gmail.com")
	fmt.Printf("输入完成.\n")

	time.Sleep(5 * time.Second)

	fmt.Printf("点击密码输入框\n")
	left, top = navigator.GetHtmlElementOffset(frame, "#password_input")
	x = left + 10 + 40
	y = top + 10
	fmt.Printf("输入框坐标:left=%v,top=%v\n", x, y)
	navigator.InjectMouseClick(browser, x, y)
	fmt.Printf("点击账号输入框\n")
	time.Sleep(3 * time.Second)
	navigator.InjectKeyPress(browser, "1")
	fmt.Printf("输入完成.\n")
	time.Sleep(5 * time.Second)
	fmt.Printf("失去焦点1.\n")
	navigator.InjectFocus(browser, false)
	fmt.Printf("失去焦点2.\n")
	//}
	//}
}
