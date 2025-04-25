import {API_URL} from './const';

const fetchBackendStatus = async () => {
    try {
        const response = await fetch(`https://${API_URL}`);
        return response.ok;
    } catch (error) {
        return false;
    }
};

export {fetchBackendStatus};
