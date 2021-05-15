package webviewloader

import (
	"fmt"

	"github.com/jchv/go-winloader"
)

func New() (winloader.Module, error) {
	dll, err := winloader.LoadFromMemory(moduleBin)
	if err != nil {
		return nil, fmt.Errorf("failed to load the Webview2Loader DLL from memory: %w", err)
	}

	return dll, nil
}
