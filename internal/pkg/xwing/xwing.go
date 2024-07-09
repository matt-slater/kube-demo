package xwing

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

type XWing struct {
	HTTPClient http.Client
	Target     string
	Commander  string
	Logger     zerolog.Logger
}

func New(target, commander string, logger zerolog.Logger) *XWing {
	return &XWing{
		HTTPClient: *http.DefaultClient,
		Target:     target,
		Commander:  commander,
		Logger:     logger,
	}
}

func (xw *XWing) Launch() error {
	xw.Logger.Info().Msg(fmt.Sprintf("x-wing launched by %s, attacking %s", xw.Commander, xw.Target))

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/", xw.Target), nil)
	if err != nil {
		return fmt.Errorf("error formulating flight plan: %w", err)
	}

	req.Header.Set("commander", xw.Commander)

	resp, err := xw.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error performing flight plan: %w", err)
	}

	defer resp.Body.Close()

	if !(resp.StatusCode == http.StatusOK) {
		return errors.New("mission failed")
	}

	return nil
}
