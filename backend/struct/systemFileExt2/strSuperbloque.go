package systemfileext2

/*
	┌───────────────────────┬────────┬───────────────────────────────────────────────────────────────┐
	│ NOMBRE                │ TIPO   │ DESCRIPCIÓN                                                   │
	├───────────────────────┼────────┼───────────────────────────────────────────────────────────────┤
	│ s_filesystem_type     │ int    │ Guarda el número que identifica el sistema de archivos (2)    │
	│ s_inodes_count        │ int    │ Guarda el número total de inodos                              │
	│ s_blocks_count        │ int    │ Guarda el número total de bloques                             │
	│ s_free_blocks_count   │ int    │ Contiene el número de bloques libres                          │
	│ s_free_inodes_count   │ int    │ Contiene el número de inodos libres                           │
	│ s_mtime               │ time   │ Última fecha en la que el sistema fue montado                 │
	│ s_umtime              │ time   │ Última fecha en que el sistema fue desmontado                 │
	│ s_mnt_count           │ int    │ Indica cuántas veces se ha montado el sistema                 │
	│ s_magic               │ int    │ Identificador del sistema de archivos (0xEF53)                │
	│ s_inode_s             │ int    │ Tamaño del inodo                                              │
	│ s_block_s             │ int    │ Tamaño del bloque                                             │
	│ s_firts_ino           │ int    │ Primer inodo libre (dirección del inodo)                      │
	│ s_first_blo           │ int    │ Primer bloque libre (dirección del bloque)                    │
	│ s_bm_inode_start      │ int    │ Inicio del bitmap de inodos                                   │
	│ s_bm_block_start      │ int    │ Inicio del bitmap de bloques                                  │
	│ s_inode_start         │ int    │ Inicio de la tabla de inodos                                  │
	│ s_block_start         │ int    │ Inicio de la tabla de bloques                                 │
	└───────────────────────┴────────┴───────────────────────────────────────────────────────────────┘
*/

// Superblock contiene información sobre el sistema de archivos EXT2
type Superblock struct {
	SFilesystemType  int32 `binary:"little"` // Número que identifica el sistema de archivos utilizado (2)
	SInodesCount     int32 `binary:"little"` // Número total de inodos
	SBlocksCount     int32 `binary:"little"` // Número total de bloques
	SFreeBlocksCount int32 `binary:"little"` // Número de bloques libres
	SFreeInodesCount int32 `binary:"little"` // Número de inodos libres
	SMtime           int32 `binary:"little"` // Fecha y hora de la última vez que se montó el sistema de archivos
	SUmtime          int32 `binary:"little"` // Fecha y hora de la última vez que se desmontó el sistema de archivos
	SMntCount        int32 `binary:"little"` // Número de veces que se ha montado el sistema de archivos
	SMagic           int32 `binary:"little"` // Identificador del sistema de archivos (0xEF53)
	SInodeSize       int32 `binary:"little"` // Tamaño del inodo
	SBlockSize       int32 `binary:"little"` // Tamaño del bloque
	SFirstInode      int32 `binary:"little"` // Primer inodo libre (dirección del inodo)
	SFirstBlock      int32 `binary:"little"` // Primer bloque libre (dirección del bloque)
	SBmInodeStart    int32 `binary:"little"` // Inicio del bitmap de inodos
	SBmBlockStart    int32 `binary:"little"` // Inicio del bitmap de bloques
	SInodeStart      int32 `binary:"little"` // Inicio de la tabla de inodos
	SBlockStart      int32 `binary:"little"` // Inicio de la tabla de bloques
}
