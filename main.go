package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Result Struct que va a albergar la informacion contenida en cada objeto del JSON
type Result struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

const Url = "https://xkcd.com"

func fetch(n int) (*Result, error) {

	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	// concatenar strings para obtener la url; ej: https://xkcd.com/571/info.0.json
	url := strings.Join([]string{Url, fmt.Sprintf("%d", n), "info.0.json"}, "/")

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("error con la peticion HTTP: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error HTTP: %v", err)
	}

	var data Result

	// error del web service, se emite un struct vacio para evitar la disrupcion del proceso
	if resp.StatusCode != http.StatusOK {
		data = Result{}
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("error con el objeto JSON: %v", err)
		}
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error al cerrar el cuerpo de la respuesta HTTP: %v", err)
	}

	return &data, nil
}

type Job struct {
	number int
}

// declaracion de los canales
var jobs = make(chan Job, 100)
var results = make(chan Result, 100)
var resultCollection []Result

func allocateJobs(noOfJobs int) {
	for i := 0; i <= noOfJobs; i++ {
		jobs <- Job{i + 1}
	}
	close(jobs)
}

// funcion que define el comportamiento de un worker dado el numero de trabajos a realizar
func worker(wg *sync.WaitGroup) {
	for job := range jobs {
		result, err := fetch(job.number)
		if err != nil {
			log.Printf("Error al obtener recurso: %v\n", err)
		}
		results <- *result
	}
	wg.Done()
}

// funcion que crea el pool de workers
func createWorkerPool(noOfWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i <= noOfWorkers; i++ {
		wg.Add(1)
		go worker(&wg)
	}
	wg.Wait()
	close(results)
}

// obtiene los resultados del web service
func getResults(done chan bool) {
	for result := range results {
		if result.Num != 0 {
			fmt.Printf("Recuperando comic #%d con titulo '%s'\n", result.Num, result.Title)
			resultCollection = append(resultCollection, result)
		}
	}
	done <- true
}

func main() {

	// obtiene el tiempo de inicio del programa
	start := time.Now()

	noOfJobs := 3000

	fmt.Println("-------------------------------------------------------------")
	fmt.Println("Trabajo final de Paradigmas y Lenguajes de Programacion 2022.")
	fmt.Println("Alumno: Paredes, Fernando Ivan.")
	fmt.Println("Demostracion de concurrencia en el lenguaje GO.")
	fmt.Printf("Descarga por medio de %d workers los metadatos de comics disponibles en XKCD.COM.\n", noOfJobs)
	fmt.Println("-------------------------------------------------------------")

	go allocateJobs(noOfJobs)

	// obtener resultados
	done := make(chan bool)
	go getResults(done)

	// crear pool de workers
	noOfWorkers := 100
	createWorkerPool(noOfWorkers)

	// espera a que se obtengan todos los resultados
	<-done

	// convierte el resultado de la coleccion al formato JSON
	data, err := json.MarshalIndent(resultCollection, "", "    ")
	if err != nil {
		log.Fatal("Error con el formato JSON: ", err)
	}

	// escribe el JSON a un archivo
	// cabe aclarar que no se va a poder ver el archivo si se ejecuta en el contenedor docker que proveo
	err = writeToFile(data)
	if err != nil {
		log.Fatal(err)
	}

	// obtiene el tiempo de ejecucion del programa
	elapsed := time.Since(start)

	fmt.Println("-------------------------------------------------------------")
	fmt.Printf("Numero de trabajos: %d. - Numero de trabajadores: %d.\n", noOfJobs, noOfWorkers)
	fmt.Printf("Tiempo de ejecucion: %s \n", elapsed)
	fmt.Println("Trabajo final de Paradigmas y Lenguajes de Programacion 2022.")
	fmt.Println("Alumno: Paredes, Fernando Ivan.")
	fmt.Println("Demostracion de concurrencia en el lenguaje GO.")
	fmt.Println("Muchas gracias Mgter. German Paustch.")
	fmt.Println("-------------------------------------------------------------")
}

// funcion que escribe el archivo
func writeToFile(data []byte) error {
	f, err := os.Create("xkcd.json")
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}
