package infrastructure

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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

func (p *URLParameter) toURL() string {
	const BASE_URI = "https://hakobus.bus-navigation.jp/wgsys/wgs/bus.htm?"

	return BASE_URI + "tabName=" + p.tabName + "&from=" + p.from + "&to=" + p.to + "&locale=" + p.locale + "&bsid=" + p.bsid
}

func fetch(url string) []byte {
	println("Access: " + url)

	response, _ := http.Get(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//connection error
		}
	}(response.Body)

	body, _ := ioutil.ReadAll(response.Body)
	return body
}

func parse(buffer []byte) BusInformation {
	var busInformation BusInformation

	reader := bytes.NewReader(buffer)
	document, _ := goquery.NewDocumentFromReader(reader)

	busInformation.RefTime = parseTime(siblingAt(1, document.Find("div.container").Find("div.label_bar").Find("div.clearfix").Find("li")).Text())

	busInformation.IsBusExist = document.Find("div#errInfo").text() == ""
	busInformation.Results = []result{}
	//60分以内にバスがない(情報が掲載されてない)とき、スクレイピングを終了する
	if !busInformation.IsBusExist {
		return busInformation
	}

	document.Find("div#buslist").Find("div.clearfix").Find("div.route_box").Each(func(i int, selection *goquery.Selection) {
		bus := selection.Find("table").Find("tbody").Find("tr")
		result := result{}

		result.Name = siblingAt(0, bus).Find("td").Find("span").Text()
		result.Via = siblingAt(1, bus).Text()
		result.Direction = siblingAt(2, bus).Text()
		result.Estimate, err = strconv.Atoi(timeParse(siblingAt(3, bus).Find("td").First().Text()))
		result.From = siblingAt(4, bus).Find("font").Text()
		departure := siblingAt(5, bus).Find("td")
		result.Departure.Schedule, err = timeParse(siblingAt(0, departure).Text())
		result.Departure.Prediction, err = timeParse(siblingAt(1, departure).Text())
		result.To = siblingAt(6, bus).Find("font").Text()
		arrive := siblingAt(7, bus).Find("td")
		result.Arrive.Schedule = timeParse(siblingAt(0, arrive).Text())
		result.Arrive.Prediction, err = timeParse(siblingAt(1, arrive).Text())
		take := siblingAt(8, bus).Find("td").First()
		if take.Text() == "まもなく発車します" {
			result.Take = 0
		} else {
			result.Take, err = strconv.Atoi(take.Find("font").Text())
			if err != nil {
				panic(err)
			}
		}

		busInformation.Results = append(busInformation.Results, result)
	})

	return busInformation
}

func parseTime(s string) *time.Time {
	const tr = regexp.MustCompile(`[^0-9:-]`)
	const TIME_PATTERN = "15:04"

	s1 := tr.ReplaceAllString(s, "")
	if strings.Contains(s1, "-") {
		return nil
	}

	t, e := time.Parse(TIME_PATTERN, s1)
	if e != nil {
		return nil
	}

	return t
}

// 同じ親のi番目の兄弟を取得します
// 配列ライクに扱えるはず？
// i : 何番目の要素か
func siblingAt(i int, s *goquery.Selection) *goquery.Selection {
	if i == 0 {
		return s.First()
	} else if i > 0 {
		return siblingAt(i-1, s.Next())
	} else {
		//動作未検証
		return siblingAt(i+1, s.Prev())
	}
}
