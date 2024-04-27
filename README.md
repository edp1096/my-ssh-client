Password input SSH Client for my personal usage

## Usage
```sh
.\ssh-client.exe -l my_account -passwd my_password ip_adress
```

## Insert key not work on windows - need asking to `go-tty` author
```go
switch vk {
case 0x12: // menu
    if kr.controlKeyState&leftAltPressed != 0 {
        tty.readNextKeyUp = true
    }
    return 0, nil
case 0x21: // page-up
    tty.rs = []rune{0x5b, 0x35, 0x7e}
    return rune(0x1b), nil

...

case 0x28: // down
    tty.rs = []rune{0x5b, 0x42}
    return rune(0x1b), nil
case 0x2D: // Insert
    tty.rs = []rune{0x5b, 0x32, 0x7e}
    return rune(0x1b), nil
case 0x2e: // delete
    tty.rs = []rune{0x5b, 0x33, 0x7e}
    return rune(0x1b), nil

...

}
```

## Source
https://gist.github.com/atotto/ba19155295d95c8d75881e145c751372