# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Additions
- Multidimensional arrays for custom types (structs)
- Print functions for multidimensional arrays (both basic and struct types)
- Functions as first-class objects
- Build tool (instead of feeding several files to the interpreter, feed it a `main.cx` that pulls every required file in a project)

## [0.2.0] - 2017-12-23
### Added
- Methods. For example:

```
package main

type Point struct {
    x i32
    y i32
}

func (p Point) PrintX {
    i32.print(p.x)
}

func main () {
    var myPoint Point
    myPoint.PrintX()
}
```

### Changed
- `time.now()` replaced by `time.Unix()`, `time.UnixMilli()` and `time.UnixNano()`, which return a Unix timestamp in seconds, milliseconds and nanoseconds, respectively.

### Fixed
- The programmer is no longer forced to write a pair of parentheses in the place of a function's output, even if the function is not returning anything. For example, in older versions:

```
func main () () {}
```

Now, you can write:

```
func main () {}
```

## [0.1.0] - 2017-12-22
### Added
- Multidimensional arrays are now supported for basic types (byte, str, bool, i32, i64, f32, f64)
