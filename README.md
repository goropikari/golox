# TLPS (Toy Language Processing System)

TLPS is a toy language processing system for me to study.
TLPS adapts [Off-side rule](https://en.wikipedia.org/wiki/Off-side_rule) like python.



```
git clone https://github.com/goropikari/tlps
cd tlps
docker build -t tlps .
docker run -it tlps
```

Architecture is based on jlox (tree-walk interpreter) by [munificent/craftinginterpreters](https://github.com/munificent/craftinginterpreters).
