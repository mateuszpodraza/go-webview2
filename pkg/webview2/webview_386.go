package webview2

import (
	"unsafe"
)

func (e *edge) Resize() error {
	if e.controller == nil {
		return errors.New("nil controller")
	}

	bounds, err := user32.GetClientRect(e.hwnd)
	if err != nil {
		return fmt.Errorf("failed to get client rect: %w", err)
	}

	_, _, err = e.controller.vtbl.PutBounds.Call(
		uintptr(unsafe.Pointer(e.controller)),
		uintptr(bounds.Left),
		uintptr(bounds.Top),
		uintptr(bounds.Right),
		uintptr(bounds.Bottom),
	)
	if err != nil && !errors.Is(err, errSuccess) {
		return fmt.Errorf("failed to put bounds: %w", err)
	}

	return nil
}
