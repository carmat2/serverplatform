package server

import (
	"encoding/json"
	"fmt"
)

// decodeMsgData is the processor decoding protocol messages.
type decoderMsgData struct {
	processorBase
	msg MessageBase
}

// NewDecoderMsgData creates a new decoderMsgData instance.
func NewDecoderMsgData() *decoderMsgData {
	d := &decoderMsgData{
		processorBase {
			0,
			nil,
			0,
			nil,
		},
		nil,
	}
	return d
}
// processReadComplete first decodes the base message json values (name, dest and payload),
// then it decodes all serverplatform dest messsages from the message payload, and requests the plugin to decode plugin dest messages.
func (d *decoderMsgData) processReadComplete() (err Errorer) {
	// reset the decoder internal state
	d.totalProcessed = 0

	var m *Message
	if errJson := json.Unmarshal(d.data, &m) ; errJson != nil {
		return protocolError {
			status: invalidMsgData,
			msg: fmt.Sprintf("decoderMsgData processReadComplete() - invalid message data %s, underlying error %s", string(d.data[:]), errJson),
		}
	}

	if m.Dest == "sp" {
		d.msg, err = m.decodeServerPlarformMessage()
		if( err != nil) {
			return err
		}
	} else if m.Dest == "pl" {
		// TODO - decode plugin message
		return protocolError {
			status: invalidMsgData,
			msg: fmt.Sprintf("decoderMsgData () - invalid plugin message %s - not implemented", string(d.data[:])),
		}
	} else {
		return protocolError {
			status: invalidMsgData,
			msg: fmt.Sprintf("decoderMsgData (processReadComplete) - invalid message dest %s ", m.Dest),
		}
	}

	return nil
}

// processReadData copies bytes from the dataRead buffer into the decoderMsgData interval buffer up to the current msg size.
// If more read bytes are required it returns a needMoreConnRead Errorrer to request a new connection read.
// If all message bytes were received it decodes the message and passes it to the processing handler.
func (d *decoderMsgData) processReadData(dataRead []byte) (int, Errorer) {
	d.msg = nil
	
	processed, err := d.processReadDataBase(dataRead)
	if err != nil {
		return processed, err
	}

	err = d.processReadComplete()
	return processed, err
}

// getNextProcessor returns the validator next connection processor - nil, to be replaced by the handler processor.
func (d *decoderMsgData) getNextProcessor() (processor) {
	return d.pNext
}