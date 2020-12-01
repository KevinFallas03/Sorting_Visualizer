package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"./algorithms"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width              = 1280
	height             = 700
	vertexShaderSource = `
		#version 410
		in vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410
		out vec4 frag_colour;
		void main() {
			frag_colour = vec4(1, 1, 1, 1.0);
		}
	` + "\x00"

	rows = 200

	threshold = 0.15
	fps       = 10
)

var (
	columns = 0
	square  = []float32{
		-0.1, 0.1, 0,
		-0.1, -0.1, 0, //-0.1, -0.5, 0,
		0.1, -0.1, 0, //0.1, -0.5, 0,

		-0.1, 0.1, 0,
		0.1, 0.1, 0,
		0.1, -0.1, 0, //0.1, -0.5, 0,
	}
)

type bar struct {
	drawable uint32
	color    bool
}

func generateList() []int {
	rand.Seed(time.Now().UnixNano())
	//size := int(rand.Int31n(10000 + 1))
	size := 50
	numberList := make([]int, size, size)
	for x := range numberList {
		numberList[x] = int(rand.Int31n(31 + 1))
	}
	return numberList
}

func main() {
	//GENERA LA LISTA DE NUMEROS
	numberList := generateList()
	columns = len(numberList) + int(float32(len(numberList))*0.05)

	//INICIALIZA DATOS PARA CADA ALGORITMO

	//BUBBLE
	bubbleList := make([]int, len(numberList), len(numberList)) //Crea una lista
	copy(bubbleList, numberList)                                //Llena la lista
	bubbleChannel := make(chan []int)                           //Crea un canal

	//SELECTION
	selectionList := make([]int, len(numberList), len(numberList)) //Crea la lista
	copy(selectionList, numberList)                                //Llena la lista
	selectionChannel := make(chan []int)                           //Crea un canal

	//INSERTION
	insertionList := make([]int, len(numberList), len(numberList)) //Crea la lista
	copy(insertionList, numberList)                                //Llena la lista
	insertionChannel := make(chan []int)                           //Crea un canal

	//HEAP
	heapList := make([]int, len(numberList), len(numberList)) //Crea la lista
	copy(heapList, numberList)                                //Llena la lista
	heapChannel := make(chan []int)                           //Crea un canal

	//INICIA CADA ALGORITMO
	go algorithms.HeapSort(heapList, heapChannel)
	go algorithms.InsertionSort(insertionList, insertionChannel)
	go algorithms.SelectionSort(selectionList, selectionChannel)
	go algorithms.BubbleSort(bubbleList, bubbleChannel)

	//MOSTRAR VENTANA
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	//CREA TEMPORALES
	bubbleTemp := make([]int, len(numberList), len(numberList))
	copy(bubbleTemp, numberList)

	selectionTemp := make([]int, len(numberList), len(numberList))
	copy(selectionTemp, numberList)

	insertionTemp := make([]int, len(numberList), len(numberList))
	copy(insertionTemp, numberList)

	heapTemp := make([]int, len(numberList), len(numberList))
	copy(heapTemp, numberList)

	timer := 0
	percentage := float32(columns) * 0.03

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		//OBTIENE INFORMACION DE LOS CANALES
		heapData := <-heapChannel
		insertionData := <-insertionChannel
		bubbleData := <-bubbleChannel
		selectionData := <-selectionChannel

		if timer%int(percentage) == 0 {

			if len(heapData) == 0 {
				setBars(10.2, heapTemp, true)
			} else {
				setBars(10.2, heapData, false)
				heapTemp = heapData
			}

			if len(insertionData) == 0 {
				setBars(6.8, insertionTemp, true)
			} else {
				setBars(6.8, insertionData, false)
				insertionTemp = insertionData
			}

			if len(bubbleData) == 0 {
				setBars(3.4, bubbleTemp, true)
			} else {
				setBars(3.4, bubbleData, false)
				bubbleTemp = bubbleData
			}

			if len(selectionData) == 0 {
				setBars(0, selectionTemp, true)
			} else {
				setBars(0, selectionData, false)
				selectionTemp = selectionData
			}
			glfw.PollEvents()
			window.SwapBuffers()
		}

		timer++
	}
}

func setBars(y float32, data []int, color bool) {
	for i := range data {
		newBar(i, y, data[i], color)
	}
}
func newBar(x int, y float32, value int, color bool) {
	points := make([]float32, len(square), len(square))
	copy(points, square)

	for i := 0; i < len(points); i++ {
		var position float32
		var size float32
		var m float32
		switch i % 3 {
		case 0:
			size = (2.0 / float32(columns)) / 2
			position = float32(x) * size
			m = 0
		case 1:
			size = (float32(value) / float32(rows)) / 2
			position = 0
			m = y / 10
		default:
			continue
		}

		if points[i] < 0 {
			points[i] = ((position * 2) - 1) + m
		} else {
			points[i] = (((position + size) * 2) - 1) + m
		}
	}

	bar := &bar{
		drawable: makeVao(points),
		color:    color,
	}
	bar.draw()
}
func (c *bar) draw() {
	gl.ColorMask(true, c.color, false, false)
	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Sorting Algorithms", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	// version := gl.GoStr(gl.GetString(gl.VERSION))
	// log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
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

// const (
// 	width  = 500
// 	height = 500
// )

// func generateList() []int {
// 	rand.Seed(time.Now().UnixNano())
// 	size := int(rand.Int31n(100 + 1))
// 	numberList := make([]int, size, size)
// 	for x := range numberList {
// 		numberList[x] = int(rand.Int31n(31 + 1))
// 	}
// 	return numberList
// }

// func main() {
// 	c := make(chan []int)
// 	m := generateList()
// 	list1 := make([]int, len(m), len(m))
// 	copy(list1, m)
// 	go selectionSort(list1, c)

// 	runtime.LockOSThread()

// 	// list2 := make([]int, len(m), len(m))
// 	// copy(list2,m)
// 	// go insertionSort(list2)

// 	window := initGlfw()
// 	defer glfw.Terminate()
// 	program := initOpenGL()

// 	for !window.ShouldClose() {
// 		x := <-c // receive from c
// 		fmt.Println(x)
// 		draw(window, program)
// 	}

// }
// func draw(window *glfw.Window, program uint32) {
// 	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
// 	gl.UseProgram(program)

// 	glfw.PollEvents()
// 	window.SwapBuffers()
// }

// // initGlfw initializes glfw and returns a Window to use.
// func initGlfw() *glfw.Window {
// 	if err := glfw.Init(); err != nil {
// 		panic(err)
// 	}

// 	glfw.WindowHint(glfw.Resizable, glfw.False)
// 	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
// 	glfw.WindowHint(glfw.ContextVersionMinor, 1)
// 	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
// 	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

// 	window, err := glfw.CreateWindow(width, height, "Conway's Game of Life", nil, nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	window.MakeContextCurrent()

// 	return window
// }

// // initOpenGL initializes OpenGL and returns an intiialized program.
// func initOpenGL() uint32 {
// 	if err := gl.Init(); err != nil {
// 		panic(err)
// 	}
// 	version := gl.GoStr(gl.GetString(gl.VERSION))
// 	log.Println("OpenGL version", version)

// 	prog := gl.CreateProgram()
// 	gl.LinkProgram(prog)
// 	return prog
// }

// func selectionSort(data []int, c chan []int) {
// 	length := len(data)
// 	for i := 0; i < length; i++ {
// 		maxIndex := 0
// 		for j := 1; j < length-i; j++ {
// 			if data[j] > data[maxIndex] {
// 				maxIndex = j
// 			}
// 		}
// 		data[length-i-1], data[maxIndex] = data[maxIndex], data[length-i-1]
// 		c <- data
// 	}
// }
// func insertionSort(data []int) []int {
// 	for i := 1; i < len(data); i++ {
// 		if data[i] < data[i-1] {
// 			j := i - 1
// 			temp := data[i]
// 			for j >= 0 && data[j] > temp {
// 				data[j+1] = data[j]
// 				j--
// 			}
// 			data[j+1] = temp
// 			fmt.Println("Insertion: ", data)
// 		}
// 	}
// 	return data
// }
