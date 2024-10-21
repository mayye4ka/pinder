package websocket

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	public_api "github.com/mayye4ka/pinder-api/api/go"
	notification_api "github.com/mayye4ka/pinder-api/notifications/go"
	"golang.org/x/sync/errgroup"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

const (
	tokenHeader             = "Authorization"
	authorizationTrimPrefix = "Bearer "
)

type WsServer struct {
	auth                 Authenticator
	notificationProducer NotificationProducer
	port                 int

	connStore   map[uint64]map[string]*websocket.Conn
	connStoreMu sync.RWMutex

	httpServer              *http.Server
	finishNotifications     chan struct{}
	finishNotificationsDone chan struct{}
}

type Authenticator interface {
	UnpackToken(ctx context.Context, token string) (uint64, error)
}

type NotificationProducer interface {
	Notifications() <-chan *notification_api.UserNotification
}

func NewWsServer(auth Authenticator, ntfcProducer NotificationProducer, port int) *WsServer {
	return &WsServer{
		auth:                 auth,
		notificationProducer: ntfcProducer,
		port:                 port,

		connStore: map[uint64]map[string]*websocket.Conn{},

		finishNotifications:     make(chan struct{}),
		finishNotificationsDone: make(chan struct{}),
	}
}

func (s *WsServer) addUser(id uint64, conn *websocket.Conn) {
	s.connStoreMu.Lock()
	if s.connStore[id] == nil {
		s.connStore[id] = map[string]*websocket.Conn{}
	}
	connId := uuid.New().String()
	s.connStore[id][connId] = conn
	s.connStoreMu.Unlock()
	go s.serveConn(id, connId, conn)
}

func (s *WsServer) serveConn(id uint64, connId string, conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("ws read error", err)
			conn.Close()
			s.connStoreMu.Lock()
			delete(s.connStore[id], connId)
			if len(s.connStore[id]) == 0 {
				delete(s.connStore, id)
			}
			s.connStoreMu.Unlock()
			break
		}
	}
}

func (s *WsServer) sendBytes(id uint64, bytes []byte) {
	s.connStoreMu.RLock()
	for _, conn := range s.connStore[id] {
		err := conn.WriteMessage(websocket.BinaryMessage, bytes)
		if err != nil {
			log.Println(err)
		}
	}
	s.connStoreMu.RUnlock()
}

func (s *WsServer) notify(userId uint64, notification *public_api.DataPackage) error {
	bytes, err := proto.Marshal(notification)
	if err != nil {
		return nil
	}
	s.sendBytes(userId, bytes)
	return nil
}

var upgrader = websocket.Upgrader{}

func (s *WsServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(tokenHeader) == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token := strings.TrimPrefix(r.Header.Get(tokenHeader), authorizationTrimPrefix)
	userId, err := s.auth.UnpackToken(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	s.addUser(userId, conn)
}

func (s *WsServer) Start(ctx context.Context) error {
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return s.startHttpServer(ectx)
	})
	eg.Go(func() error {
		return s.startNotificationSending(ectx)
	})
	return eg.Wait()
}

func (s *WsServer) Stop(ctx context.Context) error {
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return s.stopHttpServer(ectx)
	})
	eg.Go(func() error {
		return s.stopNotificationSending(ectx)
	})
	return eg.Wait()
}

func (s *WsServer) startNotificationSending(ctx context.Context) error {
	c := s.notificationProducer.Notifications()
	for {
		select {
		case <-ctx.Done():
			close(s.finishNotificationsDone)
			return nil
		case <-s.finishNotifications:
			close(s.finishNotificationsDone)
			return nil
		case n := <-c:
			err := s.notify(n.UserId, n.DataPackage)
			if err != nil {
				return err
			}
		}
	}
}

func (s *WsServer) stopNotificationSending(ctx context.Context) error {
	close(s.finishNotifications)
	select {
	case <-s.finishNotificationsDone:
	case <-ctx.Done():
	}
	return nil
}

func (s *WsServer) startHttpServer(ctx context.Context) error {
	http.HandleFunc("/ws", s.wsHandler)
	s.httpServer = &http.Server{
		Addr: fmt.Sprintf(":%d", s.port),
	}
	go func() {
		<-ctx.Done()
		s.httpServer.Shutdown(ctx)
	}()
	return s.httpServer.ListenAndServe()
}

func (s *WsServer) stopHttpServer(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
