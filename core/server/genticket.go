package server

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/graphql-go/handler"
	"github.com/tiny-sky/Tvote/core/dao"
	"github.com/tiny-sky/Tvote/core/dao/entity"
	"github.com/tiny-sky/Tvote/log"
	"github.com/tiny-sky/Tvote/router"
	"github.com/tiny-sky/Tvote/tools"
)

type GenticketSrv struct {
	listenOn       string
	interval       int
	maxTicketUsage int

	httpServer *http.Server
}

func New(settings Settings) *GenticketSrv {
	srv := &GenticketSrv{
		listenOn:       tools.FigureOutListen(settings.ListenOn),
		interval:       settings.Interval,
		maxTicketUsage: settings.MaxTicketUsage,
	}
	return srv
}

func (s *GenticketSrv) Run(ctx context.Context) error {

	go func() {
		generateNewTicket(s.interval, s.maxTicketUsage)
		ticker := time.NewTicker(time.Second * time.Duration(s.interval))
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ticket, _ := generateNewTicket(s.interval, s.maxTicketUsage)
				log.Infof("Create new ticker[%d]: %s", ticket.ID, ticket.Ticket)
			case <-ctx.Done():
				return
			}
		}
	}()

	schema, err := router.CreateSchema()
	if err != nil {
		log.Fatalf("Failed to create GraphQL schema: %v", err)
		return err
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// 设置 HTTP 路由
	mux := http.NewServeMux()
	mux.Handle("/graphql", h)

	// 创建 HTTP 服务器
	s.httpServer = &http.Server{
		Addr:    s.listenOn,
		Handler: mux,
	}

	fmt.Printf("GraphQL server running at http://localhost%s/graphql\n", s.listenOn)

	// 启动 HTTP 服务器
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	<-ctx.Done()

	return s.Stop(ctx)
}

func (s *GenticketSrv) Stop(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

func randString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateNewTicket(interval, maxTicketUsage int) (*entity.Ticket, error) {
	ticketstring := randString(10)

	now := time.Now().Unix()
	expiresAt := now + int64(interval)

	ticket := &entity.Ticket{
		Ticket:    ticketstring,
		CreatedAt: now,
		ExpiresAt: expiresAt,
		MaxUsage:  maxTicketUsage,
		UsedCount: 0,
	}

	if err := dao.GetDB().CreateTicket(ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}
