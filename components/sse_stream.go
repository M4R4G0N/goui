package components

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// SSEEvent é a mensagem enviada pelo servidor ao cliente via stream.
type SSEEvent struct {
	Data  string
	Event string // nome do evento (opcional; "" usa o padrão "message")
}

// sseStream gerencia todos os clientes conectados a um único componente.
type sseStream struct {
	mu        sync.RWMutex
	clients   map[chan SSEEvent]struct{}
	broadcast chan SSEEvent
	done      chan struct{}
}

func newSSEStream() *sseStream {
	s := &sseStream{
		clients:   make(map[chan SSEEvent]struct{}),
		broadcast: make(chan SSEEvent, 256),
		done:      make(chan struct{}),
	}
	go s.run()
	return s
}

func (s *sseStream) subscribe() chan SSEEvent {
	ch := make(chan SSEEvent, 32)
	s.mu.Lock()
	s.clients[ch] = struct{}{}
	s.mu.Unlock()
	return ch
}

func (s *sseStream) unsubscribe(ch chan SSEEvent) {
	s.mu.Lock()
	if _, ok := s.clients[ch]; ok {
		delete(s.clients, ch)
		close(ch)
	}
	s.mu.Unlock()
}

func (s *sseStream) send(event SSEEvent) {
	select {
	case s.broadcast <- event:
	default:
		// canal cheio — descarta silenciosamente para não bloquear o caller
	}
}

func (s *sseStream) run() {
	heartbeat := time.NewTicker(25 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case event := <-s.broadcast:
			s.mu.RLock()
			for ch := range s.clients {
				select {
				case ch <- event:
				default:
					// cliente lento — ignora este evento para ele
				}
			}
			s.mu.RUnlock()

		case <-heartbeat.C:
			// comentário SSE mantém a conexão viva em proxies/load-balancers
			s.send(SSEEvent{Data: "", Event: "__heartbeat__"})

		case <-s.done:
			s.mu.Lock()
			for ch := range s.clients {
				close(ch)
			}
			s.clients = make(map[chan SSEEvent]struct{})
			s.mu.Unlock()
			return
		}
	}
}

func (s *sseStream) close() {
	close(s.done)
}

// SSERegistry mantém um stream por ID de componente.
type SSERegistry struct {
	mu      sync.RWMutex
	streams map[string]*sseStream
}

var sseRegistry = &SSERegistry{
	streams: make(map[string]*sseStream),
}

// getOrCreate retorna o stream existente ou cria um novo para o ID fornecido.
func (r *SSERegistry) getOrCreate(id string) *sseStream {
	r.mu.RLock()
	s, ok := r.streams[id]
	r.mu.RUnlock()
	if ok {
		return s
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	// double-check após obter write lock
	if s, ok = r.streams[id]; ok {
		return s
	}
	s = newSSEStream()
	r.streams[id] = s
	return s
}

// SSEBroadcast envia um evento de dados para todos os clientes de um componente.
// Pode ser chamado de qualquer goroutine — é thread-safe.
func SSEBroadcast(id, data string) {
	sseRegistry.getOrCreate(id).send(SSEEvent{Data: data})
}

// SSEBroadcastEvent envia um evento nomeado (além de "message").
func SSEBroadcastEvent(id, event, data string) {
	sseRegistry.getOrCreate(id).send(SSEEvent{Event: event, Data: data})
}

func init() {
	// /api/goui/stream?id=X
	// O cliente abre EventSource neste endpoint e recebe atualizações em tempo real.
	http.HandleFunc("/api/goui/stream", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "parâmetro 'id' obrigatório", http.StatusBadRequest)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming não suportado por este servidor", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no") // desativa buffer do nginx

		stream := sseRegistry.getOrCreate(id)
		ch := stream.subscribe()
		defer stream.unsubscribe(ch)

		// confirma conexão imediatamente
		fmt.Fprintf(w, ": connected id=%s\n\n", id)
		flusher.Flush()

		ctx := r.Context()
		for {
			select {
			case <-ctx.Done():
				return

			case event, ok := <-ch:
				if !ok {
					return
				}
				if event.Event == "__heartbeat__" {
					fmt.Fprintf(w, ": heartbeat\n\n")
				} else {
					if event.Event != "" {
						fmt.Fprintf(w, "event: %s\n", event.Event)
					}
					fmt.Fprintf(w, "data: %s\n\n", event.Data)
				}
				flusher.Flush()
			}
		}
	})
}
