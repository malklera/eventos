# Eventos

This CLI is used to manage and edit videos used in the 360 Platform service by
[Mimbi Glitter & Tattoo](https://www.instagram.com/mimbi_glitterbar).

To install.

```sh
TAG=$(git describe --tags --exact-match 2>/dev/null || echo dev)
go install -ldflags "-X 'main.tag=$TAG'"
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

[ ] Improve error message for init, and provably other commands, "eventos init"
indicates to pass \[flags\] but init do not have flags, it should say "path to eventName" or something.

[x] Passing clientColor flag is only valid with clientText

[x] Passing textColor flag is only valid with eventText

[ ] Implement -run and -save flags.

[x] Use flag.Changed instead of checking for ""

[x] Check the PreRunE in edit.go, i think it can be better, see about writing
test for it.
