import semver


def load_scriptoid_from_toml(toml_data):
    if 'scriptoid' not in toml_data:
        raise KeyError("A 'scriptoid' section was not found in the given scriptoid definition.")

    scriptoid_info = toml_data['scriptoid']

    if 'name' not in scriptoid_info:
        raise KeyError("A 'name' field was not found in the given scriptoid definition.")

    if 'entry_file' not in scriptoid_info:
        raise KeyError("An 'entry_file' field was not found in the given scriptoid definition.")

    if 'version' not in scriptoid_info:
        raise KeyError("A 'version' field was not found in the given scriptoid definition.")

    try:
        version = semver.parse_version_info(scriptoid_info['version'])
    except ValueError:
        raise KeyError("The 'version' field in the given scriptoid definition is not a valid SemVer version.")

    return Scriptoid(
        scriptoid_info['name'],
        version,
        scriptoid_info['entry_file'],
        scriptoid_info.get('description', ''),
        scriptoid_info.get('script_dependencies', []),
        scriptoid_info.get('path_dependencies', [])
    )


class Scriptoid:
    def __init__(self, name, version, entry_file, description, script_dependencies, path_dependencies):
        self.name = name
        self.version = version
        self.entry_file = entry_file
        self.description = description
        self.script_dependencies = script_dependencies
        self.path_dependencies = path_dependencies

    def fulfills_requirement(self, name, version):
        return self.name == name and self.version.major == version.major
