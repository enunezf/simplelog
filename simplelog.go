package simplelog

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/enunezf/simplelog/configuration"
)

const maxfilelog = 2 * (1024 * 1024) // 2MB

var (
	// Trace logger para textos de trace
	Trace *log.Logger
	// Info logger para textos de info
	Info *log.Logger
	// Warning logger para textos de warning
	Warning *log.Logger
	// Error logger para textos de error
	Error *log.Logger
	c     = configuration.GetInstance()
)

// Inicializa package crea archivo log e instancia los logger
func init() {
	// Abrir o crear archivo log
	fileLog, err := os.OpenFile(c.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error al abrir archivo log %s", err.Error())
	}

	// Instacia de logger
	Trace = log.New(fileLog, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(fileLog, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(fileLog, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(fileLog, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// gorutina que permite evaluar de forma asincrónica tamaño de archivo y
	// archivarlo en caso que supere el tamaño máximo (maxfilelog)
	go func() {
		for {
			fi, err := fileLog.Stat()
			if err != nil {
				log.Printf("Error al leer estadísticas del archivo %s: %s", fileLog.Name(), err.Error())
			}
			if fi.Size() >= maxfilelog {
				archiveFile(fileLog)
			}

			time.Sleep(time.Minute * 5)
		}

	}()

}

// archiveFile Crear un log de backup y limpia archivo original
func archiveFile(file *os.File) {
	// Obtiene fecha actual
	today := time.Now()
	// crea nombre del archivo
	filename := fmt.Sprintf("%s.%d%d%d-%d%d%d.log", file.Name(), today.Year(), today.Month(), today.Day(), today.Hour(), today.Minute(), today.Second())

	// Lo abre para escritura
	af, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error al crear archivo %s: %s", filename, err.Error())
	}
	// Cuando termine la función se cierra el archivo
	defer af.Close()

	// Lee archivo pasado por parámetro
	fileContent, err := ioutil.ReadFile(file.Name())
	if err != nil {
		log.Printf("Error al leer contenido de archivo %s: %s", file.Name(), err.Error())
	}

	// Escribe contenido en archivo de respaldo
	af.Write(fileContent)

	// Limpia archivo original.
	err = ioutil.WriteFile(file.Name(), []byte("Se creó el archivo"+filename+"\n"), 0666)
	if err != nil {
		log.Printf("No se pudo limpiar el archivo %s: %s", file.Name(), err.Error())
	}

}

