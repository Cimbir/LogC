#!/bin/bash

# Navigate to the project directory
cd "$(dirname "$0")"

# Build the application
go build -o logc ./cmd/app

# Check if the build was successful
if [ $? -eq 0 ]; then
    echo "Build successful. Running the application..."
    # Run the application
    ./logc
else
    echo "Build failed."
fi