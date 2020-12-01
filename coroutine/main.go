package main

import (
	"math/rand"
	"runtime"
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
	columns  = 0
	finished = 0
	square   = []float32{
		-0.1, 0.1, 0,
		-0.1, -0.1, 0,
		0.1, -0.1, 0,

		-0.1, 0.1, 0,
		0.1, 0.1, 0,
		0.1, -0.1, 0,
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
	bubbleTemp := numberList
	selectionTemp := numberList
	insertionTemp := numberList
	heapTemp := numberList

	color := false
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

			heapTemp, color = checkStatus(heapData, heapTemp)
			setBars(10.2, heapTemp, color)

			insertionTemp, color = checkStatus(insertionData, insertionTemp)
			setBars(6.8, insertionTemp, color)

			bubbleTemp, color = checkStatus(bubbleData, bubbleTemp)
			setBars(3.4, bubbleTemp, color)

			selectionTemp, color = checkStatus(selectionData, selectionTemp)
			setBars(0, selectionTemp, color)

			glfw.PollEvents()
			window.SwapBuffers()
		}

		timer++
	}
	close(bubbleChannel)
	close(selectionChannel)
	close(insertionChannel)
	close(heapChannel)
}

// Evalua si lo que retorna el canal es vacio, si lo es retorna la lista
// temporal y cambia color, sino retorna la lista del canal y deja el color
// Parametros:
// 		channelData = data que viene del canal
// 		tempData = data que se guardo anteriormente
func checkStatus(channelData []int, tempData []int) ([]int, bool) {
	if len(channelData) == 0 {
		return tempData, true
	} else {
		return channelData, false
	}
}

// Recorre la lista de numeros y por cada numero crea una nueva barra
// Parametros:
//		y = posicion en el eje y
//  	data = lista de enteros
//		color = bandera para saber si ya termino para pintarlo de otro color
func setBars(y float32, data []int, color bool) {
	for x := range data {
		newBar(x, y, data[x], color)
	}
}

// Crea una barra con el valor entrante
// Parametros:
// 		x = posicion en el eje x
// 		y = posicion en el eje y
// 		value = numero que representa la barra
// 		color = bandera para saber si ya termino para pintarlo de otro color`
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

// Funcion de la estructura bar para poder dibujar
func (c *bar) draw() {
	gl.ColorMask(true, c.color, false, false)
	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}

// initGlfw inicializa glfw y retorna Window para usarla
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

// initOpenGL inicializa OpenGL y retorna un programa inicializado
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	prog := gl.CreateProgram()
	gl.LinkProgram(prog)
	return prog
}

// makeVao inicializa y retorna un vertex array con los puntos de parametro
// Parametros:
// 		points = lista de numeros flotantes
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
