package infrastructure

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type URLParameter struct {
	tabName string
	from    string
	to      string
	locale  string
	bsid    string
}

type BusInformation struct {
	RefTime    time.Time `json:"ref_time"`
	IsBusExist bool      `json:"is_bus_exist"`
	Results    []result  `json:"results"`
}

type result struct {
	Name      string  `json:"name"`      //バス系統名
	Via       string  `json:"via"`       //バス経由地
	Direction string  `json:"direction"` //バスの目的地
	From      string  `json:"from"`      //乗車バス停
	To        string  `json:"to"`        //降車バス停
	Departure busTime `json:"departure"` //バス発車
	Arrive    busTime `json:"arrive"`    //バス到着
	Take      int     `json:"take"`      //予想乗車時間
	Estimate  int     `json:"estimate"`  //あと何分後にバスが来るか
}

type busTime struct {
	Schedule   time.Time `json:"schedule"`
	Prediction time.Time `json:"prediction"`
}

func (b *busTime) delayed() int {
	p := b.Prediction.Hour()*60 + b.Prediction.Minute()
	s := b.Schedule.Hour()*60 + b.Schedule.Minute()

	return p - s
}

func newRide(from string, to string) *URLParameter {
	return &URLParameter{
		tabName: "searchTab",
		from:    from,
		to:      to,
		locale:  "ja",
		bsid:    "1",
	}
}

func (p *URLParameter) fetch() BusInformation {
	const BASE_URI = "https://hakobus.bus-navigation.jp/wgsys/wgs/bus.htm?"

	url := BASE_URI + "tabName=" + p.tabName + "&from=" + p.from + "&to=" + p.to + "&locale=" + p.locale + "&bsid=" + p.bsid

	println("Access: " + url)

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
	const TIME_PATTERN = "15:04"

	var busInformation BusInformation

	reader := bytes.NewReader(buffer)
	document, _ := goquery.NewDocumentFromReader(reader)

	refTime, err := time.Parse(
		TIME_PATTERN,
		strings.NewReplacer("\n", "", "\t", "").Replace(document.Find("div.container").Find("div.label_bar").Find("div.clearfix").Find("li").Next().First().Text()),
	)
	if err != nil {
		println(err.Error())
		panic("error")
	}
	busInformation.RefTime = refTime
	
	busInformation.IsBusExist = document.Find("div#errInfo").text() == ""
	busInformation.Results = []result{}
	//60分以内にバスがない(情報が掲載されてない)とき、スクレイピングを終了する
	if !busInformation.IsBusExist {
		return busInformation
	}

	document.Find("div#buslist").Find("div.clearfix").Find("div.route_box").Each(func(i int, selection *goquery.Selection) {
		bus := selection.Find("table").Find("tbody").Find("tr")
		result := result{}

		result.Name = nextSingle(0, bus).Find("td").Find("span").Text()
		result.Via = nextSingle(1, bus).Text()
		result.Direction = nextSingle(2, bus).Text()
		println(removeCtrlStr(strings.ReplaceAll(nextTo(3, bus).Find("td").First().Text(), TIME_REGEX, "")))
		result.Estimate, err = strconv.Atoi(removeCtrlStr(strings.ReplaceAll(nextTo(3, bus).Find("td").First().Text(), TIME_REGEX, "")))
		result.From = nextTo(4, bus).Find("font").Text()
		departure := nextTo(5, bus).Find("td")
		result.Departure.Schedule, _ = time.Parse(TIME_PATTERN, strings.ReplaceAll(nextSingle(0, departure).Text(), TIME_REGEX, ""))
		result.Departure.Prediction, _ = time.Parse(TIME_PATTERN, strings.ReplaceAll(nextSingle(1, departure).Text(), TIME_REGEX, ""))
		result.To = nextTo(6, bus).Find("font").Text()
		arrive := nextTo(7, bus).Find("td")
		result.Arrive.Schedule, _ = time.Parse(TIME_PATTERN, strings.ReplaceAll(nextSingle(0, arrive).Text(), TIME_REGEX, ""))
		result.Arrive.Prediction, _ = time.Parse(TIME_PATTERN, strings.ReplaceAll(nextSingle(1, arrive).Text(), TIME_REGEX, ""))
		take := nextTo(8, bus).Find("td").First().Text()
		if take == "まもなく発車します" {
			result.Take = 0
		} else {
			result.Take, _ = strconv.Atoi(take)
		}

		busInformation.Results = append(busInformation.Results, result)
	})

	return busInformation
}

func removeCtrlStr(s string) string {
	return strings.NewReplacer("\t", "", "\n", "").Replace(s)
}

func nextTo(i int, s *goquery.Selection) *goquery.Selection {
	if i <= 0 {
		return s
	} else {
		return nextTo(i - 1, s.Next())
	}
}

// 同じ親のi番目の兄弟を取得します
// i : 何番目の要素か
func nextSingle(i int, s *goquery.Selection) *goquery.Selection {
	return nextTo(i, s).First()
}
