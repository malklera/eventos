[14:24, 23/3/2026] Mamá Claro Particular: Video
Foto principal 2/3 segundos.
Video central 15 segundos.
Transiciones... hacia adelante, lento y acelerar, transición rápida a efecto pantalla partida en 4 por 3/4 Segundos
Transición corta 1 segundo. 
Reversa
Efecto pin pong
[14:25, 23/3/2026] Mamá Claro Particular: Foto final de institucional

Video total

Imagen-Cliente+Transicion-crossfade(imagen cliente transparenta y se ve el video)+Video-plataforma+Transicion-crossfade+Imagen-Mimbi

Musica

Fade-In+                                                                                                                        +Fade-out

Video plataforma

Static client-imagen: IMAGE_TIME
Fade-out transition: TRANSITION_DURATION
transicion abajo hacia arriba: overlay, from bottom to top, input_video_path

hacia adelante: 

"lento y acelerar" video y audio o solo video?
slow/fast motion: minterpolate, mci

efecto pantalla partida en 4: xstack

Reversa: 

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
