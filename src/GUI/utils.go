package GUI

import (
	"crypto/sha1"
	"github.com/go-gl/gl/v4.6-core/gl"
	"strings"
)

func FileIsShader(fullPath string) (bool, uint64) {
	// get file extension and check what type of shader it is
	extension := strings.Split(fullPath, ".")[len(strings.Split(fullPath, "."))-1]
	var shaderType uint64
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
		return false, 0
	}

	return true, shaderType
}

func GetHash(data []byte) []byte {
	hashCreator := sha1.New()
	_, err := hashCreator.Write(data)
	if err != nil {
		panic(err)
	}

	return hashCreator.Sum(nil)
}
