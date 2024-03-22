package GUI

import (
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"os"
	"strings"
)

type Shader struct {
	shaderFile     string
	shaderName     string
	shaderType     int
	shaderCode     string
	compiledShader uint32
}

const (
	VertexShader   = "vert"
	FragmentShader = "frag"
	ComputeShader  = "comp"
	GeometryShader = "geom"
)

func loadShaders(path string) ([]Shader, error) {
	// load all shaders from a directory
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	shaders := make([]Shader, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// get file extension and check what type of shader it is
		extension := strings.Split(file.Name(), ".")[1]
		var shaderType int
		switch extension {
		case VertexShader:
			shaderType = gl.VERTEX_SHADER
		case FragmentShader:
			shaderType = gl.FRAGMENT_SHADER
		case ComputeShader:
			shaderType = gl.COMPUTE_SHADER
		case GeometryShader:
			shaderType = gl.GEOMETRY_SHADER
		default:
			return nil, fmt.Errorf("unknown shader type: %v", extension)
		}

		shaderCode, err := os.ReadFile(path + "/" + file.Name())
		if err != nil {
			return nil, err
		}

		shaders = append(
			shaders,
			Shader{
				shaderType: shaderType,
				shaderCode: string(shaderCode),
				shaderName: file.Name(),
				shaderFile: path + "/" + file.Name(),
			},
		)
	}
	return shaders, nil

}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

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

		compiledShader, err := compileShader(shaderCode, shaderType)
		if err != nil {
			panic(err)
		}

		// if the shader is already loaded, delete the old one
		if shaderExists {
			gl.DeleteShader(shaders[shaderIndex].shaderType)
		}

		// attach the new shader to the program
		gl.AttachShader(program, compiledShader)

		// link the program
		gl.LinkProgram(program)

	}
}
