### CSquared

#### Basics

- Pure Functional Language
  - Variables are immutable.
  - Functions can be treated as variables.
  - Each function that does not depend on other will automatically run in a new thread.
  - Users can additionally create single line functions to run on GPU.
- Syntax
  - It is easy!
    - variables are initialized with a ":" to announce the data type
    - ```cgo
      x: int = add(a, b)
      printf("%d\n", x)
        ```
    - functions can be defined by:
    - ```cgo
      x: string = (a, b: int, c: string) {
        printf("%d  + %d = %d, and hello %s!", a, b, add(a, b), c) 
      }
        ```
    - Note that the functions change the arguments. Each argument is passed by reference.
    - Each independent function is detected by the program runner.
    - Operators do not exist in CSquared. Everything's a function!
    - Each file is then converted to a C code, which is run using gcc.
    - Features are underway.