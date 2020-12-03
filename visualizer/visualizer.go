package visualizer

import (
	"runtime"
	"../algorithms"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 1280
	height = 700
	rows   = 200
)

var (
	columns = 0
	square  = []float32{
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

func generateList(n int, x int, m int) []int {

	//  x es la semilla. Debe de ser primo entre [11, 101]  0 <= x < m
	//  n es la cantidad.
	//  1103515245, 12345,
	// var m int = 2048       // >= 2048   Periodo
	var a int = 1103515245 //    0 < a < m multiplicador
	var c int = 12345      //      0 <= c < m  Incremento

	var nums []int
	for i := 0; i < n; i++ {
		x = (a*x + c) % m
		nums = append(nums, x%31)
	}
	return nums
}

func Start(n int, x int, m int) {
	//GENERA LA LISTA DE NUMEROS
	numberList := generateList(n, x, m)
	columns = len(numberList) + int(float32(len(numberList))*0.05)

	//GENERA DATA PARA LOS ALGORITMOS
	var numberLists [][]int      //Lista de listas de numeros
	var tempLists [][]int        //Lista de listas temporales
	var actualLists [][]int      //Lista de listas actualizadas
	var channelList []chan []int //Lista de canales
	stopCh := make(chan struct{}) 
	for i := 0; i < 5; i++ {
		newList := make([]int, len(numberList), len(numberList))
		copy(newList, numberList)
		numberLists = append(numberLists, newList)
		channelList = append(channelList, make(chan []int))
		tempLists = append(tempLists, numberList)
		actualLists = append(actualLists, numberList)
	}

	//INICIA CADA ALGORITMO
	go algorithms.HeapSort(numberLists[0], channelList[0],stopCh)
	go algorithms.InsertionSort(numberLists[1], channelList[1],stopCh)
	go algorithms.SelectionSort(numberLists[2], channelList[2],stopCh)
	go algorithms.BubbleSort(numberLists[3], channelList[3],stopCh)
	go algorithms.QuickSort(numberLists[4], channelList[4],stopCh)

	//MOSTRAR VENTANA
	runtime.LockOSThread()
	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	color := false
	timer := 0
	percentage := float32(columns) * 0.01

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		//OBTIENE INFORMACION DE LOS CANALES
		for data := 0; data < len(channelList); data++ {
			actualLists[data] = <-channelList[data]
		}

		//CADA CIERTO TIEMPO PINTA
		if timer%int(percentage) == 0 {
			for data := 0; data < len(channelList); data++ {
				tempLists[data], color = checkStatus(actualLists[data], tempLists[data])
				setBars(3.4*float32(data), tempLists[data], color)
			}
			glfw.PollEvents()
			window.SwapBuffers()
		}
		timer++
	}
	close(stopCh)
	// for data := 0; data < len(channelList); data++ {
	// 	close(channelList[data])
	// }
}

// Evalua si lo que retorna el canal es vacio
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
		var position, size, m float32
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
