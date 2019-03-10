package server

//Stream wraps the SpacyServer_PlayServer and a channel to close the grpc conn
type Stream struct {
	SpacyServer_GetUpdatesServer
	ShutdownC chan struct{}
	ActionC   chan *Action
}
