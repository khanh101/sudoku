const screen_height = window.innerHeight;
const screen_width = window.innerWidth;
const cell_size = Math.floor(Math.min(screen_height / 11, screen_width / 9));

let waiting_panel_img = undefined;
let waiting_img = undefined;
let playing_panel_img = undefined;
let block_img = undefined;
let current_img = undefined;
let initial_img = undefined;
let violation_img = undefined;
let value_img_list = [];
let youwin_panel_img = undefined;
function preload() {
    waiting_panel_img = loadImage("assets/waiting_panel.png");
    waiting_img = loadImage("assets/waiting.png");
    playing_panel_img = loadImage("assets/playing_panel.png");
    block_img = loadImage("assets/block.png");
    current_img = loadImage("assets/current.png");
    initial_img = loadImage("assets/initial.png");
    violation_img = loadImage("assets/violation.png");
    for (let num=0; num<10; num++) {
        value_img_list.push(loadImage(`assets/${num}.png`));
    }
    youwin_panel_img = loadImage("assets/youwin_panel.png");
}

const STATE_WAITING = 0;
const STATE_PLAYING = 1;
let state = STATE_WAITING;

function setup() {
    noLoop();
    createCanvas(9 * cell_size, 11 * cell_size);
    waiting_panel_img.resize(9 * cell_size, 2 * cell_size);
    waiting_img.resize(9 * cell_size, 9 * cell_size);
    playing_panel_img.resize(9 * cell_size, 2 * cell_size);
    block_img.resize(3 * cell_size, 3 * cell_size);
    current_img.resize(cell_size, cell_size);
    initial_img.resize(cell_size, cell_size);
    violation_img.resize(cell_size, cell_size);
    for (let num=0; num<value_img_list.length; num++) {
        value_img_list[num].resize(cell_size, cell_size);
    }
    youwin_panel_img.resize(9 * cell_size, 2 * cell_size);

    httpPost("api/new", "json", {}, update_board);
    draw();
}


let current_cell = undefined;
let youwin = undefined;
let current_board = undefined;
let initial_mask = undefined;
let violation_mask = undefined;


function update_board() {
    httpGet("api/view", "json", function (response) {
        console.log(response);
        youwin = response.youwin;
        current_board = response.current_board;
        initial_mask = response.initial_mask;
        violation_mask = response.violation_mask;
        state = STATE_PLAYING;
        draw();
    });
}

function draw() {
    background(200, 200, 200);
    switch (state) {
        case STATE_WAITING:
            draw_panel_waiting();
            draw_board_waiting();
            break;
        case STATE_PLAYING:
            draw_board_playing();
            draw_panel_playing();
            break;
    }
}

function draw_panel_waiting() {
    image(waiting_panel_img, 0, 9 * cell_size);
}

function draw_board_waiting() {
    image(waiting_img, 0, 0);
}

function draw_panel_playing() {
    if (youwin) {
        image(youwin_panel_img, 0, 9 * cell_size);
    } else {
        image(playing_panel_img, 0, 9 * cell_size);
    }
}

function draw_board_playing() {
    if (current_board !== undefined && initial_mask !== undefined && violation_mask !== undefined) {
        // initial
        for (let row = 0; row < 9; row++) {
            for (let col = 0; col < 9; col++) {
                if (initial_mask[row][col]) {
                    const [x, y] = cell_to_pos_tl(row, col);
                    image(initial_img, x, y);
                }
            }
        }
        // block
        for (let row = 0; row < 9; row += 3) {
            for (let col = 0; col < 9; col += 3) {
                const [x, y] = cell_to_pos_tl(row, col);
                image(block_img, x, y);
            }
        }
        // values
        for (let row = 0; row < 9; row++) {
            for (let col = 0; col < 9; col++) {
                const value = current_board[row][col];
                const [x, y] = cell_to_pos_tl(row, col);
                image(value_img_list[value], x, y);
            }
        }
        // violation
        for (let row = 0; row < 9; row++) {
            for (let col = 0; col < 9; col++) {
                if (violation_mask[row][col]) {
                    const [x, y] = cell_to_pos_tl(row, col);
                    image(violation_img, x, y);
                }
            }
        }
        // current
        if (current_cell !== undefined) {
            const [row, col] = current_cell;
            const [x, y] = cell_to_pos_tl(row, col);
            image(current_img, x, y);
        }
    }
}

function mousePressed() {
    current_cell = pos_to_cell(mouseX, mouseY);
    draw();
}

function pos_to_cell(x, y) {
    const col = Math.floor(x / cell_size);
    const row = Math.floor(y / cell_size);
    if (!((0 <= col && col < 9) && (0 <= row && row < 9))) {
        return undefined;
    }
    return [row, col]
}

function cell_to_pos_tl(row, col) {
    const x = col * cell_size;
    const y = row * cell_size;
    return [x, y];
}

function keyPressed() {
    function key_to_value(key) {
        return key - 48;
    }
    if (keyCode === LEFT_ARROW) {
        if (current_cell !== undefined) {
            if (current_cell[1] > 0) {
                current_cell[1] -= 1;
            }
        }
    }
    if (keyCode === RIGHT_ARROW) {
        if (current_cell !== undefined) {
            if (current_cell[1] < 8) {
                current_cell[1] += 1;
            }
        }
    }
    if (keyCode === UP_ARROW) {
        if (current_cell !== undefined) {
            if (current_cell[0] > 0) {
                current_cell[0] -= 1;
            }
        }
    }
    if (keyCode === DOWN_ARROW) {
        if (current_cell !== undefined) {
            if (current_cell[0] < 8) {
                current_cell[0] += 1;
            }
        }
    }
    if (48 <= keyCode && keyCode < 58) {
        const value = key_to_value(keyCode);
        if (current_cell !== undefined) {
            const [row, col] = current_cell;
            httpPost("api/place", "json", {
                row: row,
                col: col,
                value: value,
            }, draw);
        }
    }
    if (keyCode === 88 || keyCode === 8 || keyCode === 46) {
        if (current_cell !== undefined) {
            const [row, col] = current_cell;
            httpPost("api/place", "json", {
                row: row,
                col: col,
                value: 0,
            }, draw);
        }
    }// x
    if (keyCode === 82) {
        httpPost("api/reset", "json", {}, draw);
    }// r
    draw();
}