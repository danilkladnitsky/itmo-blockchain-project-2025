import {API_URL} from './const';

const fetchBackendStatus = async () => {
    try {
        const response = await fetch(`${API_URL}/alive`, {
            method: 'GET',
            mode: 'cors',
            credentials: 'include',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
        });
        return response.ok;
    } catch (error) {
        return false;
    }
};

const fetchWalletAnalysis = async (address: string) => {
    try {
        const response = await fetch(`${API_URL}/analyze`, {
            method: 'POST',
            body: JSON.stringify({address}),
            mode: 'cors',
            credentials: 'include',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
        });
        return response.json();
    } catch (error) {
        return null;
    }
};
export {fetchBackendStatus, fetchWalletAnalysis};
