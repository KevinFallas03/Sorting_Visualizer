package main

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
	window *glfw.Window
	font   *glfont.Font
)

type bar struct {
	drawable uint32
	color    []bool
}

type graph struct {
	bars      []bar
	color     []bool
	yPosition float32
	lado      bool
	name      string
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
func main() {
	n := 200
	x := 101
	m := 2048
	msgCh := make(chan string)

	//GENERA LA LISTA DE NUMEROS
	numberList := generateList(n, x, m)
	columns = len(numberList) + int(float32(len(numberList))*0.05)

	//GENERA DATA PARA LOS ALGORITMOS
	var numberLists [][]int      //Lista de listas de numeros
	var tempLists [][]int        //Lista de listas temporales
	var actualLists [][]int      //Lista de listas actualizadas
	var channelList []chan []int //Lista de canales
	//var graphList []*graph        //Lista de graficos
	stopCh := make(chan struct{}) //Canal para detener todo
	color := []bool{false, true, false}
	algorithmsName := [6]string{"BubbleSort", "SelectionSort", "InsertionSort", "MergeSort", "QuickSort", "HeapSort"}

	//INICIALIZA TODOS LOS DATOS
	for i := 0; i < 6; i++ {
		newList := make([]int, len(numberList), len(numberList))
		copy(newList, numberList)
		numberLists = append(numberLists, newList)
		channelList = append(channelList, make(chan []int))
		tempLists = append(tempLists, numberList)
		actualLists = append(actualLists, numberList)
		// lado := true
		// if i > 2 {
		// 	lado = false
		// }
		// newGraph := createGraph(3.4*float32(i), newList, color, lado, algorithmsName[i])
		// graphList = append(graphList, newGraph)
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

	window = initGlfw()
	initOpenGL()

	font, _ = glfont.LoadFont("Roboto-Light.ttf", int32(52), width, height)

	for !window.ShouldClose() {

		//gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		select {
		case actualLists[0] = <-channelList[0]:
			tempLists[0], color = checkStatus(actualLists[0], tempLists[0])
			// graphList[0].updateBarsGraph(tempLists[0])
			// graphList[0].drawGraph()
			drawGraph(3.4*float32(0), tempLists[0], color, false, algorithmsName[0], 100)
		case actualLists[1] = <-channelList[1]:
			tempLists[1], color = checkStatus(actualLists[1], tempLists[1])
			// graphList[1].updateBarsGraph(tempLists[1])
			// graphList[1].drawGraph()
			drawGraph(3.4*float32(1), tempLists[1], color, false, algorithmsName[1], 100)
		case actualLists[2] = <-channelList[2]:
			tempLists[2], color = checkStatus(actualLists[2], tempLists[2])
			// graphList[2].updateBarsGraph(tempLists[2])
			// graphList[2].drawGraph()
			drawGraph(3.4*float32(2), tempLists[2], color, false, algorithmsName[2], 100)
		case actualLists[3] = <-channelList[3]:
			tempLists[3], color = checkStatus(actualLists[3], tempLists[3])
			// graphList[3].updateBarsGraph(tempLists[3])
			// graphList[3].drawGraph()
			drawGraph(3.4*float32(3), tempLists[3], color, true, algorithmsName[3], 800)
		case actualLists[4] = <-channelList[4]:
			tempLists[4], color = checkStatus(actualLists[4], tempLists[4])
			// graphList[4].updateBarsGraph(tempLists[4])
			// graphList[4].drawGraph()
			drawGraph(3.4*float32(4), tempLists[4], color, true, algorithmsName[4], 800)
		case actualLists[5] = <-channelList[5]:
			tempLists[5], color = checkStatus(actualLists[5], tempLists[5])
			// graphList[5].updateBarsGraph(tempLists[5])
			// graphList[5].drawGraph()
			drawGraph(3.4*float32(5), tempLists[5], color, true, algorithmsName[5], 800)
		}
		// for data := 0; data < len(channelList); data++ {
		// 	actualLists[data] = <-channelList[data]
		// 	tempLists[data], color = checkStatus(actualLists[data], tempLists[data])
		// 	if data < 3 {
		// 		drawGraph(3.4*float32(data), tempLists[data], color, false, algorithmsName[data], 100)
		// 	} else {
		// 		drawGraph(3.4*float32(data), tempLists[data], color, true, algorithmsName[data], 800)
		// 	}
		// }
		glfw.PollEvents()
		window.SwapBuffers()
	}
	close(stopCh) //Cerrando este canal cerramos los demas canales en cada algoritmo
	close(msgCh)
}

func drawGraph(y float32, data []int, color []bool, lado bool, name string, x float32) {
	setBars(y, data, color, lado)
	font.Printf(x, ((y/3.4)+0.7)*120, 1.2, name) //x,y,scale,string,printf args
}

func checkStatus(channelData []int, tempData []int) ([]int, []bool) {
	if len(channelData) == 0 {
		return tempData, []bool{true, true, true}
	} else {
		return channelData, []bool{true, false, false}
	}
}
func setBars(y float32, data []int, color []bool, lado bool) {
	for x := range data {
		newBar(x, y, data[x], color, lado)
	}
}
func newBar(x int, y float32, value int, color []bool, izqDer bool) bar {
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

	bar := bar{
		drawable: makeVao(points),
		color:    color,
	}

	bar.draw()
	return bar
}

//FUNCIONES DE LA BARRA
func (c *bar) draw() {
	gl.ColorMask(c.color[0], c.color[1], c.color[2], false)
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
func initGlfw() *glfw.Window {
	window, _ := glfw.CreateWindow(int(width), int(height), "THE BEST SORT VISUALIZER", nil, nil)
	window.MakeContextCurrent()
	return window
}
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	prog := gl.CreateProgram()
	gl.LinkProgram(prog)
	return prog
}
