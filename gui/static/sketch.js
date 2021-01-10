import {Board} from "./board.js";

let game = new Board();

const sketch = new p5(function(p5) {
    p5.preload = function() {
        game.preload(p5);
    }
    p5.setup = function () {
        game.setup(p5);
    }
    p5.draw = function () {
        game.draw_canvas(p5);
    }
});

export {game, sketch};