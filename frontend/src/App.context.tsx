import {createContext, useContext, useEffect, useState} from 'react';
import {useNavigate} from 'react-router';
import {fetchBackendStatus, fetchWalletAnalysis} from './api';
import {WalletAnalysis} from './types/walletAnalysis';

const DEFAULT_WALLET_ADDRESS = '0xcB1C1FdE09f811B294172696404e88E658659905';

type AppContextType = {
    walletAddress: string;
    isBackendAlive: boolean;
    isLoadingWalletAnalysis: boolean;
    walletAnalysis: WalletAnalysis | null;
    lastWalletsList: string[];
    error: string | null;
    setWalletAddress: (address: string) => void;
    searchWallet: (address: string) => void;
};

const AppContext = createContext<AppContextType>({} as AppContextType);

export const useAppContext = () => {
    return useContext(AppContext);
};

export const AppContextProvider = ({children}: {children: React.ReactNode}) => {
    const [walletAddress, setWalletAddress] = useState(DEFAULT_WALLET_ADDRESS);
    const [isBackendAlive, setIsBackendAlive] = useState(false);
    const [isLoadingWalletAnalysis, setIsLoadingWalletAnalysis] = useState(false);
    const [walletAnalysis, setWalletAnalysis] = useState<WalletAnalysis | null>(null);
    const [lastWalletsList, setLastWalletsList] = useState<string[]>([]);
    const [error, setError] = useState<string | null>(null);

    const navigate = useNavigate();

    const searchWallet = async (address: string) => {
        setIsLoadingWalletAnalysis(true);
        fetchWalletAnalysis(address)
            .then((data) => {
                setWalletAnalysis(data);
                setIsLoadingWalletAnalysis(false);
                setLastWalletsList([...lastWalletsList, address]);
                navigate('/analyze');
            })
            .catch((err) => {
                setIsLoadingWalletAnalysis(false);
                setError(err);
            });
    };

    useEffect(() => {
        fetchBackendStatus().then(setIsBackendAlive);
    }, []);

    return (
        <AppContext.Provider
            value={{
                walletAddress,
                setWalletAddress,
                searchWallet,
                isBackendAlive,
                isLoadingWalletAnalysis,
                walletAnalysis,
                lastWalletsList,
                error,
            }}
        >
            {children}
        </AppContext.Provider>
    );
};
