package server

import (
	"fmt"
)

// processorBase contains the common data for a connection processor.
// The received data buffer is created with the already detemined length of the data to process.
// Once the processor completes processing the data, it sets pNext as the current connection data processor.
type processorBase struct {
	dataLen int
	data []byte
	totalProcessed int
	pNext processor
}

// Processor is the interface to be implemented by all connection processors.
type processor interface {
	processReadData(dataRead []byte) (processed int, err Errorer)
	processReadComplete() (err Errorer)
	setDataLen(len int)
	getNextProcessor() (p processor)
}

// setDataLen creates the processor instance data buffer.
func (pBase *processorBase) setDataLen(len int) {
	pBase.dataLen = len
	pBase.data = make([]byte, len)
}

// processReadDataBase is a common function used by all processors to append in-progress read data to their internal data buffer.
func (pBase *processorBase) processReadDataBase(dataRead []byte) (int, Errorer) {
	logger.Tracef("processReadDataBase() input dataRead %s", string(dataRead))
	processed := copy(pBase.data[pBase.totalProcessed:], dataRead)
	pBase.totalProcessed += processed
	logger.Tracef("processReadDataBase() - copied %d new bytes to in-progress data; %d bytes copied so far %s", processed, pBase.totalProcessed, string(pBase.data))

	if pBase.totalProcessed != len(pBase.data) {
		logger.Debug("processReadDataBase() - request new connection read")
		return processed, protocolError {
			status: needMoreConnRead,
			msg: fmt.Sprintf("processReadDataBase() - need %d more input bytes", pBase.dataLen-pBase.totalProcessed),
		}
	}

	logger.Trace("processReadDataBase() - all data bytes were received")
	return processed, nil
}

