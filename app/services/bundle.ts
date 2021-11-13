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
    private url: string;
    private axios = require("axios")
    
    axios.get(url).then(function (response) {
        //正常
    }).catch(function (error) {
        //異常
    }).then(function () {
        //finally
    })
    
}
