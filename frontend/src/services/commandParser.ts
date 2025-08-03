// services/commandParser.ts

export interface ParsedCommand {
  command: string;
  parameters: Record<string, string>;
  errors: string[];
}

export interface CommandDefinition {
  name: string;
  parameters: {
    [key: string]: {
      required: boolean;
      type: 'string' | 'number' | 'enum';
      values?: string[]; // Para enums
      defaultValue?: string;
      description: string;
    };
  };
}

// Definiciones de comandos según el proyecto
export const COMMAND_DEFINITIONS: Record<string, CommandDefinition> = {
  mkdisk: {
    name: 'mkdisk',
    parameters: {
      size: {
        required: true,
        type: 'number',
        description: 'Tamaño del disco a crear (debe ser positivo y mayor que cero)'
      },
      fit: {
        required: false,
        type: 'enum',
        values: ['BF', 'FF', 'WF'],
        defaultValue: 'FF',
        description: 'Tipo de ajuste: BF (Best Fit), FF (First Fit), WF (Worst Fit)'
      },
      unit: {
        required: false,
        type: 'enum',
        values: ['K', 'M'],
        defaultValue: 'M',
        description: 'Unidades: K (Kilobytes), M (Megabytes)'
      },
      path: {
        required: true,
        type: 'string',
        description: 'Ruta donde se creará el archivo del disco'
      }
    }
  },
  rmdisk: {
    name: 'rmdisk',
    parameters: {
      path: {
        required: true,
        type: 'string',
        description: 'Ruta del archivo de disco a eliminar'
      }
    }
  },
  fdisk: {
    name: 'fdisk',
    parameters: {
      size: {
        required: true,
        type: 'number',
        description: 'Tamaño de la partición a crear'
      },
      unit: {
        required: false,
        type: 'enum',
        values: ['B', 'K', 'M'],
        defaultValue: 'K',
        description: 'Unidades: B (Bytes), K (Kilobytes), M (Megabytes)'
      },
      path: {
        required: true,
        type: 'string',
        description: 'Ruta del disco donde se creará la partición'
      },
      type: {
        required: false,
        type: 'enum',
        values: ['P', 'E', 'L'],
        defaultValue: 'P',
        description: 'Tipo de partición: P (Primaria), E (Extendida), L (Lógica)'
      },
      fit: {
        required: false,
        type: 'enum',
        values: ['BF', 'FF', 'WF'],
        defaultValue: 'WF',
        description: 'Tipo de ajuste: BF (Best Fit), FF (First Fit), WF (Worst Fit)'
      },
      name: {
        required: true,
        type: 'string',
        description: 'Nombre de la partición'
      }
    }
  },
  mount: {
    name: 'mount',
    parameters: {
      path: {
        required: true,
        type: 'string',
        description: 'Ruta del disco que se montará'
      },
      name: {
        required: true,
        type: 'string',
        description: 'Nombre de la partición a montar'
      }
    }
  },
  mkfs: {
    name: 'mkfs',
    parameters: {
      id: {
        required: true,
        type: 'string',
        description: 'ID de la partición a formatear'
      },
      type: {
        required: false,
        type: 'enum',
        values: ['full'],
        defaultValue: 'full',
        description: 'Tipo de formateo: full (completo)'
      }
    }
  },
  login: {
    name: 'login',
    parameters: {
      user: {
        required: true,
        type: 'string',
        description: 'Nombre del usuario'
      },
      pass: {
        required: true,
        type: 'string',
        description: 'Contraseña del usuario'
      },
      id: {
        required: true,
        type: 'string',
        description: 'ID de la partición'
      }
    }
  },
  logout: {
    name: 'logout',
    parameters: {}
  },
  mkgrp: {
    name: 'mkgrp',
    parameters: {
      name: {
        required: true,
        type: 'string',
        description: 'Nombre del grupo'
      }
    }
  },
  rmgrp: {
    name: 'rmgrp',
    parameters: {
      name: {
        required: true,
        type: 'string',
        description: 'Nombre del grupo a eliminar'
      }
    }
  },
  mkusr: {
    name: 'mkusr',
    parameters: {
      user: {
        required: true,
        type: 'string',
        description: 'Nombre del usuario'
      },
      pass: {
        required: true,
        type: 'string',
        description: 'Contraseña del usuario'
      },
      grp: {
        required: true,
        type: 'string',
        description: 'Grupo al que pertenece el usuario'
      }
    }
  },
  rmusr: {
    name: 'rmusr',
    parameters: {
      user: {
        required: true,
        type: 'string',
        description: 'Nombre del usuario a eliminar'
      }
    }
  },
  chgrp: {
    name: 'chgrp',
    parameters: {
      user: {
        required: true,
        type: 'string',
        description: 'Nombre del usuario'
      },
      grp: {
        required: true,
        type: 'string',
        description: 'Nuevo grupo del usuario'
      }
    }
  },
  mkfile: {
    name: 'mkfile',
    parameters: {
      path: {
        required: true,
        type: 'string',
        description: 'Ruta del archivo a crear'
      },
      r: {
        required: false,
        type: 'string',
        description: 'Crear carpetas padres si no existen'
      },
      size: {
        required: false,
        type: 'number',
        description: 'Tamaño del archivo en bytes'
      },
      cont: {
        required: false,
        type: 'string',
        description: 'Ruta del archivo fuente para el contenido'
      }
    }
  },
  mkdir: {
    name: 'mkdir',
    parameters: {
      path: {
        required: true,
        type: 'string',
        description: 'Ruta de la carpeta a crear'
      },
      p: {
        required: false,
        type: 'string',
        description: 'Crear carpetas padres si no existen'
      }
    }
  },
  cat: {
    name: 'cat',
    parameters: {
      filen: {
        required: true,
        type: 'string',
        description: 'Lista de archivos a mostrar'
      }
    }
  },
  rep: {
    name: 'rep',
    parameters: {
      name: {
        required: true,
        type: 'enum',
        values: ['mbr', 'disk', 'inode', 'block', 'bm_inode', 'bm_block', 'tree', 'sb', 'file', 'ls'],
        description: 'Tipo de reporte a generar'
      },
      path: {
        required: true,
        type: 'string',
        description: 'Ruta donde se guardará el reporte'
      },
      id: {
        required: true,
        type: 'string',
        description: 'ID de la partición'
      },
      path_file_ls: {
        required: false,
        type: 'string',
        description: 'Ruta del archivo o carpeta para reportes file y ls'
      }
    }
  },
  mounted: {
    name: 'mounted',
    parameters: {}
  }
};

export class CommandParser {
  /**
   * Parsea un comando completo en sus componentes
   */
  parseCommand(commandText: string): ParsedCommand {
    const trimmedCommand = commandText.trim();
    
    if (!trimmedCommand) {
      return {
        command: '',
        parameters: {},
        errors: ['Comando vacío']
      };
    }

    // Extraer el comando principal
    const parts = this.tokenizeCommand(trimmedCommand);
    if (parts.length === 0) {
      return {
        command: '',
        parameters: {},
        errors: ['Comando inválido']
      };
    }

    const commandName = parts[0].toLowerCase();
    const parameterTokens = parts.slice(1);

    // Verificar si el comando existe
    const commandDef = COMMAND_DEFINITIONS[commandName];
    if (!commandDef) {
      return {
        command: commandName,
        parameters: {},
        errors: [`Comando desconocido: ${commandName}`]
      };
    }

    // Parsear parámetros
    const { parameters, errors } = this.parseParameters(parameterTokens);
    
    // Validar parámetros según la definición
    const validationErrors = this.validateParameters(parameters, commandDef);
    
    // Aplicar valores por defecto
    const finalParameters = this.applyDefaults(parameters, commandDef);

    return {
      command: commandName,
      parameters: finalParameters,
      errors: [...errors, ...validationErrors]
    };
  }

  /**
   * Tokeniza el comando respetando comillas
   */
  private tokenizeCommand(command: string): string[] {
    const tokens: string[] = [];
    let current = '';
    let inQuotes = false;
    let quoteChar = '';

    for (let i = 0; i < command.length; i++) {
      const char = command[i];
      const nextChar = command[i + 1];

      if ((char === '"' || char === "'") && !inQuotes) {
        // Inicio de comillas
        inQuotes = true;
        quoteChar = char;
      } else if (char === quoteChar && inQuotes) {
        // Fin de comillas
        inQuotes = false;
        quoteChar = '';
      } else if (char === ' ' && !inQuotes) {
        // Espacio fuera de comillas - separador de tokens
        if (current.trim()) {
          tokens.push(current.trim());
          current = '';
        }
      } else {
        // Carácter normal
        current += char;
      }
    }

    // Agregar el último token si existe
    if (current.trim()) {
      tokens.push(current.trim());
    }

    return tokens;
  }

  /**
   * Parsea los parámetros del comando
   */
  private parseParameters(tokens: string[]): { parameters: Record<string, string>; errors: string[] } {
    const parameters: Record<string, string> = {};
    const errors: string[] = [];

    for (const token of tokens) {
      if (token.startsWith('-')) {
        // Es un parámetro
        const equalIndex = token.indexOf('=');
        
        if (equalIndex === -1) {
          // Parámetro sin valor (flag)
          const paramName = token.slice(1).toLowerCase();
          parameters[paramName] = 'true';
        } else {
          // Parámetro con valor
          const paramName = token.slice(1, equalIndex).toLowerCase();
          const paramValue = token.slice(equalIndex + 1);
          
          if (!paramName) {
            errors.push(`Nombre de parámetro vacío en: ${token}`);
            continue;
          }

          parameters[paramName] = paramValue;
        }
      } else {
        errors.push(`Token inválido (debe empezar con -): ${token}`);
      }
    }

    return { parameters, errors };
  }

  /**
   * Valida los parámetros según la definición del comando
   */
  private validateParameters(parameters: Record<string, string>, commandDef: CommandDefinition): string[] {
    const errors: string[] = [];

    // Verificar parámetros requeridos
    for (const [paramName, paramDef] of Object.entries(commandDef.parameters)) {
      if (paramDef.required && !(paramName in parameters)) {
        errors.push(`Parámetro requerido faltante: -${paramName}`);
      }
    }

    // Verificar parámetros válidos y tipos
    for (const [paramName, paramValue] of Object.entries(parameters)) {
      const paramDef = commandDef.parameters[paramName];
      
      if (!paramDef) {
        errors.push(`Parámetro desconocido: -${paramName}`);
        continue;
      }

      // Validar tipo
      if (paramDef.type === 'number') {
        const numValue = Number(paramValue);
        if (isNaN(numValue) || numValue <= 0) {
          errors.push(`El parámetro -${paramName} debe ser un número positivo mayor que cero`);
        }
      } else if (paramDef.type === 'enum' && paramDef.values) {
        const upperValue = paramValue.toUpperCase();
        const upperValues = paramDef.values.map(v => v.toUpperCase());
        if (!upperValues.includes(upperValue)) {
          errors.push(`Valor inválido para -${paramName}. Valores permitidos: ${paramDef.values.join(', ')}`);
        }
      }
    }

    return errors;
  }

  /**
   * Aplica valores por defecto a los parámetros
   */
  private applyDefaults(parameters: Record<string, string>, commandDef: CommandDefinition): Record<string, string> {
    const result = { ...parameters };

    for (const [paramName, paramDef] of Object.entries(commandDef.parameters)) {
      if (!(paramName in result) && paramDef.defaultValue !== undefined) {
        result[paramName] = paramDef.defaultValue;
      }
    }

    return result;
  }

  /**
   * Obtiene ayuda para un comando específico
   */
  getCommandHelp(commandName: string): string {
    const commandDef = COMMAND_DEFINITIONS[commandName.toLowerCase()];
    if (!commandDef) {
      return `Comando desconocido: ${commandName}`;
    }

    let help = `Comando: ${commandDef.name}\n\nParámetros:\n`;
    
    for (const [paramName, paramDef] of Object.entries(commandDef.parameters)) {
      const required = paramDef.required ? '(Obligatorio)' : '(Opcional)';
      const defaultVal = paramDef.defaultValue ? ` [Default: ${paramDef.defaultValue}]` : '';
      const values = paramDef.values ? ` [Valores: ${paramDef.values.join(', ')}]` : '';
      
      help += `  -${paramName} ${required}${defaultVal}${values}\n`;
      help += `    ${paramDef.description}\n\n`;
    }

    return help;
  }

  /**
   * Lista todos los comandos disponibles
   */
  listCommands(): string[] {
    return Object.keys(COMMAND_DEFINITIONS);
  }

  /**
   * Convierte un comando parseado de vuelta a string
   */
  commandToString(parsedCommand: ParsedCommand): string {
    let result = parsedCommand.command;
    
    for (const [key, value] of Object.entries(parsedCommand.parameters)) {
      if (value === 'true') {
        result += ` -${key}`;
      } else {
        // Agregar comillas si el valor contiene espacios
        const needsQuotes = value.includes(' ');
        const quotedValue = needsQuotes ? `"${value}"` : value;
        result += ` -${key}=${quotedValue}`;
      }
    }
    
    return result;
  }
}

// Instancia singleton del parser
export const commandParser = new CommandParser();
export default commandParser;

// Función de utilidad para validar comando antes de enviar
export const validateAndFormatCommand = (commandText: string): { 
  isValid: boolean; 
  formattedCommand: string; 
  errors: string[]; 
  originalCommand: string;
} => {
  const parsed = commandParser.parseCommand(commandText);
  
  return {
    isValid: parsed.errors.length === 0,
    formattedCommand: commandParser.commandToString(parsed),
    errors: parsed.errors,
    originalCommand: commandText
  };
};