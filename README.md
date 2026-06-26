# Eventos

This CLI is used to manage and edit videos used in the 360 Platform service by
(https://www.instagram.com/mimbi_glitterbar)[Mimbi Glitter & Tattoo].

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

Contact info
Instagram: https://www.instagram.com/mimbi_glitterbar
WhatsApp: 343-5362802

Image with contact info
Resolution: 1080x1920
Font: Works Sans regular
Font size: 110

---

# TODO

[ ] Add command or flag to point to the assets like "1 - Indicaciones.odt"

[x] In init command allow to pass a path.

[ ] Change the durations to return float64 instead of int.

[ ] Validate that all arguments of edit are valid, files exist, color are the correct format.

[ ] Ensure copy can accept paths. If no path is passed it is implicit ".".

[ ] Improve error message for init, and provably other commands, "eventos init"
indicates to pass \[flags\] but init do not have flags, it should say "path to eventName" or something.
