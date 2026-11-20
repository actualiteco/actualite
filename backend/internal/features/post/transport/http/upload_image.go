package post_transport

import (
	"actualiteco/actualite/internal/core/domain"
	"actualiteco/actualite/internal/core/logger"
	httprequest "actualiteco/actualite/internal/core/transport/request"
	"actualiteco/actualite/internal/core/transport/response"
	"net/http"
)

const multipartOverhead = 512 << 10

// UploadImage godoc
// @Summary Загрузить изображение
// @Description Загрузить изображение для поста (jpg или png, до 5 МБ) в хранилище
// @Tags Posts
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл изображения (jpg/png, до 5 МБ)"
// @Success 200 {object} imageUrlDTOResponse "Успешный запрос"
// @Failure 400 {object} domain.CustomError "Неверный запрос (например, файл пустой)"
// @Failure 401 {object} domain.CustomError "Пользователь неавторизован"
// @Failure 403 {object} domain.CustomError "Для действия требуется подтверждение почты"
// @Failure 413 {object} domain.CustomError "Размер файла превышает допустимый"
// @Failure 415 {object} domain.CustomError "Неподходящий формат файла"
// @Failure 500 {object}  domain.InternalError "Внутренняя ошибка сервера"
// @Router /posts/image [post]
// @Security BearerAuth
func (h *PostsHTTPHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	respWriter := response.NewResponseHandler(log, w)

	r.Body = http.MaxBytesReader(w, r.Body, domain.MaxImageSize+multipartOverhead)

	data, err := httprequest.GetFileBytes(r, "file", domain.MaxImageSize)
	if err != nil {
		respWriter.MapError(err)
		return
	}

	imageURL, err := h.postsService.UploadImage(ctx, data)
	if err != nil {
		respWriter.MapError(err)
		return
	}

	respWriter.JSONResponse(imageUrlDTOResponse{ImageURL: imageURL}, http.StatusOK)
}
