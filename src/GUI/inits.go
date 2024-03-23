package GUI

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
)

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	program := gl.CreateProgram()
	initShaders(program)
	gl.LinkProgram(program)

	return program
}

func initShadersOld(program uint32) {
	shaders, err := loadShaders("./src/GUI/shaders/")
	if err != nil {
		panic(err)
	}

	for _, rawShader := range shaders {
		shader, err := compileShader(rawShader.shaderCode, rawShader.shaderType)
		if err != nil {
			panic(err)
		}

		log.Default().Println("Compiled Shader:", rawShader.shaderName, "Type:", rawShader.shaderType, "File:", rawShader.shaderFile)

		gl.AttachShader(program, shader)
	}
}

func initShaders(program uint32) {
	loadAndHotReloadCompileShaders("./src/GUI/shaders/", program)
}
