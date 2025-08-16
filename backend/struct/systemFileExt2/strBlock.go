package systemfileext2

/*
	Un bloque es la unidad minima de almacenamiento a nivel logico,
	son un conjunto de sectores contiguos en el disco.

	En el proyecto indican que su tamaño es de 64 bytes, los
	tipos de bloques son:

	- Bloques de Carpetas: Guardan información sobre los nombres de
	  los archivos que contienen y a que Inodo apuntan.
	- Bloques de Archivos: Guardan la información del archivo,
	  incluyendo su contenido y metadatos.
	- Bloques de apuntadores: Guardan información sobre los apuntadores
	  indirectos (simples, dobles y triples).
*/

// Bloque de Carpetas
// Content representa el contenido de un bloque carpeta
type Content struct {
	Bname  [12]byte `binary:"little"` // Nombre de la carpeta o archivo
	BInodo int32    `binary:"little"` // Apunta hacia un inodo asociado
}
