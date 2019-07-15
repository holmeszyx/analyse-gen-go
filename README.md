
Generating Analyse SDK source code for events with a `toml` definition file.

## Usage

1. build cmd binary

```
go build -o als-gen cmd
```

2. writing the events definition file by `toml`. like `testdata/als.toml`

3. Executing

```bash

# in cmd dir

als-gen -c ../testdata/als.toml -t ../testdata/umeng-android-kt.tmpl -out-base-dir ../testdata  -o St.kt -use-pkg=1

```