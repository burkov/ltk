#!/bin/bash
set -e

xmodmap -e 'keycode 66 = '         # disable caps
xmodmap -e 'keycode 112 = '        # disable pg up
xmodmap -e 'keycode 117 = '        # disable pg down
xmodmap -e 'keycode 134 = Super_L' # swap left and right super keys to fix ubuntu bug
xmodmap -e 'keycode 63 = Left'     # vim-like controls
xmodmap -e 'keycode 106 = Down'    # vim-like controls
xmodmap -e 'keycode 110 = Up'      # vim-like controls
xmodmap -e 'keycode 112 = Right'   # vim-like controls
xmodmap -e 'keycode 86 = Home'     #
xmodmap -e 'keycode 82 = End'      #
