let start_time = undefined;
let total_time = undefined;
function hms_to_string(hour, min, sec) {
    function num_to_string(num) {
        if (num < 10) {
            return `0${num}`;
        } else {
            return `${num}`
        }
    }
    return num_to_string(hour) + ":" + num_to_string(min) + ":" + num_to_string(sec)
}

function time_to_string(time) {
    return hms_to_string(time.getHours(), time.getMinutes(), time.getSeconds());
}

function milliseconds_to_string(milliseconds) {
    const seconds = Math.floor(milliseconds / 1000);
    const sec = seconds % 60;
    const minutes = Math.floor(seconds / 60);
    const min = minutes % 60;
    const hours = Math.floor(minutes / 60);
    return hms_to_string(hours, min, sec);
}

function write_timer() {
    const current_time = new Date();
    if (youwin !== true) {
        total_time = current_time - start_time;
    }
    let timer = document.getElementById("timer");
    timer.innerHTML = "";
    timer.innerHTML += `start time: ${time_to_string(start_time)}` + "<br>";
    timer.innerHTML += `current time: ${time_to_string(current_time)}` + "<br>";
    timer.innerHTML += `elapsed time: ${milliseconds_to_string(total_time)}` + "<br>";
    setTimeout(write_timer, 1000);
}

function start_timer() {
    if (state === STATE_PLAYING) {
        start_time = new Date();
        write_timer();
        setTimeout(write_timer, 1000);
        return;
    }
    setTimeout(start_timer, 1000);
}

start_timer();