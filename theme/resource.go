package theme

import (
	"bytes"
	_ "embed"
	"io"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	log "github.com/sirupsen/logrus"
	"github.com/ulikunitz/xz"
)

// //go:embed fonts/wqy-microhei.ttc
// var font_wqy_microhei []byte
// var resourceFont = NewResource("wqy-microhei.ttc", font_wqy_microhei)

//go:embed fonts/wqy-microhei.ttc.xz
var font_wqy_microhei_xz []byte
var resourceFont = NewResource("wqy-microhei.ttc.xz", font_wqy_microhei_xz)

func ExtractXZ(original []byte) ([]byte, error) {
	r, err := xz.NewReader(bytes.NewReader(original))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}

func NewResource(name string, data []byte) *fyne.StaticResource {
	if strings.HasSuffix(name, ".xz") {
		//	extract if '.xz'
		t := time.Now()
		var err error
		data, err = ExtractXZ(data)
		if err != nil {
			log.Error(err)
			return nil
		}
		log.Debugf("Extracted %q from xz in %v", name, time.Since(t))
		name = strings.TrimSuffix(name, ".xz")
	}

	return fyne.NewStaticResource(name, data)
}
