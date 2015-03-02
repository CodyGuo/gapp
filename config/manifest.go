package config

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const (
	manifest_filename       = `manifest.json`
	first_page              = `first_page`
	application_title       = `application_title`
	locale                  = `locale`
	cache_path              = `cache_path`
	style                   = `style`
	width                   = `width`
	height                  = `height`
	form_fixed              = `form_fixed`
	enable_transparent      = `enable_transparent`
	browser_subprocess_path = `browser_subprocess_path`

	WindowStyleNormal      = 0
	WindowStyleCaptionLess = 1
)

var (
	kernel                = syscall.MustLoadDLL("kernel32.dll")
	getModuleFileNameProc = kernel.MustFindProc("GetModuleFileNameW")
)

type Manifest struct {
	manifest revManifest
}

type revManifest map[string]interface{}

func (a *Manifest) BrowserSubprocessPath() string {
	return a.Get(browser_subprocess_path).(string)
}

func (a *Manifest) FirstPage() string {
	return a.Get(first_page).(string)
}

func (a *Manifest) ApplicationTitle() string {
	return a.Get(application_title).(string)
}

func (a *Manifest) Locale() string {
	return a.Get(locale).(string)
}

func (a *Manifest) CachePath() string {
	return a.Get(cache_path).(string)
}

func (a *Manifest) FormFixed() bool {
	return a.Get(form_fixed).(bool)
}

func (a *Manifest) EnableTransparent() bool {
	return a.Get(enable_transparent).(bool)
}

func (a *Manifest) Style() int32 {
	return int32(a.Get(style).(float64))
}

func (a *Manifest) Width() int32 {
	return int32(a.Get(width).(float64))
}

func (a *Manifest) Height() int32 {
	return int32(a.Get(height).(float64))
}

func (a *Manifest) Get(key string) interface{} {
	v, _ := a.manifest[key]
	return v
}

func (a *Manifest) Load() *Manifest {
	manifestPath := path.Join(a.Path(), manifest_filename)
	//fmt.Printf("manifestPath=%v\n", manifestPath)
	data, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		panic("Load Manifest")
	}
	//fmt.Println(err)
	json.Unmarshal(data, &a.manifest)
	return a
}

func (a Manifest) Path() string {
	return ExePath()
}

func ExePath() string {
	exePath, _ := Executable()
	exeDir := filepath.Dir(exePath)
	return exeDir
}

// GetModuleFileName() with hModule = NULL
func Executable() (exePath string, err error) {
	return getModuleFileName()
}

func getModuleFileName() (string, error) {
	var n uint32
	b := make([]uint16, syscall.MAX_PATH)
	size := uint32(len(b))

	r0, _, e1 := getModuleFileNameProc.Call(0, uintptr(unsafe.Pointer(&b[0])), uintptr(size))
	n = uint32(r0)
	if n == 0 {
		return "", e1
	}
	s := string(utf16.Decode(b[0:n]))
	s = strings.Replace(s, "\\", "/", -1)
	return s, nil
}
