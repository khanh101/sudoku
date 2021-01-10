import {game} from "./sketch.js"

function interval_access() {
    httpPost("api/interval_access", "json", {
        key: game.key,
    }, function() {
        setTimeout(interval_access, 30000);
    });
}

function setup(response) {
    game.key = response.key;
    interval_access();
    document.getElementById("key").value = response.key;
    document.getElementById("new_board").style.display = "none"
    document.getElementById("new_random").style.display = "none"
    document.getElementById("login_key").style.display = "none"
}

function login_random() {
    const textkey = document.getElementById("key");
    httpPost("api/new", "json", {}, setup);
}

function login_key() {
    const textkey = document.getElementById("key");
    httpPost("api/login", "json", {
        key: textkey.value,
    }, setup);
}

function login_board() {
    const textboard = document.getElementById("board");
    if (textboard.value.length === 0) {
        textboard.value = new Array(81 + 1).join("0");
    }
    httpPost("api/new", "json", {
        board: textboard.value,
    }, setup);
}
window.login_with_board = login_board;
window.login_with_key = login_key;
window.login_with_random = login_random;
