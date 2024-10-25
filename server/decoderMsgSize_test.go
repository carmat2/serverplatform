package server

import (
	"testing"
)


func Test_decoderMsgSize_validMsgSize(t *testing.T) {
	dataRead := []byte {'[', '0', '1', '2', '5', ']'}
	decoderMsgSize_validMsgSize(dataRead, 125, t)

	dataRead = []byte {'[', '0', '0', '2', '5', ']'}
	decoderMsgSize_validMsgSize(dataRead, 25, t)

	dataRead = []byte {'[', '0', '0', '0', '5', ']'}
	decoderMsgSize_validMsgSize(dataRead, 5, t)

	dataRead = []byte {'[', '1', '0', '2', '4', ']'}
	decoderMsgSize_validMsgSize(dataRead, 1024, t)

	dataRead = []byte {'[', '4', '0', '9', '6', ']'}
	decoderMsgSize_validMsgSize(dataRead, 4096, t)

	dataRead = []byte {'[', '4', '0', '9', '6', ']', 'E', 'X', 'T'}
	decoderMsgSize_validMsgSize(dataRead, 4096, t)
}

func Test_decoderMsgSize_invalidMsgSize(t *testing.T) {
	dataRead := []byte {']', '0', '1', '2', '5', ']'}
	decoderMsgSize_invalidMsgSize(dataRead, invalidMsgSize, t)

	dataRead = []byte {']', '0', '1', '2', '5', '['}
	decoderMsgSize_invalidMsgSize(dataRead, invalidMsgSize, t)

	dataRead = []byte {'[', '0', '1', '2', '5', '['}
	decoderMsgSize_invalidMsgSize(dataRead, invalidMsgSize, t)

	dataRead = []byte {'[', ' ', '1', '2', '5', ']'}
	decoderMsgSize_invalidMsgSize(dataRead, invalidMsgSize, t)

	dataRead = []byte {'[', '0', '0', '0', '0', ']'}
	decoderMsgSize_invalidMsgSize(dataRead, invalidMsgSize, t)

	dataRead = []byte {'[', '0', ' ', '1', '2', ']'}
	decoderMsgSize_invalidMsgSize(dataRead, invalidMsgSize, t)

	dataRead = []byte {'[', '4', '0', '9', '7', ']'}
	decoderMsgSize_invalidMsgSize(dataRead, invalidMsgSize, t)

}

func decoderMsgSize_validMsgSize(dataRead []byte, expected int, t *testing.T) {
	logger.Tracef("****** testing decoderMsgSize - valid msg size %s", string(dataRead))
	d := NewDecoderMsgSize()
	d.setDataLen(msgSizeLength)
	processed, err := d.processReadData(dataRead)
	logger.Tracef("processed %d, err %v", processed, err)
	if processed != 6 {
		t.Error("incorrect processed value - expected 6, received", processed)
	}
	if err != nil {
		t.Error("incorrect err value - expected nil, received", err)
	}
	if(d.msgSize != expected) {
		t.Error("incorrect message size value - expected / received", expected, d.msgSize)
	}
}

func decoderMsgSize_invalidMsgSize(dataRead []byte, errExpected status, t *testing.T) {
	logger.Tracef("****** testing decoderMsgSize - invalid msg size %s, error expected %d", string(dataRead), errExpected)
	d := NewDecoderMsgSize()
	d.setDataLen(msgSizeLength)
	processed, err := d.processReadData(dataRead)
	logger.Tracef("processed %d, err %v", processed, err)
	if processed != 6 {
		t.Error("incorrect processed value - expected 6, received", processed)
	}
	if err == nil {
		t.Error("incorrect err value - expected error, received nil")
	} else if err.Status() != errExpected {
		t.Error("incorrect err value - expected / received", errExpected, err.Status())
	}
}