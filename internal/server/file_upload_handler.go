package server

import (
	"PattyWagon/internal/constants"
	"PattyWagon/observability"
	"fmt"
	"net/http"
)

const (
	SuccessUploadMsg string = "File uploaded sucessfully"
)

func (s *Server) fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := observability.Tracer.Start(r.Context(), "handler.file_upload")
	defer span.End()

	if r.Method != "POST" {
		sendErrorResponse(w, http.StatusBadRequest, "Method not allowed")
		return
	}

	if err := r.ParseMultipartForm(int64(constants.MaxUploadForm)); err != nil { // 3 MB
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error parsing multipart form: %v", err))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}
	defer file.Close()

	uploadedFile, err := s.service.UploadFile(ctx, file, header.Filename, header.Size)
	if err != nil {
		handleServiceError(w, err, nil)
		return
	}

	resp := FileUploadResponse{
		Message: SuccessUploadMsg,
		Data: FileUploadData{
			ImageUrl: uploadedFile.Uri,
		},
	}

	sendResponse(w, http.StatusOK, resp)
	defer r.Body.Close()

}
