package server

import (
	"fmt"
	"encoding/json"
	"reflect"
)

// valid serverplatform message struct types.
type msgType string
const (
	msgOpenSession msgType = "MsgOpenSession"
	msgSessionOpened msgType = "MsgSessionOpened"
	msgJoinSession msgType = "MsgJoinSession"
	msgSessionJoined msgType = "MsgSessionJoined"
)

// Message wraps the common protocol message json fields.
// Name is the unique message name.
// Dest is the message handling destination - "sp" if the mssage is to be handled by the server platform, 
// or "pl" if the message is to be decoded and handled by a Plugin.
type Message struct {
	Name	string	`json:"name"`
	Dest    string	`json:"dest"`
	Payload	[]map[string]interface{}	`json:"payload"`
}

// MessageBase is the interface to be implemented by all protocol messages.
type MessageBase interface {
	decode(payload []map[string]interface{}) (Errorer)
}

// MsgOpenSession is the request to open a session on the set Plugin using the set Token as permission validator.
type MsgOpenSession struct {
	Plugin	string	`json:"plugin"`
	Token	string	`json:"token"`
}

// MsgSessionOpenened is the response to the client after the session with the set Sid was opened.
type MsgSessionOpened struct {
	Plugin	string	`json:"plugin"`
	Sid		int	`json:"sid"`
}

// MsgJoinSession is the request to join a session with the set Sid on the set Plugin using the set Token as permission validator.
type MsgJoinSession struct {
	Plugin	string	`json:"plugin"`
	Token	string	`json:"token"`
	Sid		int		`json:"sid"`
}

// MsgSessionJoined is the response to the client after the client with the set Cid joined the Sid session.
type MsgSessionJoined struct {
	Plugin	string	`json:"plugin"`
	Sid		int		`json:"sid"`
	Cid		int		`json:"cid"`
}

// getMsgType is a generic function to return the struct name of the passed MessageBase.
func getMsgType(m MessageBase) msgType {
    if t := reflect.TypeOf(m); t.Kind() == reflect.Ptr {
        return (msgType) (t.Elem().Name())
    } else {
        return (msgType) (t.Name())
    }
}

// decode creates a message struct object form the passed payload.
func (m *MsgOpenSession) decode(payload []map[string]interface{}) (Errorer){
	arrLen := len(payload)
	if arrLen != 1 {
		return protocolError {
			status: invalidMsgData,
			msg: fmt.Sprintln("MsgOpenSession - invalid message data, payload array size is not 1"),
		}
	}

	jsonData, errJson := json.Marshal(payload[0])
    if errJson != nil {
        return protocolError {
			status: invalidMsgData,
			msg: fmt.Sprintf("MsgOpenSession - invalid message data, underlying error %s", errJson),
		}
    }
    if errJson := json.Unmarshal(jsonData, &m); errJson != nil {
        return protocolError {
			status: invalidMsgData,
			msg: fmt.Sprintf("MsgOpenSession - invalid message data, underlying error %s", errJson),
		}
    }

	return nil
}

// decodeServerMessage is the decoder for all serverplatform messages.
func (m Message) decodeServerPlarformMessage() (msg MessageBase, err Errorer) {
	switch m.Name {
	case "opensession":
		msg = &MsgOpenSession {}
		err = msg.decode(m.Payload)
	default:
		return nil, protocolError {
			status: invalidMsgData,
			msg: fmt.Sprintf("invalid message name %s", m.Name),
		}
	}

	if err != nil {
		return nil, err
	}

	logger.Infof("Success decoding message %s", m)
	return msg, nil
}
