package realtime

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"time"
)

const wsGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func computeWSAcceptKey(wsKey string) string {
	hash := sha1.New()
	_, _ = hash.Write([]byte(wsKey + wsGUID))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// readLoop lit les trames WebSocket entrantes du client (ex: LIVRET_EMARGEMENT_LIVE, PING).
func (c *Client) readLoop() {
	defer c.hub.unregister(c)

	for {
		_ = c.conn.SetReadDeadline(time.Now().Add(180 * time.Second))
		opcode, payload, err := c.readFrame()
		if err != nil {
			return
		}

		switch opcode {
		case 0x8: // Close frame
			return
		case 0x9: // Ping frame -> répondre Pong
			_ = c.writeControlFrame(0xA, payload)
		case 0x1: // Text frame JSON
			var incoming struct {
				Type        string      `json:"type"`
				AutoEcoleID string      `json:"auto_ecole_id"`
				Data        interface{} `json:"data"`
			}
			if err := json.Unmarshal(payload, &incoming); err == nil && incoming.Type != "" {
				if incoming.Type == "PING" {
					pong, _ := json.Marshal(EventMessage{
						Type:        "PONG",
						AutoEcoleID: c.autoEcoleID,
						Timestamp:   time.Now().UTC().Format(time.RFC3339Nano),
					})
					_ = c.writeTextFrame(pong)
					continue
				}
				// Relayer en temps réel strictement au tenant authentifié de la session active
				c.hub.BroadcastToTenant(c.autoEcoleID, incoming.Type, incoming.Data)
			}
		}
	}
}

// readFrame décode une trame WebSocket cliente masquée (RFC 6455).
func (c *Client) readFrame() (byte, []byte, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(c.rw, header); err != nil {
		return 0, nil, err
	}

	opcode := header[0] & 0x0F
	masked := (header[1] & 0x80) != 0
	payloadLen := int64(header[1] & 0x7F)

	switch payloadLen {
	case 126:
		ext := make([]byte, 2)
		if _, err := io.ReadFull(c.rw, ext); err != nil {
			return 0, nil, err
		}
		payloadLen = int64(binary.BigEndian.Uint16(ext))
	case 127:
		ext := make([]byte, 8)
		if _, err := io.ReadFull(c.rw, ext); err != nil {
			return 0, nil, err
		}
		payloadLen = int64(binary.BigEndian.Uint64(ext))
	}

	if payloadLen > 10*1024*1024 {
		return 0, nil, io.ErrShortBuffer
	}

	var maskKey [4]byte
	if masked {
		if _, err := io.ReadFull(c.rw, maskKey[:]); err != nil {
			return 0, nil, err
		}
	}

	payload := make([]byte, payloadLen)
	if payloadLen > 0 {
		if _, err := io.ReadFull(c.rw, payload); err != nil {
			return 0, nil, err
		}
		if masked {
			for i := range payload {
				payload[i] ^= maskKey[i%4]
			}
		}
	}

	return opcode, payload, nil
}

// writeTextFrame encode et envoie une trame WebSocket texte non-masquée du serveur vers le client.
func (c *Client) writeTextFrame(payload []byte) error {
	return c.writeFrame(0x1, payload)
}

func (c *Client) writeControlFrame(opcode byte, payload []byte) error {
	return c.writeFrame(opcode, payload)
}

func (c *Client) writeFrame(opcode byte, payload []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	if c.closed {
		return net.ErrClosed
	}

	_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	length := len(payload)
	var header []byte
	header = append(header, 0x80|opcode) // FIN bit + opcode

	if length <= 125 {
		header = append(header, byte(length))
	} else if length <= 65535 {
		header = append(header, 126, byte(length>>8), byte(length))
	} else {
		header = append(header, 127,
			byte(uint64(length)>>56),
			byte(uint64(length)>>48),
			byte(uint64(length)>>40),
			byte(uint64(length)>>32),
			byte(uint64(length)>>24),
			byte(uint64(length)>>16),
			byte(uint64(length)>>8),
			byte(uint64(length)),
		)
	}

	if _, err := c.rw.Write(header); err != nil {
		return err
	}
	if length > 0 {
		if _, err := c.rw.Write(payload); err != nil {
			return err
		}
	}
	return c.rw.Flush()
}

