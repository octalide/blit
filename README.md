# blit

A mach-native immediate-mode GUI library, written entirely in Mach. No C
dependencies, no libm — just layout, input, a draw list, and an embedded
bitmap font.

blit is **backend-agnostic**. Widgets accumulate a draw list of textured,
colored triangles and report interaction results; the consumer feeds input
each frame, calls widget functions, then uploads and renders the draw list
through whatever it likes (e.g. [mach-gl](https://github.com/octalide/mach-gl)).
blit owns no window and no GL.

```mach
use blit;

# once: upload blit.font_atlas() as a texture, build a quad shader.
# per frame:
blit.begin(?ctx, ?input, screen_w, screen_h);
if (blit.button(?ctx, "step")) { ... }
blit.slider_f(?ctx, "eta", ?eta, 0.0, 1.0);
blit.end(?ctx);
# then: upload blit.draw_verts(?ctx) and draw as triangles with the atlas.
```

## Goals

- **Pure Mach.** Every line in Mach; no C/C++ bindings, no libm.
- **Backend-agnostic.** A draw list of `(x, y, u, v, rgba)` vertices plus an
  embedded font atlas; the consumer renders. Zero coupling to GL/GLFW.
- **Immediate mode.** No retained widget tree; the UI is rebuilt each frame
  from straight-line calls, state lives in the caller.
- **Small and legible.** A handful of modules, readable over clever.

## Status

v0, under active development. Consumed first by
[co](https://github.com/octalide/co)'s 3D organism viewer.

## Layout

```
src/
  blit.mach     library surface (re-exports the submodules)
  ...           draw list, input, context, font, widgets (landing incrementally)
```

## Conventions

`main`/`dev` long-lived branches; `feat/*` and `fix/*` branch off `dev` and PR
back; `dev` integrates to `main` for releases. Conventional commits, semver
tags, `branch/main` as the published ref.
