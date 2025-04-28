package rabbitmodels

type TranscribeRequest struct {
	TaskId        int32
	AudioFileLink string
}

type ProtocolRequest struct {
	TaskId          int32
	TranscribedText string
}
