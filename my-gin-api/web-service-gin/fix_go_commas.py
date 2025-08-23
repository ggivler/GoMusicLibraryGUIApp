#!/usr/bin/env python3
"""
Script to fix Go struct literal syntax errors by adding missing commas.
"""
import re

def fix_go_file(file_path):
    """Fix Go struct literal syntax by adding missing commas."""
    
    # Read the file
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Pattern to match struct field assignments that should have commas
    # This matches lines like: \t\tFieldName: "value" (without comma at end)
    # But not the closing brace line
    pattern = r'(\t+)(AlphabetizingLetter|FullPathToFolder|OriginalFilename|SongTitle|Voicing|ComposerOrArranger|FileType|FileCreateDate|LibraryType):\s*([^,\n]+)(?!\s*,)(\s*\n)'
    
    # Add commas to these lines, except when they're followed by a closing brace
    def add_comma(match):
        indent, field, value, newline = match.groups()
        return f"{indent}{field}: {value},{newline}"
    
    # Apply the fix
    fixed_content = re.sub(pattern, add_comma, content)
    
    # Special case: remove comma from the last field in each struct (before closing brace)
    # Pattern to match: field: "value", followed by whitespace and }
    pattern_last = r'(\t+)(AlphabetizingLetter|FullPathToFolder|OriginalFilename|SongTitle|Voicing|ComposerOrArranger|FileType|FileCreateDate|LibraryType):\s*([^,\n]+),(\s*\n\s*\})'
    
    def remove_last_comma(match):
        indent, field, value, closing = match.groups()
        return f"{indent}{field}: {value}{closing}"
    
    fixed_content = re.sub(pattern_last, remove_last_comma, fixed_content)
    
    # Write the fixed content back
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(fixed_content)
    
    print(f"Fixed {file_path}")

if __name__ == "__main__":
    fix_go_file("main.go")
