#!/usr/bin/env python3
import argparse
import os
import re
import subprocess
import sys
from typing import List, Optional, Iterator


class Git:
    def __init__(self) -> None:
        self.__local_tags: Optional[List[str]] = None
        self.__remote_tags: Optional[List[str]] = None

    def _capture(self, args: List[str], check: bool = True) -> str:
        result = subprocess.run(args, check=check, stdout=subprocess.PIPE)
        return result.stdout.strip().decode()

    def current_branch(self) -> str:
        return self._capture(['git', 'rev-parse', '--abbrev-ref', 'HEAD'])

    def local_tags(self) -> List[str]:
        if self.__local_tags is None:
            output = self._capture(['git', 'show-ref', '--tags'], check=False)
            self.__local_tags = self._tags_from_refs(output)
        return self.__local_tags

    def remote_tags(self) -> List[str]:
        if self.__remote_tags is None:
            output = self._capture(['git', 'ls-remote', '--tags', 'origin'])
            self.__remote_tags = self._tags_from_refs(output)
        return self.__remote_tags

    def _tags_from_refs(self, refs: str) -> List[str]:
        lines = re.findall(r'refs/tags/([^^\n]+)', refs)
        return list(set(lines))

    def is_index_clean(self) -> bool:
        result = subprocess.run(['git', 'diff-index', '--quiet', 'HEAD', '--'], check=False, stdout=subprocess.PIPE)
        return result.returncode == 0

    def create_commit(self, message: str) -> None:
        subprocess.run(['git', 'commit', '--allow-empty', '-m', message], check=True)

    def create_annotated_tag(self, name: str, message: str) -> None:
        subprocess.run(['git', 'tag', '-a', name, '-m', message], check=True)

    def push(self) -> None:
        subprocess.run(['git', 'push', '--follow-tags'], check=True)


class InvalidVersion(ValueError):
    pass



class Version:

    _RE = re.compile(r'^(\d+).(\d+).(\d+)(-(rc)){0,1}(.(\d+)){0,1}$')

    @classmethod
    def from_string(cls, value: str) -> 'Version':
        match = cls._RE.match(value)
        if match is None:
            raise InvalidVersion(f'incorrect version: {value}')
        major, minor, patch, _, pre_type, _, pre_version = match.groups()
        return cls(int(major), int(minor), int(patch), pre_type, None if pre_version is None else int(pre_version))

    @staticmethod
    def sorted(versions: Iterator['Version']) -> List['Version']:
        return sorted(versions, key=lambda e: e.sort_key())

    def __init__(self, major: int, minor: int, patch: int, pre_type: Optional[str], pre_version: Optional[int]) -> None:
        self.major = major
        self.minor = minor
        self.patch = patch
        self.pre_type = pre_type
        self.pre_version = pre_version

    def sort_key(self) -> List:
        return [self.major, self.minor, self.patch, self.pre_type, self.pre_version]

    def next_major(self) -> 'Version':
        return Version(self.major + 1, 0, 0, None, None)

    def next_minor(self) -> 'Version':
        return Version(self.major, self.minor + 1, 0, None, None)

    def next_patch(self) -> 'Version':
        return Version(self.major, self.minor, self.patch + 1, None, None)

    def next_pre_version(self) -> 'Version':
        if self.pre_version is None:
            raise Exception(f'cannot compute the next pre-version for {self} since it is not a pre-version')
        return Version(self.major, self.minor, self.patch, self.pre_type, self.pre_version + 1)

    def is_pre_release(self) -> bool:
        return self.pre_type is not None

    def as_release_candidate(self, index: Optional[int] = 0) -> 'Version':
        return Version(self.major, self.minor, self.patch, 'rc', index)

    def __str__(self) -> str:
        s = f'{self.major}.{self.minor}.{self.patch}'
        if self.pre_type is not None:
            s += f'-{self.pre_type}'
        if self.pre_version is not None:
            s += f'.{self.pre_version}'
        return s

    def __repr__(self) -> str:
        return f'<Version {self}>'


class Releases:

    class Error(Exception):
        pass

    def __init__(self, git: Git) -> None:
        self._git = git

    def versions(self) -> List[Version]:
        def gen() -> Iterator[Version]:
            for tag in self._git.local_tags():
                try:
                    yield Version.from_string(tag.lstrip('v'))
                except InvalidVersion:
                    pass
        return Version.sorted(gen())

    def latest_release(self) -> Version:
        versions = [v for v in self.versions() if not v.is_pre_release()]
        if not versions:
            raise self.Error('no releases found from tags')
        return versions[-1]

    def next_release_candidate(self) -> Version:
        '''Create a new RC if last version is a release, increment the rc index otherwise.'''

        versions = self.versions()
        if not versions:
            raise self.Error('no releases found from tags')
        latest_release = versions[-1]

        if latest_release.is_pre_release():
            return latest_release.next_pre_version()

        return latest_release.next_minor().as_release_candidate(0)

    def next_release(self) -> Version:
        return self.latest_release().next_minor()

    def create_release(self, version: Version) -> None:
        name = f'v{version}'

        if name in self._git.remote_tags():
            raise self.Error(f'the tag {name} exists in remote')

        self._git.create_commit(f"Release {name}")
        self._git.create_annotated_tag(name, f"Release {name}")


def process(git: Git, releases: Releases, action: str) -> None:
    if git.current_branch() != 'master':
        sys.exit('not on the master branch')

    if not git.is_index_clean():
        sys.exit('uncommited changes')

    print(f'Latest release: {releases.latest_release()}')

    if action == 'release-candidate':
        version = releases.next_release_candidate()
    elif action == 'release':
        version = releases.next_release()
    else:
        raise AssertionError()

    print(f'Creating commit and tag for new version: {version}')
    releases.create_release(version)

    print(f'Pushing to remote')
    git.push()


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument('action', choices=['release-candidate', 'release'])
    args = parser.parse_args()

    git = Git()
    releases = Releases(git)

    try:
        process(git, releases, args.action)
    except Releases.Error as err:
        sys.exit(f'\033[31m{err}\033[0m')


if __name__ == "__main__":
    main()
