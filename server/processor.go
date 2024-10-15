package server

import (
	"fmt"
)

// processorBase contains the common data for a connection processor
// the received data buffer is created with the already detemined length of the data to process
// once the processor completes processing the data, it sets pNext as the current connection data processor
type processorBase struct {
	dataLen int
	data []byte
	totalProcessed int
	pNext processor
}

// processor is the interface to be implemented by all connection processors
type processor interface {
	processReadData(dataRead []byte) (processed int, err Errorer)
	processReadComplete() (err Errorer)
	setDataLen(len int)
	getNextProcessor() (p processor)
}

// setDataLen creates the processor instance data buffer
func (pBase *processorBase) setDataLen(len int) {
	pBase.dataLen = len
	pBase.data = make([]byte, len)
}

// processReadDataBase is a common function used by all processors to append in-progress read data to their internal data buffer
func (pBase *processorBase) processReadDataBase(dataRead []byte) (int, Errorer) {
	fmt.Printf("processReadDataBase() - input dataRead %v / %s \r\n", dataRead, string(dataRead))
	processed := copy(pBase.data[pBase.totalProcessed:], dataRead)
	pBase.totalProcessed += processed
	fmt.Printf("processReadDataBase() - copied %d new bytes to in-progress data; %d bytes copied so far - %v / %s \r\n", processed, pBase.totalProcessed, pBase.data, string(pBase.data))

	if pBase.totalProcessed != len(pBase.data) {
		fmt.Println("processReadDataBase() - request new connection read")
		return processed, protocolError {
			status: needMoreConnRead,
			msg: fmt.Sprintf("processReadDataBase() - need %d more input bytes", pBase.dataLen-pBase.totalProcessed),
		}
	}

	fmt.Println("processReadDataBase() - all data bytes were received")
	return processed, nil
}

