#!/bin/bash

# Input file
input_file="script_output.txt"

# Output file
output_file="script_output.md"

# Start the output file with a header
echo "# Terraform Validation Output" > "$output_file"

# Read the input file line by line
while IFS= read -r line
do
  # Add the line to the output file, formatted as a code block
  echo "    $line" >> "$output_file"
done < "$input_file"