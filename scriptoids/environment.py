import os

import shutil
import toml
from semver import VersionInfo

from scriptoids.cliutil import die
from scriptoids.scriptoid import load_scriptoid_from_toml


class ScriptoidEnvironment:
    """
    Represents the working Scriptoid environment. This is where scriptoid sources are stored, as well as the "bin"
    directory where symlinks are created to.
    """

    def __init__(self,
                 home_directory,
                 version=VersionInfo(0, 0, 0, build='dev'),
                 scriptoid_definition_filename='script_info.toml'):
        self.home_directory = home_directory
        self.bin_directory = os.path.join(home_directory, 'bin')
        self.version = version
        self.scriptoid_definition_filename = scriptoid_definition_filename

        self._init_required_directories()

    def _init_required_directories(self):
        """
        Initializes the required Scriptoid directories within this ScriptoidEnvironment.
        """
        try:
            os.makedirs(self.home_directory, exist_ok=True)
            os.makedirs(self.bin_directory, exist_ok=True)
        except OSError:
            die(f'Failed to create the following required directories: {self.home_directory}, {self.bin_directory}')

    def scriptoid_exists(self, scriptoid_name):
        """
        Determines whether or not a Scriptoid is present within this ScriptoidEnvironment.
        :param scriptoid_name: Name of the Scriptoid to check for.
        :return: Whether or not the given Scriptoid exists.
        """
        return scriptoid_name != 'bin' and os.path.isdir(
            os.path.join(self.home_directory, scriptoid_name)) and os.path.exists(
            os.path.join(self.home_directory, scriptoid_name, self.scriptoid_definition_filename))

    def get_scriptoid(self, scriptoid_name):
        """
        Gets a Scriptoid instance by name within this ScriptoidEnvironment.
        :param scriptoid_name: Name of the Scriptoid to get.
        :return: A Scriptoid instance describing the requested scriptoid.
        """
        if not self.scriptoid_exists(scriptoid_name):
            raise NameError(f'No scriptoid named {scriptoid_name} was found.')

        scriptoid_definition_path = os.path.join(self.home_directory, scriptoid_name,
                                                 self.scriptoid_definition_filename)

        with open(scriptoid_definition_path, 'r') as scriptoid_definition:
            return load_scriptoid_from_toml(toml.load(scriptoid_definition))

    def scriptoid_is_linked(self, scriptoid_name):
        """
        Determines whether or not a given scriptoid is linked within this ScriptoidEnvironment.
        :param scriptoid_name: Name of the Scriptoid to check link status of.
        :return: Whether or not the given Scriptoid is currently linked.
        """
        return self.scriptoid_exists(scriptoid_name) and os.path.islink(
            os.path.join(self.home_directory, 'bin', scriptoid_name))

    def link_scriptoid(self, scriptoid_name):
        """
        Creates a symbolic link to a given scriptoid from this ScriptoidEnvironment's bin directory. This will
        recursively link dependencies if they are present, but not linked.
        :param scriptoid_name: Name of the scriptoid to link.
        """
        scriptoid = self.get_scriptoid(scriptoid_name)

        for required_program in scriptoid.path_dependencies:
            if not shutil.which(required_program):
                raise EnvironmentError(
                    f'Scriptoid {scriptoid_name} requires {required_program} in your PATH in order to be linked, '
                    f'but it was not found.')
