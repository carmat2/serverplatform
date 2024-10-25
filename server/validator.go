package server

import(
	"fmt"
)

const signature string = "$SIGNATURE"

// validator is the process validating the connection signature on a newly accepted connection.
type validator struct {
	processorBase
}

// NewValidator creates a new validator instance
func NewValidator() *validator {
	dMsgSize := NewDecoderMsgSize()
	v := &validator{
		processorBase {
			0,
			nil,
			0,
			dMsgSize,
		},
	}
	return v
}

// processReadComplete validates the signature received in the TCP connection.
func (v *validator) processReadComplete() (err Errorer) {
	// reset the validator internal state
	v.totalProcessed = 0

	recvSignature := string(v.data[:])
	if signature != recvSignature {
		logger.Errorf("validator processReadComplete() - connection signature did not match %s", recvSignature)
		return protocolError {
			status: invalidConnSignature,
			msg: fmt.Sprintf("validator processReadComplete() - invalid signature data %s", recvSignature),
		}
	}

	logger.Debug("validator processReadComplete() - connection signature match, next connection read processor - decodeMsgSize")
	v.pNext.setDataLen(msgSizeLength)

	return nil
}

// processReadData copies bytes from the dataRead buffer into the validator interval buffer up the expected signature len.
// If more read bytes are required it returns a needMoreConnRead Errorrer to request a new connection read.
// If all signature bytes were received it validates the signature.
func (v *validator) processReadData(dataRead []byte) (int, Errorer) {
	processed, err := v.processReadDataBase(dataRead)
	if err != nil {
		return processed, err
	}

	err = v.processReadComplete()
	return processed, err
}

// getNextProcessor returns the validator next connection processor - decoderMsgSize.
func (v *validator) getNextProcessor() (processor) {
	return v.pNext
}
