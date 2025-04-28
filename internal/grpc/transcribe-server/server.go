package transcribeserver

import (
	"context"
	"fmt"

	msu_loggingv1 "github.com/makarmolochaev/msu-logging-protos/gen/go/msu-logging"
	"google.golang.org/grpc"
)

type AudioProcessor interface {
	WhenAudioTranscribed(taskId int32, transcribedText string) error
	WhenProtocolIsReady(taskId int32, protocolText string) error
}

type serverAPI struct {
	msu_loggingv1.UnimplementedTranscribeServer
	msu_loggingv1.UnimplementedProtocolServer
	audio_service AudioProcessor
}

func Register(gRPC *grpc.Server, audio_service AudioProcessor) {
	serverApi := &serverAPI{
		audio_service: audio_service,
	}
	msu_loggingv1.RegisterTranscribeServer(gRPC, serverApi)
	msu_loggingv1.RegisterProtocolServer(gRPC, serverApi)
}

func (s *serverAPI) SendTranscribeResult(
	ctx context.Context,
	req *msu_loggingv1.TranscribeResult,
) (*msu_loggingv1.Result, error) {

	fmt.Println("Recieved gRPC message SendTranscribeResult")

	if req.GetSuccess() {
		err := s.audio_service.WhenAudioTranscribed(req.GetTaskId(), req.GetResult())
		if err == nil {
			return &msu_loggingv1.Result{Success: true}, nil
		}
	}

	return &msu_loggingv1.Result{Success: false}, nil

}

func (s *serverAPI) SendProtocolResult(
	ctx context.Context,
	req *msu_loggingv1.ProtocolResult,
) (*msu_loggingv1.Result, error) {

	fmt.Println("Recieved gRPC message SendProtocolResult")

	if req.GetSuccess() {
		err := s.audio_service.WhenProtocolIsReady(req.GetTaskId(), req.GetResult())
		if err == nil {
			return &msu_loggingv1.Result{Success: true}, nil
		}
	}

	return &msu_loggingv1.Result{Success: false}, nil

}
