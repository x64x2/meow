global meow_haar_x2
global meow_haar_y2
global meow_haar_x3
global meow_haar_y3
global meow_haar_x2y2

extern meow_sum_area

section .text

meow_haar_x2:
	; Preserve registers.
	push	ebp
	mov	ebp, esp
	push	ebx
	push	esi
	
	; Calculate the sum of the left area.
	mov	eax, dword [ebp + 24]		; eax <- y0
	mov	edx, dword [ebp + 32]		; edx <- height
	add	edx, eax			; edx <- y0 + height
	push	edx				; y1	

	mov     ebx, dword [ebp + 28]		; ebx <- width
	shr	ebx, 1				; ebx <- width / 2
	mov	esi, dword [ebp + 20]		; esi <- x0
	add	ebx, esi			; ebx <- x0 + width / 2
	push	ebx				; x1

	push	eax				; y0
	push	esi				; x0

	push	dword [esp + 40]		; img.values
	push	dword [esp + 40]		; img.height
	push	dword [esp + 40]		; img.width

	call	meow_sum_area
	
	; Calculate the sum of the right area.
	mov	[esp + 12], ebx			; x0
	add	esi, dword [ebp + 28]		; esi <- x0 + width
	mov	[esp + 20], esi			; x1
	
	mov	ebx, eax			; Preserve the left sum.
	
	call	meow_sum_area
	
	; Subtract the right sum from the left sum.
	sub	eax, ebx
	neg	eax

	pop	esi
	pop	ebx
	leave
	ret

meow_haar_y2:
	; Preserve registers.
	push	ebp
	mov	ebp, esp
	push	ebx
	push	esi
	
	; Calculate the sum of the top area.	
	mov	ebx, dword [ebp + 32]		; ebx <- height
	shr	ebx, 1				; ebx <- height / 2
	mov	esi, dword [ebp + 24]		; esi <- y0
	add	ebx, esi			; ebx <- y0 + height / 2
	push	ebx				; y1
	
	mov	eax, dword [ebp + 20]		; eax <- x0
	mov	edx, dword [ebp + 28]		; edx <- width
	add	edx, eax			; edx <- x0 + width
	push	edx				; x1
	
	push	esi				; y0
	push	eax				; x0
	
	push	dword [esp + 40]		; img.values
	push	dword [esp + 40]		; img.height
	push	dword [esp + 40]		; img.width
	
	call	meow_sum_area
	
	; Calculate the sum of the bottom area.
	mov	[esp + 16], ebx			; y0
	add	esi, dword [ebp + 32]		; esi <- y0 + height
	mov	[esp + 24], esi			; y1
	
	mov	ebx, eax			; Preserve the top sum.
	
	call	meow_sum_area
	
	; Subtract the top sum from the bottom sum.
	sub	eax, ebx
	neg	eax
	
	pop	esi
	pop	ebx
	leave
	ret

meow_haar_x3:
	; Preserve registers.
	push	ebp
	mov	ebp, esp
	push	ebx
	push	esi
	push	edi
	
	; Calculate the sum of the left area.
	mov	ecx, dword [ebp + 24]		; ecx <- y0
	mov	eax, dword [ebp + 32]		; eax <- height
	add	eax, ecx			; eax <- y0 + height
	push	eax				; y1	

	mov     eax, 0xAAAAAAAB			; I generated this using AMD's udiv.exe (0.AAAA... = ⅔)
	mul	dword [ebp + 28]		; edx:eax <- width * ⅔
	shr	edx, 1				; edx <- width / 3
	mov	esi, edx			; esi <- width / 3
	
	mov	ebx, dword [ebp + 20]		; ebx <- x0
	mov	eax, ebx			; eax <- x0
	add	ebx, edx			; ebx <- x0 + width / 3
	push	ebx				; x1
	
	push	ecx				; y0
	push	eax				; x0
	
	push	dword [esp + 44]		; img.values
	push	dword [esp + 44]		; img.height
	push	dword [esp + 44]		; img.width
	
	call	meow_sum_area
	mov	edi, eax			; edi <- left_value
	neg	edi				; edi <- -left_value
	
	; Calculate the sum of the middle area.
	mov	[esp + 12], ebx			; x0 <- x1
	add	ebx, esi			; x1 <- x1 + width / 3
	mov	[esp + 20], ebx			; x1
	call	meow_sum_area
	add	edi, eax			; edi <- -left_value + middle_value
	
	; Calculate the sum of the right area.
	mov	[esp + 12], ebx			; x0 <- x1
	mov	eax, dword [ebp + 20]		; eax <- x0
	add	eax, dword [ebp + 28]		; eax <- x0 + width
	mov	[esp + 20], eax			; x1
	call	meow_sum_area
	sub	eax, edi			; eax <- right_value - (middle_value - left_value)

	pop	edi
	pop	esi
	pop	ebx
	leave
	ret

meow_haar_y3:
	; Preserve registers.
	push	ebp
	mov	ebp, esp
	push	ebx
	push	esi
	push	edi
	
	; Calculate the sum of the top area.
	mov     eax, 0xAAAAAAAB			; I generated this using AMD's udiv.exe (0.AAAA... = ⅔)
	mul	dword [ebp + 32]		; edx:eax <- height * ⅔
	shr	edx, 1				; edx <- height / 3
	mov	esi, edx			; esi <- height / 3

	mov	ebx, dword [ebp + 24]		; ebx <- y0
	mov	edx, ebx			; edx <- y0
	add	ebx, esi			; ebx <- y0 + height / 3
	push	ebx				; y1

	mov	ecx, dword [ebp + 20]		; ecx <- x0
	mov	eax, dword [ebp + 28]		; eax <- width
	add	eax, ecx			; eax <- x0 + width
	push	eax				; x1	
	
	push	edx				; y0
	push	ecx				; x0
	
	push	dword [esp + 44]		; img.values
	push	dword [esp + 44]		; img.height
	push	dword [esp + 44]		; img.width
	
	call	meow_sum_area
	mov	edi, eax			; edi <- top_value
	neg	edi				; edi <- -top_value
	
	; Calculate the sum of the middle area.
	mov	[esp + 16], ebx			; y0 <- y1
	add	ebx, esi			; y1 <- y1 + width / 3
	mov	[esp + 24], ebx			; y1
	call	meow_sum_area
	add	edi, eax			; edi <- -top_value + middle_value
	
	; Calculate the sum of the bottom area.
	mov	[esp + 16], ebx			; y0 <- y1
	mov	eax, dword [ebp + 24]		; eax <- y0
	add	eax, dword [ebp + 32]		; eax <- y0 + height
	mov	[esp + 24], eax			; y1
	call	meow_sum_area
	sub	eax, edi			; eax <- bottom_value - (middle_value - top_value)

	pop	edi
	pop	esi
	pop	ebx
	leave
	ret

meow_haar_x2y2:
	ret