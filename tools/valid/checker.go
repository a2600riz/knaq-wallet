package valid

import (
	"encoding/json"
	"fmt"
	"knaq-wallet/tools/network"
	"net/http"
)

func TradingNameChecker(s string) error {
	for _, r := range s {
		if r > 255 {
			return fmt.Errorf("trading name must contain only latin characters and special characters")
		}
	}
	return nil
}

func PostCodeChecker(postcode string) (string, error) {
	type postCodeCheckerData struct {
		Status int `json:"status"`
		Result []struct {
			Result *struct {
				Postcode string `json:"postcode"`
			} `json:"result"`
		} `json:"result"`
	}

	b, err := json.Marshal(map[string]interface{}{
		"postcodes": []string{postcode},
	})
	if err != nil {
		return "", err
	}
	req := network.Request{
		URL:         "https://api.postcodes.io/postcodes",
		Method:      http.MethodPost,
		ContentType: network.TypeJSON,
		Body:        b,
	}
	respBody, status, err := req.Send()
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		var data map[string]interface{}
		if len(respBody) != 0 {
			if err = json.Unmarshal(respBody, &data); err != nil {
				return "", err
			}
		}

		return "", fmt.Errorf("status: %d, message: %+v", status, data)
	}

	var d postCodeCheckerData
	if err = json.Unmarshal(respBody, &d); err != nil {
		return "", err
	}

	if d.Result[0].Result == nil {
		return "", fmt.Errorf("invalid postcode")
	}

	return d.Result[0].Result.Postcode, nil
}
