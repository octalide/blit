#version 460 core

in vec2 uv;

out vec4 out_color;

uniform sampler2D tex;
uniform vec4 color;

void main() {
	out_color = texture(tex, uv);
}
