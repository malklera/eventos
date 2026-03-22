Make Google Drive folder with the event-name, make two folders, "Originales",
"Editados", upload the "instrucciones.txt" file, copy URL.

```sh
cd ~/Videos/eventos
./mkdir.sh <event-name>
```

Generate the QR for the event.

```sh
cd ~/Videos/eventos/event-name
qrtool encode -o <event-name>.svg -t svg "<link>"
```
Copy videos from phone to ~/Videos/eventos/event-name/original

```sh
./rename.sh <event-name>
```

```sh
./formater.sh <event-name>
```

Cut each video manually and export them to ~/Videos/eventos/event-name/cortado

In the process take a screenshot of each for the catalog.

Put all images into the "catalogo" directory.

```sh
cd ~/Videos/eventos/<event-name>/catalogo
perl-rename -ni 's/^(\d+)-.*/$1.png/' *.png
```

See the output of that, if it look ok, remove the -n flag

```sh
./edit.sh -e <event-name> [-t "<event-text>"] -m <music> -l <logo-video> \
            [-i <client-image>] [-c "<client-text>"] [-C "<client-color>"] \
            [-T "<text-color>"] [-f <font>] [-L <left-icon>] [-R <right-icon>]
```

Construct the command into a edit.txt file on the event-name directory.


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
