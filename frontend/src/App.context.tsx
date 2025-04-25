import {createContext, useContext, useEffect, useState} from 'react';
import {useNavigate} from 'react-router';
import {fetchBackendStatus} from './api';

const DEFAULT_WALLET_ADDRESS = '0x0000000000000000000000000000000000000000';

type AppContextType = {
    walletAddress: string;
    walletType: string;
    isBackendAlive: boolean;
    setWalletAddress: (address: string) => void;
    setWalletType: (type: string) => void;
    searchWallet: () => void;
};

const AppContext = createContext<AppContextType>({} as AppContextType);

export const useAppContext = () => {
    return useContext(AppContext);
};

export const AppContextProvider = ({children}: {children: React.ReactNode}) => {
    const [walletAddress, setWalletAddress] = useState(DEFAULT_WALLET_ADDRESS);
    const [walletType, setWalletType] = useState('Unknown');
    const [isBackendAlive, setIsBackendAlive] = useState(false);
    const navigate = useNavigate();

    const searchWallet = async () => {
        navigate('/analyze');
    };

    useEffect(() => {
        fetchBackendStatus().then(setIsBackendAlive);
    }, []);

    return (
        <AppContext.Provider
            value={{
                walletAddress,
                setWalletAddress,
                walletType,
                setWalletType,
                searchWallet,
                isBackendAlive,
            }}
        >
            {children}
        </AppContext.Provider>
    );
};
