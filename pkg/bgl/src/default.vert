#version 460 core

in vec4 tex; // vertex and texture coordinates {x, y, u, v}

out vec2 uv;

uniform mat4 modl;
uniform mat4 view;
uniform mat4 proj;

void main() {
	uv = tex.zw;

	gl_Position = proj * view * modl * vec4(tex.xy, 0.0, 1.0);
}
