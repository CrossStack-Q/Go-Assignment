#!/bin/bash

# Define source files
source_files=("o1.png" "o2.png" "o3.png")

# Define target file sets
target_sets=(
    "t1.png t2.png t3.png"
    "th1.png th2.png th3.png"
    "f1.png f2.png f3.png"
    "fi1.png fi2.png fi3.png"
)

# Loop through each target set and copy the source files
for targets in "${target_sets[@]}"; do
    read -r t1 t2 t3 <<< "$targets"
    cp "${source_files[0]}" "$t1"
    cp "${source_files[1]}" "$t2"
    cp "${source_files[2]}" "$t3"
    echo "Copied ${source_files[0]} -> $t1"
    echo "Copied ${source_files[1]} -> $t2"
    echo "Copied ${source_files[2]} -> $t3"
done

echo "All files copied successfully!"

