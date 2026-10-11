---
date: 2026-09-22
---
## Variables

```python
# 型宣言不要
# num: int = 1
num = 1
name = 'Mike'
is_true = True

print(num, type(num)) # 1 <class 'int'>
print(name, type(name)) # Mike <class 'str'>
print(is_true, type(is_true)) # True <class 'bool'>

num = name
print(num, type(num)) # Mike <class 'str'>

# 型変換
str = '1'
num = int(str)
print(num, type(num)) # 1 <class 'int'>
```

### Print

```python
print('Hi', 'Mike', sep=',', end='') # Hi,Mike
print('Hi', 'Mike', sep=',', end='\n') # Hi,Mike
print('Hi', 'Mike', sep=',', end='.\n') # Hi,Mike.
```

### Number

```python
import math

print(2 + 2) # 4
print(5 - 2) # 3
print(5 * 6) # 30
print(50 - 5 * 6) # 20
print((50 - 5) * 6) # 270
print(8 / 5) # 1.6
print(type(1)) # <class 'int'>
print(type(1.6)) # <class 'float'>
print(type(.6)) # <class 'float'>
print(10 / 3) # 3.3333333333333335
print(10 // 3) # 3
print(10 % 3) # 1
print(5 ** 2) # 25

pie = 3.1415926
print(round(pie, 2)) # 3.14

result = math.sqrt(25)
print(result) # 5.0

result = math.log2(10)
print(result) # 3.321928094887362

print(help(math))
```

### String

```python
print('hello')
print("I don't know")
print('I don\'t know')
print('say "I don\'t know" ')
print("say \"I don't know\" ")
print('hello.\nHow are you')
print(r'C:\Program Files\name')

print("#####")
print("""\
line1
line2
line3\
""")
print("#####")

print('Hi.' * 3 + 'Mike.') # Hi.Hi.Hi.Mike.

print('py' + 'thon') # python
print('py''thon') # python

s = 'aaaaa'\
    'bbbbb'

print(s) # aaaaabbbbb
```

#### String のスライス

```python
word = 'python'
print(word[0]) # p
print(word[-1]) # n
print(word[0:2]) # py
print(word[:2]) # py
print(word[2:]) # thon

# word[0] = 'j' # ERROR
word = 'j' + word[1:]
print(word) # jython
print(len(word)) # 6
```

#### String のメソッド

```python
s = 'my name is Mike. Hi Mike.'

is_start = s.startswith('my')
print(is_start) # True

is_start = s.startswith('X')
print(is_start) # False

print(s.find('Mike')) # 11
print(s.rfind('Mike')) # 20
print(s.count('Mike')) # 2
print(s.capitalize()) # My name is mike. hi mike.
print(s.title()) # My Name Is Mike. Hi Mike.
print(s.upper()) # MY NAME IS MIKE. HI MIKE.
print(s.lower()) # my name is mike. hi mike.
print(s.replace('Mike', 'Nancy')) # my name is Nancy. Hi Nancy.
```

#### String の代入

```python
x = str(1)
print(type(x)) # <class 'str'>
print(type(str(3.14))) # <class 'str'>
print(type(str(True))) # <class 'str'>

a = 'a'
print(f'a is {a}')

x, y, z = 1, 2, 3
print(f'a is {x}, {y}, {z}')
print(f'a is {z}, {y}, {x}')

name = 'A'
family = 'B'
print(f'My name is {name} {family}. Watashi ha {family} {name}')
```

### List

```python
l = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
print(l)
print(l[0])
print(l[-1])
print(l[0:2])
print(l[2:])
print(l[:])
print(len(l))
print(type(l)) # <class 'list'>
print(l[::2]) # [1, 3, 5, 7, 9]
print(l[::-1]) # [10, 9, 8, 7, 6, 5, 4, 3, 2, 1]

a = ['a', 'b', 'c']
b = [1, 2, 3]
x = [a, b]
print(x) # [['a', 'b', 'c'], [1, 2, 3]]
print(x[0])
print(x[0][1]) # b
```

#### List の操作

```python
s = ['a', 'b', 'c', 'd', 'e', 'f', 'g']
s[0] = 'x'
print(s) # ['x', 'b', 'c', 'd', 'e', 'f', 'g']

s[2:5] = ['C', 'D', 'E']
print(s) # ['x', 'b', 'C', 'D', 'E', 'f', 'g']

s[2:5] = []
print(s) # ['x', 'b', 'f', 'g']

s[:] = []
print(s) # []

n = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
n.append(100)
print(n) # [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 100]

n.insert(0,200)
print(n) # [200, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 100]

n.pop()
print(n) # [200, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

n.pop(0)
print(n) # [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

del n[0]
print(n) # [2, 3, 4, 5, 6, 7, 8, 9, 10]
# del n

n.remove(2)
print(n) # [3, 4, 5, 6, 7, 8, 9, 10]

a = [1, 2, 3, 4, 5]
b = [6, 7, 8, 9, 10]
x = a + b
print(x) # [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

a += b
print(a) # [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

a = [1, 2, 3, 4, 5]
b = [6, 7, 8, 9, 10]
a.extend(b)
print(a) # [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
```

#### List のメソッド

```python
r = [1, 2, 3, 4, 5, 1, 2, 3]
print(r.index(3)) # 2
print(r.index(3, 3)) # 7
print(r.count(3)) # 2

r.sort()
print(r) # [1, 1, 2, 2, 3, 3, 4, 5]
r.sort(reverse=True)
print(r) # [5, 4, 3, 3, 2, 2, 1, 1]
r.reverse()
print(r) # [1, 1, 2, 2, 3, 3, 4, 5]

s = 'My name is Mike.'
to_split = s.split(' ')
print(to_split) # ['My', 'name', 'is', 'Mike.']

x = ' '.join(to_split)
print(x) # My name is Mike.

# print(help(list))
```

#### List のコピー

```python
i = [1, 2, 3, 4, 5]
j = i # 参照渡し (メモリのアドレスをjに入れる)
print('i = ', i)
print('j = ', j)

j[0] = 100
print('i = ', i) # i =  [100, 2, 3, 4, 5]
print('j = ', j) # j =  [100, 2, 3, 4, 5]
print(id(i[0])) # 4349599256
print(id(j[0])) # 4349599256

x = [1, 2, 3, 4, 5]
y = x.copy() # 値渡し
# y = x[:] # 値渡し
y[0] = 100
print('x = ', x) # x =  [1, 2, 3, 4, 5]
print('y = ', y) # y =  [100, 2, 3, 4, 5]

# Taxi上下車など
```

### Tuple

```python
t = (1, 2, 3, 4, 1, 2) # ()なくでもいい
print(t)
print(t[0])
print(t[-1])
print(t[2:5])
print(t.index(1))
print(t.index(1,1))
print(t.count(1))

# t[0] = 100 # Error タプルは値変更できない
t = ([1, 2, 3], [4, 5, 6])
t[0][0] = 100 # タプルの中のリストに対して値変更できる
print(t) # ([100, 2, 3], [4, 5, 6])

x = () # tuple
y = (1, 2, 3)
z = x + y # 足すことが可能
print(z) # (1, 2, 3)

x += y
print(x) # (1, 2, 3)
```

#### Tuple の Unpacking

```python
num_tuple = (10, 20)
print(num_tuple)

x, y = num_tuple # unpack tuple
print(x, y)

min, max = 0, 100 # unpack tuple
print(min, max)

i = 10
j = 20
i, j = j, i
print(i, j) # 20 10

# 三つの選択肢から二つ選ぶ
choose_from_two = ('A', 'B', 'C')
answer = []
answer.append('A')
answer.append('C')

print(choose_from_two)
print(answer)
```

### Dict

```python
d = {'x': 10, 'y': 20}
print(type(d)) # <class 'dict'>
print(d) # {'x': 10, 'y': 20}

print(d['x']) # 10
print(d['y']) # 20

d['x'] = 100
print(d) # {'x': 100, 'y': 20}

d['x'] = 'abc'
print(d) # {'x': 'abc', 'y': 20}

d['z'] = 200
print(d) # {'x': 'abc', 'y': 20, 'z': 200}

d[1] = 100
print(d) # {'x': 'abc', 'y': 20, 'z': 200, 1: 100}

x = dict(a=10, b=20)
print(x) # {'a': 10, 'b': 20}

y = dict([('a', 10), ('b', 20)])
print(y) # {'a': 10, 'b': 20}
```

#### Dict のメソッド

```python
d = {'x': 10, 'y': 20}
# print(help(d))

print(d.keys()) # dict_keys(['x', 'y'])
print(d.values()) # dict_values([10, 20])

d2 = {'x': 1000, 'j': 500}
d.update(d2) # Overwrite
print(d) # {'x': 1000, 'y': 20, 'j': 500}

print(d['x']) # 1000
print(d.get('x')) # 1000

print(d.pop('x')) # 1000 を取り出す
print(d) # {'y': 20, 'j': 500}

del d['y'] # 削除
print(d) # {'j': 500}

# del d
d = {'x': 10, 'y': 20}
print('x' in d) # True
```

#### Dict のコピー

```python
x = {'a': 10, 'b': 20}
y = x # 参照渡し
y['a'] = 1000

print(x) # {'a': 1000, 'b': 20}
print(y) # {'a': 1000, 'b': 20}

x = {'a': 10, 'b': 20}
y = x.copy() # 値渡し
y['a'] = 1000
print(x) # {'a': 10, 'b': 20}
print(y) # {'a': 1000, 'b': 20}

# 何かをすぐ取り出したい時
fruits = {
  'apple': 10,
  'banana': 20,
  'orange': 30
}

print(fruits['apple'])
```

### Set

```python
a = {1, 2, 3, 4, 4, 5, 6}
b = {2, 3, 3, 6, 7}
print(a) # {1, 2, 3, 4, 5, 6}
print(type(a)) # <class 'set'>
print(a - b) # {1, 4, 5}
print(b - a) # {7}
print(a & b) # {2, 3, 6}
print(a | b) # {1, 2, 3, 4, 5, 6, 7}
print(a ^ b) # {1, 4, 5, 7} 重複していないもの
```

#### Set のメソッド

```python
s = {1, 2, 3, 4, 5}
s.add(6)
print(s) # {1, 2, 3, 4, 5, 6}
s.remove(6)
print(s) # {1, 2, 3, 4, 5}
s.clear()
print(s) # set()
# print(help(s))

# 集合を探す
my_friends = {'A', 'C', 'D'}
A_friends = {'B', 'D', 'E'}
print(my_friends & A_friends) # {'D'}

f = ['apple', 'banana', 'apple', 'banana']
kind = set(f) # 型変換
print(kind) # {'banana', 'apple'}
```
