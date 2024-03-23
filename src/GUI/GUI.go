package GUI

import (
	"github.com/go-gl/gl/v4.6-core/gl" // OR: github.com/go-gl/gl/v2.1/gl
	"github.com/go-gl/glfw/v3.3/glfw"
	"runtime"
)

var (
	Triangle = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}

	Square = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,

		-0.5, 0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
	}
)

type Gui struct {
	window           *glfw.Window
	program          uint32
	cells            *[][]bool
	applyShadersFunc func(path string, program *uint32)

	width  int
	height int

	title string
}

func NewGUI() *Gui {
	runtime.LockOSThread()
	defer glfw.Terminate()
	return &Gui{
		window:           initGlfw(),
		program:          initOpenGL(),
		cells:            nil,
		applyShadersFunc: hotShaders,
		width:            0,
		height:           0,
		title:            "",
	}
}

func newVertices(x, y, rows, cols int) uint32 {
	points := make([]float32, len(Square), len(Square))
	copy(points, Square)

	for i := 0; i < len(points); i++ {
		var position float32
		var size float32
		switch i % 3 {
		case 0:
			size = 1.0 / float32(cols)
			position = float32(x) * size
		case 1:
			size = 1.0 / float32(rows)
			position = float32(y) * size
		default:
			continue
		}

		if points[i] < 0 {
			points[i] = (position * 2) - 1
		} else {
			points[i] = ((position + size) * 2) - 1
		}
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	return vao
}
