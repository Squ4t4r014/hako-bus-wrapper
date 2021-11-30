package infrastructure

import (
	"io"
	"io/ioutil"
	"net/http"
)

type URLParameter struct {
	tabName		string
	from		string
	to			string
	locale		string
	bsid		string
}

func newRide(from string, to string) *URLParameter {
	return &URLParameter{
		tabName: "searchTab",
		from: from,
		to: to,
		locale: "ja",
		bsid: "1",
	}
}

func (p *URLParameter) fetch() string {
	const BASE_URI = "https://hakobus.bus-navigation.jp/wgsys/wgs/bus.htm?"

	url := BASE_URI + "tabName=" + p.tabName + "&from=" + p.from + "&to=" + p.to + "&locale=" + p.locale + "&bsid=" + p.bsid

	response, _ := http.Get(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//connection error
		}
	}(response.Body)

	body, _ := ioutil.ReadAll(response.Body)
	return string(body)
}