package GUI

import (
	"GameOfLife/src/Watchdogs"
	"github.com/go-gl/gl/v4.6-core/gl"
	"os"
	"strings"
	"time"
)

func loadAndHotReloadCompileShaders(path string, program uint32) {
	// get all files in the shaders directory
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	// list of shaders that are currently loaded
	var shaders []Shader

	for _, file := range files {
		extension := strings.Split(file.Name(), ".")[1]

		// check if shader type is supported
		if extension != VertexShader &&
			extension != FragmentShader &&
			extension != ComputeShader &&
			extension != GeometryShader {
			continue
		}

		// register the shader file for hot reloading
		shader, err := os.ReadFile(path +
			"/" + file.Name())
		if err != nil {
			panic(err)
		}

		shaderCode := string(shader)

		// check if the shader is already loaded
		var shaderIndex int
		var shaderExists bool
		for i, s := range shaders {
			if s.shaderName == file.Name() {
				shaderIndex = i
				shaderExists = true
				break
			}
		}

		// if the shader is already loaded, check if it has changed
		if shaderExists {
			if shaders[shaderIndex].shaderCode == shaderCode {
				continue
			}
		}

		// if the shader is not loaded or has changed, compile it
		var shaderType uint32
		switch extension {
		case VertexShader:
			shaderType = gl.VERTEX_SHADER
		case FragmentShader:
			shaderType = gl.FRAGMENT_SHADER
		case ComputeShader:
			shaderType = gl.COMPUTE_SHADER
		case GeometryShader:
			shaderType = gl.GEOMETRY_SHADER
		}

		shaders = append(shaders, Shader{
			shaderFile: path + "/" + file.Name(),
			shaderName: file.Name(),
			shaderType: shaderType,
			shaderCode: shaderCode,
		})

		// add a file watcher to the shader file
		// if the file changes, recompile and attach all shaders
		// to the program

		_, err = Watchdogs.FileWatchdog(path+"/"+file.Name(), func(updatedData *string) error {
			shaders[shaderIndex].shaderCode = *updatedData
			recompileAndAttachShaders(shaders, program)
			return nil
		}, time.Second)
		if err != nil {
			return
		}
	}

	recompileAndAttachShaders(shaders, program)

}

func recompileAndAttachShaders(shaders []Shader, program uint32) {
	println("Recompiling shaders")
	// remove all shaders from the program
	for _, shader := range shaders {
		if shader.shaderID != 0 {
			gl.DetachShader(program, shader.shaderID)
			gl.DeleteShader(shader.shaderID)
		}
	}

	// compile all shaders
	compiledShaders := make([]uint32, 0)
	for _, shader := range shaders {
		compiledShader, err := compileShader(shader.shaderCode, shader.shaderType)
		if err != nil {
			panic(err)
		}
		compiledShaders = append(compiledShaders, compiledShader)
	}

	// attach all shaders to the program
	for index, shader := range compiledShaders {
		shaders[index].shaderID = shader
		gl.AttachShader(program, shader)
	}
}
