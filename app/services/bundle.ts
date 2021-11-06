import "jquery"
import "bootstrap-honoka"
import "bootstrap-honoka/dist/css/bootstrap.min.css"
import "animate.css"
import "./style.scss"

(function () {
    const SPEED = 0.5;

    let unit = 100,
        canvas: HTMLCanvasElement,
        context: CanvasRenderingContext2D,
        height: number,
        width: number,
        xAxis: number,
        yAxis: number;

    function init() {
        canvas = <HTMLCanvasElement> document.getElementById("canvas-wave");

        canvas.width = document.documentElement.clientWidth; 
        canvas.height = 300;

        context = canvas.getContext("2d");

        height = canvas.height;
        width = canvas.width;

        xAxis = Math.floor(height / 2);
        yAxis = 0;

        draw();
    }

    function draw() {
        context.clearRect(0, 0, width, height);

        //drawWave('#a0c0d4', 0.5, 3, SPEED,0);
        //drawWave('#7ba4d9', 0.6, 2, SPEED * 1.5, 87);
        //drawWave('#c0faff', 0.4, 1.6, SPEED * 2.0, 174);

        drawWave('#c0faff', 1, 3, SPEED,0);
        drawWave('#a0c0d4', 0.8, 2, SPEED * 1.5, 87);
        drawWave('#99c2ed', 0.8, 1.6, SPEED * 2.0, 174);

        draw.seconds = draw.seconds + 0.014;
        draw.t = draw.seconds * Math.PI;
        setTimeout(draw, 35);
    }
    draw.seconds = 0;
    draw.t = 0;

    function drawWave(color: string, alpha: number, zoom: number, speed: number, delay: number) {
        context.fillStyle = color;
        context.globalAlpha = alpha;

        context.beginPath();
        drawSine(draw.t / 0.5, zoom, speed, delay);
        context.lineTo(width + 10, height);
        context.lineTo(0, height); 
        context.closePath() 
        context.fill(); 
    }

    function drawSine(t: number, zoom: number, speed: number, delay: number) {

        let x = t; 
        let y = Math.sin(x) / zoom;
        context.moveTo(yAxis, unit * y + xAxis);
        
        for (let i = yAxis; i <= width + 10; i += 10) {
            x = (t + (-yAxis + i) / unit / zoom) * speed;
            y = Math.sin(x - delay);
            context.lineTo(i, unit * y + xAxis);
        }
    }

    init();

})();