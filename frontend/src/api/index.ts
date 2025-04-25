import {WalletAnalysis} from '@/types/walletAnalysis';
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

const fetchWalletAnalysis = (address: string): Promise<WalletAnalysis> => {
    return new Promise((resolve, reject) => {
        fetch(`${API_URL}/analyze`, {
            method: 'POST',
            body: JSON.stringify({address}),
        })
            .then((response) => {
                if (response.ok) {
                    response.json().then(resolve);
                } else {
                    response.text().then(reject);
                }
            })
            .catch(reject);
    });
};

export {fetchBackendStatus, fetchWalletAnalysis};
