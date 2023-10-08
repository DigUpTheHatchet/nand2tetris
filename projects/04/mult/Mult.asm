// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
//
// This program only needs to handle arguments that satisfy
// R0 >= 0, R1 >= 0, and R0*R1 < 32768.


// Approach is simple:
// Add R0 to itself R1 times. 

// Pseudo Code:
// sum := 0
// i := 0
// if R0 == 0 goto STOP
// LOOP:
//   if i == R1 goto STOP
//   sum = sum + R0
//   i = i+1
//   goto LOOP
// STOP:
//   R2 = sum

// ASM Code

  // sum = 0
  @sum
  M=0

  // i = 0
  @i
  M=0

  // if R0 == 0 goto STOP, result will be 0
  @R0
  D=M
  @STOP
  D;JEQ

(LOOP)
  // if i==R1 goto STOP, result is current 'sum' value
  @i
  D=M
  @R1
  D=D-M
  @STOP
  D;JEQ

  // sum += R0 (Add R0 to itself)
  @R0
  D=M

  @sum
  M=M+D

  // i += 1
  @i
  M=M+1

  // goto LOOP
  @LOOP
  0;JMP

(STOP)
  // R2 = sum (store the result in R2)
  @sum
  D=M
  @R2
  M=D

(END)
  // Infinite "End" Loop
  @END
  0;JMP