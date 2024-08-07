package structures

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// MBR estructura del MBR
type MBR struct {
	Mbr_tamano   		int32
	Mbr_fecha_creacion 	[4]byte
	Mbr_disk_signature 	int32
	Dsk_fit 			[1]byte
	Mbr_particion 		[4]Partition
}



// Partition estructura de una partición
type Partition struct {
	Part_status 	[1]byte // 0 = Inactiva, 1 = Activa
	Part_type 		[1]byte // P = Primaria, E = Extendida
	Part_fit 		[1]byte // B= Best, F = First, W = Worst
	Part_start 		int32	// indica en qué byte del disco inicia la partición
	Part_size 		int32	// tamaño de la partición en bytes
	Part_name 		[16]byte	// nombre de la partición
	Part_correlative int32	// número correlativo de la partición
	Part_id		    [4]byte // identificador único de la partición
}

func ObtenerFechaHora () [4]byte{

	timestamp := time.Now().Unix()
	var buffer [4]byte
	binary.LittleEndian.PutUint32(buffer[:], uint32(timestamp))
	return buffer
}

// funcion para obtener número aleatorio, sin repetir
func obtenerNumeroAleatorio() int32 {
	num := rand.Int31()
	//almacenamos el número aleatorio en un arreglo, para verificar si se repite, teniendo en cuenta que no sabemos cuantos números aleatorios se van a generar
	var numeros [100]int32
	//recorremos el arreglo de números aleatorios
	for i := 0; i < len(numeros); i++ {
		//verificamos si el número aleatorio generado ya existe en el arreglo
		if numeros[i] == num {
			//si el número aleatorio ya existe, generamos un nuevo número aleatorio
			num = rand.Int31()
			//reiniciamos el ciclo
			i = 0
		}
	}
	//retornamos el número aleatorio
	return num
}


// CrearMBR crea un MBR con los valores inicializados
func CrearMBR(tamano int32, fit string) MBR {
	fechaHora := ObtenerFechaHora()
	numeroAleatorio := obtenerNumeroAleatorio()
	
	var fitArray [1]byte
	copy(fitArray[:], fit)

	mbr := MBR{
		Mbr_tamano: tamano,
		Mbr_fecha_creacion: fechaHora,
		Mbr_disk_signature: numeroAleatorio,
		Dsk_fit: fitArray,

	}

	return mbr
}

// ReadMBRFromFile lee el MBR desde el archivo binario y lo retorna
func ReadMBRFromFile(path string) (MBR, error) {
	file, err := os.Open(path)
	if err != nil {
		return MBR{}, err // Si no se pudo abrir el archivo, retornar el error
	}
	defer file.Close() // Cerrar el archivo

	var mbr MBR
	err = binary.Read(file, binary.LittleEndian, &mbr)
	if err != nil {
		return MBR{}, err // Si no se pudo leer el MBR, retornar el error
	}

	return mbr, nil
}

// ConvertirFecha convierte el array de 4 bytes a una fecha legible
func ConvertirFecha(fechaBytes [4]byte) string {
	timestamp := binary.LittleEndian.Uint32(fechaBytes[:])
	t := time.Unix(int64(timestamp), 0)
	return t.Format("02-01-2006 15:04:05")
}

// PrintMBR imprime la información del MBR
func PrintMBR(mbr MBR) {
	fmt.Println("MBR Info:")
	fmt.Printf("Tamaño del disco: %d bytes\n", mbr.Mbr_tamano)
	fmt.Printf("Fecha de creación: %s\n", ConvertirFecha(mbr.Mbr_fecha_creacion))
	fmt.Printf("Disk Signature: %d\n", mbr.Mbr_disk_signature)
	fmt.Printf("Fit: %s\n", string(mbr.Dsk_fit[:]))
	for i, part := range mbr.Mbr_particion {
		fmt.Printf("Partición %d:\n", i+1)
		fmt.Printf("  Status: %s\n", string(part.Part_status[:]))
		fmt.Printf("  Type: %s\n", string(part.Part_type[:]))
		fmt.Printf("  Fit: %s\n", string(part.Part_fit[:]))
		fmt.Printf("  Start: %d\n", part.Part_start)
		fmt.Printf("  Size: %d\n", part.Part_size)
		fmt.Printf("  Name: %s\n", string(part.Part_name[:]))
		fmt.Printf("  Correlative: %d\n", part.Part_correlative)
		fmt.Printf("  ID: %s\n", string(part.Part_id[:]))
	}
}
