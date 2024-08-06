const FPS = 30;
const WINDOW_WIDTH = window.innerWidth;
const WINDOW_HEIGHT = window.innerHeight;
const BORDER = 5;
const CANVAS_EDGE = Math.min(WINDOW_HEIGHT, WINDOW_WIDTH) - BORDER * 2;

let direction = "none";
let paused = false;

const WIDTH = 400;
const HEIGHT = 400;

const canvas = document.getElementById("game");
const ctx = canvas.getContext("2d");

let posX = 500;
let posY = 500;
let angle = 0;
let SPEED = 4;
let ANGLE_SPEED = 0.1;

addEventListener("keydown", (event) => {
    // document.body.requestFullscreen();

    if (event.key == "ArrowLeft")
        direction = "left"
    if (event.key == "ArrowRight")
        direction = "right"
});

addEventListener("keyup", (event) => {
    if (event.key == "ArrowLeft" || event.key == "ArrowRight")
        direction = "none"
});

addEventListener("touchstart", (event) => {
    const touch = event.targetTouches[0] || event.changedTouches[0];

    if (touch.clientX < WINDOW_WIDTH / 2) {
        direction = "left";
    } else {
        direction = "right";
    }
});

addEventListener("touchend", (event) => {
    direction = "none";
});


function updatePosition() {
    let deltaX = SPEED * Math.cos(angle);
    let deltaY = SPEED * Math.sin(angle);

    // Update positions
    posX += deltaX;
    posY += deltaY;

    if (direction == "left") {
        angle -= ANGLE_SPEED;
    } else if (direction == "right") {
        angle += ANGLE_SPEED;
    }
}

function tick() {
    ctx.fillStyle = "red";
    ctx.beginPath();
    ctx.arc(posX,posY,5,0,Math.PI*2,true);
    ctx.fill();

    updatePosition();

    ctx.fillStyle = "yellow";
    ctx.beginPath();
    ctx.arc(posX,posY,5,0,Math.PI*2,true);
    ctx.fill();

    socket.send(JSON.stringify({direction: direction}));
}

function init() {
    canvas.height = CANVAS_EDGE;
    canvas.width = CANVAS_EDGE;
    console.log("EDGE", CANVAS_EDGE);
    ctx.fillStyle = "black";
    ctx.fillRect(0, 0, CANVAS_EDGE, CANVAS_EDGE);
}

console.log(window.location.host);
const socket = new WebSocket(`ws://${window.location.host}/websocket`);


init();

setInterval(tick, 1000/FPS);

socket.addEventListener("message", (event) => {
  console.log("Message from server ", event.data + " ");
});


console.log("requestse");