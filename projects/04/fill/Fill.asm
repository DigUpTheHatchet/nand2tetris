// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Pseudocode
// screensize = 8192 
// LOOP:
//  color = 0 (white)
//  i = 0
//  if RAM[KBD] != 0 set color = -1 (black) (KBD is 0 when no key pressed)
//  PAINT:
//   if i == screensize goto LOOP
//   RAM[SCREEN + i] = color (SCREEN is the address of top leftmost corner)
//   i += 1


// Set screensize = 8192 (256x32 16 bit words)
@8192
D=A
@screensize
M=D

(LOOP)
// Reset: color = white, i = 0
@color
M=0
@i
M=0

// if RAM[KBD] != 0 set color = -1
@KBD
D=M
@DRAW
D;JEQ
@color
M=-1

(DRAW)
// if i==screensize goto LOOP (finished drawing loop)
@i
D=M
@screensize
D=D-M
@LOOP
D;JEQ

// Find the address of the nextword *(SCREEN + i)
@i
D=M
@SCREEN
D=D+A
@nextword
M=D

// Retrieve the current color value
@color
D=M

// Set RAM[nextword] = color (0 or -1)
@nextword
A=M
M=D

// i=i+1
@i
M=M+1

// Continue DRAW loop
@DRAW
0;JMP

