package fireball

type Config struct {
    Host string
    Port int
    Log  Log
}

type Log struct {
    Request  bool
}
