// src/lib/api/client.ts
import { goto } from '$app/navigation';
import { authStore } from '$lib/stores/auth.store';
import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios';

const baseConfig = {
	baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:5005/api/v1',
	withCredentials: true,
	timeout: 10000,
	headers: {
		'Content-Type': 'application/json'
	}
};

class PublicApiClient {
	private client: AxiosInstance;

	constructor() {
		this.client = axios.create(baseConfig);
	}

	async get<T = any>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> {
		return this.client.get(url, config);
	}

	async post<T = any>(
		url: string,
		data?: any,
		config?: AxiosRequestConfig
	): Promise<AxiosResponse<T>> {
		return this.client.post(url, data, config);
	}

	async put<T = any>(
		url: string,
		data?: any,
		config?: AxiosRequestConfig
	): Promise<AxiosResponse<T>> {
		return this.client.put(url, data, config);
	}

	async delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> {
		return this.client.delete(url, config);
	}
}

class ProtectedApiClient {
	private client: AxiosInstance;
	private isRefreshing = false;
	private refreshSubscribers: ((token: string) => void)[] = [];

	constructor() {
		this.client = axios.create(baseConfig);
		this.setupInterceptors();
	}

	private setupInterceptors() {
		// Request interceptor
		console.log('Check step process: 1');
		this.client.interceptors.request.use(
			(config) => {
				return config;
			},
			(error) => {
				return Promise.reject(error);
			}
		);
		console.log('Check step process: 2');

		// Response interceptor dengan automatic token refresh
		this.client.interceptors.response.use(
			(response) => response,
			async (error) => {
				const originalRequest = error.config;
				console.log('Check step process:3');
				if (error.response?.status === 401 && !originalRequest._retry) {
					console.log('check step process: 3.1');
					if (this.isRefreshing) {
						return new Promise((resolve) => {
							this.refreshSubscribers.push(() => {
								resolve(this.client(originalRequest));
							});
						});
					}

					originalRequest._retry = true;
					this.isRefreshing = true;

					try {
						console.log('Check step process: 4');
						await this.refreshToken();
						this.onRefreshed();
						return this.client(originalRequest);
					} catch (refreshError) {
						console.log('Check step process: 5');
						this.onRefreshFailed();
						return Promise.reject(refreshError);
					} finally {
						console.log('Check step process: 6');
						this.isRefreshing = false;
					}
				}

				return Promise.reject(error);
			}
		);
	}

	private async refreshToken() {
		const response = await publicApiClient.post('/auth/refresh-token');
		if (response.data.success) {
			authStore.setUser(response.data.user);
			return response.data;
		}
		throw new Error('Token refresh failed');
	}

	private onRefreshed() {
		this.refreshSubscribers.forEach((callback) => callback(''));
		this.refreshSubscribers = [];
	}

	private onRefreshFailed() {
		this.refreshSubscribers = [];
		authStore.logout();
		goto('/signin');
	}

	// Public methods
	async get<T = any>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> {
		return this.client.get(url, config);
	}

	async post<T = any>(
		url: string,
		data?: any,
		config?: AxiosRequestConfig
	): Promise<AxiosResponse<T>> {
		return this.client.post(url, data, config);
	}

	async put<T = any>(
		url: string,
		data?: any,
		config?: AxiosRequestConfig
	): Promise<AxiosResponse<T>> {
		return this.client.put(url, data, config);
	}

	async delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> {
		return this.client.delete(url, config);
	}
}

// Export instances
export const publicApiClient = new PublicApiClient();
export const protectedApiClient = new ProtectedApiClient();

// Untuk backward compatibility
export const apiClient = protectedApiClient;
