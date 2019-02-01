### Hackme challenge

Two repositories
1. hackmebackend
- simple golang backend
2. hackmefrontend
- the react frontend

In hackmefrontend do
```
yarn install
yarn build
```

copy the whole build directory into hackmebackend folder as static directory.

`cp -R hackmefrontend/build/* hackmebackend/static/`

Build the golang server

In hackmebackend directory

`go build main.go primes.go`

- will generate an executable called `main`

Still in hackemebackend do

`./main`

to view the page

`http://localhost:3500`

