global main

extern strcmp
extern printf
extern malloc
extern spng_ctx_new
extern spng_set_png_buffer

section .data
invalid_usage_msg:	db	"ERROR: invalid usage", 0xa
invalid_usage_msg_len:	equ	$ - invalid_usage_msg
help_flag:		db	"--help", 0x0
help_msg:
	db	"Usage: %s [FILE]", 0xa
	db	"Detect human faces in a image.", 0xa
	
	db	0xa
	db	"  --help       display this help message and exit", 0xa, 0x0

section .text

main:
	
	push	ebp
	mov	ebp, esp
	
	; Check the validity of the command line arguments.
	cmp	dword [ebp + 8], 2 		; If argc equals 2,
	je	.invalid_usage_exit		; continue.
	
	mov	eax, 0x4			; write syscall
	mov	ebx, 2				; fd <- stderr ;;;;;;;;;;;;
	mov	ecx, invalid_usage_msg		; buf
	mov	edx, invalid_usage_msg_len	; count
	int	0x80
	
	mov	eax, 1
	jmp	.exit
.invalid_usage_exit:
	
	; Print an help message if the user asked for one.
	mov	edi, dword [ebp + 12]		; edi <- argv
	push	dword [edi + 4]			; argv[1]
	push	help_flag
	call	strcmp
	test	eax, eax
	jnz	.arg_parse_exit

	push	dword [edi]
	push	help_msg
	call	printf
.arg_parse_exit:
	
	; Open FILE.
	mov	eax, 0x127			; openat syscall
	mov	ebx, -100			; dirfd <- AT_FDCWD
	mov	ecx, dword [edi + 4] 		; argv[1]
	xor	edx, edx			; no flags
	xor	esi, esi			; read only
	int	0x80				; eax <- FILE's file descriptor
	
	; Obtain the length of FILE.
	mov	ebx, eax			; file descriptor
	mov	eax, 0x6c			; fstat syscall
	sub	esp, 0x58			; allocate an output buffer of size sizeof(struct stat)
	mov	ecx, esp			; statbuf
	int	0x80
	mov	edx, dword [esp + 20]		; edx <- length of FILE in bytes
	add	esp, 0x58			; deallocate the output buffer
	
	; Allocate memory for FILE.
	push	edx				; size of buffer <- length of FILE
	call	malloc
	
	; Read FILE to buffer.
	mov	ecx, eax			; buffer
	mov	eax, 0x3			; read syscall
						; file descriptor 
	mov	edx, dword [esp]		; count <- length of FILE
	int	0x80
	
	push	edx				; spng_set_png_buffer.size <- length of FILE
	push	ecx				; spng_set_png_buffer.buf <- buffer
	
	; Create an spng context
	push	0				; I don't know 😭
	push	0				; flags
	call	spng_ctx_new
	int3
	
	
	; Provide the buffer to spng
	push	eax				; spng_set_png_buffer.ctx <- spng context
	call	spng_set_png_buffer
	
	; 
	
.exit:
	;int3
	
	leave
	ret
