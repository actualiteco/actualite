package search_transport

import (
	"net/http"

	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/transport/response"
)

// SearchUsers godoc
// @Summary Поиск пользователей
// @Description Возвращает публичные данные пользователей по подстроке в никнейме или отображаемом имени без учёта регистра. Возвращает не более 20 результатов.
// @Tags Search
// @Produce json
// @Param query query string true "Подстрока для поиска"
// @Success 200 {array} UserSearchDTOResponse "Найденные пользователи"
// @Failure 400 {object} domain.CustomError "Неверный запрос"
// @Failure 401 {object} domain.CustomError "Требуется авторизация"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /search/users [get]
// @Security BearerAuth
func (h *SearchHTTPHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	respWriter := response.NewResponseHandler(logger.FromContext(ctx), w)

	users, err := h.service.SearchUsers(ctx, r.URL.Query().Get("query"))
	if err != nil {
		respWriter.MapError(err)
		return
	}

	result := make([]UserSearchDTOResponse, len(users))
	for i, user := range users {
		result[i] = UserSearchDTOResponse{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			AvatarURL:   user.AvatarURL,
		}
	}
	respWriter.JSONResponse(result, http.StatusOK)
}
