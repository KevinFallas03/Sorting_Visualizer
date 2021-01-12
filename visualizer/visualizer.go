package main

import (
	"fmt"
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
	window      *glfw.Window
	font        *glfont.Font
	channelList []chan []int //Lista de canales
	graphList   []*graph     //Lista de graficos
)

type bar struct {
	drawable uint32
	color    []bool
	value    int
	index    float32
}

type graph struct {
	bars      []bar
	color     []bool
	yPosition float32
	lado      bool
	name      string
	done      bool
}

//generateList crea una lista de numeros aleatorios con el metodo de
//congruncia lineal multiplicativa
func generateList(n int, x int, m int) [][]int {
	//  x es la semilla. Debe de ser primo entre [11, 101]  0 <= x < m
	//  n es la cantidad.
	//  1103515245, 12345,
	// var m int = 2048       // >= 2048   Periodo
	var a int = 1103515245 //    0 < a < m multiplicador
	var c int = 12345      //      0 <= c < m  Incremento

	var nums [][]int
	for i := 0; i < n; i++ {
		x = (a*x + c) % m
		nums = append(nums, []int{x % 31, i})
	}
	return nums
}

//Start ...
func main() {
	n := 100
	x := 101
	m := 2048
	msgCh := make(chan string)

	//GENERA LA LISTA DE NUMEROS ALEATORIOS
	numberList := generateList(n, x, m)
	columns = len(numberList) + int(float32(len(numberList))*0.05)

	//GENERA DATA PARA LOS ALGORITMOS
	var numberLists [][][]int     //Lista de listas de numeros
	stopCh := make(chan struct{}) //Canal para detener todo

	color := []bool{true, false, true}
	algorithmsName := [6]string{"BubbleSort", "SelectionSort", "InsertionSort", "MergeSort", "QuickSort", "HeapSort"}

	//INICIALIZA TODOS LOS DATOS
	for i := 0; i < 6; i++ {
		newList := make([][]int, len(numberList), len(numberList))
		copy(newList, numberList)
		numberLists = append(numberLists, newList)
		channelList = append(channelList, make(chan []int))
	}

	//INICIA CADA ALGORITMO CON CORRUTINAS
	//go algorithms.HeapSort(numberLists[0], channelList[0], stopCh, msgCh)
	// go algorithms.QuickSort(numberLists[1], channelList[1], stopCh, msgCh)
	// go algorithms.MergeSort(numberLists[2], channelList[2], stopCh, msgCh)
	// go algorithms.InsertionSort(numberLists[3], channelList[3], stopCh, msgCh)
	// go algorithms.SelectionSort(numberLists[4], channelList[4], stopCh, msgCh)
	go algorithms.BubbleSort(numberLists[5], channelList[5], stopCh, msgCh)

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

	//DIBUJA LAS ETIQUETAS
	font, _ = glfont.LoadFont("Roboto-Light.ttf", int32(52), width, height)
	for i := 0; i < 6; i++ {
		x := 100
		lado := false
		if i > 2 {
			lado = true
			x = 800
		}
		newGraph := createGraph(3.4*float32(i), numberLists[i], color, lado, algorithmsName[i])
		graphList = append(graphList, newGraph)

		font.Printf(float32(x), (float32(i)+0.7)*120, 1.2, algorithmsName[i])
		window.SwapBuffers()
		font.Printf(float32(x), (float32(i)+0.7)*120, 1.2, algorithmsName[i])
	}

	gl.Enable(gl.SCISSOR_TEST)

	for !window.ShouldClose() {
		select {
		case currentListHP := <-channelList[0]: //HeapSort
			drawInWindow(0, 0, currentListHP, 0)
		case currentListQS := <-channelList[1]: //QuickSort
			drawInWindow(0, 117, currentListQS, 1)
		case currentListMS := <-channelList[2]: //MergeSort
			drawInWindow(0, 234, currentListMS, 2)
		case currentListIS := <-channelList[3]: //InsertionSort
			drawInWindow(640, 351, currentListIS, 3)
		case currentListSS := <-channelList[4]: //SelectionSort
			drawInWindow(640, 468, currentListSS, 4)
		case currentListBS := <-channelList[5]: //BubbleSort
			drawInWindow(640, 585, currentListBS, 5)
		}

		// gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// for data := 0; data < len(channelList); data++ {
		// 	actualLists[data] = <-channelList[data]
		// 	tempLists[data], graphList[0].color = checkStatus(actualLists[data], tempLists[data])
		// 	graphList[data].updateGraph(tempLists[data])
		// 	graphList[data].drawGraph()
		// }

		glfw.PollEvents()

	}
	close(stopCh) //Cerrando este canal cerramos los demas canales en cada algoritmo
	close(msgCh)
}
func drawInWindow(xCut, yCut int32, currentList []int, index int) {
	if !graphList[index].done { //Si el grafico no se ha terminado de pintar
		gl.Scissor(xCut, yCut, 640, 117)                    //Seleccionamos la parte de la ventana que queremos actualizar
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) //La limpiamos
		if len(currentList) == 0 {                          //Si el algoritmo termino
			graphList[index].done = true                             //El grafico se pinto
			graphList[index].updateColor([]bool{false, true, false}) //Actualiza el color del grafico
		} else { //Si el algoritmo no ha terminado
			graphList[index].updateGraph(currentList) //Actualiza las barras
		}
		//Pinta el grafico en ambos buffers
		graphList[index].drawGraph()
		window.SwapBuffers()
	}
}

//======================FUNCIONES DEl GRAFICO======================================

// Crea un grafico nuevo con todos sus atributos.
// (Posicion en y, color, lista de barras que lo componen, lado de la pantalla donde se ubica, nombre del algoritmo que representa)
func createGraph(yPos float32, data [][]int, color []bool, lado bool, algorithmName string) *graph {
	var newBars []bar
	for i := 0; i < len(data); i++ {
		newB := createBar(float32(i), yPos, data[i][0], color, lado)
		newBars = append(newBars, newB)
	}

	newGraph := &graph{
		yPosition: yPos,
		color:     color,
		bars:      newBars,
		lado:      lado,
		name:      algorithmName,
		done:      false,
	}
	return newGraph
}

// Dibuja el grafico.
// Agarra cada barra del grafico y la pinta.
func (g *graph) drawGraph() {
	for i := 0; i < len(g.bars); i++ {
		g.bars[i].drawBar()
	}
}

// Actualiza las barras:
// V0: Actualiza la lista entera con la lista que recibe del canal del algoritmo.
func (g *graph) updateGraph(data []int) {

	// //UPDATE EACH BAR: va a funcionar cuando los algoritmos retornen solo un elemento o indice
	// for i := 0; i < len(data); i++ {
	// 	g.bars[i].setDrawable(float32(i), g.yPosition, data[i], g.lado)
	// 	g.bars[i].index = float32(i)
	// 	g.bars[i].value = data[i]
	// }
	//UPDATE EACH BAR: va a funcionar cuando los algoritmos retornen solo un elemento o indice

	fmt.Println("indices->", data)
	// fmt.Println("barras->", g.bars)
	// var toSwap []int
	// for i := 0; i < len(g.bars); i++ {
	// 	// g.bars[i].setDrawable(float32(i), g.yPosition, data[i], g.lado)
	// 	// g.bars[i].index = float32(i)
	// 	// g.bars[i].value = data[i]
	// 	if int(g.bars[i].index) == data[0] || int(g.bars[i].index) == data[1] {
	// 		toSwap = append(toSwap, i)
	// 	}
	// }
	// g.bars[toSwap[0]].setDrawable(g.bars[toSwap[1]].index, g.yPosition, g.bars[toSwap[1]].value, g.lado)
	// g.bars[toSwap[1]].setDrawable(g.bars[toSwap[0]].index, g.yPosition, g.bars[toSwap[0]].value, g.lado)

	g.bars[data[0]].drawable, g.bars[data[1]].drawable = g.bars[data[1]].drawable, g.bars[data[0]].drawable
}

//Actualiza los colores de las barras respecto a la del grafico
func (g *graph) updateColor(color []bool) {
	g.color = color
	for i := 0; i < len(g.bars); i++ {
		g.bars[i].color = g.color
	}
}

// Genera los "drawables" de cada barra: se hace por separado y no a la hora de crear la barra
// porque se necesita hacer despues de crear la ventana.
func (g *graph) setDrawables() {
	for i := 0; i < len(g.bars); i++ {
		g.bars[i].setDrawable(g.bars[i].index, g.yPosition, g.bars[i].value, g.lado)
	}
}

//======================FUNCIONES DE LA BARRA======================================

//Crea una barra con todos sus atributos: (color, valor, indice en la lista).
func createBar(x, y float32, value int, color []bool, lado bool) bar {
	bar := bar{
		color: color,
		value: value,
		index: x,
	}
	bar.setDrawable(x, y, value, lado)
	return bar
}

//Genera y establece el drawable para la barra, es el objeto que se pintara.
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

//Dibuja la barra en la pantalla
func (c *bar) drawBar() {
	gl.ColorMask(c.color[0], c.color[1], c.color[2], false)
	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(rectangle)/3))
}

//Genera el objeto VAO con la liberia OpenGL
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
