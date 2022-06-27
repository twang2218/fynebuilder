//	The default theme of Fyne is missing CJK font support, UnicodeTheme
//	 inherts the default theme with CJK characters support.
package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type UnicodeTheme struct{}

func (t *UnicodeTheme) Font(s fyne.TextStyle) fyne.Resource {
	if resourceFont == nil {
		return theme.DefaultTheme().Font(s)
	}
	return resourceFont
}

func (t *UnicodeTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (t *UnicodeTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (t *UnicodeTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
