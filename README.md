# scriptoids

> A system for managing small scripts.

## Usage

### Installation

Python 3.6 or later is required.

```sh
$ git clone https://github.com/dhsavell/scriptoids.git ~/scriptoids
$ cd ~/scriptoids && pip3 install -r requirements.txt
```

Finally, add `~/scriptoids/bin` to your PATH.

### Creating a scriptoid
```sh
$ ./scriptoid.py new your_script_name
```

### Linking and using a scriptoid
```sh
$ ./scriptoid.py link your_script_name
$ your_script_name
```

### Unlinking a scriptoid
```sh
$ ./scriptoid.py unlink your_script_name
```

## Overview

"Scriptoids" are small scripts contained in folders with `script_info.toml`
files. New scriptoids can be created with the following command:

```sh
$ ./scriptoid.py new your_script_name
```

This will generate the following structure:
```
...
 |_ your_script_name
 |   |_ script_info.toml
 |   |_ your_script_name.sh
 |_ ...
```

The default generated `script_info.toml` file will point to
`your_script_name.sh`, meaning that the scriptoid is immediately usable.

### script_info.toml

A `script_info.toml` file describes a scriptoid. The following fields are
available:

- **[script]**
  - **name** (required): Name of the script.
  - **entry** (required): "Entry point" of the script. This is the file that
    is executed when calling the scriptoid by name.
  - **description**: Description displayed while linking the scriptoid.
  - **script_dependencies**: An array of strings listing other scriptoids that
    need to be available in order for this scriptoid to run. If the necessary
    scriptoids can be found, they will automatically be linked.
  - **path_dependencies**: Programs needed in the PATH for this scriptoid to
    run. The scriptoid will fail to link if these are not found.