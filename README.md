# bt-fixture

bt-fixture is a test-fixture management library for the Cloud Bigtable.

[![GoDoc](https://godoc.org/github.com/takashabe/bt-fixture?status.svg)](https://godoc.org/github.com/takashabe/bt-fixture)
[![CircleCI](https://circleci.com/gh/takashabe/bt-fixture.svg?style=shield)](https://circleci.com/gh/takashabe/bt-fixture)
[![Go Report Card](https://goreportcard.com/badge/github.com/takashabe/bt-fixture)](https://goreportcard.com/report/github.com/takashabe/bt-fixture)

Management the fixtures of the database for testing.

## Features

- Load the Cloud Bigtable columns fixture by a yaml file
- Before cleanup a table
- Types
    - int and float types in a fixture convert to big-endian 8-byte values

## Usage

- Import `bt-fixture` package and setup a fixture file

```go
package main

import (
  fixture "github.com/takashabe/bt-fixture"
)

func main() {
  // omit the error handling
  f, _ := fixture.NewFixture("test-project", "test-instance")
  f.Load("testdata/fixture.yaml")
}
```

## Fixture file format

```yaml
table: person
column_families:
  - family: d
    columns:
      - key: 1##a
        rows:
          name: a
          age: 10
          height: 140.0
        version: 2018-05-19 00:00:00 +09:00
      - key: 1##a
        rows:
          name: a
          abe: 11
          height: 140.0
        version: 2019-05-19 00:00:00 +09:00
  - family: e
    columns:
      - key: 2##b
        rows:
          name: b
          age: 20
          height: 180.0
        version: 2018-05-19 00:00:00 +09:00
```

- Example output:

```
$cbt read person
----------------------------------------
1##a
  d:abe                                    @ 2019/05/19-00:00:00.000000
    "\x00\x00\x00\x00\x00\x00\x00\v"
  d:age                                    @ 2018/05/19-00:00:00.000000
    "\x00\x00\x00\x00\x00\x00\x00\n"
  d:height                                 @ 2019/05/19-00:00:00.000000
    "@a\x80\x00\x00\x00\x00\x00"
  d:height                                 @ 2018/05/19-00:00:00.000000
    "@a\x80\x00\x00\x00\x00\x00"
  d:name                                   @ 2019/05/19-00:00:00.000000
    "a"
  d:name                                   @ 2018/05/19-00:00:00.000000
    "a"
----------------------------------------
2##b
  e:age                                    @ 2018/05/19-00:00:00.000000
    "\x00\x00\x00\x00\x00\x00\x00\x14"
  e:height                                 @ 2018/05/19-00:00:00.000000
    "@f\x80\x00\x00\x00\x00\x00"
  e:name                                   @ 2018/05/19-00:00:00.000000
    "b"
```
