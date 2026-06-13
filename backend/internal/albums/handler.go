package albums

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/YardenDrori/music-platform/internal/apperrors"
	httputils "github.com/YardenDrori/music-platform/internal/http_utils"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) NewAlbum(w http.ResponseWriter, r *http.Request) error {
	req := &NewAlbumRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return apperrors.NewErrBadRequest("invalid request body")
	}
	if err := h.service.NewAlbum(r.Context(), req); err != nil {
		return err
	}
	return httputils.WriteResponse(w, http.StatusCreated, nil)
}

func (h *handler) GetAlbumsByName(w http.ResponseWriter, r *http.Request) error {
	name := r.URL.Query().Get("name")
	if name == "" {
		return apperrors.NewErrBadRequest("query missing name field")
	}
	albums, err := h.service.GetAlbumsByName(r.Context(), name)
	if err != nil {
		return err
	}
	return httputils.WriteResponse(w, http.StatusOK, albums)
}

func (h *handler) GetAlbumByID(w http.ResponseWriter, r *http.Request) error {
	idRaw := r.URL.Query().Get("id")
	if idRaw == "" {
		return apperrors.NewErrBadRequest("query missing id field")
	}
	id, err := uuid.Parse(idRaw)
	if err != nil {
		return apperrors.NewErrBadRequest("id provided is not a valid UUID")
	}
	album, err := h.service.GetAlbumByID(r.Context(), id)
	if err != nil {
		return err
	}
	return httputils.WriteResponse(w, http.StatusOK, album)
}

func (h *handler) UpdateAlbumDetails(w http.ResponseWriter, r *http.Request) error {
	albumIDRaw := r.PathValue("id")
	if albumIDRaw == "" {
		return apperrors.NewErrInternal().
			WithInternal("user accessed UpdateAlbumDetails handler without \"id\" in the url path")
	}
	albumID, err := uuid.Parse(albumIDRaw)
	if err != nil {
		return apperrors.NewErrBadRequest("invalid album id").WithCause(err)
	}

	req := &UpdateAlbumRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return apperrors.NewErrBadRequest("invalid request body")
	}
	req.ID = albumID

	if err := h.service.UpdateAlbumDetails(r.Context(), req); err != nil {
		return err
	}
	return httputils.WriteResponse(w, http.StatusOK, nil)
}

func (h *handler) SoftDeleteAlbum(w http.ResponseWriter, r *http.Request) error {
	idRaw := r.URL.Query().Get("id")
	if idRaw == "" {
		return apperrors.NewErrBadRequest("query missing id field")
	}
	id, err := uuid.Parse(idRaw)
	if err != nil {
		return apperrors.NewErrBadRequest("id provided is not a valid UUID")
	}
	if err := h.service.SoftDeleteAlbum(r.Context(), id); err != nil {
		return err
	}
	return httputils.WriteResponse(w, http.StatusOK, nil)
}

func (h *handler) HardDeleteAlbum(w http.ResponseWriter, r *http.Request) error {
	idRaw := r.URL.Query().Get("id")
	if idRaw == "" {
		return apperrors.NewErrBadRequest("query missing id field")
	}
	id, err := uuid.Parse(idRaw)
	if err != nil {
		return apperrors.NewErrBadRequest("id provided is not a valid UUID")
	}
	if err := h.service.HardDeleteAlbum(r.Context(), id); err != nil {
		return err
	}
	return httputils.WriteResponse(w, http.StatusOK, nil)
}

const maxAlbumArtBytes = 10 * 1024 * 1024 // 10MB

func (h *handler) UploadAlbumPicture(w http.ResponseWriter, r *http.Request) error {
	albumIDRaw := r.PathValue("id")
	if albumIDRaw == "" {
		return apperrors.NewErrInternal().
			WithInternal("user accessed UploadAlbumPicture handler without \"id\" in the url path")
	}
	albumID, err := uuid.Parse(albumIDRaw)
	if err != nil {
		return apperrors.NewErrBadRequest("invalid album id").WithCause(err)
	}

	pictureBytes, err := io.ReadAll(io.LimitReader(r.Body, maxAlbumArtBytes+1))
	if err != nil {
		return apperrors.NewErrInternal().WithCause(err)
	}
	if len(pictureBytes) > maxAlbumArtBytes {
		return apperrors.NewErrBadRequest(fmt.Sprintf("file size over %d bytes", maxAlbumArtBytes))
	}

	if err := h.service.UploadAlbumPicture(r.Context(), pictureBytes, albumID); err != nil {
		return err
	}
	return httputils.WriteResponse(w, http.StatusOK, nil)
}
