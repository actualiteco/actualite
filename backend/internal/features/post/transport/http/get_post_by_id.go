package post_transport

import (
	_ "actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/transport/response"
	httputils "actualiteco/actualite/internal/utils/http"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type GetPostByIdDTO struct {
	PostId string `json:"id"`
}

// GetPostByID godoc
// @Summary Поиск поста по ID
// @Description Ищет пост по ID и возвращает информацию о посте + краткую информацию об авторе (username, display_name, avatar_url)
// @Tags Posts
// @Accept json
// @Produce json
// @Param postID path string true "Post ID (UUID)" Format(uuid) Example(3fa85f64-5717-4562-b3fc-2c963f66afa6)
// @Success 200 {object} PostWithAuthorDTOResponse "Пост найден"
// @Failure 400 {object} domain.CustomError "Неверный запрос"
// @Failure 401 {object} domain.CustomError "Не удалось аутентифицировать пользователя"
// @Failure 404 {object} domain.CustomError "Пост не найден"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /posts/{postID} [get]
// @Security BearerAuth
func (h *PostsHTTPHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	postID, err := httputils.GetUUIDPathValue(r, "postID")
	if err != nil {
		respWriter.MapError(err)
		return
	}

	log.Debug("getting post with id", zap.String("postID", postID.String()))

	postWithAuthor, err := h.PostsService.GetPostByID(ctx, postID)
	if err != nil {
		err = fmt.Errorf("get post by id: %w", err)
		respWriter.MapError(err)
		return
	}

	respWriter.JSONResponse(postWithAuthorToResponse(postWithAuthor), http.StatusOK)
}
