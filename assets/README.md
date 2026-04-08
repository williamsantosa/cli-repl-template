# Assets

Place art files here (`.png`, `.jpg`, `.gif`, `.bmp`, `.txt`) before building.
They are embedded into the binary via `//go:embed` and loaded at runtime
without needing the original files on disk.

Set `art.source` in `config.yaml` to the filename (e.g. `"my-art.gif"`).
