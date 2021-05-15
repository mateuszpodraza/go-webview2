package webview2

import (
	"fmt"

	"github.com/mattpodraza/webview2/v2/pkg/user32"
	"golang.org/x/sys/windows"
)

type windowConfig struct {
	title               string
	width, height       int32
	maxWidth, maxHeight int32
	minWidth, minHeight int32
}

type window struct {
	config *windowConfig
	handle windows.Handle
}

func (w *window) focus() error {
	return user32.SetFocus(w.handle)
}

func (w *window) show() error {
	return user32.ShowWindow(w.handle, user32.SW_SHOW)
}

func (w *window) setTitle(title string) error {
	return user32.SetWindowTextW(w.handle, title)
}

func (w *window) center() error {
	sx, err := user32.GetSystemMetrics(user32.SystemMetricsCxScreen)
	if err != nil {
		return fmt.Errorf("failed to get the horizontal screen size: %w", err)
	}

	sy, err := user32.GetSystemMetrics(user32.SystemMetricsCyScreen)
	if err != nil {
		return fmt.Errorf("failed to get the vertical screen size: %w", err)
	}

	rect := user32.Rect{
		Left:   0,
		Top:    0,
		Right:  w.config.width,
		Bottom: w.config.height,
	}

	if err := user32.AdjustWindowRec(&rect, user32.WSOverlappedWindow, true); err != nil {
		return fmt.Errorf("failed to adjust window rect: %w", err)
	}

	rect.Left = (int32(sx) - rect.Right) / 2
	rect.Top = (int32(sy) - rect.Bottom) / 2

	err = user32.SetWindowPos(
		w.handle,
		rect.Left,
		rect.Top,
		rect.Right-rect.Left,
		rect.Bottom-rect.Top,
		user32.SWPNoZOrder|user32.SWPNoActivate|user32.SWPNoSize|user32.SWPFrameChanged,
	)

	if err != nil {
		return fmt.Errorf("failed to set the window position: %w", err)
	}

	return nil
}

func (w *window) setSize(width, height int32) error {
	rect := user32.Rect{
		Left:   0,
		Top:    0,
		Right:  width,
		Bottom: height,
	}

	if err := user32.AdjustWindowRec(&rect, user32.WSOverlappedWindow, true); err != nil {
		return fmt.Errorf("failed to adjust window rect: %w", err)
	}

	err := user32.SetWindowPos(
		w.handle,
		rect.Left,
		rect.Top,
		rect.Right-rect.Left,
		rect.Bottom-rect.Top,
		user32.SWPNoZOrder|user32.SWPNoActivate|user32.SWPNoMove|user32.SWPFrameChanged,
	)

	if err != nil {
		return fmt.Errorf("failed to set the window position: %w", err)
	}

	return nil
}
