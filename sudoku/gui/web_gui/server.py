from typing import Any

import flask
from flask import request, jsonify, send_from_directory
import numpy as np

from sudoku import board


def to_jsonifiable(data: Any) -> Any:
    if isinstance(data, tuple) or isinstance(data, list):
        return [to_jsonifiable(element) for element in data]
    if isinstance(data, bool) or isinstance(data, np.bool_):
        return int(data)
    if isinstance(data, np.ndarray):
        return data.tolist()
    return data


class Game:
    board_game: board.Game

    def __init__(self, seed: int):
        self.board_game = board.Game(seed)
        self.app = flask.Flask(__name__)
        self.app.config["DEBUG"] = True


        self.app.route("/game/<path:path>", methods=["GET"])(self.serve_static)
        self.app.route("/api/new_board", methods=["POST"])(self.new_board)
        self.app.route("/api/place", methods=["PUT"])(self.place)
        self.app.route("/api/reset", methods=["PUT"])(self.reset)
        self.app.route("/api/view", methods=["GET"])(self.view)

    def serve_static(self, path):
        return send_from_directory("./static/", path)

    def new_board(self):
        self.board_game.new_board()
        return jsonify()

    def place(self):
        try:
            body = request.json
            row, col = body["row"], body["col"]
            value = body.get("value", 0)
            self.board_game.place((row, col), value)
            return jsonify()
        except Exception:
            return jsonify(), 400

    def reset(self):
        self.board_game.reset()
        return jsonify()

    def view(self):
        return jsonify(to_jsonifiable(self.board_game.view()))

    def run(self):
        self.app.run()
