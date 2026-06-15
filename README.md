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

# once, at startup:
blit.font.build_atlas(?a);                 # rasterize the font atlas
# upload blit.font.atlas_pixels() as an RGBA8 texture (see Rendering)

# per frame:
blit.context.begin(?ctx, in, screen_w, screen_h);
val panel: usize = blit.widget.begin_panel(?ctx, 8.0::f32, 8.0::f32, 200.0::f32);
blit.widget.text(?ctx, "controls");
if (blit.widget.button(?ctx, "step")) { ... }
blit.widget.checkbox(?ctx, "running", ?running);
blit.widget.slider_f(?ctx, "rate", ?rate, 0.0::f32, 1.0::f32);
blit.widget.end_panel(?ctx, panel);
blit.context.end(?ctx);

# then upload blit.context.draw_verts(?ctx) / draw_count(?ctx) and draw as
# triangles with the atlas bound (see Rendering).
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
  draw.mach     blit.draw    — Color, Vert, and the growable DrawList
  font.mach     blit.font    — embedded 8x8 ASCII font and its RGBA8 atlas
  input.mach    blit.input   — per-frame mouse input and hit-testing
  context.mach  blit.context — context, layout cursor, ids, painter helpers
  widget.mach   blit.widget  — panel, text, button, checkbox, slider
  bin/harness.mach           — headless smoke and test host (see Testing)
```

## Surface

`use blit;` binds the surface, which re-exports each submodule under its short
name. Reach everything through the submodule it lives in:

- `blit.draw` — `Color`, `rgba`, `Vert`, `DrawList`, `push_quad`, `count`, `data`
- `blit.font` — `build_atlas`, `atlas_pixels`, `white_uv`, `glyph_uv`,
  `text_width`, `ATLAS_W`, `ATLAS_H`, `GLYPH_W`, `GLYPH_H`
- `blit.input` — `Input`, `pressed`, `released`, `in_rect`
- `blit.context` — `Context`, `init`, `free`, `begin`, `end`, `next_id`, `ok`,
  `draw_verts`, `draw_count`, and the painters `quad` / `glyph` / `text_at`
- `blit.widget` — `begin_panel`, `end_panel`, `text`, `button`, `checkbox`,
  `slider_f`

Consumers may also import a submodule directly, e.g. `use w: blit.widget;`.

## Rendering

blit produces one stream of vertices that draws both solid rectangles and text
through a single shader and texture. Solid quads sample a reserved white texel in
the atlas, so `color * texel` yields the flat color; glyph quads sample the
glyph's cell.

**Atlas.** `blit.font.build_atlas(?a)` rasterizes the font once into a
process-lifetime buffer. `blit.font.atlas_pixels()` returns it: RGBA8,
`ATLAS_W`×`ATLAS_H` (128×48) pixels, `ATLAS_W*ATLAS_H*4` bytes. Every pixel is
white rgb; alpha is 255 on lit glyph pixels and the white cell, 0 elsewhere
(white rgb everywhere avoids dark fringing under linear filtering). Nearest
filtering keeps the pixel font crisp.

**Vertex layout.** `blit.draw.Vert` is 8 `f32` = 32-byte stride:

| attribute | type | offset |
|-----------|------|--------|
| `aPos`   (x, y)       | vec2 | 0  |
| `aUV`    (u, v)       | vec2 | 8  |
| `aColor` (r, g, b, a) | vec4 | 16 |

Positions are in pixels; uv are in [0, 1]; color is straight (non-premultiplied)
rgba in [0, 1].

**Shader.** One uniform `vec2 uScreen`:

```glsl
// vertex
gl_Position = vec4(aPos.x / uScreen.x * 2.0 - 1.0,
                   1.0 - aPos.y / uScreen.y * 2.0, 0.0, 1.0);
// fragment
FragColor = aColor * texture(atlas, aUV);
```

**Per frame.** Fill an `Input` (`mx`, `my`, `down`; blit carries `prev_down`
itself), call `blit.context.begin`, run the widgets, call `blit.context.end`,
then upload `draw_verts(?ctx)` / `draw_count(?ctx)` into one dynamic VBO and draw
`GL_TRIANGLES` with alpha blending (`SRC_ALPHA`, `ONE_MINUS_SRC_ALPHA`). Because
solid quads sample the white texel, the same shader and texture draw both
rectangles and text in a single call.

## Testing

A static library carries no `_start`, so test executables cannot link against the
`[lib.blit]` artifact directly. The tests are hosted by the `[bin.blit]` smoke
harness, which imports `std.runtime`:

```
mach build              # compile clean
mach test --bin blit    # build and run every module's tests
```

## Conventions

`main`/`dev` long-lived branches; `feat/*` and `fix/*` branch off `dev` and PR
back; `dev` integrates to `main` for releases. Conventional commits, semver
tags, `branch/main` as the published ref.
