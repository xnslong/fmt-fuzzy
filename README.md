# fmt-fuzzy

can format any text into a well indented format, so that it will be friendly for reading. for example:

```
10.186.36.13(somebody@hostname:dir):#{"a":0,"b":{"c":"1013","d":"1049","e":"1453",},"f":"success"}
```

can be formatted into

```
10.186.36.13(
    somebody@hostname:dir
):#{
    "a":0,
    "b":{
        "c":"1013",
        "d":"1049",
        "e":"1453",
    },
    "f":"success"
}
```

# usage

```
$ cat input.txt | fmt-fuzzy
```

# Installation

```
go get github.com/xnslong/fmt-fuzzy
```