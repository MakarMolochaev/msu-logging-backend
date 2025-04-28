package rmqapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"msu-logging-backend/internal/config"
	rabbitmodels "msu-logging-backend/internal/domain/models"
	"os"
	"time"

	amqp "github.com/streadway/amqp"
)

type App struct {
	log     *slog.Logger
	conn    *amqp.Connection
	channel *amqp.Channel
	uri     string
}

func New(
	log *slog.Logger,
	config *config.Config,
) *App {
	return &App{
		log: log,
		uri: os.Getenv("RABBITMQ_CONN_STR"),
	}
}

func (a *App) connect() error {
	const op = "rmqapp.connect"

	var err error
	a.conn, err = amqp.Dial(a.uri)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.channel, err = a.conn.Channel()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Run() error {
	const op = "rmqapp.Run"

	log := a.log.With(slog.String("op", op))

	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err = a.connect()
		if err == nil {
			break
		}
		log.Error("Failed to connect to RabbitMQ",
			slog.Int("attempt", i+1),
			slog.String("error", err.Error()))
		time.Sleep(time.Second * time.Duration(i+1))
	}
	if err != nil {
		return fmt.Errorf("%s: failed after %d retries: %w", op, maxRetries, err)
	}

	log.Info("RabbitMQ publisher is ready")
	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() error {
	const op = "rmqapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("Stopping RabbitMQ connection")

	var err error
	if a.channel != nil {
		if closeErr := a.channel.Close(); closeErr != nil {
			err = fmt.Errorf("%s: channel close error: %w", op, closeErr)
		}
	}

	if a.conn != nil {
		if closeErr := a.conn.Close(); closeErr != nil {
			err = fmt.Errorf("%s: connection close error: %w", op, closeErr)
		}
	}

	return err
}

func (a *App) SendMessage(queueName string, message string) error {
	const op = "rmqapp.SendMessage"
	log := a.log.With(
		slog.String("op", op),
	)
	_, err := a.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error("Error in declaring queue")
		return fmt.Errorf("error in declaring queue:%v", err)
	}

	err = a.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		log.Error("Error in publishing in queue")
		return fmt.Errorf("error in publishing in queue:%v", err)
	}
	return nil
}

// for NN Service
func (a *App) SendTranscribeRequest(queueName string, transcribeRequestData rabbitmodels.TranscribeRequest) error {
	const op = "rmqapp.SendTranscribeRequest"
	log := a.log.With(
		slog.String("op", op),
	)
	_, err := a.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error("Error in declaring queue")
		return fmt.Errorf("error in declaring queue:%v", err)
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(transcribeRequestData)
	if err != nil {
		log.Error("Error in encoding JSON")
		return fmt.Errorf("error in encoding json: %v", err)
	}

	jsonBody := buf.Bytes()

	err = a.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         jsonBody,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		log.Error("Error in publishing in queue")
		return fmt.Errorf("error in publishing in queue:%v", err)
	}
	return nil
}

// for NLP Service
func (a *App) SendProtocolRequest(queueName string, protocolRequestData rabbitmodels.ProtocolRequest) error {
	const op = "rmqapp.SendProtocolRequest"

	log := a.log.With(
		slog.String("op", op),
	)
	_, err := a.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error("Error in declaring queue")
		return fmt.Errorf("error in declaring queue:%v", err)
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(protocolRequestData)
	if err != nil {
		log.Error("Error in encoding JSON")
		return fmt.Errorf("error in encoding json: %v", err)
	}

	jsonBody := buf.Bytes()

	err = a.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         jsonBody,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		log.Error("Error in publishing in queue")
		return fmt.Errorf("error in publishing in queue:%v", err)
	}
	return nil
}
