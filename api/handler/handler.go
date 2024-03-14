package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"retask/api/model"
	"retask/config"

	"github.com/sirupsen/logrus"
)

var (
	ErrNoPackages             = fmt.Errorf("provided packages are empty")
	ErrPackagesHaveDuplicates = fmt.Errorf("provided packages have duplicates")
	ErrOrderInvalid           = fmt.Errorf("provided order is negative or zero")
	ErrInternalServerError    = fmt.Errorf("internal server error")
)

type PackagingRepo interface {
	Calculate(int, []int) []int
}

type Handler struct {
	conf        *config.Config
	packageRepo PackagingRepo
}

func New(conf *config.Config, pr PackagingRepo) *Handler {
	return &Handler{
		conf:        conf,
		packageRepo: pr,
	}
}

func (h *Handler) CalculateBestPackages(rw http.ResponseWriter, req *http.Request) {
	reqID, ok := req.Context().Value("request_id").(string)
	if !ok {
		reqID = "failed_to_fetch"
	}
	logger := logrus.WithField("request_id", reqID)
	logger.Debug("fetched request id")

	r := &model.CalculateBestPackagesRequest{}
	if err := parseRequest(req, r); err != nil {
		logger.WithField("error", err).Error("failed to parse request")
	}

	if r.Order <= 0 {
		writeResponse(rw, 400, nil, ErrOrderInvalid, logger)
		return
	}

	packs := h.packageRepo.Calculate(r.Order, h.conf.GetPacks())

	res := &model.CalculateBestPackagesResponse{
		Packages: packs,
	}

	writeResponse(rw, 200, res, nil, logger)
}

func (h *Handler) UpdatePackageSizes(rw http.ResponseWriter, req *http.Request) {
	reqID, ok := req.Context().Value("request_id").(string)
	if !ok {
		reqID = "failed_to_fetch"
	}
	logger := logrus.WithField("request_id", reqID)
	logger.Debug("fetched request id")

	r := &model.UpdatePackageSizes{}
	if err := parseRequest(req, r); err != nil {
		logger.WithField("error", err).Error("failed to parse request")
	}

	if len(r.Sizes) == 0 {
		writeResponse(rw, 400, nil, ErrNoPackages, logger)
		return
	}

	checkDuplicates := make(map[int]int)
	for _, item := range r.Sizes {
		_, ok := checkDuplicates[item]
		if ok {
			writeResponse(rw, 400, nil, ErrPackagesHaveDuplicates, logger)
			return
		}

		checkDuplicates[item] += 1
	}

	packs := h.conf.SetPacks(r.Sizes)

	res := &model.UpdatePackageSizes{
		Sizes: packs,
	}

	writeResponse(rw, 200, res, nil, logger)
}

func parseRequest(request *http.Request, bodyStruct any) error {
	b, err := io.ReadAll(request.Body)
	if err != nil {
		return err
	}
	defer request.Body.Close()

	if err := json.Unmarshal(b, bodyStruct); err != nil {
		return err
	}

	return nil
}

func writeResponse(rw http.ResponseWriter, statusCode int, body any, resErr error, logger *logrus.Entry) {
	if statusCode != 200 {
		rw.WriteHeader(statusCode)
	}

	if resErr != nil {
		if _, err := rw.Write([]byte(resErr.Error())); err != nil {
			logger.WithField("error", err).Error("failed to write response")
		}

		return
	}

	b, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		logger.WithField("error", err).Error("failed to marshal json response")
		if _, err := rw.Write([]byte(ErrInternalServerError.Error())); err != nil {
			logger.WithField("error", err).Error("failed to write response")
		}
	}

	if _, err := rw.Write(b); err != nil {
		logger.WithField("error", err).Error("failed to write response")
	}
}
