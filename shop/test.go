package main

import (
	"shop/common"
	"shop/datamodels"
)

func main() {
	data := map[string]string{"ID":"1", "productName":"imooc测试结构体","productNum":"2", "prodoctImage":"123","productUrl":"http://url"}
	product := &datamodels.Product{}
	common.DataToStructByTagSql(data, product)
}
