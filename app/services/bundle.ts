import "jquery"
import "bootstrap-honoka"
import "bootstrap-honoka/dist/css/bootstrap.min.css"
import "animate.css"
import "./style.scss"

const BASE__URL = "http://127.0.0.1:8080/api?"

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
        from = "",
        to = "",
    ) {
        this.params["from"] = from;
        this.params["to"] = to;
    }

    from(from: string): UrlBuilder {
        this.params["from"] = from;
        return this;
    }

    to(to: string): UrlBuilder {
        this.params["to"] = to;
        return this;
    }

    build(): string {
        let url = BASE__URL;

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
        const data = await fetch(this.url, {
            method: "GET",
            mode: "cors",
            headers: {
                "Content-Type": "application/json",
            },
        }).then(function (response) {
            return response.text();
        }).catch(function (reason) {
            return reason.text();
        });

        console.log(data);
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
type BusInformation = {
    refTime: string,
    isBusExist: string,
    results: Result[],
}

type Result = {
    name: string,
    via: string,
    direction: string,
    from: string,
    to: string,
    departure: BusTime,
    arrive: BusTime,
    take: number,
    estimate: number,
}

type BusTime = {
    schedule: string,
    prediction: string,
}
class Parser {
    //テスト用
    text: string = "{\"reftime\":\"13:55\",\"isbusexist\":\"true\",\"results\":[{\"name\":\"N21\",\"via\":\"西浦線\",\"direction\":\"江梨\",\"from\":\"沼津駅\",\"to\":\"長井崎小中一貫学校\",\"departure\":{\"schedule\":\"13:55\",\"prediction\":\"13:55\"},\"arrive\":{\"schedule\":\"14:45\",\"prediction\":\"14:45\"},\"take\":50,\"estimate\":30}]}";

    parse<T>(json: string, type: T): T {
        return JSON.parse(json) as T;
    }
}

let urlBuilder = new UrlBuilder()
let url = urlBuilder.from("函館駅前").to("五稜郭").build()
console.log(url)

let httpClient = new HTTPClient()
httpClient.setUrl(url)
httpClient.fetch().then(r => console.log(r))
