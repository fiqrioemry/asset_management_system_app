// ========================================
// VALIDATION RULES FROM GO DTOs
// File: src/lib/utils/dto-validation.ts
// ========================================

import { type ValidationRules } from './validation';

// ========================================
// 1. AUTH VALIDATION RULES
// ========================================

// Login Request
export const loginValidationRules: ValidationRules = {
	email: {
		required: true,
		email: true,
		message: 'Please enter a valid email address'
	},
	password: {
		required: true,
		minLength: 6,
		message: 'Password must be at least 6 characters'
	}
};

// Register Request
export const registerValidationRules: ValidationRules = {
	fullname: {
		required: true,
		minLength: 2,
		maxLength: 100,
		message: 'Full name must be between 2-100 characters'
	},
	email: {
		required: true,
		email: true,
		message: 'Please enter a valid email address'
	},
	password: {
		required: true,
		minLength: 6,
		message: 'Password must be at least 6 characters'
	}
};

// Update User Request
export const updateUserValidationRules: ValidationRules = {
	fullname: {
		required: true,
		minLength: 2,
		maxLength: 100,
		message: 'Full name must be between 2-100 characters'
	}
	// avatar: handled separately (file upload)
	// avatarUrl: optional string, no validation needed
};

// Change Password Request
export const changePasswordValidationRules: ValidationRules = {
	currentPassword: {
		required: true,
		minLength: 6,
		message: 'Current password must be at least 6 characters'
	},
	newPassword: {
		required: true,
		minLength: 6,
		message: 'New password must be at least 6 characters'
	},
	confirmPassword: {
		required: true,
		minLength: 6,
		message: 'Confirm password must be at least 6 characters'
	}
};

// Forgot Password Request
export const forgotPasswordValidationRules: ValidationRules = {
	email: {
		required: true,
		email: true,
		message: 'Please enter a valid email address'
	}
};

// Reset Password Request
export const resetPasswordValidationRules: ValidationRules = {
	token: {
		required: true,
		message: 'Reset token is required'
	},
	newPassword: {
		required: true,
		minLength: 6,
		message: 'New password must be at least 6 characters'
	},
	confirmPassword: {
		required: true,
		minLength: 6,
		message: 'Confirm password must be at least 6 characters'
	}
};

// ========================================
// 2. CATEGORY VALIDATION RULES
// ========================================

// Create Category Request
export const createCategoryValidationRules: ValidationRules = {
	name: {
		required: true,
		minLength: 2,
		maxLength: 100,
		message: 'Category name must be between 2-100 characters'
	}
	// parentId: optional UUID, handled separately if needed
};

// Update Category Request
export const updateCategoryValidationRules: ValidationRules = {
	name: {
		required: true,
		minLength: 2,
		maxLength: 100,
		message: 'Category name must be between 2-100 characters'
	}
	// parentId: optional UUID, handled separately if needed
};

// ========================================
// 3. LOCATION VALIDATION RULES
// ========================================

// Create Location Request
export const createLocationValidationRules: ValidationRules = {
	name: {
		required: true,
		minLength: 2,
		maxLength: 100,
		message: 'Location name must be between 2-100 characters'
	}
};

// Update Location Request
export const updateLocationValidationRules: ValidationRules = {
	name: {
		required: true,
		minLength: 2,
		maxLength: 100,
		message: 'Location name must be between 2-100 characters'
	}
};

// ========================================
// 4. ASSET VALIDATION RULES
// ========================================

// Create Asset Request
export const createAssetValidationRules: ValidationRules = {
	name: {
		required: true,
		minLength: 1,
		maxLength: 100,
		message: 'Asset name must be between 1-100 characters'
	},
	description: {
		maxLength: 255,
		message: 'Description must be less than 255 characters'
	},
	locationId: {
		required: true,
		message: 'Please select a location'
	},
	categoryId: {
		required: true,
		message: 'Please select a category'
	},
	price: {
		required: true,
		nonNegative: true, // Based on Go binding:"required,min=0"
		decimal: 2,
		message: 'Price must be 0 or greater'
	},
	condition: {
		required: true,
		pattern: /^(new|good|fair|poor)$/,
		message: 'Condition must be one of: new, good, fair, poor'
	},
	serialNumber: {
		maxLength: 100,
		message: 'Serial number must be less than 100 characters'
	}
	// image: handled separately (file upload)
	// purchaseDate: optional date, handled separately
	// warranty: optional date, handled separately
};

// Update Asset Request
export const updateAssetValidationRules: ValidationRules = {
	name: {
		minLength: 1,
		maxLength: 100,
		message: 'Asset name must be between 1-100 characters'
	},
	description: {
		maxLength: 255,
		message: 'Description must be less than 255 characters'
	},
	// locationId: optional in update
	// categoryId: optional in update
	price: {
		nonNegative: true, // Based on Go binding:"omitempty,min=0"
		decimal: 2,
		message: 'Price must be 0 or greater'
	},
	condition: {
		pattern: /^(new|good|fair|poor)$/,
		message: 'Condition must be one of: new, good, fair, poor'
	},
	serialNumber: {
		maxLength: 100,
		message: 'Serial number must be less than 100 characters'
	}
	// Other fields are optional in update
};

// Get Assets Request (for filters)
export const getAssetsFiltersValidationRules: ValidationRules = {
	page: {
		positive: true,
		integer: true,
		message: 'Page must be a positive integer'
	},
	limit: {
		positive: true,
		integer: true,
		min: 1,
		max: 100,
		message: 'Limit must be between 1-100'
	},
	search: {
		maxLength: 100,
		message: 'Search term must be less than 100 characters'
	},
	condition: {
		pattern: /^(new|good|fair|poor)$/,
		message: 'Condition must be one of: new, good, fair, poor'
	},
	minPrice: {
		nonNegative: true,
		decimal: 2,
		message: 'Minimum price must be 0 or greater'
	},
	maxPrice: {
		nonNegative: true,
		decimal: 2,
		message: 'Maximum price must be 0 or greater'
	},
	sortBy: {
		pattern: /^(name|price|createdAt|purchaseDate)$/,
		message: 'Sort by must be one of: name, price, createdAt, purchaseDate'
	},
	sortOrder: {
		pattern: /^(asc|desc)$/,
		message: 'Sort order must be asc or desc'
	}
};

// ========================================
// 5. FORM-SPECIFIC VALIDATION GROUPS
// ========================================

// Complete Asset Form Validation (Create + Update combined)
export const assetFormValidationRules: ValidationRules = {
	...createAssetValidationRules,
	// Override price for better UX (required for create, optional for update)
	price: {
		required: true,
		positive: true, // Changed to positive since assets should have value > 0
		decimal: 2,
		message: 'Price must be greater than 0'
	}
};

// User Profile Form Validation
export const userProfileValidationRules: ValidationRules = {
	...updateUserValidationRules
};

// Category Form Validation
export const categoryFormValidationRules: ValidationRules = {
	...createCategoryValidationRules
};

// Location Form Validation
export const locationFormValidationRules: ValidationRules = {
	...createLocationValidationRules
};
