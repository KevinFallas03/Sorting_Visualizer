package visualizer

import (
	"log"
	"runtime"

	"../algorithms"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/nullboundary/glfont"
)

const (
	width  = 1280
	height = 700
	rows   = 200
)

var (
	columns   = 0
	rectangle = []float32{
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

//generateList crea una lista de numeros aleatorios con el metodo de
//congruncia lineal multiplicativa
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

//Start ...
func Start(n int, x int, m int, msgCh chan string) {
	//GENERA LA LISTA DE NUMEROS
	numberList := generateList(n, x, m)
	columns = len(numberList) + int(float32(len(numberList))*0.05)

	//GENERA DATA PARA LOS ALGORITMOS
	var numberLists [][]int       //Lista de listas de numeros
	var tempLists [][]int         //Lista de listas temporales
	var actualLists [][]int       //Lista de listas actualizadas
	var channelList []chan []int  //Lista de canales
	stopCh := make(chan struct{}) //Canal para detener todo

	//INICIALIZA TODOS LOS DATOS
	for i := 0; i < 6; i++ {
		newList := make([]int, len(numberList), len(numberList))
		copy(newList, numberList)
		numberLists = append(numberLists, newList)
		channelList = append(channelList, make(chan []int))
		tempLists = append(tempLists, numberList)
		actualLists = append(actualLists, numberList)
	}

	//INICIA CADA ALGORITMO CON CORRUTINAS
	go algorithms.HeapSort(numberLists[0], channelList[0], stopCh, msgCh)
	go algorithms.QuickSort(numberLists[1], channelList[1], stopCh, msgCh)
	go algorithms.MergeSort(numberLists[2], channelList[2], stopCh, msgCh)
	go algorithms.InsertionSort(numberLists[3], channelList[3], stopCh, msgCh)
	go algorithms.SelectionSort(numberLists[4], channelList[4], stopCh, msgCh)
	go algorithms.BubbleSort(numberLists[5], channelList[5], stopCh, msgCh)

	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window := initGlfw()
	initOpenGL()

	font, err := glfont.LoadFont("Roboto-Light.ttf", int32(52), width, height)
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}

	color := false
	timer := 0
	percentage := float32(columns) * 0.01
	if percentage < 1 {
		percentage = 1
	}
	algorithmsName := [6]string{"BubbleSort", "SelectionSort", "InsertionSort", "MergeSort", "QuickSort", "HeapSort"}
	//load font (fontfile, font scale, window width, window height
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		//OBTIENE INFORMACION DE LOS CANALES
		for data := 0; data < len(channelList); data++ {
			actualLists[data] = <-channelList[data]
		}

		//CADA CIERTO TIEMPO PINTA
		if timer%int(percentage) == 0 {
			for data := 0; data < len(channelList); data++ {
				tempLists[data], color = checkStatus(actualLists[data], tempLists[data])
				if data < 3 {
					font.Printf(100, (float32(data)+0.7)*120, 1.2, algorithmsName[data]) //x,y,scale,string,printf args
					setBars(3.4*float32(data), tempLists[data], color, false)
				} else {
					font.Printf(800, (float32(data)+0.7)*120, 1.2, algorithmsName[data]) //x,y,scale,string,printf args
					setBars(3.4*float32(data), tempLists[data], color, true)
				}
			}
			window.SwapBuffers()
			glfw.PollEvents()

		}
		timer++
	}
	close(stopCh) //Cerrando este canal cerramos los demas canales en cada algoritmo
	close(msgCh)
}
func initGlfw() *glfw.Window {
	window, _ := glfw.CreateWindow(int(width), int(height), "glfontExample", nil, nil)
	window.MakeContextCurrent()
	return window
}
func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
}
func checkStatus(channelData []int, tempData []int) ([]int, bool) {
	if len(channelData) == 0 {
		return tempData, true
	} else {
		return channelData, false
	}
}
func setBars(y float32, data []int, color bool, lado bool) {

	for x := range data {
		newBar(x, y, data[x], color, lado)
	}
}
func newBar(x int, y float32, value int, color bool, izqDer bool) {
	points := make([]float32, len(rectangle), len(rectangle))
	copy(points, rectangle)

	for i := 0; i < len(points); i++ {
		var position, size, m float32
		switch i % 3 {
		case 0:
			size = (2.0 / float32(columns)) / 2
			position = float32(x) * size / 2 // POSITION
			if izqDer == false {
				m = 0
			} else {
				m = 1

			}

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
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(rectangle)/3))
}
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
