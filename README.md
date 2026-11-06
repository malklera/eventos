# Eventos

This CLI is used to manage and edit videos used in the 360 Platform service by
[Mimbi Glitter & Tattoo](https://www.instagram.com/mimbi_glitterbar).

To install.

```sh
go install -ldflags "-X 'github.com/malklera/eventos/cmd.tag=$(git describe --tags --exact-match 2>/dev/null || git rev-parse --short HEAD)'"
```

To use read the steps given by `eventos help`.

---

Video from camera

Resolution: 1080x1920 30 fps

---

## Contact info

Instagram: https://www.instagram.com/mimbi_glitterbar

WhatsApp: 343-5362802

Image with contact info

Resolution: 1080x1920

Font: Works Sans regular

Font size: 110

---

## TODO

- [x] Rotate the pdf 90 degress

- [x] bug in "eventos catalog">"Originals" it counts all files, it has to count only images

- [ ] Implement -run and -save flags.

- [x] Improve error message for init, and provably other commands, "eventos init"
indicates to pass \[flags\] but init do not have flags, it should say "path to eventName" or something.

- [x] Move the total video from before copying to the end so i can compare originals and copied

- [x] Add a notification when copy, format or edit finish

- [ ] Replace "X" and "!" for all notifications with actual icons.

- [ ] Reduce all time outputs to show two digits precision.

- [ ] Handle Ctrl+c interrupt, currently my cli closes but ffmpeg stay running in the background i think, check first

- [ ] Replace inkscape with a c++ library to convert svg to pdf

---

| imagen | music | logo |
---|---|---|
| n | n | n |
| n | n | y |
| n | y | n |
| n | y | y |
| y | n | n |
| y | n | y |
| y | y | n |
| y | y | y |


---

go run . edit

Has to have
fade-in,all effects, fade-out

---

go run . edit --logo 500-500.png

Needs fade-in, all effects, cross-fade, logo

---

go run . edit --music ~/Videos/eventos/assets/musica/cortado/electronica/mix-electronica-yamila.opus


---

go run . edit --music ~/Videos/eventos/assets/musica/cortado/electronica/mix-electronica-yamila.opus --logo 500-500.png

---

---

go run . edit --image 500-500.png

---

---

go run . edit --image 500-500.png --logo 500-500.png

---

---

go run . edit --image 500-500.png --music ~/Videos/eventos/assets/musica/cortado/electronica/mix-electronica-yamila.opus

---

---

go run . edit --image 500-500.png --music ~/Videos/eventos/assets/musica/cortado/electronica/mix-electronica-yamila.opus --logo 500-500.png

---
