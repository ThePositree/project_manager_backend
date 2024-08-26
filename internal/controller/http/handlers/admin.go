package handlers

import (
	"errors"
	"net/http"

	"github.com/ThePositree/billing_manager/internal/controller/http/dto"
	model_billing "github.com/ThePositree/billing_manager/internal/model/billing"
	"github.com/ThePositree/billing_manager/internal/usecase/billing_managing"
	"github.com/ThePositree/billing_manager/internal/usecase/user_managing"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

func PatchBillingNextState(billingManaging billing_managing.BillingManaging, logger zerolog.Logger, password string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With().Str("Handler", "admin/billing/state/next/{id}").Str("Method", "PATCH").Logger()

		_, userPassword, ok := r.BasicAuth()
		if !ok || password != userPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			if err := WriteResponse(
				w,
				http.StatusUnauthorized,
				ResponseMessageDTO{Message: "you are unauthorized"},
			); err != nil {
				logger.Error().Err(err).Msg("Request with incorrect password")
			}
			return
		}

		ctx := r.Context()

		billingId, ok := mux.Vars(r)["id"]
		if !ok {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "billing id in path param not found"},
			); err != nil {
				logger.Error().Err(err).Msg("Request without billing id")
			}
			return
		}

		billing, err := billingManaging.NextState(ctx, billingId)
		if errors.Is(billing_managing.ErrBillingNotFound, err) {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "billing not found"},
			); err != nil {
				logger.Error().Err(err).Msg("Billing not found")
			}
			return
		}
		if errors.Is(model_billing.ErrNextCompletedState{}, err) {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "impossible to next the state from completed state"},
			); err != nil {
				logger.Error().Err(err).Msg("Next from completed state")
			}
			return
		}
		if err != nil {
			logger.Error().Err(err).Msg("Billing managing next state")
			if err := WriteResponse(
				w,
				http.StatusInternalServerError,
				ResponseMessageDTO{Message: "internal server error"},
			); err != nil {
				logger.Error().Err(err).Msg("Internal server error")
			}
			return
		}

		dto := dto.NewBillingDTOFromModel(billing)
		if err := WriteResponse(w, http.StatusOK, dto); err != nil {
			logger.Error().Err(err).Msg("Write OK response")
		}
	}
}

func PatchBillingPrevState(billingManaging billing_managing.BillingManaging, logger zerolog.Logger, password string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With().Str("Handler", "admin/billing/state/prev/{id}").Str("Method", "PATCH").Logger()
		ctx := r.Context()

		_, userPassword, ok := r.BasicAuth()
		if !ok || password != userPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			if err := WriteResponse(
				w,
				http.StatusUnauthorized,
				ResponseMessageDTO{Message: "you are unauthorized"},
			); err != nil {
				logger.Error().Err(err).Msg("Request with incorrect password")
			}
			return
		}

		billingId, ok := mux.Vars(r)["id"]
		if !ok {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "billing id in path param not found"},
			); err != nil {
				logger.Error().Err(err).Msg("Request without billing id")
			}
			return
		}

		billing, err := billingManaging.PrevState(ctx, billingId)
		if errors.Is(billing_managing.ErrBillingNotFound, err) {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "billing not found"},
			); err != nil {
				logger.Error().Err(err).Msg("Billing not found")
			}
			return
		}
		if errors.Is(model_billing.ErrPrevPendingState{}, err) {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "impossible to prev the state from pending state"},
			); err != nil {
				logger.Error().Err(err).Msg("Prev from pending state")
			}
			return
		}
		if err != nil {
			logger.Error().Err(err).Msg("Billing managing prev state")
			if err := WriteResponse(
				w,
				http.StatusInternalServerError,
				ResponseMessageDTO{Message: "internal server error"},
			); err != nil {
				logger.Error().Err(err).Msg("Internal server error")
			}
			return
		}

		dto := dto.NewBillingDTOFromModel(billing)
		if err := WriteResponse(w, http.StatusOK, dto); err != nil {
			logger.Error().Err(err).Msg("Write OK response")
		}
	}
}

func GetAllBillings(billingManaging billing_managing.BillingManaging, logger zerolog.Logger, password string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With().Str("Handler", "admin/billings").Str("Method", "GET").Logger()
		ctx := r.Context()

		_, userPassword, ok := r.BasicAuth()
		if !ok || password != userPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			if err := WriteResponse(
				w,
				http.StatusUnauthorized,
				ResponseMessageDTO{Message: "you are unauthorized"},
			); err != nil {
				logger.Error().Err(err).Msg("Request with incorrect password")
			}
			return
		}

		billings, err := billingManaging.GetAll(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("Billing managing get all")
			if err := WriteResponse(
				w,
				http.StatusInternalServerError,
				ResponseMessageDTO{Message: "internal server error"},
			); err != nil {
				logger.Error().Err(err).Msg("Internal server error")
			}
			return
		}

		var result []dto.Billing
		for _, billing := range billings {
			result = append(result, dto.NewBillingDTOFromModel(billing))
		}
		if err := WriteResponse(w, http.StatusOK, result); err != nil {
			logger.Error().Err(err).Msg("Write OK response")
		}
	}
}

func GetAllUsers(userManaging user_managing.UserManaging, logger zerolog.Logger, password string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With().Str("Handler", "admin/users").Str("Method", "GET").Logger()
		ctx := r.Context()

		_, userPassword, ok := r.BasicAuth()
		if !ok || password != userPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			if err := WriteResponse(
				w,
				http.StatusUnauthorized,
				ResponseMessageDTO{Message: "you are unauthorized"},
			); err != nil {
				logger.Error().Err(err).Msg("Request with incorrect password")
			}
			return
		}

		users, err := userManaging.GetAll(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("Billing managing get all")
			if err := WriteResponse(
				w,
				http.StatusInternalServerError,
				ResponseMessageDTO{Message: "internal server error"},
			); err != nil {
				logger.Error().Err(err).Msg("Internal server error")
			}
			return
		}

		var result []dto.User
		for _, user := range users {
			result = append(result, dto.NewUserDTOFromModel(user))
		}
		if err := WriteResponse(w, http.StatusOK, result); err != nil {
			logger.Error().Err(err).Msg("Write OK response")
		}
	}
}
