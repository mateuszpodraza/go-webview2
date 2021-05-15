# webview2

A proof of concept for using the Microsoft Edge WebView2 without cgo and with embedded copies of the webview DLL.

This is a fork of https://github.com/jchv/go-webview2, with a different API that suited my needs better than the original package.

It also uses some bits from https://github.com/Inkeliz/gowebview, specifically the way the COM procedures are called.

## Notice

This requires you to have the [WebView2 runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) installed, as it doesn't ship with Windows.

## Non-goals

* EdgeHTML fallback
* Support for other platforms than Windows on amd64
