package message_broker

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/AnatolyKoltun/calculator-storage/models"
	"github.com/AnatolyKoltun/calculator-storage/repositories"
	"github.com/nats-io/nats.go"
)

type HandleNats struct {
	nc *nats.Conn
}

var calcRepository = repositories.CalculationRepository{}

func (n *HandleNats) CreateStreamNats() {
	var err error
	// 1. Подключение к NATS
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	n.nc, err = nats.Connect(natsURL)

	if err != nil {
		log.Fatal("Ошибка подключения к NATS:", err)
	}

	js, err := n.nc.JetStream()

	if err != nil {
		log.Fatal("Ошибка JetStream:", err)
	}

	// 2. Создаем/проверяем Stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "CALCULATIONS",
		Subjects: []string{"calculations.*"},
		Storage:  nats.FileStorage,
	})

	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.Fatal("Ошибка создания Stream:", err)
	}

	// 3. Подписываемся на сообщения (durable consumer)
	_, err = js.Subscribe("calculations.*", func(msg *nats.Msg) {
		log.Printf("Получено сообщение: %s", string(msg.Data))

		// Парсим JSON
		var data *models.Calculation

		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("Ошибка парсинга: %v", err)
			msg.Nak() // Не подтверждаем, NATS отправит повторно
			return
		}

		ctx := context.Background()
		err = calcRepository.Save(ctx, data)

		if err == nil {
			log.Printf("Успешно сохранено: %+v", data)
			msg.Ack()
		} else {
			log.Printf("Ошибка сохранения: %v", err)
			msg.Nak()
		}
	}, nats.Durable("storage-consumer"), nats.ManualAck())

	if err != nil {
		log.Fatal("Ошибка подписки:", err)
	}

	log.Println("NATS consumer запущен")
}

func (n *HandleNats) Close() {
	n.nc.Close()
}
