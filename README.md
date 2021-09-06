# TLPS (Toy Language Processing System)

TLPS is a toy language processing system for me to study.
TLPS adopts [Off-side rule](https://en.wikipedia.org/wiki/Off-side_rule) like python.



```
git clone https://github.com/goropikari/tlps
cd tlps
docker build -t tlps .
docker run -it tlps
```

Architecture is based on jlox (tree-walk interpreter) by [munificent/craftinginterpreters](https://github.com/munificent/craftinginterpreters).


```
// declare variable
var x = 10; // terminal ';' is optional

// if statement
if expr:
  statements
elseif expr:
  statements
else:
  statements

// while loop
while expr:
  statements

// for loop
for var i = 0; i < 5; i = i + 1:
  print i
```

# Todo

- [ ] add tests
- [ ] escape string
- [ ] detect IndentationError
- [ ] import another tlps file
- [ ] support IO

