import "jquery"
import "bootstrap-honoka"
import "bootstrap-honoka/dist/css/bootstrap.min.css"
import "animate.css"
import "./style.scss"

import axios, {Axios} from "axios"

//For Android PWA
/*if ("serviceWorker" in navigator) {
    navigator.serviceWorker.register("sw.js").then((req) => {
        console.log("Service worker registerd.", req)
    });
}*/

interface UrlParams {
    [key: string]: string;
}

class UrlBuilder {
    private readonly params: UrlParams = {};

    constructor(
        tabName = "searchTab",
        from = "",
        to = "",
        locale = "ja",
        bsid = "1"
    ) {
        this.params["tabName"] = tabName
        this.params["from"] = from
        this.params["to"] = to
        this.params["locale"] = locale
        this.params["bsid"] = bsid
    }

    from(from: string): UrlBuilder {
        this.params["from"] = from
        return this
    }

    to(to: string): UrlBuilder {
        this.params["to"] = to
        return this
    }

    build(): string {
        let url = "https://hakobus.bus-navigation.jp/wgsys/wgs/bus.htm?";

        Object.keys(this.params).forEach(key => {
            url += key + "=" + this.params[key] + "&";
        });

        return url.slice(0, url.length - 1);
    }


}

class HTTPClient {
    private url: string = "";
    
    setUrl(url: string) {
        this.url = url;
    }
    
    async fetch(): Promise<string> {
        //const res = await axios.get(this.url);
        //return res.data

        const data = await fetch(this.url, {
            method: "GET",
            mode: "no-cors",
        }).then(function (response) {
            return response.text();
        }).catch(function (reason) {
            return reason.text()
        });

        console.log(data)
        return data;
    }
}

//usage
//var client = new HTTPClient();
//client.setUrl("https://jsonplaceholder.typicode.com/todos/1");
//client.fetch();

//https://developer.mozilla.org/ja/docs/Web/API/Fetch_API/Using_Fetch
//https://developer.mozilla.org/ja/docs/Web/API/Response/text

//https://maku.blog/p/x3ocp9a/

class Parser {
    
}

let urlBuilder = new UrlBuilder()
let url = urlBuilder.from("函館駅前").to("五稜郭").build()
console.log(url)

let httpClient = new HTTPClient()
httpClient.setUrl(url)
httpClient.fetch().then(r => console.log(r))