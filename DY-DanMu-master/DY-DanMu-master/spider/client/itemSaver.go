package client

import (
	"DY-DanMu/DMconfig/config"
	"DY-DanMu/persistServer/item"
	"DY-DanMu/persistServer/rpcsupport"
	"DY-DanMu/web/util"

	Log "github.com/sirupsen/logrus"
)

type ItemSaverStruct struct {
	ItemCountAll int
	ItemCountMin int
}

// ItemSaver:将Item放入gorutine进行分发上传到数据库服务器
func (e *ItemSaverStruct) ItemSaver(host string) (chan item.Item, error) {
	clinet, err := rpcsupport.NewClinet(host)
	if err != nil {
		return nil, err
	}
	out := make(chan item.Item)
	go func() {
		for {
			item := <-out
			Log.Debugf("当前items #%d :%v", e.ItemCountAll, item)
			e.ItemCountAll++
			e.ItemCountMin++
			// Call RPC to save item
			result := ""
			err := clinet.Call(config.SpiderConfig.ItemSaverRpc, item, &result)
			if err != nil {
				err = util.RpcClientShutDownErrorhandler(err)
				if err == nil {
					// 重新连接
					clinet, err = rpcsupport.NewClinet(host)
					if err == nil {
						err = clinet.Call(config.SpiderConfig.ItemSaverRpc, item, &result)
					}
				}
			}
			if err != nil {
				Log.Errorf("Item Saver: saving item %v:%v", item, err)
			}
		}
	}()
	return out, nil
}
