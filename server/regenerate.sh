#!/bin/bash

# Copy the original file to a temporary copy
cp gqlgen.yml gqlgen_temp.yml

# Create a temporary file to store the paths of .graphqls files
TEMP_FILE=$(mktemp)

# Find all .graphqls files in the graphql directory but exclude .history and its subdirectories
find ../graphql -name "*.graphqls" ! -path "../graphql/.history*" > $TEMP_FILE

# Generate insertion text for the YAML file
INSERT_TEXT="schema:\n"
while read -r line; do
    INSERT_TEXT="${INSERT_TEXT}  - $line\n"
done < $TEMP_FILE

# Use sed to insert the text into gqlgen_temp.yml
sed -i "/# Schema routes go here:/a $INSERT_TEXT" gqlgen_temp.yml

# Clean up the temporary file with paths
rm $TEMP_FILE

# Ensure all dependencies are available
go get github.com/99designs/gqlgen@v0.17.39

# Execute gqlgen with the temporary configuration
go run github.com/99designs/gqlgen generate --config gqlgen_temp.yml

# Remove the temporary configuration copy
rm gqlgen_temp.yml

echo "server schema generation done"