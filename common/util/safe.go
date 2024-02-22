package util

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var fileTypeMap sync.Map

const verifyByteSize = 15

func init() {
	// images
	fileTypeMap.Store("ffd8ff", "jpg")   //JPEG (jpg)
	fileTypeMap.Store("89504e47", "png") //PNG (png)
	fileTypeMap.Store("47494638", "gif") //GIF (gif)
	fileTypeMap.Store("49492a00", "tif") //TIFF (tif)
	fileTypeMap.Store("424d", "bmp")     //BMP (bmp)
	/*CAD*/
	fileTypeMap.Store("41433130", "dwg") //CAD (dwg)
	fileTypeMap.Store("38425053", "psd") //Photoshop (psd)
	/* 日记本 */
	fileTypeMap.Store("7b5c727466", "rtf")
	fileTypeMap.Store("3c3f786d6c", "xml")
	fileTypeMap.Store("68746d6c3e", "html")
	// 邮件
	fileTypeMap.Store("44656c69766572792d646174653a", "eml") //Email
	fileTypeMap.Store("d0cf11e0", "doc")                     //MS Excel 注意：word、msi 和 excel的文件头一样
	//excel2003版本文件
	fileTypeMap.Store("d0cf11e0", "xls")
	fileTypeMap.Store("5374616e64617264204a", "mdb") //MS Access (mdb)
	fileTypeMap.Store("252150532d41646f6265", "ps")
	fileTypeMap.Store("255044462d312e", "pdf") //Adobe Acrobat (pdf)
	fileTypeMap.Store("504b0304", "docx")      //docx文件
	//excel2007以上版本文件
	fileTypeMap.Store("504b0304", "xlsx")
	fileTypeMap.Store("52617221", "rar")
	fileTypeMap.Store("57415645", "wav")
	fileTypeMap.Store("41564920", "avi")
	fileTypeMap.Store("2e524d46", "rm")
	fileTypeMap.Store("000001ba", "mpg")
	fileTypeMap.Store("000001b3", "mpg")
	fileTypeMap.Store("6d6f6f76", "mov")
	fileTypeMap.Store("3026b2758e66cf11", "asf")
	fileTypeMap.Store("4d546864", "mid")
	fileTypeMap.Store("1f8b08", "gz")
}

// bytesToHexString 获取前面结果字节的二进制
func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}

// getFileType 用文件前面几个字节来判断
// fSrc: 文件字节流（就用前面几个字节）
// 这个包的用途是用来防止客户端注入攻击(上传的文件包含一段病毒代码脚本等)
func getFileType(fSrc []byte) string {
	var fileType string
	fileCode := bytesToHexString(fSrc)

	fileTypeMap.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(string)
		if strings.HasPrefix(fileCode, strings.ToLower(k)) {
			fileType = v
			return false
		}
		return true
	})
	return fileType
}

//IsSafeStaticResource true为安全，false为不安全
func IsSafeStaticResource(staticSourceURL string) (bool, error) {
	resp, err := http.Get(staticSourceURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	//读取静态资源到内存中
	fSrc, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	//第一重防范逻辑：是否是我们允许的静态资源
	if getFileType(fSrc[:verifyByteSize]) == "" {
		arr := strings.Split(staticSourceURL, ".")
		extendType := arr[len(arr)-1]
		return false, fmt.Errorf("传入的资源类型：%s，是不支持的资源类型", extendType)
	}
	//第二重防范逻辑：不能有js代码<script>标签
	if bytes.Contains(fSrc, []byte("<script>")) {
		return false, errors.New("静态资源文件里面不能含有JS代码")
	}
	return true, nil
}
