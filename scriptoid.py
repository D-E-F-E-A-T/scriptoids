#!/usr/bin/env python3

import os
import re
import stat
import sys

import click
import toml
from colorama import Fore, Style

SCRIPT_CONFIG_FILENAME = 'script_info.toml'
SCRIPT_CONFIG_REQUIRED_FIELDS = ['name', 'entry']
FILENAME_VALIDATION_REGEX = re.compile(r'^[\w\-. ]+$')


def info(msg):
    print(f'{Style.BRIGHT}{Fore.BLACK}  {msg}{Style.RESET_ALL}')


def success(msg):
    print(f'{Style.BRIGHT}{Fore.GREEN}✔ {Fore.BLACK}{msg}{Style.RESET_ALL}')


def die(msg):
    sys.stderr.write(f'{Style.BRIGHT}{Fore.RED}✘ An error occurred: {Fore.BLACK}{msg}{Style.RESET_ALL}\n')
    exit(1)


def is_script_directory(directory):
    return os.path.exists(os.path.join(directory, SCRIPT_CONFIG_FILENAME))


def is_script_info_valid(script_info):
    return 'script' in script_info and \
           all(field in script_info['script'] for field in SCRIPT_CONFIG_REQUIRED_FIELDS) and \
           FILENAME_VALIDATION_REGEX.match(script_info['script']['name'])


@click.group()
def cli():
    pass


@click.command()
@click.argument('name')
def new(name):
    """
    Scaffolds a new scriptoid
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
                'description': '',          # TODO: Use when displaying info
                'script_dependencies': [],  # TODO: Validate dependencies when linking
                'path_dependencies': []     # TODO: Enforce programs in PATH when linking
            },
        }, metadata)

    success(f'Created new scriptoid in {name}/.')


@click.command()
@click.argument('name')
def link(name):
    """
    Creates a symbolic link from a scriptoid to ./bin/
    """
    info(f'Linking script {name} to bin/...')

    if not os.path.isdir('bin'):
        info('Directory bin/ does not exist, creating...')
        os.mkdir('bin')

    if not is_script_directory(name):
        die(f'{name} does not appear to be a valid script directory.')

    script_info_toml = toml.load(open(os.path.join(name, 'script_info.toml'), 'r'))

    if not is_script_info_valid(script_info_toml):
        die(f"{name}'s {SCRIPT_CONFIG_FILENAME} file is not valid.")

    script_info = script_info_toml['script']
    entry = script_info['entry']
    dest_name = script_info['name']

    dest_path = os.path.abspath(os.path.join('bin', dest_name))

    info('Creating symlink...')
    os.symlink(os.path.abspath(os.path.join(name, entry)), dest_path)

    info('Setting executable...')
    link_st = os.stat(dest_path)
    os.chmod(dest_path, link_st.st_mode | stat.S_IEXEC)

    success(f'Linked {name} to bin/.')


@click.command()
@click.argument('name')
def unlink(name):
    """
    Removes a symbolic link from a scriptoid to ./bin/
    """
    info(f'Unlinking {name} from bin/...')

    if not os.path.isdir('bin'):
        die('No scripts are currently linked.')

    if not os.path.islink(os.path.join('bin', name)):
        die(f'Script {name} does not appear to be linked.')

    info('Removing link...')
    os.remove(os.path.join('bin', name))

    success(f'Removed link to {name}')


for command in [
    new,
    link,
    unlink
]:
    cli.add_command(command)

if __name__ == '__main__':
    cli()
