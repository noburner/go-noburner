package noburner

import (

	// standard

	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (

	// constants

	apiURL = "https://api.noburner.com"

	queryToken = "token"

	verifyTimeout = time.Second * 5
)

var (

	// Errors
	ErrFailed  = errors.New("failed")
	ErrTimeout = errors.New("timeout")
)

// Config
type Config struct {

	// publics

	Secret string

	Client *http.Client
}

// Result
type Result struct {

	// publics
}

// NoBurner
type NoBurner struct {

	// privates

	secret string

	client *http.Client
}

// New
func New(config Config) *NoBurner {

	nb := new(NoBurner)

	nb.secret = config.Secret

	if config.Client == nil {

		nb.client = http.DefaultClient

	} else {

		nb.client = config.Client
	}

	return nb
}

// Verify
func (nb *NoBurner) Verify(ctx context.Context, email string) (*Result, error) {

	if ctx == nil {

		ctx = context.Background()
	}

	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(ctx, verifyTimeout)

	defer cancel()

	url, err := url.Parse(apiURL)

	if err != nil {

		return nil, ErrFailed
	}

	url = url.JoinPath(email)

	query := url.Query()

	query.Add(queryToken, nb.secret)

	url.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)

	if err != nil {

		return nil, ErrFailed
	}

	response, err := nb.client.Do(request)

	if ctx.Err() != nil {

		return nil, ErrTimeout
	}

	if err != nil {

		return nil, ErrFailed
	}

	defer response.Body.Close()

	var result *Result

	switch response.StatusCode {

	case http.StatusOK:

		data, err := io.ReadAll(response.Body)

		if err != nil {

			return nil, ErrFailed
		}

		err = json.Unmarshal(data, &result)

		if err != nil {

			return nil, ErrFailed
		}

	default:

		return nil, ErrFailed
	}

	return result, nil
}
