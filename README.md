# Lunno

**Lunno** is a small, type-explicit, purely functional scripting
language built in Go.

Everything in Lunno is an **expression**, variables are **immutable**,
and functions are **first-class**.

---

## Features
- **Purely functional**: no mutable variables by default,  recursion replaces loops.
- **Type-explicit**: all bindings have explicit types, improving clarity and safety.
- **First-class functions**: functions can be passwd around, returned, and stored in variables.
- **Small but expressive**: a minimal set of keywords allows you to write complex programs.

---

## Example Lunno Script

```lunno
import io

print("Hello, world")
```

---