package rmqapp

import (
	"fmt"
	"log/slog"
	"msu-logging-backend/internal/config"
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

func (a *App) Publish(queueName string, body []byte) error {
	const op = "rmqapp.Publish"

	if a.channel == nil {
		return fmt.Errorf("%s: channel is not initialized", op)
	}

	_, err := a.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("%s: queue declare error: %w", op, err)
	}

	err = a.channel.Publish(
		"",
		queueName,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/octet-stream",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
