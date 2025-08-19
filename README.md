# GoMusicLibraryGUIApp

This application provides a wizard-style GUI for processing music libraries using the Fyne framework.

## Running the Application

```bash
go run process_wizard_demo.go
```

## Known Issues

### Fyne "Getting favorite locations" Error on Windows

When you click the "Select the Music Library Folder" button, you may see this error message in the console:

```
Fyne error:  Getting favorite locations
  Cause: uri is not listable
  At: C:/Users/.../go/pkg/mod/fyne.io/fyne/v2@v2.6.2/dialog/file.go:367
```

**This is a harmless error** that occurs due to a known issue in the Fyne framework on Windows systems. The error happens because:

1. Fyne tries to populate the file dialog with Windows "favorite locations" (like Quick Access items)
2. Some of these system locations may not be accessible or listable
3. Fyne's internal logging system reports this as an error

**Important Notes:**
- ✅ **The file dialog still works correctly** - you can select folders without any issues
- ✅ **No functionality is lost** - this is purely a cosmetic logging issue
- ❌ **The error cannot be easily suppressed** because it comes from Fyne's internal C code that writes directly to stderr
- ❌ **This is not a bug in our application** - it's a known limitation of the Fyne framework on Windows

## Features
### File Menu
#### File -> Settings...
Selecting this menu will open a file dialog to open the settings file
In the yaml file there will be a Settings heading
and sub-headings for the config file name and the directory
for the directory location and a 
filename for the csv file and the location of the
csv file which is the output of walking the directory
where the pdf files are

The default file name is currently config.yml

### CSV File
#### CSV File Fields
The following array contains the headers for the csv file

The default file name currently is csv_output_full.csv

keywords = ["alphabetizing letter",
"full path to folder",
"original filename",
"song title",
"voicing",
"composer or arranger",
"file type",
"file create date",
"library type"
]
