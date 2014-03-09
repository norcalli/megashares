#Megashares in Go [![GoDoc](https://godoc.org/github.com/norcalli/megashares?status.png)](https://godoc.org/github.com/norcalli/megashares)

Megashares API written in Go.

##Installation:
```
go get github.com/norcalli/megashares
```

##Example:
  ###Shitty client:
There is a [shitty example]("https://github.com/norcalli/megashares/example/shittyclient.go") (aptly named `shittyclient.go`) in the repository. It should serve as a starting point for figuring out how to use the API. Usage:

This will print the results of the query:
```
shittyclient -u "username" -p "password" -q "Firefly s01e01"
```

This will download result `n` of the query:
```
shittyclient -u "username" -p "password" -q "Firefly s01e01" -n 0
```

###Snippet:

Here is a snippet that outlines basic usage.
```
m := megashares.New()
if err := m.Login(username, password); err != nil {
  log.Fatal("Couldn't login!")
}
entries, _ := m.SearchEntries(query)
for i, entry := range entries {
  fmt.Printf("%d: %s\n", i, entry.Filename)
}
entry := entries[0] // Choose first entry.
fmt.Println(entry.Filename, ":", entry.Url)
```

##Todo:
- Add multiple page search support.
- Clean up filenames better, they sometimes are malformed.
