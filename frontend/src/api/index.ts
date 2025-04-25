import {API_URL, ML_SERVICE_URL} from './const';

const fetchBackendStatus = async () => {
    try {
        const response = await fetch(`${API_URL}/alive`);
        return response.ok;
    } catch (error) {
        return false;
    }
};

const fetchMlServiceStatus = async () => {
    try {
        const response = await fetch(`${ML_SERVICE_URL}/health`);
        return response.ok;
    } catch (error) {
        return false;
    }
};
export {fetchBackendStatus, fetchMlServiceStatus};
