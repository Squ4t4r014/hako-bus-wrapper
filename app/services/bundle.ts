import "jquery"
import "bootstrap-honoka"
import "bootstrap-honoka/dist/css/bootstrap.min.css"
import "animate.css"
import "./style.scss"

//For Android PWA
if ("serviceWorker" in navigator) {
    navigator.serviceWorker.register("sw.js").then((req) => {
        console.log("Service worker registerd.", req)
    });
}

class HTTPClient {
    private url: string = "";
    //private axios = require("axios")
    
    private tabName = "searchTab";
    private from = "";
    private to = "";
    private locale = "ja";
    private bsid = "1";
    
    reset() {
        this.tabName = "searchTab";
        this.from = "";
        this.to = "";
        this.locale = "ja";
        this.bsid = "1";
    }
    
    setUrl(url: string) {
        this.url = url;
    }
    
    from(busStop: String) {
        this.from = busStop
    }
    
    to(busStop: String) {
        this.to = busStop
    }
    
    fetch() {
        fetch(this.url)
        .then(response => response.text())
        .then(text => {
            //callback
            console.log(text)
        });
    }
}

//usage
//var client = new HTTPClient();
//client.setUrl("https://jsonplaceholder.typicode.com/todos/1");
//client.fetch();

//https://developer.mozilla.org/ja/docs/Web/API/Fetch_API/Using_Fetch
//https://developer.mozilla.org/ja/docs/Web/API/Response/text

class Parser {
    
}
