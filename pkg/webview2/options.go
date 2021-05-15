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

// Initial testing tells me that WithDevTools seems to be the only one that actually works.
func WithDevtools(enabled bool) Option {
	return func(wv *WebView) {
		wv.browser.config.devtools = enabled
	}
}

// Context menus are also disabled if you disable DevTools.
// See the comment above WithDevtools.
func WithDefaultContextMenus(enabled bool) Option {
	return func(wv *WebView) {
		wv.browser.config.defaultContextMenus = enabled
	}
}

// See the comment above WithDevtools.
func WithBuiltinErrorPage(enabled bool) Option {
	return func(wv *WebView) {
		wv.browser.config.builtInErrorPage = enabled
	}
}

// See the comment above WithDevtools.
func WithStatusBar(enabled bool) Option {
	return func(wv *WebView) {
		wv.browser.config.statusBar = enabled
	}
}

// See the comment above WithDevtools.
func WithZoomControl(enabled bool) Option {
	return func(wv *WebView) {
		wv.browser.config.zoomControl = enabled
	}
}
