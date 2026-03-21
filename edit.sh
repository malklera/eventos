#!/bin/bash

# --- Configuration ---
# Video dimension
VIDEO_WIDTH=1080
VIDEO_HEIGHT=1920

# Duration of the transition effect in second
TRANSITION_DURATION=1
 
# Common text styling
# Line spacing (pixels between lines)
LINE_SPACING=10
FONT_SIZE=80
# Font color in FFmpeg hex format: [0x|#]RRGGBB[AA]
FONT_COLOR="#E6E70F"
# Border width in pixels (creates outline around each letter)
BORDER_WIDTH=2
# Border color in FFmpeg hex format: [0x|#]RRGGBB[AA]
BORDER_COLOR="#000000"
# Box padding, some fonts get clipped otherwise
BOX_PADDING=8

# Client text styling (only appears with client image)
CLIENT_TEXT_X="(w-tw)/2"
CLIENT_TEXT_Y="(h-th)/1.25"

# Event text styling (appears during filmed video)
TEXT_X="w-tw-20"
TEXT_Y="H-th-20"

# Define your base video directory (e.g., ~/Videos)
BASE_VIDEOS_DIR="$HOME/Videos/eventos"

# Client image display duration (in seconds)
CLIENT_TIME=5
# A helper function to display usage information
usage() {
    echo "Usage: $0 -e <event_name> -m <music> -l <logo_video> [-t \"<event_text>\"] [-i <client_image>] [-c \"<client_text>\"] [-C \"<client_color>\"] [-T \"<text_color>\"] [-f <font>] [-L <left_icon>] [-R <right_icon>]"
    echo ""
    echo "  -e, --event        : Name of the event same name used with ./mkdir <event-name>."
	echo "  -m, --music        : Partial path to music file (~/Videos/eventos/assets/musica/cortado/<your path>)."
    echo "  -l, --logo         : Partial path to music file (~/Videos/eventos/assets/logo/<your path>)."
    echo "  -i, --image        : (Optional) Image file name, file has to be inside <event-name>/."
	echo "  -t, --text         : (Optional) Text to overlay on all videos for this event (enclose in quotes if it has spaces)."
    echo "  -c, --client       : (Optional) Client text to display over client image (enclose in quotes if it has spaces)."
    echo "  -C, --client-color : (Optional) Client text color in hex format (e.g., \"#FFFFFF\" or \"#E6E70F\")."
    echo "  -T, --text-color   : (Optional) Event text color in hex format (e.g., \"#FFFFFF\" or \"#E6E70F\")."
    echo "  -f, --font         : Font filename to use for event text (e.g., MyFont.ttf)."
    echo "  -L, --left         : (Optional) Path to icon image to display to the left of client text."
    echo "  -R, --right        : (Optional) Path to icon image to display to the right of client text."
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
            MUSIC="$BASE_VIDEOS_DIR/assets/musica/cortado/$2"
            shift # past argument
            ;;
        -i|--image)
            CLIENT_IMAGE="$BASE_VIDEOS_DIR/$EVENT_NAME/$2"
            shift # past argument
            ;;
        -c|--client)
            CLIENT_TEXT="$2"
            shift # past argument
            ;;
        -C|--client-color)
            CLIENT_COLOR="$2"
            shift # past argument
            ;;
        -T|--text-color)
            EVENT_COLOR="$2"
            shift # past argument
            ;;
        -l|--logo)
            LOGO="$BASE_VIDEOS_DIR/assets/logo/$2"
            shift # past argument
            ;;
        -f|--font)
            FONT="/usr/local/share/fonts/$2"
            shift # past argument
            ;;
        -L|--left)
            ICON_LEFT="$BASE_VIDEOS_DIR/assets/icono/decoraciones/$2"
            shift # past argument
            ;;
        -R|--right)
            ICON_RIGHT="$BASE_VIDEOS_DIR/assets/icono/decoraciones/$2"
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
# The user will manage what is mandatory.
# if [ -z "$EVENT_NAME" ] || [ -z "$MUSIC" ] || [ -z "$LOGO" ]; then
#     echo "Error: All required arguments (-e, -m, -l) must be provided."
#     usage
# fi

# --- Construct paths using the event name ---
if [ -n "$EVENT_NAME" ]; then
    CUTTED_DIR="$BASE_VIDEOS_DIR/$EVENT_NAME/cortado"
    EDITED_DIR="$BASE_VIDEOS_DIR/$EVENT_NAME/editado"
fi

if [ -n "$MUSIC" ] && [ ! -f "$MUSIC" ]; then
    echo "Error: Music file not found at '$MUSIC'"
    exit 1
fi

if [ -n "$LOGO" ] && [ ! -f "$LOGO" ]; then
    echo "Error: Logo video file not found at '$LOGO'"
    exit 1
fi


if [ -n "$CUTTED_DIR" ] && [ ! -d "$CUTTED_DIR" ]; then
    echo "Error: Manual cuts directory '$CUTTED_DIR' not found."
    echo "Please perform manual cuts in Shotcut first and export videos there."
    exit 1
fi
if [ -n "$FONT" ] && [ ! -f "$FONT" ]; then
    echo "Error: Font file not found at '$FONT'. Please ensure it's installed or provide a correct path."
    echo "You can list available fonts with: fc-list | rg .ttf"
    exit 1
fi

if [ -n "$CLIENT_IMAGE" ] && [ ! -f "$CLIENT_IMAGE" ]; then
    echo "Error: Client image file not found at '$CLIENT_IMAGE'"
    exit 1
fi

if [ -n "$ICON_LEFT" ] && [ ! -f "$ICON_LEFT" ]; then
    echo "Error: Left icon image file not found at '$ICON_LEFT'"
    exit 1
fi

if [ -n "$ICON_RIGHT" ] && [ ! -f "$ICON_RIGHT" ]; then
    echo "Error: Right icon image file not found at '$ICON_RIGHT'"
    exit 1
fi

# --- Validate Icons vs Client Text dependency ---
if [ -z "$CLIENT_TEXT" ]; then
    if [ -n "$ICON_LEFT" ] || [ -n "$ICON_RIGHT" ]; then
        echo "Warning: Icons (--left/-L or --right/-R) were provided but no client text (--client/-c) was specified."
        echo "Icons will be ignored as they are designed to flank the client text."
        ICON_LEFT=""
        ICON_RIGHT=""
    fi
fi

# --- Apply color overrides if provided ---
if [ -n "$CLIENT_COLOR" ]; then
    FONT_COLOR="$CLIENT_COLOR"
fi

if [ -n "$EVENT_COLOR" ]; then
    FONT_COLOR="$EVENT_COLOR"
fi

# --- Get Logo Duration ---
if [ -n "$LOGO" ]; then
    LOGO_DURATION=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "$LOGO" | cut -d'.' -f1)
    if [ -z "$LOGO_DURATION" ]; then
        echo "Error: Could not determine duration for $LOGO."
        exit 1
    fi
fi
 
echo "--- Starting Automated Processing for Event: '$EVENT_NAME' ---"
echo "Text for videos: \"$EVENT_TEXT\""
echo "Music: $MUSIC"
if [ -n "$CLIENT_IMAGE" ]; then
    echo "Client image: $CLIENT_IMAGE (will display for $CLIENT_TIME seconds)"
    if [ -n "$CLIENT_TEXT" ]; then
        echo "Client text: \"$CLIENT_TEXT\""
    fi
    if [ -n "$ICON_LEFT" ]; then
        echo "Left icon: $ICON_LEFT"
    fi
    if [ -n "$ICON_RIGHT" ]; then
        echo "Right icon: $ICON_RIGHT"
    fi
fi
echo "Logo video: $LOGO"
echo "Output to: $EDITED_DIR"
echo "------------------------------------------------------------------"

processed_count=1

# (Logo duration is now calculated earlier) 
# Loop through each .mp4 file from the manually cut videos
for input_video_path in "$CUTTED_DIR"/*.mp4; do
    if [[ -f "$input_video_path" ]]; then

        original_filename=$(basename -- "$input_video_path")
        output_filename="${EVENT_NAME}-${processed_count}.mp4"
        output_file_path="$EDITED_DIR/$output_filename"

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
        if [ -n "$CLIENT_IMAGE" ]; then
            # Case 1: Client Image + Main Video + Logo Video
            PRE_LOGO_VISUAL_DUR=$((CLIENT_TIME + VIDEO_DURATION - TRANSITION_DURATION))  # Length after first xfade
            TOTAL_DURATION=$((PRE_LOGO_VISUAL_DUR + LOGO_DURATION - TRANSITION_DURATION))
            
            CLIENT_TRANSITION_START=$((CLIENT_TIME - TRANSITION_DURATION))
            LOGO_TRANSITION_START=$((PRE_LOGO_VISUAL_DUR - TRANSITION_DURATION))
            
            # Fades apply only to the content part (client + input_video), not logo
            VIDEO_END_FADE_START=$((PRE_LOGO_VISUAL_DUR - TRANSITION_DURATION))
            
            TEXT_FADE_IN_START=$CLIENT_TRANSITION_START
            TEXT_FADE_IN_END=$CLIENT_TIME
            TEXT_FADE_OUT_START=$VIDEO_END_FADE_START
            
            echo "Client: ${CLIENT_TIME}s, Video: ${VIDEO_DURATION}s, Logo: ${LOGO_DURATION}s, Total: ${TOTAL_DURATION}s (Music: ${PRE_LOGO_VISUAL_DUR}s)"
            
            # Single-pass approach: client + input_video + logo with transitions
            # Scale and prepare all video inputs
            FILTER_COMPLEX="[0:v]scale=$VIDEO_WIDTH:$VIDEO_HEIGHT:force_original_aspect_ratio=decrease,pad=$VIDEO_WIDTH:$VIDEO_HEIGHT:(ow-iw)/2:(oh-ih)/2,fps=30,setpts=PTS-STARTPTS[client];"
            FILTER_COMPLEX+="[1:v]scale=$VIDEO_WIDTH:$VIDEO_HEIGHT,fps=30,setpts=PTS-STARTPTS[video];"
            FILTER_COMPLEX+="[2:v]fps=30,setpts=PTS-STARTPTS[logo];"
            
            # First transition: client to input_video
            # shellcheck disable=SC1087
            FILTER_COMPLEX+="[client][video]xfade=transition=fade:duration=$TRANSITION_DURATION:offset=$CLIENT_TRANSITION_START[client_video_raw];"
            
            # Second transition: client_video_raw to logo (removed initial fade-in to fix black thumbnail)
            # shellcheck disable=SC1087
            FILTER_COMPLEX+="[client_video_raw][logo]xfade=transition=fade:duration=$TRANSITION_DURATION:offset=$LOGO_TRANSITION_START[video_combined];"
            
            
            # Add event text overlay (text fades in during transition2: client->video xfade)
            FILTER_COMPLEX+="[video_combined]drawtext=text='$EVENT_TEXT':x=$TEXT_X:y=$TEXT_Y:fontfile=$FONT:fontsize=$FONT_SIZE:fontcolor=$FONT_COLOR:borderw=$BORDER_WIDTH:bordercolor=$BORDER_COLOR:box=1:boxcolor=0x00000000:boxborderw=$BOX_PADDING:line_spacing=$LINE_SPACING:alpha='if(lt(t,$TEXT_FADE_IN_START),0,if(lt(t,$TEXT_FADE_IN_END),(t-$TEXT_FADE_IN_START)/$TRANSITION_DURATION,if(gt(t,$TEXT_FADE_OUT_START),(1-(t-$TEXT_FADE_OUT_START)/$TRANSITION_DURATION),1)))'"
            
            # Estimate text width for icon positioning (if CLIENT_TEXT exists)
            if [ -n "$CLIENT_TEXT" ]; then
                # Average character width 0.37 * font_size for most fonts, this seems to work ok
                CHAR_COUNT=${#CLIENT_TEXT}
                ESTIMATED_TEXT_WIDTH=$(awk "BEGIN {printf \"%.0f\", $CHAR_COUNT * $FONT_SIZE * 0.37}")
                echo "Debug: Text '$CLIENT_TEXT' ($CHAR_COUNT chars) -> Estimated visual width: $ESTIMATED_TEXT_WIDTH px"
            fi
            
            # Add client text overlay if provided (only during client image: 0 to CLIENT_TIME)
            if [ -n "$CLIENT_TEXT" ]; then
                CLIENT_TEXT_FADE_OUT_START=$CLIENT_TRANSITION_START
                # Text appears solid immediately to match the image, only fades out
                FILTER_COMPLEX+=",drawtext=text='$CLIENT_TEXT':x=$CLIENT_TEXT_X:y=$CLIENT_TEXT_Y:fontfile=$FONT:fontsize=$FONT_SIZE:fontcolor=$FONT_COLOR:borderw=$BORDER_WIDTH:bordercolor=$BORDER_COLOR:box=1:boxcolor=0x00000000:boxborderw=$BOX_PADDING:alpha='if(lt(t,$CLIENT_TEXT_FADE_OUT_START),1,if(lt(t,$CLIENT_TIME),(1-(t-$CLIENT_TEXT_FADE_OUT_START)/$TRANSITION_DURATION),0))'"
            fi
            
            # Label the output after text overlays
            FILTER_COMPLEX+="[v_with_text];"
            
            # Determine the next input index based on what icons are provided
            NEXT_INPUT=4  # Start after [3:a] (music)
            CURRENT_STREAM="v_with_text"
            
            # Add left icon overlay if provided
            if [ -n "$ICON_LEFT" ]; then
                ICON_HEIGHT="$FONT_SIZE"  # Scale icon to match font size
                ICON_SPACING=20  # Spacing between icon and text in pixels
                
                # Calculate icon X position: center minus half text width minus icon width minus spacing
                # Use max(0, ...) to ensure icon stays on screen even if text is wide
                ICON_LEFT_X="max(0,W/2-$ESTIMATED_TEXT_WIDTH/2-w-$ICON_SPACING)"
                
                # Scale icon maintaining aspect ratio, then apply ONLY fade out to match text timing
                FILTER_COMPLEX+="[${NEXT_INPUT}:v]scale=-1:${ICON_HEIGHT}:force_original_aspect_ratio=decrease,fade=t=out:st=$CLIENT_TRANSITION_START:d=$TRANSITION_DURATION:alpha=1[icon_left_scaled];"
                # Position left icon to the left of the text
                # Y: Match CLIENT_TEXT_Y formula: (H-h)/1.25 positions at ~75% down
                FILTER_COMPLEX+="[${CURRENT_STREAM}][icon_left_scaled]overlay=x='$ICON_LEFT_X':y='(H-h)/1.25':enable='between(t,0,$CLIENT_TIME)':format=auto[v_with_left];"
                CURRENT_STREAM="v_with_left"
                ((NEXT_INPUT++))
            fi
            
            # Add right icon overlay if provided
            if [ -n "$ICON_RIGHT" ]; then
                ICON_HEIGHT="$FONT_SIZE"  # Scale icon to match font size
                ICON_SPACING=20  # Spacing between icon and text in pixels
                
                # Calculate icon X position: center plus half text width plus spacing
                # Use min(W-w, ...) to ensure icon stays on screen
                ICON_RIGHT_X="min(W-w,W/2+$ESTIMATED_TEXT_WIDTH/2+$ICON_SPACING)"
                
                # Scale icon maintaining aspect ratio, then apply ONLY fade out to match text timing
                FILTER_COMPLEX+="[${NEXT_INPUT}:v]scale=-1:${ICON_HEIGHT}:force_original_aspect_ratio=decrease,fade=t=out:st=$CLIENT_TRANSITION_START:d=$TRANSITION_DURATION:alpha=1[icon_right_scaled];"
                # Position right icon to the right of the text
                # Y: Match CLIENT_TEXT_Y formula: (H-h)/1.25 positions at ~75% down
                FILTER_COMPLEX+="[${CURRENT_STREAM}][icon_right_scaled]overlay=x='$ICON_RIGHT_X':y='(H-h)/1.25':enable='between(t,0,$CLIENT_TIME)':format=auto[v_with_right];"
                CURRENT_STREAM="v_with_right"
            fi
            
            # Set final output stream
            FILTER_COMPLEX+="[${CURRENT_STREAM}]null[v_out];"
            
            # Audio: Music fades in, plays through content, acrossfade to logo audio, logo audio fades out at end
            # Music track: fade in at start, trim to exactly PRE_LOGO_VISUAL_DUR (no fade out, acrossfade handles transition)
            FILTER_COMPLEX+="[3:a]afade=t=in:st=0:d=$TRANSITION_DURATION,atrim=0:$PRE_LOGO_VISUAL_DUR,asetpts=PTS-STARTPTS[music];"
            # Acrossfade to logo audio, then fade out at the very end
            LOGO_FADEOUT_START=$((TOTAL_DURATION - TRANSITION_DURATION))
            # shellcheck disable=SC1087
            FILTER_COMPLEX+="[music][2:a]acrossfade=d=$TRANSITION_DURATION,afade=t=out:st=$LOGO_FADEOUT_START:d=$TRANSITION_DURATION[a_out]"
            
            # Build ffmpeg command with optional icon inputs
            FFMPEG_CMD="ffmpeg -v warning"
            FFMPEG_CMD+=" -loop 1 -t \"$CLIENT_TIME\" -i \"$CLIENT_IMAGE\""
            FFMPEG_CMD+=" -i \"$input_video_path\""
            FFMPEG_CMD+=" -i \"$LOGO\""
            FFMPEG_CMD+=" -i \"$MUSIC\""
            
            # Add icon inputs if provided (loop them like client image)
            if [ -n "$ICON_LEFT" ]; then
                FFMPEG_CMD+=" -loop 1 -t \"$CLIENT_TIME\" -i \"$ICON_LEFT\""
            fi
            if [ -n "$ICON_RIGHT" ]; then
                FFMPEG_CMD+=" -loop 1 -t \"$CLIENT_TIME\" -i \"$ICON_RIGHT\""
            fi
            
            FFMPEG_CMD+=" -filter_complex \"$FILTER_COMPLEX\""
            FFMPEG_CMD+=" -map \"[v_out]\""
            FFMPEG_CMD+=" -map \"[a_out]\""
            FFMPEG_CMD+=" -t \"$TOTAL_DURATION\""
            FFMPEG_CMD+=" -c:v libx264"
            FFMPEG_CMD+=" -preset medium"
            FFMPEG_CMD+=" -crf 23"
            FFMPEG_CMD+=" -c:a aac"
            FFMPEG_CMD+=" -b:a 192k"
            FFMPEG_CMD+=" -movflags +faststart"
            FFMPEG_CMD+=" -f mp4"
            FFMPEG_CMD+=" \"$output_file_path\""
            
            # Execute the command
            eval "$FFMPEG_CMD"
                   
		else
            # Case 2: Main Video + Logo Video (no client)
            PRE_LOGO_VISUAL_DUR=$VIDEO_DURATION  # Duration before logo
            TOTAL_DURATION=$((PRE_LOGO_VISUAL_DUR + LOGO_DURATION - TRANSITION_DURATION))
            
            LOGO_TRANSITION_START=$((PRE_LOGO_VISUAL_DUR - TRANSITION_DURATION))
            VIDEO_END_FADE_START=$((PRE_LOGO_VISUAL_DUR - TRANSITION_DURATION))
            
            echo "Video: ${VIDEO_DURATION}s, Logo: ${LOGO_DURATION}s, Total: ${TOTAL_DURATION}s (Music: ${PRE_LOGO_VISUAL_DUR}s)"
            
            # Scale and prepare video inputs
            FILTER_COMPLEX="[0:v]scale=$VIDEO_WIDTH:$VIDEO_HEIGHT,fps=30,setpts=PTS-STARTPTS[video_raw];"
            FILTER_COMPLEX+="[1:v]fps=30,setpts=PTS-STARTPTS[logo];"
            
            # Transition from input_video to logo (removed initial fade-in to fix black thumbnail)
            # shellcheck disable=SC1087
            FILTER_COMPLEX+="[video_raw][logo]xfade=transition=fade:duration=$TRANSITION_DURATION:offset=$LOGO_TRANSITION_START[video_combined];"
            
            # Add text overlay (text only appears during content portion, not logo)
            FILTER_COMPLEX+="[video_combined]drawtext=text='$EVENT_TEXT':x=$TEXT_X:y=$TEXT_Y:fontfile=$FONT:fontsize=$FONT_SIZE:fontcolor=$FONT_COLOR:borderw=$BORDER_WIDTH:bordercolor=$BORDER_COLOR:box=1:boxcolor=0x00000000:boxborderw=$BOX_PADDING:line_spacing=$LINE_SPACING:alpha='if(gt(t,$VIDEO_END_FADE_START),(1-(t-$VIDEO_END_FADE_START)/$TRANSITION_DURATION),1)'[v_out];"
            
            # Audio: Music fades in, plays through content, acrossfade to logo audio, logo audio fades out at end
            # Music track: fade in at start, trim to exactly PRE_LOGO_VISUAL_DUR (no fade out, acrossfade handles transition)
            FILTER_COMPLEX+="[2:a]afade=t=in:st=0:d=$TRANSITION_DURATION,atrim=0:$PRE_LOGO_VISUAL_DUR,asetpts=PTS-STARTPTS[music];"
            # Acrossfade to logo audio, then fade out at the very end
            LOGO_FADEOUT_START=$((TOTAL_DURATION - TRANSITION_DURATION))
            # shellcheck disable=SC1087
            FILTER_COMPLEX+="[music][1:a]acrossfade=d=$TRANSITION_DURATION,afade=t=out:st=$LOGO_FADEOUT_START:d=$TRANSITION_DURATION[a_out]"

            ffmpeg -v warning \
                   -i "$input_video_path" \
                   -i "$LOGO" \
                   -i "$MUSIC" \
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

		# I think this is better because of how large the command is
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
    echo "Final videos are in: $EDITED_DIR"
fi
