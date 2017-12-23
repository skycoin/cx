# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Additions
- Multidimensional arrays for custom types (structs)
- Print functions for multidimensional arrays (both basic and struct types)
- Methods
- Functions as first-class objects
- Build tool (instead of feeding several files to the interpreter, feed it a `main.cx` that pulls every required file in a project)

## [0.1.0] - 2017-12-22
### Added
- Multidimensional arrays are now supported for basic types (byte, str, bool, i32, i64, f32, f64)
