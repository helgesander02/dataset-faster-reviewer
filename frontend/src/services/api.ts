import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080';
const api = axios.create({
    baseURL: API_BASE_URL?.endsWith('/') ? API_BASE_URL : `${API_BASE_URL}/`,
});

api.interceptors.response.use(
    (response) => response,
    (error) => {
        console.error('API Error:', error.response?.data || error.message);
        return Promise.reject(error);
    }
);

export interface JobsResponse {
    job_names: string[];
}

export interface DatasetsResponse {
    dataset_names: string[];
}

interface PageItems {
    item_dataset_name: string;
    item_image_set: unknown[];
}

interface GetAllPagesResponse {
    total_pages: number;
    pages: PageItems[];
}

interface ImageToSave {
    job: string;
    dataset: string;
    imageName: string;
    imagePath: string;
}

export interface SavePendingReviewPayload {
    images: ImageToSave[];
}


export const fetchJobs = async (): Promise<JobsResponse> => {
    try {
        const response = await api.get('/api/getJobs');
        return response.data;
        
    } catch (error) {
        console.error('Error fetching jobs:', error);
        throw error;
    }
};

export const fetchDatasets = async (job: string): Promise<DatasetsResponse> => {
    if (!job) {
        throw new Error('The "job" parameter is required to fetch datasets.');
    }
    
    try {
        const response = await api.get<GetAllPagesResponse>('/api/getAllPages', { params: { job } });
        console.log(`Fetching datasets for job "${job}"`);
        console.log(response.data);
        
        const datasetNames = response.data.pages
            .map(page => page.item_dataset_name)
            .filter(name => name && name.trim() !== ''); 
        
        return {
            dataset_names: datasetNames
        };
        
    } catch (error) {
        console.error('Error fetching datasets:', error);
        throw error;
    }
};


export const fetchBase64Images = async (job: string, pageIndex: number) => {
    try {
        const response = await api.get('/api/getBase64ImageSet', { params: { job, pageIndex } });
        return {
            image_path_set: response.data.image_path,
            base64_image_set: response.data.base64_image,
        };
    } catch (error) {
        console.error('Error fetching base64 images:', error);
        throw error;
    }
};

export const fetchImages = async (job: string, pageIndex: number) => {
    try {
        const response = await api.get('/api/getImageSet', { params: { job, pageIndex } });
        console.log(response.data);
        return {
            image_name_set: response.data.image_name,
            image_path_set: response.data.image_path,
        };
    } catch (error) {
        console.error('Error fetching base64 images:', error);
        throw error;
    }
};

export const updateALLPages = async (job: string, pageSize: number) => {
    if (!job) {
        throw new Error('The "job" parameter is required to update all pages.');
    }
    
    if (!pageSize || pageSize <= 0) {
        throw new Error('The "pageSize" parameter is required and must be greater than 0.');
    }
    
    try {
        const response = await api.post('/api/setAllPages', {job, pageSize});
        
        console.log(`Updated all pages for job ${job} with pageSize ${pageSize}:`, response.data);
        return response.data;
    } catch (error) {
        console.error('Error updating all pages:', error);
        throw error;
    }
};

export const fetchALLPages = async (job: string) => {
    if (!job) {
        throw new Error('The "job" parameter is required to fetch all pages.');
    }
    try {
        
        const response = await api.get('/api/getAllPages', { params: { job } });
        console.log(`Fetched all pages for job ${job}:`, response.data);
        return response.data;
    } catch (error) {
        console.error('Error fetching all pages:', error);
        throw error;
    }
};

// FIX 2: Updated function signature to use the new interface
export const savePendingReview = async (data: SavePendingReviewPayload) => {
    try {
        console.log('Saving pending review data:', data);
        const response = await api.post('/api/savePendingReview', data.images);
        console.log('Saved pending review data:', response.data);
        return response.data;
    } catch (error) {
        console.error('Error saving pending review data:', error);
        throw error;
    }
};

export const getPendingReview = async (flatten: boolean = false) => {
    try {
        const response = await api.get('/api/getPendingReview', { params: { flatten } });
        console.log('Fetched pending review data:', response.data);
        return response.data;
    } catch (error) {
        console.error('Error fetching pending review data:', error);
        throw error;
    }
};