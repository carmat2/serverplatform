package server

import (
	"testing"
	"time"
)


func Test_validator_validSignatureExactLength(t *testing.T) {
	logger.Trace("****** testing validator - valid signature, one read, exact data length")
	v := NewValidator()
	v.setDataLen(len(signature))
	dataRead := []byte {'$', 'S', 'I', 'G', 'N', 'A', 'T', 'U', 'R', 'E'}
	processed, err := v.processReadData(dataRead)
	logger.Tracef("processed %d, err %v", processed, err)
	if processed != 10 {
		t.Error("incorrect processed value - expected 10, received", processed)
	}
	if err != nil {
		t.Error("incorrect err value - expected nil, received", err)
	}
}

func Test_validator_validSignatureExactLengthMultipleRead(t *testing.T) {
	logger.Trace("****** testing validator - valid signature, multiple reads, exact data length")
	v := NewValidator()
	v.setDataLen(len(signature))
	dataRead1 := []byte {'$', 'S', 'I', 'G'}
	processed, err := v.processReadData(dataRead1)
	logger.Tracef("processed %d, err %v", processed, err)
	if processed != 4 {
		t.Error("incorrect processed value - expected 4, received", processed)
	}
	if err == nil || err.Status() != needMoreConnRead{
		t.Error("incorrect err value - expected protocolError.status.needMoreConnRead")
	}

	time.Sleep(5*time.Second)
	dataRead2 := []byte {'N', 'A', 'T', 'U', 'R', 'E'}
	processed, err = v.processReadData(dataRead2)
	logger.Tracef("processed %d, err %v", processed, err)
	if processed != 6 {
		t.Error("incorrect processed value - expected 6, received", processed)
	}
	if err != nil {
		t.Error("incorrect err value - expected nil, received", err)
	}
}

func Test_validator_validSignatureExtraLength(t *testing.T) {
	logger.Trace("****** testing validator - valid signature, one read, extra data length")
	v := NewValidator()
	v.setDataLen(len(signature))
	dataRead := []byte {'$', 'S', 'I', 'G', 'N', 'A', 'T', 'U', 'R', 'E', '_', 'E', 'X', 'T', 'R', 'A'}
	processed, err := v.processReadData(dataRead)
	logger.Tracef("processed %d, err %v", processed, err)
	if processed != 10 {
		t.Error("incorrect processed value - expected 10, received", processed)
	}
	if err != nil {
		t.Error("incorrect err value - expected nil, received", err)
	}
}

func Test_validator_validSignatureExtraLengthMultipleRead(t *testing.T) {
	logger.Trace("****** testing validator - valid signature, multiple read, extra data length")
	v := NewValidator()
	v.setDataLen(len(signature))
	dataRead1 := []byte {'$', 'S', 'I', 'G', 'N', 'A', 'T', 'U', 'R'}
	processed, err := v.processReadData(dataRead1)
	logger.Tracef("processed %d, err %v", processed, err)
	if processed != 9 {
		t.Error("incorrect processed value - expected 9, received", processed)
	}
	if err == nil || err.Status() != needMoreConnRead {
		t.Error("incorrect err value - expected protocolError.status.needMoreConnRead")
	}

	time.Sleep(5*time.Second)
	dataRead2 := []byte {'E', '_', 'E', 'X', 'T', 'R', 'A'}
	processed, err = v.processReadData(dataRead2)
	logger.Tracef("processed %d, err %v", processed, err)
	if processed != 1 {
		t.Error("incorrect processed value - expected 10, received", processed)
	}
	if err != nil {
		t.Error("incorrect err value - expected nil, received", err)
	}
}

func Test_validator_invalidSignature(t *testing.T) {
	logger.Trace("****** testing validator - invalid signature")
	v := NewValidator()
	v.setDataLen(len(signature))
	dataRead := []byte {'$', 'S', 'I', 'G', 'N', 'A', 'T', 'U', 'R', '$'}
	processed, err := v.processReadData(dataRead)
	logger.Tracef("processed %d, err %v", processed, err)
	if processed != 10 {
		t.Error("incorrect processed value - expected 10, received", processed)
	}
	if err == nil || err.Status() != invalidConnSignature {
		t.Error("incorrect err value - expected protocolError.status.invalidSignature")
	}
}