# Eventos

This CLI is used to manage and edit videos used in the 360 Platform service by
[Mimbi Glitter & Tattoo](https://www.instagram.com/mimbi_glitterbar).

To install.

```sh
go install -ldflags "-X 'github.com/malklera/eventos/cmd.tag=$(git describe --tags --exact-match 2>/dev/null || echo dev)'"
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

- [ ] Implement -run and -save flags.

- [ ] Improve error message for init, and provably other commands, "eventos init"
indicates to pass \[flags\] but init do not have flags, it should say "path to eventName" or something.

- [ ] Improve error for copy, add "eventos copy <source> <destination>"

- [ ] Add a notification when copy, format or edit finish

- [ ] Reduce all time outputs to show two digits precision.

- [ ] Handle Ctrl+c interrupt, currently my cli closes but ffmpeg stay running in the background i think, check first

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
