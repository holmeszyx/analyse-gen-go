
Generating Analyse SDK source code for events with a `toml` definition file.

## Usage

1. build cmd binary

```
go build -o als-gen cmd/*.go
```

2. writing the events definition file by `toml`. like `testdata/als.toml`

3. Executing

```bash

# in cmd dir

als-gen -c ../testdata/als.toml -t ../testdata/umeng-android-kt.tmpl -out-base-dir ../testdata  -o St.kt -use-pkg=1

```

## Template

Template funcs:

- firstCap

    coverting first letter as upper case.

- toFuncName

    Ouput string as Camel-Case

- safeHan

    Filtering chars which are no "number", "chinese" or one mark in "_-.+". And replace the filtered char with `"_"`。
    It is useful when export reporting events table.(Like `Umeng` event batch exporting template)