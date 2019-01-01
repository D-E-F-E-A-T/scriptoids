import os
import shutil
import unittest
from pathlib import Path

import toml
from semver import VersionInfo

from scriptoids.environment import ScriptoidEnvironment


class TestScriptoidEnvironment(unittest.TestCase):
    def setUp(self):
        self.scriptoid_home = Path('.', 'test_scriptoid_root').absolute()
        self.environment = ScriptoidEnvironment(str(self.scriptoid_home))

    def tearDown(self):
        if self.scriptoid_home.exists():
            shutil.rmtree(str(self.scriptoid_home))

    def _createFooScriptoid(self):
        Path(self.scriptoid_home, 'foo').mkdir()
        with open(str(Path(self.scriptoid_home, 'foo', self.environment.scriptoid_definition_filename)), 'w') as fp:
            toml.dump({
                'scriptoid': {
                    'name': 'foo',
                    'version': '0.0.0',
                    'entry_file': 'foo.sh',
                }
            }, fp)

    def test_init_required_directories(self):
        self.assertTrue(os.path.isdir(self.environment.home_directory))
        self.assertTrue(os.path.isdir(self.environment.bin_directory))

    def test_scriptoid_exists(self):
        Path(self.scriptoid_home, 'foo').mkdir()

    def test_scriptoid_exists_nonexistent(self):
        Path(self.scriptoid_home, 'foo', self.environment.scriptoid_definition_filename).touch()

        self.assertTrue(self.environment.scriptoid_exists('foo'))
        self.assertFalse(self.environment.scriptoid_exists('bar'))

    def test_get_scriptoid(self):
        self._createFooScriptoid()

        scriptoid = self.environment.get_scriptoid('foo')
        self.assertEqual('foo', scriptoid.name)
        self.assertEqual(VersionInfo(0, 0, 0), scriptoid.version)
        self.assertEqual('foo.sh', scriptoid.entry_file)
        self.assertEqual('', scriptoid.description)
        self.assertEqual([], scriptoid.path_dependencies)
        self.assertEqual([], scriptoid.script_dependencies)

    def test_get_scriptoid_nonexistent(self):
        self.assertRaises(NameError, self.environment.get_scriptoid('bar'))

    def test_link_scriptoid(self):
        self._createFooScriptoid()
        self.environment.link_scriptoid('foo')
