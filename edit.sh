#!/bin/bash
#
# Helper Functions for Logging
log_info() {
    echo -e "\e[34m[INFO]\e[0m $1" # Blue
}

log_success() {
    echo -e "\e[32m[SUCCESS]\e[0m $1" # Green
}

log_error() {
    echo -e "\e[31m[ERROR]\e[0m $1" # Red
}

log_warning() {
    echo -e "\e[33m[WARNING]\e[0m $1" # Yellow
}

# --- Configuration ---
# Video dimension
VIDEO_WIDTH=1080
VIDEO_HEIGHT=1920

# Duration of the transition effect in second
TRANSITION_DURATION=1

# Define your base video directory (e.g., ~/Videos)
BASE_VIDEOS_DIR="$HOME/Videos/eventos"

# Client image display duration (in seconds)
IMAGE_TIME=2


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
TEXT_X="(w-tw)/2"
TEXT_Y="H-th-100"  # Moved up from the bottom slightly for better look when centered
TEXT_MARGIN=100

# A helper function to display usage information
usage() {
    echo "Usage: $0 -e <event_name> -m <music> -l <logo_video> [-t \"<event_text>\"] [-i <client_image>] [-c \"<client_text>\"] [-u \"<upper_text>\"] [-d \"<down_text>\"] [-C \"<client_color>\"] [-T \"<text_color>\"] [-f <font>] [-L <left_icon>] [-R <right_icon>]"
    echo ""
    echo "  -e, --event        : Name of the event same name used with ./mkdir <event-name>."
    echo "  -l, --logo         : Partial path to logo image file (~/Videos/eventos/assets/logo/<your path>)."
	echo "  -m, --music        : (Optional) Partial path to music file (~/Videos/eventos/assets/musica/cortado/<your path>)."
    echo "  -i, --image        : (Optional) Client image file name, file has to be inside <event-name>/."
	echo "  -t, --text         : (Optional) Text to overlay on all videos for this event (enclose in quotes if it has spaces)."
    echo "  -c, --client       : (Optional) Client text to display over client image (enclose in quotes if it has spaces)."
    echo "  -u, --upper        : (Optional) Upper line of client text (manual split, bypasses auto-wrap)."
    echo "  -d, --down         : (Optional) Lower line of client text (manual split, bypasses auto-wrap)."
    echo "  -C, --client-color : (Optional) Client text color in hex format (e.g., \"#FFFFFF\" or \"#E6E70F\")."
    echo "  -T, --text-color   : (Optional) Event text color in hex format (e.g., \"#FFFFFF\" or \"#E6E70F\")."
    echo "  -f, --font         : (Optional) Font filename to use for event text (e.g., MyFont.ttf)."
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
            CLIENT_IMAGE_NAME="$2"
            shift # past argument
            ;;
        -c|--client)
            CLIENT_TEXT="$2"
            shift # past argument
            ;;
        -u|--upper)
            CLIENT_TEXT_UP="$2"
            shift # past argument
            ;;
        -d|--down)
            CLIENT_TEXT_DOWN="$2"
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
if [ -z "$EVENT_NAME" ] || [ -z "$LOGO" ]; then
    log_error "Error: -e <event-name> and -l <path-logo> must be provided."
    usage
fi

# --- Construct paths using the event name ---
if [ -n "$EVENT_NAME" ]; then
    CUTTED_DIR="$BASE_VIDEOS_DIR/$EVENT_NAME/cortado"
    EDITED_DIR="$BASE_VIDEOS_DIR/$EVENT_NAME/editado"
fi

if [ -n "$CLIENT_IMAGE_NAME" ]; then
    CLIENT_IMAGE="$BASE_VIDEOS_DIR/$EVENT_NAME/$CLIENT_IMAGE_NAME"
fi

if [ -n "$MUSIC" ] && [ ! -f "$MUSIC" ]; then
    log_error "Error: Music file not found at '$MUSIC'"
    exit 1
fi

if [ -n "$LOGO" ] && [ ! -f "$LOGO" ]; then
    log_error "Error: Logo video file not found at '$LOGO'"
    exit 1
fi


if [ -n "$CUTTED_DIR" ] && [ ! -d "$CUTTED_DIR" ]; then
    log_error "Error: Manual cuts directory '$CUTTED_DIR' not found."
    log_info "Please perform manual cuts in Shotcut first and export videos there."
    exit 1
fi
if [ -z "$FONT" ]; then
    FONT=$(fc-match -f "%{file}" sans 2>/dev/null)
fi

if [ -z "$FONT" ] || [ ! -f "$FONT" ]; then
    log_error "Error: Font file not found at '$FONT'. Please ensure it's installed or provide a correct path."
    log_info "You can list available fonts with: fc-list | rg .ttf"
    exit 1
fi

if [ -n "$CLIENT_IMAGE" ] && [ ! -f "$CLIENT_IMAGE" ]; then
    log_error "Error: Client image file not found at '$CLIENT_IMAGE'"
    exit 1
fi

if [ -n "$ICON_LEFT" ] && [ ! -f "$ICON_LEFT" ]; then
    log_error "Error: Left icon image file not found at '$ICON_LEFT'"
    exit 1
fi

if [ -n "$ICON_RIGHT" ] && [ ! -f "$ICON_RIGHT" ]; then
    log_error "Error: Right icon image file not found at '$ICON_RIGHT'"
    exit 1
fi

# --- Validate Icons vs Client Text dependency ---
if [ -z "$CLIENT_TEXT" ]; then
    if [ -n "$ICON_LEFT" ] || [ -n "$ICON_RIGHT" ]; then
        log_warning "Warning: Icons (--left/-L or --right/-R) were provided but no client text (--client/-c) was specified."
        log_info "Icons will be ignored as they are designed to flank the client text."
        ICON_LEFT=""
        ICON_RIGHT=""
    fi
fi

# --- Apply color overrides if provided ---
CLIENT_FONT_COLOR="${CLIENT_COLOR:-$FONT_COLOR}"
EVENT_FONT_COLOR="${EVENT_COLOR:-$FONT_COLOR}"

# --- Text Wrapping Logic ---
# Function to wrap text based on estimated character width
wrap_text() {
    local raw_text="$1"
    local max_w=$((VIDEO_WIDTH - TEXT_MARGIN))
    # Font size is 80. Use 0.7 factor (approx 56px) for an even safer wrap.
    # We now use physical newlines as FFmpeg handles them better through eval.
    local avg_char_w=$((FONT_SIZE * 7 / 10))
    local max_chars=$((max_w / avg_char_w))
    
    # Use fold to wrap with physical newlines
    echo "$raw_text" | fold -s -w "$max_chars"
}

if [ -n "$EVENT_TEXT" ]; then
    EVENT_TEXT=$(wrap_text "$EVENT_TEXT")
fi

# --- Handle Upper and Down text flags ---
if [ -n "$CLIENT_TEXT_UP" ] || [ -n "$CLIENT_TEXT_DOWN" ]; then
    # If upper/down are provided, they take precedence over -c
    if [ -n "$CLIENT_TEXT_UP" ] && [ -n "$CLIENT_TEXT_DOWN" ]; then
        CLIENT_TEXT="${CLIENT_TEXT_UP}"$'\n'"${CLIENT_TEXT_DOWN}"
    elif [ -n "$CLIENT_TEXT_UP" ]; then
        CLIENT_TEXT="$CLIENT_TEXT_UP"
    else
        CLIENT_TEXT="$CLIENT_TEXT_DOWN"
    fi
    # When using explicit upper/down, skip auto-wrapping to give user control
    SKIP_WRAP_CLIENT=true
fi

if [ -n "$CLIENT_TEXT" ] && [ "$SKIP_WRAP_CLIENT" != true ]; then
    CLIENT_TEXT=$(wrap_text "$CLIENT_TEXT")
fi

# Estimate max text width in pixels for icon positioning
if [ -n "$CLIENT_TEXT" ]; then
    LONGEST=0
    while IFS= read -r line; do
        len=${#line}
        [ "$len" -gt "$LONGEST" ] && LONGEST=$len
    done <<< "$CLIENT_TEXT"
    ESTIMATED_TEXT_WIDTH=$(( LONGEST * FONT_SIZE * 7 / 10 ))
fi

echo "--- Starting Automated Processing for Event: '$EVENT_NAME' ---"
if [ -n "$EVENT_TEXT" ]; then
	echo "Text for videos: \"$EVENT_TEXT\""
fi

if [ -n "$MUSIC" ]; then
	echo "Music: $MUSIC"
fi

if [ -n "$CLIENT_IMAGE" ]; then
    echo "Client image: $CLIENT_IMAGE"
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
echo "Logo image: $LOGO"
echo "Output to: $EDITED_DIR"
echo "------------------------------------------------------------------"

processed_count=0

# (Logo duration is now calculated earlier) 
# Loop through each .mp4 file from the manually cut videos
for input_video_path in "$CUTTED_DIR"/*.mp4; do
    if [[ -f "$input_video_path" ]]; then

        original_filename=$(basename -- "$input_video_path")
        output_filename="$original_filename"
        output_file_path="$EDITED_DIR/$output_filename"

        echo "Full input path: $input_video_path"
        echo "Full output path: $output_file_path"

        # --- Get Video Duration ---
        VIDEO_DURATION=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "$input_video_path" | cut -d'.' -f1)
        if [ -z "$VIDEO_DURATION" ]; then
            log_error "Error: Could not determine duration for $input_video_path. Skipping."
            continue
        fi

        # Determine the structure based on what's provided
        if [ -n "$CLIENT_IMAGE" ]; then
            # Case 1: Client Image + Main Video + Logo Video
            # Timeline: client(IMAGE_TIME w/ fade-out) -> video(VIDEO_DURATION w/ slide-up + slide-out) -> logo(IMAGE_TIME w/ fade-in)
            PRE_LOGO_VISUAL_DUR=$((IMAGE_TIME + VIDEO_DURATION))  # Total before logo
            TOTAL_DURATION=$((PRE_LOGO_VISUAL_DUR + IMAGE_TIME))  # No overlap (concat, not xfade)
            
            CLIENT_FADE_OUT_START=$((IMAGE_TIME - TRANSITION_DURATION))
            VIDEO_SLIDE_OUT_START=$((VIDEO_DURATION - TRANSITION_DURATION))
            
            # Fades apply only to the content part (client + input_video), not logo
            VIDEO_END_FADE_START=$((PRE_LOGO_VISUAL_DUR - TRANSITION_DURATION))
            

            
            # Divide the entire video duration (including intro/outro transitions) into 5 parts.
            # We add 2 * TRANSITION_DURATION to the total pool to account for the 2 xfade overlaps.
            P_LEN=$(awk "BEGIN {print ($VIDEO_DURATION + 2 * $TRANSITION_DURATION) / 5.0}")
            
            PART1_START=0
            PART1_END=$P_LEN
            
            # Xfade 1 setup: Part 1 and Part 2 overlap by TRANSITION_DURATION
            PART2_START=$(awk "BEGIN {print $PART1_END - $TRANSITION_DURATION}")
            PART2_END=$(awk "BEGIN {print $PART2_START + $P_LEN}")
            
            # Parts 2, 3, and 4 abut each other sequentially
            PART3_START=$PART2_END
            PART3_END=$(awk "BEGIN {print $PART3_START + $P_LEN}")
            
            PART4_START=$PART3_END
            PART4_END=$(awk "BEGIN {print $PART4_START + $P_LEN}")
            
            # Xfade 2 setup: Part 4 and Part 5 overlap by TRANSITION_DURATION
            PART5_START=$(awk "BEGIN {print $PART4_END - $TRANSITION_DURATION}")
            PART5_END=$VIDEO_DURATION  # Should be exactly PART5_START + P_LEN if math is perfect
            
            # Event text appears when main video starts (after client)
            TEXT_FADE_IN_START=$((IMAGE_TIME - TRANSITION_DURATION))
            TEXT_FADE_IN_END=$IMAGE_TIME
            TEXT_FADE_OUT_START=$VIDEO_END_FADE_START
            
            echo "Client: ${IMAGE_TIME}s, Video: ${VIDEO_DURATION}s, Logo: ${IMAGE_TIME}s, Total: ${TOTAL_DURATION}s"
            
            # Scale and prepare all video inputs
            # Client image: fade out to black at the end
            # shellcheck disable=SC1087
            FILTER_COMPLEX="[0:v]scale=$VIDEO_WIDTH:$VIDEO_HEIGHT:force_original_aspect_ratio=decrease,pad=$VIDEO_WIDTH:$VIDEO_HEIGHT:(ow-iw)/2:(oh-ih)/2,fps=30,setpts=PTS-STARTPTS,fade=t=out:st=$CLIENT_FADE_OUT_START:d=$TRANSITION_DURATION[client];"
            FILTER_COMPLEX+="[1:v]scale=$VIDEO_WIDTH:$VIDEO_HEIGHT,fps=30,setpts=PTS-STARTPTS[video];"
            # Logo: scale + fade-in from black
            FILTER_COMPLEX+="[2:v]scale=$VIDEO_WIDTH:$VIDEO_HEIGHT:force_original_aspect_ratio=decrease,pad=$VIDEO_WIDTH:$VIDEO_HEIGHT:(ow-iw)/2:(oh-ih)/2,fps=30,setpts=PTS-STARTPTS,settb=AVTB,fade=t=in:st=0:d=$TRANSITION_DURATION[logo];"
            
            # Slide-up entrance + slide-out exit for main video (both over TRANSITION_DURATION)
            # Slide-up: y goes H -> 0 in first TD seconds
            # Slide-out: y goes 0 -> -H in last TD seconds (slides up off-screen)
            FILTER_COMPLEX+="color=c=black:s=${VIDEO_WIDTH}x${VIDEO_HEIGHT}:d=${VIDEO_DURATION}:r=30[video_bg];"
            FILTER_COMPLEX+="[video_bg][video]overlay=x=0:y='if(lt(t,$TRANSITION_DURATION),H-H*t/$TRANSITION_DURATION,if(gt(t,$VIDEO_SLIDE_OUT_START),-H*(t-$VIDEO_SLIDE_OUT_START)/$TRANSITION_DURATION,0))':format=auto[video_slide_raw];"
            
            # --- Inner video effects ---
            # Split the video into Part 1, Parts 2-4, and Part 5 (for xfades)
            FILTER_COMPLEX+="[video_slide_raw]split=3[raw1][raw2to4][raw5];"
            
            # Part 1: From the beginning (including slide-up) until PART1_END
            FILTER_COMPLEX+="[raw1]trim=start=0:end=$PART1_END,setpts=PTS-STARTPTS[part1];"
            
            # Parts 2-4: Starting from PART2_START until PART4_END
            FILTER_COMPLEX+="[raw2to4]trim=start=$PART2_START:end=$PART4_END,setpts=PTS-STARTPTS[part_rest_raw];"
            
            # Part 5: Starting from PART5_START
            FILTER_COMPLEX+="[raw5]trim=start=$PART5_START,setpts=PTS-STARTPTS[part5];"
            
            # --- Part 2 Effect: Quad-split (2x2 grid) ---
            HALF_W=$((VIDEO_WIDTH / 2))
            HALF_H=$((VIDEO_HEIGHT / 2))
            FILTER_COMPLEX+="[part_rest_raw]split=2[pr_main][pr_quad_src];"
            FILTER_COMPLEX+="[pr_quad_src]scale=${HALF_W}:${HALF_H},split=4[q1][q2][q3][q4];"
            FILTER_COMPLEX+="[q1][q2]hstack[top_row];"
            FILTER_COMPLEX+="[q3][q4]hstack[bot_row];"
            FILTER_COMPLEX+="[top_row][bot_row]vstack[quad];"
            
            # Overlay quad on normal video only during Part 2's timeframe (0 to P_LEN)
            FILTER_COMPLEX+="[pr_main][quad]overlay=enable='between(t,0,$P_LEN)':eof_action=pass[part_rest_step1];"
            
            # --- Transition: Frei0r Distort0r between Part 2 and Part 3 ---
            DISTORT_START=$(awk "BEGIN {print $P_LEN - $TRANSITION_DURATION / 2.0}")
            DISTORT_END=$(awk "BEGIN {print $P_LEN + $TRANSITION_DURATION / 2.0}")
            
            FILTER_COMPLEX+="[part_rest_step1]split=2[pr_step15_main][distort_src];"
            FILTER_COMPLEX+="[distort_src]frei0r=filter_name=distort0r:filter_params=0.5|0.005|y|0.25[distorted];"
            FILTER_COMPLEX+="[pr_step15_main][distorted]overlay=enable='between(t,$DISTORT_START,$DISTORT_END)':eof_action=pass[part_rest_step15];"
            
            # --- Part 3 Effect: Slow/Fast time warp ---
            # Part 3 plays in slow-motion for the first half, speeding up for the second half
            P3_LOCAL_START=$(awk "BEGIN {print $PART3_START - $PART2_START}")
            P3_LOCAL_END=$(awk "BEGIN {print $PART3_END - $PART2_START}")
            
            FILTER_COMPLEX+="[part_rest_step15]split=2[pr_step2_main][p3_src];"
            FILTER_COMPLEX+="[p3_src]trim=start=$P3_LOCAL_START:end=$P3_LOCAL_END,setpts=PTS-STARTPTS[p3_trim];"
            
            # T is current time in seconds. First 1/4 of source duration plays at 0.5x speed (takes 1/2 of output P_LEN).
            # Remaining 3/4 plays at 1.5x speed (takes remaining 1/2 of output P_LEN).
            # PTS mapping securely warps timestamps precisely to fit exactly within P_LEN without breaking overlaps!
            FILTER_COMPLEX+="[p3_trim]setpts='if(lt(T, $P_LEN/4), 2*PTS + $P3_LOCAL_START/TB, (2/3)*PTS + ($P_LEN/3 + $P3_LOCAL_START)/TB)'[p3_warp];"
            
            # Overlay exactly during Part 3's timeframe
            FILTER_COMPLEX+="[pr_step2_main][p3_warp]overlay=enable='between(t,$P3_LOCAL_START,$P3_LOCAL_END)':eof_action=pass[part_rest_step2];"
            
            # --- Transition: Spin effect between Part 3 and Part 4 ---
            # Creates a fast 360-degree dizzy spin precisely across the explicit frame cut boundary
            SPIN_START=$(awk "BEGIN {print $P_LEN * 2 - $TRANSITION_DURATION / 2.0}")
            SPIN_END=$(awk "BEGIN {print $P_LEN * 2 + $TRANSITION_DURATION / 2.0}")
            
            # Spins dynamically reaching exactly 2*PI (a full 360 circle) completing squarely normalizing back to 0
            FILTER_COMPLEX+="[part_rest_step2]rotate=a='if(between(t,$SPIN_START,$SPIN_END), 2*PI*(t-$SPIN_START)/$TRANSITION_DURATION, 0)':ow=iw:oh=ih:c=black[part_rest_step25];"
            
            # --- Part 4 Effect: 3x3 grid ---
            # Part 4 timestamps need to be shifted because part_rest starts at PART2_START
            P4_LOCAL_START=$(awk "BEGIN {print $PART4_START - $PART2_START}")
            P4_LOCAL_END=$(awk "BEGIN {print $PART4_END - $PART2_START}")
            THIRD_W=$((VIDEO_WIDTH / 3))
            THIRD_H=$((VIDEO_HEIGHT / 3))
            
            FILTER_COMPLEX+="[part_rest_step25]split=2[pr_step3_main][pr_grid3_src];"
            FILTER_COMPLEX+="[pr_grid3_src]scale=${THIRD_W}:${THIRD_H},split=9[g1][g2][g3][g4][g5][g6][g7][g8][g9];"
            FILTER_COMPLEX+="[g1][g2][g3]hstack=inputs=3[row1];"
            FILTER_COMPLEX+="[g4][g5][g6]hstack=inputs=3[row2];"
            FILTER_COMPLEX+="[g7][g8][g9]hstack=inputs=3[row3];"
            FILTER_COMPLEX+="[row1][row2][row3]vstack=inputs=3[grid3];"
            
            # Overlay 3x3 grid only during Part 4's explicit timeframe
            FILTER_COMPLEX+="[pr_step3_main][grid3]overlay=enable='between(t,$P4_LOCAL_START,$P4_LOCAL_END)':eof_action=pass[part_rest_step4];"
            
            # --- Part 4 Reverse Boomerang Effect ---
            # Cut Part 4 in half: play the first half normally, and the first half reversed during the second half
            P4_LOCAL_MID=$(awk "BEGIN {print $P4_LOCAL_START + ($P4_LOCAL_END - $P4_LOCAL_START) / 2.0}")
            
            FILTER_COMPLEX+="[part_rest_step4]split=2[pr_step4_main][p4_to_rev];"
            # Trim the exact floating-point first half, reverse it, shift PTS to map precisely onto the second half
            FILTER_COMPLEX+="[p4_to_rev]trim=start=$P4_LOCAL_START:end=$P4_LOCAL_MID,setpts=PTS-STARTPTS,reverse,setpts=PTS+($P4_LOCAL_MID/TB)[p4_rev];"
            
            # Overlay reversed snippet safely during the second half, strictly enduring the xfade
            FILTER_COMPLEX+="[pr_step4_main][p4_rev]overlay=enable='between(t,$P4_LOCAL_MID,$P4_LOCAL_END)':eof_action=pass[part_rest];"
            
            # Xfade 1: Part 1 and Parts 2-4
            XFADE_OFFSET1=$(awk "BEGIN {print $PART1_END - $TRANSITION_DURATION}")
            FILTER_COMPLEX+="[part1][part_rest]xfade=transition=fade:duration=$TRANSITION_DURATION:offset=$XFADE_OFFSET1[video_step1];"
            
            # Xfade 2: video_step1 (which ends exactly at PART4_END) and Part 5
            XFADE_OFFSET2=$(awk "BEGIN {print $PART4_END - $TRANSITION_DURATION}")
            FILTER_COMPLEX+="[video_step1][part5]xfade=transition=fade:duration=$TRANSITION_DURATION:offset=$XFADE_OFFSET2[video_slide];"
            
            # Concatenate: client (with fade-out) -> main video (with slide-up + quad effect)
            # settb=AVTB normalizes timebase so xfade inputs match
            FILTER_COMPLEX+="[client][video_slide]concat=n=2:v=1:a=0,settb=AVTB[client_video_raw];"
            
            # Concatenate: client_video_raw -> logo (with fade-in)
            FILTER_COMPLEX+="[client_video_raw][logo]concat=n=2:v=1:a=0[video_combined];"
            
            
            # --- Dynamic Text Overlay Logic (Case 1) ---
            CUR_V="video_combined"

            # 1. EVENT_TEXT (during filmed video)
            if [ -n "$EVENT_TEXT" ]; then
                # Use literal newlines; avoiding complex escaping that causes 'n' artifacts
                ESC_EVENT_TEXT=$(echo -n "$EVENT_TEXT" | sed 's/\\/\\\\/g; s/[:'\'']/\\&/g')
                FILTER_COMPLEX+="[$CUR_V]drawtext=text='$ESC_EVENT_TEXT':x=$TEXT_X:y=$TEXT_Y:fontfile=$FONT:fontsize=$FONT_SIZE:fontcolor=$EVENT_FONT_COLOR:borderw=$BORDER_WIDTH:bordercolor=$BORDER_COLOR:box=1:boxcolor=0x00000088:boxborderw=$BOX_PADDING:line_spacing=$LINE_SPACING:text_align=center:alpha='if(lt(t,$TEXT_FADE_IN_START),0,if(lt(t,$TEXT_FADE_IN_END),(t-$TEXT_FADE_IN_START)/$TRANSITION_DURATION,if(gt(t,$TEXT_FADE_OUT_START),(1-(t-$TEXT_FADE_OUT_START)/$TRANSITION_DURATION),1)))'[v_ev];"
                CUR_V="v_ev"
            fi

            # 2. CLIENT_TEXT (during introduction image)
            if [ -n "$CLIENT_TEXT" ]; then
                CLIENT_TEXT_FADE_OUT_START=$CLIENT_FADE_OUT_START
                if [ -n "$CLIENT_TEXT_UP" ] || [ -n "$CLIENT_TEXT_DOWN" ]; then
                    # Manual split: two separate drawtext filters for total control as requested
                    if [ -n "$CLIENT_TEXT_UP" ]; then
                        FF_UP=$(echo -n "$CLIENT_TEXT_UP" | sed 's/\\/\\\\/g; s/[:'\'']/\\&/g')
                        FILTER_COMPLEX+="[$CUR_V]drawtext=text='$FF_UP':x=$CLIENT_TEXT_X:y=$CLIENT_TEXT_Y-60:fontfile=$FONT:fontsize=$FONT_SIZE:fontcolor=$CLIENT_FONT_COLOR:borderw=$BORDER_WIDTH:bordercolor=$BORDER_COLOR:box=1:boxcolor=0x00000088:boxborderw=$BOX_PADDING:alpha='if(lt(t,$CLIENT_TEXT_FADE_OUT_START),1,if(lt(t,$IMAGE_TIME),(1-(t-$CLIENT_TEXT_FADE_OUT_START)/$TRANSITION_DURATION),0))'[v_up];"
                        CUR_V="v_up"
                    fi
                    if [ -n "$CLIENT_TEXT_DOWN" ]; then
                        FF_DOWN=$(echo -n "$CLIENT_TEXT_DOWN" | sed 's/\\/\\\\/g; s/[:'\'']/\\&/g')
                        FILTER_COMPLEX+="[$CUR_V]drawtext=text='$FF_DOWN':x=$CLIENT_TEXT_X:y=$CLIENT_TEXT_Y+60:fontfile=$FONT:fontsize=$FONT_SIZE:fontcolor=$CLIENT_FONT_COLOR:borderw=$BORDER_WIDTH:bordercolor=$BORDER_COLOR:box=1:boxcolor=0x00000088:boxborderw=$BOX_PADDING:alpha='if(lt(t,$CLIENT_TEXT_FADE_OUT_START),1,if(lt(t,$IMAGE_TIME),(1-(t-$CLIENT_TEXT_FADE_OUT_START)/$TRANSITION_DURATION),0))'[v_dw];"
                        CUR_V="v_dw"
                    fi
                else
                    # Fallback to single text (auto-wrapped if needed)
                    FF_CL=$(echo -n "$CLIENT_TEXT" | sed 's/\\/\\\\/g; s/[:'\'']/\\&/g')
                    FILTER_COMPLEX+="[$CUR_V]drawtext=text='$FF_CL':x=$CLIENT_TEXT_X:y=$CLIENT_TEXT_Y:fontfile=$FONT:fontsize=$FONT_SIZE:fontcolor=$CLIENT_FONT_COLOR:borderw=$BORDER_WIDTH:bordercolor=$BORDER_COLOR:box=1:boxcolor=0x00000088:boxborderw=$BOX_PADDING:text_align=center:alpha='if(lt(t,$CLIENT_TEXT_FADE_OUT_START),1,if(lt(t,$IMAGE_TIME),(1-(t-$CLIENT_TEXT_FADE_OUT_START)/$TRANSITION_DURATION),0))'[v_cl];"
                    CUR_V="v_cl"
                fi
            fi
            FILTER_COMPLEX+="[$CUR_V]null[v_with_text];"
            
            # Identify the audio source and handle optional music/video audio
            CURRENT_STREAM="v_with_text"
            if [ -n "$MUSIC" ]; then
                # Music covers the full duration
                FILTER_COMPLEX+="[3:a]afade=t=in:st=0:d=$TRANSITION_DURATION,atrim=0:$TOTAL_DURATION,asetpts=PTS-STARTPTS,afade=t=out:st=$((TOTAL_DURATION - TRANSITION_DURATION)):d=$TRANSITION_DURATION[a_out];"
                NEXT_INPUT=4 # Icons start after [3:a]
            else
                # No music, generate silence to prevent ffmpeg crash
                FILTER_COMPLEX+="anullsrc=channel_layout=stereo:sample_rate=44100,atrim=0:$TOTAL_DURATION,asetpts=PTS-STARTPTS[a_out];"
                NEXT_INPUT=3 # Icons start after [2:v] (logo)
            fi
            
            # Add left icon overlay if provided
            if [ -n "$ICON_LEFT" ]; then
                ICON_HEIGHT="$FONT_SIZE"  # Scale icon to match font size
                ICON_SPACING=20  # Spacing between icon and text in pixels
                
                # Calculate icon X position: center minus half text width minus icon width minus spacing
                # Use max(0, ...) to ensure icon stays on screen even if text is wide
                ICON_LEFT_X="max(0,W/2-$ESTIMATED_TEXT_WIDTH/2-w-$ICON_SPACING)"
                
                # Scale icon maintaining aspect ratio, then apply ONLY fade out to match text timing
                FILTER_COMPLEX+="[${NEXT_INPUT}:v]scale=-1:${ICON_HEIGHT}:force_original_aspect_ratio=decrease,fade=t=out:st=$CLIENT_FADE_OUT_START:d=$TRANSITION_DURATION:alpha=1[icon_left_scaled];"
                # Position left icon to the left of the text
                # Y: Match CLIENT_TEXT_Y formula: (H-h)/1.25 positions at ~75% down
                FILTER_COMPLEX+="[${CURRENT_STREAM}][icon_left_scaled]overlay=x='$ICON_LEFT_X':y='(H-h)/1.25':enable='between(t,0,$IMAGE_TIME)':format=auto[v_with_left];"
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
                FILTER_COMPLEX+="[${NEXT_INPUT}:v]scale=-1:${ICON_HEIGHT}:force_original_aspect_ratio=decrease,fade=t=out:st=$CLIENT_FADE_OUT_START:d=$TRANSITION_DURATION:alpha=1[icon_right_scaled];"
                # Position right icon to the right of the text
                # Y: Match CLIENT_TEXT_Y formula: (H-h)/1.25 positions at ~75% down
                FILTER_COMPLEX+="[${CURRENT_STREAM}][icon_right_scaled]overlay=x='$ICON_RIGHT_X':y='(H-h)/1.25':enable='between(t,0,$IMAGE_TIME)':format=auto[v_with_right];"
                CURRENT_STREAM="v_with_right"
            fi
            
            # Set final output stream
            FILTER_COMPLEX+="[${CURRENT_STREAM}]null[v_out]"
            
            # Build ffmpeg command with optional icon inputs
            FFMPEG_CMD=(ffmpeg -v warning)
            FFMPEG_CMD+=(-loop 1 -t "$IMAGE_TIME" -i "$CLIENT_IMAGE")
            FFMPEG_CMD+=(-i "$input_video_path")
            FFMPEG_CMD+=(-loop 1 -t "$IMAGE_TIME" -i "$LOGO")
            
            if [ -n "$MUSIC" ]; then
                FFMPEG_CMD+=(-i "$MUSIC")
            fi
            
            # Add icon inputs if provided (loop them like client image)
            if [ -n "$ICON_LEFT" ]; then
                FFMPEG_CMD+=(-loop 1 -t "$IMAGE_TIME" -i "$ICON_LEFT")
            fi
            if [ -n "$ICON_RIGHT" ]; then
                FFMPEG_CMD+=(-loop 1 -t "$IMAGE_TIME" -i "$ICON_RIGHT")
            fi
            
            FFMPEG_CMD+=(-filter_complex "$FILTER_COMPLEX")
            FFMPEG_CMD+=(-map "[v_out]")
            FFMPEG_CMD+=(-map "[a_out]")
            FFMPEG_CMD+=(-t "$TOTAL_DURATION")
            FFMPEG_CMD+=(-c:v libx264)
            FFMPEG_CMD+=(-pix_fmt yuv420p)
            FFMPEG_CMD+=(-profile:v high)
            FFMPEG_CMD+=(-preset medium)
            FFMPEG_CMD+=(-crf 23)
            FFMPEG_CMD+=(-c:a aac)
            FFMPEG_CMD+=(-b:a 192k)
            FFMPEG_CMD+=(-movflags +faststart)
            FFMPEG_CMD+=(-f mp4)
            FFMPEG_CMD+=("$output_file_path")
            
            # Execute the command
            "${FFMPEG_CMD[@]}"
            FFMPEG_EXIT_CODE=$?
        else
            # Case 2: Main Video + Logo Video (no client)
            PRE_LOGO_VISUAL_DUR=$VIDEO_DURATION  # Duration before logo
            TOTAL_DURATION=$((PRE_LOGO_VISUAL_DUR + IMAGE_TIME - TRANSITION_DURATION))
            
            LOGO_TRANSITION_START=$((PRE_LOGO_VISUAL_DUR - TRANSITION_DURATION))
            VIDEO_END_FADE_START=$((PRE_LOGO_VISUAL_DUR - TRANSITION_DURATION))
            
            echo "Video: ${VIDEO_DURATION}s, Logo: ${IMAGE_TIME}s, Total: ${TOTAL_DURATION}s"
            
            # Scale and prepare video inputs
            FILTER_COMPLEX="[0:v]scale=$VIDEO_WIDTH:$VIDEO_HEIGHT,fps=30,setpts=PTS-STARTPTS[video_raw];"
            FILTER_COMPLEX+="[1:v]fps=30,setpts=PTS-STARTPTS[logo];"
            
            # Transition from input_video to logo
            # shellcheck disable=SC1087
            FILTER_COMPLEX+="[video_raw][logo]xfade=transition=fade:duration=$TRANSITION_DURATION:offset=$LOGO_TRANSITION_START[video_combined];"
            
            # --- Dynamic Text Overlay Logic (Case 2) ---
            CUR_V="video_combined"
            if [ -n "$EVENT_TEXT" ]; then
                ESC_EVENT_TEXT=$(echo -n "$EVENT_TEXT" | sed 's/\\/\\\\/g; s/[:'\'']/\\&/g')
                FILTER_COMPLEX+="[$CUR_V]drawtext=text='$ESC_EVENT_TEXT':x=$TEXT_X:y=$TEXT_Y:fontfile=$FONT:fontsize=$FONT_SIZE:fontcolor=$EVENT_FONT_COLOR:borderw=$BORDER_WIDTH:bordercolor=$BORDER_COLOR:box=1:boxcolor=0x00000088:boxborderw=$BOX_PADDING:line_spacing=$LINE_SPACING:text_align=center:alpha='if(gt(t,$VIDEO_END_FADE_START),(1-(t-$VIDEO_END_FADE_START)/$TRANSITION_DURATION),1)'[v_out];"
            else
                FILTER_COMPLEX+="[$CUR_V]null[v_out];"
            fi
            
            # Audio: Handle optional music and silent videos
            if [ -n "$MUSIC" ]; then
                FILTER_COMPLEX+="[2:a]afade=t=in:st=0:d=$TRANSITION_DURATION,atrim=0:$TOTAL_DURATION,asetpts=PTS-STARTPTS,afade=t=out:st=$((TOTAL_DURATION - TRANSITION_DURATION)):d=$TRANSITION_DURATION[a_out]"
            else
                # No music, generate silence
                FILTER_COMPLEX+="anullsrc=channel_layout=stereo:sample_rate=44100,atrim=0:$TOTAL_DURATION,asetpts=PTS-STARTPTS[a_out]"
            fi

            FFMPEG_CMD=(ffmpeg -v warning)
            FFMPEG_CMD+=(-i "$input_video_path")
            FFMPEG_CMD+=(-loop 1 -t "$IMAGE_TIME" -i "$LOGO")
            if [ -n "$MUSIC" ]; then
                FFMPEG_CMD+=(-i "$MUSIC")
            fi
            FFMPEG_CMD+=(-filter_complex "$FILTER_COMPLEX")
            FFMPEG_CMD+=(-map "[v_out]")
            FFMPEG_CMD+=(-map "[a_out]")
            FFMPEG_CMD+=(-t "$TOTAL_DURATION")
            FFMPEG_CMD+=(-c:v libx264)
            FFMPEG_CMD+=(-pix_fmt yuv420p)
            FFMPEG_CMD+=(-profile:v high)
            FFMPEG_CMD+=(-preset medium)
            FFMPEG_CMD+=(-crf 23)
            FFMPEG_CMD+=(-c:a aac)
            FFMPEG_CMD+=(-b:a 192k)
            FFMPEG_CMD+=(-movflags +faststart)
            FFMPEG_CMD+=(-f mp4)
            FFMPEG_CMD+=("$output_file_path")

            "${FFMPEG_CMD[@]}"
            FFMPEG_EXIT_CODE=$?
        fi

		# I think this is better because of how large the command is
        if [ $FFMPEG_EXIT_CODE -eq 0 ]; then
			((processed_count++))
            log_success "Successfully processed '$original_filename' to '$output_filename'"
        else
            log_error "Error processing: $original_filename (Check FFmpeg output above)"
        fi
        echo "------------------------------------------------------------------"
    fi
done

if [ "$processed_count" -eq 0 ]; then
    log_warning "No .mp4 files were found or processed in '$CUTTED_DIR'."
    log_info "Ensure you've made the Shotcut manual cuts and saved videos to that directory."
else
    echo "--- All automated conversions complete for event '$EVENT_NAME'! Total files processed: $processed_count ---"
    echo "Final videos are in: $EDITED_DIR"
fi
