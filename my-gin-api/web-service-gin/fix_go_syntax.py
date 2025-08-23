#!/usr/bin/env python3
"""
Script to fix Go struct literal syntax errors by replacing JSON tag names with Go field names.
"""

def fix_go_file(file_path):
    """Fix Go struct literal syntax by replacing JSON tags with Go field names."""
    
    # Mapping of JSON tag names to Go field names
    replacements = {
        '"alphabetizing letter":': 'AlphabetizingLetter:',
        '"full path to folder":': 'FullPathToFolder:',
        '"original filename":': 'OriginalFilename:',
        '"song title":': 'SongTitle:',
        '"voicing":': 'Voicing:',
        '"composer or arranger":': 'ComposerOrArranger:',
        '"file type":': 'FileType:',
        '"file create date":': 'FileCreateDate:',
        '"library type":': 'LibraryType:',
    }
    
    # Read the file
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Apply all replacements
    for old, new in replacements.items():
        content = content.replace(old, new)
    
    # Write the fixed content back
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"Fixed {file_path}")

if __name__ == "__main__":
    fix_go_file("main.go")
