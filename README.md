# SLOC

![](https://img.shields.io/badge/languages-43-blue.svg) ![](https://img.shields.io/github/issues/vicjohnson1213/sloc.svg)

Count how many lines are in your source code, broken down by the file type.

Statistics include:

- Total lines
- Code lines
- Comment lines (line and block comments)
- Empty lines

## Usage

### CLI

```
sloc [options] <file>|<directory>
```

#### Options

```
-e, --exclude <regex>          A regular expression for files to exclude.
-i, --include <regex>          A regular expression for files to include.
-f, --format <table|json|csv>  The ouput format for the counting results.
```

#### Examples

```
$ sloc src/

      Language  Files    Code  Comment  Blank
          JSON    123   18133        0    189
      Markdown    123   18588        0   7174
    JavaScript   1988  155095    52031  23627
          Html      6     433        0     10
           XML      2    1948       11    235
         Batch      6      42        0      0
    Typescript     35    4204      524    614
           CSS      3     284        1     18
  Coffeescript      7      92       11     30
         Shell      2      28        2     15

         Total   2295  198847    52580  31912
```

```
$ sloc --exclude node_modules src/

   Language  Files  Code  Comment  Blank
  JavaScript     16   557        2    151
        JSON      4   739        0      3
    Markdown      2    14        0      9

       Total     22  1310        2    163
```

### Supported Languages

- Assembly
- Bash
- Batch
- C
- C++
- C#
- Clojure
- Coffeescript
- CSS
- D
- Dart
- Erlang
- Golang
- Groovy
- Haskell
- Html
- Java
- JavaScript
- JSON
- Kotlin
- LESS
- Lisp
- Lua
- Make
- Markdown
- Objective-C
- Perl
- PHP
- Python
- R
- Ruby
- Rust
- Scala
- Scheme
- Swift
- SASS
- SCSS
- Shell
- SQL
- Typescript
- VimL
- Visual Basic
- XML

## License

[MIT](https://github.com/vicjohnson1213/sloc/blob/master/LICENSE)
