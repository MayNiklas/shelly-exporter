# -*- coding: utf-8 -*-
from setuptools import setup, find_packages

setup(
    name='shelly_exporter',
    version='1.0.0',
    url='',
    license='',
    author='MayNiklas',
    author_email='info@niklas-steffen.de',
    description='prometheus exporter for shelly plug s',
    package_dir={'': 'src'},
    packages=find_packages('src') + find_packages('test/src'),
    entry_points={
        'console_scripts': [
            'shelly_exporter=shelly_exporter:main',
        ],
    },

)
