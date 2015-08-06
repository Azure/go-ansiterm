// +build windows

package winterm

// effectiveSr gets the current effective scroll region in buffer coordinates
func (h *WindowsAnsiEventHandler) effectiveSr(window SMALL_RECT) scrollRegion {
	top := AddInRange(window.Top, h.sr.top, window.Top, window.Bottom)
	bottom := AddInRange(window.Top, h.sr.bottom, window.Top, window.Bottom)
	if top >= bottom {
		top = window.Top
		bottom = window.Bottom
	}
	return scrollRegion{top: top, bottom: bottom}
}

func (h *WindowsAnsiEventHandler) scrollPageUp() error {
	return h.scrollPage(1)
}

func (h *WindowsAnsiEventHandler) scrollPageDown() error {
	return h.scrollPage(-1)
}

func (h *WindowsAnsiEventHandler) scrollPage(param int) error {
	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil {
		return err
	}

	return h.scroll(param, scrollRegion{info.Window.Top, info.Window.Bottom}, info)
}

func (h *WindowsAnsiEventHandler) scrollUp(param int) error {
	info, err := GetConsoleScreenBufferInfo(h.fd)
	if err != nil {
		return err
	}

	sr := h.effectiveSr(info.Window)
	return h.scroll(param, sr, info)
}

func (h *WindowsAnsiEventHandler) scrollDown(param int) error {
	return h.scrollUp(-param)
}

// scroll scrolls the provided scroll region by param lines. The scroll region is in buffer coordinates.
func (h *WindowsAnsiEventHandler) scroll(param int, sr scrollRegion, info *CONSOLE_SCREEN_BUFFER_INFO) error {
	logger.Infof("scroll: scrollTop: %d, scrollBottom: %d", sr.top, sr.bottom)
	logger.Infof("scroll: windowTop: %d, windowBottom: %d", info.Window.Top, info.Window.Bottom)

	// Copy from and clip to the scroll region (full buffer width)
	scrollRect := SMALL_RECT{
		Top:    sr.top,
		Bottom: sr.bottom,
		Left:   0,
		Right:  info.Size.X - 1,
	}

	// Origin to which area should be copied
	destOrigin := COORD{
		X: 0,
		Y: sr.top - SHORT(param),
	}

	char := CHAR_INFO{
		UnicodeChar: ' ',
		Attributes:  h.attributes,
	}

	if err := ScrollConsoleScreenBuffer(h.fd, scrollRect, scrollRect, destOrigin, char); err != nil {
		return err
	}
	return nil
}
