package wsapp

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"msu-logging-backend/internal/services/audioservice"
	"msu-logging-backend/internal/storage/filerepository"
	"net/http"

	"github.com/gorilla/websocket"
)

type App struct {
	log           *slog.Logger
	server        *http.Server
	port          int
	upgrader      websocket.Upgrader
	audio_service *audioservice.AudioService
	fileRepo      *filerepository.FileRepository
	certFile      string
	keyFile       string
}

func New(
	log *slog.Logger,
	port int,
	audio_service *audioservice.AudioService,
	certFile string,
	keyFile string,
) *App {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			//origin:= r.Header.Get("Origin")
			return true
		},
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	app := &App{
		log:           log,
		server:        server,
		port:          port,
		upgrader:      upgrader,
		certFile:      certFile,
		keyFile:       keyFile,
		fileRepo:      filerepository.NewFileRepository(),
		audio_service: audio_service,
	}

	mux.HandleFunc("/ws", app.handleWebSocket)

	return app
}

func (a *App) handleWebSocketConnection(conn *websocket.Conn) error {
	filename := a.fileRepo.CreateAudioFile()
	defer func() {
		a.fileRepo.CloseAudioFile(filename)
		if err := a.audio_service.WhenWebsocketClosed(filename); err != nil {
			a.log.Error("Failed to process closed websocket", slog.String("error", err.Error()))
		}
		a.fileRepo.DeleteAudioFile(filename)
	}()

	a.log.Info("Created audio file")
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				a.log.Info("Client disconnected gracefully")
				return nil
			}
			a.log.Error("Read error", slog.String("error", err.Error()))
			return fmt.Errorf("read error: %w", err)
		}

		if messageType == websocket.BinaryMessage {
			if err := a.fileRepo.WriteAudioData(filename, p); err != nil {
				a.log.Error("Write error", slog.String("error", err.Error()))
				return fmt.Errorf("write error: %w", err)
			}
		}
	}
}

func (a *App) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	fmt.Println(token)

	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.log.Error("WebSocket upgrade failed", slog.String("error", err.Error()))
		return
	}

	if err := a.handleWebSocketConnection(conn); err != nil {
		a.log.Error("WebSocket connection error", slog.String("error", err.Error()))
	}
}

func (a *App) Run() error {
	const op = "wsapp.Run"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("Secure WebSocket server is running...")

	if err := a.server.ListenAndServeTLS(a.certFile, a.keyFile); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() error {
	const op = "wsapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("Stopping WebSocket server", slog.Int("port", a.port))

	return a.server.Close()
}
