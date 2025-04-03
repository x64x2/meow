global meow_integrate_image
global meow_sum_area

extern calloc

section .text

meow_integrate_image:
	push	ebp
	mov	ebp, esp
	push	ebx
	push	esi
	push	edi
	
	; Set the width and height of the integral image.
	mov	ebx, dword [ebp + 8]		; ebx <- dst
	mov	eax, dword [ebp + 16]		; eax <- width
	lea	esi, [eax + 1]			; esi <- width + 1
	
	mov	eax, dword [ebp + 20]		; eax <- height
	add	eax, 1				; eax <- height + 1
	
	mov	dword [ebx], esi		; dst.width <- width + 1
	mov	dword [ebx + 4], eax		; dst.height <- height + 1
	
	; Allocate memory for the integral image values.
	imul	eax, esi			; eax <- dst.width * dst.height
	push	4				; Each pixel consists of 4 bytes
	push	eax				; and there are eax pixels.
	call	calloc				; Allocate and reset the memory.
	mov	dword [ebx + 8], eax
	push	eax				; spill the address of dst.values
	
	; Integrate the image.
	shl	esi, 2				; esi <- [width + 1]
	xor	ecx, ecx			; ecx = y
.y_loop_head:
	cmp	ecx, dword [ebp + 20]		; If y equals to height,
	je	.y_loop_exit			; exit the loop.
	
	mov	edi, 4				; edi = [x + 1]
.x_loop_head:
	cmp	edi, esi			; If [x + 1] equals [width + 1],
	je	.x_loop_exit			; exit the inner loop.
	
	push	esi				; spill [width + 1]
	
	mov	ebx, ecx			; ebx <- y
	imul	ebx, esi			; ebx <- [y * (width + 1)]
	shl	ecx, 2				; ecx <- [y]
	sub	ebx, ecx			; ebx <- [y * width]
	lea	ebx, [ebx + edi - 4]		; ebx <- [y * width + x]
	lea	esi, dword [ebx + ebx*2]	; esi <- [3 * (y * width + x)]
	mov	eax, dword [ebp + 12]		; eax <- address of values
	shr	esi, 2				; esi <- 3 * (y * width + x)
	movzx	edx, byte [eax + esi] 		; edx <- values[3 * (x, y)]

	add	ebx, ecx			; ebx <- [(x, y)]
	mov	esi, dword [esp + 4]		; esi <- address of dst.values
	sub	edx, dword [esi + ebx]		; edx <- edx - dst.values[(x, y)]

	add	edx, dword [esi + ebx + 4]	; edx <- edx + dst.values[(x + 1, y)]

	pop	eax				; eax <- [width + 1]
	xchg	esi, eax			; esi <- [width + 1], eax <- address of dst.values
	add	ebx, esi			; ebx <- [(x, y + 1)]
	add	edx, dword [eax + ebx]		; edx <- edx + dst.values[(x, y + 1)]
	
	mov	dword [eax + ebx + 4], edx	; dst.values[(x + 1, y + 1)] <- edx

	shr	ecx, 2				; ecx <- y
	
	add	edi, 4				; x bytes <- x bytes + sizeof(int)
	jmp	.x_loop_head

.x_loop_exit:
	add	ecx, 1				; y <- y + 1
	jmp	.y_loop_head
.y_loop_exit:

	add	esp, 12
	pop	edi
	pop	esi
	pop	ebx
	leave
	ret

meow_sum_area:
	; Extract parameters into registers.
	push	ebx
	push	edi
	push	esi
	mov	esi, dword [esp + 28]		; esi <- x0
	mov	ecx, dword [esp + 24]		; ecx <- address of img.values
	mov	ebx, dword [esp + 36]		; ebx <- x1
	mov	edx, dword [esp + 16]		; edx <- img.width
	mov	edi, dword [esp + 40]		; edi <- y1
	
	imul	edi, edx			; edi <- img.width * y1
	lea	eax, [edi + ebx]		; ebx <- bottom-right index
	mov	eax, dword [ecx + 4*eax]	; eax <- bottom-right value

	imul	edx, dword [esp + 32]		; edx <- img.width * y0
	add	ebx, edx			; ebx <- top-right index
	mov	ebx, dword [ecx + 4*ebx]	; ebx <- top-right value

	add	edi, esi			; edi <- bottom-left index
	add	ebx, dword [ecx + 4*edi]	; ebx <- top-right value + bottom-left value

	add	edx, esi			; edx <- top-left index 
	sub	eax, ebx			; eax <- bottom-right - (top-right + bottom-left)
	add	eax, dword [ecx + 4*edx]	; eax <- sum of the area
	
	; Return.
	pop	esi
	pop	edi
	pop	ebx
	ret
