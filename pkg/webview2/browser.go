package webview2

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/mattpodraza/webview2/v2/pkg/com"
	"github.com/mattpodraza/webview2/v2/pkg/hresult"
	"github.com/mattpodraza/webview2/v2/pkg/user32"
	"golang.org/x/sys/windows"
)

type browserConfig struct {
	initialURL string

	builtInErrorPage     bool
	defaultContextMenus  bool
	defaultScriptDialogs bool
	devtools             bool
	hostObjects          bool
	script               bool
	statusBar            bool
	webMessage           bool
	zoomControl          bool
}

type browser struct {
	hwnd windows.Handle

	config     *browserConfig
	view       *com.ICoreWebView2
	controller *com.ICoreWebView2Controller
	settings   *com.ICoreWebView2Settings

	controllerCompleted int32
}

func (wv *WebView) Browser() *browser {
	return wv.browser
}

func (b *browser) embed(wv *WebView) error {
	b.hwnd = wv.window.handle

	exePath := make([]uint16, windows.MAX_PATH)

	_, err := windows.GetModuleFileName(windows.Handle(0), &exePath[0], windows.MAX_PATH)
	if err != nil {
		return fmt.Errorf("failed to get module file name: %w", err)
	}

	dataPath := filepath.Join(os.Getenv("AppData"), filepath.Base(windows.UTF16ToString(exePath)))

	r1, _, err := wv.dll.Call(0, uint64(uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(dataPath)))), 0, uint64(wv.environmentCompletedHandler()))
	hr := hresult.HRESULT(r1)

	if err != nil && err != errOK {
		return fmt.Errorf("failed to call CreateCoreWebView2EnvironmentWithOptions: %w", err)
	}

	if hr > hresult.S_OK {
		return fmt.Errorf("failed to call CreateCoreWebView2EnvironmentWithOptions: %s", hr)
	}

	for {
		if atomic.LoadInt32(&b.controllerCompleted) != 0 {
			break
		}

		msg, err := user32.GetMessageW()
		if err != nil {
			return err
		}

		if msg == nil {
			break
		}

		err = user32.TranslateMessage(msg)
		if err != nil {
			return err
		}

		err = user32.DispatchMessageW(msg)
		if err != nil {
			return err
		}
	}

	settings := new(com.ICoreWebView2Settings)

	r, _, err := syscall.Syscall(b.view.VTBL.GetSettings, 2, uintptr(unsafe.Pointer(b.view)), uintptr(unsafe.Pointer(&settings)), 0)
	if !errors.Is(err, errOK) {
		return err
	}

	hr = hresult.HRESULT(r)
	if hr > hresult.S_OK {
		return fmt.Errorf("failed to get webview settings: %s", hr)
	}

	b.settings = settings

	return nil
}

func (b *browser) resize() error {
	if b.controller == nil {
		return errors.New("nil controller")
	}

	bounds, err := user32.GetClientRect(b.hwnd)
	if err != nil {
		return fmt.Errorf("failed to get client rect: %w", err)
	}

	_, _, err = syscall.Syscall(
		b.controller.VTBL.PutBounds, 2,
		uintptr(unsafe.Pointer(b.controller)),
		uintptr(unsafe.Pointer(bounds)),
		0,
	)

	if !errors.Is(err, errOK) {
		return fmt.Errorf("failed to put bounds: %w", err)
	}

	return nil
}

func (b *browser) Navigate(url string) error {
	_, _, err := syscall.Syscall(
		b.view.VTBL.Navigate, 3,
		uintptr(unsafe.Pointer(b.view)),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(url))),
		0,
	)

	if !errors.Is(err, errOK) {
		return err
	}

	return nil
}

func (b *browser) AddScriptToExecuteOnDocumentCreated(script string) error {
	_, _, err := syscall.Syscall(
		b.view.VTBL.AddScriptToExecuteOnDocumentCreated, 3,
		uintptr(unsafe.Pointer(b.view)),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(script))),
		0,
	)

	if !errors.Is(err, errOK) {
		return err
	}

	return nil
}

func (b *browser) ExecuteScript(script string) error {
	_, _, err := syscall.Syscall(
		b.view.VTBL.ExecuteScript, 3,
		uintptr(unsafe.Pointer(b.view)),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(script))),
		0,
	)

	if !errors.Is(err, errOK) {
		return err
	}

	return nil
}

func (b *browser) saveSetting(setter uintptr, enabled bool) error {
	var flag uintptr = 0

	if enabled {
		flag = 1
	}

	_, _, err := syscall.Syscall(
		setter, 3,
		uintptr(unsafe.Pointer(b.settings)),
		flag,
		0,
	)

	if !errors.Is(err, errOK) {
		return fmt.Errorf("failed to save a setting: %w", err)
	}

	return nil
}

func (b *browser) saveSettings() error {
	if err := b.saveSetting(b.settings.VTBL.PutIsBuiltInErrorPageEnabled, b.config.builtInErrorPage); err != nil {
		return err
	}

	if err := b.saveSetting(b.settings.VTBL.PutAreDefaultContextMenusEnabled, b.config.defaultContextMenus); err != nil {
		return err
	}

	if err := b.saveSetting(b.settings.VTBL.PutAreDefaultScriptDialogsEnabled, b.config.defaultScriptDialogs); err != nil {
		return err
	}

	if err := b.saveSetting(b.settings.VTBL.PutAreDevToolsEnabled, b.config.devtools); err != nil {
		return err

	}

	if err := b.saveSetting(b.settings.VTBL.PutAreHostObjectsAllowed, b.config.hostObjects); err != nil {
		return err
	}

	if err := b.saveSetting(b.settings.VTBL.PutIsScriptEnabled, b.config.script); err != nil {
		return err
	}

	if err := b.saveSetting(b.settings.VTBL.PutIsStatusBarEnabled, b.config.statusBar); err != nil {
		return err

	}

	if err := b.saveSetting(b.settings.VTBL.PutIsWebMessageEnabled, b.config.webMessage); err != nil {
		return err
	}

	return b.saveSetting(b.settings.VTBL.PutIsZoomControlEnabled, b.config.zoomControl)
}

func (wv *WebView) environmentCompletedHandler() uintptr {
	h := &com.ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler{
		VTBL: &com.ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVTBL{
			Invoke: windows.NewCallback(func(i uintptr, p uintptr, createdEnvironment *com.ICoreWebView2Environment) uintptr {
				_, _, _ = syscall.Syscall(createdEnvironment.VTBL.CreateCoreWebView2Controller, 3, uintptr(unsafe.Pointer(createdEnvironment)), uintptr(wv.window.handle), wv.controllerCompletedHandler())
				return 0
			}),
		},
	}

	h.VTBL.BasicVTBL = com.NewBasicVTBL(&h.Basic)
	return uintptr(unsafe.Pointer(h))
}

func (wv *WebView) controllerCompletedHandler() uintptr {
	h := &com.ICoreWebView2CreateCoreWebView2ControllerCompletedHandler{
		VTBL: &com.ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerVTBL{
			Invoke: windows.NewCallback(func(i *com.ICoreWebView2CreateCoreWebView2ControllerCompletedHandler, p uintptr, createdController *com.ICoreWebView2Controller) uintptr {
				_, _, _ = syscall.Syscall(createdController.VTBL.AddRef, 1, uintptr(unsafe.Pointer(createdController)), 0, 0)
				wv.browser.controller = createdController

				createdWebView2 := new(com.ICoreWebView2)

				_, _, _ = syscall.Syscall(createdController.VTBL.GetCoreWebView2, 2, uintptr(unsafe.Pointer(createdController)), uintptr(unsafe.Pointer(&createdWebView2)), 0)
				wv.browser.view = createdWebView2

				_, _, _ = syscall.Syscall(wv.browser.view.VTBL.AddRef, 1, uintptr(unsafe.Pointer(wv.browser.view)), 0, 0)

				atomic.StoreInt32(&wv.browser.controllerCompleted, 1)

				return 0
			}),
		},
	}

	h.VTBL.BasicVTBL = com.NewBasicVTBL(&h.Basic)
	return uintptr(unsafe.Pointer(h))
}
