package transcribeserver

import (
	"context"

	msu_loggingv1 "github.com/makarmolochaev/msu-logging-protos/gen/go/msu-logging"
	"google.golang.org/grpc"
)

type serverAPI struct {
	msu_loggingv1.UnimplementedTranscribeServer
}

func Register(gRPC *grpc.Server) {
	msu_loggingv1.RegisterTranscribeServer(gRPC, &serverAPI{})
}

func (s *serverAPI) SendTranscribeResult(
	ctx context.Context,
	req *msu_loggingv1.TranscribeResult,
) (*msu_loggingv1.Result, error) {
	panic("implement me")
}
