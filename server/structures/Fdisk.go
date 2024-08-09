package structures


// FDISK estructura que representa el comando fdisk con sus parámetros
type FDISK struct {
	Size int    // Tamaño de la partición
	Unit string // Unidad de medida del tamaño (K o M)
	Fit  string // Tipo de ajuste (BF, FF, WF)
	Path string // Ruta del archivo del disco
	TypE  string // Tipo de partición (P, E, L)
	Name string // Nombre de la partición
}


