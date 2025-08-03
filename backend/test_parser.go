package main

import (
	"fmt"
	"strings"
)

func parseCommandLine(commandLine string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(commandLine); i++ {
		char := commandLine[i]

		switch char {
		case '"', '\'':
			if !inQuotes {
				// Comenzar comillas
				inQuotes = true
				quoteChar = char
				// No agregar la comilla al resultado
			} else if char == quoteChar {
				// Terminar comillas
				inQuotes = false
				quoteChar = 0
				// No agregar la comilla al resultado
			} else {
				// Comilla diferente dentro de comillas
				current.WriteByte(char)
			}
		case ' ', '\t':
			if inQuotes {
				// Espacio dentro de comillas
				current.WriteByte(char)
			} else {
				// Espacio fuera de comillas - separador
				if current.Len() > 0 {
					parts = append(parts, current.String())
					current.Reset()
				}
			}
		default:
			current.WriteByte(char)
		}
	}

	// Agregar la última parte
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

func main() {
	testCommand := `mkdisk -size=5 -unit=M -path="./Discos/mis discos/disco.mia"`
	fmt.Printf("Comando original: %s\n", testCommand)

	parts := parseCommandLine(testCommand)
	fmt.Printf("Número de partes: %d\n", len(parts))

	for i, part := range parts {
		fmt.Printf("Parte %d: '%s'\n", i, part)
	}
}
