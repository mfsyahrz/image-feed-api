#!/bin/sh

config_file="gomockhandler.json"
mocks_destination="mocks"
go_sdk="/usr/local/go/src"  # Ensure this is the correct GOROOT path for your system

echo "Mock Options"
echo "1. Mock Packages"
echo "2. Check Mocks"
echo "Please choose an option:"

read option_num

# Option 1: Mock internal packages based on file path
if [ "$option_num" = "1" ]; then
    echo "Enter file PATHNAME relative to content root (e.g., internal/service/comment/comment_service.go):"
    read pathname

    mock_path=$(echo "$pathname" | sed 's|internal/||')

    echo "Generating mocks for $pathname..."
    echo "Mocks will be placed in $mocks_destination/$mock_path"
    
    gomockhandler -config="$config_file" -source="$pathname" -destination="$mocks_destination/$mock_path"

# Option 2: Check if mocks exist
elif [ "$option_num" = "4" ]; then
    echo "Checking for existing mocks..."
    
    gomockhandler -config="$config_file" check

else
    echo "Invalid option. Please select a valid option (1, 2)."
fi
