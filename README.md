## About

A go module and a binary to get Russian Central Bank currency rates.

### Run
`go run github.com/itroot/cbrates/cmd/cbrates@latest`

```
+------------------------------------------------+
| Курсы валют ЦБ РФ на 26.04.2023                |
+------+-----+--------+----------------+---------+
| CODE | NUM | AMOUNT | NAME           | VALUE   |
+------+-----+--------+----------------+---------+
| USD  | 840 |      1 | Доллар США     | 81,5499 |
| EUR  | 978 |      1 | Евро           | 90,0332 |
| CNY  | 156 |      1 | Китайский юань | 11,7664 |
+------+-----+--------+----------------+---------+
```

### Use

```
go get github.com/itroot/cbrates
```

An example could be found [here](https://github.com/itroot/cbrates/blob/master/example_test.go).
