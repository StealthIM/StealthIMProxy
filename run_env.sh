#!/bin/bash

tmux new-session -d -s stimenv 'cd ../StealthIMDB && make'
tmux select-window -t stimenv:0
tmux new-window 'cd ../StealthIMGroupUser && make'
tmux new-window 'cd ../StealthIMSession && make'
tmux new-window 'cd ../StealthIMUser && make'
tmux new-window 'cd ../StealthIMFileStorage && poetry run python main.py'
tmux new-window 'cd ../StealthIMFileAPI && make'
tmux new-window 'cd ../StealthIMMSAP && make'
tmux attach-session -t stimenv