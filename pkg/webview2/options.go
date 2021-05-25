package webview2

type Option func(*WebView)

func WithSize(width, height int32) Option {
	return func(wv *WebView) {
		wv.window.config.width = width
		wv.window.config.height = height
	}
}

func WithMinSize(width, height int32) Option {
	return func(wv *WebView) {
		wv.window.config.minWidth = width
		wv.window.config.minHeight = height
	}
}

func WithMaxSize(width, height int32) Option {
	return func(wv *WebView) {
		wv.window.config.maxWidth = width
		wv.window.config.maxHeight = height
	}
}

func WithTitle(title string) Option {
	return func(wv *WebView) {
		wv.window.config.title = title
	}
}

func WithURL(url string) Option {
	return func(wv *WebView) {
		wv.browser.config.initialURL = url
	}
}

func WithBuiltinErrorPage(enabled bool) Option {
	return func(wv *WebView) {
		wv.browser.config.builtInErrorPage = enabled
	}
}

func WithDefaultContextMenus(enabled bool) Option {
	return func(wv *WebView) {
		wv.browser.config.defaultContextMenus = enabled
	}
}

func WithDefaultScriptDialogs(enabled bool) Option {
	return func(wv *WebView) {
		wv.browser.config.defaultScriptDialogs = enabled
	}
}

func WithDevtools(enabled bool) Option {
	return func(wv *WebView) {
		wv.browser.config.devtools = enabled
	}
}

func WithHostObjects(enabled bool) Option {
	return func(wv *WebView) {
		wv.browser.config.hostObjects = enabled
	}
}

func WithStatusBar(enabled bool) Option {
	return func(wv *WebView) {
		wv.browser.config.statusBar = enabled
	}
}

func WithScript(enabled bool) Option {
	return func(wv *WebView) {
		wv.browser.config.script = enabled
	}
}

func WithWebMessage(enabled bool) Option {
	return func(wv *WebView) {
		wv.browser.config.webMessage = enabled
	}
}

func WithZoomControl(enabled bool) Option {
	return func(wv *WebView) {
		wv.browser.config.zoomControl = enabled
	}
}
