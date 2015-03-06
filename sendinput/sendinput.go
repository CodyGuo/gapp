package sendinput

import (
	"fmt"
	"github.com/nvsoft/win"
	"strings"
	//"syscall"
	"errors"
	"strconv"
	"syscall"
	"time"
	"unsafe"
)

// http://blog.csdn.net/fangkailove/article/details/7614492

var xToVKeyMap map[uint32]uint32

func init() {
	xToVKeyMap = make(map[uint32]uint32)
	FillMap(xToVKeyMap)
}

func getVKey(keyCode uint32) (uint32, error) {
	res, ok := xToVKeyMap[keyCode]
	if !ok {
		err := errors.New("Unexpected keycode: " + strconv.Itoa(int(keyCode)))
		return 0, err
	}
	return res, nil
}

func FillMap(m map[uint32]uint32) {
	m[8] = win.VK_BACK     // backspace
	m[9] = win.VK_TAB      // tab
	m[13] = win.VK_RETURN  // enter
	m[16] = win.VK_SHIFT   // shift
	m[17] = win.VK_CONTROL // ctrl
	m[18] = win.VK_MENU    // alt key
	m[19] = win.VK_PAUSE
	m[20] = win.VK_CAPITAL // caps lock
	m[27] = win.VK_ESCAPE
	m[32] = win.VK_SPACE //space bar
	m[33] = win.VK_PRIOR //page up
	m[34] = win.VK_NEXT  //page down
	m[35] = win.VK_END
	m[36] = win.VK_HOME
	m[37] = win.VK_LEFT //arrows
	m[38] = win.VK_UP
	m[39] = win.VK_RIGHT
	m[40] = win.VK_DOWN
	m[45] = win.VK_INSERT
	m[46] = win.VK_DELETE
	m[48] = 0x30 // digits 0-9
	m[49] = 0x31
	m[50] = 0x32
	m[51] = 0x33
	m[52] = 0x34
	m[53] = 0x35
	m[54] = 0x36
	m[55] = 0x37
	m[56] = 0x38
	m[57] = 0x39 // 9
	m[65] = 0x41 // letters a-z
	m[66] = 0x42
	m[67] = 0x43
	m[68] = 0x44
	m[69] = 0x45
	m[70] = 0x46
	m[71] = 0x47
	m[72] = 0x48
	m[73] = 0x49
	m[74] = 0x4A
	m[75] = 0x4b
	m[76] = 0x4c
	m[77] = 0x4d
	m[78] = 0x4e
	m[79] = 0x4f
	m[80] = 0x50
	m[81] = 0x51
	m[82] = 0x52
	m[83] = 0x53
	m[84] = 0x54
	m[85] = 0x55
	m[86] = 0x56
	m[87] = 0x57
	m[88] = 0x58
	m[89] = 0x59
	m[90] = 0x5a        // z
	m[91] = win.VK_LWIN // left window key
	m[92] = win.VK_RWIN
	m[93] = win.VK_SELECT
	m[96] = win.VK_NUMPAD0
	m[97] = win.VK_NUMPAD1
	m[98] = win.VK_NUMPAD2
	m[99] = win.VK_NUMPAD3
	m[100] = win.VK_NUMPAD4
	m[101] = win.VK_NUMPAD5
	m[102] = win.VK_NUMPAD6
	m[103] = win.VK_NUMPAD7
	m[104] = win.VK_NUMPAD8
	m[105] = win.VK_NUMPAD9
	m[106] = win.VK_MULTIPLY // numpad operations
	m[107] = win.VK_ADD
	m[109] = win.VK_SUBTRACT
	m[110] = win.VK_DECIMAL // decimal point
	m[111] = win.VK_DIVIDE
	m[112] = win.VK_F1
	m[113] = win.VK_F2
	m[114] = win.VK_F3
	m[115] = win.VK_F4
	m[116] = win.VK_F5
	m[117] = win.VK_F6
	m[118] = win.VK_F7
	m[119] = win.VK_F8
	m[120] = win.VK_F9
	m[121] = win.VK_F10
	m[122] = win.VK_F11
	m[123] = win.VK_F12
	m[144] = win.VK_NUMLOCK
	m[145] = win.VK_SCROLL
	// tuka stava mnogo strashno
	m[186] = win.VK_OEM_1    // semicolon
	m[187] = win.VK_OEM_PLUS // equal sign or plus
	m[188] = win.VK_OEM_COMMA
	m[189] = win.VK_OEM_MINUS
	m[190] = win.VK_OEM_PERIOD
	m[191] = win.VK_OEM_2 // forward slash ('/')
	m[192] = win.VK_OEM_3 // grave accent key ('~')
	m[219] = win.VK_OEM_4 // open bracket ('[')
	m[220] = win.VK_OEM_5 // backslach ('\')
	m[221] = win.VK_OEM_6 // close bracket (']')
	m[222] = win.VK_OEM_7 // single quote (''' lol)
}

func send_key_event(keyCode uint16, isKeyUp bool) {
	var input win.KEYBD_INPUT
	input.Type = win.INPUT_KEYBOARD
	input.Ki.WVk = keyCode
	if isKeyUp {
		input.Ki.DwFlags = win.KEYEVENTF_KEYUP
	}
	win.SendInput(1, unsafe.Pointer(&input), int32(unsafe.Sizeof(input)))
	time.Sleep(3 * time.Millisecond)
}

func KeyDown(keyCode uint16) {
	send_key_event(keyCode, false)
}

func KeyUp(keyCode uint16) {
	send_key_event(keyCode, true)
}

// char in 1~255 key press
func KeyPress(keyCode uint16, isShift bool) {
	fmt.Printf("KeyPress keyCode=%v isShift=%v\n", keyCode, isShift)
	if isShift {
		send_key_event(win.VK_SHIFT, false)
	}
	send_key_event(keyCode, false)
	send_key_event(keyCode, true)
	if isShift {
		send_key_event(win.VK_SHIFT, true)
	}
}

// unicode char key press
func UniKeyPress(keyCode uint16) {
	var input = [2]win.KEYBD_INPUT{}

	input[0].Type = win.INPUT_KEYBOARD
	input[0].Ki.WVk = 0
	input[0].Ki.WScan = keyCode
	input[0].Ki.DwFlags = win.KEYEVENTF_UNICODE

	input[1].Type = win.INPUT_KEYBOARD
	input[1].Ki.WVk = 0
	input[1].Ki.WScan = keyCode
	input[1].Ki.DwFlags = win.KEYEVENTF_UNICODE | win.KEYEVENTF_KEYUP

	win.SendInput(2, unsafe.Pointer(&input[0]), int32(unsafe.Sizeof(input[0])))
	time.Sleep(5 * time.Millisecond)
}

func SendString(keys string) {
	cs := strings.Split(keys, "")
	for i := 0; i < len(cs); i++ {
		c := cs[i]
		cc := *syscall.StringToUTF16Ptr(c)
		cr := []rune(c)[0]
		//cc := uint16(cr)
		fmt.Printf("KEY c=%v cc=%v\n", c, cc)
		if cc >= 0 && cc < 256 {
			vk := win.VkKeyScan(cc)
			//vk, _ := getVKey(uint32(cc))
			fmt.Printf("vk=%v\n", vk)
			if vk == -1 {
				UniKeyPress(cc)
			} else {
				if vk < 0 {
					vk = ^vk + 0x1
				}
				fmt.Printf("vk=%v\n", vk)
				shift := (vk>>8&0x1 == 0x1)
				ks := win.GetKeyState(win.VK_CAPITAL)
				fmt.Printf("shift=%v ww=%v\n", shift, ks)
				if win.GetKeyState(win.VK_CAPITAL)&0x1 == 0x0 {
					if (cr >= 'a' && cr <= 'z') || (cr >= 'A' && cr <= 'Z') {
						//shift = !shift
					}
				}
				KeyPress(uint16(vk&0xFF), shift)
			}
		}
	}
}
