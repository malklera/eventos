```sh
$ cd ~/Videos/eventos
$ ./mkdir.sh name-event
```

Copy videos from phone to ~/Videos/eventos/name-event/original

```sh
$ ./formater.sh name-event
```

Cut each video manually and export them to ~/Videos/eventos/name-event/cortado

```sh
$ ./edit.sh -e <event_name> -t "<event_text>" -m <music> -l <logo_video> \
            [-i <client_image>] [-c "<client_text>"] [-C "<client_color>"] \
            [-T "<text_color>"] [-f <font>] [-L <left_icon>] [-R <right_icon>]
```

Construct the command here, then copy to the terminal

```sh
./edit.sh
```


---

# Estructura del video

(transicion fade-in)+(imagen cliente)   +(transicion xfade)     +(video filmado)+(transicion xfade)     +(video logo)
(transicion fade-in)+(texto cliente)    +(transicion fade-out)  
                                         (transicion fade-in)   +(texto video)  +(transicion fade-out)
(transicion fade-in)+(musica)                                                   +(transicion acrossfade)+(musica logo)+(transicion fade-out)

---

Video from camera

Resolution: 1080x1920 30 fps

---

Contact video

Resolution: 720x1280 25 fps


Contact info
@mimbi_glitterbar
343-5362802

Image with contact info
Resolution: 1440x2560
Icon size:
Font: Works Sans regular
Font size: 110


Command to make it the resolution of the video from camera

```sh
ffmpeg -i video_source_path \
       -vf "scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2,fps=30" \
       -c:v libx264 \
       -preset medium \
       -crf 23 \
       -c:a aac \
       -b:a 192k \
       -movflags +faststart \
       video_output_path
```

---

Image from client
Resolution: min 720x1280

---

Fonts

Courgette-Regular.ttf
MysteryQuest-Regular.ttf
SansitaSwashed-VariableFont_wght.ttf
