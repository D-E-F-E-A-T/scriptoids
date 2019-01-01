#!/usr/bin/env python3

import os
import re
import stat
import sys

import click
import toml
from colorama import Fore, Style
from semver import VersionInfo

SCRIPTOIDS_VERSION = VersionInfo(0, 1, 0)
SCRIPT_CONFIG_FILENAME = 'script_info.toml'
SCRIPT_CONFIG_REQUIRED_FIELDS = ['name', 'entry']
FILENAME_VALIDATION_REGEX = re.compile(r'^[\w\-. ]+$')




try:
    SCRIPTOID_HOME = os.environ['SCRIPTOID_HOME']
except KeyError:
    die('No value is set for environment variable SCRIPTOID_HOME. A suggested path is $HOME/scriptoids.')


def directory_contains_scriptoid(directory):
    """
    Determines whether or not the given directory contains a scriptoid.
    :param directory: Directory to validate.
    :return: Whether or not the given directory contains a scriptoid.
    """
    return os.path.exists(os.path.join(directory, SCRIPT_CONFIG_FILENAME))


def script_definition_is_valid(script_definition):
    """
    Determines whether or not a given script definition contains the necessary fields.
    :param script_definition: Script definition, as a dictionary, to validate.
    :return: Whether or not the given script definition is valid as of the current version.
    """
    return 'script' in script_definition and \
           all(field in script_definition['script'] for field in SCRIPT_CONFIG_REQUIRED_FIELDS) and \
           FILENAME_VALIDATION_REGEX.match(script_definition['script']['name'])


def get_scriptoid_directory(scriptoid_name):
    """
    Gets the expected directory housing a given scriptoid.
    :param scriptoid_name: Name of the scriptoid to get the directory for.
    :return: Expected directory for the given scriptoid.
    """
    return os.path.join(SCRIPTOID_HOME, scriptoid_name)


def get_bin_directory():
    """
    Gets the expected `bin` directory within the scriptoid home path.
    :return: Path to the scriptoid `bin` directory.
    """
    return os.path.join(SCRIPTOID_HOME, 'bin')


@click.group()
def cli():
    pass


@click.command()
@click.argument('name')
def new(name):
    """
    Scaffolds a new scriptoid.
    """
    info(f'Creating new scriptoid {name}...')
    scriptoid_path = get_scriptoid_directory(name)

    try:
        os.makedirs(scriptoid_path, exist_ok=True)
    except FileExistsError:
        die(f'A scriptoid named {name} already exists.')

    with open(os.path.join(scriptoid_path, f'{name}.sh'), 'w') as script:
        script.write('#!/bin/bash')

    with open(os.path.join(scriptoid_path, 'script_info.toml'), 'w') as metadata:
        toml.dump({
            'script': {
                'name': name,
                'entry': f'{name}.sh',
                'description': '',
                'script_dependencies': [],
                'path_dependencies': []
            },
        }, metadata)

    success(f'Created new scriptoid {name} in {scriptoid_path}.')


@click.command()
@click.argument('name')
def link(name):
    """
    Adds a link from a given scriptoid to the `bin` directory.
    """
    info(f'Linking scriptoid {name}...')
    bin_directory = get_bin_directory()
    scriptoid_directory = get_scriptoid_directory(name)

    if not os.path.isdir(bin_directory):
        info(f'Directory {bin_directory} does not exist, creating...')
        os.makedirs(bin_directory, exist_ok=True)

    if not directory_contains_scriptoid(scriptoid_directory):
        die(f'No scriptoid `{name}` found in directory {scriptoid_directory}.')

    full_script_definition = toml.load(open(os.path.join(scriptoid_directory, 'script_info.toml'), 'r'))

    if not script_definition_is_valid(full_script_definition):
        die(f"{name}'s {SCRIPT_CONFIG_FILENAME} file is not valid.")

    script_info = full_script_definition['script']
    entry = script_info['entry']
    dest_name = script_info['name']
    dest_path = os.path.abspath(os.path.join(bin_directory, dest_name))

    info('Creating symlink...')
    os.symlink(os.path.abspath(os.path.join(scriptoid_directory, entry)), dest_path)

    info('Setting executable...')
    link_st = os.stat(dest_path)
    os.chmod(dest_path, link_st.st_mode | stat.S_IEXEC)

    success(f'Linked scriptoid {name}.')


@click.command()
@click.argument('name')
def unlink(name):
    """
    Removes a link for a given scriptoid from the `bin` directory.
    """

    info(f'Unlinking scriptoid {name}...')

    bin_directory = get_bin_directory()

    if not os.path.isdir(bin_directory):
        die('No scripts are currently linked.')

    link_path = os.path.join(bin_directory, name)

    if not os.path.islink(link_path):
        die(f'Script `{name}` does not appear to be linked.')

    info('Removing link...')
    os.remove(link_path)

    success(f'Unlinked scriptoid {name}.')


for command in [
    new,
    link,
    unlink
]:
    cli.add_command(command)

if __name__ == '__main__':
    cli()
