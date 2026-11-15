package profile_transport

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/transport/response"
	httputils "actualiteco/actualite/internal/utils/http"
	"fmt"
	"net/http"
)

// CheckUsernameAvailability godoc
// @Summary Проверка доступности username
// @Description Проверить, занят ли username
// @Tags Profile
// @Accept json
// @Produce json
// @Param username query string true "Username"
// @Success 200 {object} FollowResponse "Username свободен"
// @Failure 400 {object} domain.CustomError "Неверный запрос (username не валиден или не указан)"
// @Failure 409 {object} domain.CustomError "Username занят"
// @Failure 500 {object}  domain.InternalError "Внутренняя ошибка сервера"
// @Router /profile/username-availability [get]
func (h *ProfileHTTPHandler) CheckUsernameAvailability(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	username := httputils.GetStringQueryParam(r, "username")
	if username == nil {
		err := fmt.Errorf("no username in query: %w", domain.ErrInvalidRequest)
		respWriter.MapError(err)
		return
	}

	isTaken, err := h.profileService.IsUsernameTaken(ctx, *username)
	if err != nil {
		respWriter.MapError(err)
		return
	}

	if isTaken {
		err = fmt.Errorf("username is not available: %w", domain.ErrUsernameAlreadyUsed)
		respWriter.MapError(err)
		return
	}

	respWriter.JSONResponse(CheckUsernameAvailabilityResponse{Available: true}, http.StatusOK)
}
