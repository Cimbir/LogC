@echo off

REM Navigate to the project directory
cd %~dp0

REM Build the application
go build -o logc.exe ./cmd/app

REM Check if the build was successful
if %errorlevel% equ 0 (
    echo Build successful. Running the application...
    REM Run the application
    logc.exe
) else (
    echo Build failed.
)