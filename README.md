# webview2

This is a fork of https://github.com/jchv/go-webview2, with a different API that suited my needs better than the original package.

A proof of concept for using the Microsoft Edge WebView2 without cgo.

It relies on the excellent [go-winloader](https://github.com/jchv/go-winloader), so be warned - no guarantees of API/runtime stability.

My initial testing deemed it _stable enough_, though.

## Notice

This requires you to have the [WebView2 runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) installed, as it doesn't ship with Windows.

## Non-goals

* EdgeHTML fallback
* Support for other platforms than Windows
