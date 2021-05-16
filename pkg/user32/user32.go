// Package user32 wraps some calls from the user32.dll to make them slightly easier to use.
package user32

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	CW_USEDEFAULT = 0x80000000

	SystemMetricsCxScreen = 0
	SystemMetricsCyScreen = 1
	SystemMetricsCxIcon   = 11
	SystemMetricsCyIcon   = 12

	GWLStyle = -16

	WSOverlapped       = 0x00000000
	WSMaximizeBox      = 0x00020000
	WSThickFrame       = 0x00040000
	WSCaption          = 0x00C00000
	WSSysMenu          = 0x00080000
	WSMinimizeBox      = 0x00020000
	WSOverlappedWindow = (WSOverlapped | WSCaption | WSSysMenu | WSThickFrame | WSMinimizeBox | WSMaximizeBox)

	SWPNoZOrder     = 0x0004
	SWPNoActivate   = 0x0010
	SWPNoSize       = 0x0001
	SWPNoMove       = 0x0002
	SWPFrameChanged = 0x0020

	WMDestroy       = 0x0002
	WMSize          = 0x0005
	WMClose         = 0x0010
	WMQuit          = 0x0012
	WMGetMinMaxInfo = 0x0024
	WMApp           = 0x8000
)

const (
	SW_HIDE = iota
	SW_SHOWNORMAL
	SW_SHOWMINIMIZED
	SW_SHOWMAXIMIZED
	SW_SHOWNOACTIVATE
	SW_SHOW
	SW_MINIMIZE
	SW_SHOWMINNOACTIVE
	SW_SHOWNA
	SW_RESTORE
	SW_SHOWDEFAULT
	SW_FORCEMINIMIZE
)

var (
	errOK = syscall.Errno(0)

	user32 = windows.NewLazySystemDLL("user32")

	loadImageW        = user32.NewProc("LoadImageW")
	getSystemMetrics  = user32.NewProc("GetSystemMetrics")
	registerClassExW  = user32.NewProc("RegisterClassExW")
	createWindowExW   = user32.NewProc("CreateWindowExW")
	destroyWindow     = user32.NewProc("DestroyWindow")
	showWindow        = user32.NewProc("ShowWindow")
	setFocus          = user32.NewProc("SetFocus")
	getMessageW       = user32.NewProc("GetMessageW")
	translateMessage  = user32.NewProc("TranslateMessage")
	dispatchMessageW  = user32.NewProc("DispatchMessageW")
	defWindowProcW    = user32.NewProc("DefWindowProcW")
	getClientRect     = user32.NewProc("GetClientRect")
	postQuitMessage   = user32.NewProc("PostQuitMessage")
	setWindowTextW    = user32.NewProc("SetWindowTextW")
	getWindowLongPtrW = user32.NewProc("GetWindowLongPtrW")
	setWindowLongPtrW = user32.NewProc("SetWindowLongPtrW")
	adjustWindowRect  = user32.NewProc("AdjustWindowRect")
	setWindowPos      = user32.NewProc("SetWindowPos")
)

type Msg struct {
	HWND     syscall.Handle
	Message  uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       Point
	LPrivate uint32
}

type Point struct {
	X, Y int32
}

type WndClassExW struct {
	CBSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CnClsExtra    int32
	CBWndExtra    int32
	HInstance     windows.Handle
	HIcon         windows.Handle
	HCursor       windows.Handle
	HbrBackground windows.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       windows.Handle
}

type Rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}
type MinMaxInfo struct {
	Reserved     Point
	MaxSize      Point
	MaxPosition  Point
	MinTrackSize Point
	MaxTrackSize Point
}

func GetMessageW() (*Msg, error) {
	var msg Msg

	r, _, _ := getMessageW.Call(
		uintptr(unsafe.Pointer(&msg)),
		0,
		0,
		0,
	)

	if int32(r) == -1 {
		return nil, windows.GetLastError()
	}

	return &msg, nil
}

func TranslateMessage(msg *Msg) error {
	_, _, err := translateMessage.Call(uintptr(unsafe.Pointer(msg)))
	if err != nil && !errors.Is(err, errOK) {
		return err
	}

	return nil
}

func DispatchMessageW(msg *Msg) error {
	_, _, err := dispatchMessageW.Call(uintptr(unsafe.Pointer(msg)))
	if err != nil && !errors.Is(err, errOK) {
		return err
	}

	return nil
}

func PostQuitMessage(exitCode int) error {
	_, _, err := postQuitMessage.Call(uintptr(exitCode))
	if err != nil && !errors.Is(err, errOK) {
		return err
	}

	return nil
}

func DestroyWindow(hwnd windows.Handle) error {
	_, _, err := destroyWindow.Call(uintptr(hwnd))
	if err != nil && !errors.Is(err, errOK) {
		return err
	}

	return nil
}

func DefWindowProcW(hwnd, msg, wp, lp uintptr) (uintptr, error) {
	r, _, err := defWindowProcW.Call(hwnd, msg, wp, lp)
	if err != nil && !errors.Is(err, errOK) {
		return 0, err
	}

	return r, nil
}

func RegisterClassExW(wc *WndClassExW) error {
	_, _, err := registerClassExW.Call(uintptr(unsafe.Pointer(wc)))
	if err != nil && !errors.Is(err, errOK) {
		return err
	}

	return nil
}

func CreateWindowExW(className, windowName string, x, y, width, height int, hInstance windows.Handle) (windows.Handle, error) {
	class, err := windows.UTF16PtrFromString(className)
	if err != nil {
		return 0, fmt.Errorf("invalid className: %w", err)
	}

	window, err := windows.UTF16PtrFromString(windowName)
	if err != nil {
		return 0, fmt.Errorf("invalid windowName: %w", err)
	}

	hwndptr, _, err := createWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(class)),
		uintptr(unsafe.Pointer(window)),
		0xCF0000, // WS_OVERLAPPEDWINDOW
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		0,
		0,
		uintptr(hInstance),
		0,
	)

	if err != nil && !errors.Is(err, errOK) {
		return 0, err
	}

	return windows.Handle(hwndptr), nil
}

func ShowWindow(hwnd windows.Handle, cmdShow int) error {
	_, _, err := showWindow.Call(uintptr(hwnd), uintptr(cmdShow))
	if err != nil && !errors.Is(err, errOK) {
		return err
	}

	return nil
}

func SetFocus(hwnd windows.Handle) error {
	_, _, err := setFocus.Call(uintptr(hwnd))
	if err != nil && !errors.Is(err, errOK) {
		return err
	}

	return nil
}

func GetSystemMetrics(metric uintptr) (uintptr, error) {
	hwnd, _, err := getSystemMetrics.Call(metric)
	if err != nil && !errors.Is(err, errOK) {
		return 0, err
	}

	return hwnd, nil
}

func LoadImageW(hInstance windows.Handle, cx, cy uintptr) (windows.Handle, error) {
	hwnd, _, err := loadImageW.Call(uintptr(hInstance), 32512, cx, cy, 0)
	if err != nil && !errors.Is(err, errOK) {
		return 0, err
	}

	return windows.Handle(hwnd), nil
}

func SetWindowTextW(hwnd windows.Handle, text string) error {
	tptr, err := windows.UTF16PtrFromString(text)
	if err != nil {
		return fmt.Errorf("invalid text: %w", err)
	}

	_, _, err = setWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(tptr)))
	if err != nil && !errors.Is(err, errOK) {
		return err
	}

	return nil
}

func GetWindowLongPtrW(hwnd windows.Handle, nIndex uintptr) (uintptr, error) {
	info, _, err := getWindowLongPtrW.Call(uintptr(hwnd), nIndex)
	if err != nil && !errors.Is(err, errOK) {
		return 0, err
	}

	return info, nil
}

func SetWindowLongPtrW(hwnd windows.Handle, nIndex uintptr, newLong uintptr) error {
	_, _, err := setWindowLongPtrW.Call(uintptr(hwnd), nIndex, newLong)
	if err != nil && !errors.Is(err, errOK) {
		return err
	}

	return nil
}

func AdjustWindowRec(rect *Rect, style uintptr, hasMenu bool) error {
	var hm uintptr
	if hasMenu {
		hm = 1
	}

	_, _, err := adjustWindowRect.Call(uintptr(unsafe.Pointer(rect)), style, hm)
	if err != nil && !errors.Is(err, errOK) {
		return err
	}

	return nil
}

func SetWindowPos(hwnd windows.Handle, x, y, cx, cy int32, flags uintptr) error {
	_, _, err := setWindowPos.Call(
		uintptr(hwnd),
		0,
		uintptr(x),
		uintptr(y),
		uintptr(cx),
		uintptr(cy),
		flags,
	)
	if err != nil && !errors.Is(err, errOK) {
		return err
	}

	return nil
}

func GetClientRect(hwnd windows.Handle) (*Rect, error) {
	var rect Rect
	_, _, err := getClientRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&rect)))
	if err != nil && !errors.Is(err, errOK) {
		return nil, err
	}

	return &rect, nil
}
