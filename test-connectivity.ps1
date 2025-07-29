# 🧪 Script para probar la conectividad backend-frontend en Windows
# Ejecutar en PowerShell: .\test-connectivity.ps1

Write-Host "🧪 Probando conectividad Backend ↔ Frontend" -ForegroundColor Cyan
Write-Host "==============================================" -ForegroundColor Cyan

# Función para hacer peticiones HTTP
function Test-Endpoint {
    param (
        [string]$Url,
        [string]$Method = "GET",
        [string]$Body = $null,
        [hashtable]$Headers = @{}
    )
    
    try {
        $params = @{
            Uri = $Url
            Method = $Method
            UseBasicParsing = $true
        }
        
        if ($Body) {
            $params.Body = $Body
            $params.ContentType = "application/json"
        }
        
        if ($Headers.Count -gt 0) {
            $params.Headers = $Headers
        }
        
        $response = Invoke-WebRequest @params
        return @{
            Success = $true
            StatusCode = $response.StatusCode
            Content = $response.Content
        }
    }
    catch {
        return @{
            Success = $false
            StatusCode = $_.Exception.Response.StatusCode.value__
            Error = $_.Exception.Message
        }
    }
}

# Verificar que el backend esté ejecutándose
Write-Host "🔍 Verificando backend en http://localhost:8080..." -ForegroundColor Yellow

# Probar endpoint de salud
Write-Host "📡 Probando endpoint /api/health..." -ForegroundColor White
$healthTest = Test-Endpoint -Url "http://localhost:8080/api/health"

if ($healthTest.Success) {
    Write-Host "✅ Backend funcionando correctamente" -ForegroundColor Green
    Write-Host "📄 Respuesta: $($healthTest.Content)" -ForegroundColor Gray
} else {
    Write-Host "❌ Backend no responde (HTTP $($healthTest.StatusCode))" -ForegroundColor Red
    Write-Host "🔧 Asegúrate de que el backend esté ejecutándose:" -ForegroundColor Yellow
    Write-Host "   cd backend; go run main.go" -ForegroundColor Gray
    exit 1
}

Write-Host ""

# Probar endpoint de sistemas de archivos
Write-Host "📡 Probando endpoint /api/filesystems..." -ForegroundColor White
$filesystemsTest = Test-Endpoint -Url "http://localhost:8080/api/filesystems"

if ($filesystemsTest.Success) {
    Write-Host "✅ Endpoint de sistemas de archivos funcionando" -ForegroundColor Green
    Write-Host "📄 Respuesta: $($filesystemsTest.Content)" -ForegroundColor Gray
} else {
    Write-Host "❌ Endpoint de sistemas de archivos no responde (HTTP $($filesystemsTest.StatusCode))" -ForegroundColor Red
}

Write-Host ""

# Probar endpoint de ejecución de comandos
Write-Host "📡 Probando endpoint /api/execute..." -ForegroundColor White
$commandTest = Test-Endpoint -Url "http://localhost:8080/api/execute" -Method "POST" -Body '{"command":"test command"}'

if ($commandTest.Success) {
    Write-Host "✅ Endpoint de ejecución de comandos funcionando" -ForegroundColor Green
    Write-Host "📄 Respuesta: $($commandTest.Content)" -ForegroundColor Gray
} else {
    Write-Host "❌ Endpoint de ejecución de comandos no responde (HTTP $($commandTest.StatusCode))" -ForegroundColor Red
}

Write-Host ""

# Verificar CORS
Write-Host "🔐 Verificando configuración CORS..." -ForegroundColor White
$corsHeaders = @{
    "Origin" = "http://localhost:3000"
    "Access-Control-Request-Method" = "GET"
}
$corsTest = Test-Endpoint -Url "http://localhost:8080/api/health" -Method "OPTIONS" -Headers $corsHeaders

if ($corsTest.Success) {
    Write-Host "✅ CORS configurado correctamente" -ForegroundColor Green
} else {
    Write-Host "⚠️  CORS podría no estar configurado correctamente" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "🌐 URLs para acceder:" -ForegroundColor Cyan
Write-Host "   📋 Backend API: http://localhost:8080/api/health" -ForegroundColor Gray
Write-Host "   🌍 Frontend: http://localhost:3000" -ForegroundColor Gray
Write-Host ""
Write-Host "🚀 Para iniciar los servicios:" -ForegroundColor Cyan
Write-Host "   Backend: cd backend; go run main.go" -ForegroundColor Gray
Write-Host "   Frontend: cd frontend; npm start" -ForegroundColor Gray
