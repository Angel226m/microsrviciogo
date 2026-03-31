// ═══════════════════════════════════════════════════════════════
// Adaptador de Mensajería NATS – Publicación asíncrona de eventos de dominio
// Implementa port.EventPublisher usando NATS como bus de mensajes
// ═══════════════════════════════════════════════════════════════
package messaging

import (
	"context"
	"encoding/json"

	"github.com/cloudmart/user-service/internal/domain/port"
	"github.com/nats-io/nats.go"
)

type natsPublisher struct {
	conn *nats.Conn
}

// NewNATSPublisher crea un nuevo publicador de eventos NATS implementando port.EventPublisher.
func NewNATSPublisher(conn *nats.Conn) port.EventPublisher {
	return &natsPublisher{conn: conn}
}

func (p *natsPublisher) Publish(ctx context.Context, subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.conn.Publish(subject, payload)
}
