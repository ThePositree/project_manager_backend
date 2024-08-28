package http_controller

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/ThePositree/billing_manager/internal/controller/http/handlers"
	"github.com/ThePositree/billing_manager/internal/usecase/billing_managing"
	"github.com/ThePositree/billing_manager/internal/usecase/user_managing"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type handlerInfo struct {
	handler http.HandlerFunc
	path    string
	method  string
}

type http_controller struct {
	port            int
	logger          zerolog.Logger
	billingManaging billing_managing.BillingManaging
	userManaging    user_managing.UserManaging
	adminPassword   string
}

func (hc http_controller) Start(ctx context.Context) {
	handlersInfo := []handlerInfo{
		{
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			},
			path:   "/ping",
			method: http.MethodGet,
		},
		{
			handler: handlers.GetAllBillings(hc.billingManaging, hc.logger, hc.adminPassword),
			path:    "/admin/billings",
			method:  http.MethodGet,
		},
		{
			handler: handlers.GetAllUsers(hc.userManaging, hc.logger, hc.adminPassword),
			path:    "/admin/users",
			method:  http.MethodGet,
		},
		{
			handler: handlers.GetBilling(hc.billingManaging, hc.logger),
			path:    "/billing",
			method:  http.MethodGet,
		},
		{
			handler: handlers.GetUserByTelegramUN(hc.userManaging, hc.logger),
			path:    "/user",
			method:  http.MethodGet,
		},
		{
			handler: handlers.PatchBilling(hc.billingManaging, hc.logger),
			path:    "/billing/{id}",
			method:  http.MethodPatch,
		},
		{
			handler: handlers.PatchBillingNextState(hc.billingManaging, hc.logger, hc.adminPassword),
			path:    "/admin/billing/state/next/{id}",
			method:  http.MethodPatch,
		},
		{
			handler: handlers.PatchBillingPrevState(hc.billingManaging, hc.logger, hc.adminPassword),
			path:    "/admin/billing/state/prev/{id}",
			method:  http.MethodPatch,
		},
		{
			handler: handlers.PostBilling(hc.billingManaging, hc.logger),
			path:    "/billing",
			method:  http.MethodPost,
		},
		{
			handler: handlers.PostUser(hc.userManaging, hc.logger),
			path:    "/user",
			method:  http.MethodPost,
		},
	}

	r := mux.NewRouter()
	for _, handlerInfo := range handlersInfo {
		r.Handle(handlerInfo.path, handlerInfo.handler).Methods(handlerInfo.method)
	}
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", hc.port),
		Handler: r,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
	go func() {
		<-ctx.Done()
		if err := httpServer.Shutdown(context.Background()); err != nil {
			hc.logger.Error().Err(err).Msg("Http server shutdown")
		}
	}()
	if err := httpServer.ListenAndServe(); err != nil {
		hc.logger.Error().Err(err).Msg("Http listen and serve")
	}
}

func New(
	logger zerolog.Logger,
	billingManaging billing_managing.BillingManaging,
	userManaging user_managing.UserManaging,
	port int,
	adminPassword string,
) http_controller {
	return http_controller{
		port:            port,
		logger:          logger,
		billingManaging: billingManaging,
		userManaging:    userManaging,
		adminPassword:   adminPassword,
	}
}
