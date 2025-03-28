package wsapp

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

type App struct {
	log      *slog.Logger
	server   *http.Server
	port     int
	upgrader websocket.Upgrader
	handler  func(*websocket.Conn)
	certFile string
	keyFile  string
}

func New(
	log *slog.Logger,
	port int,
	handler func(*websocket.Conn),
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
		log:      log,
		server:   server,
		port:     port,
		upgrader: upgrader,
		handler:  handler,
		certFile: certFile,
		keyFile:  keyFile,
	}

	mux.HandleFunc("/ws", app.handleWebSocket)

	return app
}

func (a *App) handleWebSocket(w http.ResponseWriter, r *http.Request) {

	/*
		token = r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// if !isValidToken(token) {}
	*/

	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.log.Error("WebSocket upgrade failed", slog.String("error", err.Error()))
		return
	}
	defer conn.Close()

	a.handler(conn)
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
