function image_resize(img, size) {
    img.resize(size[0], size[1]);
    return img;
}

class Board {
    constructor() {
        const pad = 10;
        const screen_height = window.innerHeight;
        const screen_width = window.innerWidth;
        this.size = Math.floor(Math.min((screen_height - 2 * pad) / 11, (screen_width - 2 * pad) / 9));
        this.cell_size = [this.size, this.size];
        this.panel_size = [9 * this.size, 2 * this.size];
        this.block_size = [3 * this.size, 3 * this.size];
        this.board_size = [9 * this.size, 9 * this.size];
        this.canvas_size = [9 * this.size, 11 * this.size];
        this.panel_topleft = [0, 9 * this.size];
        this.board_topleft = [0, 0];
        this.canvas = null;
        this.image = {
            waiting: {
                panel: null,
                board: null,
            },
            playing: {
                panel: null,
                board: {
                    block: null,
                    current: null,
                    initial: null,
                    violation: null,
                    value_list: [],
                    explanation: null,
                },
                youwin_panel: null,
            },
        };
        // game
        this.key = null;
    }

    preload(p5) {
        let self = this;
        this.image.waiting.panel = p5.loadImage("assets/waiting_panel.png", function(img) {
            img.resize(...self.panel_size);
        });
        this.image.waiting.board = p5.loadImage("assets/waiting.png", function(img) {
            img.resize(...self.board_size);
        });
        this.image.playing.panel = p5.loadImage("assets/playing_panel.png", function(img) {
            img.resize(...self.panel_size);
        });
        this.image.playing.board.block = p5.loadImage("assets/block.png", function (img) {
            img.resize(...self.block_size);
        });
        this.image.playing.board.current = p5.loadImage("assets/current.png", function(img) {
            img.resize(...self.cell_size);
        });
        this.image.playing.board.initial = p5.loadImage("assets/initial.png", function(img) {
            img.resize(...self.cell_size);
        });
        this.image.playing.board.violation = p5.loadImage("assets/violation.png", function(img) {
            img.resize(...self.cell_size);
        });
        for (let num = 0; num < this.image.playing.board.value_list.length; num++) {
            this.image.playing.board.value_list.push(p5.loadImage(`assets/${num}.png`, function(img) {
                img.resize(...self.cell_size);
            }));
        }
        this.image.playing.board.explanation = p5.loadImage("assets/explanation.png", function(img) {
            img.resize(...self.cell_size);
        });
        this.image.playing.youwin_panel = p5.loadImage("assets/youwin_panel.png", function(img) {
            img.resize(...self.panel_size);
        });
    }

    setup(p5) {
        this.canvas = p5.createCanvas(...this.canvas_size);
        this.canvas.parent("p5canvas");
        p5.frameRate(2);
    }

    draw_canvas(p5) {
        if (this.key === null) { // waiting
            p5.background(200, 200, 200);
            p5.image(this.image.waiting.panel, ...this.panel_topleft);
            p5.image(this.image.waiting.board, ...this.board_topleft);
        } else { // playing
            let self = this;
            this.update_board(p5, function (response) {
                p5.background(200, 200, 200);
                self.draw_board(
                    p5,
                    response.youwin,
                    response.initial_mask,
                    response.violation_mask,
                    response.current_board,
                    response.pointer,
                )
            });
        }
    }

    draw_board(p5, youwin, initial_mask, violation_mask, current_board, pointer) {
        // panel
        if (youwin) {
            p5.image(this.image.playing.youwin_panel, ...this.panel_topleft);
        } else {
            p5.image(this.image.playing.panel, ...this.panel_topleft);
        }
        // initial
        for (let row = 0; row < 9; row++) {
            for (let col = 0; col < 9; col++) {
                if (initial_mask[row][col]) {
                    p5.image(
                        this.image.playing.board.initial,
                        ...this.cell_to_topleft(row, col),
                    );
                }
            }
        }
        // block
        for (let row = 0; row < 9; row += 3) {
            for (let col = 0; col < 9; col += 3) {
                p5.image(
                    this.image.playing.board.block,
                    ...this.cell_to_topleft(row, col),
                );
            }
        }
        // values
        for (let row = 0; row < 9; row++) {
            for (let col = 0; col < 9; col++) {
                p5.image(
                    this.image.playing.board.value_list[current_board[row][col]],
                    ...this.cell_to_topleft(row, col),
                );
            }
        }
        // violation
        for (let row = 0; row < 9; row++) {
            for (let col = 0; col < 9; col++) {
                if (violation_mask[row][col]) {
                    p5.image(
                        this.image.playing.board.violation,
                        ...this.cell_to_topleft(row, col),
                    )
                }
            }
        }

        // current
        p5.image(
            self.image.playing.current,
            ...self.cell_to_topleft(pointer.row, pointer.col)
        )
    }

    draw_explanation(p5, explanation) {
        for (let i = 0; i < explanation.length; i++) {
            const cell = explanation[i];
            const row = cell.row;
            const col = cell.col;
            p5.image(
                this.image.playing.board.explanation,
                ...this.cell_to_topleft(row, col),
            )
        }
    }

    mousePressed(p5) {
        const cell = this.pos_to_cell(p5.mouseX, p5.mouseY);
        if (cell !== null) {
            let self = this;
            this.point(p5, ...cell, function (response) {
                self.draw_canvas(p5);
            })
        }
    }

    keyPressed(p5) {
        let self = this;

        function number_to_value(key) {
            return key - 48;
        }

        function numpad_to_value(key) {
            return key - 96;
        }

        if ((48 <= p5.keyCode && p5.keyCode < 58) || (96 <= p5.keyCode && p5.keyCode < 106) || p5.keyCode === 88 || p5.keyCode === 8 || p5.keyCode === 46) {
            // number, numpad, x, backspace, delete
            let val = null;
            if (48 <= p5.keyCode && p5.keyCode < 58) {
                val = number_to_value(p5.keyCode);
            }
            if (96 <= p5.keyCode && p5.keyCode < 106) {
                val = numpad_to_value(p5.keyCode);
            }
            val = 0;
            this.place(p5, val, function (response) {
                self.draw_canvas(p5);
            });
        }
        if (p5.keyCode === 85) {
            this.undo(p5, function (response) {
                self.draw_canvas(p5);
            });
        }// u
        if (p5.keyCode === 72) {
            this.implication(p5, function (response) {
                const {row, col, val, explanation} = response;
                self.point(p5, row, col, function(response) {
                   self.place(p5, val, function (response) {
                        self.draw_explanation(p5, explanation);
                   })
                });
            });
        }// h

    }

    // request
    update_board(p5, cb) {
        p5.httpPost("api/view", "json", {
            key: this.key,
        }, cb);
    }

    point(p5, row, col, cb) {
        p5.httpPost("api/point", "json", {
            key: this.key,
            row: row,
            col: col,
        }, cb);
    }

    place(p5, val, cb) {
        p5.httpPost("api/place", "json", {
            key: this.key,
            val: val,
        }, cb);
    }

    undo(p5, cb) {
        p5.httpPost("api/undo", "json", {
            key: key,
        }, cb);
    }

    implication(p5, cb) {
        p5.httpPost("api/implication", "json", {
            key: key,
        }, cb);
    }

    // util
    cell_to_topleft(row, col) {
        const x = col * this.cell_size[0];
        const y = row * this.cell_size[1];
        return [x, y];
    }

    pos_to_cell(x, y) {
        const row = Math.floor(y / this.cell_size[0]);
        const col = Math.floor(x / this.cell_size[1]);
        if (!((0 <= col && col < 9) && (0 <= row && row < 9))) {
            return null;
        }
        return [row, col]
    }
}

export {Board};