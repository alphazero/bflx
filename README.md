#BFLX
language specification for level extended brainfuck [rev alpha.0]


###program structure
A bflx program is a sequence of byte codes (with ASCII semantics) of minimum length of 1.
    
## memory model
  
- bflx maintains a list of data cell arrays.
- a data cell has unsigned 8-bit integer semantics.
- all bflx programs have at least one data array (index 0).
- data cell array index underflow semantics are the same as a circular buffer. 
- Index overflow will dynamically extend the given cell array.
- data cell arrays (levels) are navigated via cursor commands.
- bflx runtime maintains 10 indexed unsigned 8-bit integer registers (0-9)

## bytecodes

### data cursor commands
  
`^` : level up    

- move up to previous data array index.
- if already at level 0, then go to the last.
 
`v` : level down

  - move down to the next data array index.
  - if already at level max, allocate a new layer and move to that.
  - data array levels are stateful and maintain their own data index.
      
`<` : move back

 - decrement data index of current data array. 
 - if data index is 0, then move index to the last position (per circular buffer semantics).
     
`>` : move forward

 - increment the data index of current data array.
 - if index overflows the array, grow it.

`(` : level start

   - set data index of current level to 0.
   
`)` : level end

   - set data index of current level to last position.
     
### Data access commands

`+` : increment data cell by 1

`-` : decrement data cell by 1

`~` : invert bits of current data cell

### Indexed Register commands
`0-9` : select indexed register

`#` : copy current data cell to current register

`%` : write current register to current data cell

`@` : apply current register value as execution multiplier of following command


### literals 

`$` : embedded data

- treat all bytes until matching `$` as embedded data. 
- data index incremented by length of embedded data.
- `\$` escape sequence for embedding "$" (ex: `$max\$imum$` for string literal `max$imum`)
- `\x` escape sequence for embedding hex values in range 0-F (ex: `$\xF\xE$` embeds byte `0xF` followed by byte `0xE`)
- `\X` escape sequence for embedding hex values in range 00-FF (ex: `$\XFE$` embeds byte value`0xFE`)
 
   
### IO commands
   
`?` : read byte

- read byte from stdin unto current data cell.
- data index incremended by 1.

`!` : write byte

- write current data cell byte to stdout. 
- data index incremented by 1.


`n` : write numeric 

- write the value of current byte per printf("%d", b) i.e. byte value 27 -> "27"

`N` : write numeric 

- write the zero-padded numeric value of current byte per printf("%03d", b) i.e. byte value 27 -> "027")    
     
### control flow
`[` : loop begin

- if current cell value is 0, move to next command after the matching ']'.

`]` : loop end

- if current cell value is not 0, move to the next command after the matching '['.

## example: hello world!
  
    $hello world!\xc$<#(@!

### example walk thorough

First we embed literal string "hello world!" and literal hex value `0xc` (length of "hello world!")

	$hello world!\xc$
	
Data index is now positioned one cell ahead of `oxc` value. So step back and copy it to the (default 0) register.

    <#
    
Move to start of data array ...

    (
    
… and apply the value in (default 0) register to the following (output `!`) command. This prints 12 bytes starting from data index 0, i.e. "hello world!"

    @!