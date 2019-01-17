#!/bin/bash
set -e

xmodmap -e 'keycode 66 = '         # disable caps
xmodmap -e 'keycode 111 = '        # disable arrow up
xmodmap -e 'keycode 113 = '        # disable arrow left
xmodmap -e 'keycode 114 = '        # disable arrow right
xmodmap -e 'keycode 115 = '        # disable end
xmodmap -e 'keycode 116 = '        # disable arrow down
xmodmap -e 'keycode 117 = '        # disable pg down
xmodmap -e 'keycode 134 = Super_L' # swap left and right super keys to fix ubuntu bug

xmodmap -e 'keycode 63 = Left Left'     # vim-like controls
xmodmap -e 'keycode 106 = Down Down'    # vim-like controls
xmodmap -e 'keycode 110 = Up Up'        # vim-like controls
xmodmap -e 'keycode 112 = Right Right'  # vim-like controls
xmodmap -e 'keycode 86 = Home Home'     #
xmodmap -e 'keycode 82 = End End'       #