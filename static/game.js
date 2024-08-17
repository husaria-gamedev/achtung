const colorArray = ['#FF6633', '#FFB399', '#FF33FF', '#FFFF99', '#00B3E6', 
		  '#E6B333', '#3366E6', '#999966', '#99FF99', '#B34D4D',
		  '#80B300', '#809900', '#E6B3B3', '#6680B3', '#66991A', 
		  '#FF99E6', '#CCFF1A', '#FF1A66', '#E6331A', '#33FFCC',
		  '#66994D', '#B366CC', '#4D8000', '#B33300', '#CC80CC', 
		  '#66664D', '#991AFF', '#E666FF', '#4DB3FF', '#1AB399',
		  '#E666B3', '#33991A', '#CC9999', '#B3B31A', '#00E680', 
		  '#4D8066', '#809980', '#E6FF80', '#1AFF33', '#999933',
		  '#FF3380', '#CCCC00', '#66E64D', '#4D80CC', '#9900B3', 
		  '#E64D66', '#4DB380', '#FF4D4D', '#99E6E6', '#6666FF'];

const FPS = 50;
const WINDOW_WIDTH = window.innerWidth;
const WINDOW_HEIGHT = window.innerHeight;
const BORDER = 5;
const CANVAS_EDGE = Math.min(WINDOW_HEIGHT, WINDOW_WIDTH) - BORDER * 2;
const DIRECTION_DICT = {
    left: "l",
    right: "r",
    none: "n",
}

let direction = "n";
let paused = false;

const MAP_EDGE = 1000;
const SCALE = CANVAS_EDGE/MAP_EDGE;
const BREADTH = 4*SCALE;
const HEAD_BREADTH = 3*SCALE;

const canvas = document.getElementById("game");
const ctx = canvas.getContext("2d");

let lastData = [];
let newData = [];

addEventListener("keydown", (event) => {
    if (event.key == "ArrowLeft")
        direction = "l"
    if (event.key == "ArrowRight")
        direction = "r"
});

addEventListener("keyup", (event) => {
    if (event.key == "ArrowLeft" || event.key == "ArrowRight")
        direction = "n"
});

addEventListener("touchstart", (event) => {
    const touch = event.targetTouches[0] || event.changedTouches[0];

    if (touch.clientX < WINDOW_WIDTH / 2) {
        direction = "l";
    } else {
        direction = "r";
    }
});

addEventListener("touchend", (event) => {
    direction = "n";
});

function clearHeads() {
    lastData.forEach(p => {drawDot(p.x, p.y, colorArray[p.i], BREADTH)})
}

function drawHeads() {
    newData.forEach(p => {drawDot(p.x, p.y, "yellow", HEAD_BREADTH)})
}

function drawDot(x, y, color, breadth) {
    ctx.fillStyle = color;
    ctx.beginPath();
    ctx.arc(x*SCALE,y*SCALE,breadth,0,Math.PI*2,true);
    ctx.fill();
}

function processData() {
    clearHeads();
    drawHeads();

    lastData = newData;
}

function init() {
    canvas.height = CANVAS_EDGE;
    canvas.width = CANVAS_EDGE;
    ctx.fillStyle = "black";
    ctx.fillRect(0, 0, CANVAS_EDGE, CANVAS_EDGE);
}

const socket = new WebSocket(`ws://${window.location.host}/websocket`);

function sendControls() {
    socket.send(JSON.stringify({d: direction}));
}

socket.onopen = () => {
    init();
    setInterval(sendControls, 1000/FPS);
};

socket.onmessage = (event) => {
    data = JSON.parse(event.data);
    if (data) {
        newData = data;
        processData();
    }
};
