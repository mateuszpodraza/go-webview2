// +build windows

package webview2

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"github.com/mattpodraza/webview2/pkg/user32"
	"golang.org/x/sys/windows"
)

var (
	errSuccess = syscall.Errno(0)
)

// Hint is used to configure window sizing and resizing behavior.
type Hint int

const (
	// HintNone specifies that width and height are default size
	HintNone Hint = iota

	// HintFixed specifies that window size can not be changed by a user
	HintFixed

	// HintMin specifies that width and height are minimum bounds
	HintMin

	// HintMax specifies that width and height are maximum bounds
	HintMax
)

var (
	ole32               = windows.NewLazySystemDLL("ole32")
	ole32CoInitializeEx = ole32.NewProc("CoInitializeEx")

	windowContext     = map[uintptr]interface{}{}
	windowContextSync sync.RWMutex
)

type WebView struct {
	config *config
	edge   *edge
	hwnd   windows.Handle
}

type config struct {
	devtoolsDisabled             bool
	defaultContextMenusDisabled  bool
	defaultScriptDialogsDisabled bool
	hostObjectsDisallowed        bool
	builtInErrorPageDisabled     bool
	scriptDisabled               bool
	statusBarDisabled            bool
	webMessageDisabled           bool
	zoomControlDisabled          bool

	minSize user32.Point
	maxSize user32.Point
}

type Option func(*WebView) error

func New(options ...Option) (*WebView, error) {
	wv := &WebView{
		config: &config{},
		edge:   newEdge(),
	}

	if err := wv.create(); err != nil {
		return nil, fmt.Errorf("failed to create the WebView: %w", err)
	}

	for _, option := range options {
		if err := option(wv); err != nil {
			return nil, fmt.Errorf("failed to apply an option: %w", err)
		}
	}

	return wv, nil
}

func (wv *WebView) create() error {
	var hinstance windows.Handle

	err := windows.GetModuleHandleEx(0, nil, &hinstance)
	if err != nil {
		return fmt.Errorf("failed to get the module handle: %w", err)
	}

	icow, err := user32.GetSystemMetrics(user32.SystemMetricsCxIcon)
	if err != nil {
		return err
	}

	icoh, err := user32.GetSystemMetrics(user32.SystemMetricsCyIcon)
	if err != nil {
		return err
	}

	icon, err := user32.LoadImageW(hinstance, icow, icoh)
	if err != nil {
		return err
	}

	wc := user32.WndClassExW{
		CBSize:        uint32(unsafe.Sizeof(user32.WndClassExW{})),
		HInstance:     hinstance,
		LpszClassName: windows.StringToUTF16Ptr("webview"),
		HIcon:         icon,
		HIconSm:       icon,
		LpfnWndProc:   windows.NewCallback(wndproc),
	}

	err = user32.RegisterClassExW(&wc)
	if err != nil {
		return err
	}

	wv.hwnd, err = user32.CreateWindowExW(
		"webview",
		"",
		user32.CW_USEDEFAULT,
		user32.CW_USEDEFAULT,
		640,
		480,
		hinstance,
	)

	if err != nil {
		return fmt.Errorf("failed to create the window: %w", err)
	}

	setWindowContext(wv.hwnd, wv)

	err = user32.ShowWindow(wv.hwnd)
	if err != nil {
		return fmt.Errorf("failed to update the window: %w", err)
	}

	err = user32.SetFocus(wv.hwnd)
	if err != nil {
		return fmt.Errorf("failed to set focus: %w", err)
	}

	err = wv.edge.Embed(wv.hwnd)
	if err != nil {
		return fmt.Errorf("failed to embed the browser: %w", err)
	}

	err = wv.edge.Resize()
	if err != nil {
		return fmt.Errorf("failed to resize the browser: %w", err)
	}

	return nil
}

func (wv *WebView) SetTitle(title string) error {
	return user32.SetWindowTextW(wv.hwnd, title)
}

func (wv *WebView) SetSize(width, height int, hints Hint) error {
	index := user32.GWLStyle

	style, err := user32.GetWindowLongPtrW(wv.hwnd, uintptr(index))
	if err != nil {
		return fmt.Errorf("failed to call user32GetWindowLongPtrW: %w", err)
	}

	if hints == HintFixed {
		style &^= (user32.WSThickFrame | user32.WSMaximizeBox)
	} else {
		style |= (user32.WSThickFrame | user32.WSMaximizeBox)
	}

	err = user32.SetWindowLongPtrW(wv.hwnd, uintptr(index), style)
	if err != nil {
		return fmt.Errorf("failed to call user32SetWindowLongPtrW: %w", err)
	}

	if hints == HintMax {
		wv.config.maxSize.X = int32(width)
		wv.config.maxSize.Y = int32(height)
		return nil
	}

	if hints == HintMin {
		wv.config.minSize.X = int32(width)
		wv.config.minSize.Y = int32(height)
		return nil
	}

	rect := user32.Rect{
		Left:   0,
		Top:    0,
		Right:  int32(width),
		Bottom: int32(height),
	}

	err = user32.AdjustWindowRec(&rect, user32.WSOverlappedWindow, false)
	if err != nil {
		return fmt.Errorf("failed to adjust window rect: %w", err)
	}

	err = user32.SetWindowPos(
		wv.hwnd,
		rect.Left,
		rect.Top,
		rect.Right-rect.Left,
		rect.Bottom-rect.Top,
		user32.SWPNoZOrder|user32.SWPNoActivate|user32.SWPNoMove|user32.SWPFrameChanged,
	)
	if err != nil {
		return fmt.Errorf("failed to set the window position: %w", err)
	}

	return wv.edge.Resize()
}

func (wv *WebView) Navigate(url string) error {
	return wv.edge.Navigate(url)
}

func (wv *WebView) Run() error {
	for {
		msg, err := user32.GetMessageW()
		if err != nil {
			return fmt.Errorf("failed to get message: %w", err)
		}

		if msg.Message == user32.WMQuit {
			return nil
		}

		err = user32.TranslateMessage(msg)
		if err != nil {
			return fmt.Errorf("failed to translate message: %w", err)
		}

		err = user32.DispatchMessageW(msg)
		if err != nil {
			return fmt.Errorf("failed to dispatch message: %w", err)
		}
	}
}

func (wv *WebView) Terminate() error {
	return user32.PostQuitMessage(0)
}

func (wv *WebView) Window() windows.Handle {
	return wv.hwnd
}

func (wv *WebView) Init(js string) error {
	return wv.edge.Init(js)
}

func (wv *WebView) Eval(js string) error {
	return wv.edge.Eval(js)
}

func getWindowContext(wnd uintptr) interface{} {
	windowContextSync.RLock()
	defer windowContextSync.RUnlock()
	return windowContext[wnd]
}

func setWindowContext(wnd windows.Handle, data interface{}) {
	windowContextSync.Lock()
	defer windowContextSync.Unlock()
	windowContext[uintptr(wnd)] = data
}

func init() {
	runtime.LockOSThread()

	_, _, err := ole32CoInitializeEx.Call(0, 2)
	if err != nil && !errors.Is(err, errSuccess) {
		log.Printf("warning: CoInitializeEx call failed: %v", err)
	}
}

func wndproc(hwnd, msg, wp, lp uintptr) uintptr {
	if wv, ok := getWindowContext(hwnd).(*WebView); ok {
		switch msg {
		case user32.WMSize:
			_ = wv.edge.Resize()
		case user32.WMClose:
			_ = user32.DestroyWindow(windows.Handle(hwnd))
		case user32.WMDestroy:
			_ = wv.Terminate()
		case user32.WMGetMinMaxInfo:
			lpmmi := (*user32.MinMaxInfo)(unsafe.Pointer(lp))

			if wv.config.maxSize.X > 0 && wv.config.maxSize.Y > 0 {
				lpmmi.MaxSize = wv.config.maxSize
				lpmmi.MaxTrackSize = wv.config.maxSize
			}

			if wv.config.minSize.X > 0 && wv.config.minSize.Y > 0 {
				lpmmi.MinTrackSize = wv.config.minSize
			}
		default:
			r, _ := user32.DefWindowProcW(hwnd, msg, wp, lp)
			return r
		}
		return 0
	}

	r, _ := user32.DefWindowProcW(hwnd, msg, wp, lp)
	return r
}
