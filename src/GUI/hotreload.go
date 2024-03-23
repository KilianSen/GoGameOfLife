package GUI

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"log"
	"os"
	"time"
)

var hotLoadedShaders []hotLoadedShader
var lastTimeChecked time.Time = time.Now()

type hotLoadedShader struct {
	shaderPath string
	shaderHash string
}

func hotShaders(path string, program *uint32) {
	// check for shader changes if there are changes, remove old shaders, recompile the shaders and reattach them to the program

	// check if the delta time is greater than 1 second
	if time.Since(lastTimeChecked) < time.Second {
		return
	}
	lastTimeChecked = time.Now()

	// check if there are any new shaders
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	// add new shaders to the hotLoadedShaders list
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if isShader, _ := FileIsShader(file.Name()); !isShader {
			continue
		}

		// check if the shader is not in the hotLoadedShaders list
		var shaderExists bool
		for _, s := range hotLoadedShaders {
			if s.shaderPath == path+file.Name() {
				shaderExists = true
				break
			}
		}

		if !shaderExists {
			hotLoadedShaders = append(hotLoadedShaders, hotLoadedShader{shaderPath: path + file.Name(), shaderHash: ""})
		}
	}

	forceRecompile := false

	// remove shaders that are not in the directory anymore
	for i, s := range hotLoadedShaders {
		var shaderExists bool
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if s.shaderPath == path+file.Name() {
				shaderExists = true
				break
			}
		}

		if !shaderExists {
			forceRecompile = true
			hotLoadedShaders = append(hotLoadedShaders[:i], hotLoadedShaders[i+1:]...)
		}

	}

	// check if the shaders have changed

	recompile := false

	if !forceRecompile {
		for i, s := range hotLoadedShaders {
			fileData, err := os.ReadFile(s.shaderPath)
			if err != nil {
				panic(err)
			}
			hash := GetHash(fileData)
			if s.shaderHash == string(hash) {
				continue
			}
			hotLoadedShaders[i].shaderHash = string(hash)
			recompile = true
		}
	}

	if recompile || forceRecompile {
		// recompile the shaders and reattach them to the program
		gl.DeleteProgram(*program)
		*program = initOpenGL()

		// print every shader that is being recompiled
		debugPrint := ""
		for _, s := range hotLoadedShaders {
			debugPrint += " " + s.shaderPath
		}
		log.Println("Recompiling shaders:" + debugPrint)

		for _, s := range hotLoadedShaders {
			_, FileType := FileIsShader(s.shaderPath)

			// get last modified time
			fileInfo, err := os.Stat(s.shaderPath)
			if err != nil {
				panic(err)
			}
			for time.Since(fileInfo.ModTime()) < time.Second*2 {
				// wait for the file to be fully written
				time.Sleep(time.Millisecond * 250)
			}

			shaderCode, err := os.ReadFile(s.shaderPath)
			if err != nil {
				panic(err)
			}

			shader, err := compileShader(string(shaderCode), uint32(FileType))
			if err != nil {
				panic(err)
			}
			gl.AttachShader(*program, shader)

		}
		gl.LinkProgram(*program)
	}

}
