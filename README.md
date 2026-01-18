# snakefs
Trying to build a snake game to play on the terminal using Go



## Architectural Descisions

- we can keep track of the snake's whole body relative to each other, instead of having absolute values.

|   |    |  x (c)  |  x (b)   |   x (a)   |

- the head (a) would be a absolute value, however, (b) could be denoted as [-1,0]


also worthwhile to note that 2D is not required to denote this, but however, if I were to introduce a diagonal movement, then it would be something for me to add.
