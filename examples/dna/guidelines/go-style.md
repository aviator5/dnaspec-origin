# Go Style Guide

This document defines Go coding style conventions and best practices.

## Naming Conventions

**Files**: Use `snake_case.go` (e.g., `handler_create.go`, `handler_list.go`)

**Packages**: Use lowercase, no underscores (e.g., `rest`, `companies`, `contacts`, `database`)

**Interfaces**: Use singular, role-based names (e.g., `CompaniesLister`, `TokenProvider`)

## Import Ordering

Organize imports in three groups separated by blank lines:
1. Standard library
2. External dependencies
3. Local packages (project-specific prefix)

## Formatting

- Max line length: 140 characters

## Function Grouping and Ordering

- Functions should be sorted in rough **call order**
- Functions should be **grouped by receiver**
- Order: `struct/const/var` → `newXYZ()/NewXYZ()` → receiver methods → utility functions

```go
// DON'T: Random ordering
func (s *something) Cost() {
  return calcCost(s.weights)
}

type something struct{ ... }

func calcCost(n []int) int {...}

func (s *something) Stop() {...}

func newSomething() *something {
    return &something{}
}

// DO: Logical grouping
type something struct{ ... }

func newSomething() *something {
    return &something{}
}

func (s *something) Cost() {
  return calcCost(s.weights)
}

func (s *something) Stop() {...}

func calcCost(n []int) int {...}
```

## Reduce Nesting

- Handle error cases/special conditions first
- Return early or continue the loop
- Minimize deeply nested code blocks

```go
// DON'T: Deep nesting
for _, v := range data {
  if v.F1 == 1 {
    v = process(v)
    if err := v.Call(); err == nil {
      v.Send()
    } else {
      return err
    }
  } else {
    log.Printf("Invalid v: %v", v)
  }
}

// DO: Handle errors/special cases first, return/continue early
for _, v := range data {
  if v.F1 != 1 {
    log.Printf("Invalid v: %v", v)
    continue
  }

  v = process(v)
  if err := v.Call(); err != nil {
    return err
  }
  v.Send()
}
```

## Unnecessary Else

- If a variable is set in both branches, replace with single if
- Set the default value, then override in the if block if needed

```go
// DON'T: Unnecessary else
var a int
if b {
  a = 100
} else {
  a = 10
}

// DO: Set default, override if needed
a := 10
if b {
  a = 100
}
```

## Style

**Local Variable Declarations**:
- Use `:=` when setting explicit value
- Use `var` for zero values (makes intent clear)

```go
s := "foo"              // DO: Explicit value
var filtered []int      // DO: Zero value (nil slice)
```

**nil is a Valid Slice**:
- Return `nil` instead of `[]T{}`
- Check emptiness with `len(s) == 0`, not `s == nil`
- Zero value slices are usable immediately

```go
// DON'T
if x == "" { return []int{} }
if s == nil { /* empty check */ }

// DO
if x == "" { return nil }
if len(s) == 0 { /* empty check */ }
```

**Avoid Naked Parameters**:
- Add C-style comments for unclear parameters
- Better: Use custom types instead of bool

```go
// DON'T: Unclear what true/true means
printInfo("foo", true, true)

// DO: Add comments
printInfo("foo", true /* isLocal */, true /* done */)

// BEST: Use custom types
type Region int
const (Local Region = iota; Remote)
printInfo("foo", Local, StatusDone)
```

**Use Raw String Literals**:
- Avoid hand-escaped strings - use backticks

```go
// DON'T: Hard to read
wantError := "unknown name:\"test\""

// DO: Raw string literal
wantError := `unknown name:"test"`
```

**Use var for Zero Value Structs**:
- Differentiates zero-valued from non-zero fields

```go
// DON'T
user := User{}

// DO
var user User
```

**Initializing Maps**:
- Use `make()` for empty maps (not `map[T1]T2{}`)
- Use map literals for fixed initial values

```go
// DON'T: Visually similar to declaration
m1 := map[string]int{}

// DO: Visually distinct, allows capacity hint
m1 := make(map[string]int)
m2 := make(map[string]int, 10)  // With capacity hint

// DO: Fixed initial values
m3 := map[string]int{
  "a": 1,
  "b": 2,
}
```

}
```

### UUID Type Safety

- Services and repositories MUST use `uuid.UUID` type for UUID parameters (NOT `string`)
- Handlers parse string UUIDs from URLs/requests and pass `uuid.UUID` to services
- Database queries use `.String()` method to convert `uuid.UUID` to string
- Benefits: Type safety, early validation, clear API contracts, prevents invalid UUID propagation

### Copy Slices and Maps at Boundaries

- Slices and maps contain pointers to underlying data
- When receiving: Copy to prevent external modifications affecting internal state
- When returning: Copy to prevent external modifications of internal state
- Use `make()` + `copy()` for slices, iterate and copy for maps

```go
// DO: Defensive copy when receiving
func (d *Driver) SetTrips(trips []Trip) {
  d.trips = make([]Trip, len(trips))
  copy(d.trips, trips)
}

// DO: Defensive copy when returning
func (s *Stats) Snapshot() map[string]int {
  result := make(map[string]int, len(s.counters))
  for k, v := range s.counters {
    result[k] = v
  }
  return result
}
```

### Handle Type Assertion Failures

- Single-value type assertion panics on incorrect type
- Always use "comma ok" idiom for type safety

```go
// DON'T: Panics if i is not a string
t := i.(string)

// DO: Safely handle incorrect type
t, ok := i.(string)
if !ok {
  // handle error gracefully
}
```

### Don't Panic

- Production code must avoid panics (they cause cascading failures)
- Return errors and let caller decide how to handle
- Exception: Program initialization failures (e.g., `template.Must()`)
- In tests: Use `t.Fatal()` or `t.FailNow()` instead of panic

```go
// DON'T: Panic in production code
func run(args []string) {
  if len(args) == 0 {
    panic("an argument is required")
  }
}

// DO: Return error and let caller handle
func run(args []string) error {
  if len(args) == 0 {
    return errors.New("an argument is required")
  }
  return nil
}

// Exception: OK to panic during initialization
var _statusTemplate = template.Must(template.New("name").Parse("_statusHTML"))
```

### Avoid Embedding Types in Public Structs

- Embedding leaks implementation details and inhibits type evolution
- Embedded types become part of public API (breaking change to remove/modify)
- Write explicit delegation methods instead

```go
// DON'T: Embed in public struct (exposes AbstractList as public API)
type ConcreteList struct {
  *AbstractList  // Now part of public API forever
}

// DO: Use composition with explicit delegation
type ConcreteList struct {
  list *AbstractList  // Private field
}

func (l *ConcreteList) Add(e Entity) {
  l.list.Add(e)  // Explicit delegation
}

func (l *ConcreteList) Remove(e Entity) {
  l.list.Remove(e)
}
```

**Why**: Adding/removing methods from embedded types is a breaking change. Explicit delegation hides implementation details and preserves flexibility for future changes.

### Don't Fire-and-Forget Goroutines
- Do not leak goroutines - they cost memory and CPU
- Every goroutine must have predictable stop time OR a way to signal it to stop
- Always provide a way to wait for goroutine completion
- Never spawn goroutines in `init()`

```go
// DON'T: No way to stop this goroutine
go func() {
  for {
    flush()
    time.Sleep(delay)
  }
}()

// DO: Use channels to signal stop and wait for completion
var (
  stop = make(chan struct{})
  done = make(chan struct{})
)
go func() {
  defer close(done)
  ticker := time.NewTicker(delay)
  defer ticker.Stop()
  for {
    select {
    case <-ticker.C:
      flush()
    case <-stop:
      return
    }
  }
}()
// Later: close(stop); <-done

// DO: Use sync.WaitGroup for multiple goroutines
var wg sync.WaitGroup
for i := 0; i < N; i++ {
  wg.Add(1)
  go func() {
    defer wg.Done()
    // work
  }()
}
wg.Wait()

// DON'T: Spawn goroutines in init()
func init() {
  go doWork()  // NO!
}

// DO: Expose object managing goroutine lifecycle
type Worker struct {
  stop chan struct{}
  done chan struct{}
}

func NewWorker() *Worker {
  w := &Worker{stop: make(chan struct{}), done: make(chan struct{})}
  go w.doWork()
  return w
}

func (w *Worker) Shutdown() {
  close(w.stop)
  <-w.done
}
```

## Performance

**Prefer strconv over fmt**:
- For primitive-to-string conversions, `strconv` is faster than `fmt`

```go
// DON'T: Slower with more allocations
s := fmt.Sprint(rand.Int())        // 143 ns/op, 2 allocs/op

// DO: Faster with fewer allocations  
s := strconv.Itoa(rand.Int())      // 64.2 ns/op, 1 allocs/op
```

**Avoid Repeated String-to-Byte Conversions**:
- Convert fixed strings to bytes once and reuse

```go
// DON'T: Convert on every iteration
for i := 0; i < b.N; i++ {
  w.Write([]byte("Hello world"))   // 22.2 ns/op
}

// DO: Convert once, reuse
data := []byte("Hello world")
for i := 0; i < b.N; i++ {
  w.Write(data)                    // 3.25 ns/op
}
```

**Specify Container Capacity**:
- Provide capacity hints to reduce allocations from resizing

```go
// Maps: Approximate capacity (not guaranteed)
files, _ := os.ReadDir("./files")
m := make(map[string]os.DirEntry, len(files))  // Capacity hint

// Slices: Exact capacity (guaranteed)
data := make([]int, 0, size)  // length=0, capacity=size
for k := 0; k < size; k++ {
  data = append(data, k)      // No allocations until len > capacity
}

// DON'T: Without capacity
data := make([]int, 0)        // Multiple allocations as slice grows

// Performance impact: 2.48s vs 0.21s in benchmarks
```

## Error Handling

**Best Practices**:
- Use error wrapping with `fmt.Errorf("%w", err)`
- Use `errors.Is(err, target)` for error comparison (NOT `err == target`)
- Use `errors.As(err, &target)` for error type extraction (NOT type assertions)

### Error Types

Choose error type based on two factors:

| Error matching needed? | Message type | Use |
|------------------------|--------------|-----|
| No  | static  | `errors.New()` |
| No  | dynamic | `fmt.Errorf()` |
| Yes | static  | top-level `var` with `errors.New()` |
| Yes | dynamic | custom error type (struct) |

```go
// Static error, matching needed: use sentinel error
var ErrCompanyNotFound = errors.New("company not found")

// Dynamic error, matching needed: use custom error type
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

// Check with errors.Is() and errors.As()
if errors.Is(err, ErrCompanyNotFound) { ... }
var validationErr *ValidationError
if errors.As(err, &validationErr) { ... }
```

### Error Wrapping

Three options for propagating errors:
- Return as-is if no additional context needed
- Add context with `fmt.Errorf(..., %w, err)` to allow caller to match underlying error
- Add context with `fmt.Errorf(..., %v, err)` to obfuscate underlying error

**Keep context succinct** - avoid "failed to" phrases that pile up:

```go
// DON'T: Verbose error context
return fmt.Errorf("failed to create new store: %w", err)
// Result: "failed to x: failed to y: failed to create new store: the error"

// DO: Succinct error context
return fmt.Errorf("new store: %w", err)
// Result: "x: y: new store: the error"
```

### Error Naming

**Sentinel errors**: Use `Err` prefix for exported, `err` prefix for unexported
**Custom error types**: Use `Error` suffix

```go
var (
    ErrBrokenLink   = errors.New("link is broken")    // exported
    ErrCouldNotOpen = errors.New("could not open")    // exported
    errNotFound     = errors.New("not found")         // unexported
)

type NotFoundError struct { File string }             // exported
type resolveError struct { Path string }              // unexported
```

### Handle Errors Once

Each error should be handled **only once**. Don't log AND return - let caller decide.

```go
// DON'T: Log and return (causes duplicate logs)
if err != nil {
    log.Printf("Could not get user: %v", err)
    return err  // Caller will likely log too
}

// DO: Either wrap and return
if err != nil {
    return fmt.Errorf("get user %q: %w", id, err)
}

// DO: Or log and degrade gracefully (no return)
if err := emitMetrics(); err != nil {
    log.Printf("Could not emit metrics: %v", err)
    // Continue execution
}

// DO: Or match specific errors and handle appropriately
if err != nil {
    if errors.Is(err, ErrUserNotFound) {
        tz = time.UTC  // Use default
    } else {
        return fmt.Errorf("get user %q: %w", id, err)
    }
}
```

## Code Quality & Linting

**Tool**: `golangci-lint v2.x`

**Best Practices**:
- Run linters locally before committing
- Run linters periodically during development to catch issues early
- Fix linter warnings before creating pull requests
- All code must pass linting in CI/CD pipeline
