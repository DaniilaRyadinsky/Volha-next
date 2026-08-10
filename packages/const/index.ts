declare const process: {
	env: {
		NEXT_PUBLIC_BASE_URL?: string;
	};
};

export const BASE_URL = process.env.NEXT_PUBLIC_BASE_URL;
export const API_URL = `${process.env.NEXT_PUBLIC_BASE_URL}/api`;