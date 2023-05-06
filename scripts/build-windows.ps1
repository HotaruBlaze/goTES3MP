# Commit hash
$commitHash = (git rev-parse HEAD)

# Path Resolving
$goFiles = Get-ChildItem -Path 'src\' -Name '*.go' -File
$goFiles = $goFiles | ForEach-Object {"src/$_"}

# Build variables for go build
Set-Variable GOOS=windows
Set-Variable GOARCH=amd64

# Build command
go build -ldflags="-X 'main.Build=v0.0.0-Dev' -X 'main.GitCommit=$commitHash'" -o build/goTES3MP-Windows.exe $goFiles