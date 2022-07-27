Upd. : cli added

# deltagopher

Command line(in future) utility to get the delta between 2 file.
Includes signature, delta and patch(in future) commands.


# Usage

## generate a signature 
```shell
./deltagopher signature --block-size 3 --filename old.txt

```

## generate a delta
```shell
./deltagopher delta --sigfile signature.yml new.txt

```

# Installation

You need a Golang version 1.16 or greater installed on your PC.
1. To install dependencies, run from root directory:
```shell
go get
```

2. To build the CLI
```shell
go build
```

# Approach

**deltagopher** uses the [Rolling Hash](https://en.wikipedia.org/wiki/Rolling_hash) hash function to find a delta, it 
provides a way to efficiently traverse the sequence and calculate a new hash in O(1) time using the old hash. The total time
for calculating a delta is O(n <u>old file</u>) + O(n <u>changed file</u>), where n in a length of the file.

Here is how it's done:

1. Obtain a **Signature** from the first file. **Signature** represents a hashed chunks 
of the file with some metadata of chunks size, hashing algorythms and so on. Each chunk
hashed in 2 ways - _**weak**_ and _**strong**_. Weak hash is a **adler32** checksum, strong hash - **md5** hash.
```go
type Checksum struct {
	WeakCheaksum   uint32   `yaml:"weak"`
	StrongChecksum [16]byte `yaml:"strong"`
	Start          int      `yaml:"start"`
	End            int      `yaml:"end"`
}
```
2. **Signature** can be saved into a file (json/yml), and used to generate a **Delta**. **Delta** represents 
the difference between a new changed version of the file and the **Signature** of old file.
```go
type Delta struct {
    Inserted []*SingleDelta `yaml:"insert,omitempty"`
    Deleted  []*SingleDelta `yaml:"delete,omitempty"`
    Copied   []*SingleDelta `yaml:"copy,omitempty"`
}
type SingleDelta struct {
    WeakCheaksum   uint32    `yaml:"weak,omitempty"`
    StrongChecksum *[16]byte `yaml:"strong,omitempty"`
    Start          int       `yaml:"start"`
    End            int       `yaml:"end"`
    
    // For Inserted
    DiffBytes []byte         `yaml:"diff,omitempty"`
}
```
3. **Delta** is generated using the following algorithm: 
   1. Load the first chunk of the file into buffer length **Signature.blockSize**
   2. Generate the weak hash for the chunk
   3. Search for this hash in the **Signature.checksums** from old file
      1. If it exists in Signature - it's a potential candidate to be a hash of the same text as for an old file. 
      To see if it's exactly the text we had in the old file we compare strong hashes:
         1. if md5 hash old == md5 hash new => mathing, _to be copied_ from the old file.
         2. if md5 hash old != md5 hash new => consider this block as new and _to be inserted_
         3. If it's not presented in old hashes - this block of the file to be inserted to the Delta
         4. **Seek the reader** pointer to size of the chunk to skip the chunk we already processed.
      2. If not - start rolling across the file content using **Rolling Hash (adler32)**. On each step 
      generate a new hash of the sliding window (it moves byte by byte), search for it in old checksums and generate 
      a Delta on the fly.  
   4. If some of the old hashes weren't found in the new file at all - they are _to be deleted_.
