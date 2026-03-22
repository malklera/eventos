#!/bin/bash

# Check if an event name was provided as an argument
if [ -z "$1" ]; then
    echo "Usage: ./mkdir.sh <event_name>"
    exit 1
fi

# Get the event name from the first command-line argument
EVENT_NAME="$1"

# Define your base video directory (e.g., ~/Videos)
BASE_VIDEOS_DIR="$HOME/Videos/eventos"

# Construct the full paths using the event name
EVENT_DIR="$BASE_VIDEOS_DIR/$EVENT_NAME"
ORIGINAL_DIR="$EVENT_DIR/original"
FORMATED_DIR="$EVENT_DIR/formateado"
CUTTED_DIR="$EVENT_DIR/cortado"
EDITED_DIR="$EVENT_DIR/editado"
CATALOG_DIR="$EVENT_DIR/catalogo"
RENAMED_DIR="$EVENT_DIR/renombrado"

# Check if the 'EVENT_NAME' directory exists
if [ -d "$EVENT_DIR" ]; then
    echo "Error: Source directory '$EVENT_DIR' already exist."
    exit 1
fi

mkdir "$EVENT_DIR"
mkdir "$ORIGINAL_DIR"
mkdir "$FORMATED_DIR"
mkdir "$EDITED_DIR"
mkdir "$CUTTED_DIR"
mkdir "$CATALOG_DIR"
mkdir "$RENAMED_DIR"
