#!/bin/bash

# Configuration
WIDTH=4
HEIGHT=4

# Give yourself time to click into the Discord text box
echo "⚠️  Focus the Discord text box NOW."
echo "⚠️  Starting typing in 5 seconds..."
sleep 5

for (( x=0; x<$WIDTH; x++ )); do
    for (( y=0; y<$HEIGHT; y++ )); do
        echo "$x/%WIDTH - $y/%HEIGHT"
        CMD=(
            "type /draw"
            "key 15:1 15:0"       # Tab
            "type $x"
            "key 15:1 15:0"       # Tab
            "type $y"
            "key 15:1 15:0"       # Tab
            "type #${x}${y}0000"
            "key 28:1 28:0"       # Enter
        )

        for cmd in "${CMD[@]}"; do
            ydotool $cmd
            sleep 0.1
        done
        sleep 35
        
    done
done

echo "Done!"