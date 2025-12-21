#!/bin/bash

# Check if an event name was provided as an argument
if [ -z "$1" ]; then
    echo "Usage: ./ffmpeg-script.sh <event_name>"
    exit 1
fi

# Get the event name from the first command-line argument
EVENT_NAME="$1"

# Define your base video directory (e.g., ~/Videos)
BASE_VIDEOS_DIR="$HOME/Videos/eventos"

# Construct the full paths using the event name
EVENT_DIR="$BASE_VIDEOS_DIR/$EVENT_NAME"
ORIGINALES_DIR="$EVENT_DIR/original"
FORMATED_DIR="$EVENT_DIR/formateado"

# Check if the 'EVENT_DIR' directory exists
if [ ! -d "$EVENT_DIR" ]; then
    echo "Error: Source directory '$EVENT_DIR' not found."
    exit 1
fi

echo "Starting video conversion for event: '$EVENT_NAME'"
echo "Source directory: $ORIGINALES_DIR"
echo "Output directory: $FORMATED_DIR"
echo "----------------------------------------------------"

# Initialize a counter for processed files
processed_count=1

# Loop through each .mp4 file in the originales directory
for input_file_path in "$ORIGINALES_DIR"/*.mp4; do
    # Check if a file was actually found (in case no .mp4 files exist)
    # The glob "$ORIGINALES_DIR"/*.mp4 will expand to itself if no matches,
    # so we need to ensure it's a regular file and not the literal pattern.
    if [[ -f "$input_file_path" ]]; then
        filename=$(basename -- "$input_file_path")
        output_file_path="$FORMATED_DIR/${EVENT_NAME}-${processed_count}.mp4"

        echo "Processing: $filename"
        echo "Input: $input_file_path"
        echo "Output: $output_file_path"

        # Run the ffmpeg command
        ffmpeg -v warning \
               -i "$input_file_path" \
               -vf "fps=30" \
               -c:v libx264 \
               -preset medium \
               -crf 23 \
               -movflags +faststart \
			   -an \
               "$output_file_path"

        if [ $? -eq 0 ]; then
            echo "Successfully processed: $filename"
            ((processed_count++))
        else
            echo "Error processing: $filename"
        fi
        echo "----------------------------------------------------"
    fi
done

if [ "$processed_count" -eq 0 ]; then
    echo "No .mp4 files were found or processed in '$ORIGINALES_DIR'."
else
	((processed_count--))
    echo "All conversions complete for event '$EVENT_NAME'! Total files processed: $processed_count"
fi
