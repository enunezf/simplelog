package configuration

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

// Estructura que contiene la configuración desde JSON
type conf struct {
	LogFile string `json:"log_file"`
	Level   string `json:"level"`
	MaxFileLog int64 `json:"max_file_log"`
}

var (
	// Variable del tipo que contiene datos de configuración
	config conf
	// Once es un objeto que se ejecutará solo una vez
	once sync.Once
)

// GetInstance Función expuesta que lee y retorna los valores del archivo de configuración
func GetInstance() conf {
	once.Do(loadConf)
	return config
}

// loadConf lee archivo de configuración ubicado en la raiz de la aplicación
func loadConf() {
	b, err := ioutil.ReadFile("./logconf.json")
	if err != nil {
		log.Fatalf("Error al leer archivo de configuración: %s", err.Error())
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatalf("Error parsear archivo de configuración: %s", err.Error())
	}
}

