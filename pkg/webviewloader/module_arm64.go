package webviewloader

import _ "embed"

//go:embed arm64/WebView2Loader.dll
var moduleBin []byte
