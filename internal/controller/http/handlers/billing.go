package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ThePositree/billing_manager/internal/controller/http/dto"
	model_billing "github.com/ThePositree/billing_manager/internal/model/billing"
	"github.com/ThePositree/billing_manager/internal/usecase/billing_managing"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

func GetBilling(billingManaging billing_managing.BillingManaging, logger zerolog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With().Str("Handler", "billing").Str("Method", "GET").Logger()

		ctx := r.Context()

		queryParams := r.URL.Query()
		ok := queryParams.Has("user_id")
		if !ok {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "user_id in query param not found"},
			); err != nil {
				logger.Error().Err(err).Msg("Request without user_id")
			}
			return
		}

		userID := queryParams.Get("user_id")
		billings, err := billingManaging.GetAllByUserId(ctx, userID)
		if errors.Is(billing_managing.ErrNoUser, err) {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "user not found"},
			); err != nil {
				logger.Error().Err(err).Msg("User not found")
			}
			return
		}
		if err != nil {
			logger.Error().Err(err).Msg("Billing managing get all by user id")
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
			billingDTO := dto.NewBillingDTOFromModel(billing)
			result = append(result, billingDTO)
		}
		if err := WriteResponse(w, http.StatusOK, result); err != nil {
			logger.Error().Err(err).Msg("Write OK response")
		}
	}
}

func PostBilling(billingManaging billing_managing.BillingManaging, logger zerolog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With().Str("Handler", "billing").Str("Method", "POST").Logger()
		ctx := r.Context()

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error().Err(err).Msg("Read body")
			if err := WriteResponse(
				w,
				http.StatusInternalServerError,
				ResponseMessageDTO{Message: "internal server error"},
			); err != nil {
				logger.Error().Err(err).Msg("Internal server error")
			}
			return
		}

		var billingInfo dto.CreateBillingInfo

		err = json.Unmarshal(bytes, &billingInfo)
		if err != nil {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: fmt.Sprintf("wrong structure body: %s", err.Error())},
			); err != nil {
				logger.Error().Err(err).Msg("Json unmarshal")
			}
			return
		}
		if billingInfo.UserId == "" {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "user_id cannot be empty"},
			); err != nil {
				logger.Error().Err(err).Msg("user_id cannot be empty")
			}
			return
		}

		billing, err := billingManaging.Create(ctx, billingInfo.UserId)
		if errors.Is(billing_managing.ErrNoUser, err) {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "user not found"},
			); err != nil {
				logger.Error().Err(err).Msg("User not found")
			}
			return
		}
		if err != nil {
			logger.Error().Err(err).Msg("Billing managing create")
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

func PatchBilling(billingManaging billing_managing.BillingManaging, logger zerolog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With().Str("Handler", "billing/{id}").Str("Method", "PATCH").Logger()
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

		billing, err := billingManaging.GetById(ctx, billingId)
		if errors.Is(billing_managing.ErrNoBilling, err) {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "billing not found"},
			); err != nil {
				logger.Error().Err(err).Msg("Billing not found")
			}
			return
		}
		if err != nil {
			logger.Error().Err(err).Msg("Billing managing get by id")
			if err := WriteResponse(
				w,
				http.StatusInternalServerError,
				ResponseMessageDTO{Message: "internal server error"},
			); err != nil {
				logger.Error().Err(err).Msg("Internal server error")
			}
			return
		}

		if billing.GetBriefInfo().Username != "" {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "brief info already existing"},
			); err != nil {
				logger.Error().Err(err).Msg("Brief info already existing")
			}
			return
		}

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error().Err(err).Msg("Read body")
			if err := WriteResponse(
				w,
				http.StatusInternalServerError,
				ResponseMessageDTO{Message: "internal server error"},
			); err != nil {
				logger.Error().Err(err).Msg("Internal server error")
			}
			return
		}

		var briefInfo dto.BriefInfo

		err = json.Unmarshal(bytes, &briefInfo)
		if err != nil {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: fmt.Sprintf("wrong structure body: %s", err.Error())},
			); err != nil {
				logger.Error().Err(err).Msg("Json unmarshal")
			}
			return
		}
		if briefInfo.Username == "" {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "username cannot be empty"},
			); err != nil {
				logger.Error().Err(err).Msg("username cannot be empty")
			}
			return
		}

		billing, err = billingManaging.SetBriefInfo(ctx, billingId, briefInfo.Username)
		if err != nil {
			logger.Error().Err(err).Msg("Billing managing set brief info")
			if err := WriteResponse(
				w,
				http.StatusInternalServerError,
				ResponseMessageDTO{Message: "internal server error"},
			); err != nil {
				logger.Error().Err(err).Msg("Internal server error")
			}
			return
		}

		billing, err = billingManaging.NextState(ctx, billingId)
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
