from typing import Any

import flask
from flask import request, jsonify
import numpy as np

from sudoku.board import Game

def to_jsonifiable(data: Any) -> Any:
    if isinstance(data, tuple) or isinstance(data, list):
        return [to_jsonifiable(element) for element in data]
    if isinstance(data, bool) or isinstance(data, np.bool_):
        return int(data)
    if isinstance(data, np.ndarray):
        return data.tolist()
    return data


class GameServer:
    game: Game

    def __init__(self, seed: int):
        self.game = Game(seed)
        self.app = flask.Flask(__name__)
        self.app.config["DEBUG"] = True

        self.app.route("/new_board", methods=["POST"])(self.new_board)
        self.app.route("/place", methods=["PUT"])(self.place)
        self.app.route("/reset", methods=["PUT"])(self.reset)
        self.app.route("/view", methods=["GET"])(self.view)

    def new_board(self):
        self.game.new_board()
        return jsonify()

    def place(self):
        try:
            body = request.json
            row, col = body["row"], body["col"]
            value = body.get("value", 0)
            self.game.place((row, col), value)
            return jsonify()
        except Exception:
            return jsonify(), 400

    def reset(self):
        self.game.reset()
        return jsonify()

    def view(self):
        return jsonify(to_jsonifiable(self.game.view()))

    def run(self):
        self.app.run()
