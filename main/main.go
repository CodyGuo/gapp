package main

import (
	"errors"
	"fmt"
	"github.com/nvsoft/cef"
	"github.com/nvsoft/gapp/config"
	"github.com/nvsoft/win"
	"net/http"
	"os"
	"runtime"
	"syscall"
	"unsafe"
	//"time"
)

const (
	ICON_MAIN       = 100
	nguiWindowClass = `\o/ Gapp_Window_Class \o/`
)

var (
	hInstance win.HINSTANCE
	manifest  config.Manifest
)

var wndProc = syscall.NewCallback(WndProc)
var browserSettings = cef.BrowserSettings{}
var windowHolders map[string]win.HWND

func init() {
	hInstance := win.GetModuleHandle(nil)
	if hInstance == 0 {
		panic("GetModuleHandle")
	}

	MustRegisterWindowClass(nguiWindowClass)

	windowHolders = make(map[string]win.HWND)

	manifest.Load()
}

func main() {
	runtime.GOMAXPROCS(4)

	exitCode := cef.ExecuteProcess(unsafe.Pointer(hInstance))
	if exitCode >= 0 {
		return
	}

	// 启动本地Web服务器
	go func() {
		startAppServer()
	}()

	cef.OnContextInitialized = func() {
		window := createWindow()
		createBrowser(window, manifest.FirstPage()) // http://www.baidu.com/
	}

	cef.OnOpenWindow = func(url string) {
		window := createWindow()
		createBrowser(window, url)
		windowHolders[url] = window
	}

	settings := cef.Settings{}
	//settings.SingleProcess = 1                     // 单进程模式
	settings.CachePath = manifest.CachePath()      // Set to empty to disable
	settings.LogSeverity = cef.LOGSEVERITY_DEFAULT // LOGSEVERITY_VERBOSE
	//settings.LocalesDirPath = releasePath + "/locales"
	settings.Locale = manifest.Locale() //"zh-CN"
	settings.BrowserSubprocessPath = manifest.BrowserSubprocessPath()
	//settings.RemoteDebuggingPort = 7000
	cef.Initialize(settings)

	cef.RunMessageLoop()
	cef.Shutdown()
	os.Exit(0)
}

func startAppServer() error {
	staticHandler := http.FileServer(http.Dir(config.Config["static_root"]))
	http.Handle("/", http.StripPrefix("/", staticHandler))

	fmt.Printf("Starting server on port %s\n", config.Config["listen_address"])

	return http.ListenAndServe(config.Config["listen_address"], nil)
}

func createWindow() win.HWND {
	var dwExStyle, dwStyle uint32 = 0, 0

	dwStyle = win.WS_OVERLAPPEDWINDOW

	if !manifest.FormFixed() {
		dwStyle |= win.WS_SIZEBOX
	}

	// 获取屏幕宽度和高度
	var x, y int32
	var width, height int32

	width = manifest.Width()
	height = manifest.Height()
	x = (win.GetSystemMetrics(win.SM_CXSCREEN) - width) / 2
	y = (win.GetSystemMetrics(win.SM_CYSCREEN)-height)/2 - 2

	renderWindow := win.CreateWindowEx(
		dwExStyle,
		syscall.StringToUTF16Ptr(nguiWindowClass),
		nil,
		dwStyle, //|win.WS_CLIPSIBLINGS,
		x,       //win.CW_USEDEFAULT,
		y,       //win.CW_USEDEFAULT,
		width,   //win.CW_USEDEFAULT,
		height,  //win.CW_USEDEFAULT,
		0,       //hwndParent
		0,
		0, //hInstance
		nil)
	if renderWindow == 0 {
		//err := errors.New("CreateWindowEx")
		return win.HWND(0)
	}

	fmt.Printf("CreateWindow x=%v y=%v width=%v height=%v renderWindow=%v renderWindow=%v\n", x, y, width, height, renderWindow, unsafe.Pointer(renderWindow))

	win.ShowWindow(renderWindow, win.SW_SHOW) //win.SW_SHOW
	win.UpdateWindow(renderWindow)

	return renderWindow
}

func createBrowser(renderWindow win.HWND, url string) {
	//browser := cef.CreateBrowser(unsafe.Pointer(hwnd), &browserSettings, url, false)
	//browserSettings := cef.BrowserSettings{}
	cef.CreateBrowser(unsafe.Pointer(renderWindow), &browserSettings, url, false)

	//m_dwStyle = WS_CHILD | WS_CLIPCHILDREN | WS_CLIPSIBLINGS | WS_TABSTOP |
	//		WS_VISIBLE;
	cef.WindowResized(unsafe.Pointer(renderWindow))

	//cef.WindowResized(unsafe.Pointer(renderWindow))

	//cef.WindowResized(unsafe.Pointer(renderWindow))
	// It should be enough to call WindowResized after 10ms,
	// though to be sure let's extend it to 100ms.
	//time.AfterFunc(time.Millisecond*100, func() {
	//cef.WindowResized(unsafe.Pointer(renderWindow))
	//})
}

func WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_CREATE:
		result = win.DefWindowProc(hwnd, msg, wParam, lParam)
	case win.WM_SIZE:
		// 最小化时不能调整Cef窗体，否则恢复时界面一片空白
		//if wParam == win.SIZE_RESTORED || wParam == win.SIZE_MAXIMIZED {
		cef.WindowResized(unsafe.Pointer(hwnd))
		//}
	case win.WM_CLOSE:
		win.DestroyWindow(hwnd)
	case win.WM_DESTROY:
		cef.QuitMessageLoop()
	default:
		result = win.DefWindowProc(hwnd, msg, wParam, lParam)
	}
	return
}

func MustRegisterWindowClass(className string) {
	hInstance := win.GetModuleHandle(nil)
	if hInstance == 0 {
		panic("GetModuleHandle")
	}
	hIcon := win.LoadIcon(hInstance, (*uint16)(unsafe.Pointer(uintptr(ICON_MAIN))))
	//hIcon, _ := NewIconFromResource(hInstance, ICON_MAIN)
	if hIcon == 0 {
		panic("LoadIcon")
	}

	hCursor := win.LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(win.IDC_ARROW))))
	if hCursor == 0 {
		panic("LoadCursor")
	}

	var wc win.WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.LpfnWndProc = wndProc
	wc.HInstance = hInstance
	wc.HIcon = hIcon
	wc.HCursor = hCursor
	wc.HbrBackground = win.COLOR_WINDOW + 1 //COLOR_BTNFACE
	wc.LpszClassName = syscall.StringToUTF16Ptr(className)

	if atom := win.RegisterClassEx(&wc); atom == 0 {
		panic("RegisterClassEx")
	}
}

func NewIconFromResource(instance win.HINSTANCE, resId uint16) (ico win.HICON, err error) {
	if ico = win.LoadIcon(instance, win.MAKEINTRESOURCE(uintptr(resId))); ico == 0 {
		err = errors.New(fmt.Sprintf("Cannot load icon from resource with id %v", resId))
	}

	return ico, err
}
