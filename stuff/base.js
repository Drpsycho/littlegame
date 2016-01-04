var ws

window.Player = [];

function makeid() {
    var text = "";
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

    for (var i = 0; i < 5; i++)
        text += possible.charAt(Math.floor(Math.random() * possible.length));

    return text;
}

function getRandomColor() {
    var letters = '0123456789ABCDEF'.split('');
    var color = '#';
    for (var i = 0; i < 6; i++) {
        color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
}

var Player = function(_name, _color, _x, _y) {
    this.Name = _name;
    this.Color = _color;

    this.delta_x = _x;
    this.delta_y = _y;
    this.initialize(this.Color, this.delta_x, this.delta_y);
}

var p = Player.prototype = new createjs.Shape();

p.initialize = function(_color, _d_x, _d_y) {
    this.graphics.beginFill(_color).drawCircle(0, 0, 40);
    this.pixelsPerSecond = 1000;
    this.y = _d_y;
    this.x = _d_x;
}

p.sendToServer = function() {
    ws.send(JSON.stringify({
        User: this.Name,
        U_y: this.delta_y,
        U_x: this.delta_x,
        Status: true,
    }));
};

p.moveUp = function(delta) {
    this.delta_y -= 10.0;
};
p.moveDown = function(delta) {
    this.delta_y += 10.0;
};
p.moveLeft = function(delta) {
    this.delta_x -= 10.0;
};
p.moveRight = function(delta) {
    this.delta_x += 10.0;
};

p.setY = function(delta) {
    this.y = delta;
};

p.setX = function(delta) {
    this.x = delta;
};

p.getName = function() {
    return this.Name;
};


function init() {

    canvas = document.getElementById("testCanvas");
    stage = new createjs.Stage(canvas);

    var Graphics = createjs.Graphics;

    timeCircle = stage.addChild(new Player(makeid(), '#F94F70', 100, 100));

    fpsLabel = new createjs.Text("-- fps", "bold 14px Arial", "#FFF");
    stage.addChild(fpsLabel);
    fpsLabel.x = 10;
    fpsLabel.y = 20;

    createjs.Ticker.timingMode = createjs.Ticker.RAF;
    createjs.Ticker.addEventListener("tick", tick);

    var domain = document.location.hostname + (document.location.port ? ':' + document.location.port : '')

    ws = new WebSocket("ws://" + domain + "/handler");

    ws.onopen = function() {
        timeCircle.sendToServer();
    };

    var players = new Array(timeCircle)

    ws.onmessage = function(e) {

        var msg = JSON.parse(e.data);
        var res;
        res = false;

        players.forEach(function(item, i, arr) {
            if (msg.User == item.getName()) {
                if (msg.Status) {
                    item.setX(msg.U_x);
                    item.setY(msg.U_y);
                    res = true;
                    console.log("user -" + msg.User + msg.U_x + msg.U_y);
                } else {
                    console.log("user -" + msg.User + " REMOVED! ");
                    stage.removeChild( players[i] );
                    delete players[i];
                }
            }
        });

        if ((!res) && (msg.Status)) {
            console.log("added new player");
            var rnd_color = getRandomColor();
            players.push(stage.addChild(new Player(msg.User, rnd_color, msg.U_x, msg.U_y)));
        }
    };

    ws.onclose = function() {
        alert("closed");
    };
}

function tick(event) {
    fpsLabel.text = Math.round(createjs.Ticker.getMeasuredFPS()) + " fps";

    if (key.isPressed('up') || key.isPressed('w')) {
        timeCircle.moveUp(event.delta);
        timeCircle.sendToServer();
    }
    if (key.isPressed('down') || key.isPressed('s')) {
        timeCircle.moveDown(event.delta);
        timeCircle.sendToServer();
    }
    if (key.isPressed('left') || key.isPressed('a')) {
        timeCircle.moveLeft(event.delta);
        timeCircle.sendToServer();
    }
    if (key.isPressed('right') || key.isPressed('d')) {
        timeCircle.moveRight(event.delta);
        timeCircle.sendToServer();
    }

    stage.update(event);
}
