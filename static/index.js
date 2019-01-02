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

    // C: Reverse left
    67: "REVERSE_LEFT",
    // Z: Reverse right
    89: "REVERSE_RIGHT",
}
var directionLabel = document.getElementById("direction-label")

var socket = io();

socket.on("directionchanged", function(newDirection){
    directionLabel.innerText = newDirection;
});

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
        var direction = directions[evt.keyCode];
        if (direction) {
            evt.preventDefault();
            moveDirection(STOP);
        }
    }


    // // MOBILE GYROSCOPE

    // // Event listeners: see https://stackoverflow.com/a/4378439
    // if (window.DeviceOrientationEvent) {
    //     window.addEventListener("deviceorientation", function () {
    //         tilt(event.beta, event.gamma);
    //     }, true);
    // } else if (window.DeviceMotionEvent) {
    //     window.addEventListener('devicemotion', function () {
    //         tilt(event.acceleration.x * 2, event.acceleration.y * 2);
    //     }, true);
    // } else {
    //     window.addEventListener("MozOrientation", function () {
    //         tilt(orientation.x * 50, orientation.y * 50);
    //     }, true);
    // }

    // // FrontBack: Positive: backwards; Negative: forewards
    // // LeftRight: Positive: right; Negative: left
    // function tilt(frontback, leftright) {
    //     // Next lines can contain floating-point errors
    //     //directionLabel.innerText = Math.round(frontback) + ", " + Math.round(leftright);

    //     var direction = null;
    //     // Use certain tresholds to make sure we can stop
    //     if (frontback > 25) {
    //         direction = "BACK";
    //     } else if (frontback < -10) {
    //         direction = "FRONT";
    //     } else if (leftright > 10) {
    //         direction = "RIGHT";
    //     } else if (leftright < -10) {
    //         direction = "LEFT";
    //     } else {
    //         // STOP when we aren't in the area
    //         direction = STOP;
    //     }

    //     moveDirection(direction)
    // }
});


var lastDirection = null;
// GENERAL CONTROLS
// moveDirection sends wanted direction to the server and updates the label
function moveDirection(direction) {
    if (direction == lastDirection) {
        return;
    }
    lastDirection = direction;

    if (socket.disconnected) {
        return;
    }

    socket.emit("direction", {
        "date": new Date(),
        "direction": direction,
    });
}