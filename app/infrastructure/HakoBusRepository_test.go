package infrastructure

import (
  "testing"
)

func TestNewRideSuccess(t *testing.T) {
    ans := &URLParamater{
    	tabName: "searchTab",
	    from:    "函館駅前",
		to:      "西高校前",
		locale:  "ja",
		bsid:    "1",
    }
    ride := newRide("函館駅前", "西高校前")
    
    if ans != ride {
        t.Fatal("NewRide")
    }
}

func TestToURLSuccess(t *testing.T) {
    const url := "https://hakobus.bus-navigation.jp/wgsys/wgs/bus.htm?tabName=searchTab&from=函館駅前&to=西高校前&locale=ja&bsid=1"
    ride := newRide("函館駅前", "西高校前").toURL()
    if url != ride {
        t.Fatal("ToURL")
    }
}
