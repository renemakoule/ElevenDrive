package realtime

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"autoschool/backend/internal/models"
)

// EventMessage représente un événement temps réel diffusé sur le WebSocket.
type EventMessage struct {
	Type        string      `json:"type"`
	AutoEcoleID string      `json:"auto_ecole_id"`
	Timestamp   string      `json:"timestamp"`
	Data        interface{} `json:"data"`
}

// Client représente une connexion WebSocket active rattachée à une auto-école.
type Client struct {
	hub         *Hub
	conn        net.Conn
	rw          *bufio.ReadWriter
	writeMu     sync.Mutex
	autoEcoleID string
	userID      string
	role        string
	closed      bool
}

// Hub orchestre les salons multi-tenant WebSocket en temps réel.
type Hub struct {
	mu        sync.RWMutex
	rooms     map[string]map[*Client]struct{}
	jwtSecret string
}

var defaultHub *Hub

// NewHub initialise le Hub WebSocket temps réel.
func NewHub(jwtSecret string) *Hub {
	h := &Hub{
		rooms:     make(map[string]map[*Client]struct{}),
		jwtSecret: jwtSecret,
	}
	defaultHub = h
	return h
}

// GetDefaultHub renvoie l'instance globale du Hub WebSocket.
func GetDefaultHub() *Hub {
	return defaultHub
}

// Broadcast diffuse instantanément un événement temps réel à tous les écrans connectés d'une auto-école.
func Broadcast(autoEcoleID string, eventType string, data interface{}) {
	if defaultHub == nil || strings.TrimSpace(autoEcoleID) == "" {
		return
	}
	defaultHub.BroadcastToTenant(autoEcoleID, eventType, data)
}

// BroadcastToSuperAdmin diffuse un événement temps réel à tous les super-administrateurs connectés
func BroadcastToSuperAdmin(eventType string, data interface{}) {
	if defaultHub == nil {
		return
	}
	defaultHub.BroadcastToTenant("SUPER_ADMIN_ROOM", eventType, data)
}

// BroadcastToRole diffuse un événement à un rôle spécifique d'une auto-école (ex: directeurs uniquement)
func BroadcastToRole(autoEcoleID string, targetRole string, eventType string, data interface{}) {
	if defaultHub == nil || strings.TrimSpace(autoEcoleID) == "" {
		return
	}
	defaultHub.BroadcastToRoleInTenant(autoEcoleID, targetRole, eventType, data)
}

// BroadcastToUser diffuse un événement à un utilisateur précis (ex: caissier ou moniteur)
func BroadcastToUser(autoEcoleID string, targetUserID string, eventType string, data interface{}) {
	if defaultHub == nil || strings.TrimSpace(autoEcoleID) == "" || strings.TrimSpace(targetUserID) == "" {
		return
	}
	defaultHub.BroadcastToUserInTenant(autoEcoleID, targetUserID, eventType, data)
}

// BroadcastGlobal diffuse un événement à absolument tous les clients connectés
func BroadcastGlobal(eventType string, data interface{}) {
	if defaultHub == nil {
		return
	}
	defaultHub.BroadcastToAll(eventType, data)
}

func marshalEventPayload(autoEcoleID, eventType string, data interface{}) ([]byte, bool) {
	msg := EventMessage{
		Type:        eventType,
		AutoEcoleID: autoEcoleID,
		Timestamp:   time.Now().UTC().Format(time.RFC3339Nano),
		Data:        data,
	}
	payload, err := json.Marshal(msg)
	return payload, err == nil
}

func (h *Hub) sendToTargets(targets []*Client, payload []byte) {
	for _, client := range targets {
		if err := client.writeTextFrame(payload); err != nil {
			h.unregister(client)
		}
	}
}

func (h *Hub) collectTenantTargets(autoEcoleID string, filter func(*Client) bool) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clientsMap, ok := h.rooms[autoEcoleID]
	if !ok || len(clientsMap) == 0 {
		return nil
	}
	targets := make([]*Client, 0, len(clientsMap))
	for c := range clientsMap {
		if filter == nil || filter(c) {
			targets = append(targets, c)
		}
	}
	return targets
}

// BroadcastToTenant envoie un message JSON à tous les clients WebSocket d'un tenant (auto_ecole_id).
func (h *Hub) BroadcastToTenant(autoEcoleID string, eventType string, data interface{}) {
	payload, ok := marshalEventPayload(autoEcoleID, eventType, data)
	if !ok {
		return
	}
	if targets := h.collectTenantTargets(autoEcoleID, nil); len(targets) > 0 {
		h.sendToTargets(targets, payload)
	}
}

// BroadcastToRoleInTenant filtre la diffusion par rôle au sein d'une auto-école
func (h *Hub) BroadcastToRoleInTenant(autoEcoleID string, targetRole string, eventType string, data interface{}) {
	payload, ok := marshalEventPayload(autoEcoleID, eventType, data)
	if !ok {
		return
	}
	targetRoleNorm := models.NormalizeRole(targetRole)
	isBroadcastAll := targetRoleNorm == "" || strings.ToUpper(strings.TrimSpace(targetRole)) == "ALL"

	targets := h.collectTenantTargets(autoEcoleID, func(c *Client) bool {
		return isBroadcastAll || models.NormalizeRole(c.role) == targetRoleNorm
	})
	if len(targets) > 0 {
		h.sendToTargets(targets, payload)
	}
}

// BroadcastToUserInTenant cible un collaborateur précis
func (h *Hub) BroadcastToUserInTenant(autoEcoleID string, targetUserID string, eventType string, data interface{}) {
	payload, ok := marshalEventPayload(autoEcoleID, eventType, data)
	if !ok {
		return
	}
	targets := h.collectTenantTargets(autoEcoleID, func(c *Client) bool {
		return c.userID == targetUserID
	})
	if len(targets) > 0 {
		h.sendToTargets(targets, payload)
	}
}

// BroadcastToAll diffuse à tous les salons
func (h *Hub) BroadcastToAll(eventType string, data interface{}) {
	payload, ok := marshalEventPayload("GLOBAL", eventType, data)
	if !ok {
		return
	}

	h.mu.RLock()
	var targets []*Client
	for _, clientsMap := range h.rooms {
		for c := range clientsMap {
			targets = append(targets, c)
		}
	}
	h.mu.RUnlock()

	h.sendToTargets(targets, payload)
}

func (h *Hub) register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.rooms[c.autoEcoleID]; !ok {
		h.rooms[c.autoEcoleID] = make(map[*Client]struct{})
	}
	h.rooms[c.autoEcoleID][c] = struct{}{}
}

func (h *Hub) unregister(c *Client) {
	h.mu.Lock()
	if clients, ok := h.rooms[c.autoEcoleID]; ok {
		delete(clients, c)
		if len(clients) == 0 {
			delete(h.rooms, c.autoEcoleID)
		}
	}
	h.mu.Unlock()

	c.writeMu.Lock()
	if !c.closed {
		c.closed = true
		_ = c.conn.Close()
	}
	c.writeMu.Unlock()
}

func validateAccessToken(tokenString string, secret string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode de signature inattendue: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("jeton d'accès invalide")
	}
	return claims, nil
}

// ServeWS gère la mise à niveau HTTP -> WebSocket (RFC 6455) et maintient le flux temps réel.
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		http.Error(w, "Mise à niveau WebSocket requise", http.StatusBadRequest)
		return
	}

	wsKey := strings.TrimSpace(r.Header.Get("Sec-WebSocket-Key"))
	if wsKey == "" {
		http.Error(w, "Sec-WebSocket-Key manquant", http.StatusBadRequest)
		return
	}

	token := strings.TrimSpace(r.URL.Query().Get("token"))
	if token == "" {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			token = strings.TrimSpace(authHeader[7:])
		}
	}
	if token == "" {
		http.Error(w, "Connexion WebSocket refusée : jeton d'authentification Bearer ou token requis", http.StatusUnauthorized)
		return
	}

	claims, err := validateAccessToken(token, h.jwtSecret)
	if err != nil || claims == nil {
		http.Error(w, "Connexion WebSocket refusée : jeton JWT invalide ou expiré", http.StatusUnauthorized)
		return
	}

	role := strings.ToUpper(strings.TrimSpace(claims.Role))
	autoEcoleID := claims.AutoEcoleID
	userID := claims.UserID

	if role == "SUPER_ADMIN" {
		autoEcoleID = "SUPER_ADMIN_ROOM"
	} else if strings.TrimSpace(autoEcoleID) == "" {
		http.Error(w, "Connexion WebSocket refusée : auto_ecole_id manquant", http.StatusUnauthorized)
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking non supporté par le serveur", http.StatusInternalServerError)
		return
	}

	conn, rw, err := hijacker.Hijack()
	if err != nil {
		log.Printf("⚠️ Erreur Hijack WebSocket: %v\n", err)
		return
	}

	acceptKey := computeWSAcceptKey(wsKey)
	responseHeaders := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Accept: " + acceptKey + "\r\n\r\n"

	if _, err := rw.WriteString(responseHeaders); err != nil {
		_ = conn.Close()
		return
	}
	if err := rw.Flush(); err != nil {
		_ = conn.Close()
		return
	}

	client := &Client{
		hub:         h,
		conn:        conn,
		rw:          rw,
		autoEcoleID: autoEcoleID,
		userID:      userID,
		role:        role,
	}

	h.register(client)

	welcome, _ := json.Marshal(EventMessage{
		Type:        "WS_CONNECTED",
		AutoEcoleID: autoEcoleID,
		Timestamp:   time.Now().UTC().Format(time.RFC3339Nano),
		Data: map[string]interface{}{
			"status":        "LIVE",
			"auto_ecole_id": autoEcoleID,
			"user_id":       userID,
			"role":          role,
		},
	})
	_ = client.writeTextFrame(welcome)

	go client.readLoop()
}
