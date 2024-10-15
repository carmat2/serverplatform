# serverplatform
TCP Server Platform written in Golang to provide basic common functionality for conferencing-type applications.
It exposes an interface to be implemented by **Plugins** providing business logic.


The serverplatform creates an **Acceptor** to Listen for client connections on port 8081. 
Once a **Connection** is accepted and created, two goroutines handle reads from and writes to the connection.

The read loop uses **processors** to handle and decode received data. One processor instance is active and reading connection data at any given time, and once the processing is completed, it passes the handling to the next processor in the protocol logic.

The connection read loop processors are:
1. **Validator** - validates the connection first 10 bytes to match the signature 
2. **DecoderMsgSize** - decodes the next 6 bytes to decode the next message size (in bytes)
3. **DecoderMsgData** - decodes the next protocol **message** json

   
## serverplatform protocol

10 bytes - connection signature
6 bytes - next protocol message size in the following format [0084]. The maximum size for a messsage is 4096 bytes
msgSize bytes - next protocol message data as json

The **message** common fields are *name*, *dest* and *payload*
- *name* is a unique message name
- *dest* must be set to *sp* if the message is to be decoded by the serverplatform, or *pl* if the message is to be decoded by the Plugin
- *payload* is an array containing json fields with message-specific data
