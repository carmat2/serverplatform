package server

type status int

// valid protocolError const values.
const (
	needMoreConnRead status = iota + 1
	needMoreConnWrite
	invalidConnSignature
	invalidMsgSize
	invalidMsgData
	connReadTimeout
	connWriteTimeout
	pluginError
)

// Errorer is the interface implemented by a protocolError instance.
type Errorer interface {
	Error() string
	Status() status
}

// protocolError is the serverplatform internal error.
type protocolError struct {
	status status
	msg string
}

// Error returns the error message for a protocolError.
func (pe protocolError) Error() string {
	return pe.msg
}

// Status returns the status value for a protocolError.
func (pe protocolError) Status() status {
	return pe.status
}