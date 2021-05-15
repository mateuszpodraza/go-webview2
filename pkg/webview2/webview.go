package webview2

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/jchv/go-winloader"
	"github.com/mattpodraza/webview2/v2/pkg/user32"
	"github.com/mattpodraza/webview2/v2/pkg/webviewloader"
	"golang.org/x/sys/windows"
)

var (
	ole32               = windows.NewLazySystemDLL("ole32")
	ole32CoInitializeEx = ole32.NewProc("CoInitializeEx")

	errOK = syscall.Errno(0)
)

func init() {
	runtime.LockOSThread()

	_, _, err := ole32CoInitializeEx.Call(0, 2)
	if err != nil && !errors.Is(err, errOK) {
		log.Printf("warning: CoInitializeEx call failed: %v", err)
	}
}

type WebView struct {
	dll winloader.Proc

	window  *window
	browser *browser
}

func New(options ...Option) (*WebView, error) {
	wv := &WebView{
		window: &window{
			config: &windowConfig{
				width:  640,
				height: 480,
				title:  "Webview",
			},
		},
		browser: &browser{
			config: &browserConfig{
				initialURL:          "about:blank",
				devtools:            true,
				defaultContextMenus: true,
				builtInErrorPage:    true,
				statusBar:           true,
				zoomControl:         true,
			},
		},
	}

	for _, s := range []string{"WEBVIEW2_BROWSER_EXECUTABLE_FOLDER", "WEBVIEW2_USER_DATA_FOLDER", "WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS", "WEBVIEW2_RELEASE_CHANNEL_PREFERENCE"} {
		os.Unsetenv(s)
	}

	dll, err := webviewloader.New()
	if err != nil {
		return nil, err
	}

	wv.dll = dll.Proc("CreateCoreWebView2EnvironmentWithOptions")

	if err := wv.createWindow(); err != nil {
		return nil, fmt.Errorf("failed to create the window: %w", err)
	}

	for _, option := range options {
		option(wv)
	}

	if err := wv.initializeWindow(); err != nil {
		return nil, fmt.Errorf("failed to initialize the window: %w", err)
	}

	if err := wv.browser.Navigate(wv.browser.config.initialURL); err != nil {
		return nil, fmt.Errorf("failed at the initial navigation: %w", err)
	}

	return wv, nil
}

func (wv *WebView) createWindow() error {
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

	wv.window.handle, err = user32.CreateWindowExW(
		"webview",
		"",
		user32.CW_USEDEFAULT,
		user32.CW_USEDEFAULT,
		int(wv.window.config.width),
		int(wv.window.config.height),
		hinstance,
	)

	if err != nil {
		return fmt.Errorf("failed to create the window: %w", err)
	}

	webviewContext.set(wv.window.handle, wv)

	return nil
}

func (wv *WebView) initializeWindow() error {
	if err := wv.window.SetTitle(wv.window.config.title); err != nil {
		return fmt.Errorf("failed to set the window title: %w", err)
	}

	if err := wv.window.SetSize(wv.window.config.width, wv.window.config.height); err != nil {
		return fmt.Errorf("failed to set the window size: %w", err)
	}

	if err := wv.window.Center(); err != nil {
		return fmt.Errorf("failed to center the window: %w", err)
	}

	if err := wv.window.Show(); err != nil {
		return fmt.Errorf("failed to show the window: %w", err)
	}

	if err := wv.window.Focus(); err != nil {
		return fmt.Errorf("failed to set focus: %w", err)
	}

	if err := wv.browser.embed(wv); err != nil {
		return fmt.Errorf("failed to embed the browser: %w", err)
	}

	if err := wv.browser.resize(); err != nil {
		return fmt.Errorf("failed to resize the browser: %w", err)
	}

	if err := wv.browser.saveSettings(); err != nil {
		return fmt.Errorf("failed to save browser settings: %w", err)
	}

	return nil
}

func (wv *WebView) Terminate() error {
	return user32.PostQuitMessage(0)
}

func wndproc(hwnd, msg, wp, lp uintptr) uintptr {
	if wv, ok := webviewContext.get(windows.Handle(hwnd)); ok {
		switch msg {
		case user32.WMSize:
			_ = wv.browser.resize()
		case user32.WMClose:
			_ = user32.DestroyWindow(windows.Handle(hwnd))
		case user32.WMDestroy:
			_ = wv.Terminate()
		case user32.WMGetMinMaxInfo:
			lpmmi := (*user32.MinMaxInfo)(unsafe.Pointer(lp))

			if wv.window.config.maxWidth > 0 && wv.window.config.maxHeight > 0 {
				maxSize := user32.Point{
					X: wv.window.config.maxWidth,
					Y: wv.window.config.maxHeight,
				}

				lpmmi.MaxSize = maxSize
				lpmmi.MaxTrackSize = maxSize
			}

			if wv.window.config.minWidth > 0 && wv.window.config.minHeight > 0 {
				lpmmi.MinTrackSize = user32.Point{
					X: wv.window.config.minWidth,
					Y: wv.window.config.minHeight,
				}
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

func (wv *WebView) Run() error {
	for {
		msg, err := user32.GetMessageW()
		if err != nil {
			return fmt.Errorf("failed to get message: %w", err)
		}

		if msg == nil || msg.Message == user32.WMQuit {
			return nil
		}

		err = user32.TranslateMessage(msg)
		if err != nil {
			return fmt.Errorf("failed to translate message: %w", err)
		}

		// TODO: Closing the window while it's trying to dispatch the image causes an error here.
		// We should probably ignore it.
		err = user32.DispatchMessageW(msg)
		if err != nil {
			return fmt.Errorf("failed to dispatch message: %w", err)
		}
	}
}
