package webview2

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"unsafe"

	"github.com/mattpodraza/webview2/pkg/user32"
	"golang.org/x/sys/windows"
)

type edge struct {
	hwnd                windows.Handle
	controller          *iCoreWebView2Controller
	webview             *iCoreWebView2
	inited              uintptr
	envCompleted        *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler
	controllerCompleted *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler
	webMessageReceived  *iCoreWebView2WebMessageReceivedEventHandler
	permissionRequested *iCoreWebView2PermissionRequestedEventHandler
	msgcb               func(string)
}

func newEdge() *edge {
	e := &edge{}
	e.envCompleted = newICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler(e)
	e.controllerCompleted = newICoreWebView2CreateCoreWebView2ControllerCompletedHandler(e)
	e.webMessageReceived = newICoreWebView2WebMessageReceivedEventHandler(e)
	e.permissionRequested = newICoreWebView2PermissionRequestedEventHandler(e)
	return e
}

func (e *edge) Embed(hwnd windows.Handle) error {
	e.hwnd = hwnd

	currentExePath := make([]uint16, windows.MAX_PATH)

	_, err := windows.GetModuleFileName(windows.Handle(0), &currentExePath[0], windows.MAX_PATH)
	if err != nil {
		return fmt.Errorf("failed to get module file name: %w", err)
	}

	dataPath := filepath.Join(os.Getenv("AppData"), filepath.Base(windows.UTF16ToString(currentExePath)))

	res, err := createCoreWebView2EnvironmentWithOptions(nil, windows.StringToUTF16Ptr(dataPath), 0, e.envCompleted)
	if err != nil {
		return fmt.Errorf("failed to call WebView2Loader: %w", err)
	} else if res != 0 {
		return fmt.Errorf("invalid createCoreWebView2EnvironmentWithOptions result: %08x", res)
	}

	for {
		if atomic.LoadUintptr(&e.inited) != 0 {
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

	return e.Init("window.external={invoke:s=>window.chrome.webview.postMessage(s)}")
}

func (e *edge) Navigate(url string) error {
	_, _, err := e.webview.vtbl.Navigate.Call(
		uintptr(unsafe.Pointer(e.webview)),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(url))),
	)
	if err != nil && !errors.Is(err, errSuccess) {
		return err
	}

	return nil
}

func (e *edge) Init(script string) error {
	_, _, err := e.webview.vtbl.AddScriptToExecuteOnDocumentCreated.Call(
		uintptr(unsafe.Pointer(e.webview)),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(script))),
		0,
	)
	if err != nil && !errors.Is(err, errSuccess) {
		return err
	}

	return nil
}

func (e *edge) Eval(script string) error {
	_, _, err := e.webview.vtbl.ExecuteScript.Call(
		uintptr(unsafe.Pointer(e.webview)),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(script))),
		0,
	)
	if err != nil && !errors.Is(err, errSuccess) {
		return err
	}

	return nil
}

func (e *edge) QueryInterface(refiid, object uintptr) uintptr {
	return 0
}

func (e *edge) AddRef() uintptr {
	return 1
}

func (e *edge) Release() uintptr {
	return 1
}

func (e *edge) EnvironmentCompleted(res uintptr, env *iCoreWebView2Environment) uintptr {
	if int64(res) < 0 {
		return res
	}

	// TODO: Stop ignoring these errors
	_, _, _ = env.vtbl.CreateCoreWebView2Controller.Call(
		uintptr(unsafe.Pointer(env)),
		uintptr(e.hwnd),
		uintptr(unsafe.Pointer(e.controllerCompleted)),
	)

	return 0
}

func (e *edge) ControllerCompleted(res uintptr, controller *iCoreWebView2Controller) uintptr {
	if int64(res) < 0 {
		log.Fatalf("Creating controller failed with %08x", res)
	}

	// TODO: Stop ignoring these errors
	_, _, _ = controller.vtbl.AddRef.Call(uintptr(unsafe.Pointer(controller)))
	e.controller = controller

	var token _EventRegistrationToken
	// TODO: Stop ignoring these errors
	_, _, _ = controller.vtbl.GetCoreWebView2.Call(
		uintptr(unsafe.Pointer(controller)),
		uintptr(unsafe.Pointer(&e.webview)),
	)

	// TODO: Stop ignoring these errors
	_, _, _ = e.webview.vtbl.AddRef.Call(
		uintptr(unsafe.Pointer(e.webview)),
	)

	// TODO: Stop ignoring these errors
	_, _, _ = e.webview.vtbl.AddWebMessageReceived.Call(
		uintptr(unsafe.Pointer(e.webview)),
		uintptr(unsafe.Pointer(e.webMessageReceived)),
		uintptr(unsafe.Pointer(&token)),
	)

	// TODO: Stop ignoring these errors
	_, _, _ = e.webview.vtbl.AddPermissionRequested.Call(
		uintptr(unsafe.Pointer(e.webview)),
		uintptr(unsafe.Pointer(e.permissionRequested)),
		uintptr(unsafe.Pointer(&token)),
	)

	atomic.StoreUintptr(&e.inited, 1)

	return 0
}

func (e *edge) MessageReceived(sender *iCoreWebView2, args *iCoreWebView2WebMessageReceivedEventArgs) uintptr {
	var message *uint16
	// TODO: Stop ignoring these errors
	_, _, _ = args.vtbl.TryGetWebMessageAsString.Call(
		uintptr(unsafe.Pointer(args)),
		uintptr(unsafe.Pointer(message)),
	)

	e.msgcb(windows.UTF16PtrToString(message))

	// TODO: Stop ignoring these errors
	_, _, _ = sender.vtbl.PostWebMessageAsString.Call(
		uintptr(unsafe.Pointer(sender)),
		uintptr(unsafe.Pointer(message)),
	)

	windows.CoTaskMemFree(unsafe.Pointer(message))
	return 0
}

func (e *edge) PermissionRequested(sender *iCoreWebView2, args *iCoreWebView2PermissionRequestedEventArgs) uintptr {
	var kind _CoreWebView2PermissionKind
	// TODO: Stop ignoring these errors
	_, _, _ = args.vtbl.GetPermissionKind.Call(
		uintptr(unsafe.Pointer(args)),
		uintptr(kind),
	)
	if kind == _CoreWebView2PermissionKindClipboardRead {
		// TODO: Stop ignoring these errors
		_, _, _ = args.vtbl.PutState.Call(
			uintptr(unsafe.Pointer(args)),
			uintptr(_CoreWebView2PermissionStateAllow),
		)
	}
	return 0
}
