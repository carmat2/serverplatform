package server

// Plugin is the interface to be implemented by all serverplatform Plugins
type Plugin interface {
	OnSessionAboutToOpen(msg MsgOpenSession)
	OnSessionOpenened(msg MsgOpenSession)
}