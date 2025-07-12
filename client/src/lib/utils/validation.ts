export interface ValidationRule {
	// Existing rules
	required?: boolean;
	minLength?: number;
	maxLength?: number;
	email?: boolean;
	pattern?: RegExp;
	message?: string;

	// New numeric rules
	min?: number; // Minimum value for numbers
	max?: number; // Maximum value for numbers
	integer?: boolean; // Must be integer (no decimals)
	positive?: boolean; // Must be positive (> 0)
	nonNegative?: boolean; // Must be >= 0
	decimal?: number; // Max decimal places
	currency?: boolean; // Currency format validation
}

export interface ValidationRules {
	[key: string]: ValidationRule;
}

export interface ValidationErrors {
	[key: string]: string;
}

export function validateForm(data: Record<string, any>, rules: ValidationRules): ValidationErrors {
	const errors: ValidationErrors = {};

	for (const [field, rule] of Object.entries(rules)) {
		const value = data[field];

		// Required validation
		if (rule.required && (!value || value.toString().trim() === '')) {
			errors[field] = rule.message || `${field} is required`;
			continue;
		}

		// Skip other validations if value is empty and not required
		if (!value || value.toString().trim() === '') {
			continue;
		}

		// String length validations
		if (rule.minLength && value.length < rule.minLength) {
			errors[field] = rule.message || `${field} must be at least ${rule.minLength} characters`;
			continue;
		}

		if (rule.maxLength && value.length > rule.maxLength) {
			errors[field] = rule.message || `${field} must be less than ${rule.maxLength} characters`;
			continue;
		}

		// Email validation
		if (rule.email && !isValidEmail(value)) {
			errors[field] = rule.message || 'Please enter a valid email address';
			continue;
		}

		// Pattern validation
		if (rule.pattern && !rule.pattern.test(value)) {
			errors[field] = rule.message || `${field} format is invalid`;
			continue;
		}

		// Numeric validations
		const numValue = parseFloat(value);

		// Check if it's a valid number for numeric rules
		if (
			(rule.min !== undefined ||
				rule.max !== undefined ||
				rule.positive ||
				rule.nonNegative ||
				rule.integer ||
				rule.decimal !== undefined ||
				rule.currency) &&
			isNaN(numValue)
		) {
			errors[field] = rule.message || `${field} must be a valid number`;
			continue;
		}

		// Min value validation
		if (rule.min !== undefined && numValue < rule.min) {
			errors[field] = rule.message || `${field} must be at least ${rule.min}`;
			continue;
		}

		// Max value validation
		if (rule.max !== undefined && numValue > rule.max) {
			errors[field] = rule.message || `${field} must be no more than ${rule.max}`;
			continue;
		}

		// Positive validation (> 0)
		if (rule.positive && numValue <= 0) {
			errors[field] = rule.message || `${field} must be greater than 0`;
			continue;
		}

		// Non-negative validation (>= 0)
		if (rule.nonNegative && numValue < 0) {
			errors[field] = rule.message || `${field} must be 0 or greater`;
			continue;
		}

		// Integer validation
		if (rule.integer && !Number.isInteger(numValue)) {
			errors[field] = rule.message || `${field} must be a whole number`;
			continue;
		}

		// Decimal places validation
		if (rule.decimal !== undefined && !isValidDecimalPlaces(value.toString(), rule.decimal)) {
			errors[field] =
				rule.message || `${field} must have no more than ${rule.decimal} decimal places`;
			continue;
		}

		// Currency validation
		if (rule.currency && !isValidCurrency(value.toString())) {
			errors[field] = rule.message || `${field} must be a valid currency amount`;
			continue;
		}
	}

	return errors;
}

// Helper functions
function isValidEmail(email: string): boolean {
	const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
	return emailRegex.test(email);
}

function isValidDecimalPlaces(value: string, maxDecimals: number): boolean {
	const decimalIndex = value.indexOf('.');
	if (decimalIndex === -1) return true; // No decimals

	const decimalPlaces = value.length - decimalIndex - 1;
	return decimalPlaces <= maxDecimals;
}

function isValidCurrency(value: string): boolean {
	const currencyRegex = /^\d+(\.\d{1,2})?$/;
	return currencyRegex.test(value.replace(/,/g, '')); // Remove commas first
}

export const ValidationPresets = {
	// Price validations
	price: {
		required: true,
		positive: true,
		decimal: 2,
		message: 'Price must be greater than 0'
	} as ValidationRule,

	priceOptional: {
		positive: true,
		decimal: 2,
		message: 'Price must be greater than 0'
	} as ValidationRule,

	currency: {
		required: true,
		currency: true,
		positive: true,
		message: 'Please enter a valid price amount'
	} as ValidationRule,

	// Quantity validations
	quantity: {
		required: true,
		integer: true,
		positive: true,
		message: 'Quantity must be a positive whole number'
	} as ValidationRule,

	stock: {
		required: true,
		integer: true,
		nonNegative: true,
		message: 'Stock must be 0 or greater'
	} as ValidationRule,

	// Common field validations
	name: {
		required: true,
		minLength: 2,
		maxLength: 100,
		message: 'Name must be between 2-100 characters'
	} as ValidationRule,

	description: {
		maxLength: 500,
		message: 'Description must be less than 500 characters'
	} as ValidationRule,

	email: {
		required: true,
		email: true,
		message: 'Please enter a valid email address'
	} as ValidationRule,

	// Age/Number ranges
	age: {
		required: true,
		integer: true,
		min: 1,
		max: 150,
		message: 'Age must be between 1-150'
	} as ValidationRule,

	percentage: {
		min: 0,
		max: 100,
		decimal: 2,
		message: 'Percentage must be between 0-100'
	} as ValidationRule
};
