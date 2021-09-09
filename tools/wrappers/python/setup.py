#!/usr/bin/env python3

from setuptools import find_packages, setup
setup(
    name='chevronlib',
    packages=find_packages(include=['chevron']),
    version='1.3.1',
    description='Chevron GPG Library Wrapper for Python3',
    author='Quanto',
    license='MIT',
    install_requires=[],
    setup_requires=['pytest-runner'],
    tests_require=['pytest==4.4.1'],
    test_suite='tests',
)