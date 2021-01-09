function interval_access() {
    httpPost("api/interval_access", "json", {
        key: key,
    }, function() {
        setTimeout(interval_access, 30000);
    });
}

function login_random_and_draw() {
    const textkey = document.getElementById("key");
    httpPost("api/new", "json", {}, function(response) {
        key = response.key;
        textkey.value = key;
        interval_access();
        draw();
    });
}

function login_key_and_draw() {
    const textkey = document.getElementById("key");
    httpPost("api/login", "json", {
        key: textkey.value,
    }, function(response) {
        key = response.key;
        textkey.value = key;
        interval_access();
        draw();
    });
}

function login_board_and_draw() {
    const textkey = document.getElementById("key");
    const textboard = document.getElementById("board");
    if (textboard.value.length === 0) {
        textboard.value = new Array(81 + 1).join("0");
    }
    httpPost("api/new", "json", {
        board: textboard.value,
    }, function(response) {
        key = response.key;
        textkey.value = key;
        interval_access();
        draw();
    });
}



function update_board(cb) {
    httpPost("api/view", "json", {
        key: key,
    }, function (response) {
        if (youwin !== true) {
            youwin = response.youwin;
        }
        current_board = response.current_board;
        initial_mask = response.initial_mask;
        violation_mask = response.violation_mask;
        current_cell = [response.pointer.row, response.pointer.col];
        state = STATE_PLAYING;
        document.getElementById("board").value = get_current_board_string();
        cb();
    });
}

function place(row, col, value) {
    current_explanation = null;
    httpPost("api/place", "json", {
        key: key,
        row: row,
        col: col,
        value: value,
    }, draw);
}

function undo_button() {
    httpPost("api/undo", "json", {
        key: key,
    }, function (response) {
        if (response === null) {
            document.getElementById("undo").textContent = `undo: could not undo`;
            return;
        }
        const {row, col, value} = response;
        current_cell = [row, col]
        place(row, col, 0);
        document.getElementById("undo").textContent = `undo: found {row: ${row}, col: ${col}, value ${value}}`;
    });
}

function implication_button() {
    httpPost("api/implication", "json", {
        key: key,
    }, function (response) {
        if (response === null) {
            document.getElementById("implication").textContent = `implication: could not find`;
            return;
        }
        const {row, col, value} = response;
        current_cell = [row, col]
        place(row, col, value);
        current_explanation = response.explanation;
        document.getElementById("implication").textContent = `implication: found {row: ${row}, col: ${col}} is ${value}`;
    })
}

function pointer(row, col, cb) {
    console.log(row, col)
    httpPost("api/point", "json", {
        key: key,
        row: row,
        col: col,
    }, cb);
}