package server

import (
	"fmt"
	"strconv"
)

const msgSizeLength int = 6

// decoderMsgSize is the processor reading and decoding the next message size
// a valid message size of 74 is [0074]
// the maximum message size is 4096 bytes
type decoderMsgSize struct {
	processorBase
	msgSize int
}

// NewDecoderMsgSize creates a new decoderMsgSize instance
func NewDecoderMsgSize() *decoderMsgSize {
	dMsgData := NewDecoderMsgData()
	
	d := &decoderMsgSize{
		processorBase {
			0,
			nil,
			0,
			dMsgData,
		},
		0,
	}
	dMsgData.pNext = d
	return d
}

// processReadComplete validates the msg size format and retrieves the value
// if success it creates the next connection processor instance - decoderMsgData
func (d *decoderMsgSize) processReadComplete() (Errorer) {
	// reset the decoder internal state
	d.totalProcessed = 0
	
	// validate msg size format to be like [0125], with max value 4096
	valid := true
	
	if d.data[0] != '[' || d.data[len(d.data)-1] != ']' {
		fmt.Println( "decoderMsgSize processReadComplete() - message size format is invalid, no start and/or end []")
		valid = false
	}

	msgSize, err := strconv.Atoi(string(d.data[1:len(d.data)-1]))
	if err != nil {
		fmt.Printf( "decoderMsgSize processReadComplete() - message size is invalid %s\r\n", err)
		valid = false
	} else if msgSize <= 0 || msgSize > 4096 {
		fmt.Printf("decoderMsgSize processReadComplete() - message value is invalid %d\r\n", msgSize)
		valid = false
		} else { 
			d.msgSize = msgSize
	}

	if !valid {
		fmt.Println("decoderMsgSize processReadComplete() - message size is not valid")
		return protocolError {
			status: invalidMsgSize,
			msg: fmt.Sprintf("decoderMsgSize processReadComplete() - invalid message size data %s\r\n", string(d.data[:])),
		}
	}

	fmt.Printf("decoderMsgSize processReadComplete() - message size value %d, next connection read processor - decodeMsgData\r\n", d.msgSize)

	// set the message size value in the next processor decoding the message data
	d.pNext.setDataLen(d.msgSize)
	return nil
}

// processReadData copies bytes from the dataRead buffer into the decoderMsgSize interval buffer up to the preset msg size len (6 bytes)
// if more read bytes are required it returns a needMoreConnRead Errorrer to request a new connection read
// if all msg size bytes were received it validates the format and decodes the msg size value
func (d *decoderMsgSize) processReadData(dataRead []byte) (int, Errorer) {
	d.msgSize = 0

	processed, err := d.processReadDataBase(dataRead)
	if err != nil {
		return processed, err
	}

	err = d.processReadComplete()
	return processed, err
}

// getNextProcessor returns the decoderMsgSize next connection processor - decoderMsgData
func (d *decoderMsgSize) getNextProcessor() (processor) {
	return d.pNext
}