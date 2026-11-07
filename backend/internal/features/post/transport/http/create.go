package post_transport

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/middleware"
	httprequest "actualiteco/actualite/internal/core/transport/request"
	"actualiteco/actualite/internal/core/transport/response"
	"context"
	"fmt"
	"net/http"
	"uuid"

	_ "actualiteco/actualite/internal/core/domain"

	"go.uber.org/zap"
)

// Create godoc
// @Summary Создание поста
// @Description Создает пост длиной <= 280 символов
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body CreatePostDTO true "parent_id, image_url - необязательные параметры"
// @Success 201 {object} PostDTOResponse "Пост создан"
// @Failure 400 {object} domain.CustomError "Неверный запрос (content > 280 символов; userID неправильного формата)"
// @Failure 401 {object} domain.CustomError "Неверные учетные данные"
// @Failure 404 {object} domain.CustomError "Не удалось найти пост с данным parent_id"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /posts/create [post]
// @Security BearerAuth
func (h *PostsHTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	var req CreatePostDTO
	if err := httprequest.DecodeAndValidate(r, &req); err != nil {
		log.Debug("decode or validate create post request", zap.Error(err))
		respWriter.MapError(err)
		return
	}

	userID, err := getUserUUIDFromContext(ctx)
	if err != nil {
		log.Debug("get userID from path value", zap.Error(err))
		respWriter.MapError(err)
		return
	}

	post, err := h.PostsService.Create(ctx, userID, req.Content, req.ParentID, req.ImageURL)
	if err != nil {
		log.Error("create post", zap.Error(err))
		respWriter.MapError(err)
		return
	}

	respWriter.JSONResponse(postToResponse(post), http.StatusCreated)
}

func getUserUUIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID := middleware.GetUserIDFromCtx(ctx)
	if userID == "" {
		return uuid.Nil(), fmt.Errorf("no userID in context: %w", domain.ErrInternal)
	}

	return uuid.Parse(userID)
}
