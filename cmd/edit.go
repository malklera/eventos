package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
)

var (
	logo           string
	music          string
	image          string
	eventText      string
	clientText     string
	clientTextUp   string
	clientTextDown string
	clientColor    string
	textColor      string
	font           string
	iconLeft       string
	iconRight      string
	save           bool
	run            bool
)

var (
	ErrClientWithoutImage       = errors.New("'client' flag is only allowed when 'image' is passed")
	ErrUpWithoutImage           = errors.New("'up' flag is only allowed when 'image' is passed")
	ErrLeftWithoutClient        = errors.New("'left' flag is only allowed when 'client' is passed")
	ErrRightWithoutClient       = errors.New("'right' flag is only allowed when 'client' is passed")
	ErrClientColorWithoutClient = errors.New("'client-color' flag is only allowed when 'client' is passed")
	ErrTextColorWithoutText     = errors.New("'text-color' flag is only allowed when 'text' is passed")
	ErrFontWithoutTextOrClient  = errors.New("'font' flag is only allowed when 'text' or 'client' is passed")
)

var colorRe = regexp.MustCompile(`^0x[0-9A-Fa-f]{6}$`)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Command to edit the videos",
	Args:  cobra.ExactArgs(0),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return preRun(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validArgs(cmd); err != nil {
			return err
		}
		videoW := 1080
		videoH := 1920
		transitionT := 1.0
		imageT := 2.0

		// Text styling
		lineSpacing := 10
		fontSize := 80
		fontColor := "0xE6E70F"
		borderW := 2
		borderColor := "0x000000"
		boxPadding := 8
		clientTextX := "(w-tw)/2"
		clientTextY := "(h-th)/1.25"
		eventTextX := "(w-tw)/2"
		eventTextY := "H-th-100"
		eventTextMargin := 100

		cuttedDir := "cortado"

		// If eventText is not empty, wrap it as needed
		if cmd.Flags().Changed("text") {
			eventText = wrapText(eventText, fontSize, videoW, eventTextMargin)
		}
		// If clientText is not empty, wrap it as needed
		if cmd.Flags().Changed("client") {
			clientText = wrapText(clientText, fontSize, videoW, eventTextMargin)
		} else {
			// If it is empty, add the up and down text
			clientText = clientTextUp + "\n" + clientTextDown
		}

		textW := maxLenght(clientText, fontSize)

		cuttedVideos, err := listVideos(cuttedDir)
		if err != nil {
			return fmt.Errorf("listVideos(%s): %w", cuttedDir, err)
		}

		editedCount := 0
		start := time.Now()
		dst := "editado"
		musicT := 0.0
		if cmd.Flags().Changed("music") {
			var err error
			musicT, err = fileDuration(music)
			if err != nil {
				return fmt.Errorf("fileDuration(%s): %w", music, err)
			}
		}

		for _, file := range cuttedVideos {
			videoPath := filepath.Join(cuttedDir, file.Name())
			videoT, err := fileDuration(videoPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error getting the duration of '%s': %v", videoPath, err)
				continue
			}

			// CMD part
			dstF := filepath.Join(dst, file.Name())
			fmt.Println("Formateando:", videoPath)
			ffmpegEdit := []string{"ffmpeg", "-v", "error"}
			ffmpegEdit = append(ffmpegEdit, "-loop", "1", "-t", strconv.FormatFloat(imageT, 'f', -1, 64), "-i", image)
			ffmpegEdit = append(ffmpegEdit, "-i", videoPath)
			ffmpegEdit = append(ffmpegEdit, "-loop", "1", "-t", strconv.FormatFloat(imageT, 'f', -1, 64), "-i", logo)

			if cmd.Flags().Changed("music") {
				if musicT < (videoT + 2*imageT) {
					fmt.Fprintf(os.Stderr, "Error: music too short, musicT: %f, totalT: %f\n", musicT, (videoT + 2*imageT))
					continue
				}
				ffmpegEdit = append(ffmpegEdit, "-i", music)
			}

			if cmd.Flags().Changed("left") {
				ffmpegEdit = append(ffmpegEdit, "-loop", "1", "-t", strconv.FormatFloat(imageT, 'f', -1, 64), "-i", iconLeft)
			}
			if cmd.Flags().Changed("right") {
				ffmpegEdit = append(ffmpegEdit, "-loop", "1", "-t", strconv.FormatFloat(imageT, 'f', -1, 64), "-i", iconRight)
			}

			totalT := 0.0

			// client image + crossfade + main video + fade out + logo
			if cmd.Flags().Changed("image") {
				preLogoVisualT := imageT + videoT
				totalT = preLogoVisualT + imageT
				clientFadeOutStart := imageT - transitionT
				videoFadeOutStart := videoT - transitionT
				videoEndFadeStart := preLogoVisualT - transitionT
				pLen := (videoT + 2*transitionT) / 5

				part1End := pLen

				// xfade 1 setup: part1 and part2 overlap by transitionT
				part2Start := part1End - transitionT
				part2End := part2Start + pLen

				part3Start := part2End
				part3End := part3Start + pLen

				part4Start := part3End
				part4End := part4Start + pLen

				part5Start := part4End - transitionT

				textFadeInStart := imageT - transitionT
				textFadeInEnd := imageT
				textFadeOutStart := videoEndFadeStart

				fmt.Printf("Client: %f, Video: %f, Logo: %f, Total: %f\n", imageT, videoT, imageT, (imageT + videoT + imageT))

				// Scale and prepare all video inputs
				// Client image fade out to black at the end
				filterComplex := fmt.Sprintf("[0:v]scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2,fps=30,setpts=PTS-STARTPTS,fade=t=out:st=%f:d=%f[client];", videoW, videoH, videoW, videoH, clientFadeOutStart, transitionT)
				filterComplex += fmt.Sprintf("[1:v]scale=%d:%d,fps=30,setpts=PTS-STARTPTS[video];", videoW, videoH)
				// TODO: is this needed? logo is an image i control, with the correct dimensions
				// use this??
				// filterComplex += fmt.Sprintf("[2:v]fps=30,setpts=PTS-STARTPTS,fade=t=in:st=0:d=%d[logo];", transitionT)
				// logo: scale + fade-in from black
				filterComplex += fmt.Sprintf("[2:v]scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2,fps=30,setpts=PTS-STARTPTS,settb=AVTB,fade=t=in:st=0:d=%f[logo];", videoW, videoH, videoW, videoH, transitionT)
				// Slide-up entrance + slide-out exit for main video (both over TRANSITION_DURATION)
				// Slide-up: y goes H -> 0 in first TD seconds
				// Slide-out: y goes 0 -> -H in last TD seconds (slides up off-screen)
				filterComplex += fmt.Sprintf("color=c=black:s=%dx%d:d=%f:r=30[video_bg];", videoW, videoH, videoT)

				filterComplex += fmt.Sprintf("[video_bg][video]overlay=x=0:y='if(lt(t,%f),H-H*t/%f,if(gt(t,%f),-H*(t-%f)/%f,0))':format=auto[video_slide_raw];", transitionT, transitionT, videoFadeOutStart, videoFadeOutStart, transitionT)

				// split video int part1, part2-4 and part5 for xfades
				filterComplex += "[video_slide_raw]split=3[raw1][raw2to4][raw5];"
				// Part 1: From the beginning (including slide-up) until PART1_END
				filterComplex += fmt.Sprintf("[raw1]trim=start=0:end=%f,setpts=PTS-STARTPTS[part1];", part1End)
				// Parts 2-4: Starting from PART2_START until PART4_END
				filterComplex += fmt.Sprintf("[raw2to4]trim=start=%f:end=%f,setpts=PTS-STARTPTS[part_rest_raw];", part2Start, part4End)
				// Part 5: Starting from PART5_START
				filterComplex += fmt.Sprintf("[raw5]trim=start=%f,setpts=PTS-STARTPTS[part5];", part5Start)

				// --- Part 2 Effect: Quad-split (2x2 grid) ---
				filterComplex += "[part_rest_raw]split=2[pr_main][pr_quad_src];"
				filterComplex += fmt.Sprintf("[pr_quad_src]scale=%d:%d,split=4[q1][q2][q3][q4];", (videoW / 2), (videoH / 2))
				filterComplex += "[q1][q2]hstack[top_row];"
				filterComplex += "[q3][q4]hstack[bot_row];"
				filterComplex += "[top_row][bot_row]vstack[quad];"
				// Overlay quad on normal video only during Part 2's timeframe (0 to P_LEN)
				filterComplex += fmt.Sprintf("[pr_main][quad]overlay=enable='between(t,0,%f)':eof_action=pass[part_rest_step1];", pLen)
				// --- Transition: Frei0r Distort0r between Part 2 and Part 3 ---
				filterComplex += "[part_rest_step1]split=2[pr_step15_main][distort_src];"
				filterComplex += "[distort_src]frei0r=filter_name=distort0r:filter_params=0.5|0.005|y|0.25[distorted];"
				filterComplex += fmt.Sprintf("[pr_step15_main][distorted]overlay=enable='between(t,%f,%f)':eof_action=pass[part_rest_step15];", (pLen - transitionT/2), (pLen + transitionT/2))
				// --- Part 3 Effect: Slow/Fast time warp ---
				// Part 3 plays in slow-motion for the first half, speeding up for the second half
				p3LocalStart := part3Start - part2Start
				p3LocalEnd := part3End - part2Start
				filterComplex += "[part_rest_step15]split=2[pr_step2_main][p3_src];"
				filterComplex += fmt.Sprintf("[p3_src]trim=start=%f:end=%f,setpts=PTS-STARTPTS[p3_trim];", p3LocalStart, p3LocalEnd)
				// T is current time in seconds. First 1/4 of source duration plays at 0.5x speed (takes 1/2 of output P_LEN).
				// Remaining 3/4 plays at 1.5x speed (takes remaining 1/2 of output P_LEN).
				// PTS mapping securely warps timestamps precisely to fit exactly within P_LEN without breaking overlaps!
				filterComplex += fmt.Sprintf("[p3_trim]setpts='if(lt(T, %f/4), 2*PTS + %f/TB, (2/3)*PTS + (%f/3 + %f)/TB)'[p3_warp];", pLen, p3LocalStart, pLen, p3LocalStart)
				// Overlay exactly during Part 3's timeframe
				filterComplex += fmt.Sprintf("[pr_step2_main][p3_warp]overlay=enable='between(t,%f,%f)':eof_action=pass[part_rest_step2];", p3LocalStart, p3LocalEnd)
				// --- Transition: Spin effect between Part 3 and Part 4 ---
				// Creates a fast 360-degree dizzy spin precisely across the explicit frame cut boundary
				filterComplex += fmt.Sprintf("[part_rest_step2]rotate=a='if(between(t,%f,%f), 2*PI*(t-%f)/%f, 0)':ow=iw:oh=ih:c=black[part_rest_step25];", (pLen*2 - transitionT/2), (pLen*2 + transitionT/2), (pLen*2 - transitionT/2), transitionT)
				// --- Part 4 Effect: 3x3 grid ---
				// Part 4 timestamps need to be shifted because part_rest starts at PART2_START
				p4LocalStart := part4Start - part2Start
				p4LocalEnd := part4End - part2Start
				filterComplex += "[part_rest_step25]split=2[pr_step3_main][pr_grid3_src];"
				filterComplex += fmt.Sprintf("[pr_grid3_src]scale=%d:%d,split=9[g1][g2][g3][g4][g5][g6][g7][g8][g9];", (videoW / 3), (videoH / 3))
				filterComplex += "[g1][g2][g3]hstack=inputs=3[row1];"
				filterComplex += "[g4][g5][g6]hstack=inputs=3[row2];"
				filterComplex += "[g7][g8][g9]hstack=inputs=3[row3];"
				filterComplex += "[row1][row2][row3]vstack=inputs=3[grid3];"
				// Overlay 3x3 grid only during Part 4's explicit timeframe
				filterComplex += fmt.Sprintf("[pr_step3_main][grid3]overlay=enable='between(t,%f,%f)':eof_action=pass[part_rest_step4];", p4LocalStart, p4LocalEnd)
				// --- Part 4 Reverse Boomerang Effect ---
				// Cut Part 4 in half: play the first half normally, and the first half reversed during the second half
				p4LocalMid := p4LocalStart + (p4LocalEnd-p4LocalStart)/2
				filterComplex += "[part_rest_step4]split=2[pr_step4_main][p4_to_rev];"
				// Trim the exact floating-point first half, reverse it, shift PTS to map precisely onto the second half
				filterComplex += fmt.Sprintf("[p4_to_rev]trim=start=%f:end=%f,setpts=PTS-STARTPTS,reverse,setpts=PTS+(%f/TB)[p4_rev];", p4LocalStart, p4LocalMid, p4LocalMid)
				// Overlay reversed snippet safely during the second half, strictly enduring the xfade
				filterComplex += fmt.Sprintf("[pr_step4_main][p4_rev]overlay=enable='between(t,%f,%f)':eof_action=pass[part_rest];", p4LocalMid, p4LocalEnd)
				// Xfade 1: Part 1 and Parts 2-4
				xFadeOffset1 := part1End - transitionT
				filterComplex += fmt.Sprintf("[part1][part_rest]xfade=transition=fade:duration=%f:offset=%f[video_step1];", transitionT, xFadeOffset1)
				// Xfade 2: video_step1 (which ends exactly at PART4_END) and Part 5
				xFadeOffset2 := part4End - transitionT
				filterComplex += fmt.Sprintf("[video_step1][part5]xfade=transition=fade:duration=%f:offset=%f[video_slide];", transitionT, xFadeOffset2)
				// Concatenate: client (with fade-out) -> main video (with slide-up + quad effect)
				// settb=AVTB normalizes timebase so xfade inputs match
				filterComplex += "[client][video_slide]concat=n=2:v=1:a=0,settb=AVTB[client_video_raw];"
				// Concatenate: client_video_raw -> logo (with fade-in)
				filterComplex += "[client_video_raw][logo]concat=n=2:v=1:a=0[video_combined];"

				// --- Dynamic Text Overlay Logic (Case 1) ---
				curV := "video_combined"

				// 1. Event text
				if cmd.Flags().Changed("text") {
					filterComplex += fmt.Sprintf("[%s]drawtext=text='%s':x=%s:y=%s:fontfile=%s:fontsize=%d:fontcolor=%s:borderw=%d:bordercolor=%s:box=1:boxcolor=0x00000088:boxborderw=%d:line_spacing=%d:text_align=center:alpha='if(lt(t,%f),0,if(lt(t,%f),(t-%f)/%f,if(gt(t,%f),(1-(t-%f)/%f),1)))'[v_ev];", curV, eventText, eventTextX, eventTextY, font, fontSize, fontColor, borderW, borderColor, boxPadding, lineSpacing, textFadeInStart, textFadeInEnd, textFadeInStart, transitionT, textFadeOutStart, textFadeOutStart, transitionT)
					curV = "v_ev"
				}

				// 2. CLIENT_TEXT (during introduction image)
				// TODO: this conditional is wrong, clientText and clientTextUp/Down
				// are mutually exlclusive or should be
				if cmd.Flags().Changed("client") {
					if cmd.Flags().Changed("up") || cmd.Flags().Changed("down") {
						if cmd.Flags().Changed("up") {
							filterComplex += fmt.Sprintf("[%s]drawtext=text='%s':x=%s:y=%s-60:fontfile=%s:fontsize=%d:fontcolor=%s:borderw=%d:bordercolor=%s:box=1:boxcolor=0x00000088:boxborderw=%d:alpha='if(lt(t,%f),1,if(lt(t,%f),(1-(t-%f)/%f),0))'[v_up];", curV, clientTextUp, clientTextX, clientTextY, font, fontSize, fontColor, borderW, borderColor, boxPadding, clientFadeOutStart, imageT, clientFadeOutStart, transitionT)
							curV = "v_up"
						}
						if cmd.Flags().Changed("down") {
							filterComplex += fmt.Sprintf("[%s]drawtext=text='%s':x=%s:y=%s+60:fontfile=%s:fontsize=%d:fontcolor=%s:borderw=%d:bordercolor=%s:box=1:boxcolor=0x00000088:boxborderw=%d:alpha='if(lt(t,%f),1,if(lt(t,%f),(1-(t-%f)/%f),0))'[v_dw];", curV, clientTextUp, clientTextX, clientTextY, font, fontSize, fontColor, borderW, borderColor, boxPadding, clientFadeOutStart, imageT, clientFadeOutStart, transitionT)
							curV = "v_dw"
						}
					} else {
						// Fallback to single text (auto-wrapped if needed)
						filterComplex += fmt.Sprintf("[%s]drawtext=text='%s':x=%s:y=%s:fontfile=%s:fontsize=%d:fontcolor=%s:borderw=%d:bordercolor=%s:box=1:boxcolor=0x00000088:boxborderw=%d:text_align=center:alpha='if(lt(t,%f),1,if(lt(t,%f),(1-(t-%f)/%f),0))'[v_cl];", curV, clientText, clientTextX, clientTextY, font, fontSize, fontColor, borderW, borderColor, boxPadding, clientFadeOutStart, imageT, clientFadeOutStart, transitionT)
						curV = "v_cl"
					}
				}
				filterComplex += fmt.Sprintf("[%s]null[v_with_text];", curV)

				// Identify the audio source and handle optional music/video audio
				currentStream := "v_with_text"
				nextInput := 0
				if cmd.Flags().Changed("music") {
					// Music covers the full duration
					filterComplex += fmt.Sprintf("[3:a]afade=t=in:st=0:d=%f,atrim=0:%f,asetpts=PTS-STARTPTS,afade=t=out:st=%f:d=%f[a_out];", transitionT, totalT, (totalT - transitionT), transitionT)
					// Icons start after [3:a]
					nextInput = 4
				} else {
					// No music, generate silence to prevent ffmpeg crash
					filterComplex += fmt.Sprintf("anullsrc=channel_layout=stereo:sample_rate=44100,atrim=0:%f,asetpts=PTS-STARTPTS[a_out];", totalT)
					// Icons start after [2:v] (logo)
					nextInput = 3
				}

				iconSpacing := 20
				if cmd.Flags().Changed("left") {
					// Calculate icon X position: center minus half text width minus icon width minus spacing
					// Use max(0, ...) to ensure icon stays on screen even if text is wide
					iconLeftX := fmt.Sprintf("max(0,W/2-%d/2-w-%d)", textW, iconSpacing)

					// Scale icon maintaining aspect ratio, then apply ONLY fade out to match text timing
					filterComplex += fmt.Sprintf("[%d:v]scale=-1:%d:force_original_aspect_ratio=decrease,fade=t=out:st=%f:d=%f:alpha=1[icon_left_scaled];", nextInput, fontSize, clientFadeOutStart, transitionT)
					// Position left icon to the left of the text
					// Y: Match CLIENT_TEXT_Y formula: (H-h)/1.25 positions at ~75% down
					filterComplex += fmt.Sprintf("[%s][icon_left_scaled]overlay=x='%s':y='(H-h)/1.25':enable='between(t,0,%f)':format=auto[v_with_left];", currentStream, iconLeftX, imageT)
					currentStream = "v_with_left"
					nextInput++
				}

				if cmd.Flags().Changed("right") {
					// Calculate icon X position: center plus half text width plus spacing
					// Use min(W-w, ...) to ensure icon stays on screen
					iconRightX := fmt.Sprintf("min(W-w,W/2+%d/2+%d)", textW, iconSpacing)

					// Scale icon maintaining aspect ratio, then apply ONLY fade out to match text timing
					filterComplex += fmt.Sprintf("[%d:v]scale=-1:%d:force_original_aspect_ratio=decrease,fade=t=out:st=%f:d=%f:alpha=1[icon_right_scaled];", nextInput, fontSize, clientFadeOutStart, transitionT)
					// Position right icon to the right of the text
					// Y: Match CLIENT_TEXT_Y formula: (H-h)/1.25 positions at ~75% down
					filterComplex += fmt.Sprintf("[%s][icon_right_scaled]overlay=x='%s':y='(H-h)/1.25':enable='between(t,0,%f)':format=auto[v_with_right];", currentStream, iconRightX, imageT)
					currentStream = "v_with_right"
				}

				filterComplex += fmt.Sprintf("[%s]null[v_out]", currentStream)

				ffmpegEdit = append(ffmpegEdit, "-filter_complex", filterComplex)

			} else {
				// Case 2: main video + fade out + logo video (no client)
				totalT = videoT + transitionT + imageT
				logoTransitionStart := videoT + transitionT
				videoEndFadeStart := videoT + transitionT

				fmt.Printf("Client: %f, Video: %f, Logo: %f, Total: %f\n", imageT, videoT, imageT, (imageT + videoT + imageT))

				// Scale and prepare video inputs
				filterComplex := fmt.Sprintf("[0:v]scale=%d:%d,fps=30,setpts=PTS-STARTPTS[video_raw];", videoW, videoH)
				filterComplex += "[1:v]fps=30,setpts=PTS-STARTPTS[logo];"

				// Transition from input_video to logo
				filterComplex += fmt.Sprintf("[video_raw][logo]xfade=transition=fade:duration=%f:offset=%f[video_combined];", transitionT, logoTransitionStart)

				// --- Dynamic Text Overlay Logic (Case 2) ---
				curV := "video_combined"
				if cmd.Flags().Changed("text") {
					filterComplex += fmt.Sprintf("[%s]drawtext=text='%s':x=%s:y=%s:fontfile=%s:fontsize=%d:fontcolor=%s:borderw=%d:bordercolor=%s:box=1:boxcolor=0x00000088:boxborderw=%d:line_spacing=%d:text_align=center:alpha='if(gt(t,%f),(1-(t-%f)/%f),1)'[v_out];", curV, eventText, eventTextX, eventTextY, font, fontSize, fontColor, borderW, borderColor, boxPadding, lineSpacing, videoEndFadeStart, videoEndFadeStart, transitionT)
				} else {
					filterComplex += fmt.Sprintf("[%s]null[v_out];", curV)
				}

				// Audio: Handle optional music and silent videos
				if cmd.Flags().Changed("music") {
					// Music covers the full duration
					filterComplex += fmt.Sprintf("[2:a]afade=t=in:st=0:d=%f,atrim=0:%f,asetpts=PTS-STARTPTS,afade=t=out:st=%f:d=%f[a_out];", transitionT, totalT, (totalT - transitionT), transitionT)
				} else {
					// No music, generate silence to prevent ffmpeg crash
					filterComplex += fmt.Sprintf("anullsrc=channel_layout=stereo:sample_rate=44100,atrim=0:%f,asetpts=PTS-STARTPTS[a_out];", totalT)
				}

				ffmpegEdit = append(ffmpegEdit, "-filter_complex", filterComplex)
			}

			ffmpegEdit = append(ffmpegEdit, "-map", "[v_out]")
			ffmpegEdit = append(ffmpegEdit, "-map", "[a_out]")
			ffmpegEdit = append(ffmpegEdit, "-t", strconv.FormatFloat(totalT, 'f', -1, 64))
			ffmpegEdit = append(ffmpegEdit, "-c:v", "libx264")
			ffmpegEdit = append(ffmpegEdit, "-pix_fmt", "yuv420p")
			ffmpegEdit = append(ffmpegEdit, "-profile:v", "high")
			ffmpegEdit = append(ffmpegEdit, "-preset", "medium")
			ffmpegEdit = append(ffmpegEdit, "-crf", "23")
			ffmpegEdit = append(ffmpegEdit, "-c:a", "aac")
			ffmpegEdit = append(ffmpegEdit, "-b:a", "192k")
			ffmpegEdit = append(ffmpegEdit, "-movflags", "+faststart")
			ffmpegEdit = append(ffmpegEdit, "-f", "mp4")
			ffmpegEdit = append(ffmpegEdit, dstF)

			ffmpeg := exec.Command(ffmpegEdit[0], ffmpegEdit[1:]...)
			ffmpeg.Stdout = os.Stdout
			ffmpeg.Stderr = os.Stderr
			err = ffmpeg.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error running %s: %v\n", ffmpeg.String(), err)
			} else {
				editedCount++
			}
		}
		fmt.Println("Total videos:", len(cuttedVideos))
		fmt.Println("Edited videos:", editedCount)
		fmt.Println("Time taken:", time.Since(start))
		return nil
	},
}

func init() {
	editCmd.Flags().StringVarP(&logo, "logo", "l", "", "Path to logo image file (~/Videos/eventos/assets/logo/myLogo.png).")
	editCmd.Flags().StringVarP(&music, "music", "m", "", "Path to music file (~/Videos/eventos/assets/musica/cortado/myMusic.mp3).")
	editCmd.Flags().StringVarP(&image, "image", "i", "", "Client image file name, file has to be inside <event-name>/.")
	editCmd.Flags().StringVarP(&eventText, "text", "t", "", "Text to overlay on the whole video.")

	editCmd.Flags().StringVarP(&clientText, "client", "c", "", "Client text to display over client image.")
	editCmd.Flags().StringVarP(&clientTextUp, "up", "u", "", "Upper line of client text (manual split, bypasses auto-wrap).")
	editCmd.MarkFlagsMutuallyExclusive("client", "up")

	editCmd.Flags().StringVarP(&clientTextDown, "down", "d", "", "Lower line of client text (manual split, bypasses auto-wrap).")
	editCmd.Flags().StringVarP(&clientColor, "client-color", "C", "0xE6E70F", "Client text color in hex format (e.g., \"0xFFFFFF\" or \"0xE6E70F\").")
	editCmd.Flags().StringVarP(&textColor, "text-color", "T", "0xE6E70F", "Event text color in hex format (e.g., \"0xFFFFFF\" or \"0xE6E70F\").")
	editCmd.Flags().StringVarP(&font, "font", "f", "/usr/local/share/fonts/Courgette-Regular.ttf", "Path to font to use for all text (e.g., /usr/share/fonts/TTF/MyFont.ttf).")
	editCmd.Flags().StringVarP(&iconLeft, "left", "L", "", "Path to icon image to display to the left of client text.")
	editCmd.Flags().StringVarP(&iconRight, "right", "R", "", "Path to icon image to display to the right of client text.")

	editCmd.Flags().BoolVarP(&save, "save", "s", false, "Save the current command.")
	editCmd.Flags().BoolVarP(&run, "run", "r", false, "Run the saved command.")
	editCmd.MarkFlagsMutuallyExclusive("save", "run")
	// You either `run` the saved command or executed/`save` the current one, making `logo`
	// required.
	editCmd.MarkFlagsOneRequired("run", "logo")

	rootCmd.AddCommand(editCmd)
}

func preRun(cmd *cobra.Command) error {
	if !cmd.Flags().Changed("image") && cmd.Flags().Changed("client") {
		return ErrClientWithoutImage
	}
	if !cmd.Flags().Changed("image") && cmd.Flags().Changed("up") {
		return ErrUpWithoutImage
	}
	if !cmd.Flags().Changed("client") && cmd.Flags().Changed("left") {
		return ErrLeftWithoutClient
	}
	if !cmd.Flags().Changed("client") && cmd.Flags().Changed("right") {
		return ErrRightWithoutClient
	}
	if !cmd.Flags().Changed("client") && cmd.Flags().Changed("client-color") {
		return ErrClientColorWithoutClient
	}
	if !cmd.Flags().Changed("text") && cmd.Flags().Changed("text-color") {
		return ErrTextColorWithoutText
	}
	if (!cmd.Flags().Changed("text") && !cmd.Flags().Changed("client")) && cmd.Flags().Changed("font") {
		return ErrFontWithoutTextOrClient
	}
	return nil
}

func wrapText(text string, fontSize int, videoW int, margin int) string {
	// Use a 0.7 factor (aprox 56px) for safer wrap.
	return wordwrap.WrapString(text, uint((videoW-margin)/(fontSize*7/10)))
}

// fileDuration take a path and run ffprobe to learn its duration, return the quotient and any error
func fileDuration(path string) (float64, error) {
	ffprobe := exec.Command("ffprobe",
		"-loglevel",
		"error",
		"-show_entries",
		"format=duration",
		"-print_format",
		"default=noprint_wrappers=1:nokey=1",
		path)
	ffprobe.Stderr = os.Stderr
	out, err := ffprobe.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe.Run(%s): %w", ffprobe.String(), err)
	}
	dur, err := strconv.ParseFloat(string(out), 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.Atoi(string(%v)): %w", out, err)
	}
	return dur, nil
}

// validArgs ensure all passed arguments exists
func validArgs(cmd *cobra.Command) error {
	_, err := os.Stat(logo)
	if err != nil {
		return err
	}
	if cmd.Flags().Changed("image") {
		_, err = os.Stat(image)
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("music") {
		_, err = os.Stat(music)
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("font") {
		_, err = os.Stat(font)
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("left") {
		_, err = os.Stat(iconLeft)
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("right") {
		_, err = os.Stat(iconRight)
		if err != nil {
			return err
		}
	}

	if cmd.Flags().Changed("client-color") {
		if !colorRe.MatchString(clientColor) {
			return fmt.Errorf("'%s' is an invalid color, use e.g. '0xRRGGBB'", clientColor)
		}
	}

	if cmd.Flags().Changed("text-color") {
		if !colorRe.MatchString(textColor) {
			return fmt.Errorf("'%s' is an invalid color, use e.g. '0xRRGGBB'", textColor)
		}
	}

	return nil
}

// maxLenght takes a possibly newline separated string and measure the longest
// line in it
func maxLenght(t string, fontSize int) int {
	parts := strings.Split(t, "\n")
	lenght := 0
	for _, p := range parts {
		if lenght < len(p) {
			lenght = len(p)
		}
	}
	return lenght * fontSize * 7 / 10
}
