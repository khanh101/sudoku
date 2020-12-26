import time
from typing import Optional

import flask
from flask import request, jsonify, send_from_directory

from sudoku import board



class Game:
    board_game: Optional[board.Game]

    def __init__(self):
        self.board_game = None
        self.app = flask.Flask(__name__)
        self.app.config["DEBUG"] = True


        self.app.route("/<path:path>", methods=["GET"])(self.serve_static)
        self.app.route("/api/view", methods=["GET"])(self.view)
        self.app.route("/api/new", methods=["POST"])(self.new)
        self.app.route("/api/place", methods=["POST"])(self.place)
        self.app.route("/api/reset", methods=["POST"])(self.reset)
        self.app.route("/api/implication", methods=["GET"])(self.implication)

    def serve_static(self, path):
        return send_from_directory("./static/", path)

    def new(self):
        self.board_game = board.Game(int(time.time()))
        return jsonify()

    def implication(self):
        '''
        {
            row: int,
            col: int,
            value: int
        }
        '''
        if self.board_game is None:
            return jsonify(), 400
        implication = self.board_game.implication()
        if implication is None:
            return jsonify(None)
        return jsonify(implication)

    def view(self):
        '''
        {
            youwin: bool,
            current_board: list of list of int,
            initial_mask: list of list of bool,
            violation_mask: list of list of bool,
        }
        '''
        if self.board_game is None:
            return jsonify(), 400
        return jsonify(self.board_game.view().marshal())

    def place(self):
        '''
        {
            "row": 1,
            "col": 2,
            "value": 3,
        }
        '''
        if self.board_game is None:
            return jsonify(), 400
        try:
            body = request.json
            row, col = body["row"], body["col"]
            value = body.get("value", 0)
            self.board_game.place((row, col), value)
            return jsonify()
        except Exception:
            return jsonify(), 400

    def reset(self):
        if self.board_game is None:
            return jsonify(), 400
        self.board_game.reset()
        return jsonify()

    def run(self, *args, **kwargs):
        self.app.run(*args, **kwargs)
