#!/bin/bash

# Base directory containing your Terraform modules
base_dir="/Users/martin.drotar/Student/open_banking/EKS-BankingKube/modules"

# Array of module names
modules=("databases" "networking" "security" "eks" "node_groups" "storage")

# Loop over each module
for module_name in "${modules[@]}"; do
  # Full path to the module
  module="$base_dir/$module_name"

  if [ -d "$module" ]; then
    echo "Formatting module: $module"
    
    # Change to the module directory
    cd "$module" || exit

    # Format the Terraform code and echo the result
    echo "$(terraform fmt)"

    # Change back to the original directory
    cd - || exit
  fi
done