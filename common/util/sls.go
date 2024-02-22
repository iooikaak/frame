package util

import (
	"encoding/json"
	"github.com/iooikaak/frame/xlog"
)

const TimeoutSet = 200

type SlsMessage struct {
	Key              string `json:"key"`
	Url              string `json:"url"`
	SupplierId       int64  `json:"supplier_id"`
	UniqueId         string `json:"unique_id"`
	Action           string `json:"action"`
	Username         string `json:"username"`
	UserId           int64  `json:"user_id"`
	GoodsId          int64  `json:"goods_id"`
	GoodsName        string `json:"goods_name"`
	SupplierSkuId    int64  `json:"supplier_sku_id"`
	SupplierAttrId   int64  `json:"supplier_attr_id"`
	SupplierAttrName string `json:"supplier_attr_name"`
	GoodsSkuId       int64  `json:"product_sku_id"`
	GoodsAttrId      int64  `json:"goods_attr_id"`
	GoodsAttrName    string `json:"goods_attr_name"`
}

func Info(msg *SlsMessage) {
	b, err := json.Marshal(msg)
	if err != nil {
		xlog.Infof("Sls info failed err: %v", err)
		return
	}
	xlog.Infof(string(b))
}
