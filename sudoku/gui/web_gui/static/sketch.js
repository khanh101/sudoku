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
}

const STATE_WAITING = 0;
const STATE_PLAYING = 1;
const STATE_QUIT = 2;
const STATE_ENDED = 3;
let state = STATE_WAITING;
let current_cell = undefined;

let youwin = false;
let initial_board = undefined;
let current_board = undefined;

function update_board() {
    httpGet("api/view", "json", function(viewResponse) {
        youwin = viewResponse.youwin;
        initial_board = viewResponse.initial_board;
        current_board = viewResponse.current_board;
        state = STATE_PLAYING;
        if (youwin) {
            state = STATE_ENDED;
        }
        draw();
    });
}
function draw() {
    background(200, 200, 200);
    switch (state) {
        case STATE_WAITING:
            image(waiting_img, 0, 0);
            image(waiting_panel_img, 0, 9 * cell_size);
            httpPost("api/new_board", "json", {}, update_board);
            break;
        case STATE_PLAYING:
            draw_board();
            image(playing_panel_img, 0, 9 * cell_size);
            break;
        case STATE_ENDED:
            draw_board();
            image(youwin_panel_img, 0, 9 * cell_size);
            break;
    }
}

function draw_board() {
    if (initial_board !== undefined && current_board !== undefined) {
        // initial
        for (let row = 0; row < 9; row++) {
            for (let col = 0; col < 9; col++) {
                if (initial_board[row][col] !== 0) {
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
        const violation_board = get_violation_board();
        for (let row = 0; row < 9; row++) {
            for (let col = 0; col < 9; col++) {
                if (violation_board[row][col] !== 0) {
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

function get_violation_cell_list(row, col) {
    let violation_cell_list = []
    for (let r=0; r<9; r++) {
        if (r !== row) {
            violation_cell_list.push([r, col]);
        }
    }
    for (let c=0; c<9; c++) {
        if (c !== col) {
            violation_cell_list.push([row, c]);
        }
    }
    const btlr = 3 * Math.floor(row / 3);
    const btlc = 3 * Math.floor(col / 3);
    for (let r=0; r<3; r++) {
        for (let c=0; c<3; c++) {
            if ((btlr + r === row) || (btlc + c === col)) {
                continue;
            }
            violation_cell_list.push([btlr+r, btlc+c]);
        }
    }
    return violation_cell_list;
}

function get_violation_board() {
    let violation_board = [];
    for (let row=0; row<9; row++) {
        violation_board.push([]);
        for (let col=0; col<9; col++) {
            violation_board[row].push(0)
        }
    }
    for (let row=0; row<9; row++) {
        for (let col = 0; col < 9; col++) {
            if (violation_board[row][col] !== 0) {
                continue;
            }
            const value = current_board[row][col];
            if (value !== 0) {
                const violation_cell_list = get_violation_cell_list(row, col);
                for (let vcellidx=0; vcellidx<violation_cell_list.length; vcellidx++) {
                    const [vrow, vcol] = violation_cell_list[vcellidx];
                    if (current_board[vrow][vcol] === value) {
                        violation_board[row][col] = 1;
                        violation_board[vrow][vcol] = 1;
                    }
                }
            }
        }
    }
    return violation_board;
}

function mousePressed() {
    switch (state) {
        case STATE_WAITING:
            break;
        case STATE_PLAYING:
            current_cell = pos_to_cell(mouseX, mouseY);
            draw();
            break;
        case STATE_ENDED:
            break;
    }
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
    switch (state) {
        case STATE_WAITING:
            break;
        case STATE_PLAYING:
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
                if (current_cell != undefined) {
                    const [row, col] = current_cell;
                    httpPost("api/place", "json", {
                        row: row,
                        col: col,
                        value: value,
                    }, update_board);
                }
            }
            if (keyCode === 88 || keyCode === 8 || keyCode === 46) {
                if (current_cell != undefined) {
                    const [row, col] = current_cell;
                    httpPost("api/place", "json", {
                        row: row,
                        col: col,
                        value: 0,
                    }, update_board);
                }
            }// x
            if (keyCode === 82) {
                httpPost("api/reset", "json", {}, update_board);
            }// r
            draw();
            break;
        case STATE_ENDED:
            if (keyCode === 78) {
                state = STATE_WAITING;
                draw();
            }// n
            if (keyCode === 82) {
                httpPost("api/reset", "json", {}, update_board);
            } // r
            break;
    }

}