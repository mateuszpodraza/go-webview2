package webview2

func WithTitle(title string) Option {
	return func(wv *WebView) error {
		return wv.SetTitle(title)
	}
}

func WithSize(width, height int, hint Hint) Option {
	return func(wv *WebView) error {
		return wv.SetSize(width, height, hint)
	}
}

func WithURL(url string) Option {
	return func(wv *WebView) error {
		return wv.Navigate(url)
	}
}

// The following are not yet implemented, they do nothing.
func WithDevtools(enabled bool) Option {
	return func(wv *WebView) error {
		wv.config.devtoolsDisabled = !enabled
		return nil
	}
}

func WithDefaultContextMenus(enabled bool) Option {
	return func(wv *WebView) error {
		wv.config.defaultContextMenusDisabled = !enabled
		return nil
	}
}

func WithDefaultScriptDialogs(enabled bool) Option {
	return func(wv *WebView) error {
		wv.config.defaultScriptDialogsDisabled = !enabled
		return nil
	}
}

func WithHostObjects(enabled bool) Option {
	return func(wv *WebView) error {
		wv.config.hostObjectsDisallowed = !enabled
		return nil
	}
}

func WithBuiltInErrorPage(enabled bool) Option {
	return func(wv *WebView) error {
		wv.config.builtInErrorPageDisabled = !enabled
		return nil
	}
}

func WithScript(enabled bool) Option {
	return func(wv *WebView) error {
		wv.config.scriptDisabled = !enabled
		return nil
	}
}

func WithStatusBar(enabled bool) Option {
	return func(wv *WebView) error {
		wv.config.statusBarDisabled = !enabled
		return nil
	}
}

func WithWebMessage(enabled bool) Option {
	return func(wv *WebView) error {
		wv.config.webMessageDisabled = !enabled
		return nil
	}
}

func WithZoomControl(enabled bool) Option {
	return func(wv *WebView) error {
		wv.config.zoomControlDisabled = !enabled
		return nil
	}
}
