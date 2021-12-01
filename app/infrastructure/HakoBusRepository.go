package infrastructure

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/PuerkitoBio/goquery"
)

type URLParameter struct {
	tabName		string
	from		string
	to			string
	locale		string
	bsid		string
}

type BusInformation struct {
	refTime		time.Time
	isBusExist	bool
	results		[]result
}

type result struct {
	name		string //バス系統名
	via			string //バス経由地
	direction	string //バスの目的地
	from		string //乗車バス停
	to			string //降車バス停
	departure	busTime //バス発車
	arrive		busTime //バス到着
	take		int //予想乗車時間
	estimate	int //あと何分後にバスが来るか
}

type busTime struct {
	schedule	time.Time
	prediction	time.Time
}

func (b *busTime) delayed() int {
	p := b.prediction.Hour() * 60 + b.prediction.Minute()
	s := b.schedule.Hour() * 60 + b.schedule.Minute()

	return p - s
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

func (p *URLParameter) fetch() BusInformation {
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
	return parse(body)
}

func parse(buffer []byte) BusInformation {
	const TIME_REGEX = "[^0-9:-]"
	const TIME_PATTERN = "HH:mm"

	var busInformation BusInformation

	reader := bytes.NewReader(buffer)
	document, _ := goquery.NewDocumentFromReader(reader)

	refTime, err := time.Parse(
		TIME_PATTERN,
		document.Find("div.container").Find("div.label_bar").Find("div.clearfix").Find("ul").Find("li").Next().Text(),
	)
	if err != nil {

	}

	//60分以内にバスがない(情報が掲載されてない)とき、スクレイピングを終了する
	if document.Find("div#errInfo") != nil {
		busInformation = BusInformation{
			refTime: refTime,
			isBusExist: false,
			results: []result{},
		}
	}

	busList := document.Find("div#buslist").Find("div.clearfix").Find("div.route_box")
	busList.Each(func(i int, selection *goquery.Selection) {
		bus := selection.Find("table").Find("tbody").Find("tr")
		result := result{}

		result.name = bus.Find("td").Find("span").Text();bus.Next()
		result.via = bus.Find("td").Text();bus.Next()
		result.direction = bus.Find("td").Text();bus.Next()
		result.estimate, _ = strconv.Atoi(strings.ReplaceAll(bus.Find("td").Text(), TIME_REGEX, ""));bus.Next()
		result.from = bus.Find("td").Find("span").Next().Text();bus.Next()
		//result.departure
		departure := bus.Find("td");bus.Next()
		result.departure.schedule, _ = time.Parse(TIME_PATTERN, strings.ReplaceAll(departure.Text(), TIME_REGEX, ""));departure.Next()
		result.departure.prediction, _ = time.Parse(TIME_PATTERN, strings.ReplaceAll(departure.Text(), TIME_REGEX, ""))
		result.to = bus.Find("td").Find("span").Next().Text();bus.Next()
		arrive := bus.Find("td");bus.Next()
		result.arrive.schedule, _ = time.Parse(TIME_PATTERN, strings.ReplaceAll(departure.Text(), TIME_REGEX, ""));arrive.Next()
		result.arrive.prediction, _ = time.Parse(TIME_PATTERN, strings.ReplaceAll(departure.Text(), TIME_REGEX, ""))
		take := bus.Find("td").Text()
		if bus.Find("td").Text() == "まもなく発車します" {
			result.take = 0
		} else {
			result.take, _ = strconv.Atoi(take)
		}

		busInformation.results = append(busInformation.results, result)
	})

	return busInformation
}