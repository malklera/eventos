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

[x] Add command or flag to point to the assets like "1 - Indicaciones.odt". Change
my mind, not going to do this.

[x] In init command allow to pass a path.

[x] Change the durations to return float64 instead of int.

[x] Validate that all arguments of edit are valid, files exist, color are the correct format.

[x] Ensure copy can accept paths. Both path are mandatory.

[x] Passing clientColor flag is only valid with clientText

[x] Passing textColor flag is only valid with eventText

[x] Use flag.Changed instead of checking for ""

[x] Check the PreRunE in edit.go, i think it can be better, see about writing
test for it.

[x] Get ride of the client text related flags and code.

[x] Handle incorrect image size.

[ ] Implement -run and -save flags.

[ ] Improve error message for init, and provably other commands, "eventos init"
indicates to pass \[flags\] but init do not have flags, it should say "path to eventName" or something.

[ ] Change generated qr to 302,362*302,362 px or 80 mm

[ ] Change -logo to be optional.

[ ] Add a notification when copy, format or edit finish

[ ] Improve error for copy, add "eventos copy <source> <destination>"

[ ] Update the message when calling catalog subcommand

[ ] Reduce times show to two digits precision.

---

| imagen | music | logo |
-
| n | n | n |
| n | n | y |
| n | y | n |
| n | y | y |
| y | n | n |
| y | n | y |
| y | y | n |
| y | y | y |
-


---

go run . edit

Has to have
fade-in,all effects, fade-out

---

go run . edit -logo 500-500.png

Needs fade-in, all effects, cross-fade, logo

---

go run . edit -music ~/Videos/eventos/assets/musica/cortado/electronica/mix-electronica-yamila.opus

Cannot test, it just error with this

Error: accepts 0 arg(s), received 1

---

go run . edit -music ~/Videos/eventos/assets/musica/cortado/electronica/mix-electronica-yamila.opus -logo 500-500.png

---

---

go run . edit -image 500-500.png

---

---

go run . edit -image 500-500.png -logo 500-500.png

---

---

go run . edit -image 500-500.png -music ~/Videos/eventos/assets/musica/cortado/electronica/mix-electronica-yamila.opus

---

---

go run . edit -image 500-500.png -music ~/Videos/eventos/assets/musica/cortado/electronica/mix-electronica-yamila.opus -logo 500-500.png

---
