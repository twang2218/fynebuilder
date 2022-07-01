package fynebuilder

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/url"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage/repository"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"
)

func newWidget(tag ObjectTag, res ResourceDict) fyne.CanvasObject {
	switch tag.Tag {
	//	container
	case "Max":
		c := container.NewMax()
		size := tag.Attributes.GetSize()
		c.Resize(size)
		c.Move(tag.Attributes.GetPosition())
		return c
	case "Center":
		c := container.NewCenter()
		size := tag.Attributes.GetSize()
		c.Resize(size)
		c.Move(tag.Attributes.GetPosition())
		return c
	case "VBox":
		c := container.NewVBox()
		size := tag.Attributes.GetSize()
		c.Resize(size)
		c.Move(tag.Attributes.GetPosition())
		return c
	case "HBox":
		c := container.NewHBox()
		size := tag.Attributes.GetSize()
		c.Resize(size)
		c.Move(tag.Attributes.GetPosition())
		return c
	case "Grid":
		cols := tag.Attributes.GetInt("cols")
		rows := tag.Attributes.GetInt("rows")
		adaptive := tag.Attributes.GetBool("adaptive")

		var c *fyne.Container
		if rows == 0 && cols > 0 {
			c = container.NewGridWithColumns(cols)
		} else if rows > 0 && cols == 0 {
			c = container.NewGridWithRows(rows)
		} else if rows > 0 && cols > 0 && !adaptive {
			c = container.NewGridWrap(fyne.NewSize(float32(cols), float32(rows)))
		} else if adaptive {
			c = container.NewAdaptiveGrid(cols + rows)
		}
		size := tag.Attributes.GetSize()
		c.Resize(size)
		c.Move(tag.Attributes.GetPosition())
		return c
	case "Padded":
		c := container.NewPadded()
		size := tag.Attributes.GetSize()
		c.Resize(size)
		c.Move(tag.Attributes.GetPosition())
		return c
	// case "Border":
	//	<Border><Left>xxxx</Left><Right>yyyyy</Right></Border>
	// 	return container.NewBorder()
	case "FormLayout":
		//	tag <Form> conflicts with widget.Form, so use `<FormLayout>` instead.
		c := container.New(layout.NewFormLayout())
		size := tag.Attributes.GetSize()
		c.Resize(size)
		c.Move(tag.Attributes.GetPosition())
		return c
	case "Spacer":
		c := layout.NewSpacer()
		size := tag.Attributes.GetSize()
		c.Resize(size)
		c.Move(tag.Attributes.GetPosition())
		return c
	case "AppTab":
		var loc container.TabLocation
		switch strings.ToLower(tag.Attributes.GetString("location")) {
		default:
		case "top":
			loc = container.TabLocationTop
		case "bottom":
			loc = container.TabLocationBottom
		case "leading":
			loc = container.TabLocationLeading
		case "trailing":
			loc = container.TabLocationTrailing
		}
		tabs := container.NewAppTabs()
		tabs.SetTabLocation(loc)
		return tabs
	case "DocTab":
		return container.NewDocTabs()
	case "TabItem":
		//	TabItem is not a CanvasObject, it will be created in 'attach()'.
		return nil
	case "Scroll":
		//	missing object
		return container.NewScroll(nil)
	case "HScroll":
		return container.NewHScroll(nil)
	case "VScroll":
		return container.NewVScroll(nil)
	case "HSplit":
		return container.NewHSplit(nil, nil)
	case "VSplit":
		return container.NewVSplit(nil, nil)

	//	widget
	case "Accordion":
		return widget.NewAccordion()
	case "AccordionItem":
		//	AccordionItem is not a CanvasObject, it will be created in 'attach()'.
		return nil
	case "Button":
		var b *widget.Button
		//	label
		label := tag.Attributes.GetString("label")
		//	icon
		icon := tag.Attributes.GetResource("icon", res)
		if icon != nil {
			b = widget.NewButtonWithIcon(label, icon, nil)
		} else {
			b = widget.NewButton(tag.Attributes.GetString("label"), nil)
		}
		//	align
		switch strings.ToLower(tag.Attributes.GetString("align")) {
		default:
		case "center":
			b.Alignment = widget.ButtonAlignCenter
		case "leading":
			b.Alignment = widget.ButtonAlignLeading
		case "trailing":
			b.Alignment = widget.ButtonAlignTrailing
		}
		//	placement (icon)
		switch strings.ToLower(tag.Attributes.GetString("placement")) {
		case "leading":
			b.IconPlacement = widget.ButtonIconLeadingText
		case "trailing":
			b.IconPlacement = widget.ButtonIconTrailingText
		}
		//	importance
		switch strings.ToLower(tag.Attributes.GetString("importance")) {
		default:
		case "medium":
			b.Importance = widget.MediumImportance
		case "high":
			b.Importance = widget.HighImportance
		case "low":
			b.Importance = widget.LowImportance
		}
		//	return
		return b
	case "Card":
		//	title
		title := tag.Attributes.GetString("title")
		//	subtitle
		subtitle := tag.Attributes.GetString("subtitle")
		//	image
		image := tag.Attributes.GetImage("image", res)
		//	content will be set in `attach()`
		c := widget.NewCard(title, subtitle, nil)
		c.Image = image
		return c
	case "Check":
		//	label
		label := tag.Attributes.GetString("label")
		//	checked
		checked := tag.Attributes.GetBool("checked")
		c := widget.NewCheck(label, nil)
		c.SetChecked(checked)
		return c
	case "CheckGroup":
		//	options
		opts := tag.Attributes.GetList("options")
		//	selected
		selected := tag.Attributes.GetList("selected")
		c := widget.NewCheckGroup(opts, nil)
		c.SetSelected(selected)
		return c
	case "Entry":
		var e *widget.Entry
		switch tag.Attributes.GetString("type") {
		case "multiline":
			e = widget.NewMultiLineEntry()
		case "password":
			e = widget.NewPasswordEntry()
		default:
			e = widget.NewEntry()
		}
		if tag.Attributes.ContainsStyle() {
			e.TextStyle = tag.Attributes.GetTextStyle()
		}
		return e
	case "FileIcon":
		src := tag.Attributes.GetString("src")
		return widget.NewFileIcon(repository.NewFileURI(src))
	case "Form":
		submit := tag.Attributes.GetString("submit")
		cancel := tag.Attributes.GetString("cancel")
		//	form content will be attached in `attach()`.
		f := widget.NewForm()
		f.SubmitText = submit
		f.CancelText = cancel
		return f
	case "Link":
		url, err := url.Parse(tag.Attributes.GetString("url"))
		if err != nil {
			log.Errorf("<Link>: %v", err)
		}
		//	Link text will be added in `setContent()`
		if tag.Attributes.Contains("align") || tag.Attributes.ContainsStyle() {
			align := tag.Attributes.GetTextAlignment()
			style := tag.Attributes.GetTextStyle()
			return widget.NewHyperlinkWithStyle("", url, align, style)
		} else {
			return widget.NewHyperlink("", url)
		}
	case "Icon":
		icon := tag.Attributes.GetResource("src", res)
		return widget.NewIcon(icon)
	case "Label":
		//	content will be added in `setContent()`
		if tag.Attributes.Contains("align") || tag.Attributes.ContainsStyle() {
			align := tag.Attributes.GetTextAlignment()
			style := tag.Attributes.GetTextStyle()
			return widget.NewLabelWithStyle("", align, style)
		} else {
			return widget.NewLabel("")
		}
	// case "Menu":
	// case "PopUp":
	// case "PopUpMenu":
	case "ProgressBar":
		min := tag.Attributes.GetFloat64("min")
		max := tag.Attributes.GetFloat64("max")
		value := tag.Attributes.GetFloat64("value")
		p := widget.NewProgressBar()
		p.Min = min
		p.Max = max
		p.Value = value
		return p
	case "ProgressBarInfinite":
		return widget.NewProgressBarInfinite()
	case "RadioGroup":
		opts := tag.Attributes.GetList("options")
		r := widget.NewRadioGroup(opts, nil)
		if tag.Attributes.Contains("selected") {
			r.SetSelected(tag.Attributes.GetString("selected"))
		}
	case "RichText":
		switch tag.Attributes.GetString("type") {
		default:
		case "segment":
			return widget.NewRichText()
		case "markdown":
			//	attach the content in setContent()
			return widget.NewRichTextFromMarkdown("")
		case "text":
			//	attach the content in setContent()
			return widget.NewRichTextWithText("")
		}
	case "Select":
		opts := tag.Attributes.GetList("options")
		s := widget.NewSelect(opts, nil)
		if tag.Attributes.Contains("align") {
			s.Alignment = tag.Attributes.GetTextAlignment()
		}
		if tag.Attributes.Contains("selected") {
			s.Selected = tag.Attributes.GetString("selected")
		}
		return s
	case "SelectEntry":
		opts := tag.Attributes.GetList("options")
		e := widget.NewSelectEntry(opts)
		if tag.Attributes.ContainsStyle() {
			e.TextStyle = tag.Attributes.GetTextStyle()
		}
	case "Separator":
		return widget.NewSeparator()
	case "Slider":
		min := tag.Attributes.GetFloat64("min")
		max := tag.Attributes.GetFloat64("max")
		step := tag.Attributes.GetFloat64("step")
		value := tag.Attributes.GetFloat64("value")
		s := widget.NewSlider(min, max)
		s.Step = step
		s.SetValue(value)
		return s
	// case "Table":
	case "TextGrid":
		//	content will be filled in `setContent()`
		t := widget.NewTextGrid()
		t.ShowLineNumbers = tag.Attributes.GetBool("line_number")
		t.ShowWhitespace = tag.Attributes.GetBool("whitespace")
		t.TabWidth = tag.Attributes.GetInt("tabwidth")
		return t
	case "TextGridRow":
		//	TextGridRow is not CanvasObject, it will be created in `attach()`
		return nil
	case "TextGridCell":
		//	TextGridCell is not CanvasObject, it will be created in `attach()`
		return nil
	case "Toolbar":
		//	ToolbarAction,ToolbarSeparator is not CanvasObject, it will be created in `attach()`
		bar := widget.NewToolbar()
		return bar

	// canvas
	case "Circle":
		color := tag.Attributes.GetColor("color")
		c := canvas.NewCircle(color)
		c.Resize(tag.Attributes.GetSize())
		c.Move(tag.Attributes.GetPosition())
		return c
	case "Image":
		img := tag.Attributes.GetImage("src", res)
		size := tag.Attributes.GetSize()
		img.SetMinSize(size)
		img.Resize(size)
		img.Move(tag.Attributes.GetPosition())
		return img
	case "Line":
		color := tag.Attributes.GetColor("color")
		l := canvas.NewLine(color)
		l.Resize(tag.Attributes.GetSize())
		l.Move(tag.Attributes.GetPosition())
		strokeWidth := tag.Attributes.GetFloat32("stroke_width")
		if strokeWidth > 1 {
			l.StrokeWidth = strokeWidth
		}
		return l
	case "LinearGradient":
		start := tag.Attributes.GetColor("start_color")
		end := tag.Attributes.GetColor("stop_color")
		t := tag.Attributes.GetString("type")
		var g *canvas.LinearGradient
		switch t {
		case "horizontal":
			g = canvas.NewHorizontalGradient(start, end)
		case "vertical":
			g = canvas.NewVerticalGradient(start, end)
		default:
			angle := tag.Attributes.GetFloat64("angle")
			g = canvas.NewLinearGradient(start, end, angle)
		}
		size := tag.Attributes.GetSize()
		g.Resize(size)
		g.SetMinSize(size)
		g.Move(tag.Attributes.GetPosition())
		return g
	case "RadialGradient":
		start := tag.Attributes.GetColor("start_color")
		end := tag.Attributes.GetColor("stop_color")
		r := canvas.NewRadialGradient(start, end)
		r.CenterOffsetX = tag.Attributes.GetFloat64("offsetX")
		r.CenterOffsetY = tag.Attributes.GetFloat64("offsetY")
		size := tag.Attributes.GetSize()
		r.Resize(size)
		r.SetMinSize(size)
		r.Move(tag.Attributes.GetPosition())
		return r
	case "Raster":
		src := tag.Attributes.GetString("src")
		f, err := os.Open(src)
		if err != nil {
			log.Error(err)
		}
		img, _, err := image.Decode(f)
		if err != nil {
			log.Error(err)
		}
		r := canvas.NewRasterFromImage(img)
		r.Translucency = 1 - tag.Attributes.GetFloat64("alpha")
		size := tag.Attributes.GetSize()
		r.Resize(size)
		r.SetMinSize(size)
		r.Move(tag.Attributes.GetPosition())
		return r
	case "Rectangle":
		color := tag.Attributes.GetColor("color")
		r := canvas.NewRectangle(color)
		r.StrokeColor = tag.Attributes.GetColor("stroke_color")
		r.StrokeWidth = tag.Attributes.GetFloat32("stroke_width")
		size := tag.Attributes.GetSize()
		r.Resize(size)
		r.SetMinSize(size)
		r.Move(tag.Attributes.GetPosition())
		return r
	case "Text":
		c := tag.Attributes.GetColor("color")
		t := canvas.NewText("", c)
		t.Alignment = tag.Attributes.GetTextAlignment()
		t.TextSize = tag.Attributes.GetFloat32("text_size")
		t.TextStyle = tag.Attributes.GetTextStyle()
		size := tag.Attributes.GetSize()
		t.Resize(size)
		t.SetMinSize(size)
		t.Move(tag.Attributes.GetPosition())
		return t
	}
	return nil
}

func attach(current, parent, grand ObjectTag, res ResourceDict) {
	switch parent.Tag {
	case "TabItem":
		label := parent.Attributes.GetString("label")
		item := container.NewTabItem(label, current.Object)
		if grand.Object != nil {
			if apptabs, ok := grand.Object.(*container.AppTabs); ok {
				apptabs.Append(item)
			} else if doctabs, ok := grand.Object.(*container.DocTabs); ok {
				doctabs.Append(item)
			}
		}
	case "Scroll":
		fallthrough
	case "HScroll":
		fallthrough
	case "VScroll":
		if scroll, ok := parent.Object.(*container.Scroll); ok {
			if scroll.Content == nil {
				scroll.Content = current.Object
			}
		}
	case "HSplit":
		fallthrough
	case "VSplit":
		if split, ok := parent.Object.(*container.Split); ok {
			if split.Leading == nil {
				split.Leading = current.Object
			} else {
				split.Trailing = current.Object
			}
		}
	case "AccordionItem":
		title := parent.Attributes.GetString("title")
		if grand.Object != nil {
			if accordion, ok := grand.Object.(*widget.Accordion); ok {
				item := widget.NewAccordionItem(title, current.Object)
				accordion.Append(item)
			}
		}
	case "Card":
		if card, ok := parent.Object.(*widget.Card); ok {
			card.SetContent(current.Object)
		}
	case "FormItem":
		label := parent.Attributes.GetString("label")
		if form, ok := grand.Object.(*widget.Form); ok {
			item := widget.NewFormItem(label, current.Object)
			form.AppendItem(item)
		}
	case "TextGridRow":
		if grid, ok := grand.Object.(*widget.TextGrid); ok {
			row := widget.TextGridRow{}
			grid.Rows = append(grid.Rows, row)
		}
	case "Toolbar":
		if bar, ok := parent.Object.(*widget.Toolbar); ok {
			switch current.Tag {
			case "ToolbarAction":
				icon := current.Attributes.GetResource("icon", res)
				action := widget.NewToolbarAction(icon, nil)
				bar.Append(action)
			case "ToolbarSeparator":
				sep := widget.NewToolbarSeparator()
				bar.Append(sep)
			case "ToolbarSpacer":
				spacer := widget.NewToolbarSpacer()
				bar.Append(spacer)
			}
		}
	default:
		if c, ok := parent.Object.(*fyne.Container); ok {
			c.Add(current.Object)
			// fmt.Printf("Add %s to %s\n", current.Tag, parent.Tag)
		}
	}
}

func setContent(obj ObjectTag, content string) {
	switch o := obj.Object.(type) {
	case *widget.Label:
		//	<Label>...</Label>
		o.SetText(content)
		// log.Tracef("<Label>.SetText(%s)", t)
	case *widget.Entry:
		//	<Entry>...</Entry>
		o.SetText(content)
		// log.Tracef("<Entry>.SetText(%s)", t)
	case *widget.Hyperlink:
		//	<Link>...</Link>
		o.SetText(content)
	case *widget.RichText:
		switch obj.Attributes.GetString("type") {
		case "markdown":
			//	attach the content in setContent()
			if r, ok := obj.Object.(*widget.RichText); ok {
				r.ParseMarkdown(content)
			}
		case "text":
			//	attach the content in setContent()
			if r, ok := obj.Object.(*widget.RichText); ok {
				if len(r.Segments) > 0 {
					if s, ok := r.Segments[0].(*widget.TextSegment); ok {
						s.Text = content
					}
				}
			}
		}
	case *widget.TextGrid:
		//	<TextGrid>...</TextGrid>
		o.SetText(content)
	case *canvas.Text:
		//	<Text>...</Text>
		o.Text = content
		// log.Tracef("<Text>.Text = %s", t)
	}
}
