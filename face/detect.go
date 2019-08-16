package face

import (
	"github.com/wailovet/osmanthuswine/src/core"
	"github.com/wailovet/pigo-web/pigoutil"
	"io/ioutil"
)

type Detect struct {
	core.Controller
}

func (that *Detect) Index() {
	if that.Request.FILE == nil {
		that.DisplayByError("请上传文件", 500)
	}
	file, err := that.Request.FILE.Open()
	that.CheckErrDisplayByError(err)
	data, err := ioutil.ReadAll(file)
	that.CheckErrDisplayByError(err)
	detects, err := pigoutil.DetectFaces(data)
	that.CheckErrDisplayByError(err)
	result := pigoutil.FilterMarker(detects, 3)
	that.DisplayByData(result)
}
