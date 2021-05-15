package webview2

import (
	"sync"

	"golang.org/x/sys/windows"
)

type webviewContextStore struct {
	mu    sync.RWMutex
	store map[windows.Handle]*WebView
}

var webviewContext = &webviewContextStore{
	store: map[windows.Handle]*WebView{},
}

func (wcs *webviewContextStore) set(hwnd windows.Handle, wv *WebView) {
	wcs.mu.Lock()
	defer wcs.mu.Unlock()

	wcs.store[hwnd] = wv
}

func (wcs *webviewContextStore) get(hwnd windows.Handle) (*WebView, bool) {
	wcs.mu.Lock()
	defer wcs.mu.Unlock()

	wv, ok := wcs.store[hwnd]
	return wv, ok
}
