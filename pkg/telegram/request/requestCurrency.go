package requestProject

import (
	"encoding/json"
	"io"
	"net/http"
)

type QuoteGet struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func GetNeedQuotes() ([]QuoteGet, error) {
	response, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbols=[%22BTCUSDT%22,%22ETHUSDT%22,%22BNBUSDT%22,%22ADAUSDT%22,%22SOLUSDT%22,%22DOGEUSDT%22,%22DOTUSDT%22,%22MATICUSDT%22,%22TRXUSDT%22,%22ETCUSDT%22,%22LTCUSDT%22,%22XMRUSDT%22,%22BCHUSDT%22,%22XRPUSDT%22]")
	if err != nil {
		return []QuoteGet{}, err
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return []QuoteGet{}, err
		}
		var quoteGetResponse []QuoteGet
		err = json.Unmarshal([]byte(body), &quoteGetResponse)
		if err != nil {
			return []QuoteGet{}, err
		}
		if len(quoteGetResponse) > 0 {
			return quoteGetResponse, nil
		}
	}

	return []QuoteGet{}, nil
}
