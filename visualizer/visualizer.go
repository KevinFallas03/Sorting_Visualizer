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
	value    int
	index    float32
}

type graph struct {
	bars      []*bar
	color     []bool
	yPosition float32
	lado      bool
	name      string
	dataCh    chan []int
	area      []int32
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

	//GENERA LA LISTA DE NUMEROS ALEATORIOS
	numberList := generateList(n, x, m)
	columns = len(numberList) + int(float32(len(numberList))*0.05)

	//GENERA DATA PARA LOS ALGORITMOS
	var numberLists [][]int      //Lista de listas de numeros
	var tempLists [][]int        //Lista de listas temporales
	var actualLists [][]int      //Lista de listas actualizadas
	var channelList []chan []int //Lista de canales
	stopCh := make(chan int)     //Canal para detener todo
	var graphList []*graph       //Lista de graficos

	//INICIALIZA TODOS LOS DATOS
	for i := 0; i < 6; i++ {
		newList := make([]int, len(numberList), len(numberList))
		copy(newList, numberList)
		numberLists = append(numberLists, newList)
		channelList = append(channelList, make(chan []int))
		tempLists = append(tempLists, numberList)
		actualLists = append(actualLists, numberList)
		graphList = append(graphList, &graph{dataCh: channelList[i]})
	}

	go graphList[0].drawGraph()
	//INICIA CADA ALGORITMO CON CORRUTINAS
	go algorithms.HeapSort(numberLists[0], channelList[0], stopCh, msgCh)
	// go algorithms.QuickSort(numberLists[1], channelList[1], stopCh, msgCh)
	// go algorithms.MergeSort(numberLists[2], channelList[2], stopCh, msgCh)
	// go algorithms.InsertionSort(numberLists[3], channelList[3], stopCh, msgCh)
	// go algorithms.SelectionSort(numberLists[4], channelList[4], stopCh, msgCh)
	// go algorithms.BubbleSort(numberLists[5], channelList[5], stopCh, msgCh)

	//INICIA LA VENTANA
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

	//INICIA LOS GRAFICOS
	font, _ = glfont.LoadFont("Roboto-Light.ttf", int32(52), width, height)
	color := []bool{true, false, true}
	algorithmsName := [6]string{"BubbleSort", "SelectionSort", "InsertionSort", "MergeSort", "QuickSort", "HeapSort"}

	for i := 0; i < 6; i++ {
		lado := false
		x := 100
		xI := 0
		if i > 2 {
			lado = true
			x = 800
			xI = 640
		}
		graphList[i].setData(3.4*float32(i), numberLists[i], color, lado, algorithmsName[i])
		//newGraph := createGraph(3.4*float32(i), numberLists[i], color, lado, algorithmsName[i])
		//graphList = append(graphList, newGraph)
		//graphList[i].dataCh = channelList[i]
		graphList[i].area = []int32{int32(xI), int32(117 * i), 640, 117}

		font.Printf(float32(x), (float32(i)+0.7)*120, 1.2, algorithmsName[i])
		window.SwapBuffers()
		font.Printf(float32(x), (float32(i)+0.7)*120, 1.2, algorithmsName[i])
	}

	gl.Enable(gl.SCISSOR_TEST)
	for !window.ShouldClose() {
		log.Println("Ventana")
		stopCh <- 2
		// select {
		// case actualLists[0] = <-channelList[0]:
		// 	gl.Scissor(0, 0, 640, 117)
		// 	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// 	tempLists[0], graphList[0].color = checkStatus(actualLists[0], tempLists[0])
		// 	graphList[0].updateGraph(tempLists[0])
		// 	graphList[0].drawGraph()
		// case actualLists[1] = <-channelList[1]:
		// 	gl.Scissor(0, 117, 640, 117)
		// 	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// 	tempLists[1], graphList[0].color = checkStatus(actualLists[1], tempLists[1])
		// 	graphList[1].updateGraph(tempLists[1])
		// 	graphList[1].drawGraph()
		// case actualLists[2] = <-channelList[2]:
		// 	gl.Scissor(0, 234, 640, 117)
		// 	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// 	tempLists[2], graphList[0].color = checkStatus(actualLists[2], tempLists[2])
		// 	graphList[2].updateGraph(tempLists[2])
		// 	graphList[2].drawGraph()
		// case actualLists[3] = <-channelList[3]:
		// 	gl.Scissor(640, 351, 640, 117)
		// 	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// 	tempLists[3], graphList[0].color = checkStatus(actualLists[3], tempLists[3])
		// 	graphList[3].updateGraph(tempLists[3])
		// 	graphList[3].drawGraph()
		// case actualLists[4] = <-channelList[4]:
		// 	gl.Scissor(640, 468, 640, 117)
		// 	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// 	tempLists[4], graphList[0].color = checkStatus(actualLists[4], tempLists[4])
		// 	graphList[4].updateGraph(tempLists[4])
		// 	graphList[4].drawGraph()
		// case actualLists[5] = <-channelList[5]:
		// 	gl.Scissor(640, 585, 640, 117)
		// 	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// 	tempLists[5], graphList[0].color = checkStatus(actualLists[5], tempLists[5])
		// 	graphList[5].updateGraph(tempLists[5])
		// 	graphList[5].drawGraph()
		// }

		// gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// for data := 0; data < len(channelList); data++ {
		// 	actualLists[data] = <-channelList[data]
		// 	tempLists[data], graphList[0].color = checkStatus(actualLists[data], tempLists[data])
		// 	graphList[data].updateGraph(tempLists[data])
		// 	graphList[data].drawGraph()
		// }
		glfw.PollEvents()
		window.SwapBuffers()
	}
	close(stopCh) //Cerrando este canal cerramos los demas canales en cada algoritmo
	close(msgCh)
}

func checkStatus(channelData []int, tempData []int) ([]int, []bool) {
	if len(channelData) == 0 {
		return tempData, []bool{false, true, false}
	} else {
		return channelData, []bool{true, false, true}
	}
}

//======================FUNCIONES DEl GRAFICO======================================

func (g *graph) setData(yPos float32, data []int, color []bool, lado bool, algorithmName string) {
	var newBars []*bar
	for i := 0; i < len(data); i++ {
		newB := createBar(float32(i), yPos, data[i], color, lado)
		newBars = append(newBars, newB)
	}

	g.yPosition = yPos
	g.color = color
	g.bars = newBars
	g.lado = lado
	g.name = algorithmName
	//return newGraph
}
func (g *graph) drawGraph() {
	for {
		log.Println("pintado")
		select {
		case currentList := <-g.dataCh:
			log.Println("Llego data")
			if len(currentList) == 0 {
				return
			} else {
				gl.Scissor(g.area[0], g.area[1], g.area[2], g.area[3])
				gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
				g.updateGraph(currentList)

				for i := 0; i < len(g.bars); i++ {
					g.bars[i].drawBar()
				}
			}
		}
	}
}
func (g *graph) updateGraph(data []int) {

	//UPDATE EACH BAR: va a funcionar cuando los algoritmos retornen solo un elemento o indice
	for i := 0; i < len(data); i++ {
		g.bars[i].setDrawable(float32(i), g.yPosition, data[i], g.lado)
		g.bars[i].index = float32(i)
		g.bars[i].value = data[i]
	}
}

//======================FUNCIONES DE LA BARRA======================================

func createBar(x, y float32, value int, color []bool, lado bool) *bar {
	bar := bar{
		color: color,
		value: value,
		index: x,
	}
	bar.setDrawable(x, y, value, lado)
	return &bar
}

//Antes newBar
func (c *bar) setDrawable(x, y float32, value int, izqDer bool) {
	points := make([]float32, len(rectangle), len(rectangle))
	copy(points, rectangle)

	for i := 0; i < len(points); i++ {
		var position, size, m float32
		switch i % 3 {
		case 0:
			size = (2.0 / float32(columns)) / 2
			position = x * size / 2 // POSITION
			m = 1
			if izqDer == false {
				m = 0
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
	c.drawable = makeVao(points)
}
func (c *bar) drawBar() {
	gl.ColorMask(true, false, true, false)
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

//======================FUNCIONES DE LA VENTANA======================================

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
