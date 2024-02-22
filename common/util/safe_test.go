package util

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetFileType(t *testing.T) {

	//安全的静态资源类型比如 .jpg
	staticSourceURL := "https://www.twle.cn/static/i/img1.jpg"
	//不支持的静态资源类型比如 .mp4
	//staticSourceURL := "https://video_shejigao.redocn.com/video/202106/20210615/Redcon_202103101235254412763135.mp4"
	resp, err := http.Get(staticSourceURL)
	if err != nil {
		t.Logf("open error: %#v", err)
		return
	}
	defer resp.Body.Close()

	fSrc, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Logf("ioutil.ReadAll error:%#v", err)
		return
	}
	t.Log(getFileType(fSrc[:verifyByteSize]))
}

func TestIsSafeStaticResource(t *testing.T) {
	//安全的静态资源类型比如 .jpg
	staticSourceURL := "https://www.twle.cn/static/i/img1.jpg"
	//不全的静态资源类型比如 .mp4
	//staticSourceURL := "https://video_shejigao.redocn.com/video/202106/20210615/Redcon_202103101235254412763135.mp4"
	//不安全的html静态资源里面有js代码
	//staticSourceURL := "https://wanzhoumo-cdn.wanzhoumo.com/data/public/user_icon/2021/04/19_10/16187992497477038.html"
	boolean, err := IsSafeStaticResource(staticSourceURL)
	t.Logf("IsSafeStaticResource return bool is :%#v", boolean)
	t.Logf("IsSafeStaticResource return error is :%#v", err)
}
