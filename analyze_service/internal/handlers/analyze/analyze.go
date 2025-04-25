package analyze

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"

	"analyze-service/internal/config"
	er "analyze-service/internal/lib/api/error"
	resp "analyze-service/internal/lib/api/response"
	"analyze-service/internal/lib/logger/sl"
	"analyze-service/internal/storage/moralis"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type (
	Request struct {
		Address string `json:"address"` // Ethereum address to analyze
	}

	Response struct {
		Status          string                   `json:"status"`
		Message         string                   `json:"message,omitempty"`
		Cursor          string                   `json:"cursor,omitempty"`
		PageSize        int                      `json:"page_size,omitempty"`
		Transactions    []moralis.Transaction    `json:"transactions,omitempty"`
		NativeTransfers []moralis.NativeTransfer `json:"native_transfers,omitempty"`
		MLResult        interface{}              `json:"ml_result,omitempty"` // Result from ML service
	}

	MLRequest struct {
		Address      string          `json:"address"`
		Transactions []MLTransaction `json:"transactions"`
	}

	MLTransaction struct {
		Timestamp int64   `json:"timestamp"`
		Value     float64 `json:"value"`
		Method    string  `json:"method"`
		To        string  `json:"to"`
	}

	Repository interface {
		GetTransactions(address string) (*moralis.MoralisResponse, error)
	}
)

const (
	ethAddressRegex = "^0x[0-9a-fA-F]{40}$"
)

// Analyze handles Ethereum address analysis.
//
//	@Tags			analyze
//	@Summary		Analyze Ethereum address
//	@Description	Get transactions for Ethereum address using Moralis API and analyze with ML service
//	@Accept			json
//	@Produce		json
//	@Param			request	body		Request				true	"Ethereum address to analyze"
//	@Success		200		{object}	Response			"Transactions retrieved and analyzed successfully"
//	@Failure		400		{object}	er.BadRequestErr	"Invalid Ethereum address"
//	@Failure		500		{object}	er.InternalErr		"Internal server error"
//	@Router			/api/v1/analyze [post]
func Analyze(log *slog.Logger, repo Repository, c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.analyze.Analyze"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		reqCtx := r.Context()

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")
				resp.SendBad(reqCtx, er.BadRequestErr{Message: errors.New("empty request")})
				return
			}
			log.Error("failed to decode request body", sl.Err(err))
			resp.SendBad(reqCtx, er.BadRequestErr{Message: errors.New("failed to decode request")})
			return
		}

		if !isValidEthAddress(req.Address) {
			log.Error("invalid Ethereum address", slog.String("address", req.Address))
			resp.SendBad(reqCtx, er.BadRequestErr{Message: errors.New("invalid Ethereum address")})
			return
		}

		moralisResp, err := repo.GetTransactions(req.Address)
		if err != nil {
			log.Error("failed to get transactions", sl.Err(err))
			resp.SendBad(reqCtx, er.InternalErr{Message: errors.New("failed to get transactions")})
			return
		}

		var nativeTransfers []moralis.NativeTransfer
		for _, tx := range moralisResp.Result {
			nativeTransfers = append(nativeTransfers, tx.NativeTransfers...)
		}

		log.Info("successfully retrieved transactions",
			slog.Int("count", len(moralisResp.Result)),
			slog.Int("native_transfers", len(nativeTransfers)))

		// Prepare data for ML service
		mlRequest := prepareMLRequest(req.Address, moralisResp.Result)
		mlResult, err := sendToMLService(c.MLServiceAdress, mlRequest)
		if err != nil {
			log.Error("failed to get ML analysis", sl.Err(err))
			// We still return the transactions even if ML failed
			mlResult = map[string]interface{}{"error": "ML analysis failed"}
		}

		render.JSON(w, r, Response{
			Status:          "success",
			Message:         fmt.Sprintf("Found %d transactions", len(moralisResp.Result)),
			Cursor:          moralisResp.Cursor,
			PageSize:        moralisResp.PageSize,
			Transactions:    moralisResp.Result,
			NativeTransfers: nativeTransfers,
			MLResult:        mlResult,
		})
	}
}

func isValidEthAddress(address string) bool {
	re := regexp.MustCompile(ethAddressRegex)
	return re.MatchString(address)
}

func prepareMLRequest(address string, transactions []moralis.Transaction) MLRequest {
	mlTxs := make([]MLTransaction, 0, len(transactions))

	for _, tx := range transactions {
		mlTx := MLTransaction{
			Timestamp: tx.BlockTimestamp.Unix(),
			Value:     parseValue(tx.Value),
			Method:    tx.MethodLabel,
			To:        tx.ToAddress,
		}
		mlTxs = append(mlTxs, mlTx)
	}

	return MLRequest{
		Address:      address,
		Transactions: mlTxs,
	}
}

func parseValue(valueStr string) float64 {
	// Implement your value parsing logic here
	// This is a placeholder implementation
	var value float64
	_, err := fmt.Sscanf(valueStr, "%f", &value)
	if err != nil {
		return 0
	}
	return value
}

func sendToMLService(mlServiceURL string, request MLRequest) (interface{}, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ML request: %w", err)
	}

	resp, err := http.Post(mlServiceURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send request to ML service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service returned status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode ML response: %w", err)
	}

	return result, nil
}
