package structs

import (
	"smsService/models/www"
)

type InsertData struct {
	InserData www.CjSmsLog
	IsCount   bool
}

var InsertDataChan = make(chan InsertData, 1000)

func (i *InsertData) WhriteInsertData(data www.CjSmsLog, IsCount bool) {
	data.Create()

}

func (i *InsertData) LoopInsertData(data chan InsertData) {
	for {
		select {
		//InsertData
		case insertdata := <-data:
			go i.WhriteInsertData(insertdata.InserData, insertdata.IsCount)
		}
	}

}
