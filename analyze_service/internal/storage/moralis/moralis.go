package moralis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type (
	MoralisRepository struct {
		apiKey     string
		baseURL    string
		httpClient *http.Client
	}
	Transaction struct {
		Hash             string           `json:"hash"`
		Nonce            string           `json:"nonce,omitempty"`
		TransactionIndex string           `json:"transaction_index,omitempty"`
		FromAddress      string           `json:"from_address"`
		FromAddressLabel string           `json:"from_address_label,omitempty"`
		ToAddress        string           `json:"to_address"`
		ToAddressLabel   string           `json:"to_address_label,omitempty"`
		Value            string           `json:"value"`
		Gas              string           `json:"gas,omitempty"`
		GasPrice         string           `json:"gas_price,omitempty"`
		ReceiptStatus    string           `json:"receipt_status,omitempty"`
		BlockTimestamp   time.Time        `json:"block_timestamp"`
		BlockNumber      string           `json:"block_number"`
		BlockHash        string           `json:"block_hash,omitempty"`
		TransactionFee   string           `json:"transaction_fee,omitempty"`
		MethodLabel      string           `json:"method_label,omitempty"`
		NativeTransfers  []NativeTransfer `json:"native_transfers,omitempty"`
		Summary          string           `json:"summary,omitempty"`
		PossibleSpam     bool             `json:"possible_spam,omitempty"`
		Category         string           `json:"category,omitempty"`
	}

	NativeTransfer struct {
		FromAddress      string `json:"from_address"`
		FromAddressLabel string `json:"from_address_label,omitempty"`
		ToAddress        string `json:"to_address"`
		ToAddressLabel   string `json:"to_address_label,omitempty"`
		Value            string `json:"value"`
		ValueFormatted   string `json:"value_formatted,omitempty"`
		Direction        string `json:"direction,omitempty"`
		TokenSymbol      string `json:"token_symbol,omitempty"`
		TokenLogo        string `json:"token_logo,omitempty"`
	}

	MoralisResponse struct {
		Cursor   string        `json:"cursor"`
		PageSize int           `json:"page_size"`
		Limit    string        `json:"limit"`
		Result   []Transaction `json:"result"`
	}
)

func NewMoralisRepository(apiKey string) *MoralisRepository {
	return &MoralisRepository{
		apiKey:  apiKey,
		baseURL: "https://deep-index.moralis.io/api/v2",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (m *MoralisRepository) GetTransactions(address string) (*MoralisResponse, error) {
	url := fmt.Sprintf("%s/%s?chain=eth", m.baseURL, address)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-API-Key", m.apiKey)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("moralis api error: %s, body: %s", resp.Status, string(body))
	}

	var moralisResponse MoralisResponse
	if err := json.NewDecoder(resp.Body).Decode(&moralisResponse); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &moralisResponse, nil
}
