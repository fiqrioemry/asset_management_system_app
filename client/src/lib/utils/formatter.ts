export function formatDate(dateString: any) {
	return new Date(dateString).toLocaleDateString('id-ID', {
		year: 'numeric',
		month: 'long',
		day: 'numeric'
	});
}

export function getAvatarInitials(fullname: string) {
	if (!fullname) return '';
	return fullname
		.split(' ')
		.map((name) => name.charAt(0))
		.join('')
		.toUpperCase()
		.substring(0, 2);
}

export function formatPrice(price: number): string {
	return new Intl.NumberFormat('id-ID', {
		style: 'currency',
		currency: 'IDR'
	}).format(price);
}
