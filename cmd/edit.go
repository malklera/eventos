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
	logo      string
	music     string
	image     string
	eventText string
	textColor string
	font      string
	save      bool
	run       bool
)

var (
	ErrTextColorWithoutText = errors.New("'text-color' flag is only allowed when 'text' is passed")
	ErrFontWithoutText      = errors.New("'font' flag is only allowed when 'text' is passed")
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
		eventTextX := "(w-tw)/2"
		eventTextY := "H-th-100"
		eventTextMargin := 100

		cuttedDir := "cortado"

		// If eventText is not empty, wrap it as needed
		if cmd.Flags().Changed("text") {
			eventText = wrapText(eventText, fontSize, videoW, eventTextMargin)
		}

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

			hasImage := cmd.Flags().Changed("image")
			hasLogo := cmd.Flags().Changed("logo")
			hasMusic := cmd.Flags().Changed("music")
			hasText := cmd.Flags().Changed("text")

			// Compute totalT based on which inputs are present
			var totalT float64
			switch {
			case hasImage && hasLogo:
				totalT = imageT + videoT + imageT
			case hasImage && !hasLogo:
				totalT = imageT + videoT
			case !hasImage && hasLogo:
				totalT = videoT + imageT - transitionT
			default:
				totalT = videoT
			}

			// Check music length
			if hasMusic && musicT < totalT {
				fmt.Fprintf(os.Stderr, "Error: music too short, musicT: %f, totalT: %f\n", musicT, totalT)
				continue
			}

			// Build ffmpeg inputs and track stream indices
			ffmpegEdit := []string{"ffmpeg", "-v", "error"}
			nextIdx := 0
			imgIdx := -1
			videoIdx := -1
			logoIdx := -1
			musicIdx := -1

			if hasImage {
				ffmpegEdit = append(ffmpegEdit, "-loop", "1", "-t", strconv.FormatFloat(imageT, 'f', -1, 64), "-i", image)
				imgIdx = nextIdx
				nextIdx++
			}
			ffmpegEdit = append(ffmpegEdit, "-i", videoPath)
			videoIdx = nextIdx
			nextIdx++
			if hasLogo {
				ffmpegEdit = append(ffmpegEdit, "-loop", "1", "-t", strconv.FormatFloat(imageT, 'f', -1, 64), "-i", logo)
				logoIdx = nextIdx
				nextIdx++
			}
			if hasMusic {
				ffmpegEdit = append(ffmpegEdit, "-i", music)
				musicIdx = nextIdx
				nextIdx++
			}

			// Build filter_complex: always apply effects pipeline, optionally prepend/append image/logo
			filterComplex := ""

			// --- Always build effects pipeline for the video ---
			videoFadeOutStart := videoT - transitionT
			pLen := (videoT + 2*transitionT) / 5
			part1End := pLen
			part2Start := part1End - transitionT
			part2End := part2Start + pLen
			part3Start := part2End
			part3End := part3Start + pLen
			part4Start := part3End
			part4End := part4Start + pLen
			part5Start := part4End - transitionT

			fmt.Printf("Image: %v, Logo: %v, Video: %f, Total: %f\n", hasImage, hasLogo, videoT, totalT)

			// Scale video
			filterComplex = fmt.Sprintf("[%d:v]scale=%d:%d,fps=30,setpts=PTS-STARTPTS[video];", videoIdx, videoW, videoH)
			// Video background for slide effects
			filterComplex += fmt.Sprintf("color=c=black:s=%dx%d:d=%f:r=30[video_bg];", videoW, videoH, videoT)
			// Video slide overlay (slide-up entrance + slide-out exit)
			filterComplex += fmt.Sprintf("[video_bg][video]overlay=x=0:y='if(lt(t,%f),H-H*t/%f,if(gt(t,%f),-H*(t-%f)/%f,0))':format=auto[video_slide_raw];", transitionT, transitionT, videoFadeOutStart, videoFadeOutStart, transitionT)
			// Split into parts for effects
			filterComplex += "[video_slide_raw]split=3[raw1][raw2to4][raw5];"
			filterComplex += fmt.Sprintf("[raw1]trim=start=0:end=%f,setpts=PTS-STARTPTS[part1];", part1End)
			filterComplex += fmt.Sprintf("[raw2to4]trim=start=%f:end=%f,setpts=PTS-STARTPTS[part_rest_raw];", part2Start, part4End)
			filterComplex += fmt.Sprintf("[raw5]trim=start=%f,setpts=PTS-STARTPTS[part5];", part5Start)

			// Part 2 Effect: Quad-split (2x2 grid)
			filterComplex += "[part_rest_raw]split=2[pr_main][pr_quad_src];"
			filterComplex += fmt.Sprintf("[pr_quad_src]scale=%d:%d,split=4[q1][q2][q3][q4];", (videoW / 2), (videoH / 2))
			filterComplex += "[q1][q2]hstack[top_row];"
			filterComplex += "[q3][q4]hstack[bot_row];"
			filterComplex += "[top_row][bot_row]vstack[quad];"
			filterComplex += fmt.Sprintf("[pr_main][quad]overlay=enable='between(t,0,%f)':eof_action=pass[part_rest_step1];", pLen)
			// Transition: Frei0r Distort0r between Part 2 and Part 3
			filterComplex += "[part_rest_step1]split=2[pr_step15_main][distort_src];"
			filterComplex += "[distort_src]frei0r=filter_name=distort0r:filter_params=0.5|0.005|y|0.25[distorted];"
			filterComplex += fmt.Sprintf("[pr_step15_main][distorted]overlay=enable='between(t,%f,%f)':eof_action=pass[part_rest_step15];", (pLen - transitionT/2), (pLen + transitionT/2))
			// Part 3 Effect: Slow/Fast time warp
			p3LocalStart := part3Start - part2Start
			p3LocalEnd := part3End - part2Start
			filterComplex += "[part_rest_step15]split=2[pr_step2_main][p3_src];"
			filterComplex += fmt.Sprintf("[p3_src]trim=start=%f:end=%f,setpts=PTS-STARTPTS[p3_trim];", p3LocalStart, p3LocalEnd)
			filterComplex += fmt.Sprintf("[p3_trim]setpts='if(lt(T, %f/4), 2*PTS + %f/TB, (2/3)*PTS + (%f/3 + %f)/TB)'[p3_warp];", pLen, p3LocalStart, pLen, p3LocalStart)
			filterComplex += fmt.Sprintf("[pr_step2_main][p3_warp]overlay=enable='between(t,%f,%f)':eof_action=pass[part_rest_step2];", p3LocalStart, p3LocalEnd)
			// Transition: Spin effect between Part 3 and Part 4
			filterComplex += fmt.Sprintf("[part_rest_step2]rotate=a='if(between(t,%f,%f), 2*PI*(t-%f)/%f, 0)':ow=iw:oh=ih:c=black[part_rest_step25];", (pLen*2 - transitionT/2), (pLen*2 + transitionT/2), (pLen*2 - transitionT/2), transitionT)
			// Part 4 Effect: 3x3 grid
			p4LocalStart := part4Start - part2Start
			p4LocalEnd := part4End - part2Start
			filterComplex += "[part_rest_step25]split=2[pr_step3_main][pr_grid3_src];"
			filterComplex += fmt.Sprintf("[pr_grid3_src]scale=%d:%d,split=9[g1][g2][g3][g4][g5][g6][g7][g8][g9];", (videoW / 3), (videoH / 3))
			filterComplex += "[g1][g2][g3]hstack=inputs=3[row1];"
			filterComplex += "[g4][g5][g6]hstack=inputs=3[row2];"
			filterComplex += "[g7][g8][g9]hstack=inputs=3[row3];"
			filterComplex += "[row1][row2][row3]vstack=inputs=3[grid3];"
			filterComplex += fmt.Sprintf("[pr_step3_main][grid3]overlay=enable='between(t,%f,%f)':eof_action=pass[part_rest_step4];", p4LocalStart, p4LocalEnd)
			// Part 4 Reverse Boomerang Effect
			p4LocalMid := p4LocalStart + (p4LocalEnd-p4LocalStart)/2
			filterComplex += "[part_rest_step4]split=2[pr_step4_main][p4_to_rev];"
			filterComplex += fmt.Sprintf("[p4_to_rev]trim=start=%f:end=%f,setpts=PTS-STARTPTS,reverse,setpts=PTS+(%f/TB)[p4_rev];", p4LocalStart, p4LocalMid, p4LocalMid)
			filterComplex += fmt.Sprintf("[pr_step4_main][p4_rev]overlay=enable='between(t,%f,%f)':eof_action=pass[part_rest];", p4LocalMid, p4LocalEnd)
			// Xfade 1: Part 1 and Parts 2-4
			xFadeOffset1 := part1End - transitionT
			filterComplex += fmt.Sprintf("[part1][part_rest]xfade=transition=fade:duration=%f:offset=%f[video_step1];", transitionT, xFadeOffset1)
			// Xfade 2: video_step1 and Part 5
			xFadeOffset2 := part4End - transitionT
			filterComplex += fmt.Sprintf("[video_step1][part5]xfade=transition=fade:duration=%f:offset=%f[video_effects];", transitionT, xFadeOffset2)

			// --- Build list of streams to concatenate ---
			var concatStreams []string

			// If hasImage, build client image part
			if hasImage {
				clientFadeOutStart := imageT - transitionT
				filterComplex += fmt.Sprintf("[%d:v]scale=%d:%d:force_original_aspect_ratio=disable,setsar=1,fps=30,setpts=PTS-STARTPTS,fade=t=out:st=%f:d=%f[client];", imgIdx, videoW, videoH, clientFadeOutStart, transitionT)
				concatStreams = append(concatStreams, "client")
			}

			// Add video effects
			concatStreams = append(concatStreams, "video_effects")

			// If hasLogo, build logo part
			if hasLogo {
				filterComplex += fmt.Sprintf("[%d:v]scale=%d:%d:force_original_aspect_ratio=disable,setsar=1,fps=30,setpts=PTS-STARTPTS,settb=AVTB,fade=t=in:st=0:d=%f[logo];", logoIdx, videoW, videoH, transitionT)
				concatStreams = append(concatStreams, "logo")
			}

			// --- Concatenate all streams ---
			n := len(concatStreams)
			if n > 1 {
				var concatInputs string
				for _, s := range concatStreams {
					concatInputs += fmt.Sprintf("[%s]", s)
				}
				filterComplex += fmt.Sprintf("%sconcat=n=%d:v=1:a=0,settb=AVTB[video_combined];", concatInputs, n)
			} else {
				filterComplex += "[video_effects]null[video_combined];"
			}

			// --- Text overlay ---
			curV := "video_combined"
			if hasText {
				var textFadeInStart, textFadeInEnd, textFadeOutStart float64
				if hasImage {
					textFadeInStart = imageT - transitionT
					textFadeInEnd = imageT
				}
				if hasLogo {
					if hasImage {
						textFadeOutStart = imageT + videoT - transitionT
					} else {
						textFadeOutStart = videoT - transitionT
					}
				}

				var alpha string
				switch {
				case hasImage && hasLogo:
					alpha = fmt.Sprintf("alpha='if(lt(t,%f),0,if(lt(t,%f),(t-%f)/%f,if(gt(t,%f),(1-(t-%f)/%f),1)))'", textFadeInStart, textFadeInEnd, textFadeInStart, transitionT, textFadeOutStart, textFadeOutStart, transitionT)
				case hasImage && !hasLogo:
					alpha = fmt.Sprintf("alpha='if(lt(t,%f),0,if(lt(t,%f),(t-%f)/%f,1))'", textFadeInStart, textFadeInEnd, textFadeInStart, transitionT)
				case !hasImage && hasLogo:
					alpha = fmt.Sprintf("alpha='if(gt(t,%f),(1-(t-%f)/%f),1)'", textFadeOutStart, textFadeOutStart, transitionT)
				}

				if alpha != "" {
					filterComplex += fmt.Sprintf("[%s]drawtext=text='%s':x=%s:y=%s:fontfile=%s:fontsize=%d:fontcolor=%s:borderw=%d:bordercolor=%s:box=1:boxcolor=0x00000088:boxborderw=%d:line_spacing=%d:text_align=center:%s[v_ev];", curV, eventText, eventTextX, eventTextY, font, fontSize, fontColor, borderW, borderColor, boxPadding, lineSpacing, alpha)
				} else {
					filterComplex += fmt.Sprintf("[%s]drawtext=text='%s':x=%s:y=%s:fontfile=%s:fontsize=%d:fontcolor=%s:borderw=%d:bordercolor=%s:box=1:boxcolor=0x00000088:boxborderw=%d:line_spacing=%d:text_align=center[v_ev];", curV, eventText, eventTextX, eventTextY, font, fontSize, fontColor, borderW, borderColor, boxPadding, lineSpacing)
				}
				curV = "v_ev"
			}
			filterComplex += fmt.Sprintf("[%s]null[v_with_text];", curV)

			// --- Audio ---
			if hasMusic {
				filterComplex += fmt.Sprintf("[%d:a]afade=t=in:st=0:d=%f,atrim=0:%f,asetpts=PTS-STARTPTS,afade=t=out:st=%f:d=%f[a_out];", musicIdx, transitionT, totalT, (totalT - transitionT), transitionT)
			} else {
				filterComplex += fmt.Sprintf("anullsrc=channel_layout=stereo:sample_rate=44100,atrim=0:%f,asetpts=PTS-STARTPTS[a_out];", totalT)
			}
			filterComplex += "[v_with_text]null[v_out]"

			ffmpegEdit = append(ffmpegEdit, "-filter_complex", filterComplex)

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
	editCmd.Flags().StringVarP(&textColor, "text-color", "T", "0xE6E70F", "Event text color in hex format (e.g., \"0xFFFFFF\" or \"0xE6E70F\").")
	editCmd.Flags().StringVarP(&font, "font", "f", "/usr/local/share/fonts/Courgette-Regular.ttf", "Path to font to use for all text (e.g., /usr/share/fonts/TTF/MyFont.ttf).")

	editCmd.Flags().BoolVarP(&save, "save", "s", false, "Save the current command.")
	editCmd.Flags().BoolVarP(&run, "run", "r", false, "Run the saved command.")
	editCmd.MarkFlagsMutuallyExclusive("save", "run")

	rootCmd.AddCommand(editCmd)
}

func preRun(cmd *cobra.Command) error {
	if !cmd.Flags().Changed("text") && cmd.Flags().Changed("text-color") {
		return ErrTextColorWithoutText
	}
	if !cmd.Flags().Changed("text") && cmd.Flags().Changed("font") {
		return ErrFontWithoutText
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
	o, _ := strings.CutSuffix(string(out), "\n")
	dur, err := strconv.ParseFloat(o, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseFloat(string(%v), 64): %w", out, err)
	}
	return dur, nil
}

// validArgs ensure all passed arguments are valid, mostly if they exist
func validArgs(cmd *cobra.Command) error {
	if cmd.Flags().Changed("logo") {
		_, err := os.Stat(logo)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
	}

	if cmd.Flags().Changed("image") {
		_, err := os.Stat(image)
		if err != nil {
			return err
		}
	}

	if cmd.Flags().Changed("music") {
		_, err := os.Stat(music)
		if err != nil {
			return err
		}
	}

	if cmd.Flags().Changed("font") {
		_, err := os.Stat(font)
		if err != nil {
			return err
		}
	}

	if cmd.Flags().Changed("text-color") {
		if !colorRe.MatchString(textColor) {
			return fmt.Errorf("'%s' is an invalid color, use e.g. '0xRRGGBB'", textColor)
		}
	}

	return nil
}
