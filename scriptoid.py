#!/usr/bin/env python3

import os
import re
import stat
import sys

import click
import toml
from colorama import Fore, Style
from semver import VersionInfo

SCRIPTOID_HOME = os.environ['SCRIPTOID_HOME']
SCRIPTOIDS_VERSION = VersionInfo(0, 1, 0)
SCRIPT_CONFIG_FILENAME = 'script_info.toml'
SCRIPT_CONFIG_REQUIRED_FIELDS = ['name', 'entry']
FILENAME_VALIDATION_REGEX = re.compile(r'^[\w\-. ]+$')


def info(msg):
    """
    Prints an informational message.
    :param msg: Message to display.
    """
    print(f'{Style.BRIGHT}{Fore.BLACK}  {msg}{Style.RESET_ALL}')


def success(msg):
    """
    Prints a message indicating success.
    :param msg: Message to display.
    """
    print(f'{Style.BRIGHT}{Fore.GREEN}✔ {Fore.BLACK}{msg}{Style.RESET_ALL}')


def die(msg, exit_code=1):
    """
    Prints a message denoting a failure and exits with a given exit code.
    :param msg: Message to display.
    :param exit_code: Code to return when exiting.
    """
    sys.stderr.write(f'{Style.BRIGHT}{Fore.RED}✘ An error occurred: {Fore.BLACK}{msg}{Style.RESET_ALL}\n')
    exit(exit_code)


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

    try:
        os.mkdir(name)
    except FileExistsError:
        die(f'A directory named {name} already exists.')

    with open(os.path.join(name, f'{name}.sh'), 'w') as script:
        script.write('#!/bin/bash')

    with open(os.path.join(name, 'script_info.toml'), 'w') as metadata:
        toml.dump({
            'script': {
                'name': name,
                'entry': f'{name}.sh',
                'description': '',
                'script_dependencies': [],
                'path_dependencies': []
            },
        }, metadata)

    success(f'Created new scriptoid in {name}/.')


@click.command()
@click.argument('name')
def link(name):
    """
    Adds a link from a given scriptoid to the `bin` directory.
    """
    info(f'Linking script {name} to bin/...')

    if not os.path.isdir('bin'):
        info('Directory bin/ does not exist, creating...')
        os.mkdir('bin')

    if not directory_contains_scriptoid(name):
        die(f'{name} does not appear to be a valid script directory.')

    full_script_definition = toml.load(open(os.path.join(name, 'script_info.toml'), 'r'))

    if not script_definition_is_valid(full_script_definition):
        die(f"{name}'s {SCRIPT_CONFIG_FILENAME} file is not valid.")

    script_info = full_script_definition['script']
    entry = script_info['entry']
    dest_name = script_info['name']

    dest_path = os.path.abspath(os.path.join('bin', dest_name))

    info('Creating symlink...')
    os.symlink(os.path.abspath(os.path.join(name, entry)), dest_path)

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

    info(f'Unlinking {name} from bin/...')

    if not os.path.isdir('bin'):
        die('No scripts are currently linked.')

    if not os.path.islink(os.path.join('bin', name)):
        die(f'Script {name} does not appear to be linked.')

    info('Removing link...')
    os.remove(os.path.join('bin', name))

    success(f'Unlinked scriptoid {name}.')


for command in [
    new,
    link,
    unlink
]:
    cli.add_command(command)

if __name__ == '__main__':
    if not SCRIPTOID_HOME:
        die('No value is set for environment variable SCRIPTOID_HOME. A suggested path is $HOME/scriptoids.')

    cli()
