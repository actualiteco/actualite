package profile_transport

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	"actualiteco/actualite/internal/core/transport/response"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

// запас на служебные данные multipart поверх самого файла
const multipartOverhead = 512 << 10

// UploadAvatar godoc
// @Summary Загрузить аватар
// @Description Загрузить аватар пользователя (jpg или png, до 2 МБ) в хранилище и обновить профиль
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл аватара (jpg/png, до 2 МБ)"
// @Success 200 {object} UploadAvatarResponse "Успешный запрос"
// @Failure 400 {object} domain.CustomError "Неверный запрос (неподходящий формат или размер файла)"
// @Failure 401 {object} domain.CustomError "Пользователь неавторизован"
// @Failure 500 {object}  domain.InternalError "Внутренняя ошибка сервера"
// @Router /profile/avatar [post]
// @Security BearerAuth
func (h *ProfileHTTPHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	userID, err := getUserUUIDFromContext(ctx)
	if err != nil {
		log.Error("get userID from context", zap.Error(err))
		respWriter.MapError(err)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, domain.MaxAvatarSize+multipartOverhead)

	file, _, err := r.FormFile("file")
	if err != nil {
		log.Debug("failed to get file from form", zap.Error(err))
		respWriter.MapError(fmt.Errorf("get file from form: %w: %w", err, domain.ErrInvalidRequest))
		return
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, domain.MaxAvatarSize+1))
	if err != nil {
		log.Debug("failed to read file", zap.Error(err))
		respWriter.MapError(fmt.Errorf("read file: %w: %w", err, domain.ErrInvalidRequest))
		return
	}

	avatarURL, err := h.service.UploadAvatar(ctx, userID, data)
	if err != nil {
		log.Error("failed to upload avatar", zap.Error(err))
		respWriter.MapError(err)
		return
	}

	respWriter.JSONResponse(UploadAvatarResponse{AvatarURL: avatarURL}, http.StatusOK)
}
