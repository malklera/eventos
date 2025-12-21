#!/bin/bash

# --- Configuration ---
# Video dimension
VIDEO_WIDTH=1080
VIDEO_HEIGHT=1920

# Duration of the transition effect in second
TRANSITION_DURATION=1
 
# Client text styling (only appears with client image)
CLIENT_TEXT_X="(w-tw)/2"
CLIENT_TEXT_Y="(h-th)/1.25"
CLIENT_FONT_FILE="/usr/local/share/fonts/Kapakana-VariableFont_wght.ttf"
CLIENT_FONT_SIZE=100
CLIENT_FONT_COLOR="#FFFFFF"
CLIENT_BORDER_WIDTH=3
CLIENT_BORDER_COLOR="#000000"

# Event text styling (appears during filmed video)

# Text position
TEXT_X="w-tw-20"
# H-text_height-20
TEXT_Y="H-th-20"
# Line spacing (pixels between lines)
LINE_SPACING=10
# Full path to font file
FONT_FILE="/usr/local/share/fonts/Kapakana-VariableFont_wght.ttf"
FONT_SIZE=80
# Font color in FFmpeg hex format: [0x|#]RRGGBB[AA]
FONT_COLOR="#E6E70F"
# Border width in pixels (creates outline around each letter)
BORDER_WIDTH=2
# Border color in FFmpeg hex format: [0x|#]RRGGBB[AA]
BORDER_COLOR="#000000"

# Define your base video directory (e.g., ~/Videos)
BASE_VIDEOS_DIR="$HOME/Videos/eventos"

# Client image display duration (in seconds)
CLIENT_TIME=5

# --- Argument Parsing ---
# Initialize variables with default or empty values
EVENT_NAME=""
EVENT_TEXT=""
MUSIC_PATH=""
CLIENT_PATH=""
CLIENT_TEXT=""
CLIENT_COLOR_ARG=""
EVENT_COLOR_ARG=""
LOGO_PATH=""

# A helper function to display usage information
usage() {
    echo "Usage: $0 -e <event_name> -t \"<event_text>\" -m <music_path> -l <logo_video_path> [-i <client_image_path>] [-c \"<client_text>\"] [-C \"<client_color>\"] [-T \"<text_color>\"]"
    echo ""
    echo "  -e, --event        : Name of the event ."
    echo "                       This is used for directory paths and output filenames."
    echo "  -t, --text         : Text to overlay on all videos for this event (enclose in quotes if it has spaces)."
    echo "  -m, --music        : Full path to the background music file (e.g., ~/Videos/eventos/assets/event_jingle.mp3)."
    echo "  -i, --image        : (Optional) Full path to client image to display at the beginning (e.g., ~/Videos/eventos/evento1/client_logo.png)."
    echo "  -c, --client       : (Optional) Client text to display over client image (enclose in quotes if it has spaces)."
    echo "  -C, --client-color : (Optional) Client text color in hex format (e.g., \"#FFFFFF\" or \"#E6E70F\")."
    echo "  -T, --text-color   : (Optional) Event text color in hex format (e.g., \"#FFFFFF\" or \"#E6E70F\")."
    echo "  -l, --logo         : Full path to logo video to display at the end (e.g., ~/Videos/eventos/assets/logo1/logo1.mp4)."
    echo ""
    exit 2
}

# Parse named arguments
while [[ "$#" -gt 0 ]]; do
    case "$1" in
        -e|--event)
            EVENT_NAME="$2"
            shift # past argument
            ;;
        -t|--text)
            EVENT_TEXT="$2"
            shift # past argument
            ;;

        -m|--music)
            MUSIC_PATH="$2"
            shift # past argument
            ;;
        -i|--image)
            CLIENT_PATH="$2"
            shift # past argument
            ;;
        -c|--client)
            CLIENT_TEXT="$2"
            shift # past argument
            ;;
        -C|--client-color)
            CLIENT_COLOR_ARG="$2"
            shift # past argument
            ;;
        -T|--text-color)
            EVENT_COLOR_ARG="$2"
            shift # past argument
            ;;
        -l|--logo)
            LOGO_PATH="$2"
            shift # past argument
            ;;
        -h|--help)
            usage
            ;;
        *)
            echo "Unknown parameter passed: $1"
            usage
            ;;
    esac
    shift # past argument or value
done

# --- Validate Required Arguments ---
if [ -z "$EVENT_NAME" ] || [ -z "$EVENT_TEXT" ] || [ -z "$MUSIC_PATH" ] || [ -z "$LOGO_PATH" ]; then
    echo "Error: All required arguments (-e, -t, -m, -l) must be provided."
    usage
fi

# --- Construct paths using the event name ---
CUTTED_DIR="$BASE_VIDEOS_DIR/$EVENT_NAME/cortado"
EDITADOS_DIR="$BASE_VIDEOS_DIR/$EVENT_NAME/editado"

if [ ! -f "$MUSIC_PATH" ]; then
    echo "Error: Music file not found at '$MUSIC_PATH'"
    exit 1
fi
if [ ! -d "$CUTTED_DIR" ]; then
    echo "Error: Manual cuts directory '$CUTTED_DIR' not found."
    echo "Please perform manual cuts in Shotcut first and export videos there."
    exit 1
fi
if [ ! -f "$FONT_FILE" ]; then
    echo "Error: Font file not found at '$FONT_FILE'. Please ensure it's installed or provide a correct path."
    echo "You can list available fonts with: fc-list | rg .ttf"
    exit 1
fi

if [ -n "$CLIENT_PATH" ] && [ ! -f "$CLIENT_PATH" ]; then
    echo "Error: Client image file not found at '$CLIENT_PATH'"
    exit 1
fi

if [ ! -f "$LOGO_PATH" ]; then
    echo "Error: Logo video file not found at '$LOGO_PATH'"
    exit 1
fi

# --- Apply color overrides if provided ---
if [ -n "$CLIENT_COLOR_ARG" ]; then
    CLIENT_FONT_COLOR="$CLIENT_COLOR_ARG"
fi

if [ -n "$EVENT_COLOR_ARG" ]; then
    FONT_COLOR="$EVENT_COLOR_ARG"
fi
 
echo "--- Starting Automated Processing for Event: '$EVENT_NAME' ---"
echo "Text for videos: \"$EVENT_TEXT\""
echo "Music: $MUSIC_PATH"
if [ -n "$CLIENT_PATH" ]; then
    echo "Client image: $CLIENT_PATH (will display for $CLIENT_TIME seconds)"
    if [ -n "$CLIENT_TEXT" ]; then
        echo "Client text: \"$CLIENT_TEXT\""
    fi
fi
echo "Logo video: $LOGO_PATH"
echo "Output to: $EDITADOS_DIR"
echo "------------------------------------------------------------------"

processed_count=1

# --- Get Logo Duration if provided ---
LOGO_DURATION=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "$LOGO_PATH" | cut -d'.' -f1)
if [ -z "$LOGO_DURATION" ]; then
	echo "Error: Could not determine duration for $LOGO_PATH."
	exit 1
fi
 
# Loop through each .mp4 file from the manually cut videos
for input_video_path in "$CUTTED_DIR"/*.mp4; do
    if [[ -f "$input_video_path" ]]; then

        original_filename=$(basename -- "$input_video_path")
        output_filename="${EVENT_NAME}-${processed_count}.mp4"
        output_file_path="$EDITADOS_DIR/$output_filename"

        echo "Processing original: $original_filename"
        echo "Input video: $input_video_path"
        echo "Output will be: $output_filename"
        echo "Full Output Path: $output_file_path"

        # --- Get Video Duration ---
        VIDEO_DURATION=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "$input_video_path" | cut -d'.' -f1)
        if [ -z "$VIDEO_DURATION" ]; then
            echo "Error: Could not determine duration for $input_video_path. Skipping."
            continue # Skip to next file
        fi
       
        # Determine the structure based on what's provided
        if [ -n "$CLIENT_PATH" ]; then
            # Case 1: Client Image + Main Video + Logo Video
            # Logo is already upscaled to match input_video (1080x1920 @ 30fps)
            CONTENT_DURATION=$((CLIENT_TIME + VIDEO_DURATION))  # Duration that music covers (client + input_video)
            TOTAL_DURATION=$((CONTENT_DURATION + LOGO_DURATION))
            
            CLIENT_TRANSITION_START=$((CLIENT_TIME - TRANSITION_DURATION))
            LOGO_TRANSITION_START=$((CONTENT_DURATION - TRANSITION_DURATION))
            
            # Fades apply only to the content part (client + input_video), not logo
            VIDEO_END_FADE_START=$((CONTENT_DURATION - TRANSITION_DURATION))
            AUDIO_END_FADE_START=$((CONTENT_DURATION - TRANSITION_DURATION))
            
            TEXT_FADE_IN_START=$CLIENT_TRANSITION_START
            TEXT_FADE_IN_END=$CLIENT_TIME
            TEXT_FADE_OUT_START=$VIDEO_END_FADE_START
            
            echo "Client: ${CLIENT_TIME}s, Video: ${VIDEO_DURATION}s, Logo: ${LOGO_DURATION}s, Total: ${TOTAL_DURATION}s (Music: ${CONTENT_DURATION}s)"
            
            # Single-pass approach: client + input_video + logo with transitions
            # Scale and prepare all video inputs
            FILTER_COMPLEX="[0:v]scale=$VIDEO_WIDTH:$VIDEO_HEIGHT:force_original_aspect_ratio=decrease,pad=$VIDEO_WIDTH:$VIDEO_HEIGHT:(ow-iw)/2:(oh-ih)/2,fps=30,setpts=PTS-STARTPTS[client];"
            FILTER_COMPLEX+="[1:v]scale=$VIDEO_WIDTH:$VIDEO_HEIGHT,fps=30,setpts=PTS-STARTPTS[video];"
            # Logo video is already at correct resolution and fps, just ensure timestamps
            FILTER_COMPLEX+="[2:v]fps=30,setpts=PTS-STARTPTS[logo];"
            
            # First transition: client to input_video
            FILTER_COMPLEX+="[client][video]xfade=transition=fade:duration=$TRANSITION_DURATION:offset=$CLIENT_TRANSITION_START[client_video_raw];"
            
            # Apply fade in to content portion (no fade out, to allow xfade to logo)
            FILTER_COMPLEX+="[client_video_raw]fade=t=in:st=0:d=$TRANSITION_DURATION[client_video];"
            
            # Second transition: faded client_video to logo
            FILTER_COMPLEX+="[client_video][logo]xfade=transition=fade:duration=$TRANSITION_DURATION:offset=$LOGO_TRANSITION_START[video_combined];"
            
            # Add event text overlay (text fades in during transition2: client->video xfade)
            FILTER_COMPLEX+="[video_combined]drawtext=text='$EVENT_TEXT':x=$TEXT_X:y=$TEXT_Y:fontfile=$FONT_FILE:fontsize=$FONT_SIZE:fontcolor=$FONT_COLOR:borderw=$BORDER_WIDTH:bordercolor=$BORDER_COLOR:line_spacing=$LINE_SPACING:alpha='if(lt(t,$TEXT_FADE_IN_START),0,if(lt(t,$TEXT_FADE_IN_END),(t-$TEXT_FADE_IN_START)/$TRANSITION_DURATION,if(gt(t,$TEXT_FADE_OUT_START),(1-(t-$TEXT_FADE_OUT_START)/$TRANSITION_DURATION),1)))'"
            
            # Add client text overlay if provided (only during client image: 0 to CLIENT_TIME)
            if [ -n "$CLIENT_TEXT" ]; then
                CLIENT_TEXT_FADE_OUT_START=$CLIENT_TRANSITION_START
                FILTER_COMPLEX+=",drawtext=text='$CLIENT_TEXT':x=$CLIENT_TEXT_X:y=$CLIENT_TEXT_Y:fontfile=$CLIENT_FONT_FILE:fontsize=$CLIENT_FONT_SIZE:fontcolor=$CLIENT_FONT_COLOR:borderw=$CLIENT_BORDER_WIDTH:bordercolor=$CLIENT_BORDER_COLOR:alpha='if(lt(t,$TRANSITION_DURATION),t/$TRANSITION_DURATION,if(lt(t,$CLIENT_TEXT_FADE_OUT_START),1,if(lt(t,$CLIENT_TIME),(1-(t-$CLIENT_TEXT_FADE_OUT_START)/$TRANSITION_DURATION),0)))'"
            fi
            
            FILTER_COMPLEX+="[v_out];"
            
            # Audio: Music fades in, plays through content, acrossfade to logo audio, logo audio fades out at end
            # Music track: fade in at start, trim to exactly CONTENT_DURATION (no fade out, acrossfade handles transition)
            FILTER_COMPLEX+="[3:a]afade=t=in:st=0:d=$TRANSITION_DURATION,atrim=0:$CONTENT_DURATION,asetpts=PTS-STARTPTS[music];"
            # Acrossfade to logo audio, then fade out at the very end
            LOGO_FADEOUT_START=$((TOTAL_DURATION - TRANSITION_DURATION))
            FILTER_COMPLEX+="[music][2:a]acrossfade=d=$TRANSITION_DURATION,afade=t=out:st=$LOGO_FADEOUT_START:d=$TRANSITION_DURATION[a_out]"
            
            ffmpeg -v warning \
                   -loop 1 -t "$CLIENT_TIME" -i "$CLIENT_PATH" \
                   -i "$input_video_path" \
                   -i "$LOGO_PATH" \
                   -i "$MUSIC_PATH" \
                   -filter_complex "$FILTER_COMPLEX" \
                   -map "[v_out]" \
                   -map "[a_out]" \
                   -t "$TOTAL_DURATION" \
                   -c:v libx264 \
                   -preset medium \
                   -crf 23 \
                   -c:a aac \
                   -b:a 192k \
                   -movflags +faststart \
                   -f mp4 \
                   "$output_file_path"
                   
		else
            # Case 2: Main Video + Logo Video (no client)
            # Logo is already upscaled to match input_video (1080x1920 @ 30fps)
            CONTENT_DURATION=$VIDEO_DURATION  # Duration that music covers
            TOTAL_DURATION=$((CONTENT_DURATION + LOGO_DURATION))
            
            LOGO_TRANSITION_START=$((CONTENT_DURATION - TRANSITION_DURATION))
            VIDEO_END_FADE_START=$((CONTENT_DURATION - TRANSITION_DURATION))
            AUDIO_END_FADE_START=$((CONTENT_DURATION - TRANSITION_DURATION))
            
            echo "Video: ${VIDEO_DURATION}s, Logo: ${LOGO_DURATION}s, Total: ${TOTAL_DURATION}s (Music: ${CONTENT_DURATION}s)"
            
            # Scale and prepare video inputs
            FILTER_COMPLEX="[0:v]scale=$VIDEO_WIDTH:$VIDEO_HEIGHT,fps=30,setpts=PTS-STARTPTS[video_raw];"
            # Logo video is already at correct resolution and fps
            FILTER_COMPLEX+="[1:v]fps=30,setpts=PTS-STARTPTS[logo];"
            
            # Apply fade in to video BEFORE combining with logo (no fade out)
            FILTER_COMPLEX+="[video_raw]fade=t=in:st=0:d=$TRANSITION_DURATION[video];"
            
            # Transition from input_video to logo
            FILTER_COMPLEX+="[video][logo]xfade=transition=fade:duration=$TRANSITION_DURATION:offset=$LOGO_TRANSITION_START[video_combined];"
            
            # Add text overlay (text only appears during content portion, not logo)
            FILTER_COMPLEX+="[video_combined]drawtext=text='$EVENT_TEXT':x=$TEXT_X:y=$TEXT_Y:fontfile=$FONT_FILE:fontsize=$FONT_SIZE:fontcolor=$FONT_COLOR:borderw=$BORDER_WIDTH:bordercolor=$BORDER_COLOR:line_spacing=$LINE_SPACING:alpha='if(lt(t,$TRANSITION_DURATION),t/$TRANSITION_DURATION,if(gt(t,$VIDEO_END_FADE_START),(1-(t-$VIDEO_END_FADE_START)/$TRANSITION_DURATION),1))'[v_out];"
            
            # Audio: Music fades in, plays through content, acrossfade to logo audio, logo audio fades out at end
            # Music track: fade in at start, trim to exactly CONTENT_DURATION (no fade out, acrossfade handles transition)
            FILTER_COMPLEX+="[2:a]afade=t=in:st=0:d=$TRANSITION_DURATION,atrim=0:$CONTENT_DURATION,asetpts=PTS-STARTPTS[music];"
            # Acrossfade to logo audio, then fade out at the very end
            LOGO_FADEOUT_START=$((TOTAL_DURATION - TRANSITION_DURATION))
            FILTER_COMPLEX+="[music][1:a]acrossfade=d=$TRANSITION_DURATION,afade=t=out:st=$LOGO_FADEOUT_START:d=$TRANSITION_DURATION[a_out]"

            ffmpeg -v warning \
                   -i "$input_video_path" \
                   -i "$LOGO_PATH" \
                   -i "$MUSIC_PATH" \
                   -filter_complex "$FILTER_COMPLEX" \
                   -map "[v_out]" \
                   -map "[a_out]" \
                   -t "$TOTAL_DURATION" \
                   -c:v libx264 \
                   -preset medium \
                   -crf 23 \
                   -c:a aac \
                   -b:a 192k \
                   -movflags +faststart \
                   -f mp4 \
                   "$output_file_path"
        fi

        if [ $? -eq 0 ]; then
			((processed_count++))
            echo "Successfully processed '$original_filename' to '$output_filename'"
        else
            echo "Error processing: $original_filename (Check FFmpeg output above)"
        fi
        echo "------------------------------------------------------------------"
    fi
done

if [ "$processed_count" -eq 0 ]; then
    echo "No .mp4 files were found or processed in '$CUTTED_DIR'."
    echo "Ensure you've made the Shotcut manual cuts and saved videos to that directory."
else
	((processed_count--))
    echo "--- All automated conversions complete for event '$EVENT_NAME'! Total files processed: $processed_count ---"
    echo "Final videos are in: $EDITADOS_DIR"
fi
