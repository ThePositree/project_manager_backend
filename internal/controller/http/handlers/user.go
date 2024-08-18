package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ThePositree/billing_manager/internal/controller/http/dto"
	"github.com/ThePositree/billing_manager/internal/usecase/user_managing"
	"github.com/rs/zerolog"
)

func GetUserByTelegramUN(userManaging user_managing.UserManaging, logger zerolog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With().Str("Handler", "user").Str("Method", "GET").Logger()

		ctx := r.Context()

		queryParams := r.URL.Query()
		ok := queryParams.Has("telegram_username")
		if !ok {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "telegram_username in query param not found"},
			); err != nil {
				logger.Error().Err(err).Msg("Request without telegram_username")
			}
			return
		}

		telegramUN := queryParams.Get("telegram_username")
		user, err := userManaging.GetByTelegramUN(ctx, telegramUN)
		if err != nil {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "user not found"},
			); err != nil {
				logger.Error().Err(err).Msg("User not found")
			}
			return
		}
		dto := dto.NewUserDTOFromModel(user)
		if err := WriteResponse(w, http.StatusOK, dto); err != nil {
			logger.Error().Err(err).Msg("Write OK response")
		}
	}
}

func PostUser(userManaging user_managing.UserManaging, logger zerolog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With().Str("Handler", "user").Str("Method", "POST").Logger()

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

		var userInfo dto.CreateUserInfo

		err = json.Unmarshal(bytes, &userInfo)
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
		if userInfo.TelegramUN == "" {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "telegram_username cannot be empty"},
			); err != nil {
				logger.Error().Err(err).Msg("telegram_username cannot be empty")
			}
			return
		}

		user, err := userManaging.Create(ctx, userInfo.TelegramUN)
		if errors.Is(user_managing.ErrExistingUser, err) {
			if err := WriteResponse(
				w,
				http.StatusBadRequest,
				ResponseMessageDTO{Message: "user is existing"},
			); err != nil {
				logger.Error().Err(err).Msg("User is existing")
			}
			return
		}
		if err != nil {
			logger.Error().Err(err).Msg("User managing create")
			if err := WriteResponse(
				w,
				http.StatusInternalServerError,
				ResponseMessageDTO{Message: "internal server error"},
			); err != nil {
				logger.Error().Err(err).Msg("Internal server error")
			}
			return
		}
		dto := dto.NewUserDTOFromModel(user)
		if err := WriteResponse(w, http.StatusOK, dto); err != nil {
			logger.Error().Err(err).Msg("Write OK response")
		}
	}
}
