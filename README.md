# SLOC

![](https://img.shields.io/badge/languages-44-blue.svg) ![](https://img.shields.io/github/issues/vicjohnson1213/sloc.svg)

Count how many lines are in your source code, broken down by the file type.

Statistics include:

- Total lines
- Code lines
- Comment lines (line and block comments)
- Mixed lines (code and comments)
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

      Language  Files    Code  Comment  Mixed  Blank
         Batch      6      42        0      0      0
           CSS      3     284        1      0     18
  Coffeescript      7      88       11      4     30
          Html      6     433        0      0     10
          JSON    123   18133        0      0    189
    JavaScript   1988  161269    41616   4241  23627
      Markdown    123   18588        0      0   7174
         Shell      2      28        2      0     15
    Typescript     35    4161      522     45    614
           XML      2    1948       11      0    235

         Total   2295  204974    42163   4290  31912
```

```
$ sloc --exclude node_modules src/

    Language  Files  Code  Comment  Mixed  Blank
        JSON      4   739        0      0      3
  JavaScript     16   552        2      5    151
    Markdown      2    14        0      0      9

       Total     22  1305        2      5    163
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
- XAML
- XML

## License

[MIT](https://github.com/vicjohnson1213/sloc/blob/master/LICENSE)
