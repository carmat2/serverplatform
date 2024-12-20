package server

import (
	"testing"
)


func Test_decoderMsgData_validMsg(t *testing.T) {
	const jsonMsgOpenSession = `{"name":"opensession","dest":"sp","payload":[{"plugin":"helloworld","token":"123"}]}`
	dataRead := []byte(jsonMsgOpenSession)

	logger.Tracef("****** testing decoderMsgData - valid msg %s", string(dataRead))
	d := NewDecoderMsgData()
	d.setDataLen(len(dataRead))
	processed, err := d.processReadData(dataRead)
	logger.Tracef("processed %d, err %v", processed, err)
	if processed != len(dataRead) {
		t.Error("incorrect processed value - expected / received", len(dataRead), processed)
	}
	if err != nil {
		t.Error("incorrect err value - expected nil, received", err)
	}
	if d.msg == nil {
		t.Error("incorrect d.msg value - expected MsgSessionOpen, received nil")
	}
	
	msgType := getMsgType(d.msg)
	if msgType != msgOpenSession {
		t.Error("incorrect d.msg type - expected MsgSessionOpen, received", msgType)
	}

}