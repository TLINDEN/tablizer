# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) and this project adheres to [Semantic Versioning](http://semver.org).

## [v1.0.11](https://github.com/TLINDEN/tablizer/tree/v1.0.11) - 2022-10-19

[Full Changelog](https://github.com/TLINDEN/tablizer/compare/v1.0.10...v1.0.11)

### Added

- Added CI job golinter to regularly check for common mistakes.

- Added YAML output mode.

- Added more unit tests, we're over 95% in the lib module.

### Changed

- do  not use any global  variables anymore, makes the  code easier to
  maintain, understand and test

- using io.Writer  in print* functions, which is easier  to test, also
  re-implemented the print tests.

- replaced go-str2duration with my own implementation `duration2int()`.



## [v1.0.10](https://github.com/TLINDEN/tablizer/tree/v1.0.10) - 2022-10-15

[Full Changelog](https://github.com/TLINDEN/tablizer/compare/v1.0.9...v1.0.10)

### Added

- Added various sort modes: sort by time, by duration, numerical (-a -t -i)

- Added possibility to modify sort order to descending (-D)

- Added support  to specify a regexp in column  selector -c, which can
  also be mixed with numerical column spec
  
- More unit tests

### Fixed

- Column specification  allowed to specify duplicate  columns like `-c
  1,2,1,2` unchecked. Now this list will be deduplicated before use.



## [v1.0.9](https://github.com/TLINDEN/tablizer/tree/v1.0.9) - 2022-10-14

[Full Changelog](https://github.com/TLINDEN/tablizer/compare/v1.0.8...v1.0.9)

### Added

- Added Changelog, Contribution guidelines and no COC.

### Changed

- some minor changes to satisfy linter.



## [v1.0.8](https://github.com/TLINDEN/tablizer/tree/v1.0.8) - 2022-10-13

[Full Changelog](https://github.com/TLINDEN/tablizer/compare/v1.0.7...v1.0.8)

### Added

- Added sort support with the new parameter -k (like sort(1)).



## [v1.0.7](https://github.com/TLINDEN/tablizer/tree/v1.0.7) - 2022-10-11

[Full Changelog](https://github.com/TLINDEN/tablizer/compare/v1.0.6...v1.0.7)

### Added

- Added pattern highlighting support.

- Added more unit tests.

### Fixed

- Fixed extended more output in combination with -c.

- Fixed issue #4, the version string was missing.



## [v1.0.6](https://github.com/TLINDEN/tablizer/tree/v1.0.6) - 2022-10-05

[Full Changelog](https://github.com/TLINDEN/tablizer/compare/v1.0.5...v1.0.6)

### Added

- Added documentation about regexp syntax in the manpage.

- Added more unit tests.

### Changed

- Rewrote the input parser.

- Some more refactoring work has been done.



## [v1.0.5](https://github.com/TLINDEN/tablizer/tree/v1.0.5) - 2022-10-05

[Full Changelog](https://github.com/TLINDEN/tablizer/compare/v1.0.4...v1.0.5)

### Added

- A  new option has been  added: --invert-match -v which  behaves like
  the same option in grep(1): it inverts the pattern match.
  
- A few more unit tests have been added.

### Fixed

- Pattern  matching did  not work, because  the (new)  help subcommand
  lead  to  cobra  taking  care  of  the  first  arg  to  the  program
  (argv[1]).  So now  there's a  new parameter  -m which  displays the
  manpage and no more subcommands.
  


## [v1.0.4](https://github.com/TLINDEN/tablizer/tree/v1.0.4) - 2022-10-04

[Full Changelog](https://github.com/TLINDEN/tablizer/compare/v1.0.3...v1.0.4)

### Added

- Development version of the compiled binary now uses git vars
  in addition to program version.
  
- Added  an option to display  the manual page (compiled  in) as text:
  --help, for cases where a user just installed the binary.
  
### Changed

- Fixed go module namespace.



## [v1.0.3](https://github.com/TLINDEN/tablizer/tree/v1.0.3) - 2022-10-03

[Full Changelog](https://github.com/TLINDEN/tablizer/compare/v1.0.2...v1.0.3)

### Added

- Added a new output mode: shell mode, which allows the user
  to use the output in a shell eval loop to further process
  the data.
  
### Changed

- More refactoring work has been done.



## [v1.0.2](https://github.com/TLINDEN/tablizer/tree/v1.0.2) - 2022-10-02

[Full Changelog](https://github.com/TLINDEN/tablizer/compare/v1.0.1...v1.0.2)

### Added

- Added some basic unit tests.

### Changed

- Code has been refactored to be more efficient.

- Replaced table generation code with Tablewriter.





## [v1.0.1](https://github.com/TLINDEN/tablizer/tree/v1.0.1) - 2022-09-30

[Full Changelog](https://github.com/TLINDEN/tablizer/compare/v1.0.0...v1.0.1)

### Added

- Added a unix manual page.

- Added release builder to Makefile

### Changed

- Various minor fixes.



## [v1.0.0](https://github.com/TLINDEN/tablizer/tree/v1.0.0) - 2022-09-28

[Full Changelog](https://github.com/TLINDEN/tablizer/compare/02a64a5c3fe4220df2c791ff1421d16ebd428c19...v1.0.0)

Initial release.
