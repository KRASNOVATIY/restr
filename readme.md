## restr 

Module for generating random strings based on regular expressions. 

Exported Instances:
1) var MaxRepeat - limit for * and + literals
2) fun Rstr - construct random string by regexpr
3) fun RandomString - returns a random string from a sequence "array"
4) fun RegisterName - register data source for Named Capture Group "?P<name>"
5) type MarkovGen - Markov model text generator
    1) method NewMarkovGen - return new MarkovGen
    2) method ModelApply - expand the model with text-based content
    3) method Generate - generate text based on a model

### Usage

code:
```go
import "restr"

main() {
    re := `\A(?P<help>\[\d\])[\w]+\s\d*\(\w\)(Xor)?>\09{2,5}(\x55|[^a0-9]\*).{2}[a-f][[:space:]]Uu$`
    fmt.Println(Rstr(re))

    // something like an email
    re = `^[a-z][a-z\-\d]{3,14}@[a-z\d]{2,15}\.[a-z]{2,5}$`
    fmt.Println(Rstr(re))

    // something like an email with named groups
    // pay attention to the "r" group, the literal "dot" is not processed in it
    re := `Login: (?P<name>.{5,15})@(?P<dom>.+)\.(?P<r>[a-z]{2,5}); Password: (?P<passw>.{13,18})`
    fmt.Println(Rstr(re))

    // create Markov text generator
    markov := NewMarkovGen(4, []rune{' ', ','})
	markov.ModelApply("title1", "hi my name is norton , wona be frieds", 1)
    markov.ModelApply("title2", "the lion was so tickled by the idea of the mouse being able to help him that he lifted his paw and let him go", 1)
    // register the generator in the model for the "name" group
    RegisterName("name", markov.Generate(15))
    // register the built-in RandomString function in the model for the "dom" group
    RegisterName("dom", RandomString([]string{"yandex", "google", "mail", "proton"}))
    // register the built-in RandomString function in the model for the "r" group
    RegisterName("r", RandomString([]string{"com", "eu", "uk", "ru"}))
    // processing for the "r" group is registered in the model, but is not applied due to insufficient soft expression (the literal "dot" should be used instead)
    fmt.Println(Rstr(re))

    // it is possible to register your own methods for processing named groups
    RegisterName("passw", func() string { return strings.Repeat("secret", rand.Intn(4)) })
    fmt.Println(Rstr(re))
}
```

output:
```txt
1. [0]uFT3knvKjUxoxpzKjT_vJncFi6mYbef6gS3RoCKAEIg 5568853724798647844345752909(W)>9999󲛩*f	Uu
2. zvbe7ag0ukj282@rr54u.uzuio
3. Login: ͉ύѥˉϻȤj»˹Jϐɥ@ԳºɿȤύ̮͇̻ˣʇΙŝżʀɪήGƭӲƗσǀгѣӫǦƗȋʯłחǣϱɿƆ̀̅ɌȚ̈́ӶҒ̡לʘʣ˞ͪןΑɠɆ̦˫ȵəɉ҅юɀĄ}ɿӗɻѷȽҦŹ̬a.nb; Password: ŋז˰ȯՒȔԣӨʬʪïęĕĻь
4. Login: friedsotickledb@yandex.ubl; Password: ЖZɇȨӄǟҕŤʿнЌ҉ȱ$ϼǀב
5. Login: tohelphimynamei@mail.ejvgb; Password: secretsecretsecret
```

### TODO

Need to improve:
1) Add arguments to exclude and include specific characters
2) Check syntactic completeness