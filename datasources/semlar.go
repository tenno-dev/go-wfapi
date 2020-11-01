package datasources

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Anomalydata Kuvadata result
var Anomalydata []byte

// LoadAnomalydata  Load Kuvadata from https://semlar.com/anomaly.json
func LoadAnomalydata() (ret []byte) {
	client := &http.Client{}

	url := "https://semlar.com/anomaly.json"
	fmt.Println("url:", url)
	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	_, _ = io.Copy(ioutil.Discard, res.Body)
	Anomalydata = body[:]
	return
}
