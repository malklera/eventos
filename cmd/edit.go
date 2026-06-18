package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	leftIcon       string
	rightIcon      string
	save           bool
	run            bool
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Command to edit the videos",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := checkArgs(); err != nil {
			return err
		}
		videoW := 1080
		videoH := 1920
		transitionT := 1
		imageT := 2

		// Text styling
		lineSpacing := 10
		fontSize := 80
		fontColor := "#E6E70F"
		borderW := 2
		borderColor := "#000000"
		boxPadding := 8
		clientTextX := "(w-tw)/2"
		clientTextY := "(h-th)/1.25"
		eventTextX := "(w-tw)/2"
		eventTextY := "H-th-100"
		eventTextMargin := 100

		cuttedDir := "cortado"
		// editedDir := "editado"

		// If eventText is not empty, wrap it as needed
		if eventText != "" {
			eventText = wrapText(eventText, fontSize, videoW, eventTextMargin)
		}
		// If clientText is not empty, wrap it as needed
		if clientText != "" {
			clientText = wrapText(clientText, fontSize, videoW, eventTextMargin)
		} else {
			// If it is empty, add the up and down text
			clientText = clientTextUp + "\n" + clientTextDown
		}

		// textW := maxLenght(clientText, fontSize)

		cuttedVideos, err := listVideos(cuttedDir)
		if err != nil {
			return fmt.Errorf("listVideos(%s): %w", cuttedDir, err)
		}

		editedCount := 0
		start := time.Now()
		dst := "editado"

		for _, file := range cuttedVideos {
			videoPath := filepath.Join(cuttedDir, file.Name())
			videoT, err := videoDuration(videoPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error getting the duration of '%s': %v", videoPath, err)
				continue
			}

			// client image + main video + logo
			if image != "" {
				preLogoVisualT := videoT + imageT
				totalT := preLogoVisualT + imageT
				clientFadeOutStart := imageT - transitionT
				videoFadeOutStart := videoT - transitionT
				videoEndFadeStart := preLogoVisualT - transitionT
				// segmentY := (videoT - transitionT) / 8
				pLen := (videoT + 2*transitionT) / 5

				// part1Start := 0
				part1End := pLen

				// xfade 1 setup: part1 and part2 overlap by transitionDuration
				part2Start := part1End - transitionT
				part2End := part2Start + pLen

				part3Start := part2End
				part3End := part3Start + pLen

				part4Start := part3End
				part4End := part4Start + pLen

				part5Start := part4End - transitionT
				// part5End := videoT

				textFadeInStart := imageT - transitionT
				textFadeInEnd := imageT
				textFadeOutStart := videoEndFadeStart
				fmt.Printf("Client: %d, Video: %d, Logo: %d, Total: %d\n", imageT, videoT, imageT, (imageT + videoT + imageT))

				// Scale and prepare all video inputs
				// Client image fade out to black at the end
				filterComplex := fmt.Sprintf("[0:v]scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2,fps=30,setpts=PTS-STARTPTS,fade=t=out:st=%d:d=%d[client];", videoW, videoH, videoW, videoH, clientFadeOutStart, transitionT)
				filterComplex += fmt.Sprintf("[1:v]scale=%d:%d,fps=30,setpts=PTS-STARTPTS[video];", videoW, videoH)
				// TODO: is this needed? logo is an image i control, with the correct dimensions
				// logo: scale + fade-in from black
				filterComplex += fmt.Sprintf("[2:v]scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2,fps=30,setpts=PTS-STARTPTS,settb=AVTB,fade=t=in:st=0:d=%d[logo];", videoW, videoH, videoW, videoH, transitionT)
				// Slide-up entrance + slide-out exit for main video (both over TRANSITION_DURATION)
				// Slide-up: y goes H -> 0 in first TD seconds
				// Slide-out: y goes 0 -> -H in last TD seconds (slides up off-screen)
				filterComplex += fmt.Sprintf("color=c=black:s=%dx%d:d=%d:r=30[video_bg];", videoW, videoH, videoT)

				filterComplex += fmt.Sprintf("[video_bg][video]overlay=x=0:y='if(lt(t,%d),H-H*t/%d,if(gt(t,%d),-H*(t-%d)/%d,0))':format=auto[video_slide_raw];", transitionT, transitionT, videoFadeOutStart, videoFadeOutStart, transitionT)

				// split video int part1, part2-4 and part5 for xfades
				filterComplex += "[video_slide_raw]split=3[raw1][raw2to4][raw5];"
				// Part 1: From the beginning (including slide-up) until PART1_END
				filterComplex += fmt.Sprintf("[raw1]trim=start=0:end=%d,setpts=PTS-STARTPTS[part1];", part1End)
				// Parts 2-4: Starting from PART2_START until PART4_END
				filterComplex += fmt.Sprintf("[raw2to4]trim=start=%d:end=%d,setpts=PTS-STARTPTS[part_rest_raw];", part2Start, part4End)
				// Part 5: Starting from PART5_START
				filterComplex += fmt.Sprintf("[raw5]trim=start=%d,setpts=PTS-STARTPTS[part5];", part5Start)

				// --- Part 2 Effect: Quad-split (2x2 grid) ---
				filterComplex += "[part_rest_raw]split=2[pr_main][pr_quad_src];"
				filterComplex += fmt.Sprintf("[pr_quad_src]scale=%d:%d,split=4[q1][q2][q3][q4];", (videoW / 2), (videoH / 2))
				filterComplex += "[q1][q2]hstack[top_row];"
				filterComplex += "[q3][q4]hstack[bot_row];"
				filterComplex += "[top_row][bot_row]vstack[quad];"
				// Overlay quad on normal video only during Part 2's timeframe (0 to P_LEN)
				filterComplex += fmt.Sprintf("[pr_main][quad]overlay=enable='between(t,0,%d)':eof_action=pass[part_rest_step1];", pLen)
				// --- Transition: Frei0r Distort0r between Part 2 and Part 3 ---
				filterComplex += "[part_rest_step1]split=2[pr_step15_main][distort_src];"
				filterComplex += "[distort_src]frei0r=filter_name=distort0r:filter_params=0.5|0.005|y|0.25[distorted];"
				filterComplex += fmt.Sprintf("[pr_step15_main][distorted]overlay=enable='between(t,%d,%d)':eof_action=pass[part_rest_step15];", (pLen - transitionT/2), (pLen + transitionT/2))
				// --- Part 3 Effect: Slow/Fast time warp ---
				// Part 3 plays in slow-motion for the first half, speeding up for the second half
				p3LocalStart := part3Start - part2Start
				p3LocalEnd := part3End - part2Start
				filterComplex += "[part_rest_step15]split=2[pr_step2_main][p3_src];"
				filterComplex += fmt.Sprintf("[p3_src]trim=start=%d:end=%d,setpts=PTS-STARTPTS[p3_trim];", p3LocalStart, p3LocalEnd)
				// T is current time in seconds. First 1/4 of source duration plays at 0.5x speed (takes 1/2 of output P_LEN).
				// Remaining 3/4 plays at 1.5x speed (takes remaining 1/2 of output P_LEN).
				// PTS mapping securely warps timestamps precisely to fit exactly within P_LEN without breaking overlaps!
				filterComplex += fmt.Sprintf("[p3_trim]setpts='if(lt(T, %d/4), 2*PTS + %d/TB, (2/3)*PTS + (%d/3 + %d)/TB)'[p3_warp];", pLen, p3LocalStart, pLen, p3LocalStart)
				// Overlay exactly during Part 3's timeframe
				filterComplex += fmt.Sprintf("[pr_step2_main][p3_warp]overlay=enable='between(t,%d,%d)':eof_action=pass[part_rest_step2];", p3LocalStart, p3LocalEnd)
				// --- Transition: Spin effect between Part 3 and Part 4 ---
				// Creates a fast 360-degree dizzy spin precisely across the explicit frame cut boundary
				filterComplex += fmt.Sprintf("[part_rest_step2]rotate=a='if(between(t,%d,%d), 2*PI*(t-%d)/%d, 0)':ow=iw:oh=ih:c=black[part_rest_step25];", (pLen*2 - transitionT/2), (pLen*2 + transitionT/2), (pLen*2 - transitionT/2), transitionT)
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
				filterComplex += fmt.Sprintf("[pr_step3_main][grid3]overlay=enable='between(t,%d,%d)':eof_action=pass[part_rest_step4];", p4LocalStart, p4LocalEnd)
				// --- Part 4 Reverse Boomerang Effect ---
				// Cut Part 4 in half: play the first half normally, and the first half reversed during the second half
				p4LocalMid := p4LocalStart + (p4LocalEnd-p4LocalStart)/2
				filterComplex += "[part_rest_step4]split=2[pr_step4_main][p4_to_rev];"
				// Trim the exact floating-point first half, reverse it, shift PTS to map precisely onto the second half
				filterComplex += fmt.Sprintf("[p4_to_rev]trim=start=%d:end=%d,setpts=PTS-STARTPTS,reverse,setpts=PTS+(%d/TB)[p4_rev];", p4LocalStart, p4LocalMid, p4LocalMid)
				// Overlay reversed snippet safely during the second half, strictly enduring the xfade
				filterComplex += fmt.Sprintf("[pr_step4_main][p4_rev]overlay=enable='between(t,%d,%d)':eof_action=pass[part_rest];", p4LocalMid, p4LocalEnd)
				// Xfade 1: Part 1 and Parts 2-4
				xFadeOffset1 := part1End - transitionT
				filterComplex += fmt.Sprintf("[part1][part_rest]xfade=transition=fade:duration=%d:offset=%d[video_step1];", transitionT, xFadeOffset1)
				// Xfade 2: video_step1 (which ends exactly at PART4_END) and Part 5
				xFadeOffset2 := part4End - transitionT
				filterComplex += fmt.Sprintf("[video_step1][part5]xfade=transition=fade:duration=%d:offset=%d[video_slide];", transitionT, xFadeOffset2)
				// Concatenate: client (with fade-out) -> main video (with slide-up + quad effect)
				// settb=AVTB normalizes timebase so xfade inputs match
				filterComplex += "[client][video_slide]concat=n=2:v=1:a=0,settb=AVTB[client_video_raw];"
				// Concatenate: client_video_raw -> logo (with fade-in)
				filterComplex += "[client_video_raw][logo]concat=n=2:v=1:a=0[video_combined];"

				// --- Dynamic Text Overlay Logic (Case 1) ---
				curV := "video_combined"

				// 1. Event text
				if eventText != "" {
					filterComplex += fmt.Sprintf("[%s]drawtext=text='%s':x=%s:y=%s:fontfile=%s:fontsize=%d:fontcolor=%s:borderw=%d:bordercolor=%s:box=1:boxcolor=0x00000088:boxborderw=%d:line_spacing=%d:text_align=center:alpha='if(lt(t,%d),0,if(lt(t,%d),(t-%d)/%d,if(gt(t,%d),(1-(t-%d)/%d),1)))'[v_ev];", curV, eventText, eventTextX, eventTextY, font, fontSize, fontColor, borderW, borderColor, boxPadding, lineSpacing, textFadeInStart, textFadeInEnd, textFadeInStart, transitionT, textFadeOutStart, textFadeOutStart, transitionT)
					curV = "v_ev"
				}

				// 2. CLIENT_TEXT (during introduction image)
				// TODO: this conditional is wrong, clientText and clientTextUp/Down
				// are mutually exlclusive or should be
				if clientText != "" {
					if clientTextUp != "" || clientTextDown != "" {
						if clientTextUp != "" {
							filterComplex += fmt.Sprintf("[%s]drawtext=text='%s':x=%s:y=%s-60:fontfile=%s:fontsize=%d:fontcolor=%s:borderw=%d:bordercolor=%s:box=1:boxcolor=0x00000088:boxborderw=%d:alpha='if(lt(t,%d),1,if(lt(t,%d),(1-(t-%d)/%d),0))'[v_up];", curV, clientTextUp, clientTextX, clientTextY, font, fontSize, fontColor, borderW, borderColor, boxPadding, clientFadeOutStart, imageT, clientFadeOutStart, transitionT)
							curV = "v_up"
						}
						if clientTextDown != "" {
							filterComplex += fmt.Sprintf("[%s]drawtext=text='%s':x=%s:y=%s+60:fontfile=%s:fontsize=%d:fontcolor=%s:borderw=%d:bordercolor=%s:box=1:boxcolor=0x00000088:boxborderw=%d:alpha='if(lt(t,%d),1,if(lt(t,%d),(1-(t-%d)/%d),0))'[v_dw];", curV, clientTextUp, clientTextX, clientTextY, font, fontSize, fontColor, borderW, borderColor, boxPadding, clientFadeOutStart, imageT, clientFadeOutStart, transitionT)
							curV = "v_dw"
						}
					} else {
						// Fallback to single text (auto-wrapped if needed)
						filterComplex += fmt.Sprintf("[%s]drawtext=text='%s':x=%s:y=%s:fontfile=%s:fontsize=%d:fontcolor=%s:borderw=%d:bordercolor=%s:box=1:boxcolor=0x00000088:boxborderw=%d:text_align=center:alpha='if(lt(t,%d),1,if(lt(t,%d),(1-(t-%d)/%d),0))'[v_cl];", curV, clientText, clientTextX, clientTextY, font, fontSize, fontColor, borderW, borderColor, boxPadding, clientFadeOutStart, imageT, clientFadeOutStart, transitionT)
						curV = "v_cl"
					}
				}
				filterComplex += fmt.Sprintf("[%s]null[v_with_text];", curV)

				// Identify the audio source and handle optional music/video audio
				currentStream := "v_with_text"
				// nextInput := 0
				if music != "" {
					// Music covers the full duration
					filterComplex += fmt.Sprintf("[3:a]afade=t=in:st=0:d=%d,atrim=0:%d,asetpts=PTS-STARTPTS,afade=t=out:st=%d:d=%d[a_out];", transitionT, totalT, (totalT - transitionT), transitionT)
					// Icons start after [3:a]
					// nextInput = 4
				} else {
					// No music, generate silence to prevent ffmpeg crash
					filterComplex += fmt.Sprintf("anullsrc=channel_layout=stereo:sample_rate=44100,atrim=0:%d,asetpts=PTS-STARTPTS[a_out];", totalT)
					// Icons start after [2:v] (logo)
					// nextInput = 3
				}

				// TODO: need to copy the icons

				filterComplex += fmt.Sprintf("[%s]null[v_out]", currentStream)

				// CMD part
				dstF := filepath.Join(dst, file.Name())
				fmt.Println("Formateando:", videoPath)
				ffmpegEdit := []string{"ffmpeg", "-v", "error"}
				// ffmpegEdit = append(ffmpegEdit, "ffmpeg", "-v", "warning")
				ffmpegEdit = append(ffmpegEdit, "-loop", "1", "-t", strconv.Itoa(imageT), "-i", image)
				ffmpegEdit = append(ffmpegEdit, "-i", videoPath)
				ffmpegEdit = append(ffmpegEdit, "-loop", "1", "-t", strconv.Itoa(imageT), "-i", logo)

				if music != "" {
					ffmpegEdit = append(ffmpegEdit, "-i", music)
				}

				if leftIcon != "" {
					ffmpegEdit = append(ffmpegEdit, "-loop", "1", "-t", strconv.Itoa(imageT), "-i", leftIcon)
				}
				if rightIcon != "" {
					ffmpegEdit = append(ffmpegEdit, "-loop", "1", "-t", strconv.Itoa(imageT), "-i", rightIcon)
				}

				ffmpegEdit = append(ffmpegEdit, "-filter_complex", filterComplex)
				ffmpegEdit = append(ffmpegEdit, "-map", "[v_out]")
				ffmpegEdit = append(ffmpegEdit, "-map", "[a_out]")
				ffmpegEdit = append(ffmpegEdit, "-t", strconv.Itoa(totalT))
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
				err := ffmpeg.Run()
				if err != nil {
					fmt.Fprintf(os.Stderr, "error running %s: %v\n", ffmpeg.String(), err)
				} else {
					editedCount++
				}
			} else {
				// TODO: copy this later
				continue
			}
		}
		fmt.Println("Total videos:", len(cuttedVideos))
		fmt.Println("Edited videos:", editedCount)
		fmt.Println("Time taken:", time.Since(start))
		return nil
	},
}

func init() {
	// Do not really need this i think, at least for now.
	// editCmd.Flags().String("event", "e", "", "Name of the event.")
	editCmd.Flags().StringVarP(&logo, "logo", "l", "", "Path to logo image file (~/Videos/eventos/assets/logo/myLogo.png).")
	editCmd.Flags().StringVarP(&music, "music", "m", "", "Path to music file (~/Videos/eventos/assets/musica/cortado/myMusic.mp3).")
	editCmd.Flags().StringVarP(&image, "image", "i", "", "Client image file name, file has to be inside <event-name>/.")
	editCmd.Flags().StringVarP(&eventText, "text", "t", "", "Text to overlay on the whole video.")

	editCmd.Flags().StringVarP(&clientText, "client", "c", "", "Client text to display over client image.")
	editCmd.Flags().StringVarP(&clientTextUp, "up", "u", "", "Upper line of client text (manual split, bypasses auto-wrap).")
	editCmd.MarkFlagsMutuallyExclusive("client", "up")

	editCmd.Flags().StringVarP(&clientTextDown, "down", "d", "", "Lower line of client text (manual split, bypasses auto-wrap).")
	editCmd.Flags().StringVarP(&clientColor, "client-color", "C", "#E6E70F", "Client text color in hex format (e.g., \"#FFFFFF\" or \"#E6E70F\").")
	editCmd.Flags().StringVarP(&textColor, "text-color", "T", "#E6E70F", "Event text color in hex format (e.g., \"#FFFFFF\" or \"#E6E70F\").")
	editCmd.Flags().StringVarP(&font, "font", "f", "/usr/local/share/fonts/Courgette-Regular.ttf", "Path to font to use for all text (e.g., /usr/share/fonts/TTF/MyFont.ttf).")
	editCmd.Flags().StringVarP(&leftIcon, "left", "L", "", "Path to icon image to display to the left of client text.")
	editCmd.Flags().StringVarP(&rightIcon, "right", "R", "", "Path to icon image to display to the right of client text.")

	editCmd.Flags().BoolVarP(&save, "save", "s", false, "Save the current command.")
	editCmd.Flags().BoolVarP(&run, "run", "r", false, "Run the saved command.")
	editCmd.MarkFlagsMutuallyExclusive("save", "run")
	// You either `run` the saved command or executed/`save` the current one, making `logo`
	// required.
	editCmd.MarkFlagsOneRequired("run", "logo")

	rootCmd.AddCommand(editCmd)
}

func wrapText(text string, fontSize int, videoW int, margin int) string {
	// Use a 0.7 factor (aprox 56px) for safer wrap.
	return wordwrap.WrapString(text, uint((videoW-margin)/(fontSize*7/10)))
}

// videoDuration take a path and run ffprobe to learn its duration, return the quotient and any error
func videoDuration(path string) (int, error) {
	// TODO: change to return a float
	ffprobe := exec.Command("ffprobe",
		"-loglevel",
		"error",
		"-show_entries",
		"format=duration",
		"-output_format",
		"default=noprint_wrappers=1:nokey=1",
		path)
	ffprobe.Stderr = os.Stderr
	out, err := ffprobe.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe.Run(%s): %w", ffprobe.String(), err)
	}
	str, _, _ := strings.Cut(string(out), ".")
	dur, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("strconv.Atoi(string(%v)): %w", out, err)
	}
	return dur, nil
}

// TODO: change checkArgs for validateArgs, and check that the colors have the
// correct format

// checkArgs ensure all passed arguments exists
func checkArgs() error {
	_, err := os.Stat(logo)
	if err != nil {
		return err
	}
	if image != "" {
		_, err = os.Stat(image)
		if err != nil {
			return err
		}
	}
	if music != "" {
		_, err = os.Stat(music)
		if err != nil {
			return err
		}
	}
	if font != "" {
		_, err = os.Stat(font)
		if err != nil {
			return err
		}
	}
	if leftIcon != "" {
		_, err = os.Stat(leftIcon)
		if err != nil {
			return err
		}
	}
	if rightIcon != "" {
		_, err = os.Stat(rightIcon)
		if err != nil {
			return err
		}
	}
	return nil
}

// maxLenght takes a possibly newline separated string and measure the longest
// line in it
// func maxLenght(t string, fontSize int) int {
// 	parts := strings.Split(t, "\n")
// 	lenght := 0
// 	for _, p := range parts {
// 		if lenght < len(p) {
// 			lenght = len(p)
// 		}
// 	}
// 	return lenght * fontSize * 7 / 10
// }
