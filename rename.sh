#!/bin/bash

# Check if an event name was provided as an argument
if [ -z "$1" ]; then
    echo "Usage: ./rename.sh <event_name>"
    exit 1
fi

# Get the event name from the first command-line argument
EVENT_NAME="$1"

# Define your base video directory (e.g., ~/Videos)
BASE_VIDEO_DIR="$HOME/Videos/eventos"

# Construct the full paths using the event name
EVENT_DIR="$BASE_VIDEO_DIR/$EVENT_NAME"
ORIGINAL_DIR="$EVENT_DIR/original"
RENAMED_DIR="$EVENT_DIR/renombrado"

# Check if original directory exists
if [ ! -d "$ORIGINAL_DIR" ]; then
    echo "Error: Directory '$ORIGINAL_DIR' does not exist."
    exit 1
fi

# Create a new directory for the renamed files to keep originals safe
mkdir -p "$RENAMED_DIR"

COUNT=1

# Iterate over files sorted by modification time (oldest first: oldest filmed -> first event)
# Using ls -tr to sort by time, and while read to handle spaces in filenames
while IFS= read -r file; do
    # Ensure the file exists (handles the case where no files match the wildcard)
    if [ -f "$file" ]; then
        # Format the new filename (e.g., EventName-1.mp4)
        new_filename="${EVENT_NAME}-${COUNT}.mp4"
        
        echo "Copying: '$(basename "$file")' -> '$new_filename'"
        # Using cp -p to copy safely and preserve timestamps. 
        # If you prefer to move and rename in place to save space, change 'cp -p' to 'mv'
        cp -p "$file" "$RENAMED_DIR/$new_filename"
        
        ((COUNT++))
    fi
done < <(ls -tr "$ORIGINAL_DIR"/*.mp4 "$ORIGINAL_DIR"/*.MP4 2>/dev/null)

total=$((COUNT - 1))
echo "--------------------------------------------------------"
echo "Done! Processed $total files."
echo "The renamed files are located in: $RENAMED_DIR"
