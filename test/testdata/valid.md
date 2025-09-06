# Valid Mermaid Diagrams

This file contains valid mermaid diagrams for testing.

## Simple Graph

```mermaid
graph TD
    A[Start] --> B{Decision}
    B -->|Yes| C[OK]
    B -->|No| D[End]
```

## Sequence Diagram

Some text between diagrams.

```mermaid
sequenceDiagram
    participant Alice
    participant Bob
    Alice->>Bob: Hello Bob, how are you?
    Bob-->>Alice: Great!
```

## Flowchart

```mermaid
flowchart LR
    A[Hard edge] -->|Link text| B(Round edge)
    B --> C{Decision}
    C -->|One| D[Result one]
    C -->|Two| E[Result two]
```

## Regular Code Block

This is not a mermaid block:

```javascript
console.log("Hello, world!");
```

End of document.