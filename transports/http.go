package transports

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/GorillaPool/go-junglebus/models"
	"github.com/mdp/qrterminal/v3"
)

// TransportHTTP is the struct for HTTP
type TransportHTTP struct {
	debug      bool
	httpClient *http.Client
	server     string
	token      string
	useSSL     bool
	version    string
}

// SetDebug turn the debugging on or off
func (h *TransportHTTP) SetDebug(debug bool) {
	h.debug = debug
}

// IsDebug return the debugging status
func (h *TransportHTTP) IsDebug() bool {
	return h.debug
}

// UseSSL turn the SSL on or off
func (h *TransportHTTP) UseSSL(useSSL bool) {
	h.useSSL = useSSL
}

// IsSSL return the SSL status
func (h *TransportHTTP) IsSSL() bool {
	return h.useSSL
}

// SetToken sets the token to use for all requests manually
func (h *TransportHTTP) SetToken(token string) {
	h.token = token
}

// GetToken gets the token to use for all requests
func (h *TransportHTTP) GetToken() string {
	return h.token
}

// GetSubscriptionToken gets a token based on the subscription ID
func (h *TransportHTTP) GetSubscriptionToken(ctx context.Context, subscriptionID string) (string, error) {

	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldSubscriptionID: subscriptionID,
	})
	if err != nil {
		return "", err
	}

	var response LoginResponse
	if err = h.doHTTPJsonRequest(
		ctx, http.MethodPost, `/user/subscription-token`, jsonStr, &response,
	); err != nil {
		return "", err
	}

	return response.Token, nil
}

// RefreshToken gets a new  token to use for all requests
func (h *TransportHTTP) RefreshToken(ctx context.Context) (string, error) {
	var response LoginResponse
	if err := h.doHTTPJsonRequest(
		ctx, http.MethodGet, `/user/refresh-token`, nil, &response,
	); err != nil {
		return "", err
	}

	return response.Token, nil
}

// SetVersion sets the version to use for all calls
func (h *TransportHTTP) SetVersion(version string) {
	h.version = version
}

// GetServerURL get the server URL for this transport
func (h *TransportHTTP) GetServerURL() string {
	return h.server
}

func (h *TransportHTTP) Login(ctx context.Context, username string, password string) error {

	jsonStr, err := json.Marshal(map[string]interface{}{
		FieldUsername: username,
		FieldPassword: password,
	})
	if err != nil {
		return err
	}

	var loginResponse map[string]interface{}
	if err = h.doHTTPJsonRequest(
		ctx, http.MethodPost, `/user/login`, jsonStr, &loginResponse,
	); err != nil {
		return err
	}
	if h.debug {
		log.Printf("Login: %v\n", loginResponse)
	}

	if token, ok := loginResponse["token"]; ok {
		h.SetToken(token.(string))
		return nil
	}

	return ErrFailedLogin
}

// GetTransaction will get a transaction by ID
func (h *TransportHTTP) GetTransaction(ctx context.Context, txID string) (transaction *models.Transaction, err error) {

	if err = h.doHTTPJsonRequest(
		ctx, http.MethodGet, "/transaction/get/"+txID, nil, &transaction,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("Transaction: %v\n", transaction)
	}

	return transaction, nil
}

// GetTransaction will get a transaction by ID
func (h *TransportHTTP) GetRawTransaction(ctx context.Context, txID string) (rawtx []byte, err error) {

	if rawtx, err = h.doHTTPGetRequest(ctx, "/transaction/get/"+txID+"/bin"); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("Transaction: %x\n", rawtx)
	}

	return rawtx, nil
}

func (h *TransportHTTP) GetFromBlock(ctx context.Context, subscriptionID string, height uint32, idx uint64) (transactions []*models.Transaction, err error) {
	url := fmt.Sprintf(
		"/transaction/from-block?subscription_id=%s&block_height=%d&last_idx=%d",
		subscriptionID, height, idx,
	)
	if err = h.doHTTPJsonRequest(
		ctx, http.MethodGet, url, nil, &transactions,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("Transactions: %v\n", transactions)
	}

	return transactions, nil
}

func (h *TransportHTTP) GetLiteFromBlock(ctx context.Context, subscriptionID string, height uint32, idx uint64) (transactions []*models.TransactionResponse, err error) {
	url := fmt.Sprintf(
		"/transaction/from-block-lite?subscription_id=%s&block_height=%d&last_idx=%d",
		subscriptionID,
		height,
		idx,
	)
	if err = h.doHTTPJsonRequest(
		ctx, http.MethodGet, url, nil, &transactions,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("Transactions: %v\n", transactions)
	}

	return transactions, nil
}

// GetAddressTransactions will get the metadata of all transaction related to the given address
func (h *TransportHTTP) GetAddressTransactions(ctx context.Context, address string, fromHeight uint32) (addr []*models.AddressTx, err error) {
	if err = h.doHTTPJsonRequest(
		ctx, http.MethodGet, fmt.Sprintf("/address/get/%s/%d", address, fromHeight), nil, &addr,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("Address transactions: %v\n", addr)
	}

	return addr, nil
}

// GetAddressTransactionDetails will get all transactions related to the given address
func (h *TransportHTTP) GetAddressTransactionDetails(ctx context.Context, address string, fromHeight uint32) (transactions []*models.Transaction, err error) {

	if err = h.doHTTPJsonRequest(
		ctx, http.MethodGet, fmt.Sprintf("/address/transactions/%s/%d", address, fromHeight), nil, &transactions,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("transactions: %d\n", len(transactions))
	}

	return transactions, nil
}

// GetBlockHeader will get the given block header details
// Can pass either the block hash or the block height (as a string)
func (h *TransportHTTP) GetBlockHeader(ctx context.Context, block string) (blockHeader *models.BlockHeader, err error) {
	if err = h.doHTTPJsonRequest(
		ctx, http.MethodGet, "/block_header/get/"+block, nil, &blockHeader,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("block header: %v\n", blockHeader)
	}

	return blockHeader, nil
}

// GetBlockHeaders will get all block headers from the given block, limited by limit
// Can pass either the block hash or the block height (as a string)
func (h *TransportHTTP) GetBlockHeaders(ctx context.Context, fromBlock string, limit uint) (blockHeaders []*models.BlockHeader, err error) {
	if err = h.doHTTPJsonRequest(
		ctx, http.MethodGet, fmt.Sprintf("/block_header/list/%s?limit=%d", fromBlock, limit), nil, &blockHeaders,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("block headers: %v\n", blockHeaders)
	}

	return blockHeaders, nil
}

func (h *TransportHTTP) GetChaintip(ctx context.Context) (blockHeader *models.BlockHeader, err error) {
	if err = h.doHTTPJsonRequest(
		ctx, http.MethodGet, "/block_header/tip", nil, &blockHeader,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("Chaintip: %v\n", blockHeader)
	}

	return blockHeader, nil
}

func (h *TransportHTTP) GetUser(ctx context.Context) (*models.User, error) {
	user := &models.User{}
	if err := h.doHTTPJsonRequest(
		ctx, http.MethodGet, "/user/get", nil, user,
	); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("user: %v\n", user)
	}

	return user, nil
}

func (h *TransportHTTP) GetTxo(ctx context.Context, txID string, vout uint32) (txo []byte, err error) {
	if txo, err = h.doHTTPGetRequest(ctx, fmt.Sprintf("/txo/get/%s_%d", txID, vout)); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("Txo: %x\n", txo)
	}

	return txo, nil
}

func (h *TransportHTTP) GetSpend(ctx context.Context, txID string, vout uint32) (spend []byte, err error) {
	if spend, err = h.doHTTPGetRequest(ctx, fmt.Sprintf("/txo/spend/%s_%d", txID, vout)); err != nil {
		return nil, err
	}
	if h.debug {
		log.Printf("Spend: %x\n", spend)
	}

	return spend, nil
}

// doHTTPJsonRequest will create and submit the HTTP request
func (h *TransportHTTP) doHTTPJsonRequest(ctx context.Context, method string, path string, rawJSON []byte, responseJSON interface{}) error {

	protocol := "https"
	if !h.useSSL {
		protocol = "http"
	}
	serverRequest := fmt.Sprintf("%s://%s/%s%s", protocol, h.server, h.version, path)
	req, err := http.NewRequestWithContext(ctx, method, serverRequest, bytes.NewBuffer(rawJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", h.token)

	var resp *http.Response
	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if resp, err = h.httpClient.Do(req); err != nil {
		return err
	}
	if resp.StatusCode == http.StatusPaymentRequired {
		address := resp.Header.Get("jb-fund-address")
		buf := bytes.NewBuffer([]byte{})
		qrterminal.Generate(address, qrterminal.L, buf)
		log.Printf("WARNING. Free tier exhausted. Please fund your account to continue recieving data: %s\n %s", address, buf.String())
	} else if resp.StatusCode >= http.StatusBadRequest {
		return errors.New("server error: " + strconv.Itoa(resp.StatusCode) + " - " + resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(&responseJSON)
}

// doHTTPRequest will create and submit the HTTP request
func (h *TransportHTTP) doHTTPGetRequest(ctx context.Context, path string) (body []byte, err error) {

	protocol := "https"
	if !h.useSSL {
		protocol = "http"
	}
	serverRequest := fmt.Sprintf("%s://%s/%s%s", protocol, h.server, h.version, path)
	req, err := http.NewRequestWithContext(ctx, "GET", serverRequest, bytes.NewBuffer([]byte{}))
	if err != nil {
		return
	}
	req.Header.Set("token", h.token)

	var resp *http.Response
	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if resp, err = h.httpClient.Do(req); err != nil {
		return
	}
	if resp.StatusCode == http.StatusPaymentRequired {
		// log.Println("Payment Required: ", resp.Header.Get("jb-fund-address"))
	} else if resp.StatusCode >= http.StatusBadRequest {
		err = errors.New("server error: " + strconv.Itoa(resp.StatusCode) + " - " + resp.Status)
		return
	}

	return io.ReadAll(resp.Body)
}
