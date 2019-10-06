## restr 

Module for generating random strings based on regular expressions. 

It does not depend on third-party libraries.

### Usage

```go
import "restr"

main() {
    re := `\A(?P<help>\[\d\])[\w]+\s\d*\(\w\)(Xor)?>\09{2,5}(\x55|[^a0-9]\*).{2}[a-f][[:space:]]Uu$`
    fmt.Println(Rstr(re))

    // will print:
    // [0]uFT3knvKjUxoxpzKjT_vJncFi6mYbef6gS3RoCKAEIg 5568853724798647844345752909(W)>9999󲛩*f	Uu

    re = `^[a-z][a-z\-\d]{3,14}@[a-z\d]{2,15}\.[a-z]{2,5}$`
    fmt.Println(Rstr(re))

    // will print:
    // hltod4cxl@aq.gqdw
}
```

### TODO

Need to improve:
1) Add Markov Generator for Named Groups
2) Add arguments to exclude and include specific characters
3) Check syntactic completeness