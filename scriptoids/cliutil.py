import sys

from colorama import Style, Fore


def info(msg, use_format=True):
    """
    Prints an informational message.
    :param use_format: Whether or not to use ANSI formatting.
    :param msg: Message to display.
    """
    print(f'{Style.BRIGHT}{Fore.BLACK}  {msg}{Style.RESET_ALL}' if use_format else f'Info: {msg}')


def success(msg, use_format=True):
    """
    Prints a message indicating success.
    :param use_format: Whether or not to use ANSI formatting.
    :param msg: Message to display.
    """
    print(f'{Style.BRIGHT}{Fore.GREEN}✔ {Fore.BLACK}{msg}{Style.RESET_ALL}' if use_format else f'Success: {msg}')


def die(msg, exit_code=1, use_format=True):
    """
    Prints a message denoting a failure and exits with a given exit code.
    :param msg: Message to display.
    :param exit_code: Code to return when exiting.
    :param use_format: Whether or not to use ANSI formatting.
    """
    sys.stderr.write(
        f'{Style.BRIGHT}{Fore.RED}✘ An error occurred: {Fore.BLACK}{msg}{Style.RESET_ALL}\n' if use_format
        else f'An error occurred: {msg}')
    exit(exit_code)
