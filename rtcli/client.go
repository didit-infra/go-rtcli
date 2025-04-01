package rtcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"
)

/*
 *
 */
type Client struct {
	baseURL    string
	httpClient *http.Client
	auth       *authConfig
	logger     *slog.Logger
}

/*
 *
 */
type authConfig struct {
	username string
	password string
	token    string
}

/*
 *
 */
type ClientOptions struct {
	APIURL     string
	Username   string
	Password   string
	Token      string
	Timeout    time.Duration
	Debug      bool
	LogEnabled bool
}

func (c *ClientOptions) validate() error {
	if c.APIURL == "" {
		return fmt.Errorf("RT API URL is required")
	}
	return nil
}

/*
 *
 */
func NewClient(opts ClientOptions) (*Client, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}

	// Crear el cliente HTTP con timeout
	httpClient := &http.Client{
		Timeout: opts.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}
	// Crear el cliente base
	clientRT := &Client{
		baseURL:    opts.APIURL,
		httpClient: httpClient,
		auth: &authConfig{
			username: opts.Username,
			password: opts.Password,
			token:    opts.Token,
		},
		logger: makeLogger(opts.Debug, opts.LogEnabled),
	}
	return clientRT, nil
}

func makeLogger(debug bool, logEnabled bool) *slog.Logger {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	var w io.Writer
	w = os.Stdout
	if !logEnabled {
		w = io.Discard
	}
	handler := slog.NewTextHandler(w, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}

func (c *Client) doRequest(method string, endpoint string, body any, params map[string]string) ([]byte, error) {
	var err error
	// Construir URL completa
	baseURL := fmt.Sprintf("%s/%s", c.baseURL, endpoint)

	// Crear URL con parametros
	requestURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	// Add parms if exists
	if len(params) > 0 {
		q := requestURL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		requestURL.RawQuery = q.Encode()
	}

	var req *http.Request
	if body != nil {
		jsonBody, errb := json.Marshal(body)
		if errb != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", errb)
		}

		// Debug logging del body si está habilitado
		c.logger.Debug("Request body", "body", string(jsonBody))

		req, err = http.NewRequest(method, requestURL.String(), bytes.NewBuffer(jsonBody))
	} else {
		req, err = http.NewRequest(method, requestURL.String(), nil)
	}

	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	// Añadir headers
	req.Header.Set("Content-Type", "application/json")

	// Configurar autenticación
	if c.auth.token != "" {
		req.Header.Set("Authorization", "token "+c.auth.token)
	} else {
		req.SetBasicAuth(c.auth.username, c.auth.password)
	}

	// Debug logging
	c.logger.Debug("Request", "method", method, "url", requestURL.String())

	// Realizar la petición
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Verificar el código de estado
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s",
			resp.StatusCode, string(respBody))
	}
	c.logger.Debug("Response", "body", string(respBody))

	return respBody, nil
}
