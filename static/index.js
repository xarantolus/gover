var STOP = "STOP";
// directions is a map of keycode:directionName
var directions = {
    // Arrow keys
    37: "LEFT",
    38: "FORWARD",
    39: "RIGHT",
    40: "REVERSE",

    // WASD
    65: "LEFT",
    87: "FORWARD",
    68: "RIGHT",
    83: "REVERSE",

    // Q: Pivot left
    81: "PIVOT_LEFT",
    // E: Pivot right
    69: "PIVOT_RIGHT",

    // Z: Reverse left
    89: "REVERSE_LEFT",
    // C: Reverse right
    67: "REVERSE_RIGHT",

    // Space bar is the brake
    32: STOP,
}

var directionDeg = {
    "LEFT": -20,
    "RIGHT": 20,

    "FORWARD": 0,
    "REVERSE": 180,

    "PIVOT_LEFT": -120,
    "PIVOT_RIGHT": 120,

    "REVERSE_LEFT": 20,
    "REVERSE_RIGHT": -20,

    STOP: 0,
}

var directionLabel = document.getElementById("direction-label")

var socket = io();

socket.on("directionchanged", function (newDirection) {
    directionLabel.innerText = newDirection;
    rotateCar(newDirection);
});

var frontSensor = document.getElementById("front-cm")
socket.on("sensor:front", function (frontUpdate) {
    var text = "No data available";
    if (frontUpdate.ok) {
        text = frontUpdate.distance.toFixed(2).toString() + "cm"
    }

    frontSensor.innerText = text;
})

// Register events when connected
socket.on("connect", function () {
    directionLabel.innerText = "Connected!";

    // DESKTOP KEYDOWN AND KEYUP

    // keydown starts the rover in the specified direction
    window.onkeydown = function (evt) {
        var direction = directions[evt.keyCode];
        if (direction) {
            evt.preventDefault();
            moveDirection(direction);
        }
    }

    // keyup stops the rover
    window.onkeyup = function (evt) {
        moveDirection(STOP);
    }
    window.onbeforeunload = window.onkeyup;


    // MOBILE GYROSCOPE

    // Event listeners: see https://stackoverflow.com/a/4378439
    if (window.DeviceOrientationEvent) {
        window.addEventListener("deviceorientation", function () {
            tilt(event.beta, event.gamma);
        }, true);
    } else if (window.DeviceMotionEvent) {
        window.addEventListener('devicemotion', function () {
            tilt(event.acceleration.x * 2, event.acceleration.y * 2);
        }, true);
    } else {
        window.addEventListener("MozOrientation", function () {
            tilt(orientation.x * 50, orientation.y * 50);
        }, true);
    }

    // FrontBack: Positive: backwards; Negative: forewards
    // LeftRight: Positive: right; Negative: left
    function tilt(frontback, leftright) {
        // Next lines can contain floating-point errors
        //directionLabel.innerText = Math.round(frontback) + ", " + Math.round(leftright);

        var direction = null;
        // Use certain tresholds to make sure we can stop
        if (frontback > 25) {
            direction = "REVERSE";
        } else if (frontback < -15) {
            direction = "FORWARD";
        } else if (leftright > 15) {
            direction = "RIGHT";
        } else if (leftright < -15) {
            direction = "LEFT";
        } else {
            // STOP when we aren't in the area
            direction = STOP;
        }

        moveDirection(direction)
    }
});

socket.on('disconnect', function () {
    directionLabel.innerText = "Disconnected";
    frontSensor.innerText =  "No data available";
});

var lastDirection = null;
var lastSendTime = new Date();
// GENERAL CONTROLS
// moveDirection sends wanted direction to the server and updates the label
function moveDirection(direction) {
    if (socket.disconnected) {
        return;
    }

    if (new Date() - lastSendTime < 500 && direction == lastDirection) { // Only send every half second if it's the same
        return;
    }
    lastDirection = direction;
    lastSendTime = new Date();


    socket.emit("direction", {
        "date": lastSendTime,
        "direction": direction,
    });
}

var carImage = document.getElementById("car-image")

function rotateCar(newDir) {
    var newDeg = directionDeg[newDir];
    
    if (newDeg != null) {
        var transfromString = "rotate(" + newDeg + "deg)";
        console.log(transfromString);

        carImage.style.webkitTransform = transfromString;
        carImage.style.MozTransform = transfromString;
        carImage.style.msTransform = transfromString;
        carImage.style.OTransform = transfromString;
        carImage.style.transform = transfromString;
    }
}