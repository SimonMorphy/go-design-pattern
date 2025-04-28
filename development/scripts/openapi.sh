#!/bin/bash

# Set base paths
API_DIR="../api/openapi"
OUTPUT_DIR="../internal/ports"

# Ensure output directory exists
mkdir -p $OUTPUT_DIR

# Check if any YAML files exist
if [ ! "$(ls -A $API_DIR/*.yml 2>/dev/null)" ]; then
    echo "No YAML files found in $API_DIR"
    echo "Processing $API_DIR/user.yml directly..."
    
    # Process the known file directly
    filename="user"
    output_file="$OUTPUT_DIR/${filename}s.gen.go"
    
    # Execute code generation
    oapi-codegen -o "$output_file" -generate "server,types" -package "${filename}s" "$API_DIR/${filename}.yml"
    
    if [ $? -eq 0 ]; then
        echo "Successfully generated code: $output_file"
    else
        echo "Failed to generate code: $API_DIR/${filename}.yml"
        exit 1
    fi
else
    # Find all YAML files and generate code
    for yaml_file in $API_DIR/*.yml; do
        # Get filename without extension
        filename=$(basename -- "$yaml_file")
        filename_no_ext="${filename%.*}"
        
        # Set output file path
        output_file="$OUTPUT_DIR/${filename_no_ext}s.gen.go"
        
        echo "Processing $yaml_file..."
        
        # Execute code generation
        oapi-codegen -o "$output_file" -generate "server,types" -package "${filename_no_ext}s" "$yaml_file"
        
        if [ $? -eq 0 ]; then
            echo "Successfully generated code: $output_file"
        else
            echo "Failed to generate code: $yaml_file"
            exit 1
        fi
    done
fi

echo "All OpenAPI code generation completed!"