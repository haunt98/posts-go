# Real World Crypto 101

My notes when reading
[Real-World Cryptography](https://www.manning.com/books/real-world-cryptography)

**Hash** function convert from input to digest:

- Pre-image resistance: Given digest, can not find input
- Second pre-image resistance: Given input, digest, can not find another input
  produce same digest. Small change to input make digest big change.
- Collision resistance: Can not find 2 input produce same digest.

**MAC** aka Message Authentication Code produce from key, message to
authentication tag. **HMAC** is MAC using hash.

- A send B message with MAC (generate from message and A key).
- B double check message with MAC (generate from receive message and B key).
- A and B use same key.

```mermald
sequenceDiagram
    participant alice
    participant bob

    alice ->> bob: send username, password
    bob -->> alice: return alice|mac(private_key, alice)
    alice ->> bob: send alice|mac(private_key, alice)
    bob -->> alice: return OK
    alice ->> bob: send bob|mac(private_key, alice)
    bob -->> alice: return ERROR
```
