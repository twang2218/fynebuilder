> 中文 [English](README.md)

# Fyne 图形界面构造器

这是一个 Fyne 的图形界面构造工具，它将描述界面构成的标记语言转换为实体界面对象。由于使用了独立的界面描述文件，因此它将界面形态和其代码逻辑进行隔离，从而简化代码，并易于编辑。此外，在设计阶段使用热加载则成为可能。

### 示例 UI XML

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Max>
    <Image id="bg" width="1000" height="600" src="embed:background.jpg" />
    <Center>
        <VBox>
            <HBox>
                <Image id="id" width="250" height="150" src="embed:idcard.jpg" />
                <Label> </Label>
                <Image id="qr" width="150" height="150" src="embed:qrcode.png" />
            </HBox>
            <Label> </Label>
            <Text id="description" color="black">机读操作简单易行，只需将第二代身份证放在阅读器读卡区</Text>
        </VBox>
    </Center>
</Max>
```

### 用以加载 XML 的 Go 代码

```go
func main() {
	a := app.New()
	a.Settings().SetTheme(&theme.UnicodeTheme{})

	var res = fynebuilder.ResourceDict{
		"icon.png":       fyne.NewStaticResource("icon.png", icon),
		"idcard.jpg":     fyne.NewStaticResource("idcard.jpg", idcard),
		"background.jpg": fyne.NewStaticResource("background.jpg", background),
		"qrcode.png":     fyne.NewStaticResource("qrcode.png", qrcode),
	}

	a.SetIcon(res["icon.png"])

	//      窗口
	w := a.NewWindow("访客身份校验")

	watcher := fynebuilder.NewWatcher("demo.ui", res, func(objs fynebuilder.ObjectDict) {
		w.SetContent(objs.GetTop())
	})
	defer watcher.Close()

	w.ShowAndRun()
}
```

### 生成的界面

![Screenshot](cmd/demo/assets/screenshot.png)