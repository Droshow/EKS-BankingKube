#!/bin/bash

# Base directory containing your Terraform modules
base_dir="/Users/martin.drotar/Student/open_banking/EKS-BankingKube/modules"

# Array of module names
modules=("databases" "networking" "security" "eks" "node_groups" "storage")

# Full path to the output file
output_file="$base_dir/script_output.txt"

# Start output redirection
exec > "$output_file" 2>&1

# Loop over each module
for module_name in "${modules[@]}"; do
  # Full path to the module
  module="$base_dir/$module_name"

  if [ -d "$module" ]; then
    echo "Checking module: $module"
    
    # Change to the module directory
    cd "$module" || exit

    # Initialize Terraform
    terraform init -no-color -backend=false

    # Validate the module
    terraform validate -no-color

    # Change back to the original directory
    cd - || exit
  fi
done

# End output redirection
exec &> /dev/tty